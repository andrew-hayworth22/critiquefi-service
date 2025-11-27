package sys

import (
	"github.com/andrew-hayworth22/critiquefi-service/internal/store"
	"github.com/go-chi/chi/v5"
)

type SysApp struct {
	db store.SysStore
}

func NewSysApp(db store.SysStore) *SysApp {
	return &SysApp{
		db: db,
	}
}

func (app *SysApp) Router() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/liveness", app.liveness)
	r.Get("/readiness", app.readiness)
	return r
}
