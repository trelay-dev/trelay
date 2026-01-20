package analytics

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"time"

	"github.com/aftaab/trelay/internal/core/domain"
	"github.com/aftaab/trelay/internal/core/port"
)

// Service handles analytics and click tracking.
type Service struct {
	clickRepo       port.ClickRepository
	anonymizeIP     bool
	enabled         bool
}

// NewService creates a new analytics service.
func NewService(clickRepo port.ClickRepository, anonymizeIP, enabled bool) *Service {
	return &Service{
		clickRepo:   clickRepo,
		anonymizeIP: anonymizeIP,
		enabled:     enabled,
	}
}

// RecordClick records a click event for a link.
func (s *Service) RecordClick(ctx context.Context, linkID int64, ip, userAgent, referrer string) error {
	if !s.enabled {
		return nil
	}

	click := &domain.Click{
		LinkID:     linkID,
		Timestamp:  time.Now().UTC(),
		Referrer:   normalizeReferrer(referrer),
		DeviceHash: hashDeviceInfo(userAgent),
		UserAgent:  userAgent,
		IPHash:     s.hashIP(ip),
	}

	return s.clickRepo.Record(ctx, click)
}

// GetStats retrieves analytics for a link.
func (s *Service) GetStats(ctx context.Context, linkID int64, filter domain.StatsFilter) (*domain.ClickStats, error) {
	return s.clickRepo.GetStatsByLinkID(ctx, linkID, filter)
}

// GetClicksByDay retrieves daily click data.
func (s *Service) GetClicksByDay(ctx context.Context, linkID int64, days int) ([]domain.DayStats, error) {
	if days <= 0 {
		days = 30
	}
	if days > 365 {
		days = 365
	}
	return s.clickRepo.GetClicksByDay(ctx, linkID, days)
}

// GetClicksByMonth retrieves monthly click data.
func (s *Service) GetClicksByMonth(ctx context.Context, linkID int64, months int) ([]domain.MonthStats, error) {
	if months <= 0 {
		months = 12
	}
	if months > 24 {
		months = 24
	}
	return s.clickRepo.GetClicksByMonth(ctx, linkID, months)
}

// GetTopReferrers retrieves top referrer sources.
func (s *Service) GetTopReferrers(ctx context.Context, linkID int64, limit int) ([]domain.ReferrerStats, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	return s.clickRepo.GetTopReferrers(ctx, linkID, limit)
}

// DeleteStats removes all analytics for a link (GDPR compliance).
func (s *Service) DeleteStats(ctx context.Context, linkID int64) error {
	return s.clickRepo.DeleteByLinkID(ctx, linkID)
}

// hashIP creates an anonymized hash of an IP address.
func (s *Service) hashIP(ip string) string {
	if !s.anonymizeIP || ip == "" {
		return ""
	}

	// For IPv4, hash only the first 3 octets
	// For IPv6, hash only the first 48 bits
	parts := strings.Split(ip, ".")
	if len(parts) == 4 {
		ip = strings.Join(parts[:3], ".") + ".0"
	}

	hash := sha256.Sum256([]byte(ip))
	return hex.EncodeToString(hash[:8])
}

// hashDeviceInfo creates a privacy-preserving device fingerprint.
func hashDeviceInfo(userAgent string) string {
	if userAgent == "" {
		return ""
	}

	// Extract general device type without unique identifiers
	deviceType := extractDeviceType(userAgent)
	hash := sha256.Sum256([]byte(deviceType))
	return hex.EncodeToString(hash[:8])
}

// extractDeviceType determines general device category from user agent.
func extractDeviceType(userAgent string) string {
	ua := strings.ToLower(userAgent)

	switch {
	case strings.Contains(ua, "mobile") || strings.Contains(ua, "android"):
		return "mobile"
	case strings.Contains(ua, "tablet") || strings.Contains(ua, "ipad"):
		return "tablet"
	case strings.Contains(ua, "bot") || strings.Contains(ua, "crawler") || strings.Contains(ua, "spider"):
		return "bot"
	default:
		return "desktop"
	}
}

// normalizeReferrer cleans up referrer URL for storage.
func normalizeReferrer(referrer string) string {
	if referrer == "" {
		return "direct"
	}

	// Remove query parameters and fragments for privacy
	if idx := strings.Index(referrer, "?"); idx > 0 {
		referrer = referrer[:idx]
	}
	if idx := strings.Index(referrer, "#"); idx > 0 {
		referrer = referrer[:idx]
	}

	// Limit length
	if len(referrer) > 500 {
		referrer = referrer[:500]
	}

	return referrer
}

// IsBot checks if a user agent appears to be a bot.
func IsBot(userAgent string) bool {
	ua := strings.ToLower(userAgent)
	botIndicators := []string{
		"bot", "crawler", "spider", "slurp", "facebook",
		"twitter", "linkedin", "pinterest", "whatsapp",
		"telegram", "preview", "fetch", "curl", "wget",
	}

	for _, indicator := range botIndicators {
		if strings.Contains(ua, indicator) {
			return true
		}
	}
	return false
}
