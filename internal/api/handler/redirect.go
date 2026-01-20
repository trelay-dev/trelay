package handler

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/aftaab/trelay/internal/api/response"
	"github.com/aftaab/trelay/internal/core/analytics"
	"github.com/aftaab/trelay/internal/core/domain"
	"github.com/aftaab/trelay/internal/core/link"
)

// RedirectHandler handles URL redirects.
type RedirectHandler struct {
	linkService      *link.Service
	analyticsService *analytics.Service
}

// NewRedirectHandler creates a new redirect handler.
func NewRedirectHandler(linkService *link.Service, analyticsService *analytics.Service) *RedirectHandler {
	return &RedirectHandler{
		linkService:      linkService,
		analyticsService: analyticsService,
	}
}

// Redirect handles GET /{slug} for URL redirects.
func (h *RedirectHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	if slug == "" {
		response.NotFound(w, "")
		return
	}

	linkData, err := h.linkService.GetForRedirect(r.Context(), slug)
	if err != nil {
		h.handleError(w, err)
		return
	}

	if linkData.HasPassword {
		password := r.URL.Query().Get("p")
		if password == "" {
			response.Error(w, http.StatusUnauthorized, "password_required", "this link requires a password")
			return
		}

		linkData, err = h.linkService.Get(r.Context(), slug, password)
		if err != nil {
			h.handleError(w, err)
			return
		}

		_ = h.linkService.IncrementClick(r.Context(), linkData.ID)
	}

	go func() {
		if analytics.IsBot(r.UserAgent()) {
			return
		}
		_ = h.analyticsService.RecordClick(r.Context(), linkData.ID, getClientIP(r), r.UserAgent(), r.Referer())
	}()

	http.Redirect(w, r, linkData.OriginalURL, http.StatusMovedPermanently)
}

func (h *RedirectHandler) handleError(w http.ResponseWriter, err error) {
	switch err {
	case domain.ErrLinkNotFound:
		response.NotFound(w, "link not found")
	case domain.ErrLinkExpired:
		response.Error(w, http.StatusGone, "link_expired", "this link has expired")
	case domain.ErrLinkDeleted:
		response.NotFound(w, "link not found")
	case domain.ErrPasswordIncorrect:
		response.Error(w, http.StatusUnauthorized, "password_incorrect", "incorrect password")
	default:
		response.InternalError(w)
	}
}

func getClientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// X-Forwarded-For can contain multiple IPs; the first one is the client
		if idx := strings.Index(xff, ","); idx != -1 {
			return strings.TrimSpace(xff[:idx])
		}
		return strings.TrimSpace(xff)
	}
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}
	return r.RemoteAddr
}
