package auth

import (
	"net/http"
	"time"

	"github.com/andrew-hayworth22/critiquefi-service/internal/app/sdk"
)

func (a *App) Refresh(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie(a.refreshTokenCookieName)
	if err != nil {
		sdk.HandleError(w, sdk.NewError(http.StatusUnauthorized, "invalid session"))
		return
	}

	old := c.Value
	oldRefreshToken, err := a.db.GetRefreshToken(r.Context(), old)
	if err != nil || oldRefreshToken == nil || oldRefreshToken.Revoked || time.Now().UTC().After(oldRefreshToken.ExpiresAt) {
		sdk.HandleError(w, sdk.NewError(http.StatusUnauthorized, "invalid session"))
		return
	}

	if err := a.db.RevokeRefreshToken(r.Context(), old); err != nil {
		sdk.HandleError(w, err)
		return
	}

	user, err := a.db.GetUserByID(r.Context(), oldRefreshToken.UserID)
	if err != nil {
		sdk.HandleError(w, err)
		return
	}

	a.issueTokensAndRespond(r.Context(), w, user, oldRefreshToken.UserAgent, true)
}
