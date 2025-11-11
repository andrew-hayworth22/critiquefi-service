package auth

import (
	"github.com/andrew-hayworth22/critiquefi-service/internal/app/sdk"
	"github.com/andrew-hayworth22/critiquefi-service/internal/store"
	"github.com/go-chi/chi/v5"
)

type App struct {
	repo store.Repo
	jwt  *sdk.JWTManager
}

func NewAuthApp(r store.Repo, jwt *sdk.JWTManager) *App {
	return &App{
		repo: r,
		jwt:  jwt,
	}
}

func (app *App) Router() *chi.Mux {
	r := chi.NewRouter()
	r.Post("/register", app.Register)
	return r
}
