package handler

import (
	"fmt"
	"html"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/aftaab/trelay/internal/api/response"
	"github.com/aftaab/trelay/internal/core/analytics"
	"github.com/aftaab/trelay/internal/core/domain"
	"github.com/aftaab/trelay/internal/core/link"
)

type RedirectHandler struct {
	linkService      *link.Service
	analyticsService *analytics.Service
}

func NewRedirectHandler(linkService *link.Service, analyticsService *analytics.Service) *RedirectHandler {
	return &RedirectHandler{
		linkService:      linkService,
		analyticsService: analyticsService,
	}
}

// Redirect handles GET /{slug} (short link redirect or password gate).
func (h *RedirectHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.Error(w, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed")
		return
	}

	slug := chi.URLParam(r, "slug")
	if slug == "" {
		response.NotFound(w, "")
		return
	}

	password := r.URL.Query().Get("p")
	h.handleRedirect(w, r, slug, password)
}

// RedirectPost handles POST /{slug} (password form submission for protected links).
func (h *RedirectHandler) RedirectPost(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	if slug == "" {
		response.NotFound(w, "")
		return
	}

	if err := r.ParseForm(); err != nil {
		response.BadRequest(w, "invalid form body")
		return
	}

	password := r.FormValue("password")
	if password == "" {
		password = r.FormValue("p")
	}

	h.handleRedirect(w, r, slug, password)
}

func wantsRedirectJSON(r *http.Request) bool {
	if strings.Contains(r.Header.Get("Accept"), "application/json") {
		return true
	}
	if r.URL.Query().Get("format") == "json" {
		return true
	}
	return false
}

func (h *RedirectHandler) handleRedirect(w http.ResponseWriter, r *http.Request, slug, password string) {
	linkData, err := h.linkService.GetForRedirect(r.Context(), slug)
	if err != nil {
		h.handleError(w, err)
		return
	}

	if linkData.Domain != "" {
		requestHost := r.Host
		if idx := strings.Index(requestHost, ":"); idx != -1 {
			requestHost = requestHost[:idx]
		}
		if requestHost != linkData.Domain {
			response.NotFound(w, "link not found")
			return
		}
	}

	if !linkData.HasPassword {
		if linkData.IsOneTime {
			_ = h.linkService.Burn(r.Context(), linkData.ID)
		}
		h.recordAnalyticsAsync(r, linkData.ID)
		http.Redirect(w, r, linkData.OriginalURL, http.StatusMovedPermanently)
		return
	}

	if password == "" {
		if wantsRedirectJSON(r) {
			response.Error(w, http.StatusUnauthorized, "password_required", "this link requires a password")
			return
		}
		h.writePasswordPage(w, slug, false)
		return
	}

	linkData, err = h.linkService.Get(r.Context(), slug, password)
	if err != nil {
		if err == domain.ErrPasswordIncorrect {
			if wantsRedirectJSON(r) {
				response.Error(w, http.StatusUnauthorized, "password_incorrect", "incorrect password")
				return
			}
			h.writePasswordPage(w, slug, true)
			return
		}
		h.handleError(w, err)
		return
	}

	_ = h.linkService.IncrementClick(r.Context(), linkData.ID)

	if linkData.IsOneTime {
		_ = h.linkService.Burn(r.Context(), linkData.ID)
	}

	h.recordAnalyticsAsync(r, linkData.ID)
	http.Redirect(w, r, linkData.OriginalURL, http.StatusMovedPermanently)
}

func (h *RedirectHandler) recordAnalyticsAsync(r *http.Request, linkID int64) {
	go func() {
		if analytics.IsBot(r.UserAgent()) {
			return
		}
		_ = h.analyticsService.RecordClick(r.Context(), linkID, getClientIP(r), r.UserAgent(), r.Referer())
	}()
}

func (h *RedirectHandler) writePasswordPage(w http.ResponseWriter, slug string, wrongPassword bool) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusUnauthorized)

	title := html.EscapeString(slug)
	errMsg := ""
	if wrongPassword {
		errMsg = `<p class="err">Incorrect password. Try again.</p>`
	}

	escSlug := html.EscapeString(slug)
	fmt.Fprintf(w, `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8"/>
<meta name="viewport" content="width=device-width, initial-scale=1"/>
<title>Password required — %s</title>
<style>
body{font-family:system-ui,-apple-system,sans-serif;background:#0f1419;color:#e6edf3;margin:0;min-height:100vh;display:flex;align-items:center;justify-content:center;padding:24px;}
.card{max-width:400px;width:100%%;background:#161b22;border:1px solid #30363d;border-radius:12px;padding:28px;}
h1{font-size:1.125rem;margin:0 0 8px;font-weight:600;}
p.sub{margin:0 0 20px;color:#8b949e;font-size:0.875rem;}
label{display:block;font-size:0.75rem;color:#8b949e;margin-bottom:6px;}
input[type=password]{width:100%%;box-sizing:border-box;padding:10px 12px;border-radius:8px;border:1px solid #30363d;background:#0d1117;color:#e6edf3;font-size:1rem;}
input[type=password]:focus{outline:2px solid #388bfd;outline-offset:0;border-color:#388bfd;}
button{margin-top:16px;width:100%%;padding:10px 16px;border:none;border-radius:8px;background:#238636;color:#fff;font-weight:600;font-size:0.9375rem;cursor:pointer;}
button:hover{background:#2ea043;}
.err{color:#f85149;font-size:0.875rem;margin:0 0 12px;}
.hint{font-size:0.75rem;color:#6e7681;margin-top:16px;}
</style>
</head>
<body>
<div class="card">
<h1>Protected link</h1>
<p class="sub">/%s requires a password to continue.</p>
%s
<form method="post" action="/%s" autocomplete="current-password">
<label for="password">Password</label>
<input id="password" name="password" type="password" required autofocus/>
<button type="submit">Continue</button>
</form>
<p class="hint">You can still open this link with <code>?p=…</code> in the URL if you prefer.</p>
</div>
</body>
</html>`, title, escSlug, errMsg, escSlug)
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
