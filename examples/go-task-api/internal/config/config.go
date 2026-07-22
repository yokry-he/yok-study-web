package config

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	defaultEnvironment       = "development"
	defaultHTTPAddr          = ":8080"
	defaultReadHeaderTimeout = "5s"
	defaultReadTimeout       = "10s"
	defaultWriteTimeout      = "15s"
	defaultIdleTimeout       = "60s"
	defaultRequestTimeout    = "10s"
	defaultShutdownTimeout   = "15s"
	defaultMaxOpenConns      = "20"
	defaultMaxIdleConns      = "10"
	defaultConnMaxLifetime   = "30m"
	defaultConnMaxIdleTime   = "5m"
	defaultLogLevel          = "info"
)

var ErrDatabaseURLRequired = errors.New("DATABASE_URL is required")

type Config struct {
	Environment string
	LogLevel    slog.Level
	HTTP        HTTPConfig
	Database    DatabaseConfig
}

type HTTPConfig struct {
	Addr              string
	ReadHeaderTimeout time.Duration
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	RequestTimeout    time.Duration
	ShutdownTimeout   time.Duration
}

type DatabaseConfig struct {
	URL             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

func Load() (Config, error) {
	databaseURL, exists := os.LookupEnv("DATABASE_URL")
	databaseURL = strings.TrimSpace(databaseURL)
	if !exists || databaseURL == "" {
		return Config{}, ErrDatabaseURLRequired
	}

	environment := strings.TrimSpace(envOrDefault("APP_ENV", defaultEnvironment))
	if environment == "" {
		return Config{}, errors.New("invalid APP_ENV: must not be blank")
	}

	httpAddr := strings.TrimSpace(envOrDefault("HTTP_ADDR", defaultHTTPAddr))
	if httpAddr == "" {
		return Config{}, errors.New("invalid HTTP_ADDR: must not be blank")
	}

	readHeaderTimeout, err := positiveDuration("HTTP_READ_HEADER_TIMEOUT", defaultReadHeaderTimeout)
	if err != nil {
		return Config{}, err
	}
	readTimeout, err := positiveDuration("HTTP_READ_TIMEOUT", defaultReadTimeout)
	if err != nil {
		return Config{}, err
	}
	writeTimeout, err := positiveDuration("HTTP_WRITE_TIMEOUT", defaultWriteTimeout)
	if err != nil {
		return Config{}, err
	}
	idleTimeout, err := positiveDuration("HTTP_IDLE_TIMEOUT", defaultIdleTimeout)
	if err != nil {
		return Config{}, err
	}
	requestTimeout, err := positiveDuration("HTTP_REQUEST_TIMEOUT", defaultRequestTimeout)
	if err != nil {
		return Config{}, err
	}
	shutdownTimeout, err := positiveDuration("HTTP_SHUTDOWN_TIMEOUT", defaultShutdownTimeout)
	if err != nil {
		return Config{}, err
	}

	maxOpenConns, err := positiveInteger("DB_MAX_OPEN_CONNS", defaultMaxOpenConns)
	if err != nil {
		return Config{}, err
	}
	maxIdleConns, err := nonNegativeInteger("DB_MAX_IDLE_CONNS", defaultMaxIdleConns)
	if err != nil {
		return Config{}, err
	}
	if maxIdleConns > maxOpenConns {
		return Config{}, errors.New("invalid DB_MAX_IDLE_CONNS: must not exceed DB_MAX_OPEN_CONNS")
	}

	connMaxLifetime, err := positiveDuration("DB_CONN_MAX_LIFETIME", defaultConnMaxLifetime)
	if err != nil {
		return Config{}, err
	}
	connMaxIdleTime, err := positiveDuration("DB_CONN_MAX_IDLE_TIME", defaultConnMaxIdleTime)
	if err != nil {
		return Config{}, err
	}

	logLevel, err := parseLogLevel()
	if err != nil {
		return Config{}, err
	}

	return Config{
		Environment: environment,
		LogLevel:    logLevel,
		HTTP: HTTPConfig{
			Addr:              httpAddr,
			ReadHeaderTimeout: readHeaderTimeout,
			ReadTimeout:       readTimeout,
			WriteTimeout:      writeTimeout,
			IdleTimeout:       idleTimeout,
			RequestTimeout:    requestTimeout,
			ShutdownTimeout:   shutdownTimeout,
		},
		Database: DatabaseConfig{
			URL:             databaseURL,
			MaxOpenConns:    maxOpenConns,
			MaxIdleConns:    maxIdleConns,
			ConnMaxLifetime: connMaxLifetime,
			ConnMaxIdleTime: connMaxIdleTime,
		},
	}, nil
}

func envOrDefault(name, fallback string) string {
	value, exists := os.LookupEnv(name)
	if !exists {
		return fallback
	}
	return value
}

func positiveInteger(name, fallback string) (int, error) {
	value, err := strconv.Atoi(strings.TrimSpace(envOrDefault(name, fallback)))
	if err != nil || value <= 0 {
		return 0, fmt.Errorf("invalid %s: must be a positive integer", name)
	}
	return value, nil
}

func nonNegativeInteger(name, fallback string) (int, error) {
	value, err := strconv.Atoi(strings.TrimSpace(envOrDefault(name, fallback)))
	if err != nil || value < 0 {
		return 0, fmt.Errorf("invalid %s: must be a non-negative integer", name)
	}
	return value, nil
}

func positiveDuration(name, fallback string) (time.Duration, error) {
	value, err := time.ParseDuration(strings.TrimSpace(envOrDefault(name, fallback)))
	if err != nil || value <= 0 {
		return 0, fmt.Errorf("invalid %s: must be a positive duration", name)
	}
	return value, nil
}

func parseLogLevel() (slog.Level, error) {
	switch strings.ToLower(strings.TrimSpace(envOrDefault("LOG_LEVEL", defaultLogLevel))) {
	case "debug":
		return slog.LevelDebug, nil
	case "info":
		return slog.LevelInfo, nil
	case "warn":
		return slog.LevelWarn, nil
	case "error":
		return slog.LevelError, nil
	default:
		return 0, errors.New("invalid LOG_LEVEL: must be one of debug, info, warn, error")
	}
}
