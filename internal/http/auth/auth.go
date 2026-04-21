// Package auth provides auth-related HTTP handlers.
package auth

import (
	"errors"
	"net/http"

	"github.com/andrew-hayworth22/critiquefi-service/internal/appcontext"
	authBus "github.com/andrew-hayworth22/critiquefi-service/internal/business/auth"
	"github.com/andrew-hayworth22/critiquefi-service/internal/models"
	"github.com/andrew-hayworth22/critiquefi-service/pkg/httputil"
)

// Handler exposes HTTP endpoints related to authentication
type Handler struct {
	bus                      *authBus.Bus
	refreshTokenCookieName   string
	refreshTokenCookieDomain string
}

// New creates a new auth HTTP handler
func New(bus *authBus.Bus, refreshTokenCookieName, refreshTokenCookieDomain string) *Handler {
	return &Handler{
		bus:                      bus,
		refreshTokenCookieName:   refreshTokenCookieName,
		refreshTokenCookieDomain: refreshTokenCookieDomain,
	}
}

// RegisterRequest represents a request to register a new user
type RegisterRequest struct {
	Email           string `json:"email"`
	DisplayName     string `json:"display_name"`
	Name            string `json:"name"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
	Remember        bool   `json:"remember"`
}

// ToModel converts a RegisterRequest to a NewUserRequest
func (r *RegisterRequest) ToModel() models.NewUserRequest {
	return models.NewUserRequest{
		Email:           r.Email,
		DisplayName:     r.DisplayName,
		Name:            r.Name,
		Password:        r.Password,
		ConfirmPassword: r.ConfirmPassword,
	}
}

// LoginRequest represents a request to login a user
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Remember bool   `json:"remember"`
}

// AuthenticationResponse represents a response to an authentication request
type AuthenticationResponse struct {
	AccessToken string `json:"access_token"`
}

// Register handles a request to register a new user
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	logger := appcontext.GetLogger(r.Context())

	var req RegisterRequest
	if !httputil.DecodeRequest(w, r, &req) {
		return
	}

	accessToken, refreshToken, err := h.bus.Register(r.Context(), req.ToModel(), r.UserAgent(), req.Remember)
	if err != nil {
		switch {
		case errors.As(err, &models.ValidationErrors{}):

			httputil.WriteUnprocessable(w, err)
			return
		case errors.Is(err, authBus.ErrDuplicate):
			httputil.WriteConflict(w)
			return
		default:
			logger.Error("registration failed", "error", err)
			httputil.WriteInternalError(w)
			return
		}
	}

	if refreshToken != "" {
		h.setRefreshTokenCookie(w, refreshToken)
	}

	httputil.WriteJSON(w, http.StatusCreated, AuthenticationResponse{
		AccessToken: accessToken,
	})
}

// Login handles a request to login a user
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	logger := appcontext.GetLogger(r.Context())

	var req LoginRequest
	if !httputil.DecodeRequest(w, r, &req) {
		return
	}

	accessToken, refreshToken, err := h.bus.Login(r.Context(), req.Email, req.Password, r.UserAgent(), req.Remember)
	if err != nil {
		switch {
		case errors.Is(err, authBus.ErrInvalidCredentials):
			httputil.WriteUnauthorized(w)
			return
		default:
			logger.Error("login failed", "error", err)
			httputil.WriteInternalError(w)
			return
		}
	}

	if refreshToken != "" {
		h.setRefreshTokenCookie(w, refreshToken)
	}

	httputil.WriteJSON(w, http.StatusOK, AuthenticationResponse{
		AccessToken: accessToken,
	})
}

// Logout handles a request to logout a user
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	logger := appcontext.GetLogger(r.Context())

	cookie, err := r.Cookie(h.refreshTokenCookieName)
	if err != nil {
		httputil.WriteUnauthorized(w)
		return
	}

	if err := h.bus.Logout(r.Context(), cookie.Value); err != nil {
		logger.Error("logout failed", "error", err)
		httputil.WriteInternalError(w)
		return
	}

	h.clearRefreshTokenCookie(w)
	httputil.WriteNoContent(w)
}

// Refresh handles a request to refresh an access token
func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(h.refreshTokenCookieName)
	if err != nil {
		httputil.WriteUnauthorized(w)
		return
	}

	accessToken, refreshToken, err := h.bus.Refresh(r.Context(), cookie.Value)
	if err != nil {
		httputil.WriteUnauthorized(w)
		return
	}

	if refreshToken != "" {
		h.setRefreshTokenCookie(w, refreshToken)
	}

	httputil.WriteJSON(w, http.StatusOK, AuthenticationResponse{
		AccessToken: accessToken,
	})
}

// setRefreshTokenCookie sets the refresh token cookie
func (h *Handler) setRefreshTokenCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     h.refreshTokenCookieName,
		Value:    token,
		Domain:   h.refreshTokenCookieDomain,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})
}

// clearRefreshTokenCookie clears the refresh token cookie
func (h *Handler) clearRefreshTokenCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     h.refreshTokenCookieName,
		Domain:   h.refreshTokenCookieDomain,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1,
	})
}
