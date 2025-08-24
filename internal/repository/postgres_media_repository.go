package repository

import (
	"context"
	"errors"

	"thamaniyah/internal/domain"
	"thamaniyah/pkg/database"

	"gorm.io/gorm"
)

// postgresMediaRepository implements MediaRepository using PostgreSQL
type postgresMediaRepository struct {
	db *gorm.DB
}

// NewPostgresMediaRepository creates a new PostgreSQL media repository
func NewPostgresMediaRepository(conn *database.Connection) MediaRepository {
	return &postgresMediaRepository{
		db: conn.DB,
	}
}

// Create creates a new media record
func (r *postgresMediaRepository) Create(ctx context.Context, media *domain.Media) error {
	model := &database.MediaModel{}
	model.FromDomain(media)

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}

	// Update the media with the generated ID
	media.ID = model.ID
	media.CreatedAt = model.CreatedAt
	media.UpdatedAt = model.UpdatedAt

	return nil
}

// GetByID retrieves a media record by ID
func (r *postgresMediaRepository) GetByID(ctx context.Context, id string) (*domain.Media, error) {
	var model database.MediaModel

	err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrMediaNotFound
		}
		return nil, err
	}

	return model.ToDomain(), nil
}

// GetAll retrieves all media records with pagination
func (r *postgresMediaRepository) GetAll(ctx context.Context, limit, offset int) ([]*domain.Media, error) {
	var models []database.MediaModel

	err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&models).Error
	if err != nil {
		return nil, err
	}

	result := make([]*domain.Media, len(models))
	for i, model := range models {
		result[i] = model.ToDomain()
	}

	return result, nil
}

// Update updates an existing media record
func (r *postgresMediaRepository) Update(ctx context.Context, media *domain.Media) error {
	model := &database.MediaModel{}
	model.FromDomain(media)

	result := r.db.WithContext(ctx).
		Model(&database.MediaModel{}).
		Where("id = ?", media.ID).
		Updates(model)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return domain.ErrMediaNotFound
	}

	// Fetch the updated record to get the updated_at timestamp
	var updatedModel database.MediaModel
	if err := r.db.WithContext(ctx).Where("id = ?", media.ID).First(&updatedModel).Error; err == nil {
		media.UpdatedAt = updatedModel.UpdatedAt
	}

	return nil
}

// Delete soft deletes a media record by ID
func (r *postgresMediaRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&database.MediaModel{}, "id = ?", id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return domain.ErrMediaNotFound
	}

	return nil
}

// GetByStatus retrieves media records by status
func (r *postgresMediaRepository) GetByStatus(ctx context.Context, status domain.MediaStatus, limit, offset int) ([]*domain.Media, error) {
	var models []database.MediaModel

	err := r.db.WithContext(ctx).
		Where("status = ?", string(status)).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&models).Error
	if err != nil {
		return nil, err
	}

	result := make([]*domain.Media, len(models))
	for i, model := range models {
		result[i] = model.ToDomain()
	}

	return result, nil
}

// UpdateStatus updates only the status of a media record
func (r *postgresMediaRepository) UpdateStatus(ctx context.Context, id string, status domain.MediaStatus) error {
	result := r.db.WithContext(ctx).
		Model(&database.MediaModel{}).
		Where("id = ?", id).
		Update("status", string(status))

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return domain.ErrMediaNotFound
	}

	return nil
}

// GetTotal returns the total count of media records
func (r *postgresMediaRepository) GetTotal(ctx context.Context) (int64, error) {
	var count int64

	err := r.db.WithContext(ctx).
		Model(&database.MediaModel{}).
		Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}
