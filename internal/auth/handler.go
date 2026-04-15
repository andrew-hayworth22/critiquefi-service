package auth

import (
	"errors"
	"net/http"

	"github.com/andrew-hayworth22/critiquefi-service/internal/appcontext"
	"github.com/andrew-hayworth22/critiquefi-service/internal/models"
	"github.com/andrew-hayworth22/critiquefi-service/pkg/httputil"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

type RegisterRequest struct {
	Email           string `json:"email"`
	DisplayName     string `json:"display_name"`
	Name            string `json:"name"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
	Remember        bool   `json:"remember"`
}

func (r *RegisterRequest) ToModel() models.NewUserRequest {
	return models.NewUserRequest{
		Email:           r.Email,
		DisplayName:     r.DisplayName,
		Name:            r.Name,
		Password:        r.Password,
		ConfirmPassword: r.ConfirmPassword,
	}
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Remember bool   `json:"remember"`
}

type AuthResponse struct {
	AccessToken string `json:"access_token"`
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	logger := appcontext.GetLogger(r.Context())

	var req RegisterRequest
	if !httputil.DecodeRequest(w, r, &req) {
		return
	}

	accessToken, refreshToken, err := h.service.Register(r.Context(), req.ToModel(), r.UserAgent(), req.Remember)
	if err != nil {
		switch {
		case errors.As(err, &models.ValidationErrors{}):

			httputil.WriteUnprocessable(w, err)
			return
		case errors.Is(err, ErrDuplicate):
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

	httputil.WriteJSON(w, http.StatusCreated, AuthResponse{
		AccessToken: accessToken,
	})
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	logger := appcontext.GetLogger(r.Context())

	var req LoginRequest
	if !httputil.DecodeRequest(w, r, &req) {
		return
	}

	accessToken, refreshToken, err := h.service.Login(r.Context(), req.Email, req.Password, r.UserAgent(), req.Remember)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidCredentials):

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

	httputil.WriteJSON(w, http.StatusOK, AuthResponse{
		AccessToken: accessToken,
	})
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	logger := appcontext.GetLogger(r.Context())

	cookie, err := r.Cookie(h.service.refreshTokenCookieName)
	if err != nil {
		httputil.WriteUnauthorized(w)
		return
	}

	if err := h.service.Logout(r.Context(), cookie.Value); err != nil {
		logger.Error("logout failed", "error", err)
		httputil.WriteInternalError(w)
		return
	}

	h.clearRefreshTokenCookie(w)
	httputil.WriteNoContent(w)
}

func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(h.service.refreshTokenCookieName)
	if err != nil {
		httputil.WriteUnauthorized(w)
		return
	}

	accessToken, refreshToken, err := h.service.Refresh(r.Context(), cookie.Value)
	if err != nil {
		httputil.WriteUnauthorized(w)
		return
	}

	if refreshToken != "" {
		h.setRefreshTokenCookie(w, refreshToken)
	}

	httputil.WriteJSON(w, http.StatusOK, AuthResponse{
		AccessToken: accessToken,
	})
}

func (h *Handler) setRefreshTokenCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     h.service.refreshTokenCookieName,
		Value:    token,
		Domain:   h.service.refreshTokenCookieDomain,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})
}

func (h *Handler) clearRefreshTokenCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     h.service.refreshTokenCookieName,
		Domain:   h.service.refreshTokenCookieDomain,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1,
	})
}
