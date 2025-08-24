package service

import (
	"context"
	"thamaniyah/internal/domain"
)

// MediaService defines the contract for media business logic
type MediaService interface {
	// CreateUploadURL generates a presigned URL for media upload
	CreateUploadURL(ctx context.Context, req *domain.UploadRequest) (*domain.UploadURL, error)

	// ConfirmUpload confirms that a file has been uploaded successfully
	ConfirmUpload(ctx context.Context, mediaID string) error

	// GetMedia retrieves a media record by ID
	GetMedia(ctx context.Context, id string) (*domain.Media, error)

	// GetAllMedia retrieves all media records with pagination
	GetAllMedia(ctx context.Context, limit, offset int) ([]*domain.Media, int64, error)

	// UpdateMedia updates media metadata
	UpdateMedia(ctx context.Context, id string, req *domain.UpdateMediaRequest) (*domain.Media, error)

	// DeleteMedia soft deletes a media record
	DeleteMedia(ctx context.Context, id string) error

	// ProcessMedia processes uploaded media (extract metadata, etc.)
	ProcessMedia(ctx context.Context, mediaID string) error
}
