package media

import (
	"net/http"

	"github.com/andrew-hayworth22/critiquefi-service/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
)

type App struct {
	jwtHandler *jwtauth.JWTAuth
	db         store.MediaStore
}

func NewApp(jwtHandler *jwtauth.JWTAuth, db store.MediaStore) *App {
	return &App{jwtHandler: jwtHandler, db: db}
}

func (a *App) Router() *chi.Mux {
	r := chi.NewRouter()

	// protected routes
	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Authenticator(a.jwtHandler))

		r.Post("/film", a.CreateFilm)
	})

	// public routes
	r.Group(func(r chi.Router) {
		r.Get("/", a.GetMedia)
		r.Get("/film/{id}", a.GetFilm)
	})

	return r
}

func (a *App) GetMedia(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}
