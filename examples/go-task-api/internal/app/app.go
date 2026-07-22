package app

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/yokry-he/yok-study-web/examples/go-task-api/internal/config"
	taskdomain "github.com/yokry-he/yok-study-web/examples/go-task-api/internal/task"
	userdomain "github.com/yokry-he/yok-study-web/examples/go-task-api/internal/user"
)

var (
	ErrNilApp             = errors.New("app is nil")
	ErrNilRunContext      = errors.New("app run context is nil")
	ErrInvalidShutdownTTL = errors.New("shutdown timeout must be positive")
)

type App struct {
	Handler   http.Handler
	server    *http.Server
	readiness *Readiness
	db        *sql.DB
	logger    *slog.Logger
}

func New(cfg config.Config, logger *slog.Logger, db *sql.DB) *App {
	usersRepository := userdomain.NewPostgresRepository(db)
	usersService := userdomain.NewService(usersRepository)
	usersHandler := userdomain.NewHandler(usersService)

	tasksRepository := taskdomain.NewPostgresRepository(db)
	tasksService := taskdomain.NewService(tasksRepository, usersRepository)
	tasksHandler := taskdomain.NewHandler(tasksService)

	readiness := NewReadiness(db)
	handler := NewRouter(logger, readiness, usersHandler, tasksHandler, cfg.HTTP.RequestTimeout)
	return &App{
		Handler:   handler,
		server:    NewHTTPServer(cfg.HTTP, handler),
		readiness: readiness,
		db:        db,
		logger:    logger,
	}
}

func (a *App) Run(ctx context.Context, shutdownTimeout time.Duration) error {
	if a == nil || a.server == nil {
		return ErrNilApp
	}
	if ctx == nil {
		return ErrNilRunContext
	}
	if shutdownTimeout <= 0 {
		return ErrInvalidShutdownTTL
	}
	if a.logger != nil {
		a.logger.Info("HTTP 服务开始监听", "address", a.server.Addr)
	}

	serveResult := make(chan error, 1)
	go func() {
		serveResult <- a.server.ListenAndServe()
	}()

	select {
	case err := <-serveResult:
		return errors.Join(normalizeServerError(err), a.closeDatabase())
	case <-ctx.Done():
		if a.logger != nil {
			a.logger.Info("HTTP 服务开始关闭", "timeout", shutdownTimeout)
		}
		if a.readiness != nil {
			a.readiness.StartShutdown()
		}
		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		shutdownErr := a.server.Shutdown(shutdownCtx)
		cancel()
		if shutdownErr != nil {
			shutdownErr = errors.Join(shutdownErr, a.server.Close())
		}
		serveErr := <-serveResult
		result := errors.Join(shutdownErr, normalizeServerError(serveErr), a.closeDatabase())
		if a.logger != nil {
			a.logger.Info("HTTP 服务关闭完成", "success", result == nil)
		}
		return result
	}
}

func normalizeServerError(err error) error {
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}

func (a *App) closeDatabase() error {
	if a == nil || a.db == nil {
		return nil
	}
	return a.db.Close()
}
