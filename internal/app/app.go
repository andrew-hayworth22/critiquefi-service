package app

import (
	"github.com/andrew-hayworth22/critiquefi-service/internal/app/handlers/sys"
	"github.com/andrew-hayworth22/critiquefi-service/internal/store"
	"github.com/go-chi/chi/v5"
)

type App struct {
	Server chi.Router
	db     store.Repo
}

func NewApp(db store.Repo) *App {
	r := chi.NewRouter()

	sysApp := sys.NewSysApp()

	r.Mount("/", sysApp.Router())

	return &App{
		db:     db,
		Server: r,
	}
}
