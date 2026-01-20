package handler

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/aftaab/trelay/internal/api/response"
	"github.com/aftaab/trelay/internal/core/analytics"
	"github.com/aftaab/trelay/internal/core/domain"
	"github.com/aftaab/trelay/internal/core/link"
)

type StatsHandler struct {
	linkService      *link.Service
	analyticsService *analytics.Service
}

func NewStatsHandler(linkService *link.Service, analyticsService *analytics.Service) *StatsHandler {
	return &StatsHandler{
		linkService:      linkService,
		analyticsService: analyticsService,
	}
}

func (h *StatsHandler) GetStats(w http.ResponseWriter, r *http.Request) {
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

	filter := domain.StatsFilter{}
	if period := r.URL.Query().Get("period"); period != "" {
		filter.Period = domain.StatsPeriod(period)
	}

	stats, err := h.analyticsService.GetStats(r.Context(), linkData.ID, filter)
	if err != nil {
		response.InternalError(w)
		return
	}

	exportFormat := r.URL.Query().Get("export")
	switch exportFormat {
	case "csv":
		h.exportCSV(w, slug, stats)
	case "json":
		h.exportJSON(w, slug, stats)
	default:
		response.JSON(w, http.StatusOK, stats)
	}
}

func (h *StatsHandler) exportCSV(w http.ResponseWriter, slug string, stats *domain.ClickStats) {
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s-stats.csv", slug))

	writer := csv.NewWriter(w)
	defer writer.Flush()

	writer.Write([]string{"metric", "value"})
	writer.Write([]string{"total_clicks", strconv.FormatInt(stats.TotalClicks, 10)})
	writer.Write([]string{})

	if len(stats.ClicksByDay) > 0 {
		writer.Write([]string{"date", "clicks"})
		for _, d := range stats.ClicksByDay {
			writer.Write([]string{d.Date, strconv.FormatInt(d.Clicks, 10)})
		}
		writer.Write([]string{})
	}

	if len(stats.TopReferrers) > 0 {
		writer.Write([]string{"referrer", "clicks"})
		for _, r := range stats.TopReferrers {
			writer.Write([]string{r.Referrer, strconv.FormatInt(r.Clicks, 10)})
		}
	}
}

func (h *StatsHandler) exportJSON(w http.ResponseWriter, slug string, stats *domain.ClickStats) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s-stats.json", slug))

	json.NewEncoder(w).Encode(stats)
}

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
	linkData, err := h.linkService.Get(r.Context(), slug, "")
	if err == domain.ErrPasswordRequired {
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
