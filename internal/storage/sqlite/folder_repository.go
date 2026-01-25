package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/aftaab/trelay/internal/core/domain"
)

type FolderRepository struct {
	db *DB
}

func NewFolderRepository(db *DB) *FolderRepository {
	return &FolderRepository{db: db}
}

func (r *FolderRepository) Create(ctx context.Context, folder *domain.Folder) (*domain.Folder, error) {
	query := `INSERT INTO folders (name, parent_id, created_at) VALUES (?, ?, ?)`

	result, err := r.db.ExecContext(ctx, query, folder.Name, folder.ParentID, folder.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create folder: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}

	folder.ID = id
	return folder, nil
}

func (r *FolderRepository) GetByID(ctx context.Context, id int64) (*domain.Folder, error) {
	query := `SELECT id, name, parent_id, created_at FROM folders WHERE id = ?`

	folder := &domain.Folder{}
	var parentID sql.NullInt64

	err := r.db.QueryRowContext(ctx, query, id).Scan(&folder.ID, &folder.Name, &parentID, &folder.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrFolderNotFound
		}
		return nil, fmt.Errorf("failed to get folder: %w", err)
	}

	if parentID.Valid {
		folder.ParentID = &parentID.Int64
	}

	return folder, nil
}

func (r *FolderRepository) List(ctx context.Context) ([]*domain.Folder, error) {
	query := `SELECT id, name, parent_id, created_at FROM folders ORDER BY name ASC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list folders: %w", err)
	}
	defer rows.Close()

	var folders []*domain.Folder
	for rows.Next() {
		folder := &domain.Folder{}
		var parentID sql.NullInt64

		if err := rows.Scan(&folder.ID, &folder.Name, &parentID, &folder.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan folder: %w", err)
		}

		if parentID.Valid {
			folder.ParentID = &parentID.Int64
		}

		folders = append(folders, folder)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate folders: %w", err)
	}

	return folders, nil
}

func (r *FolderRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM folders WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete folder: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return domain.ErrFolderNotFound
	}

	return nil
}
