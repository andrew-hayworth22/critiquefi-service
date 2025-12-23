package main

import (
	"context"
	"errors"
	"io/fs"
	"log"
	"net/http"

	critiquefi_service "github.com/andrew-hayworth22/critiquefi-service"
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
	authApp := auth.NewApp(authDb, jwt, cfg.RefreshTokenTTL, cfg.RefreshTokenCookieName, cfg.RefreshTokenCookieDomain)
	sysApp := sys.NewSysApp(sysDb)

	// Create application
	a := app.NewApp(
		sysApp,
		authApp,
	)

	// Serve frontend
	fileSystem, err := fs.Sub(critiquefi_service.FrontendFiles, "frontend-svelte/dist")
	if err != nil {
		log.Fatal(err)
	}

	a.HandleFunc("/*", frontendHandler(fileSystem))

	err = http.ListenAndServe(cfg.Port, a)
	if err != nil {
		log.Fatalf("error starting server: %v", err)
	}
}
