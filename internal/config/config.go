// Package config imports environment variables to configure the application
package config

import (
	"errors"
	"os"
	"strings"
	"time"
)

// Config represents all the configuration values for the application
type Config struct {
	Port           string
	DatabaseURL    string
	JWTSecret      string
	AccessTokenTTL time.Duration
	CORSOrigins    []string
	CookieDomain   string
}

// Load gets the configurations values from the application environment
func Load() (Config, error) {
	accessTokenTTL, err := time.ParseDuration(getenv("ACCESS_TOKEN_TTL", "15m"))
	if err != nil {
		return Config{}, err
	}

	cfg := Config{
		Port:           ":" + getenv("PORT", "8080"),
		DatabaseURL:    getenv("DB_URL", ""),
		JWTSecret:      getenv("JWT_SECRET", ""),
		AccessTokenTTL: accessTokenTTL,
		CORSOrigins:    splitCSV(getenv("CORS_ORIGINS", "http://localhost:3000")),
		CookieDomain:   getenv("COOKIE_DOMAIN", "localhost"),
	}

	if cfg.JWTSecret == "" {
		return Config{}, errors.New("JWT secret not set; run \"make genkey\" and set the JWT_KEY environment variable")
	}

	return cfg, nil
}

// getenv fetches a value from the environment and returns the default if a value does not exist
func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

// splitCSV splits comma-separated values
func splitCSV(s string) []string {
	return strings.Split(s, ",")
}
