package main

import (
	"errors"
	"github.com/andrew-hayworth22/critiquefi-service/internal/app"
	"github.com/andrew-hayworth22/critiquefi-service/internal/config"
	"github.com/andrew-hayworth22/critiquefi-service/internal/store/postgres"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"log"
	"net/http"
)

func main() {
	// Load environment
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading .env file: %v", err)
	}

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

	// Create application
	a := app.NewApp(db.Repo())

	err = http.ListenAndServe(cfg.Port, a.Server)
	if err != nil {
		log.Fatalf("error starting server: %v", err)
	}
}
