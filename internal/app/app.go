package app

import (
	"github.com/andrew-hayworth22/critiquefi-service/internal/app/handlers/auth"
	"github.com/andrew-hayworth22/critiquefi-service/internal/app/handlers/sys"
	"github.com/andrew-hayworth22/critiquefi-service/internal/app/sdk"
	"github.com/andrew-hayworth22/critiquefi-service/internal/store"
	"github.com/go-chi/chi/v5"
)

type App struct {
	Server chi.Router
	db     store.Repo
}

func NewApp(db store.Repo, jwt *sdk.JWTManager) *App {
	r := chi.NewRouter()

	sysApp := sys.NewSysApp(db)
	authApp := auth.NewAuthApp(db, jwt)

	r.Mount("/sys", sysApp.Router())
	r.Mount("/auth", authApp.Router())

	return &App{
		db:     db,
		Server: r,
	}
}
