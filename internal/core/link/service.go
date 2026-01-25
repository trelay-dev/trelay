package link

import (
	"context"
	"time"

	"github.com/aftaab/trelay/internal/core/auth"
	"github.com/aftaab/trelay/internal/core/domain"
	"github.com/aftaab/trelay/internal/core/port"
	"github.com/aftaab/trelay/internal/core/slug"
	"github.com/aftaab/trelay/internal/core/url"
)

// Service handles link business logic.
type Service struct {
	repo         port.LinkRepository
	slugGen      *slug.Generator
	urlValidator *url.Validator
}

// NewService creates a new link service.
func NewService(repo port.LinkRepository, slugLength int, selfDomains []string) *Service {
	return &Service{
		repo:         repo,
		slugGen:      slug.NewGenerator(slugLength),
		urlValidator: url.NewValidator(0, selfDomains),
	}
}

func (s *Service) Create(ctx context.Context, req domain.CreateLinkRequest) (*domain.Link, error) {
	normalizedURL, err := s.urlValidator.Normalize(req.URL)
	if err != nil {
		return nil, err
	}

	if err := s.urlValidator.Validate(normalizedURL); err != nil {
		return nil, err
	}

	linkSlug := req.Slug
	if linkSlug == "" {
		linkSlug, err = s.slugGen.Generate()
		if err != nil {
			return nil, err
		}
	} else {
		linkSlug = slug.Normalize(linkSlug)
		if err := s.slugGen.Validate(linkSlug); err != nil {
			return nil, err
		}
	}

	exists, err := s.repo.SlugExists(ctx, linkSlug)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, domain.ErrSlugTaken
	}

	var passwordHash string
	if req.Password != "" {
		passwordHash, err = auth.HashPassword(req.Password)
		if err != nil {
			return nil, err
		}
	}

	var expiresAt *time.Time
	if req.TTLHours > 0 {
		t := time.Now().Add(time.Duration(req.TTLHours) * time.Hour)
		expiresAt = &t
	}

	now := time.Now()
	link := &domain.Link{
		Slug:         linkSlug,
		OriginalURL:  normalizedURL,
		Domain:       req.Domain,
		PasswordHash: passwordHash,
		HasPassword:  passwordHash != "",
		IsOneTime:    req.IsOneTime,
		ExpiresAt:    expiresAt,
		Tags:         req.Tags,
		FolderID:     req.FolderID,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	return s.repo.Create(ctx, link)
}

// Get retrieves a link by slug with optional password verification.
func (s *Service) Get(ctx context.Context, linkSlug, password string) (*domain.Link, error) {
	link, err := s.repo.GetBySlug(ctx, linkSlug)
	if err != nil {
		return nil, err
	}

	if link.IsDeleted() {
		return nil, domain.ErrLinkDeleted
	}

	if link.IsExpired() {
		return nil, domain.ErrLinkExpired
	}

	if link.HasPassword {
		if password == "" {
			return nil, domain.ErrPasswordRequired
		}
		if !auth.VerifyPassword(password, link.PasswordHash) {
			return nil, domain.ErrPasswordIncorrect
		}
	}

	return link, nil
}

func (s *Service) GetForRedirect(ctx context.Context, linkSlug string) (*domain.Link, error) {
	link, err := s.repo.GetBySlug(ctx, linkSlug)
	if err != nil {
		return nil, err
	}

	if link.IsDeleted() {
		return nil, domain.ErrLinkDeleted
	}

	if link.IsExpired() {
		return nil, domain.ErrLinkExpired
	}

	if !link.HasPassword {
		_ = s.repo.IncrementClickCount(ctx, link.ID)
	}

	return link, nil
}

// GetByID retrieves a link by ID (admin/owner access).
func (s *Service) GetByID(ctx context.Context, id int64) (*domain.Link, error) {
	return s.repo.GetByID(ctx, id)
}

// Update modifies an existing link.
func (s *Service) Update(ctx context.Context, linkSlug string, req domain.UpdateLinkRequest) (*domain.Link, error) {
	link, err := s.repo.GetBySlug(ctx, linkSlug)
	if err != nil {
		return nil, err
	}

	if req.URL != nil {
		normalizedURL, err := s.urlValidator.Normalize(*req.URL)
		if err != nil {
			return nil, err
		}
		if err := s.urlValidator.Validate(normalizedURL); err != nil {
			return nil, err
		}
		link.OriginalURL = normalizedURL
	}

	if req.Password != nil {
		if *req.Password == "" {
			link.PasswordHash = ""
			link.HasPassword = false
		} else {
			hash, err := auth.HashPassword(*req.Password)
			if err != nil {
				return nil, err
			}
			link.PasswordHash = hash
			link.HasPassword = true
		}
	}

	if req.TTLHours != nil {
		if *req.TTLHours <= 0 {
			link.ExpiresAt = nil
		} else {
			t := time.Now().Add(time.Duration(*req.TTLHours) * time.Hour)
			link.ExpiresAt = &t
		}
	}

	if req.Tags != nil {
		link.Tags = *req.Tags
	}

	if req.FolderID != nil {
		link.FolderID = req.FolderID
	}

	link.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, link); err != nil {
		return nil, err
	}

	return link, nil
}

// Delete soft-deletes a link.
func (s *Service) Delete(ctx context.Context, linkSlug string) error {
	return s.repo.Delete(ctx, linkSlug)
}

// HardDelete permanently removes a link.
func (s *Service) HardDelete(ctx context.Context, linkSlug string) error {
	return s.repo.HardDelete(ctx, linkSlug)
}

// Restore recovers a soft-deleted link.
func (s *Service) Restore(ctx context.Context, linkSlug string) error {
	return s.repo.Restore(ctx, linkSlug)
}

// List retrieves links matching filter criteria.
func (s *Service) List(ctx context.Context, filter domain.ListLinksFilter) ([]*domain.Link, error) {
	if filter.Limit <= 0 {
		filter.Limit = 50
	}
	if filter.Limit > 200 {
		filter.Limit = 200
	}
	return s.repo.List(ctx, filter)
}

// Count returns total links matching filter.
func (s *Service) Count(ctx context.Context, filter domain.ListLinksFilter) (int64, error) {
	return s.repo.Count(ctx, filter)
}

// IncrementClick increments click count (for password-protected links after auth).
func (s *Service) IncrementClick(ctx context.Context, linkID int64) error {
	return s.repo.IncrementClickCount(ctx, linkID)
}

// Burn marks a one-time link as consumed (soft-delete).
func (s *Service) Burn(ctx context.Context, linkID int64) error {
	return s.repo.Burn(ctx, linkID)
}
