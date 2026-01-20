package handler

import (
	"net/http"

	"github.com/aftaab/trelay/internal/api/response"
	"github.com/aftaab/trelay/internal/core/preview"
)

type PreviewHandler struct {
	service *preview.Service
}

func NewPreviewHandler(service *preview.Service) *PreviewHandler {
	return &PreviewHandler{service: service}
}

type previewRequest struct {
	URL string `json:"url"`
}

func (h *PreviewHandler) Fetch(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	if url == "" {
		response.BadRequest(w, "url query parameter is required")
		return
	}

	previewData, err := h.service.Fetch(r.Context(), url)
	if err != nil {
		response.Error(w, http.StatusBadGateway, "fetch_failed", "failed to fetch preview")
		return
	}

	response.JSON(w, http.StatusOK, previewData)
}
