package auth

import (
	"encoding/json"
	"net/http"

	"github.com/andrew-hayworth22/critiquefi-service/internal/app/sdk"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Remember bool   `json:"remember"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
}

func (l *LoginRequest) Decode(data []byte) error {
	return json.Unmarshal(data, l)
}

func (app *App) Login(w http.ResponseWriter, r *http.Request) {
	var request LoginRequest
	if err := sdk.Decode(r, &request); err != nil {
		sdk.HandleError(w, err)
		return
	}

	user, err := app.db.GetUserByEmail(r.Context(), request.Email)
	if err != nil {
		sdk.HandleError(w, err)
		return
	}

	if user == nil {
		sdk.HandleError(w, sdk.NewError(http.StatusUnauthorized, "invalid credentials"))
		return
	}

	if sdk.ComparePasswords(user.PasswordHash, request.Password) != nil {
		sdk.HandleError(w, sdk.NewError(http.StatusUnauthorized, "invalid credentials"))
		return
	}

	app.issueTokensAndRespond(r.Context(), w, user, r.UserAgent(), request.Remember)
}
