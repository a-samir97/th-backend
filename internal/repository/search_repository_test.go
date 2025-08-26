package repository

import (
	"context"
	"testing"

	"thamaniyah/internal/domain"
)

// TestSearchRepositoryInterface is an interface test to ensure all implementations
// satisfy the SearchRepository interface
func TestSearchRepositoryInterface(t *testing.T) {
	// This is a compile-time check to ensure our interface is properly defined
	var _ SearchRepository = (*MockSearchRepository)(nil)
}

// MockSearchRepository can be used in tests
type MockSearchRepository struct{}

func (m *MockSearchRepository) Search(ctx context.Context, req *domain.SearchRequest) ([]*domain.SearchResult, int64, error) {
	return nil, 0, nil
}

func (m *MockSearchRepository) Suggest(ctx context.Context, req *domain.SuggestRequest) ([]*domain.Suggestion, error) {
	return nil, nil
}

func (m *MockSearchRepository) IndexMedia(ctx context.Context, media *domain.Media) error {
	return nil
}

func (m *MockSearchRepository) RemoveFromIndex(ctx context.Context, mediaID string) error {
	return nil
}

func (m *MockSearchRepository) ReindexAll(ctx context.Context, mediaList []*domain.Media) error {
	return nil
}
