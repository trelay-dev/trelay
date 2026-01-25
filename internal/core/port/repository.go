package port

import (
	"context"

	"github.com/aftaab/trelay/internal/core/domain"
)

// LinkRepository defines the interface for link persistence operations.
type LinkRepository interface {
	// Create stores a new link and returns the created link with ID.
	Create(ctx context.Context, link *domain.Link) (*domain.Link, error)

	// GetBySlug retrieves a link by its slug.
	GetBySlug(ctx context.Context, slug string) (*domain.Link, error)

	// GetByID retrieves a link by its ID.
	GetByID(ctx context.Context, id int64) (*domain.Link, error)

	// Update modifies an existing link.
	Update(ctx context.Context, link *domain.Link) error

	// Delete soft-deletes a link by its slug.
	Delete(ctx context.Context, slug string) error

	// HardDelete permanently removes a link.
	HardDelete(ctx context.Context, slug string) error

	// Restore recovers a soft-deleted link.
	Restore(ctx context.Context, slug string) error

	// List retrieves links matching the filter criteria.
	List(ctx context.Context, filter domain.ListLinksFilter) ([]*domain.Link, error)

	// Count returns the total number of links matching the filter.
	Count(ctx context.Context, filter domain.ListLinksFilter) (int64, error)

	// SlugExists checks if a slug is already in use.
	SlugExists(ctx context.Context, slug string) (bool, error)

	// IncrementClickCount atomically increments the click count for a link.
	IncrementClickCount(ctx context.Context, linkID int64) error

	// Burn marks a one-time link as used (soft-delete).
	Burn(ctx context.Context, linkID int64) error
}

// ClickRepository defines the interface for click/analytics persistence.
type ClickRepository interface {
	// Record stores a new click event.
	Record(ctx context.Context, click *domain.Click) error

	// GetByLinkID retrieves all clicks for a specific link.
	GetByLinkID(ctx context.Context, linkID int64, filter domain.StatsFilter) ([]*domain.Click, error)

	// GetStatsByLinkID retrieves aggregated stats for a link.
	GetStatsByLinkID(ctx context.Context, linkID int64, filter domain.StatsFilter) (*domain.ClickStats, error)

	// GetClicksByDay retrieves daily click counts for a link.
	GetClicksByDay(ctx context.Context, linkID int64, days int) ([]domain.DayStats, error)

	// GetClicksByMonth retrieves monthly click counts for a link.
	GetClicksByMonth(ctx context.Context, linkID int64, months int) ([]domain.MonthStats, error)

	// GetTopReferrers retrieves the most common referrers for a link.
	GetTopReferrers(ctx context.Context, linkID int64, limit int) ([]domain.ReferrerStats, error)

	// DeleteByLinkID removes all clicks for a link (for GDPR compliance).
	DeleteByLinkID(ctx context.Context, linkID int64) error
}

// FolderRepository defines the interface for folder persistence operations.
type FolderRepository interface {
	Create(ctx context.Context, folder *domain.Folder) (*domain.Folder, error)
	GetByID(ctx context.Context, id int64) (*domain.Folder, error)
	List(ctx context.Context) ([]*domain.Folder, error)
	Delete(ctx context.Context, id int64) error
}

// ConfigRepository defines the interface for application config persistence.
type ConfigRepository interface {
	// Get retrieves a config value by key.
	Get(ctx context.Context, key string) (string, error)

	// Set stores a config value.
	Set(ctx context.Context, key, value string) error

	// Delete removes a config entry.
	Delete(ctx context.Context, key string) error

	// GetAll retrieves all config entries.
	GetAll(ctx context.Context) (map[string]string, error)
}
