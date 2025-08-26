package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUploadRequest_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		request  UploadRequest
		expected bool
	}{
		{
			name: "valid video request",
			request: UploadRequest{
				Title:    "Test Video",
				Filename: "test.mp4",
				FileSize: 1024 * 1024, // 1MB
				Type:     TypeVideo,
			},
			expected: true,
		},
		{
			name: "valid podcast request",
			request: UploadRequest{
				Title:       "Test Podcast",
				Description: "A great podcast",
				Filename:    "test.mp3",
				FileSize:    512 * 1024, // 512KB
				Type:        TypePodcast,
			},
			expected: true,
		},
		{
			name: "missing title",
			request: UploadRequest{
				Filename: "test.mp4",
				FileSize: 1024 * 1024,
				Type:     TypeVideo,
			},
			expected: false,
		},
		{
			name: "empty title",
			request: UploadRequest{
				Title:    "",
				Filename: "test.mp4",
				FileSize: 1024 * 1024,
				Type:     TypeVideo,
			},
			expected: false,
		},
		{
			name: "missing filename",
			request: UploadRequest{
				Title:    "Test Video",
				FileSize: 1024 * 1024,
				Type:     TypeVideo,
			},
			expected: false,
		},
		{
			name: "empty filename",
			request: UploadRequest{
				Title:    "Test Video",
				Filename: "",
				FileSize: 1024 * 1024,
				Type:     TypeVideo,
			},
			expected: false,
		},
		{
			name: "zero file size",
			request: UploadRequest{
				Title:    "Test Video",
				Filename: "test.mp4",
				FileSize: 0,
				Type:     TypeVideo,
			},
			expected: false,
		},
		{
			name: "negative file size",
			request: UploadRequest{
				Title:    "Test Video",
				Filename: "test.mp4",
				FileSize: -1,
				Type:     TypeVideo,
			},
			expected: false,
		},
		{
			name: "invalid media type",
			request: UploadRequest{
				Title:    "Test Video",
				Filename: "test.mp4",
				FileSize: 1024 * 1024,
				Type:     MediaType("invalid"),
			},
			expected: false,
		},
		{
			name: "file too large",
			request: UploadRequest{
				Title:    "Large Video",
				Filename: "large.mp4",
				FileSize: 6 * 1024 * 1024 * 1024, // 6GB (over 5GB limit)
				Type:     TypeVideo,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.request.IsValid()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestUploadRequest_ToMedia(t *testing.T) {
	// Given
	request := &UploadRequest{
		Title:       "Test Video",
		Description: "A test video",
		Filename:    "test.mp4",
		FileSize:    1024 * 1024,
		Type:        TypeVideo,
	}
	id := "media-123"
	filePath := "/uploads/media-123.mp4"

	// When
	media := request.ToMedia(id, filePath)

	// Then
	assert.NotNil(t, media)
	assert.Equal(t, id, media.ID)
	assert.Equal(t, request.Title, media.Title)
	assert.Equal(t, request.Description, media.Description)
	assert.Equal(t, filePath, media.FilePath)
	assert.Equal(t, request.FileSize, media.FileSize)
	assert.Equal(t, request.Type, media.Type)
	assert.Equal(t, StatusUploading, media.Status)
	assert.False(t, media.CreatedAt.IsZero())
	assert.False(t, media.UpdatedAt.IsZero())
}

func TestUpdateMediaRequest_ApplyTo(t *testing.T) {
	// Given
	originalTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	media := &Media{
		ID:          "123",
		Title:       "Original Title",
		Description: "Original Description",
		UpdatedAt:   originalTime,
	}

	newTitle := "Updated Title"
	newDescription := "Updated Description"
	updateRequest := &UpdateMediaRequest{
		Title:       &newTitle,
		Description: &newDescription,
	}

	beforeUpdate := time.Now()

	// When
	updateRequest.ApplyTo(media)

	// Then
	assert.Equal(t, "Updated Title", media.Title)
	assert.Equal(t, "Updated Description", media.Description)
	assert.True(t, media.UpdatedAt.After(beforeUpdate))
	assert.True(t, media.UpdatedAt.After(originalTime))
}

func TestUpdateMediaRequest_ApplyTo_PartialUpdate(t *testing.T) {
	// Given
	media := &Media{
		ID:          "123",
		Title:       "Original Title",
		Description: "Original Description",
	}

	newTitle := "Updated Title"
	updateRequest := &UpdateMediaRequest{
		Title: &newTitle,
		// Description is nil, should not be updated
	}

	// When
	updateRequest.ApplyTo(media)

	// Then
	assert.Equal(t, "Updated Title", media.Title)
	assert.Equal(t, "Original Description", media.Description) // Unchanged
}

func TestUpdateMediaRequest_ApplyTo_EmptyUpdate(t *testing.T) {
	// Given
	originalTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	media := &Media{
		ID:          "123",
		Title:       "Original Title",
		Description: "Original Description",
		UpdatedAt:   originalTime,
	}

	updateRequest := &UpdateMediaRequest{
		// Both fields are nil
	}

	beforeUpdate := time.Now()

	// When
	updateRequest.ApplyTo(media)

	// Then
	assert.Equal(t, "Original Title", media.Title)
	assert.Equal(t, "Original Description", media.Description)
	assert.True(t, media.UpdatedAt.After(beforeUpdate)) // UpdatedAt should still be updated
}

func TestUploadURL_Structure(t *testing.T) {
	// Test that UploadURL has the expected fields
	uploadURL := UploadURL{
		MediaID:   "media-123",
		URL:       "https://example.com/upload",
		ExpiresAt: time.Now().Add(time.Hour),
	}

	assert.Equal(t, "media-123", uploadURL.MediaID)
	assert.Equal(t, "https://example.com/upload", uploadURL.URL)
	assert.False(t, uploadURL.ExpiresAt.IsZero())
}
