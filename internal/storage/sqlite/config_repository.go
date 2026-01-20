package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// ConfigRepository implements port.ConfigRepository for SQLite.
type ConfigRepository struct {
	db *DB
}

// NewConfigRepository creates a new SQLite config repository.
func NewConfigRepository(db *DB) *ConfigRepository {
	return &ConfigRepository{db: db}
}

// Get retrieves a config value by key.
func (r *ConfigRepository) Get(ctx context.Context, key string) (string, error) {
	query := `SELECT value FROM config WHERE key = ?`

	var value string
	err := r.db.QueryRowContext(ctx, query, key).Scan(&value)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}
		return "", fmt.Errorf("failed to get config: %w", err)
	}

	return value, nil
}

// Set stores a config value.
func (r *ConfigRepository) Set(ctx context.Context, key, value string) error {
	query := `
		INSERT INTO config (key, value, updated_at)
		VALUES (?, ?, ?)
		ON CONFLICT(key) DO UPDATE SET value = excluded.value, updated_at = excluded.updated_at
	`

	_, err := r.db.ExecContext(ctx, query, key, value, time.Now())
	if err != nil {
		return fmt.Errorf("failed to set config: %w", err)
	}

	return nil
}

// Delete removes a config entry.
func (r *ConfigRepository) Delete(ctx context.Context, key string) error {
	query := `DELETE FROM config WHERE key = ?`

	_, err := r.db.ExecContext(ctx, query, key)
	if err != nil {
		return fmt.Errorf("failed to delete config: %w", err)
	}

	return nil
}

// GetAll retrieves all config entries.
func (r *ConfigRepository) GetAll(ctx context.Context) (map[string]string, error) {
	query := `SELECT key, value FROM config`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all config: %w", err)
	}
	defer rows.Close()

	config := make(map[string]string)
	for rows.Next() {
		var key, value string
		if err := rows.Scan(&key, &value); err != nil {
			return nil, fmt.Errorf("failed to scan config: %w", err)
		}
		config[key] = value
	}

	return config, rows.Err()
}
