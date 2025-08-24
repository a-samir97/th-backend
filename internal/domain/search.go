package domain

import "time"

// SearchRequest represents a search request
type SearchRequest struct {
	Query  string `json:"query" form:"query" binding:"required"`
	Type   string `json:"type,omitempty" form:"type"`     // video, podcast, or empty for all
	Limit  int    `json:"limit,omitempty" form:"limit"`   // default 20
	Offset int    `json:"offset,omitempty" form:"offset"` // default 0
}

// SearchResult represents a search result item
type SearchResult struct {
	Media *Media  `json:"media"`
	Score float64 `json:"score"` // relevance score
}

// SearchResponse represents the search response
type SearchResponse struct {
	Results []*SearchResult `json:"results"`
	Total   int64           `json:"total"`
	Query   string          `json:"query"`
	Limit   int             `json:"limit"`
	Offset  int             `json:"offset"`
}

// SuggestRequest represents a suggestion request
type SuggestRequest struct {
	Query string `json:"query" form:"query" binding:"required"`
	Limit int    `json:"limit,omitempty" form:"limit"` // default 10
}

// Suggestion represents a search suggestion
type Suggestion struct {
	Text  string `json:"text"`
	Count int    `json:"count"` // how many matches
}

// SuggestResponse represents the suggestions response
type SuggestResponse struct {
	Suggestions []*Suggestion `json:"suggestions"`
	Query       string        `json:"query"`
}

// SearchIndex represents a search index entry in the database
type SearchIndex struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	MediaID     string    `json:"media_id" gorm:"index;not null"`
	Title       string    `json:"title" gorm:"not null"`
	Description string    `json:"description"`
	Content     string    `json:"content"`                      // combined searchable text
	Type        MediaType `json:"type" gorm:"type:varchar(20)"` // video, podcast
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName specifies the table name for SearchIndex
func (SearchIndex) TableName() string {
	return "search_index"
}
