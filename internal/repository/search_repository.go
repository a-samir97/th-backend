package repository

import (
	"context"
	"thamaniyah/internal/domain"
)

// SearchRepository defines the contract for search operations
type SearchRepository interface {
	// Index indexes a media document in the search engine
	Index(ctx context.Context, media *domain.Media) error
	
	// Search performs a search query and returns results
	Search(ctx context.Context, query *domain.SearchQuery) (*domain.SearchResult, error)
	
	// SearchWithFilters performs a search with additional filters
	SearchWithFilters(ctx context.Context, query *domain.SearchQuery, filters *domain.SearchFilters) (*domain.SearchResult, error)
	
	// Delete removes a document from the search index
	Delete(ctx context.Context, id string) error
	
	// Update updates a document in the search index
	Update(ctx context.Context, media *domain.Media) error
	
	// GetSuggestions returns search suggestions based on partial query
	GetSuggestions(ctx context.Context, partial string, limit int) ([]string, error)
	
	// IsHealthy checks if the search service is healthy
	IsHealthy(ctx context.Context) error
}
