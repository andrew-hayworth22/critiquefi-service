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
	err := sdk.Respond(w, response, http.StatusOK)
	if err != nil {
		sdk.HandleError(w, err)
	}
}
