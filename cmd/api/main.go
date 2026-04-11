package main

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/andrew-hayworth22/critiquefi-service/internal/app"
	"github.com/andrew-hayworth22/critiquefi-service/internal/app/handlers/auth"
	"github.com/andrew-hayworth22/critiquefi-service/internal/app/handlers/media"
	"github.com/andrew-hayworth22/critiquefi-service/internal/app/handlers/sys"
	"github.com/andrew-hayworth22/critiquefi-service/internal/app/sdk"
	"github.com/andrew-hayworth22/critiquefi-service/internal/config"
	"github.com/andrew-hayworth22/critiquefi-service/internal/store/postgres"
	"github.com/go-chi/jwtauth/v5"
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

	// Create JWT packages
	jwtManager := sdk.NewJWTManager(cfg.JWTSecret, cfg.AccessTokenTTL)
	jwtHandler := jwtauth.New("HS256", []byte(cfg.JWTSecret), nil)

	// Build storage packages
	sysDb := postgres.NewSysPG(db)
	authDb := postgres.NewAuthPG(db)
	mediaDb := postgres.NewMediaPG(db)

	// Build application packages
	authApp := auth.NewApp(authDb, jwtManager, jwtHandler, cfg.RefreshTokenTTL, cfg.RefreshTokenCookieName, cfg.RefreshTokenCookieDomain)
	sysApp := sys.NewApp(sysDb)
	mediaApp := media.NewApp(jwtHandler, mediaDb)

	// Create application
	a := app.NewApp(
		jwtHandler,
		sysApp,
		authApp,
		mediaApp,
	)

	err = http.ListenAndServe(cfg.Port, a)
	if err != nil {
		log.Fatalf("error starting server: %v", err)
	}
}
