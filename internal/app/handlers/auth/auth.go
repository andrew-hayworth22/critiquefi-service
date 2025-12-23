package auth

import (
	"context"
	"net/http"
	"time"

	"github.com/andrew-hayworth22/critiquefi-service/internal/app/sdk"
	"github.com/andrew-hayworth22/critiquefi-service/internal/store"
	"github.com/andrew-hayworth22/critiquefi-service/internal/store/types"
	"github.com/go-chi/chi/v5"
)

type App struct {
	db                       store.AuthStore
	jwt                      *sdk.JWTManager
	refreshTokenTTL          time.Duration
	refreshTokenCookieName   string
	refreshTokenCookieDomain string
}

func NewApp(db store.AuthStore, jwt *sdk.JWTManager, refreshTokenTTL time.Duration, refreshTokenCookieName string, refreshTokenCookieDomain string) *App {
	return &App{
		db:                       db,
		jwt:                      jwt,
		refreshTokenTTL:          refreshTokenTTL,
		refreshTokenCookieName:   refreshTokenCookieName,
		refreshTokenCookieDomain: refreshTokenCookieDomain,
	}
}

func (app *App) Router() *chi.Mux {
	r := chi.NewRouter()
	r.Post("/register", app.Register)
	r.Post("/login", app.Login)
	r.Post("/refresh", app.Refresh)
	return r
}

type authenticationResponse struct {
	AccessToken string `json:"access_token"`
}

func (app *App) issueTokensAndRespond(ctx context.Context, w http.ResponseWriter, user *types.User, userAgent string, remember bool) {
	accessToken, err := app.jwt.GenerateToken(user)
	if err != nil {
		sdk.HandleError(w, err)
		return
	}

	response := authenticationResponse{
		AccessToken: accessToken,
	}

	if !remember {
		_ = sdk.Respond(w, response, http.StatusOK)
		return
	}

	token, err := sdk.RandomString(32)
	if err != nil {
		sdk.HandleError(w, err)
		return
	}

	refreshToken := &types.RefreshToken{
		Token:     token,
		UserID:    user.ID,
		UserAgent: userAgent,
		ExpiresAt: time.Now().Add(app.refreshTokenTTL).UTC(),
		CreatedAt: time.Now().UTC(),
	}

	if err := app.db.CreateRefreshToken(ctx, refreshToken); err != nil {
		sdk.HandleError(w, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     app.refreshTokenCookieName,
		Value:    token,
		Expires:  refreshToken.ExpiresAt,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		Domain:   app.refreshTokenCookieDomain,
	})

	_ = sdk.Respond(w, response, http.StatusOK)
}
