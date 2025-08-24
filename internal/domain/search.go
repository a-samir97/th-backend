package domain

import (
	"strconv"
	"strings"
)

// SearchQuery represents a search request
type SearchQuery struct {
	Query    string    `json:"query"`
	Type     MediaType `json:"type,omitempty"`
	Tags     []string  `json:"tags,omitempty"`
	Limit    int       `json:"limit"`
	Offset   int       `json:"offset"`
	SortBy   string    `json:"sort_by"`   // title, created_at, duration, file_size
	SortDesc bool      `json:"sort_desc"` // true for descending
}

// NewSearchQuery creates a search query with default values
func NewSearchQuery() *SearchQuery {
	return &SearchQuery{
		Limit:    20,
		Offset:   0,
		SortBy:   "created_at",
		SortDesc: true,
	}
}

// IsValid validates the search query
func (sq *SearchQuery) IsValid() bool {
	if sq.Limit <= 0 || sq.Limit > 100 {
		sq.Limit = 20
	}
	if sq.Offset < 0 {
		sq.Offset = 0
	}
	
	// Validate sort field
	validSortFields := map[string]bool{
		"title":      true,
		"created_at": true,
		"duration":   true,
		"file_size":  true,
	}
	
	if !validSortFields[sq.SortBy] {
		sq.SortBy = "created_at"
	}
	
	return true
}

// SearchResult represents the result of a search operation
type SearchResult struct {
	Items      []*Media `json:"items"`
	Total      int64    `json:"total"`
	Limit      int      `json:"limit"`
	Offset     int      `json:"offset"`
	Query      string   `json:"query"`
	SearchTime string   `json:"search_time"` // e.g., "15ms"
}

// SearchDocument represents a media document in Elasticsearch
type SearchDocument struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Type        MediaType `json:"type"`
	Tags        []string  `json:"tags"`
	Duration    int       `json:"duration"`
	FileSize    int64     `json:"file_size"`
	Format      string    `json:"format"`
	CreatedAt   string    `json:"created_at"`
}

// FromMedia converts Media entity to SearchDocument
func (sd *SearchDocument) FromMedia(media *Media) {
	sd.ID = media.ID
	sd.Title = media.Title
	sd.Description = media.Description
	sd.Type = media.Type
	sd.Tags = media.Tags
	sd.Duration = media.Duration
	sd.FileSize = media.FileSize
	sd.Format = media.Format
	sd.CreatedAt = media.CreatedAt.Format("2006-01-02T15:04:05Z")
}

// ToMedia converts SearchDocument back to Media entity (partial)
func (sd *SearchDocument) ToMedia() *Media {
	return &Media{
		ID:          sd.ID,
		Title:       sd.Title,
		Description: sd.Description,
		Type:        sd.Type,
		Tags:        sd.Tags,
		Duration:    sd.Duration,
		FileSize:    sd.FileSize,
		Format:      sd.Format,
		Status:      StatusReady, // Search results are always ready
	}
}

// SearchFilters represents additional search filters
type SearchFilters struct {
	MinDuration int    `json:"min_duration,omitempty"` // in seconds
	MaxDuration int    `json:"max_duration,omitempty"` // in seconds
	MinFileSize int64  `json:"min_file_size,omitempty"`
	MaxFileSize int64  `json:"max_file_size,omitempty"`
	Format      string `json:"format,omitempty"`
}

// ParseFiltersFromQuery parses filters from query parameters
func ParseFiltersFromQuery(params map[string]string) *SearchFilters {
	filters := &SearchFilters{}
	
	if minDur := params["min_duration"]; minDur != "" {
		if val, err := strconv.Atoi(minDur); err == nil {
			filters.MinDuration = val
		}
	}
	
	if maxDur := params["max_duration"]; maxDur != "" {
		if val, err := strconv.Atoi(maxDur); err == nil {
			filters.MaxDuration = val
		}
	}
	
	if minSize := params["min_file_size"]; minSize != "" {
		if val, err := strconv.ParseInt(minSize, 10, 64); err == nil {
			filters.MinFileSize = val
		}
	}
	
	if maxSize := params["max_file_size"]; maxSize != "" {
		if val, err := strconv.ParseInt(maxSize, 10, 64); err == nil {
			filters.MaxFileSize = val
		}
	}
	
	if format := params["format"]; format != "" {
		filters.Format = strings.ToLower(format)
	}
	
	return filters
}
