package sys

import (
	"github.com/andrew-hayworth22/critiquefi-service/internal/app/sdk"
	"net/http"
)

func (app *App) liveness(w http.ResponseWriter, r *http.Request) {
	response := struct {
		Status string `json:"status"`
	}{
		"ok",
	}

	if err := sdk.Respond(w, response, http.StatusOK); err != nil {
		sdk.HandleError(w, err)
	}
}

func (app *App) readiness(w http.ResponseWriter, r *http.Request) {
	response := struct {
		Status string `json:"status"`
	}{
		"ok",
	}

	if err := app.repo.Sys().Ping(); err != nil {
		sdk.HandleError(w, err)
		return
	}

	if err := sdk.Respond(w, response, http.StatusOK); err != nil {
		sdk.HandleError(w, err)
		return
	}
}
