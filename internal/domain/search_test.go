package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSearchRequest_Structure(t *testing.T) {
	// Test SearchRequest structure
	req := SearchRequest{
		Query:  "test query",
		Type:   "video",
		Limit:  10,
		Offset: 0,
	}

	assert.Equal(t, "test query", req.Query)
	assert.Equal(t, "video", req.Type)
	assert.Equal(t, 10, req.Limit)
	assert.Equal(t, 0, req.Offset)
}

func TestSearchResult_Structure(t *testing.T) {
	// Test SearchResult structure
	media := &Media{
		ID:    "media-123",
		Title: "Test Video",
	}

	result := SearchResult{
		Media: media,
		Score: 2.5,
	}

	assert.Equal(t, media, result.Media)
	assert.Equal(t, 2.5, result.Score)
}

func TestSearchResponse_Structure(t *testing.T) {
	// Test SearchResponse structure
	media := &Media{ID: "media-123", Title: "Test Video"}
	results := []*SearchResult{
		{Media: media, Score: 2.5},
	}

	response := SearchResponse{
		Results: results,
		Total:   1,
		Query:   "test",
		Limit:   10,
		Offset:  0,
	}

	assert.Equal(t, results, response.Results)
	assert.Equal(t, int64(1), response.Total)
	assert.Equal(t, "test", response.Query)
	assert.Equal(t, 10, response.Limit)
	assert.Equal(t, 0, response.Offset)
	assert.Len(t, response.Results, 1)
}

func TestSuggestRequest_Structure(t *testing.T) {
	// Test SuggestRequest structure
	req := SuggestRequest{
		Query: "test",
		Limit: 5,
	}

	assert.Equal(t, "test", req.Query)
	assert.Equal(t, 5, req.Limit)
}

func TestSuggestion_Structure(t *testing.T) {
	// Test Suggestion structure
	suggestion := Suggestion{
		Text:  "golang tutorial",
		Count: 5,
	}

	assert.Equal(t, "golang tutorial", suggestion.Text)
	assert.Equal(t, 5, suggestion.Count)
}

func TestSuggestResponse_Structure(t *testing.T) {
	// Test SuggestResponse structure
	suggestions := []*Suggestion{
		{Text: "golang tutorial", Count: 5},
		{Text: "golang basics", Count: 3},
	}

	response := SuggestResponse{
		Suggestions: suggestions,
		Query:       "golang",
	}

	assert.Equal(t, suggestions, response.Suggestions)
	assert.Equal(t, "golang", response.Query)
	assert.Len(t, response.Suggestions, 2)
}

func TestSearchIndex_Structure(t *testing.T) {
	// Test SearchIndex structure
	searchIndex := SearchIndex{
		ID:          "search-123",
		MediaID:     "media-123",
		Title:       "Test Video",
		Description: "A test video",
		Content:     "Test Video A test video",
		Type:        TypeVideo,
	}

	assert.Equal(t, "search-123", searchIndex.ID)
	assert.Equal(t, "media-123", searchIndex.MediaID)
	assert.Equal(t, "Test Video", searchIndex.Title)
	assert.Equal(t, "A test video", searchIndex.Description)
	assert.Equal(t, "Test Video A test video", searchIndex.Content)
	assert.Equal(t, TypeVideo, searchIndex.Type)
}

func TestSearchIndex_TableName(t *testing.T) {
	// Test SearchIndex table name
	searchIndex := SearchIndex{}
	tableName := searchIndex.TableName()

	assert.Equal(t, "search_index", tableName)
}

func TestSearchRequest_DefaultValues(t *testing.T) {
	// Test SearchRequest with default/empty values
	req := SearchRequest{}

	assert.Empty(t, req.Query)
	assert.Empty(t, req.Type)
	assert.Zero(t, req.Limit)
	assert.Zero(t, req.Offset)
}

func TestSuggestRequest_DefaultValues(t *testing.T) {
	// Test SuggestRequest with default values
	req := SuggestRequest{}

	assert.Empty(t, req.Query)
	assert.Zero(t, req.Limit)
}

func TestSearchResponse_EmptyResults(t *testing.T) {
	// Test SearchResponse with empty results
	response := SearchResponse{
		Results: []*SearchResult{},
		Total:   0,
		Query:   "no results",
		Limit:   10,
		Offset:  0,
	}

	assert.Empty(t, response.Results)
	assert.Zero(t, response.Total)
	assert.Equal(t, "no results", response.Query)
}

func TestSuggestResponse_EmptySuggestions(t *testing.T) {
	// Test SuggestResponse with empty suggestions
	response := SuggestResponse{
		Suggestions: []*Suggestion{},
		Query:       "no suggestions",
	}

	assert.Empty(t, response.Suggestions)
	assert.Equal(t, "no suggestions", response.Query)
}
