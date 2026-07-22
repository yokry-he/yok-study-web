package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/yokry-he/yok-study-web/examples/go-task-api/internal/config"
	"github.com/yokry-he/yok-study-web/examples/go-task-api/internal/platform/database"
)

var errUsage = errors.New("用法：migrate up | migrate down 1 | migrate version")

func main() {
	if err := run(os.Args[1:], os.Stdout); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string, output io.Writer) (runErr error) {
	operation, err := parseOperation(args)
	if err != nil {
		return err
	}
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), cfg.HTTP.ShutdownTimeout)
	defer cancel()
	db, err := database.Open(ctx, cfg.Database)
	if err != nil {
		return err
	}
	defer func() { runErr = errors.Join(runErr, db.Close()) }()

	migrator, err := database.NewMigratorContext(ctx, db)
	if err != nil {
		return err
	}
	defer func() { runErr = errors.Join(runErr, migrator.Close()) }()

	switch operation {
	case "up":
		return migrator.Up()
	case "down-one":
		return migrator.DownOne()
	case "version":
		version, dirty, err := migrator.Version()
		if err != nil {
			return err
		}
		_, err = fmt.Fprintf(output, "version=%d dirty=%t\n", version, dirty)
		return err
	default:
		return errUsage
	}
}

func parseOperation(args []string) (string, error) {
	switch {
	case len(args) == 1 && args[0] == "up":
		return "up", nil
	case len(args) == 2 && args[0] == "down" && args[1] == "1":
		return "down-one", nil
	case len(args) == 1 && args[0] == "version":
		return "version", nil
	default:
		return "", errUsage
	}
}
