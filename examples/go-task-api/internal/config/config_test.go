package config

import (
	"errors"
	"log/slog"
	"os"
	"strings"
	"testing"
	"time"
)

const testDatabaseURL = "postgres://app:super-secret@localhost:5432/taskdb?sslmode=disable"

var configEnvironmentVariables = []string{
	"APP_ENV",
	"HTTP_ADDR",
	"HTTP_READ_HEADER_TIMEOUT",
	"HTTP_READ_TIMEOUT",
	"HTTP_WRITE_TIMEOUT",
	"HTTP_IDLE_TIMEOUT",
	"HTTP_REQUEST_TIMEOUT",
	"HTTP_SHUTDOWN_TIMEOUT",
	"DATABASE_URL",
	"DB_MAX_OPEN_CONNS",
	"DB_MAX_IDLE_CONNS",
	"DB_CONN_MAX_LIFETIME",
	"DB_CONN_MAX_IDLE_TIME",
	"LOG_LEVEL",
}

func TestLoadUsesDefaults(t *testing.T) {
	resetConfigEnvironment(t)
	t.Setenv("DATABASE_URL", testDatabaseURL)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() returned an unexpected error: %v", err)
	}

	want := Config{
		Environment: "development",
		LogLevel:    slog.LevelInfo,
		HTTP: HTTPConfig{
			Addr:              ":8080",
			ReadHeaderTimeout: 5 * time.Second,
			ReadTimeout:       10 * time.Second,
			WriteTimeout:      15 * time.Second,
			IdleTimeout:       60 * time.Second,
			RequestTimeout:    10 * time.Second,
			ShutdownTimeout:   15 * time.Second,
		},
		Database: DatabaseConfig{
			URL:             testDatabaseURL,
			MaxOpenConns:    20,
			MaxIdleConns:    10,
			ConnMaxLifetime: 30 * time.Minute,
			ConnMaxIdleTime: 5 * time.Minute,
		},
	}
	assertConfigEqual(t, cfg, want)
}

func TestLoadRejectsMissingDatabaseURL(t *testing.T) {
	resetConfigEnvironment(t)

	_, err := Load()
	if !errors.Is(err, ErrDatabaseURLRequired) {
		t.Fatalf("Load() error does not match ErrDatabaseURLRequired")
	}
}

func TestLoadRejectsBlankDatabaseURL(t *testing.T) {
	resetConfigEnvironment(t)
	t.Setenv("DATABASE_URL", " \t\n")

	_, err := Load()
	if !errors.Is(err, ErrDatabaseURLRequired) {
		t.Fatalf("Load() error does not match ErrDatabaseURLRequired")
	}
}

func TestLoadRejectsInvalidValues(t *testing.T) {
	tests := []struct {
		name        string
		field       string
		value       string
		wantMessage string
	}{
		{
			name:        "non-integer max open connections",
			field:       "DB_MAX_OPEN_CONNS",
			value:       "many",
			wantMessage: "invalid DB_MAX_OPEN_CONNS: must be a positive integer",
		},
		{
			name:        "zero max open connections",
			field:       "DB_MAX_OPEN_CONNS",
			value:       "0",
			wantMessage: "invalid DB_MAX_OPEN_CONNS: must be a positive integer",
		},
		{
			name:        "negative max idle connections",
			field:       "DB_MAX_IDLE_CONNS",
			value:       "-1",
			wantMessage: "invalid DB_MAX_IDLE_CONNS: must be a non-negative integer",
		},
		{
			name:        "max idle exceeds max open",
			field:       "DB_MAX_IDLE_CONNS",
			value:       "21",
			wantMessage: "invalid DB_MAX_IDLE_CONNS: must not exceed DB_MAX_OPEN_CONNS",
		},
		{
			name:        "invalid duration",
			field:       "HTTP_READ_TIMEOUT",
			value:       "soon",
			wantMessage: "invalid HTTP_READ_TIMEOUT: must be a positive duration",
		},
		{
			name:        "zero read header timeout",
			field:       "HTTP_READ_HEADER_TIMEOUT",
			value:       "0s",
			wantMessage: "invalid HTTP_READ_HEADER_TIMEOUT: must be a positive duration",
		},
		{
			name:        "zero read timeout",
			field:       "HTTP_READ_TIMEOUT",
			value:       "0s",
			wantMessage: "invalid HTTP_READ_TIMEOUT: must be a positive duration",
		},
		{
			name:        "zero write timeout",
			field:       "HTTP_WRITE_TIMEOUT",
			value:       "0s",
			wantMessage: "invalid HTTP_WRITE_TIMEOUT: must be a positive duration",
		},
		{
			name:        "zero idle timeout",
			field:       "HTTP_IDLE_TIMEOUT",
			value:       "0s",
			wantMessage: "invalid HTTP_IDLE_TIMEOUT: must be a positive duration",
		},
		{
			name:        "zero request timeout",
			field:       "HTTP_REQUEST_TIMEOUT",
			value:       "0s",
			wantMessage: "invalid HTTP_REQUEST_TIMEOUT: must be a positive duration",
		},
		{
			name:        "zero shutdown timeout",
			field:       "HTTP_SHUTDOWN_TIMEOUT",
			value:       "0s",
			wantMessage: "invalid HTTP_SHUTDOWN_TIMEOUT: must be a positive duration",
		},
		{
			name:        "zero connection max lifetime",
			field:       "DB_CONN_MAX_LIFETIME",
			value:       "0s",
			wantMessage: "invalid DB_CONN_MAX_LIFETIME: must be a positive duration",
		},
		{
			name:        "zero connection max idle time",
			field:       "DB_CONN_MAX_IDLE_TIME",
			value:       "0s",
			wantMessage: "invalid DB_CONN_MAX_IDLE_TIME: must be a positive duration",
		},
		{
			name:        "invalid log level",
			field:       "LOG_LEVEL",
			value:       "verbose",
			wantMessage: "invalid LOG_LEVEL: must be one of debug, info, warn, error",
		},
		{
			name:        "blank environment",
			field:       "APP_ENV",
			value:       " \t",
			wantMessage: "invalid APP_ENV: must not be blank",
		},
		{
			name:        "blank HTTP address",
			field:       "HTTP_ADDR",
			value:       " \n",
			wantMessage: "invalid HTTP_ADDR: must not be blank",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			resetConfigEnvironment(t)
			t.Setenv("DATABASE_URL", testDatabaseURL)
			t.Setenv(tc.field, tc.value)

			_, err := Load()
			if err == nil {
				t.Fatalf("Load() unexpectedly succeeded for %s", tc.field)
			}
			if err.Error() != tc.wantMessage {
				t.Errorf("Load() returned the wrong reason for %s", tc.field)
			}
			if strings.Contains(err.Error(), testDatabaseURL) || strings.Contains(err.Error(), "super-secret") {
				t.Errorf("Load() error exposed DATABASE_URL credentials")
			}
		})
	}
}

func TestLoadAcceptsAndNormalizesOverrides(t *testing.T) {
	resetConfigEnvironment(t)
	overrides := map[string]string{
		"APP_ENV":                  " staging ",
		"HTTP_ADDR":                " 127.0.0.1:9090 ",
		"HTTP_READ_HEADER_TIMEOUT": "1s",
		"HTTP_READ_TIMEOUT":        "2s",
		"HTTP_WRITE_TIMEOUT":       "3s",
		"HTTP_IDLE_TIMEOUT":        "4s",
		"HTTP_REQUEST_TIMEOUT":     "5s",
		"HTTP_SHUTDOWN_TIMEOUT":    "6s",
		"DATABASE_URL":             " " + testDatabaseURL + " ",
		"DB_MAX_OPEN_CONNS":        "40",
		"DB_MAX_IDLE_CONNS":        "12",
		"DB_CONN_MAX_LIFETIME":     "2h",
		"DB_CONN_MAX_IDLE_TIME":    "45m",
		"LOG_LEVEL":                " WaRn ",
	}
	for field, value := range overrides {
		t.Setenv(field, value)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() returned an unexpected error: %v", err)
	}

	want := Config{
		Environment: "staging",
		LogLevel:    slog.LevelWarn,
		HTTP: HTTPConfig{
			Addr:              "127.0.0.1:9090",
			ReadHeaderTimeout: time.Second,
			ReadTimeout:       2 * time.Second,
			WriteTimeout:      3 * time.Second,
			IdleTimeout:       4 * time.Second,
			RequestTimeout:    5 * time.Second,
			ShutdownTimeout:   6 * time.Second,
		},
		Database: DatabaseConfig{
			URL:             testDatabaseURL,
			MaxOpenConns:    40,
			MaxIdleConns:    12,
			ConnMaxLifetime: 2 * time.Hour,
			ConnMaxIdleTime: 45 * time.Minute,
		},
	}
	assertConfigEqual(t, cfg, want)
}

func assertConfigEqual(t *testing.T, got, want Config) {
	t.Helper()

	if got.Environment != want.Environment {
		t.Errorf("Environment = %q, want %q", got.Environment, want.Environment)
	}
	if got.LogLevel != want.LogLevel {
		t.Errorf("LogLevel = %s, want %s", got.LogLevel, want.LogLevel)
	}
	if got.HTTP != want.HTTP {
		t.Errorf("HTTP = %+v, want %+v", got.HTTP, want.HTTP)
	}
	if got.Database.URL != want.Database.URL {
		t.Errorf("Database.URL does not match the expected value")
	}
	gotDatabase := got.Database
	wantDatabase := want.Database
	gotDatabase.URL = ""
	wantDatabase.URL = ""
	if gotDatabase != wantDatabase {
		t.Errorf("Database settings = %+v, want %+v", gotDatabase, wantDatabase)
	}
}

func resetConfigEnvironment(t *testing.T) {
	t.Helper()

	for _, name := range configEnvironmentVariables {
		value, exists := os.LookupEnv(name)
		if err := os.Unsetenv(name); err != nil {
			t.Fatalf("unset %s: %v", name, err)
		}
		t.Cleanup(func() {
			var err error
			if exists {
				err = os.Setenv(name, value)
			} else {
				err = os.Unsetenv(name)
			}
			if err != nil {
				t.Errorf("restore %s: %v", name, err)
			}
		})
	}
}
