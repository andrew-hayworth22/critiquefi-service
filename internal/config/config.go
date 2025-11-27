// Package config imports environment variables to configure the application
package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config represents all the configuration values for the application
type Config struct {
	Port              string
	Env               string
	DatabaseURL       string
	MaxDBConns        int32
	MinDBConns        int32
	MaxConnLifetime   time.Duration
	HealthCheckPeriod time.Duration
	JWTSecret         string
	AccessTokenTTL    time.Duration
	CORSOrigins       []string
	CookieDomain      string
}

// Load gets the configuration values from the application environment
func Load() (Config, error) {
	jwtSecret, err := must("JWT_SECRET")
	if err != nil {
		return Config{}, err
	}
	maxDBConns, err := getInt("MAX_DB_CONNS", 5)
	if err != nil {
		return Config{}, err
	}
	minDBConns, err := getInt("MIN_DB_CONNS", 1)
	if err != nil {
		return Config{}, err
	}
	maxConnLifetime, err := getDuration("MAX_CONN_LIFETIME", 30*time.Minute)
	if err != nil {
		return Config{}, err
	}
	healthCheckPeriod, err := getDuration("HEALTH_CHECK_PERIOD", 1*time.Minute)
	if err != nil {
		return Config{}, err
	}
	accessTokenTTL, err := getDuration("ACCESS_TOKEN_TTL", 15*time.Minute)
	if err != nil {
		return Config{}, err
	}

	cfg := Config{
		Port:              ":" + get("PORT", "8080"),
		DatabaseURL:       get("DB_URL", ""),
		MaxDBConns:        int32(maxDBConns),
		MinDBConns:        int32(minDBConns),
		MaxConnLifetime:   maxConnLifetime,
		HealthCheckPeriod: healthCheckPeriod,
		JWTSecret:         jwtSecret,
		AccessTokenTTL:    accessTokenTTL,
		CORSOrigins:       getCSV("CORS_ORIGINS", "http://localhost:3000"),
		CookieDomain:      get("COOKIE_DOMAIN", "localhost"),
	}

	if cfg.JWTSecret == "" {
		return Config{}, errors.New("JWT secret not set; run \"make genkey\" and set the JWT_KEY environment variable")
	}

	return cfg, nil
}

// get fetches a value from the environment and returns the default if a value does not exist
func get(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

// getInt fetches an integer value from the environment and returns the default if a value does not exist
func getInt(key string, def int) (int, error) {
	if v := os.Getenv(key); v != "" {
		return strconv.Atoi(v)
	}
	return def, nil
}

// getDuration fetches a duration value from the environment and returns the default if a value does not exist
func getDuration(key string, def time.Duration) (time.Duration, error) {
	if v := os.Getenv(key); v != "" {
		d, err := time.ParseDuration(v)
		if err != nil {
			return 0, err
		}
		return d, nil
	}
	return def, nil
}

// getCSV gets a comma-separated list of values from the environment
func getCSV(key string, def string) []string {
	s := get(key, def)
	return strings.Split(s, ",")
}

// must fetch a value from the environment and errors if a value does not exist
func must(key string) (string, error) {
	v := os.Getenv(key)
	if v == "" {
		return "", fmt.Errorf("%s is required", key)
	}
	return v, nil
}
