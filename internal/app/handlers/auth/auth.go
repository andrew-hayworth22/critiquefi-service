package auth

import (
	"github.com/andrew-hayworth22/critiquefi-service/internal/app/sdk"
	"github.com/andrew-hayworth22/critiquefi-service/internal/store"
	"github.com/go-chi/chi/v5"
)

type AuthApp struct {
	db  store.AuthStore
	jwt *sdk.JWTManager
}

func NewAuthApp(db store.AuthStore, jwt *sdk.JWTManager) *AuthApp {
	return &AuthApp{
		db:  db,
		jwt: jwt,
	}
}

func (app *AuthApp) Router() *chi.Mux {
	r := chi.NewRouter()
	r.Post("/register", app.Register)
	return r
}
