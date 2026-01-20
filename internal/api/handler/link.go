package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/aftaab/trelay/internal/api/response"
	"github.com/aftaab/trelay/internal/core/domain"
	"github.com/aftaab/trelay/internal/core/link"
)

type LinkHandler struct {
	service *link.Service
}

func NewLinkHandler(service *link.Service) *LinkHandler {
	return &LinkHandler{service: service}
}

func (h *LinkHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	if req.URL == "" {
		response.ValidationError(w, "url", "url is required")
		return
	}

	createdLink, err := h.service.Create(r.Context(), req)
	if err != nil {
		h.handleError(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, createdLink)
}

func (h *LinkHandler) Get(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	if slug == "" {
		response.BadRequest(w, "slug is required")
		return
	}

	password := r.URL.Query().Get("password")

	linkData, err := h.service.Get(r.Context(), slug, password)
	if err != nil {
		h.handleError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, linkData)
}

func (h *LinkHandler) Update(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	if slug == "" {
		response.BadRequest(w, "slug is required")
		return
	}

	var req domain.UpdateLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	updatedLink, err := h.service.Update(r.Context(), slug, req)
	if err != nil {
		h.handleError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, updatedLink)
}

func (h *LinkHandler) Delete(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	if slug == "" {
		response.BadRequest(w, "slug is required")
		return
	}

	permanent := r.URL.Query().Get("permanent") == "true"

	var err error
	if permanent {
		err = h.service.HardDelete(r.Context(), slug)
	} else {
		err = h.service.Delete(r.Context(), slug)
	}

	if err != nil {
		h.handleError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, map[string]bool{"deleted": true})
}

func (h *LinkHandler) Restore(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	if slug == "" {
		response.BadRequest(w, "slug is required")
		return
	}

	if err := h.service.Restore(r.Context(), slug); err != nil {
		h.handleError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, map[string]bool{"restored": true})
}

type bulkDeleteRequest struct {
	Slugs     []string `json:"slugs"`
	Permanent bool     `json:"permanent"`
}

type bulkDeleteResponse struct {
	Deleted []string `json:"deleted"`
	Failed  []string `json:"failed"`
}

func (h *LinkHandler) BulkDelete(w http.ResponseWriter, r *http.Request) {
	var req bulkDeleteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	if len(req.Slugs) == 0 {
		response.ValidationError(w, "slugs", "at least one slug is required")
		return
	}

	if len(req.Slugs) > 100 {
		response.ValidationError(w, "slugs", "maximum 100 slugs per request")
		return
	}

	result := bulkDeleteResponse{
		Deleted: make([]string, 0),
		Failed:  make([]string, 0),
	}

	for _, slug := range req.Slugs {
		var err error
		if req.Permanent {
			err = h.service.HardDelete(r.Context(), slug)
		} else {
			err = h.service.Delete(r.Context(), slug)
		}

		if err != nil {
			result.Failed = append(result.Failed, slug)
		} else {
			result.Deleted = append(result.Deleted, slug)
		}
	}

	response.JSON(w, http.StatusOK, result)
}

func (h *LinkHandler) List(w http.ResponseWriter, r *http.Request) {
	filter := domain.ListLinksFilter{
		Search: r.URL.Query().Get("search"),
		Domain: r.URL.Query().Get("domain"),
	}

	if tags := r.URL.Query()["tags"]; len(tags) > 0 {
		filter.Tags = tags
	}

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			filter.Limit = limit
		}
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil {
			filter.Offset = offset
		}
	}

	filter.IncludeDeleted = r.URL.Query().Get("include_deleted") == "true"

	links, err := h.service.List(r.Context(), filter)
	if err != nil {
		h.handleError(w, err)
		return
	}

	total, err := h.service.Count(r.Context(), filter)
	if err != nil {
		h.handleError(w, err)
		return
	}

	response.JSONWithMeta(w, http.StatusOK, links, &response.Meta{
		Total:  total,
		Limit:  filter.Limit,
		Offset: filter.Offset,
	})
}

func (h *LinkHandler) handleError(w http.ResponseWriter, err error) {
	switch err {
	case domain.ErrLinkNotFound:
		response.NotFound(w, "link not found")
	case domain.ErrLinkExpired:
		response.Error(w, http.StatusGone, "link_expired", "this link has expired")
	case domain.ErrLinkDeleted:
		response.NotFound(w, "link has been deleted")
	case domain.ErrSlugTaken:
		response.Error(w, http.StatusConflict, "slug_taken", "this slug is already in use")
	case domain.ErrSlugInvalid:
		response.ValidationError(w, "slug", "slug contains invalid characters")
	case domain.ErrSlugTooShort:
		response.ValidationError(w, "slug", "slug is too short (minimum 4 characters)")
	case domain.ErrSlugTooLong:
		response.ValidationError(w, "slug", "slug is too long (maximum 32 characters)")
	case domain.ErrURLInvalid:
		response.ValidationError(w, "url", "URL is invalid")
	case domain.ErrPasswordRequired:
		response.Error(w, http.StatusUnauthorized, "password_required", "this link requires a password")
	case domain.ErrPasswordIncorrect:
		response.Error(w, http.StatusUnauthorized, "password_incorrect", "incorrect password")
	default:
		if ve, ok := err.(domain.ValidationError); ok {
			response.ValidationError(w, ve.Field, ve.Message)
			return
		}
		response.InternalError(w)
	}
}
