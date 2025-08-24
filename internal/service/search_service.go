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
	return s.reindexWithPagination(ctx)
}

// reindexWithPagination handles large datasets by paginating through CMS data
func (s *SearchServiceImpl) reindexWithPagination(ctx context.Context) error {
	const batchSize = 100
	var offset int
	var allMedia []*domain.Media

	for {
		// Fetch batch of media from CMS service
		url := fmt.Sprintf("/api/v1/media?limit=%d&offset=%d", batchSize, offset)
		mediaListResponse, err := s.cmsClient.Get(ctx, url)
		if err != nil {
			return fmt.Errorf("failed to fetch media from CMS service at offset %d: %w", offset, err)
		}

		// Parse response
		var cmsResponse struct {
			Items []*domain.Media `json:"items"`
			Total int64           `json:"total"`
		}

		if err := json.Unmarshal(mediaListResponse, &cmsResponse); err != nil {
			return fmt.Errorf("failed to parse CMS response: %w", err)
		}

		// Add to collection
		allMedia = append(allMedia, cmsResponse.Items...)

		// Check if we've fetched all data
		if len(cmsResponse.Items) < batchSize {
			break
		}

		offset += batchSize

		// Prevent infinite loops
		if offset > int(cmsResponse.Total) {
			break
		}
	}

	// Reindex all media in batches
	return s.searchRepo.ReindexAll(ctx, allMedia)
}
