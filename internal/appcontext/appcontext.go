package appcontext

import (
	"context"
	"log/slog"

	"github.com/andrew-hayworth22/critiquefi-service/internal/models"
)

type contextKey string

const (
	claimsKey contextKey = "claims"
	loggerKey contextKey = "logger"
)

func SetClaims(ctx context.Context, claims models.Claims) context.Context {
	return context.WithValue(ctx, claimsKey, claims)
}

func GetClaims(ctx context.Context) (models.Claims, bool) {
	claims, ok := ctx.Value(claimsKey).(models.Claims)
	return claims, ok
}

func SetLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

func GetLogger(ctx context.Context) *slog.Logger {
	logger, ok := ctx.Value(loggerKey).(*slog.Logger)
	if !ok {
		return slog.Default()
	}
	return logger
}
