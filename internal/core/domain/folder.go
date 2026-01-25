package domain

import "time"

type Folder struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	ParentID  *int64    `json:"parent_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateFolderRequest struct {
	Name     string `json:"name"`
	ParentID *int64 `json:"parent_id,omitempty"`
}
