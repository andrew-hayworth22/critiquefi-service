package main

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/andrew-hayworth22/critiquefi-service/internal/app"
	"github.com/andrew-hayworth22/critiquefi-service/internal/app/handlers/auth"
	"github.com/andrew-hayworth22/critiquefi-service/internal/app/handlers/sys"
	"github.com/andrew-hayworth22/critiquefi-service/internal/app/sdk"
	"github.com/andrew-hayworth22/critiquefi-service/internal/config"
	"github.com/andrew-hayworth22/critiquefi-service/internal/store/postgres"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	ctx := context.Background()

	// Build configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	// Establish database connection
	db := postgres.NewPool(ctx, cfg.DatabaseURL, cfg.MaxDBConns, cfg.MinDBConns, cfg.MaxConnLifetime, cfg.HealthCheckPeriod)

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

	// Build storage packages
	sysDb := postgres.NewSysPG(db)
	authDb := postgres.NewAuthPG(db)

	// Build application packages
	authApp := auth.NewAuthApp(authDb, jwt)
	sysApp := sys.NewSysApp(sysDb)

	// Create application
	a := app.NewApp(
		cfg,
		sysApp,
		authApp,
	)

	err = http.ListenAndServe(cfg.Port, a)
	if err != nil {
		log.Fatalf("error starting server: %v", err)
	}
}
