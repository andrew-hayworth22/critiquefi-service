package auth

import (
	"net/http"

	"github.com/andrew-hayworth22/critiquefi-service/internal/app/sdk"
	"github.com/go-chi/jwtauth/v5"
)

func (a *App) Me(w http.ResponseWriter, r *http.Request) {
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		sdk.HandleError(w, sdk.NewError(http.StatusUnauthorized, "invalid session"))
		return
	}

	uid := int64(claims["uid"].(float64))

	if uid <= 0 {
		sdk.HandleError(w, sdk.NewError(http.StatusUnauthorized, "invalid session"))
		return
	}

	user, err := a.db.GetUserByID(r.Context(), uid)
	if err != nil {
		sdk.HandleError(w, err)
		return
	}
	_ = sdk.Respond(w, user, http.StatusOK)
}
