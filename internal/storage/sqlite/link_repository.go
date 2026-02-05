package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aftaab/trelay/internal/core/domain"
)

// LinkRepository implements port.LinkRepository for SQLite.
type LinkRepository struct {
	db *DB
}

// NewLinkRepository creates a new SQLite link repository.
func NewLinkRepository(db *DB) *LinkRepository {
	return &LinkRepository{db: db}
}

// Create stores a new link and returns the created link with ID.
func (r *LinkRepository) Create(ctx context.Context, link *domain.Link) (*domain.Link, error) {
	tagsJSON, err := link.TagsJSON()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal tags: %w", err)
	}

	// Remove any soft-deleted link with the same slug to allow reuse
	_, _ = r.db.ExecContext(ctx, `DELETE FROM links WHERE slug = ? AND deleted_at IS NOT NULL`, link.Slug)

	query := `
		INSERT INTO links (slug, original_url, domain, password_hash, expires_at, tags, folder_id, is_one_time, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.ExecContext(ctx, query,
		link.Slug,
		link.OriginalURL,
		link.Domain,
		link.PasswordHash,
		link.ExpiresAt,
		tagsJSON,
		link.FolderID,
		link.IsOneTime,
		link.CreatedAt,
		link.UpdatedAt,
	)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return nil, domain.ErrSlugTaken
		}
		return nil, fmt.Errorf("failed to create link: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}

	link.ID = id
	link.HasPassword = link.PasswordHash != ""
	return link, nil
}

// GetBySlug retrieves a link by its slug.
func (r *LinkRepository) GetBySlug(ctx context.Context, slug string) (*domain.Link, error) {
	query := `
		SELECT id, slug, original_url, domain, password_hash, expires_at, tags, folder_id, is_one_time, click_count, created_at, updated_at, deleted_at
		FROM links
		WHERE slug = ?
	`

	link := &domain.Link{}
	var tagsJSON string
	var expiresAt, deletedAt sql.NullTime
	var folderID sql.NullInt64

	err := r.db.QueryRowContext(ctx, query, slug).Scan(
		&link.ID,
		&link.Slug,
		&link.OriginalURL,
		&link.Domain,
		&link.PasswordHash,
		&expiresAt,
		&tagsJSON,
		&folderID,
		&link.IsOneTime,
		&link.ClickCount,
		&link.CreatedAt,
		&link.UpdatedAt,
		&deletedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrLinkNotFound
		}
		return nil, fmt.Errorf("failed to get link: %w", err)
	}

	if expiresAt.Valid {
		link.ExpiresAt = &expiresAt.Time
	}
	if deletedAt.Valid {
		link.DeletedAt = &deletedAt.Time
	}
	if folderID.Valid {
		link.FolderID = &folderID.Int64
	}

	if err := link.ParseTagsJSON(tagsJSON); err != nil {
		return nil, fmt.Errorf("failed to parse tags: %w", err)
	}

	link.HasPassword = link.PasswordHash != ""
	return link, nil
}

// GetByID retrieves a link by its ID.
func (r *LinkRepository) GetByID(ctx context.Context, id int64) (*domain.Link, error) {
	query := `
		SELECT id, slug, original_url, domain, password_hash, expires_at, tags, folder_id, is_one_time, click_count, created_at, updated_at, deleted_at
		FROM links
		WHERE id = ?
	`

	link := &domain.Link{}
	var tagsJSON string
	var expiresAt, deletedAt sql.NullTime
	var folderID sql.NullInt64

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&link.ID,
		&link.Slug,
		&link.OriginalURL,
		&link.Domain,
		&link.PasswordHash,
		&expiresAt,
		&tagsJSON,
		&folderID,
		&link.IsOneTime,
		&link.ClickCount,
		&link.CreatedAt,
		&link.UpdatedAt,
		&deletedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrLinkNotFound
		}
		return nil, fmt.Errorf("failed to get link: %w", err)
	}

	if expiresAt.Valid {
		link.ExpiresAt = &expiresAt.Time
	}
	if deletedAt.Valid {
		link.DeletedAt = &deletedAt.Time
	}
	if folderID.Valid {
		link.FolderID = &folderID.Int64
	}

	if err := link.ParseTagsJSON(tagsJSON); err != nil {
		return nil, fmt.Errorf("failed to parse tags: %w", err)
	}

	link.HasPassword = link.PasswordHash != ""
	return link, nil
}

// Update modifies an existing link.
func (r *LinkRepository) Update(ctx context.Context, link *domain.Link) error {
	tagsJSON, err := link.TagsJSON()
	if err != nil {
		return fmt.Errorf("failed to marshal tags: %w", err)
	}

	query := `
		UPDATE links
		SET original_url = ?, domain = ?, password_hash = ?, expires_at = ?, tags = ?, folder_id = ?, updated_at = ?
		WHERE id = ?
	`

	result, err := r.db.ExecContext(ctx, query,
		link.OriginalURL,
		link.Domain,
		link.PasswordHash,
		link.ExpiresAt,
		tagsJSON,
		link.FolderID,
		time.Now(),
		link.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update link: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return domain.ErrLinkNotFound
	}

	return nil
}

// Delete soft-deletes a link by its slug.
func (r *LinkRepository) Delete(ctx context.Context, slug string) error {
	query := `UPDATE links SET deleted_at = ? WHERE slug = ? AND deleted_at IS NULL`

	result, err := r.db.ExecContext(ctx, query, time.Now(), slug)
	if err != nil {
		return fmt.Errorf("failed to delete link: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return domain.ErrLinkNotFound
	}

	return nil
}

// HardDelete permanently removes a link.
func (r *LinkRepository) HardDelete(ctx context.Context, slug string) error {
	query := `DELETE FROM links WHERE slug = ?`

	result, err := r.db.ExecContext(ctx, query, slug)
	if err != nil {
		return fmt.Errorf("failed to hard delete link: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return domain.ErrLinkNotFound
	}

	return nil
}

// Restore recovers a soft-deleted link.
func (r *LinkRepository) Restore(ctx context.Context, slug string) error {
	query := `UPDATE links SET deleted_at = NULL, updated_at = ? WHERE slug = ? AND deleted_at IS NOT NULL`

	result, err := r.db.ExecContext(ctx, query, time.Now(), slug)
	if err != nil {
		return fmt.Errorf("failed to restore link: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return domain.ErrLinkNotFound
	}

	return nil
}

// List retrieves links matching the filter criteria.
func (r *LinkRepository) List(ctx context.Context, filter domain.ListLinksFilter) ([]*domain.Link, error) {
	var conditions []string
	var args []interface{}

	if filter.OnlyDeleted {
		conditions = append(conditions, "deleted_at IS NOT NULL")
	} else if !filter.IncludeDeleted {
		conditions = append(conditions, "deleted_at IS NULL")
	}

	if filter.Search != "" {
		conditions = append(conditions, "(slug LIKE ? OR original_url LIKE ?)")
		searchTerm := "%" + filter.Search + "%"
		args = append(args, searchTerm, searchTerm)
	}

	if filter.Domain != "" {
		conditions = append(conditions, "domain = ?")
		args = append(args, filter.Domain)
	}

	if filter.FolderID != nil {
		conditions = append(conditions, "folder_id = ?")
		args = append(args, *filter.FolderID)
	}

	if len(filter.Tags) > 0 {
		for _, tag := range filter.Tags {
			conditions = append(conditions, "tags LIKE ?")
			args = append(args, "%\""+tag+"\"%")
		}
	}

	if filter.CreatedAfter != "" {
		conditions = append(conditions, "created_at >= ?")
		args = append(args, filter.CreatedAfter)
	}

	if filter.CreatedBefore != "" {
		conditions = append(conditions, "created_at <= ?")
		args = append(args, filter.CreatedBefore)
	}

	query := `
		SELECT id, slug, original_url, domain, password_hash, expires_at, tags, folder_id, is_one_time, click_count, created_at, updated_at, deleted_at
		FROM links
	`

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY created_at DESC"

	if filter.Limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", filter.Limit)
	}

	if filter.Offset > 0 {
		query += fmt.Sprintf(" OFFSET %d", filter.Offset)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list links: %w", err)
	}
	defer rows.Close()

	var links []*domain.Link
	for rows.Next() {
		link := &domain.Link{}
		var tagsJSON string
		var expiresAt, deletedAt sql.NullTime
		var folderID sql.NullInt64

		err := rows.Scan(
			&link.ID,
			&link.Slug,
			&link.OriginalURL,
			&link.Domain,
			&link.PasswordHash,
			&expiresAt,
			&tagsJSON,
			&folderID,
			&link.IsOneTime,
			&link.ClickCount,
			&link.CreatedAt,
			&link.UpdatedAt,
			&deletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan link: %w", err)
		}

		if expiresAt.Valid {
			link.ExpiresAt = &expiresAt.Time
		}
		if deletedAt.Valid {
			link.DeletedAt = &deletedAt.Time
		}
		if folderID.Valid {
			link.FolderID = &folderID.Int64
		}

		if err := link.ParseTagsJSON(tagsJSON); err != nil {
			return nil, fmt.Errorf("failed to parse tags: %w", err)
		}

		link.HasPassword = link.PasswordHash != ""
		links = append(links, link)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate links: %w", err)
	}

	return links, nil
}

// Count returns the total number of links matching the filter.
func (r *LinkRepository) Count(ctx context.Context, filter domain.ListLinksFilter) (int64, error) {
	var conditions []string
	var args []interface{}

	if filter.OnlyDeleted {
		conditions = append(conditions, "deleted_at IS NOT NULL")
	} else if !filter.IncludeDeleted {
		conditions = append(conditions, "deleted_at IS NULL")
	}

	if filter.Search != "" {
		conditions = append(conditions, "(slug LIKE ? OR original_url LIKE ?)")
		searchTerm := "%" + filter.Search + "%"
		args = append(args, searchTerm, searchTerm)
	}

	if filter.Domain != "" {
		conditions = append(conditions, "domain = ?")
		args = append(args, filter.Domain)
	}

	if filter.FolderID != nil {
		conditions = append(conditions, "folder_id = ?")
		args = append(args, *filter.FolderID)
	}

	query := `SELECT COUNT(*) FROM links`

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	var count int64
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count links: %w", err)
	}

	return count, nil
}

// SlugExists checks if a slug is already in use.
func (r *LinkRepository) SlugExists(ctx context.Context, slug string) (bool, error) {
	// Only check non-deleted links so slugs can be reused after deletion
	query := `SELECT EXISTS(SELECT 1 FROM links WHERE slug = ? AND deleted_at IS NULL)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, slug).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check slug existence: %w", err)
	}

	return exists, nil
}

// IncrementClickCount atomically increments the click count for a link.
func (r *LinkRepository) IncrementClickCount(ctx context.Context, linkID int64) error {
	query := `UPDATE links SET click_count = click_count + 1 WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query, linkID)
	if err != nil {
		return fmt.Errorf("failed to increment click count: %w", err)
	}

	return nil
}

func (r *LinkRepository) Burn(ctx context.Context, linkID int64) error {
	query := `UPDATE links SET deleted_at = ? WHERE id = ? AND deleted_at IS NULL`

	_, err := r.db.ExecContext(ctx, query, time.Now(), linkID)
	if err != nil {
		return fmt.Errorf("failed to burn one-time link: %w", err)
	}

	return nil
}
