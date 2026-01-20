package sqlite

import (
	"context"
	"fmt"
	"time"

	"github.com/aftaab/trelay/internal/core/domain"
)

// ClickRepository implements port.ClickRepository for SQLite.
type ClickRepository struct {
	db *DB
}

// NewClickRepository creates a new SQLite click repository.
func NewClickRepository(db *DB) *ClickRepository {
	return &ClickRepository{db: db}
}

// Record stores a new click event.
func (r *ClickRepository) Record(ctx context.Context, click *domain.Click) error {
	query := `
		INSERT INTO clicks (link_id, timestamp, referrer, device_hash, user_agent, ip_hash)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.ExecContext(ctx, query,
		click.LinkID,
		click.Timestamp,
		click.Referrer,
		click.DeviceHash,
		click.UserAgent,
		click.IPHash,
	)
	if err != nil {
		return fmt.Errorf("failed to record click: %w", err)
	}

	return nil
}

// GetByLinkID retrieves all clicks for a specific link.
func (r *ClickRepository) GetByLinkID(ctx context.Context, linkID int64, filter domain.StatsFilter) ([]*domain.Click, error) {
	query := `
		SELECT id, link_id, timestamp, referrer, device_hash
		FROM clicks
		WHERE link_id = ?
	`
	args := []interface{}{linkID}

	if filter.StartDate != nil {
		query += " AND timestamp >= ?"
		args = append(args, *filter.StartDate)
	}

	if filter.EndDate != nil {
		query += " AND timestamp <= ?"
		args = append(args, *filter.EndDate)
	}

	query += " ORDER BY timestamp DESC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get clicks: %w", err)
	}
	defer rows.Close()

	var clicks []*domain.Click
	for rows.Next() {
		click := &domain.Click{}
		err := rows.Scan(
			&click.ID,
			&click.LinkID,
			&click.Timestamp,
			&click.Referrer,
			&click.DeviceHash,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan click: %w", err)
		}
		clicks = append(clicks, click)
	}

	return clicks, rows.Err()
}

// GetStatsByLinkID retrieves aggregated stats for a link.
func (r *ClickRepository) GetStatsByLinkID(ctx context.Context, linkID int64, filter domain.StatsFilter) (*domain.ClickStats, error) {
	stats := &domain.ClickStats{}

	// Get total clicks
	totalQuery := `SELECT COUNT(*) FROM clicks WHERE link_id = ?`
	if err := r.db.QueryRowContext(ctx, totalQuery, linkID).Scan(&stats.TotalClicks); err != nil {
		return nil, fmt.Errorf("failed to get total clicks: %w", err)
	}

	// Get clicks by day (last 30 days)
	dayStats, err := r.GetClicksByDay(ctx, linkID, 30)
	if err != nil {
		return nil, err
	}
	stats.ClicksByDay = dayStats

	// Get top referrers
	referrerStats, err := r.GetTopReferrers(ctx, linkID, 10)
	if err != nil {
		return nil, err
	}
	stats.TopReferrers = referrerStats

	return stats, nil
}

// GetClicksByDay retrieves daily click counts for a link.
func (r *ClickRepository) GetClicksByDay(ctx context.Context, linkID int64, days int) ([]domain.DayStats, error) {
	startDate := time.Now().AddDate(0, 0, -days).Format("2006-01-02")

	query := `
		SELECT DATE(timestamp) as date, COUNT(*) as clicks
		FROM clicks
		WHERE link_id = ? AND DATE(timestamp) >= ?
		GROUP BY DATE(timestamp)
		ORDER BY date DESC
	`

	rows, err := r.db.QueryContext(ctx, query, linkID, startDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get clicks by day: %w", err)
	}
	defer rows.Close()

	var stats []domain.DayStats
	for rows.Next() {
		var s domain.DayStats
		if err := rows.Scan(&s.Date, &s.Clicks); err != nil {
			return nil, fmt.Errorf("failed to scan day stats: %w", err)
		}
		stats = append(stats, s)
	}

	return stats, rows.Err()
}

// GetClicksByMonth retrieves monthly click counts for a link.
func (r *ClickRepository) GetClicksByMonth(ctx context.Context, linkID int64, months int) ([]domain.MonthStats, error) {
	startDate := time.Now().AddDate(0, -months, 0).Format("2006-01")

	query := `
		SELECT strftime('%Y-%m', timestamp) as month, COUNT(*) as clicks
		FROM clicks
		WHERE link_id = ? AND strftime('%Y-%m', timestamp) >= ?
		GROUP BY strftime('%Y-%m', timestamp)
		ORDER BY month DESC
	`

	rows, err := r.db.QueryContext(ctx, query, linkID, startDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get clicks by month: %w", err)
	}
	defer rows.Close()

	var stats []domain.MonthStats
	for rows.Next() {
		var s domain.MonthStats
		if err := rows.Scan(&s.Month, &s.Clicks); err != nil {
			return nil, fmt.Errorf("failed to scan month stats: %w", err)
		}
		stats = append(stats, s)
	}

	return stats, rows.Err()
}

// GetTopReferrers retrieves the most common referrers for a link.
func (r *ClickRepository) GetTopReferrers(ctx context.Context, linkID int64, limit int) ([]domain.ReferrerStats, error) {
	query := `
		SELECT referrer, COUNT(*) as clicks
		FROM clicks
		WHERE link_id = ?
		GROUP BY referrer
		ORDER BY clicks DESC
		LIMIT ?
	`

	rows, err := r.db.QueryContext(ctx, query, linkID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get top referrers: %w", err)
	}
	defer rows.Close()

	var stats []domain.ReferrerStats
	for rows.Next() {
		var s domain.ReferrerStats
		if err := rows.Scan(&s.Referrer, &s.Clicks); err != nil {
			return nil, fmt.Errorf("failed to scan referrer stats: %w", err)
		}
		stats = append(stats, s)
	}

	return stats, rows.Err()
}

// DeleteByLinkID removes all clicks for a link.
func (r *ClickRepository) DeleteByLinkID(ctx context.Context, linkID int64) error {
	query := `DELETE FROM clicks WHERE link_id = ?`

	_, err := r.db.ExecContext(ctx, query, linkID)
	if err != nil {
		return fmt.Errorf("failed to delete clicks: %w", err)
	}

	return nil
}
