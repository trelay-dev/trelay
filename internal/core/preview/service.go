package preview

import (
	"context"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

const (
	fetchTimeout   = 10 * time.Second
	maxBodySize    = 1024 * 1024 // 1MB
	maxTitleLen    = 200
	maxDescLen     = 500
	maxImageURLLen = 2048
)

var (
	titleRegex    = regexp.MustCompile(`(?i)<title[^>]*>([^<]+)</title>`)
	ogTitleRegex  = regexp.MustCompile(`(?i)<meta[^>]+property=["']og:title["'][^>]+content=["']([^"']+)["']`)
	ogTitleRegex2 = regexp.MustCompile(`(?i)<meta[^>]+content=["']([^"']+)["'][^>]+property=["']og:title["']`)
	ogDescRegex   = regexp.MustCompile(`(?i)<meta[^>]+property=["']og:description["'][^>]+content=["']([^"']+)["']`)
	ogDescRegex2  = regexp.MustCompile(`(?i)<meta[^>]+content=["']([^"']+)["'][^>]+property=["']og:description["']`)
	ogImageRegex  = regexp.MustCompile(`(?i)<meta[^>]+property=["']og:image["'][^>]+content=["']([^"']+)["']`)
	ogImageRegex2 = regexp.MustCompile(`(?i)<meta[^>]+content=["']([^"']+)["'][^>]+property=["']og:image["']`)
	descRegex     = regexp.MustCompile(`(?i)<meta[^>]+name=["']description["'][^>]+content=["']([^"']+)["']`)
	descRegex2    = regexp.MustCompile(`(?i)<meta[^>]+content=["']([^"']+)["'][^>]+name=["']description["']`)
)

type Preview struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	ImageURL    string `json:"image_url,omitempty"`
	FetchedAt   time.Time `json:"fetched_at"`
}

type Service struct {
	client *http.Client
}

func NewService() *Service {
	return &Service{
		client: &http.Client{
			Timeout: fetchTimeout,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				if len(via) >= 5 {
					return http.ErrUseLastResponse
				}
				return nil
			},
		},
	}
}

func (s *Service) Fetch(ctx context.Context, url string) (*Preview, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Trelay/1.0 (Link Preview)")
	req.Header.Set("Accept", "text/html")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return &Preview{FetchedAt: time.Now()}, nil
	}

	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "text/html") {
		return &Preview{FetchedAt: time.Now()}, nil
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, maxBodySize))
	if err != nil {
		return nil, err
	}

	html := string(body)
	preview := &Preview{
		Title:       extractTitle(html),
		Description: extractDescription(html),
		ImageURL:    extractImage(html),
		FetchedAt:   time.Now(),
	}

	return preview, nil
}

func extractTitle(html string) string {
	if matches := ogTitleRegex.FindStringSubmatch(html); len(matches) > 1 {
		return truncate(decodeHTML(matches[1]), maxTitleLen)
	}
	if matches := ogTitleRegex2.FindStringSubmatch(html); len(matches) > 1 {
		return truncate(decodeHTML(matches[1]), maxTitleLen)
	}
	if matches := titleRegex.FindStringSubmatch(html); len(matches) > 1 {
		return truncate(decodeHTML(matches[1]), maxTitleLen)
	}
	return ""
}

func extractDescription(html string) string {
	if matches := ogDescRegex.FindStringSubmatch(html); len(matches) > 1 {
		return truncate(decodeHTML(matches[1]), maxDescLen)
	}
	if matches := ogDescRegex2.FindStringSubmatch(html); len(matches) > 1 {
		return truncate(decodeHTML(matches[1]), maxDescLen)
	}
	if matches := descRegex.FindStringSubmatch(html); len(matches) > 1 {
		return truncate(decodeHTML(matches[1]), maxDescLen)
	}
	if matches := descRegex2.FindStringSubmatch(html); len(matches) > 1 {
		return truncate(decodeHTML(matches[1]), maxDescLen)
	}
	return ""
}

func extractImage(html string) string {
	if matches := ogImageRegex.FindStringSubmatch(html); len(matches) > 1 {
		return truncate(decodeHTML(matches[1]), maxImageURLLen)
	}
	if matches := ogImageRegex2.FindStringSubmatch(html); len(matches) > 1 {
		return truncate(decodeHTML(matches[1]), maxImageURLLen)
	}
	return ""
}

func decodeHTML(s string) string {
	s = strings.ReplaceAll(s, "&amp;", "&")
	s = strings.ReplaceAll(s, "&lt;", "<")
	s = strings.ReplaceAll(s, "&gt;", ">")
	s = strings.ReplaceAll(s, "&quot;", "\"")
	s = strings.ReplaceAll(s, "&#39;", "'")
	s = strings.ReplaceAll(s, "&apos;", "'")
	return strings.TrimSpace(s)
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
