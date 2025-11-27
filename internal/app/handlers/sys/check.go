package sys

import (
	"net/http"

	"github.com/andrew-hayworth22/critiquefi-service/internal/app/sdk"
)

func (app *SysApp) liveness(w http.ResponseWriter, r *http.Request) {
	response := struct {
		Status string `json:"status"`
	}{
		"ok",
	}

	if err := sdk.Respond(w, response, http.StatusOK); err != nil {
		sdk.HandleError(w, err)
	}
}

func (app *SysApp) readiness(w http.ResponseWriter, r *http.Request) {
	response := struct {
		Status string `json:"status"`
	}{
		"ok",
	}

	if err := app.db.Ping(r.Context()); err != nil {
		sdk.HandleError(w, err)
		return
	}

	if err := sdk.Respond(w, response, http.StatusOK); err != nil {
		sdk.HandleError(w, err)
		return
	}
}
