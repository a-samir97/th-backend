package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMedia_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		media    Media
		expected bool
	}{
		{
			name: "valid media",
			media: Media{
				ID:     "123",
				Title:  "Test Video",
				Type:   TypeVideo,
				Status: StatusReady,
			},
			expected: true,
		},
		{
			name: "missing title",
			media: Media{
				ID:     "123",
				Type:   TypeVideo,
				Status: StatusReady,
			},
			expected: false,
		},
		{
			name: "missing type",
			media: Media{
				ID:     "123",
				Title:  "Test Video",
				Status: StatusReady,
			},
			expected: false,
		},
		{
			name: "missing status",
			media: Media{
				ID:    "123",
				Title: "Test Video",
				Type:  TypeVideo,
			},
			expected: false,
		},
		{
			name: "empty title",
			media: Media{
				ID:     "123",
				Title:  "",
				Type:   TypeVideo,
				Status: StatusReady,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.media.IsValid()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMedia_IsProcessed(t *testing.T) {
	tests := []struct {
		name     string
		status   MediaStatus
		expected bool
	}{
		{
			name:     "ready status",
			status:   StatusReady,
			expected: true,
		},
		{
			name:     "uploading status",
			status:   StatusUploading,
			expected: false,
		},
		{
			name:     "failed status",
			status:   StatusFailed,
			expected: false,
		},
		{
			name:     "deleted status",
			status:   StatusDeleted,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			media := Media{Status: tt.status}
			result := media.IsProcessed()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMedia_CanBeSearched(t *testing.T) {
	tests := []struct {
		name     string
		status   MediaStatus
		expected bool
	}{
		{
			name:     "ready status",
			status:   StatusReady,
			expected: true,
		},
		{
			name:     "uploading status",
			status:   StatusUploading,
			expected: false,
		},
		{
			name:     "failed status",
			status:   StatusFailed,
			expected: false,
		},
		{
			name:     "deleted status",
			status:   StatusDeleted,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			media := Media{Status: tt.status}
			result := media.CanBeSearched()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMedia_UpdateStatus(t *testing.T) {
	// Given
	media := &Media{
		ID:        "123",
		Status:    StatusUploading,
		UpdatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	beforeTime := time.Now()

	// When
	media.UpdateStatus(StatusReady)

	// Then
	assert.Equal(t, StatusReady, media.Status)
	assert.True(t, media.UpdatedAt.After(beforeTime))
}

func TestMedia_TableName(t *testing.T) {
	// Given
	media := Media{}

	// When
	tableName := media.TableName()

	// Then
	assert.Equal(t, "media_files", tableName)
}

func TestMediaType_Constants(t *testing.T) {
	// Test that constants are defined correctly
	assert.Equal(t, MediaType("video"), TypeVideo)
	assert.Equal(t, MediaType("podcast"), TypePodcast)
}

func TestMediaStatus_Constants(t *testing.T) {
	// Test that constants are defined correctly
	assert.Equal(t, MediaStatus("uploading"), StatusUploading)
	assert.Equal(t, MediaStatus("ready"), StatusReady)
	assert.Equal(t, MediaStatus("failed"), StatusFailed)
	assert.Equal(t, MediaStatus("deleted"), StatusDeleted)
}
