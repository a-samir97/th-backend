package database

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"thamaniyah/internal/domain"

	"gorm.io/gorm"
)

// MediaModel represents the database model for media
type MediaModel struct {
	ID          string         `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Title       string         `gorm:"not null;size:255"`
	Description string         `gorm:"type:text"`
	FilePath    string         `gorm:"not null;size:500"`
	FileSize    int64          `gorm:"not null"`
	Duration    int            `gorm:"default:0"`
	Format      string         `gorm:"size:50"`
	Type        string         `gorm:"not null;size:20"`
	Status      string         `gorm:"not null;size:20;default:'uploading'"`
	Tags        StringArray    `gorm:"type:jsonb"`
	CreatedAt   time.Time      `gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

// TableName returns the table name for MediaModel
func (MediaModel) TableName() string {
	return "media_files"
}

// ToDomain converts MediaModel to domain.Media
func (m *MediaModel) ToDomain() *domain.Media {
	return &domain.Media{
		ID:          m.ID,
		Title:       m.Title,
		Description: m.Description,
		FilePath:    m.FilePath,
		FileSize:    m.FileSize,
		Duration:    m.Duration,
		Format:      m.Format,
		Type:        domain.MediaType(m.Type),
		Status:      domain.MediaStatus(m.Status),
		Tags:        []string(m.Tags),
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

// FromDomain converts domain.Media to MediaModel
func (m *MediaModel) FromDomain(media *domain.Media) {
	m.ID = media.ID
	m.Title = media.Title
	m.Description = media.Description
	m.FilePath = media.FilePath
	m.FileSize = media.FileSize
	m.Duration = media.Duration
	m.Format = media.Format
	m.Type = string(media.Type)
	m.Status = string(media.Status)
	m.Tags = StringArray(media.Tags)
	m.CreatedAt = media.CreatedAt
	m.UpdatedAt = media.UpdatedAt
}

// StringArray is a custom type to handle []string in JSONB
type StringArray []string

// Value implements the driver.Valuer interface
func (sa StringArray) Value() (driver.Value, error) {
	if sa == nil {
		return nil, nil
	}
	return json.Marshal(sa)
}

// Scan implements the sql.Scanner interface
func (sa *StringArray) Scan(value interface{}) error {
	if value == nil {
		*sa = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan StringArray")
	}

	return json.Unmarshal(bytes, sa)
}

// GormDataType returns the data type for GORM
func (StringArray) GormDataType() string {
	return "jsonb"
}
