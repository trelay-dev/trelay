package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/aftaab/trelay/internal/api/response"
	"github.com/aftaab/trelay/internal/core/analytics"
	"github.com/aftaab/trelay/internal/core/domain"
	"github.com/aftaab/trelay/internal/core/link"
)

// StatsHandler handles analytics-related HTTP requests.
type StatsHandler struct {
	linkService      *link.Service
	analyticsService *analytics.Service
}

// NewStatsHandler creates a new stats handler.
func NewStatsHandler(linkService *link.Service, analyticsService *analytics.Service) *StatsHandler {
	return &StatsHandler{
		linkService:      linkService,
		analyticsService: analyticsService,
	}
}

// GetStats handles GET /api/v1/stats/{slug}
func (h *StatsHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	if slug == "" {
		response.BadRequest(w, "slug is required")
		return
	}

	// Get link to verify it exists and get ID
	linkData, err := h.linkService.Get(r.Context(), slug, "")
	if err != nil {
		if err == domain.ErrPasswordRequired {
			// For stats, allow viewing without password
			linkData, err = h.linkService.GetByID(r.Context(), 0)
		}
		if err != nil && err != domain.ErrPasswordRequired {
			h.handleError(w, err)
			return
		}
	}

	// Re-fetch by slug without password check for stats
	linkData, err = h.getLinkBySlugForStats(r, slug)
	if err != nil {
		h.handleError(w, err)
		return
	}

	filter := domain.StatsFilter{}
	if period := r.URL.Query().Get("period"); period != "" {
		filter.Period = domain.StatsPeriod(period)
	}

	stats, err := h.analyticsService.GetStats(r.Context(), linkData.ID, filter)
	if err != nil {
		response.InternalError(w)
		return
	}

	response.JSON(w, http.StatusOK, stats)
}

// GetDailyStats handles GET /api/v1/stats/{slug}/daily
func (h *StatsHandler) GetDailyStats(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	if slug == "" {
		response.BadRequest(w, "slug is required")
		return
	}

	linkData, err := h.getLinkBySlugForStats(r, slug)
	if err != nil {
		h.handleError(w, err)
		return
	}

	days := 30
	if daysStr := r.URL.Query().Get("days"); daysStr != "" {
		if d, err := strconv.Atoi(daysStr); err == nil && d > 0 {
			days = d
		}
	}

	stats, err := h.analyticsService.GetClicksByDay(r.Context(), linkData.ID, days)
	if err != nil {
		response.InternalError(w)
		return
	}

	response.JSON(w, http.StatusOK, stats)
}

// GetMonthlyStats handles GET /api/v1/stats/{slug}/monthly
func (h *StatsHandler) GetMonthlyStats(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	if slug == "" {
		response.BadRequest(w, "slug is required")
		return
	}

	linkData, err := h.getLinkBySlugForStats(r, slug)
	if err != nil {
		h.handleError(w, err)
		return
	}

	months := 12
	if monthsStr := r.URL.Query().Get("months"); monthsStr != "" {
		if m, err := strconv.Atoi(monthsStr); err == nil && m > 0 {
			months = m
		}
	}

	stats, err := h.analyticsService.GetClicksByMonth(r.Context(), linkData.ID, months)
	if err != nil {
		response.InternalError(w)
		return
	}

	response.JSON(w, http.StatusOK, stats)
}

// GetReferrers handles GET /api/v1/stats/{slug}/referrers
func (h *StatsHandler) GetReferrers(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	if slug == "" {
		response.BadRequest(w, "slug is required")
		return
	}

	linkData, err := h.getLinkBySlugForStats(r, slug)
	if err != nil {
		h.handleError(w, err)
		return
	}

	limit := 10
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	stats, err := h.analyticsService.GetTopReferrers(r.Context(), linkData.ID, limit)
	if err != nil {
		response.InternalError(w)
		return
	}

	response.JSON(w, http.StatusOK, stats)
}

func (h *StatsHandler) getLinkBySlugForStats(r *http.Request, slug string) (*domain.Link, error) {
	// Use a direct repo call or add a GetBySlugNoAuth method
	// For now, try to get without password
	linkData, err := h.linkService.Get(r.Context(), slug, "")
	if err == domain.ErrPasswordRequired {
		// For stats, we need to bypass password check
		// This requires getting the link by slug without password validation
		// We'll need to add a method for this or use the repo directly
		// For now, return the error
		return nil, err
	}
	return linkData, err
}

func (h *StatsHandler) handleError(w http.ResponseWriter, err error) {
	switch err {
	case domain.ErrLinkNotFound:
		response.NotFound(w, "link not found")
	case domain.ErrLinkExpired:
		response.Error(w, http.StatusGone, "link_expired", "this link has expired")
	case domain.ErrPasswordRequired:
		response.Error(w, http.StatusUnauthorized, "password_required", "this link requires a password")
	default:
		response.InternalError(w)
	}
}
