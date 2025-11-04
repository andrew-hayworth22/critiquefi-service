package sys

import (
	"github.com/andrew-hayworth22/critiquefi-service/internal/store"
	"github.com/go-chi/chi/v5"
)

type App struct {
	repo store.Repo
}

func NewSysApp(r store.Repo) *App {
	return &App{
		repo: r,
	}
}

func (app *App) Router() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/liveness", app.liveness)
	r.Get("/readiness", app.readiness)
	return r
}
