package folder

import (
	"context"
	"time"

	"github.com/aftaab/trelay/internal/core/domain"
	"github.com/aftaab/trelay/internal/core/port"
)

type Service struct {
	repo port.FolderRepository
}

func NewService(repo port.FolderRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, req domain.CreateFolderRequest) (*domain.Folder, error) {
	if req.Name == "" {
		return nil, domain.NewValidationError("name", "folder name is required")
	}

	// Validate parent folder exists if provided
	if req.ParentID != nil {
		_, err := s.repo.GetByID(ctx, *req.ParentID)
		if err != nil {
			if err == domain.ErrFolderNotFound {
				return nil, domain.ErrParentFolderNotFound
			}
			return nil, err
		}
	}

	folder := &domain.Folder{
		Name:      req.Name,
		ParentID:  req.ParentID,
		CreatedAt: time.Now(),
	}

	return s.repo.Create(ctx, folder)
}

func (s *Service) List(ctx context.Context) ([]*domain.Folder, error) {
	return s.repo.List(ctx)
}

func (s *Service) Get(ctx context.Context, id int64) (*domain.Folder, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
