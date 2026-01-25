package domain

import (
	"encoding/json"
	"time"
)

// Link represents a shortened URL with its metadata.
type Link struct {
	ID           int64      `json:"id"`
	Slug         string     `json:"slug"`
	OriginalURL  string     `json:"original_url"`
	Domain       string     `json:"domain,omitempty"`
	PasswordHash string     `json:"-"`
	HasPassword  bool       `json:"has_password"`
	IsOneTime    bool       `json:"is_one_time,omitempty"`
	ExpiresAt    *time.Time `json:"expires_at,omitempty"`
	Tags         []string   `json:"tags,omitempty"`
	FolderID     *int64     `json:"folder_id,omitempty"`
	ClickCount   int64      `json:"click_count"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"`
}

// IsExpired checks if the link has expired.
func (l *Link) IsExpired() bool {
	if l.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*l.ExpiresAt)
}

// IsDeleted checks if the link is soft-deleted.
func (l *Link) IsDeleted() bool {
	return l.DeletedAt != nil
}

// TagsJSON returns tags as JSON for database storage.
func (l *Link) TagsJSON() (string, error) {
	if l.Tags == nil {
		return "[]", nil
	}
	data, err := json.Marshal(l.Tags)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// ParseTagsJSON parses JSON tags from database.
func (l *Link) ParseTagsJSON(data string) error {
	if data == "" || data == "null" {
		l.Tags = nil
		return nil
	}
	return json.Unmarshal([]byte(data), &l.Tags)
}

// LinkPreview contains Open Graph metadata for a link.
type LinkPreview struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	ImageURL    string `json:"image_url,omitempty"`
	FetchedAt   time.Time `json:"fetched_at"`
}

// CreateLinkRequest contains data for creating a new link.
type CreateLinkRequest struct {
	URL       string   `json:"url"`
	Slug      string   `json:"slug,omitempty"`
	Domain    string   `json:"domain,omitempty"`
	Password  string   `json:"password,omitempty"`
	TTLHours  int      `json:"ttl_hours,omitempty"`
	Tags      []string `json:"tags,omitempty"`
	FolderID  *int64   `json:"folder_id,omitempty"`
	IsOneTime bool     `json:"is_one_time,omitempty"`
}

// UpdateLinkRequest contains data for updating an existing link.
type UpdateLinkRequest struct {
	URL      *string   `json:"url,omitempty"`
	Password *string   `json:"password,omitempty"`
	TTLHours *int      `json:"ttl_hours,omitempty"`
	Tags     *[]string `json:"tags,omitempty"`
	FolderID *int64    `json:"folder_id,omitempty"`
}

// ListLinksFilter contains filter options for listing links.
type ListLinksFilter struct {
	Search    string   `json:"search,omitempty"`
	Tags      []string `json:"tags,omitempty"`
	FolderID  *int64   `json:"folder_id,omitempty"`
	Domain    string   `json:"domain,omitempty"`
	Limit     int      `json:"limit,omitempty"`
	Offset    int      `json:"offset,omitempty"`
	IncludeDeleted bool `json:"include_deleted,omitempty"`
}
