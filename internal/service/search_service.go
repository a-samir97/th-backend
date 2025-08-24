package service

import (
	"context"
	"thamaniyah/internal/domain"
)

// SearchService defines the contract for search business logic
type SearchService interface {
	// Search performs a search across media content
	Search(ctx context.Context, query *domain.SearchQuery) (*domain.SearchResult, error)
	
	// SearchWithFilters performs a search with additional filters
	SearchWithFilters(ctx context.Context, query *domain.SearchQuery, filters *domain.SearchFilters) (*domain.SearchResult, error)
	
	// GetSuggestions returns search suggestions
	GetSuggestions(ctx context.Context, partial string, limit int) ([]string, error)
	
	// IndexMedia indexes media content for search
	IndexMedia(ctx context.Context, media *domain.Media) error
	
	// RemoveFromIndex removes media from search index
	RemoveFromIndex(ctx context.Context, mediaID string) error
	
	// UpdateIndex updates media in search index
	UpdateIndex(ctx context.Context, media *domain.Media) error
}
