package url

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/aftaab/trelay/internal/core/domain"
)

const (
	MaxURLLength     = 2048
	ReachableTimeout = 5 * time.Second
)

var (
	// allowedSchemes for shortened URLs
	allowedSchemes = map[string]bool{
		"http":  true,
		"https": true,
	}

	// blockedHosts that should not be shortened (prevent redirect loops)
	blockedHosts = map[string]bool{
		"localhost": true,
		"127.0.0.1": true,
	}
)

// Validator handles URL validation and normalization.
type Validator struct {
	maxLength    int
	blockedHosts map[string]bool
	selfDomains  []string
}

// NewValidator creates a new URL validator.
func NewValidator(maxLength int, selfDomains []string) *Validator {
	if maxLength <= 0 {
		maxLength = MaxURLLength
	}

	// Merge blocked hosts with self domains to prevent loops
	blocked := make(map[string]bool)
	for host := range blockedHosts {
		blocked[host] = true
	}
	for _, domain := range selfDomains {
		blocked[strings.ToLower(domain)] = true
	}

	return &Validator{
		maxLength:    maxLength,
		blockedHosts: blocked,
		selfDomains:  selfDomains,
	}
}

// Validate checks if a URL is valid for shortening.
func (v *Validator) Validate(rawURL string) error {
	if rawURL == "" {
		return domain.NewValidationError("url", "URL is required")
	}

	if len(rawURL) > v.maxLength {
		return domain.NewValidationError("url", "URL exceeds maximum length")
	}

	parsed, err := url.Parse(rawURL)
	if err != nil {
		return domain.ErrURLInvalid
	}

	if !allowedSchemes[strings.ToLower(parsed.Scheme)] {
		return domain.NewValidationError("url", "URL must use http or https scheme")
	}

	if parsed.Host == "" {
		return domain.ErrURLInvalid
	}

	host := strings.ToLower(parsed.Hostname())
	if v.blockedHosts[host] {
		return domain.NewValidationError("url", "this host cannot be shortened")
	}

	return nil
}

// Normalize ensures a URL has a scheme and is properly formatted.
func (v *Validator) Normalize(rawURL string) (string, error) {
	rawURL = strings.TrimSpace(rawURL)

	// Add https:// if no scheme present
	if !strings.Contains(rawURL, "://") {
		rawURL = "https://" + rawURL
	}

	parsed, err := url.Parse(rawURL)
	if err != nil {
		return "", domain.ErrURLInvalid
	}

	// Ensure path is at least "/"
	if parsed.Path == "" {
		parsed.Path = "/"
	}

	return parsed.String(), nil
}

// CheckReachable attempts to reach the URL (optional validation).
func (v *Validator) CheckReachable(ctx context.Context, rawURL string) error {
	ctx, cancel := context.WithTimeout(ctx, ReachableTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodHead, rawURL, nil)
	if err != nil {
		return domain.ErrURLInvalid
	}

	client := &http.Client{
		Timeout: ReachableTimeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 10 {
				return http.ErrUseLastResponse
			}
			return nil
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return domain.ErrURLUnreachable
	}
	defer resp.Body.Close()

	// Accept 2xx, 3xx, and some 4xx status codes
	if resp.StatusCode >= 500 {
		return domain.ErrURLUnreachable
	}

	return nil
}

// ExtractHost extracts the hostname from a URL.
func ExtractHost(rawURL string) (string, error) {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	return parsed.Hostname(), nil
}
