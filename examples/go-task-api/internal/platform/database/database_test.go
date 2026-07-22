package database

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/yokry-he/yok-study-web/examples/go-task-api/internal/config"
)

func TestOpenValidatesConfigurationBeforeConnecting(t *testing.T) {
	t.Setenv("PGHOST", "127.0.0.1")
	t.Setenv("PGPORT", "1")
	t.Setenv("PGUSER", "environment-user")

	tests := []struct {
		name   string
		mutate func(*config.DatabaseConfig)
		want   error
	}{
		{
			name: "blank URL cannot fall back to PG environment",
			mutate: func(cfg *config.DatabaseConfig) {
				cfg.URL = " \t\n "
			},
			want: ErrDatabaseURLRequired,
		},
		{
			name: "zero max open connections",
			mutate: func(cfg *config.DatabaseConfig) {
				cfg.MaxOpenConns = 0
			},
			want: ErrInvalidMaxOpenConns,
		},
		{
			name: "negative max open connections",
			mutate: func(cfg *config.DatabaseConfig) {
				cfg.MaxOpenConns = -1
			},
			want: ErrInvalidMaxOpenConns,
		},
		{
			name: "negative max idle connections",
			mutate: func(cfg *config.DatabaseConfig) {
				cfg.MaxIdleConns = -1
			},
			want: ErrInvalidMaxIdleConns,
		},
		{
			name: "max idle exceeds max open",
			mutate: func(cfg *config.DatabaseConfig) {
				cfg.MaxIdleConns = cfg.MaxOpenConns + 1
			},
			want: ErrInvalidMaxIdleConns,
		},
		{
			name: "zero connection lifetime",
			mutate: func(cfg *config.DatabaseConfig) {
				cfg.ConnMaxLifetime = 0
			},
			want: ErrInvalidConnMaxLifetime,
		},
		{
			name: "negative connection lifetime",
			mutate: func(cfg *config.DatabaseConfig) {
				cfg.ConnMaxLifetime = -time.Second
			},
			want: ErrInvalidConnMaxLifetime,
		},
		{
			name: "zero connection idle time",
			mutate: func(cfg *config.DatabaseConfig) {
				cfg.ConnMaxIdleTime = 0
			},
			want: ErrInvalidConnMaxIdleTime,
		},
		{
			name: "negative connection idle time",
			mutate: func(cfg *config.DatabaseConfig) {
				cfg.ConnMaxIdleTime = -time.Second
			},
			want: ErrInvalidConnMaxIdleTime,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := validDatabaseConfig()
			tt.mutate(&cfg)

			db, err := Open(context.Background(), cfg)
			if db != nil {
				_ = db.Close()
				t.Fatal("Open returned a database for invalid configuration")
			}
			if !errors.Is(err, tt.want) {
				t.Fatalf("Open error = %v, want errors.Is(%v)", err, tt.want)
			}
			if strings.Contains(err.Error(), cfg.URL) && strings.TrimSpace(cfg.URL) != "" {
				t.Fatalf("validation error exposes URL: %v", err)
			}
		})
	}
}

func TestOpenRejectsNilContext(t *testing.T) {
	db, err := Open(nil, validDatabaseConfig())
	if db != nil {
		_ = db.Close()
		t.Fatal("Open returned a database for nil context")
	}
	if !errors.Is(err, ErrNilContext) {
		t.Fatalf("Open error = %v, want errors.Is(ErrNilContext)", err)
	}
}

func TestOpenPreservesContextCancellationWithoutExposingDSN(t *testing.T) {
	const secret = "context-secret-password"
	cfg := validDatabaseConfig()
	cfg.URL = "postgres://app:" + secret + "@127.0.0.1:1/taskdb?sslmode=disable"

	tests := []struct {
		name    string
		context func() (context.Context, context.CancelFunc)
		want    error
	}{
		{
			name: "canceled",
			context: func() (context.Context, context.CancelFunc) {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx, func() {}
			},
			want: context.Canceled,
		},
		{
			name: "deadline exceeded",
			context: func() (context.Context, context.CancelFunc) {
				return context.WithDeadline(context.Background(), time.Now().Add(-time.Second))
			},
			want: context.DeadlineExceeded,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := tt.context()
			defer cancel()

			db, err := Open(ctx, cfg)
			if db != nil {
				_ = db.Close()
				t.Fatal("Open returned a database after context cancellation")
			}
			if !errors.Is(err, ErrOpenPostgres) || !errors.Is(err, tt.want) {
				t.Fatalf("Open error = %v, want ErrOpenPostgres and %v", err, tt.want)
			}
			if strings.Contains(err.Error(), secret) || strings.Contains(err.Error(), cfg.URL) {
				t.Fatalf("Open error exposes connection details: %v", err)
			}
		})
	}
}

func TestOpenClosesPoolAfterPingFailure(t *testing.T) {
	conn := &failingPingConn{}
	factory := func(string) (*sql.DB, error) {
		return sql.OpenDB(failingPingConnector{conn: conn}), nil
	}

	db, err := openWithFactory(context.Background(), validDatabaseConfig(), factory)
	if db != nil {
		_ = db.Close()
		t.Fatal("openWithFactory returned a database after Ping failure")
	}
	if !errors.Is(err, ErrOpenPostgres) {
		t.Fatalf("openWithFactory error = %v, want errors.Is(ErrOpenPostgres)", err)
	}
	if !conn.closed.Load() {
		t.Fatal("database connection was not closed after Ping failure")
	}
}

func validDatabaseConfig() config.DatabaseConfig {
	return config.DatabaseConfig{
		URL:             "postgres://app:secret@127.0.0.1:1/taskdb?sslmode=disable",
		MaxOpenConns:    4,
		MaxIdleConns:    2,
		ConnMaxLifetime: 5 * time.Minute,
		ConnMaxIdleTime: time.Minute,
	}
}

type failingPingConnector struct {
	conn *failingPingConn
}

func (c failingPingConnector) Connect(context.Context) (driver.Conn, error) {
	return c.conn, nil
}

func (c failingPingConnector) Driver() driver.Driver {
	return failingPingDriver{}
}

type failingPingDriver struct{}

func (failingPingDriver) Open(string) (driver.Conn, error) {
	return nil, errors.New("测试驱动不支持 Open")
}

type failingPingConn struct {
	closed atomic.Bool
}

func (*failingPingConn) Prepare(string) (driver.Stmt, error) {
	return nil, errors.New("测试连接不支持 Prepare")
}

func (c *failingPingConn) Close() error {
	c.closed.Store(true)
	return nil
}

func (*failingPingConn) Begin() (driver.Tx, error) {
	return nil, errors.New("测试连接不支持 Begin")
}

func (*failingPingConn) Ping(context.Context) error {
	return errors.New("包含不应公开内容的底层 Ping 错误")
}
