package domain

import (
	"time"
)

// UploadURL represents a presigned URL for file upload
type UploadURL struct {
	MediaID   string    `json:"media_id"`
	URL       string    `json:"upload_url"`
	ExpiresAt time.Time `json:"expires_at"`
}

// UploadRequest represents a request to initiate file upload
type UploadRequest struct {
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description"`
	Filename    string    `json:"filename" binding:"required"`
	FileSize    int64     `json:"file_size" binding:"required"`
	Type        MediaType `json:"type" binding:"required"`
	Tags        []string  `json:"tags"`
}

// IsValid validates the upload request
func (ur *UploadRequest) IsValid() bool {
	if ur.Title == "" || ur.Filename == "" || ur.FileSize <= 0 {
		return false
	}

	// Validate media type
	if ur.Type != TypeVideo && ur.Type != TypePodcast {
		return false
	}

	// Basic file size validation (max 5GB)
	maxFileSize := int64(5 * 1024 * 1024 * 1024) // 5GB
	if ur.FileSize > maxFileSize {
		return false
	}

	return true
}

// ToMedia converts UploadRequest to Media entity
func (ur *UploadRequest) ToMedia(id, filePath string) *Media {
	return &Media{
		ID:          id,
		Title:       ur.Title,
		Description: ur.Description,
		FilePath:    filePath,
		FileSize:    ur.FileSize,
		Type:        ur.Type,
		Status:      StatusUploading,
		Tags:        ur.Tags,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// UpdateMediaRequest represents a request to update media metadata
type UpdateMediaRequest struct {
	Title       *string  `json:"title,omitempty"`
	Description *string  `json:"description,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}

// ApplyTo applies the update request to a media entity
func (umr *UpdateMediaRequest) ApplyTo(media *Media) {
	if umr.Title != nil {
		media.Title = *umr.Title
	}
	if umr.Description != nil {
		media.Description = *umr.Description
	}
	if umr.Tags != nil {
		media.Tags = umr.Tags
	}
	media.UpdatedAt = time.Now()
}
