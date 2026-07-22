package database

import (
	"context"
	"database/sql"
	"errors"
	"sync"

	"github.com/golang-migrate/migrate/v4"
	migratepostgres "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/yokry-he/yok-study-web/examples/go-task-api/migrations"
)

var (
	ErrNilDatabase    = errors.New("数据库连接不能为空")
	ErrMigratorClosed = errors.New("数据库迁移器已关闭")
)

// Migrator 负责执行嵌入在程序中的 PostgreSQL 迁移。
type Migrator struct {
	migrate   *migrate.Migrate
	mu        sync.Mutex
	closeOnce sync.Once
	closeErr  error
	closed    bool
}

func NewMigrator(db *sql.DB) (*Migrator, error) {
	return NewMigratorContext(context.Background(), db)
}

func NewMigratorContext(ctx context.Context, db *sql.DB) (*Migrator, error) {
	if ctx == nil {
		return nil, ErrNilContext
	}
	if db == nil {
		return nil, ErrNilDatabase
	}

	sourceDriver, err := iofs.New(migrations.Files, ".")
	if err != nil {
		return nil, err
	}

	conn, err := db.Conn(ctx)
	if err != nil {
		return nil, errors.Join(err, sourceDriver.Close())
	}
	databaseDriver, err := migratepostgres.WithConnection(ctx, conn, &migratepostgres.Config{})
	if err != nil {
		return nil, errors.Join(err, conn.Close(), sourceDriver.Close())
	}

	engine, err := migrate.NewWithInstance("iofs", sourceDriver, "postgres", databaseDriver)
	if err != nil {
		return nil, errors.Join(err, databaseDriver.Close(), sourceDriver.Close())
	}

	return &Migrator{migrate: engine}, nil
}

// Up 执行所有尚未应用的迁移；没有新迁移也视为成功。
func (m *Migrator) Up() error {
	if m == nil {
		return ErrMigratorClosed
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.closed {
		return ErrMigratorClosed
	}

	err := m.migrate.Up()
	if errors.Is(err, migrate.ErrNoChange) {
		return nil
	}
	return err
}

// DownOne 回滚最近应用的一次迁移。
func (m *Migrator) DownOne() error {
	if m == nil {
		return ErrMigratorClosed
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.closed {
		return ErrMigratorClosed
	}

	return m.migrate.Steps(-1)
}

// Version 返回当前迁移版本；全新数据库使用零版本且不是 dirty 状态。
func (m *Migrator) Version() (uint, bool, error) {
	if m == nil {
		return 0, false, ErrMigratorClosed
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.closed {
		return 0, false, ErrMigratorClosed
	}

	version, dirty, err := m.migrate.Version()
	if errors.Is(err, migrate.ErrNilVersion) {
		return 0, false, nil
	}
	return version, dirty, err
}

// Close 释放迁移源和迁移器占用的数据库资源。
func (m *Migrator) Close() error {
	if m == nil {
		return ErrMigratorClosed
	}
	m.closeOnce.Do(func() {
		m.mu.Lock()
		defer m.mu.Unlock()

		m.closed = true
		sourceErr, databaseErr := m.migrate.Close()
		m.closeErr = errors.Join(sourceErr, databaseErr)
	})
	return m.closeErr
}
