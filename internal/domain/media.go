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
	ID          string      `json:"id" gorm:"primaryKey"`
	Title       string      `json:"title" gorm:"not null"`
	Description string      `json:"description"`
	FilePath    string      `json:"file_path"`
	FileSize    int64       `json:"file_size"`
	Duration    int         `json:"duration"` // in seconds
	Format      string      `json:"format"`   // mp4, mp3, etc
	Type        MediaType   `json:"type" gorm:"type:varchar(20)"`
	Status      MediaStatus `json:"status" gorm:"type:varchar(20)"`
	CreatedAt   time.Time   `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time   `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   *time.Time  `json:"deleted_at,omitempty" gorm:"index"`
}

// TableName specifies the table name for Media
func (Media) TableName() string {
	return "media_files"
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
