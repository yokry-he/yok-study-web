package database

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/yokry-he/yok-study-web/examples/go-task-api/internal/config"
)

var (
	ErrOpenPostgres           = errors.New("连接 PostgreSQL 失败")
	ErrNilContext             = errors.New("上下文不能为空")
	ErrDatabaseURLRequired    = errors.New("数据库连接地址不能为空")
	ErrInvalidMaxOpenConns    = errors.New("MaxOpenConns 必须大于零")
	ErrInvalidMaxIdleConns    = errors.New("MaxIdleConns 必须大于等于零且不超过 MaxOpenConns")
	ErrInvalidConnMaxLifetime = errors.New("ConnMaxLifetime 必须大于零")
	ErrInvalidConnMaxIdleTime = errors.New("ConnMaxIdleTime 必须大于零")
)

type dbFactory func(connectionString string) (*sql.DB, error)

// Open 创建并验证 PostgreSQL 连接池。连接失败时不返回底层错误，避免 DSN 中的凭据泄露。
func Open(ctx context.Context, cfg config.DatabaseConfig) (*sql.DB, error) {
	return openWithFactory(ctx, cfg, func(connectionString string) (*sql.DB, error) {
		return sql.Open("pgx", connectionString)
	})
}

func openWithFactory(ctx context.Context, cfg config.DatabaseConfig, factory dbFactory) (*sql.DB, error) {
	if ctx == nil {
		return nil, ErrNilContext
	}
	cfg.URL = strings.TrimSpace(cfg.URL)
	if cfg.URL == "" {
		return nil, ErrDatabaseURLRequired
	}
	if cfg.MaxOpenConns <= 0 {
		return nil, ErrInvalidMaxOpenConns
	}
	if cfg.MaxIdleConns < 0 || cfg.MaxIdleConns > cfg.MaxOpenConns {
		return nil, ErrInvalidMaxIdleConns
	}
	if cfg.ConnMaxLifetime <= 0 {
		return nil, ErrInvalidConnMaxLifetime
	}
	if cfg.ConnMaxIdleTime <= 0 {
		return nil, ErrInvalidConnMaxIdleTime
	}

	db, err := factory(cfg.URL)
	if err != nil {
		return nil, ErrOpenPostgres
	}
	if db == nil {
		return nil, ErrOpenPostgres
	}
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	db.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		if errors.Is(err, context.Canceled) {
			return nil, errors.Join(ErrOpenPostgres, context.Canceled)
		}
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, errors.Join(ErrOpenPostgres, context.DeadlineExceeded)
		}
		return nil, ErrOpenPostgres
	}

	return db, nil
}
