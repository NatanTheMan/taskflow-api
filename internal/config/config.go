// Package config loads .env vars
package config

import (
	"errors"
	"fmt"
	"os"
	"time"
)

type Config struct {
	Env             string
	HTTPPort        string
	LogLevel        string
	JwtSecret       string
	DatabaseURL     string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
	ShutdownTimeout time.Duration
}

func Load() (*Config, error) {
	var errs []error

	accesTTL, err := durationEnv("ACCESS_TOKEN_TTL", 15*time.Minute)
	errs = append(errs, err)
	refreshTTL, err := durationEnv("REFRESH_TOKEN_TTL", 168*time.Hour)
	errs = append(errs, err)
	shutdownTimeout, err := durationEnv("SHUTDOWN_TIMEOUT", 10*time.Second)
	errs = append(errs, err)

	cfg := &Config{
		Env:             stringEnv("APP_ENV", "development"),
		HTTPPort:        stringEnv("HTTP_PORT", "8080"),
		LogLevel:        stringEnv("LOG_LEVEL", "info"),
		JwtSecret:       os.Getenv("JWT_SECRET"),
		DatabaseURL:     os.Getenv("DATABASE_URL"),
		AccessTokenTTL:  accesTTL,
		RefreshTokenTTL: refreshTTL,
		ShutdownTimeout: shutdownTimeout,
	}

	if len(cfg.JwtSecret) < 32 {
		errs = append(errs, errors.New("JWT_SECRET must have at least 32 characters"))
	}

	if cfg.DatabaseURL == "" {
		errs = append(errs, errors.New("DATABASE_URL is required"))
	}

	if err := errors.Join(errs...); err != nil {
		return nil, err
	}

	return cfg, nil
}

func stringEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func durationEnv(key string, def time.Duration) (time.Duration, error) {
	v := os.Getenv(key)

	if v == "" {
		return def, nil
	}

	d, err := time.ParseDuration(v)
	if err != nil {
		return 0, fmt.Errorf("%s: invalid duration %q", key, v)
	}

	return d, nil
}
