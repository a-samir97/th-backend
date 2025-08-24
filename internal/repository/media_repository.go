package repository

import (
	"context"
	"thamaniyah/internal/domain"
)

// MediaRepository defines the contract for media data access
type MediaRepository interface {
	// Create creates a new media record
	Create(ctx context.Context, media *domain.Media) error

	// GetByID retrieves a media record by ID
	GetByID(ctx context.Context, id string) (*domain.Media, error)

	// GetAll retrieves all media records with pagination
	GetAll(ctx context.Context, limit, offset int) ([]*domain.Media, error)

	// Update updates an existing media record
	Update(ctx context.Context, media *domain.Media) error

	// Delete soft deletes a media record by ID
	Delete(ctx context.Context, id string) error

	// GetByStatus retrieves media records by status
	GetByStatus(ctx context.Context, status domain.MediaStatus, limit, offset int) ([]*domain.Media, error)

	// UpdateStatus updates only the status of a media record
	UpdateStatus(ctx context.Context, id string, status domain.MediaStatus) error

	// GetTotal returns the total count of media records
	GetTotal(ctx context.Context) (int64, error)
}
