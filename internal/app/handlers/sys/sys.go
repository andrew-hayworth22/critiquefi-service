package sys

import "github.com/go-chi/chi/v5"

type App struct{}

func NewSysApp() *App {
	return &App{}
}

func (app *App) Router() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/liveness", app.liveness)
	return r
}
