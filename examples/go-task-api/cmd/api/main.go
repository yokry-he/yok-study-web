package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/yokry-he/yok-study-web/examples/go-task-api/internal/app"
	"github.com/yokry-he/yok-study-web/examples/go-task-api/internal/config"
	"github.com/yokry-he/yok-study-web/examples/go-task-api/internal/platform/database"
)

func main() {
	if err := run(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "API 启动失败：%v\n", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: cfg.LogLevel}))
	db, err := database.Open(ctx, cfg.Database)
	if err != nil {
		return err
	}
	application := app.New(cfg, logger, db)
	return application.Run(ctx, cfg.HTTP.ShutdownTimeout)
}
