//go:build integration

package tests_test

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"sync"
	"testing"
	"time"
	"unicode"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/yokry-he/yok-study-web/examples/go-task-api/internal/config"
	"github.com/yokry-he/yok-study-web/examples/go-task-api/internal/platform/database"
	taskdomain "github.com/yokry-he/yok-study-web/examples/go-task-api/internal/task"
	userdomain "github.com/yokry-he/yok-study-web/examples/go-task-api/internal/user"
)

type postgresFixture struct {
	db               *sql.DB
	migrator         *database.Migrator
	connectionString string
}

func startPostgres(t *testing.T) *postgresFixture {
	t.Helper()

	ctx := context.Background()
	container, err := tcpostgres.Run(ctx, "postgres:18.4",
		tcpostgres.WithDatabase("taskdb"),
		tcpostgres.WithUsername("app"),
		tcpostgres.WithPassword("app"),
		tcpostgres.BasicWaitStrategies(),
	)
	if err != nil {
		t.Fatalf("启动 PostgreSQL 18.4 容器: %v", err)
	}
	t.Cleanup(func() {
		if err := testcontainers.TerminateContainer(container); err != nil {
			t.Errorf("终止 PostgreSQL 测试容器: %v", err)
		}
	})

	connectionString, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("获取 PostgreSQL 连接串: %v", err)
	}

	db, err := database.Open(ctx, testDatabaseConfig(connectionString))
	if err != nil {
		t.Fatalf("打开 PostgreSQL: %v", err)
	}
	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Errorf("关闭 PostgreSQL 连接池: %v", err)
		}
	})

	if err := db.PingContext(ctx); err != nil {
		t.Fatalf("Ping PostgreSQL: %v", err)
	}
	if got := db.Stats().MaxOpenConnections; got != 4 {
		t.Fatalf("MaxOpenConnections = %d, want 4", got)
	}

	migrator, err := database.NewMigrator(db)
	if err != nil {
		t.Fatalf("创建迁移器: %v", err)
	}
	t.Cleanup(func() {
		if err := migrator.Close(); err != nil {
			t.Errorf("关闭迁移器: %v", err)
		}
	})
	if err := migrator.Up(); err != nil {
		t.Fatalf("执行向上迁移: %v", err)
	}

	return &postgresFixture{
		db:               db,
		migrator:         migrator,
		connectionString: connectionString,
	}
}

func TestNewMigratorRejectsNilDatabase(t *testing.T) {
	if _, err := database.NewMigrator(nil); !errors.Is(err, database.ErrNilDatabase) {
		t.Fatalf("NewMigrator(nil) error = %v, want ErrNilDatabase", err)
	}
}

func TestNewMigratorContextRejectsNilContext(t *testing.T) {
	if _, err := database.NewMigratorContext(nil, nil); !errors.Is(err, database.ErrNilContext) {
		t.Fatalf("NewMigratorContext(nil, nil) error = %v, want ErrNilContext", err)
	}
}

func TestOpenDoesNotExposeDSN(t *testing.T) {
	const secret = "never-print-this-password"
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	db, err := database.Open(ctx, config.DatabaseConfig{
		URL:             "postgres://app:" + secret + "@127.0.0.1:1/taskdb?sslmode=disable",
		MaxOpenConns:    1,
		MaxIdleConns:    0,
		ConnMaxLifetime: time.Minute,
		ConnMaxIdleTime: time.Minute,
	})
	if db != nil {
		_ = db.Close()
	}
	if err == nil {
		t.Fatal("Open error = nil, want connection error")
	}
	if strings.Contains(err.Error(), secret) {
		t.Fatalf("Open error exposes DSN password: %v", err)
	}
}

func TestMigratorCloseKeepsCallerPoolOpen(t *testing.T) {
	fixture := startPostgres(t)

	firstErr := fixture.migrator.Close()
	secondErr := fixture.migrator.Close()
	if firstErr != secondErr {
		t.Fatalf("Close errors differ: first=%v second=%v", firstErr, secondErr)
	}
	if err := fixture.db.PingContext(context.Background()); err != nil {
		t.Fatalf("caller pool cannot Ping after Migrator.Close: %v", err)
	}

	if err := fixture.migrator.Up(); !errors.Is(err, database.ErrMigratorClosed) {
		t.Fatalf("Up after Close error = %v, want ErrMigratorClosed", err)
	}
	if err := fixture.migrator.DownOne(); !errors.Is(err, database.ErrMigratorClosed) {
		t.Fatalf("DownOne after Close error = %v, want ErrMigratorClosed", err)
	}
	if _, _, err := fixture.migrator.Version(); !errors.Is(err, database.ErrMigratorClosed) {
		t.Fatalf("Version after Close error = %v, want ErrMigratorClosed", err)
	}
}

func TestNewMigratorContextReleasesConnectionsOnFailure(t *testing.T) {
	fixture := startPostgres(t)

	t.Run("canceled context", func(t *testing.T) {
		before := fixture.db.Stats().InUse
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		migrator, err := database.NewMigratorContext(ctx, fixture.db)
		if migrator != nil {
			_ = migrator.Close()
			t.Fatal("NewMigratorContext returned a migrator for canceled context")
		}
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("NewMigratorContext error = %v, want context.Canceled", err)
		}
		assertEventuallyInUse(t, fixture.db, before)
	})

	t.Run("driver initialization failure", func(t *testing.T) {
		const schema = "migrator_initialization_failure"
		if _, err := fixture.db.Exec(`create schema ` + schema); err != nil {
			t.Fatalf("create isolated schema: %v", err)
		}

		readOnlyURL := withRuntimeParams(t, fixture.connectionString, map[string]string{
			"search_path":                   schema,
			"default_transaction_read_only": "on",
		})
		readOnlyDB, err := database.Open(context.Background(), testDatabaseConfig(readOnlyURL))
		if err != nil {
			t.Fatalf("open read-only pool: %v", err)
		}
		defer readOnlyDB.Close()
		assertEventuallyInUse(t, readOnlyDB, 0)

		migrator, err := database.NewMigratorContext(context.Background(), readOnlyDB)
		if migrator != nil {
			_ = migrator.Close()
			t.Fatal("NewMigratorContext returned a migrator after initialization failure")
		}
		if err == nil {
			t.Fatal("NewMigratorContext initialization error = nil")
		}
		assertEventuallyInUse(t, readOnlyDB, 0)
		if err := readOnlyDB.PingContext(context.Background()); err != nil {
			t.Fatalf("caller pool cannot Ping after initialization failure: %v", err)
		}
	})
}

func TestPostgresMigrationContract(t *testing.T) {
	fixture := startPostgres(t)

	t.Run("schema objects have Chinese comments", func(t *testing.T) {
		assertSchemaComments(t, fixture.db)
	})

	t.Run("schema columns match contract", func(t *testing.T) {
		assertSchemaColumns(t, fixture.db)
	})

	t.Run("business indexes match contract", func(t *testing.T) {
		assertBusinessIndexes(t, fixture.db)
	})

	t.Run("task owner foreign key uses restrict", func(t *testing.T) {
		var deleteAction string
		err := fixture.db.QueryRow(`
			select con.confdeltype::text
			from pg_constraint con
			join pg_namespace n on n.oid = con.connamespace
			where n.nspname = 'public'
			  and con.conname = 'fk_tasks_owner'`).Scan(&deleteAction)
		if err != nil {
			t.Fatalf("query fk_tasks_owner delete action: %v", err)
		}
		if deleteAction != "r" {
			t.Fatalf("fk_tasks_owner confdeltype = %q, want %q (RESTRICT)", deleteAction, "r")
		}
	})

	t.Run("migration is idempotent and clean at version one", func(t *testing.T) {
		if err := fixture.migrator.Up(); err != nil {
			t.Fatalf("second Up() = %v, want nil", err)
		}
		version, dirty, err := fixture.migrator.Version()
		if err != nil {
			t.Fatalf("Version() error = %v", err)
		}
		if version != 1 || dirty {
			t.Fatalf("Version() = (%d, %t), want (1, false)", version, dirty)
		}
	})

	t.Run("email uniqueness is case insensitive", func(t *testing.T) {
		tx := beginTx(t, fixture.db)
		defer tx.Rollback()

		if _, err := tx.Exec(`insert into users(name, email) values ('用户甲', 'USER@example.com')`); err != nil {
			t.Fatalf("insert first user: %v", err)
		}
		_, err := tx.Exec(`insert into users(name, email) values ('用户乙', 'user@example.com')`)
		assertPGCode(t, err, "23505")
	})

	t.Run("invalid user status violates check constraint", func(t *testing.T) {
		tx := beginTx(t, fixture.db)
		defer tx.Rollback()

		_, err := tx.Exec(`insert into users(name, email, status) values ('用户甲', 'invalid-user@example.com', 'LOCKED')`)
		assertPGCode(t, err, "23514")
	})

	t.Run("invalid task status violates check constraint", func(t *testing.T) {
		tx := beginTx(t, fixture.db)
		defer tx.Rollback()

		ownerID := insertUser(t, tx, "task-status@example.com")
		_, err := tx.Exec(`insert into tasks(owner_id, title, status) values ($1, '错误状态任务', 'ARCHIVED')`, ownerID)
		assertPGCode(t, err, "23514")
	})

	t.Run("task owner cannot be deleted", func(t *testing.T) {
		tx := beginTx(t, fixture.db)
		defer tx.Rollback()

		ownerID := insertUser(t, tx, "owner-delete@example.com")
		if _, err := tx.Exec(`insert into tasks(owner_id, title) values ($1, '受保护任务')`, ownerID); err != nil {
			t.Fatalf("insert task: %v", err)
		}
		_, err := tx.Exec(`delete from users where id = $1`, ownerID)
		assertPGCode(t, err, "23001")
	})

	t.Run("negative versions violate check constraints", func(t *testing.T) {
		t.Run("user", func(t *testing.T) {
			tx := beginTx(t, fixture.db)
			defer tx.Rollback()

			_, err := tx.Exec(`insert into users(name, email, version) values ('用户甲', 'negative-user-version@example.com', -1)`)
			assertPGCode(t, err, "23514")
		})

		t.Run("task", func(t *testing.T) {
			tx := beginTx(t, fixture.db)
			defer tx.Rollback()

			ownerID := insertUser(t, tx, "negative-task-version@example.com")
			_, err := tx.Exec(`insert into tasks(owner_id, title, version) values ($1, '负版本任务', -1)`, ownerID)
			assertPGCode(t, err, "23514")
		})
	})

}

func TestPostgresDownOneRemovesSchemaAndClearsVersion(t *testing.T) {
	fixture := startPostgres(t)

	if err := fixture.migrator.DownOne(); err != nil {
		t.Fatalf("DownOne() error = %v", err)
	}
	assertTableMissing(t, fixture.db, "users")
	assertTableMissing(t, fixture.db, "tasks")

	version, dirty, err := fixture.migrator.Version()
	if err != nil {
		t.Fatalf("Version() after DownOne error = %v", err)
	}
	if version != 0 || dirty {
		t.Fatalf("Version() after DownOne = (%d, %t), want (0, false)", version, dirty)
	}
}

func TestRepositoryConstructorsHandleNilDatabase(t *testing.T) {
	ctx := context.Background()

	users := userdomain.NewPostgresRepository(nil)
	if _, err := users.Get(ctx, 1); err == nil {
		t.Fatal("user repository Get with nil database error = nil")
	}

	tasks := taskdomain.NewPostgresRepository(nil)
	if _, err := tasks.Get(ctx, 1); err == nil {
		t.Fatal("task repository Get with nil database error = nil")
	}
}

func TestUserPostgresRepository(t *testing.T) {
	fixture := startPostgres(t)
	repository := userdomain.NewPostgresRepository(fixture.db)

	t.Run("create and get return every field", func(t *testing.T) {
		resetRepositoryTables(t, fixture.db)

		created, err := repository.Create(context.Background(), userdomain.CreateParams{
			Name:   "张三",
			Email:  "USER@example.com",
			Status: userdomain.StatusActive,
		})
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}
		if created.ID <= 0 || created.Name != "张三" || created.Email != "USER@example.com" ||
			created.Status != userdomain.StatusActive || created.Version != 0 ||
			created.CreatedAt.IsZero() || created.UpdatedAt.IsZero() {
			t.Fatalf("Create() = %+v, want complete persisted user", created)
		}

		got, err := repository.Get(context.Background(), created.ID)
		if err != nil {
			t.Fatalf("Get() error = %v", err)
		}
		if got != created {
			t.Fatalf("Get() = %+v, want %+v", got, created)
		}
	})

	t.Run("email conflict uses PostgreSQL error code", func(t *testing.T) {
		resetRepositoryTables(t, fixture.db)

		_, err := repository.Create(context.Background(), userdomain.CreateParams{
			Name: "用户甲", Email: "CaseSensitive@example.com", Status: userdomain.StatusActive,
		})
		if err != nil {
			t.Fatalf("first Create() error = %v", err)
		}
		_, err = repository.Create(context.Background(), userdomain.CreateParams{
			Name: "用户乙", Email: "casesensitive@example.com", Status: userdomain.StatusActive,
		})
		if !errors.Is(err, userdomain.ErrEmailConflict) {
			t.Fatalf("duplicate Create() error = %v, want ErrEmailConflict", err)
		}
	})

	t.Run("list filters pages stably and keeps total past last page", func(t *testing.T) {
		resetRepositoryTables(t, fixture.db)

		statuses := []userdomain.Status{
			userdomain.StatusActive,
			userdomain.StatusDisabled,
			userdomain.StatusActive,
			userdomain.StatusDisabled,
		}
		created := make([]userdomain.User, 0, len(statuses))
		for index, status := range statuses {
			value, err := repository.Create(context.Background(), userdomain.CreateParams{
				Name:   fmt.Sprintf("用户%d", index+1),
				Email:  fmt.Sprintf("user-list-%d@example.com", index+1),
				Status: status,
			})
			if err != nil {
				t.Fatalf("Create(%d) error = %v", index, err)
			}
			created = append(created, value)
		}
		setRepositoryCreatedAt(t, fixture.db, "users")

		firstPage, err := repository.List(context.Background(), userdomain.ListFilter{Limit: 2, Offset: 0})
		if err != nil {
			t.Fatalf("List(first page) error = %v", err)
		}
		assertUserIDs(t, firstPage.Items, []int64{created[3].ID, created[2].ID})
		if firstPage.Total != 4 {
			t.Fatalf("List(first page) total = %d, want 4", firstPage.Total)
		}

		disabled := userdomain.StatusDisabled
		filtered, err := repository.List(context.Background(), userdomain.ListFilter{
			Status: &disabled,
			Limit:  10,
			Offset: 0,
		})
		if err != nil {
			t.Fatalf("List(filtered) error = %v", err)
		}
		assertUserIDs(t, filtered.Items, []int64{created[3].ID, created[1].ID})
		if filtered.Total != 2 {
			t.Fatalf("List(filtered) total = %d, want 2", filtered.Total)
		}

		pastEnd, err := repository.List(context.Background(), userdomain.ListFilter{Limit: 2, Offset: 100})
		if err != nil {
			t.Fatalf("List(past end) error = %v", err)
		}
		if len(pastEnd.Items) != 0 || pastEnd.Total != 4 {
			t.Fatalf("List(past end) = %+v, want empty items and total 4", pastEnd)
		}
	})

	t.Run("missing rows and canceled context retain identity", func(t *testing.T) {
		resetRepositoryTables(t, fixture.db)

		if _, err := repository.Get(context.Background(), 999); !errors.Is(err, userdomain.ErrNotFound) {
			t.Fatalf("Get(missing) error = %v, want ErrNotFound", err)
		}
		if _, err := repository.UpdateStatus(context.Background(), 999, userdomain.StatusDisabled, 0); !errors.Is(err, userdomain.ErrNotFound) {
			t.Fatalf("UpdateStatus(missing) error = %v, want ErrNotFound", err)
		}

		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		if _, err := repository.Get(ctx, 1); !errors.Is(err, context.Canceled) {
			t.Fatalf("Get(canceled) error = %v, want context.Canceled", err)
		}
		if _, err := repository.List(ctx, userdomain.ListFilter{Limit: 20}); !errors.Is(err, context.Canceled) {
			t.Fatalf("List(canceled) error = %v, want context.Canceled", err)
		}
		if _, err := repository.UpdateStatus(ctx, 1, userdomain.StatusDisabled, 0); !errors.Is(err, context.Canceled) {
			t.Fatalf("UpdateStatus(canceled) error = %v, want context.Canceled", err)
		}
	})

	t.Run("closed pool error preserves identity without leaking connection details", func(t *testing.T) {
		closedDB, err := sql.Open("pgx", fixture.connectionString)
		if err != nil {
			t.Fatalf("sql.Open(closed pool): %v", err)
		}
		if err := closedDB.Close(); err != nil {
			t.Fatalf("close isolated pool: %v", err)
		}

		closedRepository := userdomain.NewPostgresRepository(closedDB)
		_, err = closedRepository.Get(context.Background(), 1)
		assertSafeRepositoryError(t, err, fixture.connectionString, "select", "users")
		if errors.Unwrap(err) == nil {
			t.Fatalf("closed pool error = %v, want an unwrap-able cause", err)
		}
	})
}

func TestOptimisticUpdateAllowsOnlyOneWriter(t *testing.T) {
	fixture := startPostgres(t)
	repository := userdomain.NewPostgresRepository(fixture.db)
	created, err := repository.Create(context.Background(), userdomain.CreateParams{
		Name: "并发用户", Email: "optimistic-user@example.com", Status: userdomain.StatusActive,
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	statuses := []userdomain.Status{userdomain.StatusDisabled, userdomain.StatusActive}
	results := make(chan error, len(statuses))
	var wait sync.WaitGroup
	for _, status := range statuses {
		wait.Add(1)
		go func(next userdomain.Status) {
			defer wait.Done()
			_, updateErr := repository.UpdateStatus(context.Background(), created.ID, next, created.Version)
			results <- updateErr
		}(status)
	}
	wait.Wait()
	close(results)

	var successes, conflicts int
	for result := range results {
		switch {
		case result == nil:
			successes++
		case errors.Is(result, userdomain.ErrVersionConflict):
			conflicts++
		default:
			t.Fatalf("UpdateStatus() error = %v, want nil or ErrVersionConflict", result)
		}
	}
	if successes != 1 || conflicts != 1 {
		t.Fatalf("concurrent results: successes=%d conflicts=%d, want 1 and 1", successes, conflicts)
	}
}

func TestUserOptimisticWriteClassifiesConcurrentDeleteAsConflict(t *testing.T) {
	fixture := startPostgres(t)
	repository := userdomain.NewPostgresRepository(fixture.db)
	created, err := repository.Create(context.Background(), userdomain.CreateParams{
		Name: "并发删除用户", Email: "delete-race-user@example.com", Status: userdomain.StatusActive,
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	deleteTx := beginTx(t, fixture.db)
	if _, err := deleteTx.Exec(`delete from users where id = $1`, created.ID); err != nil {
		_ = deleteTx.Rollback()
		t.Fatalf("delete user in concurrent transaction: %v", err)
	}

	result := make(chan error, 1)
	go func() {
		_, updateErr := repository.UpdateStatus(
			context.Background(),
			created.ID,
			userdomain.StatusDisabled,
			created.Version,
		)
		result <- updateErr
	}()
	waitForBlockedRepositoryQuery(t, fixture.db, "update users")
	if err := deleteTx.Commit(); err != nil {
		t.Fatalf("commit concurrent delete: %v", err)
	}

	select {
	case err := <-result:
		if !errors.Is(err, userdomain.ErrVersionConflict) {
			t.Fatalf("UpdateStatus after concurrent delete error = %v, want ErrVersionConflict", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("UpdateStatus remained blocked after concurrent delete committed")
	}
}

func TestTaskPostgresRepository(t *testing.T) {
	fixture := startPostgres(t)
	users := userdomain.NewPostgresRepository(fixture.db)
	tasks := taskdomain.NewPostgresRepository(fixture.db)

	t.Run("create and get preserve nullable values without aliases", func(t *testing.T) {
		resetRepositoryTables(t, fixture.db)
		owner := createRepositoryUser(t, users, "task-null-owner@example.com")
		dueAt := time.Date(2032, time.March, 4, 5, 6, 7, 0, time.UTC)
		description := "原始描述"

		created, err := tasks.Create(context.Background(), taskdomain.CreateParams{
			OwnerID: owner.ID, Title: "完整任务", Description: &description,
			Status: taskdomain.StatusTodo, DueAt: &dueAt,
		})
		if err != nil {
			t.Fatalf("Create(non-null) error = %v", err)
		}
		description = "调用方修改"
		dueAt = dueAt.Add(24 * time.Hour)
		if created.Description == nil || *created.Description != "原始描述" {
			t.Fatalf("Create() description = %v, want independent original value", created.Description)
		}
		wantDueAt := time.Date(2032, time.March, 4, 5, 6, 7, 0, time.UTC)
		if created.DueAt == nil || !created.DueAt.Equal(wantDueAt) {
			t.Fatalf("Create() dueAt = %v, want %v", created.DueAt, wantDueAt)
		}
		if created.ID <= 0 || created.Version != 0 || created.CreatedAt.IsZero() || created.UpdatedAt.IsZero() {
			t.Fatalf("Create() = %+v, want complete persisted task", created)
		}

		got, err := tasks.Get(context.Background(), created.ID)
		if err != nil {
			t.Fatalf("Get() error = %v", err)
		}
		if got.Description == nil || *got.Description != "原始描述" || got.DueAt == nil || !got.DueAt.Equal(wantDueAt) {
			t.Fatalf("Get() = %+v, want persisted nullable values", got)
		}
		*got.Description = "本地修改"
		gotAgain, err := tasks.Get(context.Background(), created.ID)
		if err != nil {
			t.Fatalf("second Get() error = %v", err)
		}
		if gotAgain.Description == nil || *gotAgain.Description != "原始描述" {
			t.Fatalf("second Get() description = %v, want independent scan value", gotAgain.Description)
		}

		nullTask, err := tasks.Create(context.Background(), taskdomain.CreateParams{
			OwnerID: owner.ID, Title: "空值任务", Status: taskdomain.StatusTodo,
		})
		if err != nil {
			t.Fatalf("Create(null) error = %v", err)
		}
		if nullTask.Description != nil || nullTask.DueAt != nil {
			t.Fatalf("Create(null) = %+v, want nil description and dueAt", nullTask)
		}
	})

	t.Run("missing owner maps foreign key error", func(t *testing.T) {
		resetRepositoryTables(t, fixture.db)

		_, err := tasks.Create(context.Background(), taskdomain.CreateParams{
			OwnerID: 999, Title: "无负责人任务", Status: taskdomain.StatusTodo,
		})
		if !errors.Is(err, taskdomain.ErrOwnerNotFound) {
			t.Fatalf("Create(missing owner) error = %v, want ErrOwnerNotFound", err)
		}
	})

	t.Run("list filters pages stably and keeps total past last page", func(t *testing.T) {
		resetRepositoryTables(t, fixture.db)
		ownerA := createRepositoryUser(t, users, "task-list-owner-a@example.com")
		ownerB := createRepositoryUser(t, users, "task-list-owner-b@example.com")
		params := []taskdomain.CreateParams{
			{OwnerID: ownerA.ID, Title: "任务一", Status: taskdomain.StatusTodo},
			{OwnerID: ownerA.ID, Title: "任务二", Status: taskdomain.StatusDoing},
			{OwnerID: ownerB.ID, Title: "任务三", Status: taskdomain.StatusTodo},
			{OwnerID: ownerA.ID, Title: "任务四", Status: taskdomain.StatusTodo},
		}
		created := make([]taskdomain.Task, 0, len(params))
		for index, item := range params {
			value, err := tasks.Create(context.Background(), item)
			if err != nil {
				t.Fatalf("Create(%d) error = %v", index, err)
			}
			created = append(created, value)
		}
		setRepositoryCreatedAt(t, fixture.db, "tasks")

		page, err := tasks.List(context.Background(), taskdomain.ListFilter{Limit: 2})
		if err != nil {
			t.Fatalf("List(first page) error = %v", err)
		}
		assertTaskIDs(t, page.Items, []int64{created[3].ID, created[2].ID})
		if page.Total != 4 {
			t.Fatalf("List(first page) total = %d, want 4", page.Total)
		}

		status := taskdomain.StatusTodo
		filtered, err := tasks.List(context.Background(), taskdomain.ListFilter{
			OwnerID: &ownerA.ID, Status: &status, Limit: 10,
		})
		if err != nil {
			t.Fatalf("List(filtered) error = %v", err)
		}
		assertTaskIDs(t, filtered.Items, []int64{created[3].ID, created[0].ID})
		if filtered.Total != 2 {
			t.Fatalf("List(filtered) total = %d, want 2", filtered.Total)
		}

		ownerOnly, err := tasks.List(context.Background(), taskdomain.ListFilter{
			OwnerID: &ownerA.ID, Limit: 10,
		})
		if err != nil {
			t.Fatalf("List(owner only) error = %v", err)
		}
		assertTaskIDs(t, ownerOnly.Items, []int64{created[3].ID, created[1].ID, created[0].ID})
		if ownerOnly.Total != 3 {
			t.Fatalf("List(owner only) total = %d, want 3", ownerOnly.Total)
		}

		statusOnly, err := tasks.List(context.Background(), taskdomain.ListFilter{
			Status: &status, Limit: 10,
		})
		if err != nil {
			t.Fatalf("List(status only) error = %v", err)
		}
		assertTaskIDs(t, statusOnly.Items, []int64{created[3].ID, created[2].ID, created[0].ID})
		if statusOnly.Total != 3 {
			t.Fatalf("List(status only) total = %d, want 3", statusOnly.Total)
		}

		pastEnd, err := tasks.List(context.Background(), taskdomain.ListFilter{Limit: 2, Offset: 100})
		if err != nil {
			t.Fatalf("List(past end) error = %v", err)
		}
		if len(pastEnd.Items) != 0 || pastEnd.Total != 4 {
			t.Fatalf("List(past end) = %+v, want empty items and total 4", pastEnd)
		}
	})

	t.Run("update status and delete enforce versions", func(t *testing.T) {
		resetRepositoryTables(t, fixture.db)
		owner := createRepositoryUser(t, users, "task-write-owner@example.com")
		created, err := tasks.Create(context.Background(), taskdomain.CreateParams{
			OwnerID: owner.ID, Title: "初始任务", Status: taskdomain.StatusTodo,
		})
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		description := "更新后的描述"
		dueAt := time.Date(2033, time.April, 5, 6, 7, 8, 0, time.UTC)
		updated, err := tasks.Update(context.Background(), created.ID, taskdomain.UpdateParams{
			Title: "更新任务", Description: &description, DueAt: &dueAt,
		}, created.Version)
		if err != nil {
			t.Fatalf("Update() error = %v", err)
		}
		if updated.Version != 1 || updated.Title != "更新任务" || updated.Description == nil ||
			*updated.Description != description || updated.DueAt == nil || !updated.DueAt.Equal(dueAt) {
			t.Fatalf("Update() = %+v, want updated fields and version 1", updated)
		}
		cleared, err := tasks.Update(context.Background(), created.ID, taskdomain.UpdateParams{
			Title: "清空可空字段",
		}, updated.Version)
		if err != nil {
			t.Fatalf("Update(clear nullable fields) error = %v", err)
		}
		if cleared.Version != 2 || cleared.Description != nil || cleared.DueAt != nil {
			t.Fatalf("Update(clear nullable fields) = %+v, want nil fields at version 2", cleared)
		}
		if _, err := tasks.Update(context.Background(), created.ID, taskdomain.UpdateParams{Title: "旧写入"}, 0); !errors.Is(err, taskdomain.ErrVersionConflict) {
			t.Fatalf("stale Update() error = %v, want ErrVersionConflict", err)
		}
		if _, err := tasks.UpdateStatus(context.Background(), created.ID, taskdomain.StatusDoing, updated.Version); !errors.Is(err, taskdomain.ErrVersionConflict) {
			t.Fatalf("stale UpdateStatus() error = %v, want ErrVersionConflict", err)
		}

		changed, err := tasks.UpdateStatus(context.Background(), created.ID, taskdomain.StatusDoing, cleared.Version)
		if err != nil {
			t.Fatalf("UpdateStatus() error = %v", err)
		}
		if changed.Status != taskdomain.StatusDoing || changed.Version != 3 {
			t.Fatalf("UpdateStatus() = %+v, want DOING at version 3", changed)
		}
		if err := tasks.Delete(context.Background(), created.ID, cleared.Version); !errors.Is(err, taskdomain.ErrVersionConflict) {
			t.Fatalf("stale Delete() error = %v, want ErrVersionConflict", err)
		}
		if err := tasks.Delete(context.Background(), created.ID, changed.Version); err != nil {
			t.Fatalf("Delete() error = %v", err)
		}
		if _, err := tasks.Get(context.Background(), created.ID); !errors.Is(err, taskdomain.ErrNotFound) {
			t.Fatalf("Get(deleted) error = %v, want ErrNotFound", err)
		}
	})

	t.Run("missing rows and canceled context retain identity", func(t *testing.T) {
		resetRepositoryTables(t, fixture.db)

		if _, err := tasks.Get(context.Background(), 999); !errors.Is(err, taskdomain.ErrNotFound) {
			t.Fatalf("Get(missing) error = %v, want ErrNotFound", err)
		}
		if _, err := tasks.Update(context.Background(), 999, taskdomain.UpdateParams{Title: "不存在"}, 0); !errors.Is(err, taskdomain.ErrNotFound) {
			t.Fatalf("Update(missing) error = %v, want ErrNotFound", err)
		}
		if _, err := tasks.UpdateStatus(context.Background(), 999, taskdomain.StatusDoing, 0); !errors.Is(err, taskdomain.ErrNotFound) {
			t.Fatalf("UpdateStatus(missing) error = %v, want ErrNotFound", err)
		}
		if err := tasks.Delete(context.Background(), 999, 0); !errors.Is(err, taskdomain.ErrNotFound) {
			t.Fatalf("Delete(missing) error = %v, want ErrNotFound", err)
		}

		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		if _, err := tasks.Get(ctx, 1); !errors.Is(err, context.Canceled) {
			t.Fatalf("Get(canceled) error = %v, want context.Canceled", err)
		}
		if _, err := tasks.List(ctx, taskdomain.ListFilter{Limit: 20}); !errors.Is(err, context.Canceled) {
			t.Fatalf("List(canceled) error = %v, want context.Canceled", err)
		}
		if _, err := tasks.Update(ctx, 1, taskdomain.UpdateParams{Title: "取消更新"}, 0); !errors.Is(err, context.Canceled) {
			t.Fatalf("Update(canceled) error = %v, want context.Canceled", err)
		}
		if _, err := tasks.UpdateStatus(ctx, 1, taskdomain.StatusDoing, 0); !errors.Is(err, context.Canceled) {
			t.Fatalf("UpdateStatus(canceled) error = %v, want context.Canceled", err)
		}
		if err := tasks.Delete(ctx, 1, 0); !errors.Is(err, context.Canceled) {
			t.Fatalf("Delete(canceled) error = %v, want context.Canceled", err)
		}
	})

	t.Run("closed pool error preserves identity without leaking connection details", func(t *testing.T) {
		closedDB, err := sql.Open("pgx", fixture.connectionString)
		if err != nil {
			t.Fatalf("sql.Open(closed pool): %v", err)
		}
		if err := closedDB.Close(); err != nil {
			t.Fatalf("close isolated pool: %v", err)
		}

		closedRepository := taskdomain.NewPostgresRepository(closedDB)
		_, err = closedRepository.Get(context.Background(), 1)
		assertSafeRepositoryError(t, err, fixture.connectionString, "select", "tasks")
		if errors.Unwrap(err) == nil {
			t.Fatalf("closed pool error = %v, want an unwrap-able cause", err)
		}
	})
}

func TestTaskOptimisticDeleteClassifiesConcurrentDeleteAsConflict(t *testing.T) {
	fixture := startPostgres(t)
	users := userdomain.NewPostgresRepository(fixture.db)
	tasks := taskdomain.NewPostgresRepository(fixture.db)
	owner := createRepositoryUser(t, users, "delete-race-task-owner@example.com")
	created, err := tasks.Create(context.Background(), taskdomain.CreateParams{
		OwnerID: owner.ID, Title: "并发删除任务", Status: taskdomain.StatusTodo,
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	deleteTx := beginTx(t, fixture.db)
	if _, err := deleteTx.Exec(`delete from tasks where id = $1`, created.ID); err != nil {
		_ = deleteTx.Rollback()
		t.Fatalf("delete task in concurrent transaction: %v", err)
	}

	result := make(chan error, 1)
	go func() {
		result <- tasks.Delete(context.Background(), created.ID, created.Version)
	}()
	waitForBlockedRepositoryQuery(t, fixture.db, "delete from tasks")
	if err := deleteTx.Commit(); err != nil {
		t.Fatalf("commit concurrent delete: %v", err)
	}

	select {
	case err := <-result:
		if !errors.Is(err, taskdomain.ErrVersionConflict) {
			t.Fatalf("Delete after concurrent delete error = %v, want ErrVersionConflict", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Delete remained blocked after concurrent delete committed")
	}
}

func resetRepositoryTables(t *testing.T, db *sql.DB) {
	t.Helper()
	if _, err := db.Exec(`truncate table tasks, users restart identity`); err != nil {
		t.Fatalf("reset repository tables: %v", err)
	}
}

func setRepositoryCreatedAt(t *testing.T, db *sql.DB, table string) {
	t.Helper()
	if table != "users" && table != "tasks" {
		t.Fatalf("unsupported repository table %q", table)
	}
	query := fmt.Sprintf(`update %s set created_at = '2030-01-02T03:04:05Z'`, table)
	if _, err := db.Exec(query); err != nil {
		t.Fatalf("set %s created_at: %v", table, err)
	}
}

func createRepositoryUser(t *testing.T, repository *userdomain.PostgresRepository, email string) userdomain.User {
	t.Helper()
	created, err := repository.Create(context.Background(), userdomain.CreateParams{
		Name: "任务负责人", Email: email, Status: userdomain.StatusActive,
	})
	if err != nil {
		t.Fatalf("create repository user %s: %v", email, err)
	}
	return created
}

func assertUserIDs(t *testing.T, users []userdomain.User, want []int64) {
	t.Helper()
	got := make([]int64, len(users))
	for index := range users {
		got[index] = users[index].ID
	}
	if fmt.Sprint(got) != fmt.Sprint(want) {
		t.Fatalf("user IDs = %v, want %v", got, want)
	}
}

func assertTaskIDs(t *testing.T, tasks []taskdomain.Task, want []int64) {
	t.Helper()
	got := make([]int64, len(tasks))
	for index := range tasks {
		got[index] = tasks[index].ID
	}
	if fmt.Sprint(got) != fmt.Sprint(want) {
		t.Fatalf("task IDs = %v, want %v", got, want)
	}
}

type expectedColumn struct {
	dataType        string
	maxLength       int64
	nullable        bool
	defaultSemantic string
	identity        bool
}

func assertSchemaColumns(t *testing.T, db *sql.DB) {
	t.Helper()

	want := map[string]expectedColumn{
		"users.id":         {dataType: "bigint", identity: true},
		"users.name":       {dataType: "character varying", maxLength: 64},
		"users.email":      {dataType: "character varying", maxLength: 254},
		"users.status":     {dataType: "character varying", maxLength: 16, defaultSemantic: "ACTIVE"},
		"users.version":    {dataType: "bigint", defaultSemantic: "zero"},
		"users.created_at": {dataType: "timestamp with time zone", defaultSemantic: "now"},
		"users.updated_at": {dataType: "timestamp with time zone", defaultSemantic: "now"},
		"tasks.id":         {dataType: "bigint", identity: true},
		"tasks.owner_id":   {dataType: "bigint"},
		"tasks.title":      {dataType: "character varying", maxLength: 128},
		"tasks.description": {
			dataType: "text",
			nullable: true,
		},
		"tasks.status":     {dataType: "character varying", maxLength: 16, defaultSemantic: "TODO"},
		"tasks.due_at":     {dataType: "timestamp with time zone", nullable: true},
		"tasks.version":    {dataType: "bigint", defaultSemantic: "zero"},
		"tasks.created_at": {dataType: "timestamp with time zone", defaultSemantic: "now"},
		"tasks.updated_at": {dataType: "timestamp with time zone", defaultSemantic: "now"},
	}

	rows, err := db.Query(`
		select table_name, column_name, data_type, character_maximum_length,
		       is_nullable, column_default, is_identity, identity_generation
		from information_schema.columns
		where table_schema = 'public'
		  and table_name in ('users', 'tasks')
		order by table_name, ordinal_position`)
	if err != nil {
		t.Fatalf("query schema columns: %v", err)
	}
	defer rows.Close()

	seen := make(map[string]struct{}, len(want))
	for rows.Next() {
		var tableName, columnName, dataType, nullable, identity string
		var maxLength sql.NullInt64
		var columnDefault, identityGeneration sql.NullString
		if err := rows.Scan(
			&tableName,
			&columnName,
			&dataType,
			&maxLength,
			&nullable,
			&columnDefault,
			&identity,
			&identityGeneration,
		); err != nil {
			t.Fatalf("scan schema column: %v", err)
		}

		name := tableName + "." + columnName
		expected, ok := want[name]
		if !ok {
			t.Errorf("unexpected schema column %s", name)
			continue
		}
		seen[name] = struct{}{}

		if dataType != expected.dataType {
			t.Errorf("%s data_type = %q, want %q", name, dataType, expected.dataType)
		}
		if expected.maxLength == 0 {
			if maxLength.Valid {
				t.Errorf("%s character_maximum_length = %d, want NULL", name, maxLength.Int64)
			}
		} else if !maxLength.Valid || maxLength.Int64 != expected.maxLength {
			t.Errorf("%s character_maximum_length = %v, want %d", name, maxLength, expected.maxLength)
		}
		if gotNullable := nullable == "YES"; gotNullable != expected.nullable {
			t.Errorf("%s nullable = %t, want %t", name, gotNullable, expected.nullable)
		}
		if gotIdentity := identity == "YES"; gotIdentity != expected.identity {
			t.Errorf("%s identity = %t, want %t", name, gotIdentity, expected.identity)
		}
		if expected.identity && (!identityGeneration.Valid || identityGeneration.String != "ALWAYS") {
			t.Errorf("%s identity_generation = %v, want ALWAYS", name, identityGeneration)
		}
		assertColumnDefault(t, name, columnDefault, expected.defaultSemantic)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("iterate schema columns: %v", err)
	}

	for name := range want {
		if _, ok := seen[name]; !ok {
			t.Errorf("schema column %s is missing", name)
		}
	}
	if len(seen) != len(want) {
		t.Errorf("schema column count = %d, want %d", len(seen), len(want))
	}
}

func assertColumnDefault(t *testing.T, name string, got sql.NullString, semantic string) {
	t.Helper()

	if semantic == "" {
		if got.Valid {
			t.Errorf("%s default = %q, want no default", name, got.String)
		}
		return
	}
	if !got.Valid {
		t.Errorf("%s has no default, want %s semantic", name, semantic)
		return
	}

	value := strings.TrimSpace(got.String)
	switch semantic {
	case "zero":
		if value != "0" {
			t.Errorf("%s default = %q, want numeric zero", name, got.String)
		}
	case "now":
		if !strings.Contains(strings.ToLower(value), "now()") {
			t.Errorf("%s default = %q, want current-time semantic", name, got.String)
		}
	default:
		if !strings.Contains(value, "'"+semantic+"'") {
			t.Errorf("%s default = %q, want %q semantic", name, got.String, semantic)
		}
	}
}

type indexElement struct {
	position   int
	columnName string
	definition string
	unique     bool
	descending bool
}

type expectedIndexColumn struct {
	name       string
	descending bool
}

func assertBusinessIndexes(t *testing.T, db *sql.DB) {
	t.Helper()

	rows, err := db.Query(`
		select index_class.relname,
		       key.position,
		       coalesce(attribute.attname, ''),
		       pg_get_indexdef(index_data.indexrelid, key.position::integer, true),
		       index_data.indisunique,
		       (key.options & 1) = 1
		from pg_index index_data
		join pg_class index_class on index_class.oid = index_data.indexrelid
		join pg_class table_class on table_class.oid = index_data.indrelid
		join pg_namespace namespace on namespace.oid = table_class.relnamespace
		cross join lateral unnest(
		  index_data.indkey::smallint[],
		  index_data.indoption::smallint[]
		) with ordinality as key(attnum, options, position)
		left join pg_attribute attribute
		  on attribute.attrelid = table_class.oid
		 and attribute.attnum = key.attnum
		where namespace.nspname = 'public'
		  and index_class.relname in (
		    'uk_users_email_lower',
		    'idx_users_created_id',
		    'idx_users_status_created_id',
		    'idx_tasks_owner_status',
		    'idx_tasks_status_due_at',
		    'idx_tasks_created_id',
		    'idx_tasks_owner_created_id',
		    'idx_tasks_status_created_id',
		    'idx_tasks_owner_status_created_id'
		  )
		order by index_class.relname, key.position`)
	if err != nil {
		t.Fatalf("query business indexes: %v", err)
	}
	defer rows.Close()

	got := make(map[string][]indexElement, 9)
	for rows.Next() {
		var name string
		var element indexElement
		if err := rows.Scan(
			&name,
			&element.position,
			&element.columnName,
			&element.definition,
			&element.unique,
			&element.descending,
		); err != nil {
			t.Fatalf("scan business index: %v", err)
		}
		got[name] = append(got[name], element)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("iterate business indexes: %v", err)
	}

	assertExpressionIndex(t, got["uk_users_email_lower"])
	assertColumnIndex(t, "idx_tasks_owner_status", got["idx_tasks_owner_status"], []string{"owner_id", "status", "id"})
	assertColumnIndex(t, "idx_tasks_status_due_at", got["idx_tasks_status_due_at"], []string{"status", "due_at", "id"})
	assertOrderedColumnIndex(t, "idx_users_created_id", got["idx_users_created_id"], []expectedIndexColumn{
		{name: "created_at", descending: true}, {name: "id", descending: true},
	})
	assertOrderedColumnIndex(t, "idx_users_status_created_id", got["idx_users_status_created_id"], []expectedIndexColumn{
		{name: "status"}, {name: "created_at", descending: true}, {name: "id", descending: true},
	})
	assertOrderedColumnIndex(t, "idx_tasks_created_id", got["idx_tasks_created_id"], []expectedIndexColumn{
		{name: "created_at", descending: true}, {name: "id", descending: true},
	})
	assertOrderedColumnIndex(t, "idx_tasks_owner_created_id", got["idx_tasks_owner_created_id"], []expectedIndexColumn{
		{name: "owner_id"}, {name: "created_at", descending: true}, {name: "id", descending: true},
	})
	assertOrderedColumnIndex(t, "idx_tasks_status_created_id", got["idx_tasks_status_created_id"], []expectedIndexColumn{
		{name: "status"}, {name: "created_at", descending: true}, {name: "id", descending: true},
	})
	assertOrderedColumnIndex(t, "idx_tasks_owner_status_created_id", got["idx_tasks_owner_status_created_id"], []expectedIndexColumn{
		{name: "owner_id"}, {name: "status"}, {name: "created_at", descending: true}, {name: "id", descending: true},
	})
	if len(got) != 9 {
		t.Errorf("business index count = %d, want 9; got %v", len(got), got)
	}
}

func assertExpressionIndex(t *testing.T, elements []indexElement) {
	t.Helper()

	if len(elements) != 1 {
		t.Fatalf("uk_users_email_lower element count = %d, want 1", len(elements))
	}
	element := elements[0]
	definition := strings.ToLower(element.definition)
	if element.position != 1 || element.columnName != "" || !element.unique {
		t.Errorf("uk_users_email_lower metadata = %+v, want one unique expression", element)
	}
	if !strings.Contains(definition, "lower(") || !strings.Contains(definition, "email") {
		t.Errorf("uk_users_email_lower expression = %q, want lower(email) semantic", element.definition)
	}
}

func assertColumnIndex(t *testing.T, name string, elements []indexElement, wantColumns []string) {
	t.Helper()

	if len(elements) != len(wantColumns) {
		t.Fatalf("%s element count = %d, want %d", name, len(elements), len(wantColumns))
	}
	for index, wantColumn := range wantColumns {
		element := elements[index]
		if element.position != index+1 || element.columnName != wantColumn || element.unique {
			t.Errorf(
				"%s element %d = %+v, want non-unique column %s at position %d",
				name,
				index,
				element,
				wantColumn,
				index+1,
			)
		}
	}
}

func assertOrderedColumnIndex(t *testing.T, name string, elements []indexElement, want []expectedIndexColumn) {
	t.Helper()

	if len(elements) != len(want) {
		t.Fatalf("%s element count = %d, want %d", name, len(elements), len(want))
	}
	for index, expected := range want {
		element := elements[index]
		if element.position != index+1 || element.columnName != expected.name || element.unique || element.descending != expected.descending {
			t.Errorf(
				"%s element %d = %+v, want column=%s descending=%t position=%d",
				name,
				index,
				element,
				expected.name,
				expected.descending,
				index+1,
			)
		}
	}
}

func assertSchemaComments(t *testing.T, db *sql.DB) {
	t.Helper()

	assertCatalogComments(t, db, `
		select c.relname, obj_description(c.oid, 'pg_class')
		from pg_class c
		join pg_namespace n on n.oid = c.relnamespace
		where n.nspname = 'public'
		  and c.relkind = 'r'
		  and c.relname in ('users', 'tasks')`, []string{"tasks", "users"})

	assertCatalogComments(t, db, `
		select c.relname || '.' || a.attname, col_description(c.oid, a.attnum)
		from pg_class c
		join pg_namespace n on n.oid = c.relnamespace
		join pg_attribute a on a.attrelid = c.oid
		where n.nspname = 'public'
		  and c.relname in ('users', 'tasks')
		  and a.attnum > 0
		  and not a.attisdropped`, []string{
		"tasks.created_at", "tasks.description", "tasks.due_at", "tasks.id", "tasks.owner_id",
		"tasks.status", "tasks.title", "tasks.updated_at", "tasks.version",
		"users.created_at", "users.email", "users.id", "users.name", "users.status",
		"users.updated_at", "users.version",
	})

	assertCatalogComments(t, db, `
		select con.conname, obj_description(con.oid, 'pg_constraint')
		from pg_constraint con
		join pg_namespace n on n.oid = con.connamespace
		where n.nspname = 'public'
		  and con.conname in (
		    'pk_users', 'ck_users_status', 'ck_users_version',
		    'pk_tasks', 'fk_tasks_owner', 'ck_tasks_status', 'ck_tasks_version'
		  )`, []string{
		"ck_tasks_status", "ck_tasks_version", "ck_users_status", "ck_users_version",
		"fk_tasks_owner", "pk_tasks", "pk_users",
	})

	assertCatalogComments(t, db, `
		select c.relname, obj_description(c.oid, 'pg_class')
		from pg_class c
		join pg_namespace n on n.oid = c.relnamespace
		where n.nspname = 'public'
		  and c.relkind = 'i'
		  and c.relname in (
		    'pk_users', 'uk_users_email_lower', 'pk_tasks',
		    'idx_users_created_id', 'idx_users_status_created_id',
		    'idx_tasks_owner_status', 'idx_tasks_status_due_at',
		    'idx_tasks_created_id', 'idx_tasks_owner_created_id',
		    'idx_tasks_status_created_id', 'idx_tasks_owner_status_created_id'
		  )`, []string{
		"idx_tasks_created_id", "idx_tasks_owner_created_id", "idx_tasks_owner_status",
		"idx_tasks_owner_status_created_id", "idx_tasks_status_created_id", "idx_tasks_status_due_at",
		"idx_users_created_id", "idx_users_status_created_id", "pk_tasks", "pk_users", "uk_users_email_lower",
	})
}

func assertSafeRepositoryError(t *testing.T, err error, connectionString string, forbidden ...string) {
	t.Helper()
	if err == nil {
		t.Fatal("repository error = nil")
	}
	message := strings.ToLower(err.Error())
	if strings.Contains(message, strings.ToLower(connectionString)) || strings.Contains(message, "postgres://") || strings.Contains(message, "password") {
		t.Fatalf("repository error leaked connection details: %v", err)
	}
	for _, value := range forbidden {
		if strings.Contains(message, strings.ToLower(value)) {
			t.Fatalf("repository error leaked %q: %v", value, err)
		}
	}
}

func waitForBlockedRepositoryQuery(t *testing.T, db *sql.DB, queryFragment string) {
	t.Helper()

	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		var blocked bool
		err := db.QueryRow(`
			select exists (
				select 1
				from pg_stat_activity
				where pid <> pg_backend_pid()
				  and datname = current_database()
				  and state = 'active'
				  and wait_event_type = 'Lock'
				  and position(lower($1) in lower(query)) > 0
			)`, queryFragment).Scan(&blocked)
		if err != nil {
			t.Fatalf("inspect blocked repository query: %v", err)
		}
		if blocked {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("repository query containing %q did not block before timeout", queryFragment)
}

func assertCatalogComments(t *testing.T, db *sql.DB, query string, wantNames []string) {
	t.Helper()

	rows, err := db.Query(query)
	if err != nil {
		t.Fatalf("query catalog comments: %v", err)
	}
	defer rows.Close()

	got := make(map[string]string, len(wantNames))
	for rows.Next() {
		var name string
		var comment sql.NullString
		if err := rows.Scan(&name, &comment); err != nil {
			t.Fatalf("scan catalog comment: %v", err)
		}
		if !comment.Valid || strings.TrimSpace(comment.String) == "" {
			t.Errorf("%s has no comment", name)
			continue
		}
		if !containsHan(comment.String) {
			t.Errorf("%s comment is not Chinese: %q", name, comment.String)
		}
		got[name] = comment.String
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("iterate catalog comments: %v", err)
	}

	for _, name := range wantNames {
		if _, ok := got[name]; !ok {
			t.Errorf("catalog object %s missing or lacks a Chinese comment", name)
		}
	}
	if len(got) != len(wantNames) {
		t.Errorf("catalog object count = %d, want %d; got %v", len(got), len(wantNames), got)
	}
}

func containsHan(value string) bool {
	for _, r := range value {
		if unicode.Is(unicode.Han, r) {
			return true
		}
	}
	return false
}

func beginTx(t *testing.T, db *sql.DB) *sql.Tx {
	t.Helper()

	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("begin transaction: %v", err)
	}
	return tx
}

func insertUser(t *testing.T, tx *sql.Tx, email string) int64 {
	t.Helper()

	var id int64
	if err := tx.QueryRow(`insert into users(name, email) values ('测试用户', $1) returning id`, email).Scan(&id); err != nil {
		t.Fatalf("insert user %s: %v", email, err)
	}
	return id
}

func assertPGCode(t *testing.T, err error, want string) {
	t.Helper()

	if err == nil {
		t.Fatalf("SQL error = nil, want PostgreSQL code %s", want)
	}
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		t.Fatalf("SQL error type = %T (%v), want *pgconn.PgError", err, err)
	}
	if pgErr.Code != want {
		t.Fatalf("PostgreSQL code = %s (%s), want %s", pgErr.Code, pgErr.Message, want)
	}
}

func assertTableMissing(t *testing.T, db *sql.DB, table string) {
	t.Helper()

	var exists bool
	query := fmt.Sprintf(`select to_regclass('public.%s') is not null`, table)
	if err := db.QueryRow(query).Scan(&exists); err != nil {
		t.Fatalf("check table %s: %v", table, err)
	}
	if exists {
		t.Fatalf("table %s still exists after DownOne", table)
	}
}

func testDatabaseConfig(connectionString string) config.DatabaseConfig {
	return config.DatabaseConfig{
		URL:             connectionString,
		MaxOpenConns:    4,
		MaxIdleConns:    2,
		ConnMaxLifetime: 5 * time.Minute,
		ConnMaxIdleTime: time.Minute,
	}
}

func withRuntimeParams(t *testing.T, connectionString string, params map[string]string) string {
	t.Helper()

	parsed, err := url.Parse(connectionString)
	if err != nil {
		t.Fatalf("parse PostgreSQL connection string: %v", err)
	}
	query := parsed.Query()
	for name, value := range params {
		query.Set(name, value)
	}
	parsed.RawQuery = query.Encode()
	return parsed.String()
}

func assertEventuallyInUse(t *testing.T, db *sql.DB, want int) {
	t.Helper()

	deadline := time.Now().Add(time.Second)
	for {
		if got := db.Stats().InUse; got == want {
			return
		}
		if time.Now().After(deadline) {
			t.Fatalf("database InUse = %d, want %d", db.Stats().InUse, want)
		}
		time.Sleep(10 * time.Millisecond)
	}
}
