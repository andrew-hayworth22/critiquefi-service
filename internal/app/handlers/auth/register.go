package auth

import (
	"encoding/json"
	"net/http"
	"net/mail"
	"time"

	"github.com/andrew-hayworth22/critiquefi-service/internal/app/sdk"
	"github.com/andrew-hayworth22/critiquefi-service/internal/store/types"
)

type RegisterRequest struct {
	Email                string `json:"email"`
	DisplayName          string `json:"display_name"`
	Name                 string `json:"name"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"confirm_password"`
}

type RegisterResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

func (r *RegisterRequest) Decode(data []byte) error {
	return json.Unmarshal(data, r)
}

func (r *RegisterRequest) Validate() error {
	reqErr := sdk.NewError(http.StatusBadRequest, "invalid request")

	if _, err := mail.ParseAddress(r.Email); err != nil {
		reqErr.Message = "invalid email address"
		return reqErr
	}

	displayNameRegexpr := `^[A-Za-z0-9._-]+$`
	if err := sdk.ValidateText(r.DisplayName, displayNameRegexpr, 3, 30); err != nil {
		reqErr.Message = "invalid display name: " + err.Error()
		return reqErr
	}

	if err := sdk.ValidateText(r.Name, "", 3, 50); err != nil {
		reqErr.Message = "invalid name: " + err.Error()
		return reqErr
	}

	if err := sdk.ValidateText(r.Password, "", 8, 60); err != nil {
		reqErr.Message = "invalid password: " + err.Error()
		return reqErr
	}

	if r.Password != r.PasswordConfirmation {
		reqErr.Message = "passwords do not match"
		return reqErr
	}

	return nil
}

func (app *AuthApp) Register(w http.ResponseWriter, r *http.Request) {
	var request RegisterRequest
	if err := sdk.Decode(r, &request); err != nil {
		sdk.HandleError(w, err)
		return
	}

	conflictingUser, err := app.db.GetUserByEmail(r.Context(), request.Email)
	if err != nil {
		sdk.HandleError(w, err)
		return
	}

	if conflictingUser != nil {
		sdk.HandleError(w, sdk.NewError(http.StatusConflict, "email already registered"))
		return
	}

	hash, err := sdk.HashPassword(request.Password)
	if err != nil {
		sdk.HandleError(w, err)
		return
	}

	user := types.User{
		Email:        request.Email,
		DisplayName:  request.DisplayName,
		Name:         request.Name,
		PasswordHash: hash,
		IsAdmin:      false,
		LastLogin:    types.NullableTime{},
		IsActive:     false,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    types.NullableTime{},
	}

	id, err := app.db.CreateUser(r.Context(), user)
	if err != nil {
		sdk.HandleError(w, err)
		return
	}
	user.ID = id

	accessToken, err := app.jwt.GenerateToken(&user)
	if err != nil {
		sdk.HandleError(w, err)
		return
	}

	response := RegisterResponse{
		AccessToken: accessToken,
		ExpiresIn:   int64(app.jwt.AccessTokenTTL.Seconds()),
		TokenType:   "bearer",
	}

	_ = sdk.Respond(w, response, http.StatusOK)
}
