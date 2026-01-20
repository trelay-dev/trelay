package handler

import (
	"net/http"

	"github.com/aftaab/trelay/internal/api/response"
)

// HealthHandler handles health check requests.
type HealthHandler struct{}

// NewHealthHandler creates a new health handler.
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

type healthResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
}

// Health handles GET /healthz
func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, healthResponse{
		Status:  "ok",
		Version: "1.0.0",
	})
}

// Ready handles GET /readyz
func (h *HealthHandler) Ready(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, healthResponse{
		Status:  "ready",
		Version: "1.0.0",
	})
}
