package handler

import (
	"net/http"

	"github.com/aftaab/trelay/internal/api/response"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

type healthResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
}

func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, healthResponse{
		Status:  "ok",
		Version: "2.0.0",
	})
}

func (h *HealthHandler) Ready(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, healthResponse{
		Status:  "ready",
		Version: "2.0.0",
	})
}
