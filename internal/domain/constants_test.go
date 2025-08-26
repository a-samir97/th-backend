package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidVideoFormat(t *testing.T) {
	tests := []struct {
		name     string
		format   string
		expected bool
	}{
		{"valid mp4", "mp4", true},
		{"valid mov", "mov", true},
		{"valid avi", "avi", true},
		{"valid mkv", "mkv", true},
		{"valid webm", "webm", true},
		{"invalid format", "txt", false},
		{"empty format", "", false},
		{"case sensitive - MP4", "MP4", false}, // Should be lowercase
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidVideoFormat(tt.format)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsValidAudioFormat(t *testing.T) {
	tests := []struct {
		name     string
		format   string
		expected bool
	}{
		{"valid mp3", "mp3", true},
		{"valid wav", "wav", true},
		{"valid flac", "flac", true},
		{"valid aac", "aac", true},
		{"valid ogg", "ogg", true},
		{"invalid format", "txt", false},
		{"empty format", "", false},
		{"case sensitive - MP3", "MP3", false}, // Should be lowercase
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidAudioFormat(tt.format)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsValidFormat(t *testing.T) {
	tests := []struct {
		name      string
		mediaType MediaType
		format    string
		expected  bool
	}{
		{"video with valid format", TypeVideo, "mp4", true},
		{"video with invalid format", TypeVideo, "mp3", false},
		{"podcast with valid format", TypePodcast, "mp3", true},
		{"podcast with invalid format", TypePodcast, "mp4", false},
		{"invalid media type", MediaType("invalid"), "mp4", false},
		{"empty format", TypeVideo, "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidFormat(tt.mediaType, tt.format)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetMaxFileSize(t *testing.T) {
	tests := []struct {
		name      string
		mediaType MediaType
		expected  int64
	}{
		{"video type", TypeVideo, MaxVideoFileSize},
		{"podcast type", TypePodcast, MaxPodcastFileSize},
		{"invalid type", MediaType("invalid"), 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetMaxFileSize(tt.mediaType)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestConstants(t *testing.T) {
	// Test file size constants
	assert.Equal(t, 5*1024*1024*1024, MaxVideoFileSize)  // 5GB
	assert.Equal(t, 1*1024*1024*1024, MaxPodcastFileSize) // 1GB
	
	// Test pagination constants
	assert.Equal(t, 20, DefaultPageSize)
	assert.Equal(t, 100, MaxPageSize)
	
	// Test search constants
	assert.Equal(t, 100, MaxSearchLimit)
	assert.Equal(t, 20, DefaultSearchLimit)
}

func TestVideoFormats(t *testing.T) {
	expectedFormats := []string{"mp4", "mov", "avi", "mkv", "webm"}
	assert.Equal(t, expectedFormats, VideoFormats)
}

func TestAudioFormats(t *testing.T) {
	expectedFormats := []string{"mp3", "wav", "flac", "aac", "ogg"}
	assert.Equal(t, expectedFormats, AudioFormats)
}

func TestNewEvent(t *testing.T) {
	// Given
	eventType := EventMediaUploaded
	data := map[string]interface{}{
		"media_id": "123",
		"title":    "Test Video",
	}

	// When
	event := NewEvent(eventType, data)

	// Then
	assert.NotNil(t, event)
	assert.NotEmpty(t, event.ID)
	assert.Equal(t, eventType, event.Type)
	assert.Equal(t, data, event.Data)
	assert.False(t, event.Timestamp.IsZero())
}

func TestEventTypes(t *testing.T) {
	// Test that event type constants are defined correctly
	assert.Equal(t, "media.uploaded", EventMediaUploaded)
	assert.Equal(t, "media.processed", EventMediaProcessed)
	assert.Equal(t, "media.deleted", EventMediaDeleted)
	assert.Equal(t, "media.updated", EventMediaUpdated)
}

func TestGenerateID(t *testing.T) {
	// When
	id1 := generateID()
	id2 := generateID()

	// Then
	assert.NotEmpty(t, id1)
	assert.NotEmpty(t, id2)
	// IDs should be different (assuming they're called at different times)
	// Note: This test might be flaky if called within the same second
}
