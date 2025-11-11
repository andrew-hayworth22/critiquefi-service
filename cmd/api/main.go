package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/andrew-hayworth22/critiquefi-service/internal/app"
	"github.com/andrew-hayworth22/critiquefi-service/internal/app/sdk"
	"github.com/andrew-hayworth22/critiquefi-service/internal/config"
	"github.com/andrew-hayworth22/critiquefi-service/internal/store/postgres"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	// Build configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	// Establish database connection
	db, err := postgres.Open(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("error connecting to database: %v", err)
	}

	// Run migrations
	m, err := migrate.New("file://migrations", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("error connecting to database for migrations: %v", err)
	}
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("error running migrations: %v", err)
	}

	// Create JWT package
	jwt := sdk.NewJWTManager(cfg.JWTSecret, cfg.AccessTokenTTL)

	// Create application
	a := app.NewApp(db.Repo(), jwt)

	err = http.ListenAndServe(cfg.Port, a.Server)
	if err != nil {
		log.Fatalf("error starting server: %v", err)
	}
}
