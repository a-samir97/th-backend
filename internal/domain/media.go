package domain

import (
	"time"
)

// MediaStatus represents the current state of a media file
type MediaStatus string

const (
	StatusUploading MediaStatus = "uploading"
	StatusReady     MediaStatus = "ready"
	StatusFailed    MediaStatus = "failed"
	StatusDeleted   MediaStatus = "deleted"
)

// MediaType represents the type of media content
type MediaType string

const (
	TypeVideo   MediaType = "video"
	TypePodcast MediaType = "podcast"
)

// Media represents a media file entity
type Media struct {
	ID          string      `json:"id" db:"id"`
	Title       string      `json:"title" db:"title"`
	Description string      `json:"description" db:"description"`
	FilePath    string      `json:"file_path" db:"file_path"`
	FileSize    int64       `json:"file_size" db:"file_size"`
	Duration    int         `json:"duration" db:"duration"` // in seconds
	Format      string      `json:"format" db:"format"`     // mp4, mp3, etc
	Type        MediaType   `json:"type" db:"type"`
	Status      MediaStatus `json:"status" db:"status"`
	Tags        []string    `json:"tags" db:"tags"`
	CreatedAt   time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at" db:"updated_at"`
}

// IsValid checks if the media entity has valid required fields
func (m *Media) IsValid() bool {
	return m.Title != "" && m.Type != "" && m.Status != ""
}

// IsProcessed returns true if the media has been successfully processed
func (m *Media) IsProcessed() bool {
	return m.Status == StatusReady
}

// CanBeSearched returns true if the media can appear in search results
func (m *Media) CanBeSearched() bool {
	return m.Status == StatusReady
}

// UpdateStatus updates the media status and timestamp
func (m *Media) UpdateStatus(status MediaStatus) {
	m.Status = status
	m.UpdatedAt = time.Now()
}
