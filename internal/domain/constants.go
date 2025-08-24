package domain

import (
	"time"
)

// Constants for business logic
const (
	// File size limits
	MaxVideoFileSize   = 5 * 1024 * 1024 * 1024 // 5GB
	MaxPodcastFileSize = 1 * 1024 * 1024 * 1024 // 1GB

	// Upload URL expiration
	UploadURLTTL = 1 * time.Hour

	// Search limits
	MaxSearchLimit     = 100
	DefaultSearchLimit = 20

	// Pagination
	DefaultPageSize = 20
	MaxPageSize     = 100
)

// Supported file formats
var (
	VideoFormats = []string{"mp4", "mov", "avi", "mkv", "webm"}
	AudioFormats = []string{"mp3", "wav", "flac", "aac", "ogg"}
)

// IsValidVideoFormat checks if the format is a valid video format
func IsValidVideoFormat(format string) bool {
	for _, validFormat := range VideoFormats {
		if format == validFormat {
			return true
		}
	}
	return false
}

// IsValidAudioFormat checks if the format is a valid audio format
func IsValidAudioFormat(format string) bool {
	for _, validFormat := range AudioFormats {
		if format == validFormat {
			return true
		}
	}
	return false
}

// IsValidFormat checks if the format is valid for the given media type
func IsValidFormat(mediaType MediaType, format string) bool {
	switch mediaType {
	case TypeVideo:
		return IsValidVideoFormat(format)
	case TypePodcast:
		return IsValidAudioFormat(format)
	default:
		return false
	}
}

// GetMaxFileSize returns the maximum allowed file size for a media type
func GetMaxFileSize(mediaType MediaType) int64 {
	switch mediaType {
	case TypeVideo:
		return MaxVideoFileSize
	case TypePodcast:
		return MaxPodcastFileSize
	default:
		return 0
	}
}

// Event represents a domain event for async processing
type Event struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Data      map[string]interface{} `json:"data"`
	Timestamp time.Time              `json:"timestamp"`
}

// Event types
const (
	EventMediaUploaded  = "media.uploaded"
	EventMediaProcessed = "media.processed"
	EventMediaDeleted   = "media.deleted"
	EventMediaUpdated   = "media.updated"
)

// NewEvent creates a new domain event
func NewEvent(eventType string, data map[string]interface{}) *Event {
	return &Event{
		ID:        generateID(), // We'll implement this later
		Type:      eventType,
		Data:      data,
		Timestamp: time.Now(),
	}
}

// Temporary ID generator (will be replaced with proper UUID)
func generateID() string {
	return time.Now().Format("20060102150405")
}
