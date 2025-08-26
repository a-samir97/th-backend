package service

import (
	"context"
	"errors"
	"testing"

	"thamaniyah/internal/domain"
	"thamaniyah/pkg/httpclient"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockSearchRepository is a mock implementation of SearchRepository
type MockSearchRepository struct {
	mock.Mock
}

func (m *MockSearchRepository) Search(ctx context.Context, req *domain.SearchRequest) ([]*domain.SearchResult, int64, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*domain.SearchResult), args.Get(1).(int64), args.Error(2)
}

func (m *MockSearchRepository) Suggest(ctx context.Context, req *domain.SuggestRequest) ([]*domain.Suggestion, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Suggestion), args.Error(1)
}

func (m *MockSearchRepository) IndexMedia(ctx context.Context, media *domain.Media) error {
	args := m.Called(ctx, media)
	return args.Error(0)
}

func (m *MockSearchRepository) RemoveFromIndex(ctx context.Context, mediaID string) error {
	args := m.Called(ctx, mediaID)
	return args.Error(0)
}

func (m *MockSearchRepository) ReindexAll(ctx context.Context, media []*domain.Media) error {
	args := m.Called(ctx, media)
	return args.Error(0)
}

func TestSearchService_Search(t *testing.T) {
	tests := []struct {
		name        string
		request     *domain.SearchRequest
		setupMock   func(*MockSearchRepository)
		expectError bool
		errorCode   string
	}{
		{
			name: "successful search",
			request: &domain.SearchRequest{
				Query:  "golang tutorial",
				Type:   "video",
				Limit:  10,
				Offset: 0,
			},
			setupMock: func(mockRepo *MockSearchRepository) {
				results := []*domain.SearchResult{
					{
						Media: &domain.Media{
							ID:    "media-1",
							Title: "Golang Tutorial",
							Type:  domain.TypeVideo,
						},
						Score: 2.5,
					},
				}
				mockRepo.On("Search", mock.Anything, mock.AnythingOfType("*domain.SearchRequest")).
					Return(results, int64(1), nil)
			},
			expectError: false,
		},
		{
			name: "empty query",
			request: &domain.SearchRequest{
				Query: "",
				Limit: 10,
			},
			setupMock: func(mockRepo *MockSearchRepository) {
				// No expectations as validation should fail before repository call
			},
			expectError: true,
			errorCode:   "INVALID_SEARCH_QUERY",
		},
		{
			name: "search with defaults",
			request: &domain.SearchRequest{
				Query:  "test",
				Limit:  0,  // Should default to 20
				Offset: -1, // Should default to 0
			},
			setupMock: func(mockRepo *MockSearchRepository) {
				mockRepo.On("Search", mock.Anything, mock.MatchedBy(func(req *domain.SearchRequest) bool {
					return req.Query == "test" && req.Limit == 20 && req.Offset == 0
				})).Return([]*domain.SearchResult{}, int64(0), nil)
			},
			expectError: false,
		},
		{
			name: "repository error",
			request: &domain.SearchRequest{
				Query: "test",
				Limit: 10,
			},
			setupMock: func(mockRepo *MockSearchRepository) {
				mockRepo.On("Search", mock.Anything, mock.AnythingOfType("*domain.SearchRequest")).
					Return(nil, int64(0), errors.New("elasticsearch error"))
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			mockRepo := new(MockSearchRepository)
			tt.setupMock(mockRepo)
			service := NewSearchService(mockRepo, &httpclient.Client{})
			ctx := context.Background()

			// When
			result, err := service.Search(ctx, tt.request)

			// Then
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
				if tt.errorCode != "" {
					var businessErr *domain.BusinessError
					if errors.As(err, &businessErr) {
						assert.Equal(t, tt.errorCode, businessErr.Code)
					}
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.request.Query, result.Query)
				// Check that defaults were applied
				if tt.request.Limit <= 0 {
					assert.Equal(t, 20, result.Limit)
				}
				if tt.request.Offset < 0 {
					assert.Equal(t, 0, result.Offset)
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestSearchService_Suggest(t *testing.T) {
	tests := []struct {
		name        string
		request     *domain.SuggestRequest
		setupMock   func(*MockSearchRepository)
		expectError bool
		errorCode   string
	}{
		{
			name: "successful suggestions",
			request: &domain.SuggestRequest{
				Query: "golang",
				Limit: 5,
			},
			setupMock: func(mockRepo *MockSearchRepository) {
				suggestions := []*domain.Suggestion{
					{Text: "golang tutorial", Count: 10},
					{Text: "golang basics", Count: 5},
				}
				mockRepo.On("Suggest", mock.Anything, mock.AnythingOfType("*domain.SuggestRequest")).
					Return(suggestions, nil)
			},
			expectError: false,
		},
		{
			name: "empty query",
			request: &domain.SuggestRequest{
				Query: "",
			},
			setupMock: func(mockRepo *MockSearchRepository) {
				// No expectations as validation should fail before repository call
			},
			expectError: true,
			errorCode:   "INVALID_SUGGEST_QUERY",
		},
		{
			name: "suggest with defaults",
			request: &domain.SuggestRequest{
				Query: "test",
				Limit: 0, // Should default to 10
			},
			setupMock: func(mockRepo *MockSearchRepository) {
				mockRepo.On("Suggest", mock.Anything, mock.MatchedBy(func(req *domain.SuggestRequest) bool {
					return req.Query == "test" && req.Limit == 10
				})).Return([]*domain.Suggestion{}, nil)
			},
			expectError: false,
		},
		{
			name: "repository error",
			request: &domain.SuggestRequest{
				Query: "test",
				Limit: 5,
			},
			setupMock: func(mockRepo *MockSearchRepository) {
				mockRepo.On("Suggest", mock.Anything, mock.AnythingOfType("*domain.SuggestRequest")).
					Return(nil, errors.New("elasticsearch error"))
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			mockRepo := new(MockSearchRepository)
			tt.setupMock(mockRepo)
			service := NewSearchService(mockRepo, &httpclient.Client{})
			ctx := context.Background()

			// When
			result, err := service.Suggest(ctx, tt.request)

			// Then
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
				if tt.errorCode != "" {
					var businessErr *domain.BusinessError
					if errors.As(err, &businessErr) {
						assert.Equal(t, tt.errorCode, businessErr.Code)
					}
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.request.Query, result.Query)
				// Check that defaults were applied
				if tt.request.Limit <= 0 {
					// Note: the service uses 10 as default, not reflected in response
					assert.NotNil(t, result.Suggestions)
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestSearchService_Reindex(t *testing.T) {
	// Note: This is a simplified test since mocking HTTP client requires more setup
	// In a real application, you would inject an HTTP client interface for better testability
	
	t.Run("reindex with mock repository", func(t *testing.T) {
		// Given
		mockRepo := new(MockSearchRepository)
		// Setup mock to expect ReindexAll to be called (though HTTP call will fail)
		mockRepo.On("ReindexAll", mock.Anything, mock.AnythingOfType("[]*domain.Media")).Return(nil)
		
		// Create service - note this will try to make HTTP calls
		service := NewSearchService(mockRepo, httpclient.NewClient("http://localhost:8080"))
		ctx := context.Background()

		// When - this will fail due to HTTP connection, which is expected in unit tests
		err := service.Reindex(ctx)

		// Then - we expect an error since HTTP client can't connect
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to fetch media from CMS service")
	})
}

func TestNewSearchService(t *testing.T) {
	// Given
	mockRepo := new(MockSearchRepository)
	mockClient := &httpclient.Client{}

	// When
	service := NewSearchService(mockRepo, mockClient)

	// Then
	assert.NotNil(t, service)
	assert.IsType(t, &SearchServiceImpl{}, service)
}
