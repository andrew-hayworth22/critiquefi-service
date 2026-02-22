package app

import (
	"net/http"
	"time"

	"github.com/andrew-hayworth22/critiquefi-service/internal/app/handlers/auth"
	"github.com/andrew-hayworth22/critiquefi-service/internal/app/handlers/media"
	"github.com/andrew-hayworth22/critiquefi-service/internal/app/handlers/sys"
	"github.com/andrew-hayworth22/critiquefi-service/internal/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
)

type App struct {
	cfg    config.Config
	Server *http.Server
}

func NewApp(
	jwtHandler *jwtauth.JWTAuth,
	sys *sys.App,
	auth *auth.App,
	media *media.App,
) chi.Router {

	r := chi.NewRouter()

	r.Use(
		middleware.RequestID,
		middleware.RealIP,
		middleware.Logger,
		middleware.Recoverer,
		middleware.Timeout(60*time.Second),
		jwtauth.Verifier(jwtHandler),
	)

	r.Mount("/api/sys", sys.Router())
	r.Mount("/api/auth", auth.Router())
	r.Mount("/api/media", media.Router())

	return r
}
