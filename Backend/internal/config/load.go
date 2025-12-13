package config

import (
	"errors"
	"os"
	"strconv"
	"time"
)

func Load() (Config, error) {
	cfg := Config{
		Env:             getenv("APP_ENV", "dev"),
		ServerHost:      getenv("SERVER_HOST", "0.0.0.0"),
		ServerPort:      getenv("SERVER_PORT", "8080"),
		DatabaseURL:     getenv("DATABASE_URL", ""),
		JWTSecret:       getenv("JWT_SECRET", ""),
		JWTIssuer:       getenv("JWT_ISSUER", "pulsementor"),
		AccessTokenTTL:  mustDuration(getenv("ACCESS_TOKEN_TTL", "24h")),
		RefreshTokenTTL: mustDuration(getenv("REFRESH_TOKEN_TTL", "168h")), // 7 days default
		BcryptCost:      mustInt(getenv("BCRYPT_COST", "12")),
	}

	if cfg.DatabaseURL == "" {
		return Config{}, errors.New("DATABASE_URL is required")
	}
	if cfg.JWTSecret == "" {
		return Config{}, errors.New("JWT_SECRET is required")
	}
	if cfg.BcryptCost < 10 || cfg.BcryptCost > 15 {
		return Config{}, errors.New("BCRYPT_COST must be between 10 and 15")
	}
	return cfg, nil
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func mustDuration(v string) time.Duration {
	d, err := time.ParseDuration(v)
	if err != nil {
		return 24 * time.Hour
	}
	return d
}

func mustInt(v string) int {
	i, err := strconv.Atoi(v)
	if err != nil {
		return 12
	}
	return i
}
