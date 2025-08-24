package service

import (
	"context"
	"encoding/json"
	"fmt"

	"thamaniyah/internal/domain"
	"thamaniyah/internal/repository"
	"thamaniyah/pkg/httpclient"
)

// SearchService defines search operations
type SearchService interface {
	// Search performs search with the given request
	Search(ctx context.Context, req *domain.SearchRequest) (*domain.SearchResponse, error)

	// Suggest provides search suggestions
	Suggest(ctx context.Context, req *domain.SuggestRequest) (*domain.SuggestResponse, error)

	// Reindex rebuilds the search index by fetching data from CMS service
	Reindex(ctx context.Context) error
}

// SearchServiceImpl implements SearchService
type SearchServiceImpl struct {
	searchRepo repository.SearchRepository
	cmsClient  *httpclient.Client
}

// NewSearchService creates a new search service
func NewSearchService(searchRepo repository.SearchRepository, cmsClient *httpclient.Client) SearchService {
	return &SearchServiceImpl{
		searchRepo: searchRepo,
		cmsClient:  cmsClient,
	}
}

// Search performs search operation
func (s *SearchServiceImpl) Search(ctx context.Context, req *domain.SearchRequest) (*domain.SearchResponse, error) {
	// Validate request
	if req.Query == "" {
		return nil, &domain.BusinessError{
			Code:    "INVALID_SEARCH_QUERY",
			Message: "Search query cannot be empty",
		}
	}

	// Set defaults
	if req.Limit <= 0 {
		req.Limit = 20
	}
	if req.Offset < 0 {
		req.Offset = 0
	}

	// Perform search
	results, total, err := s.searchRepo.Search(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("search failed: %w", err)
	}

	// Build response
	response := &domain.SearchResponse{
		Results: results,
		Total:   total,
		Query:   req.Query,
		Limit:   req.Limit,
		Offset:  req.Offset,
	}

	return response, nil
}

// Suggest provides search suggestions
func (s *SearchServiceImpl) Suggest(ctx context.Context, req *domain.SuggestRequest) (*domain.SuggestResponse, error) {
	// Validate request
	if req.Query == "" {
		return nil, &domain.BusinessError{
			Code:    "INVALID_SUGGEST_QUERY",
			Message: "Suggest query cannot be empty",
		}
	}

	// Set defaults
	if req.Limit <= 0 {
		req.Limit = 10
	}

	// Get suggestions
	suggestions, err := s.searchRepo.Suggest(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("suggest failed: %w", err)
	}

	// Build response
	response := &domain.SuggestResponse{
		Suggestions: suggestions,
		Query:       req.Query,
	}

	return response, nil
}

// Reindex rebuilds the search index by fetching data from CMS service
func (s *SearchServiceImpl) Reindex(ctx context.Context) error {
	// Fetch all media from CMS service
	mediaListResponse, err := s.cmsClient.Get(ctx, "/api/v1/media?limit=1000")
	if err != nil {
		return fmt.Errorf("failed to fetch media from CMS service: %w", err)
	}

	// Parse response
	var cmsResponse struct {
		Items []*domain.Media `json:"items"`
		Total int64           `json:"total"`
	}

	if err := json.Unmarshal(mediaListResponse, &cmsResponse); err != nil {
		return fmt.Errorf("failed to parse CMS response: %w", err)
	}

	// Reindex all media
	if err := s.searchRepo.ReindexAll(ctx, cmsResponse.Items); err != nil {
		return fmt.Errorf("failed to reindex: %w", err)
	}

	return nil
}
