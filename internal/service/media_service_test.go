package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"thamaniyah/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockMediaRepository is a mock implementation of MediaRepository
type MockMediaRepository struct {
	mock.Mock
}

func (m *MockMediaRepository) Create(ctx context.Context, media *domain.Media) error {
	args := m.Called(ctx, media)
	return args.Error(0)
}

func (m *MockMediaRepository) GetByID(ctx context.Context, id string) (*domain.Media, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Media), args.Error(1)
}

func (m *MockMediaRepository) GetAll(ctx context.Context, limit, offset int) ([]*domain.Media, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Media), args.Error(1)
}

func (m *MockMediaRepository) Update(ctx context.Context, media *domain.Media) error {
	args := m.Called(ctx, media)
	return args.Error(0)
}

func (m *MockMediaRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockMediaRepository) GetByStatus(ctx context.Context, status domain.MediaStatus, limit, offset int) ([]*domain.Media, error) {
	args := m.Called(ctx, status, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Media), args.Error(1)
}

func (m *MockMediaRepository) UpdateStatus(ctx context.Context, id string, status domain.MediaStatus) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockMediaRepository) GetTotal(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

func TestMediaService_CreateUploadURL(t *testing.T) {
	tests := []struct {
		name        string
		request     *domain.UploadRequest
		setupMock   func(*MockMediaRepository)
		expectError bool
		errorType   string
	}{
		{
			name: "successful upload URL creation",
			request: &domain.UploadRequest{
				Title:    "Test Video",
				Filename: "test.mp4",
				FileSize: 1024 * 1024, // 1MB
				Type:     domain.TypeVideo,
			},
			setupMock: func(mockRepo *MockMediaRepository) {
				mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Media")).Return(nil)
			},
			expectError: false,
		},
		{
			name: "invalid request - missing title",
			request: &domain.UploadRequest{
				Filename: "test.mp4",
				FileSize: 1024 * 1024,
				Type:     domain.TypeVideo,
			},
			setupMock: func(mockRepo *MockMediaRepository) {
				// No expectations as validation should fail before repository call
			},
			expectError: true,
			errorType:   "INVALID_REQUEST",
		},
		{
			name: "invalid request - file too large",
			request: &domain.UploadRequest{
				Title:    "Large Video",
				Filename: "large.mp4",
				FileSize: 6 * 1024 * 1024 * 1024, // 6GB
				Type:     domain.TypeVideo,
			},
			setupMock: func(mockRepo *MockMediaRepository) {
				// No expectations as validation should fail before repository call
			},
			expectError: true,
			errorType:   "INVALID_REQUEST",
		},
		{
			name: "repository error",
			request: &domain.UploadRequest{
				Title:    "Test Video",
				Filename: "test.mp4",
				FileSize: 1024 * 1024,
				Type:     domain.TypeVideo,
			},
			setupMock: func(mockRepo *MockMediaRepository) {
				mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Media")).
					Return(errors.New("database error"))
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			mockRepo := new(MockMediaRepository)
			tt.setupMock(mockRepo)
			service := NewMediaService(mockRepo)
			ctx := context.Background()

			// When
			result, err := service.CreateUploadURL(ctx, tt.request)

			// Then
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
				if tt.errorType != "" {
					var businessErr *domain.BusinessError
					if errors.As(err, &businessErr) {
						assert.Equal(t, tt.errorType, businessErr.Code)
					}
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.NotEmpty(t, result.MediaID)
				assert.NotEmpty(t, result.URL)
				assert.True(t, result.ExpiresAt.After(time.Now()))
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestMediaService_ConfirmUpload(t *testing.T) {
	tests := []struct {
		name        string
		mediaID     string
		setupMock   func(*MockMediaRepository)
		expectError bool
		errorType   string
	}{
		{
			name:    "successful upload confirmation",
			mediaID: "media-123",
			setupMock: func(mockRepo *MockMediaRepository) {
				media := &domain.Media{
					ID:     "media-123",
					Status: domain.StatusUploading,
				}
				mockRepo.On("GetByID", mock.Anything, "media-123").Return(media, nil)
				mockRepo.On("UpdateStatus", mock.Anything, "media-123", domain.StatusReady).Return(nil)
			},
			expectError: false,
		},
		{
			name:    "media not found",
			mediaID: "non-existent",
			setupMock: func(mockRepo *MockMediaRepository) {
				mockRepo.On("GetByID", mock.Anything, "non-existent").
					Return(nil, domain.ErrMediaNotFound)
			},
			expectError: true,
		},
		{
			name:    "invalid status - already ready",
			mediaID: "media-123",
			setupMock: func(mockRepo *MockMediaRepository) {
				media := &domain.Media{
					ID:     "media-123",
					Status: domain.StatusReady,
				}
				mockRepo.On("GetByID", mock.Anything, "media-123").Return(media, nil)
			},
			expectError: true,
			errorType:   "INVALID_STATUS",
		},
		{
			name:    "update status error",
			mediaID: "media-123",
			setupMock: func(mockRepo *MockMediaRepository) {
				media := &domain.Media{
					ID:     "media-123",
					Status: domain.StatusUploading,
				}
				mockRepo.On("GetByID", mock.Anything, "media-123").Return(media, nil)
				mockRepo.On("UpdateStatus", mock.Anything, "media-123", domain.StatusReady).
					Return(errors.New("database error"))
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			mockRepo := new(MockMediaRepository)
			tt.setupMock(mockRepo)
			service := NewMediaService(mockRepo)
			ctx := context.Background()

			// When
			err := service.ConfirmUpload(ctx, tt.mediaID)

			// Then
			if tt.expectError {
				assert.Error(t, err)
				if tt.errorType != "" {
					var businessErr *domain.BusinessError
					if errors.As(err, &businessErr) {
						assert.Equal(t, tt.errorType, businessErr.Code)
					}
				}
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestMediaService_GetMedia(t *testing.T) {
	tests := []struct {
		name      string
		mediaID   string
		setupMock func(*MockMediaRepository)
		expected  *domain.Media
		wantErr   bool
	}{
		{
			name:    "successful get media",
			mediaID: "media-123",
			setupMock: func(mockRepo *MockMediaRepository) {
				expectedMedia := &domain.Media{
					ID:    "media-123",
					Title: "Test Video",
					Type:  domain.TypeVideo,
				}
				mockRepo.On("GetByID", mock.Anything, "media-123").Return(expectedMedia, nil)
			},
			expected: &domain.Media{
				ID:    "media-123",
				Title: "Test Video",
				Type:  domain.TypeVideo,
			},
			wantErr: false,
		},
		{
			name:    "media not found",
			mediaID: "non-existent",
			setupMock: func(mockRepo *MockMediaRepository) {
				mockRepo.On("GetByID", mock.Anything, "non-existent").
					Return(nil, domain.ErrMediaNotFound)
			},
			expected: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			mockRepo := new(MockMediaRepository)
			tt.setupMock(mockRepo)
			service := NewMediaService(mockRepo)
			ctx := context.Background()

			// When
			result, err := service.GetMedia(ctx, tt.mediaID)

			// Then
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestMediaService_GetAllMedia(t *testing.T) {
	tests := []struct {
		name      string
		limit     int
		offset    int
		setupMock func(*MockMediaRepository)
		wantErr   bool
	}{
		{
			name:   "successful get all media",
			limit:  20,
			offset: 0,
			setupMock: func(mockRepo *MockMediaRepository) {
				expectedMedia := []*domain.Media{
					{ID: "media-1", Title: "Video 1"},
					{ID: "media-2", Title: "Video 2"},
				}
				mockRepo.On("GetAll", mock.Anything, 20, 0).Return(expectedMedia, nil)
				mockRepo.On("GetTotal", mock.Anything).Return(int64(100), nil)
			},
			wantErr: false,
		},
		{
			name:   "invalid limit - use default",
			limit:  0,
			offset: 0,
			setupMock: func(mockRepo *MockMediaRepository) {
				expectedMedia := []*domain.Media{}
				mockRepo.On("GetAll", mock.Anything, domain.DefaultPageSize, 0).Return(expectedMedia, nil)
				mockRepo.On("GetTotal", mock.Anything).Return(int64(0), nil)
			},
			wantErr: false,
		},
		{
			name:   "limit too large - use default",
			limit:  200,
			offset: 0,
			setupMock: func(mockRepo *MockMediaRepository) {
				expectedMedia := []*domain.Media{}
				mockRepo.On("GetAll", mock.Anything, domain.DefaultPageSize, 0).Return(expectedMedia, nil)
				mockRepo.On("GetTotal", mock.Anything).Return(int64(0), nil)
			},
			wantErr: false,
		},
		{
			name:   "negative offset - use zero",
			limit:  20,
			offset: -10,
			setupMock: func(mockRepo *MockMediaRepository) {
				expectedMedia := []*domain.Media{}
				mockRepo.On("GetAll", mock.Anything, 20, 0).Return(expectedMedia, nil)
				mockRepo.On("GetTotal", mock.Anything).Return(int64(0), nil)
			},
			wantErr: false,
		},
		{
			name:   "repository error",
			limit:  20,
			offset: 0,
			setupMock: func(mockRepo *MockMediaRepository) {
				mockRepo.On("GetAll", mock.Anything, 20, 0).
					Return(nil, errors.New("database error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			mockRepo := new(MockMediaRepository)
			tt.setupMock(mockRepo)
			service := NewMediaService(mockRepo)
			ctx := context.Background()

			// When
			mediaList, total, err := service.GetAllMedia(ctx, tt.limit, tt.offset)

			// Then
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, mediaList)
				assert.Zero(t, total)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, mediaList)
				assert.GreaterOrEqual(t, total, int64(0))
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestMediaService_UpdateMedia(t *testing.T) {
	tests := []struct {
		name        string
		mediaID     string
		request     *domain.UpdateMediaRequest
		setupMock   func(*MockMediaRepository)
		expectError bool
	}{
		{
			name:    "successful media update",
			mediaID: "media-123",
			request: &domain.UpdateMediaRequest{
				Title:       stringPtr("Updated Title"),
				Description: stringPtr("Updated Description"),
			},
			setupMock: func(mockRepo *MockMediaRepository) {
				originalMedia := &domain.Media{
					ID:          "media-123",
					Title:       "Original Title",
					Description: "Original Description",
				}
				mockRepo.On("GetByID", mock.Anything, "media-123").Return(originalMedia, nil)
				mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.Media")).Return(nil)
			},
			expectError: false,
		},
		{
			name:    "media not found",
			mediaID: "non-existent",
			request: &domain.UpdateMediaRequest{
				Title: stringPtr("Updated Title"),
			},
			setupMock: func(mockRepo *MockMediaRepository) {
				mockRepo.On("GetByID", mock.Anything, "non-existent").
					Return(nil, domain.ErrMediaNotFound)
			},
			expectError: true,
		},
		{
			name:    "update error",
			mediaID: "media-123",
			request: &domain.UpdateMediaRequest{
				Title: stringPtr("Updated Title"),
			},
			setupMock: func(mockRepo *MockMediaRepository) {
				originalMedia := &domain.Media{
					ID:    "media-123",
					Title: "Original Title",
				}
				mockRepo.On("GetByID", mock.Anything, "media-123").Return(originalMedia, nil)
				mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.Media")).
					Return(errors.New("database error"))
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			mockRepo := new(MockMediaRepository)
			tt.setupMock(mockRepo)
			service := NewMediaService(mockRepo)
			ctx := context.Background()

			// When
			result, err := service.UpdateMedia(ctx, tt.mediaID, tt.request)

			// Then
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				if tt.request.Title != nil {
					assert.Equal(t, *tt.request.Title, result.Title)
				}
				if tt.request.Description != nil {
					assert.Equal(t, *tt.request.Description, result.Description)
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestMediaService_DeleteMedia(t *testing.T) {
	tests := []struct {
		name        string
		mediaID     string
		setupMock   func(*MockMediaRepository)
		expectError bool
	}{
		{
			name:    "successful media deletion",
			mediaID: "media-123",
			setupMock: func(mockRepo *MockMediaRepository) {
				existingMedia := &domain.Media{
					ID:    "media-123",
					Title: "Test Video",
				}
				mockRepo.On("GetByID", mock.Anything, "media-123").Return(existingMedia, nil)
				mockRepo.On("Delete", mock.Anything, "media-123").Return(nil)
			},
			expectError: false,
		},
		{
			name:    "media not found",
			mediaID: "non-existent",
			setupMock: func(mockRepo *MockMediaRepository) {
				mockRepo.On("GetByID", mock.Anything, "non-existent").
					Return(nil, domain.ErrMediaNotFound)
			},
			expectError: true,
		},
		{
			name:    "delete error",
			mediaID: "media-123",
			setupMock: func(mockRepo *MockMediaRepository) {
				existingMedia := &domain.Media{
					ID:    "media-123",
					Title: "Test Video",
				}
				mockRepo.On("GetByID", mock.Anything, "media-123").Return(existingMedia, nil)
				mockRepo.On("Delete", mock.Anything, "media-123").Return(errors.New("database error"))
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			mockRepo := new(MockMediaRepository)
			tt.setupMock(mockRepo)
			service := NewMediaService(mockRepo)
			ctx := context.Background()

			// When
			err := service.DeleteMedia(ctx, tt.mediaID)

			// Then
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestMediaService_ProcessMedia(t *testing.T) {
	tests := []struct {
		name        string
		mediaID     string
		setupMock   func(*MockMediaRepository)
		expectError bool
	}{
		{
			name:    "successful media processing",
			mediaID: "media-123",
			setupMock: func(mockRepo *MockMediaRepository) {
				media := &domain.Media{
					ID:       "media-123",
					Type:     domain.TypeVideo,
					Duration: 0, // Will be set by extractMetadata
				}
				mockRepo.On("GetByID", mock.Anything, "media-123").Return(media, nil)
				mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.Media")).Return(nil)
			},
			expectError: false,
		},
		{
			name:    "media not found",
			mediaID: "non-existent",
			setupMock: func(mockRepo *MockMediaRepository) {
				mockRepo.On("GetByID", mock.Anything, "non-existent").
					Return(nil, domain.ErrMediaNotFound)
			},
			expectError: true,
		},
		{
			name:    "update error after processing",
			mediaID: "media-123",
			setupMock: func(mockRepo *MockMediaRepository) {
				media := &domain.Media{
					ID:   "media-123",
					Type: domain.TypeVideo,
				}
				mockRepo.On("GetByID", mock.Anything, "media-123").Return(media, nil)
				mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.Media")).
					Return(errors.New("database error"))
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			mockRepo := new(MockMediaRepository)
			tt.setupMock(mockRepo)
			service := NewMediaService(mockRepo)
			ctx := context.Background()

			// When
			err := service.ProcessMedia(ctx, tt.mediaID)

			// Then
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

// Helper function to create string pointers
func stringPtr(s string) *string {
	return &s
}
