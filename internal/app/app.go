package app

import (
	"net/http"
	"time"

	"github.com/andrew-hayworth22/critiquefi-service/internal/app/handlers/auth"
	"github.com/andrew-hayworth22/critiquefi-service/internal/app/handlers/sys"
	"github.com/andrew-hayworth22/critiquefi-service/internal/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type App struct {
	cfg    config.Config
	Server *http.Server
}

func NewApp(
	cfg config.Config,
	sys *sys.SysApp,
	auth *auth.AuthApp,
) chi.Router {

	r := chi.NewRouter()

	r.Use(
		middleware.RequestID,
		middleware.RealIP,
		middleware.Logger,
		middleware.Recoverer,
		middleware.Timeout(60*time.Second),
	)

	r.Mount("/sys", sys.Router())
	r.Mount("/auth", auth.Router())

	return r
}
