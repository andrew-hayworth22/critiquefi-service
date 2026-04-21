// Package sys contains HTTP endpoints related to system checks
package sys

import (
	"net/http"

	sysBus "github.com/andrew-hayworth22/critiquefi-service/internal/business/sys"
	"github.com/andrew-hayworth22/critiquefi-service/pkg/httputil"
)

// Handler exposes HTTP endpoints related to system checks
type Handler struct {
	service *sysBus.Bus
}

// NewHandler creates a new system HTTP handler
func NewHandler(service *sysBus.Bus) *Handler {
	return &Handler{service: service}
}

// SystemCheckResponse represents the response of a system check
type SystemCheckResponse struct {
	Status string `json:"status"`
}

// Liveness checks the liveness of the HTTP server
func (h *Handler) Liveness(w http.ResponseWriter, r *http.Request) {
	httputil.WriteJSON(w, http.StatusOK, SystemCheckResponse{
		Status: "ok",
	})
}

// Readiness checks the readiness of the HTTP server and the database connection
func (h *Handler) Readiness(w http.ResponseWriter, r *http.Request) {
	if err := h.service.Ping(r.Context()); err != nil {
		httputil.WriteInternalError(w)
		return
	}
	httputil.WriteJSON(w, http.StatusOK, SystemCheckResponse{
		Status: "ok",
	})
}
