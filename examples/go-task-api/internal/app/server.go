package app

import (
	"net/http"

	"github.com/yokry-he/yok-study-web/examples/go-task-api/internal/config"
)

func NewHTTPServer(cfg config.HTTPConfig, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:              cfg.Addr,
		Handler:           handler,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
		ReadTimeout:       cfg.ReadTimeout,
		WriteTimeout:      cfg.WriteTimeout,
		IdleTimeout:       cfg.IdleTimeout,
	}
}
