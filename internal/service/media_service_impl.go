package service

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"thamaniyah/internal/domain"
	"thamaniyah/internal/repository"

	"github.com/google/uuid"
)

// mediaService implements MediaService interface
type mediaService struct {
	mediaRepo repository.MediaRepository
}

// NewMediaService creates a new media service
func NewMediaService(mediaRepo repository.MediaRepository) MediaService {
	return &mediaService{
		mediaRepo: mediaRepo,
	}
}

// CreateUploadURL generates a presigned URL for media upload
func (s *mediaService) CreateUploadURL(ctx context.Context, req *domain.UploadRequest) (*domain.UploadURL, error) {
	// Validate the request
	if !req.IsValid() {
		return nil, domain.NewBusinessError("INVALID_REQUEST", "Upload request validation failed")
	}

	// Generate unique media ID
	mediaID := uuid.New().String()

	// Generate file path (in production this would be S3 path)
	filePath := s.generateFilePath(req.Filename, mediaID)

	// Create media record in uploading state
	media := req.ToMedia(mediaID, filePath)
	if err := s.mediaRepo.Create(ctx, media); err != nil {
		return nil, fmt.Errorf("failed to create media record: %w", err)
	}

	// Generate upload URL (simulated presigned URL for local development)
	uploadURL := s.generateUploadURL(filePath)

	return &domain.UploadURL{
		MediaID:   mediaID,
		URL:       uploadURL,
		ExpiresAt: time.Now().Add(domain.UploadURLTTL),
	}, nil
}

// ConfirmUpload confirms that a file has been uploaded successfully
func (s *mediaService) ConfirmUpload(ctx context.Context, mediaID string) error {
	// Get the media record
	media, err := s.mediaRepo.GetByID(ctx, mediaID)
	if err != nil {
		return err
	}

	// Check if media is in uploading state
	if media.Status != domain.StatusUploading {
		return domain.NewBusinessError("INVALID_STATUS",
			fmt.Sprintf("Media is in %s state, expected uploading", media.Status))
	}

	// In production, we would verify the file exists in S3
	// For now, we'll simulate file validation
	if err := s.validateUploadedFile(media.FilePath); err != nil {
		// Mark as failed
		s.mediaRepo.UpdateStatus(ctx, mediaID, domain.StatusFailed)
		return err
	}

	// Update status to ready
	if err := s.mediaRepo.UpdateStatus(ctx, mediaID, domain.StatusReady); err != nil {
		return fmt.Errorf("failed to update media status: %w", err)
	}

	return nil
}

// GetMedia retrieves a media record by ID
func (s *mediaService) GetMedia(ctx context.Context, id string) (*domain.Media, error) {
	return s.mediaRepo.GetByID(ctx, id)
}

// GetAllMedia retrieves all media records with pagination
func (s *mediaService) GetAllMedia(ctx context.Context, limit, offset int) ([]*domain.Media, int64, error) {
	// Validate pagination parameters
	if limit <= 0 || limit > domain.MaxPageSize {
		limit = domain.DefaultPageSize
	}
	if offset < 0 {
		offset = 0
	}

	// Get media records
	mediaList, err := s.mediaRepo.GetAll(ctx, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get media list: %w", err)
	}

	// Get total count
	total, err := s.mediaRepo.GetTotal(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get total count: %w", err)
	}

	return mediaList, total, nil
}

// UpdateMedia updates media metadata
func (s *mediaService) UpdateMedia(ctx context.Context, id string, req *domain.UpdateMediaRequest) (*domain.Media, error) {
	// Get existing media
	media, err := s.mediaRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Apply updates
	req.ApplyTo(media)

	// Update in database
	if err := s.mediaRepo.Update(ctx, media); err != nil {
		return nil, fmt.Errorf("failed to update media: %w", err)
	}

	return media, nil
}

// DeleteMedia soft deletes a media record
func (s *mediaService) DeleteMedia(ctx context.Context, id string) error {
	// Check if media exists
	_, err := s.mediaRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Soft delete the record
	if err := s.mediaRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete media: %w", err)
	}

	// In production, we would also delete the file from S3 here
	// For now, we just log it
	fmt.Printf("Media %s marked for deletion\n", id)

	return nil
}

// ProcessMedia processes uploaded media (extract metadata, etc.)
func (s *mediaService) ProcessMedia(ctx context.Context, mediaID string) error {
	media, err := s.mediaRepo.GetByID(ctx, mediaID)
	if err != nil {
		return err
	}

	// Simulate metadata extraction
	if err := s.extractMetadata(media); err != nil {
		// Mark as failed
		s.mediaRepo.UpdateStatus(ctx, mediaID, domain.StatusFailed)
		return fmt.Errorf("failed to extract metadata: %w", err)
	}

	// Update the media record with extracted metadata
	if err := s.mediaRepo.Update(ctx, media); err != nil {
		return fmt.Errorf("failed to update media with metadata: %w", err)
	}

	return nil
}

// Helper methods

// generateFilePath creates a file path for the uploaded media
func (s *mediaService) generateFilePath(filename, mediaID string) string {
	ext := filepath.Ext(filename)
	return fmt.Sprintf("/uploads/%s%s", mediaID, ext)
}

// generateUploadURL creates a presigned URL (simulated for local development)
func (s *mediaService) generateUploadURL(filePath string) string {
	// In production, this would generate an actual S3 presigned URL
	return fmt.Sprintf("http://localhost:8080/upload%s", filePath)
}

// validateUploadedFile validates that the uploaded file exists and is valid
func (s *mediaService) validateUploadedFile(filePath string) error {
	// In production, we would check if file exists in S3 and validate its size/format
	// For now, we'll just simulate validation
	fmt.Printf("Validating uploaded file: %s\n", filePath)
	return nil
}

// extractMetadata extracts metadata from the uploaded media file
func (s *mediaService) extractMetadata(media *domain.Media) error {
	// In production, we would use libraries like FFmpeg to extract video metadata
	// For now, we'll simulate metadata extraction

	switch media.Type {
	case domain.TypeVideo:
		// Simulate video metadata extraction
		if media.Duration == 0 {
			media.Duration = 300 // 5 minutes default
		}
	case domain.TypePodcast:
		// Simulate audio metadata extraction
		if media.Duration == 0 {
			media.Duration = 1800 // 30 minutes default
		}
	}

	fmt.Printf("Extracted metadata for media %s: duration=%d seconds\n", media.ID, media.Duration)
	return nil
}
