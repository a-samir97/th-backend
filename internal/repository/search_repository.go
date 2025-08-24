package repository

import (
	"context"
	"fmt"

	"thamaniyah/internal/domain"
	"thamaniyah/pkg/database"
)

// SearchRepository defines search operations
type SearchRepository interface {
	// Search performs full-text search on indexed media
	Search(ctx context.Context, req *domain.SearchRequest) ([]*domain.SearchResult, int64, error)

	// Suggest provides search suggestions based on query
	Suggest(ctx context.Context, req *domain.SuggestRequest) ([]*domain.Suggestion, error)

	// IndexMedia adds or updates media in search index
	IndexMedia(ctx context.Context, media *domain.Media) error

	// RemoveFromIndex removes media from search index
	RemoveFromIndex(ctx context.Context, mediaID string) error

	// ReindexAll rebuilds the entire search index
	ReindexAll(ctx context.Context, mediaList []*domain.Media) error
}

// PostgresSearchRepository implements SearchRepository using PostgreSQL
type PostgresSearchRepository struct {
	conn *database.Connection
}

// NewPostgresSearchRepository creates a new PostgreSQL search repository
func NewPostgresSearchRepository(conn *database.Connection) SearchRepository {
	return &PostgresSearchRepository{
		conn: conn,
	}
}

// Search performs full-text search using PostgreSQL's text search capabilities
func (r *PostgresSearchRepository) Search(ctx context.Context, req *domain.SearchRequest) ([]*domain.SearchResult, int64, error) {
	var results []*domain.SearchResult
	var total int64

	// Build the search query
	query := r.conn.DB.WithContext(ctx).Model(&domain.SearchIndex{})

	// Full-text search on content field
	if req.Query != "" {
		query = query.Where("to_tsvector('english', content) @@ plainto_tsquery('english', ?)", req.Query).
			Select("*, ts_rank(to_tsvector('english', content), plainto_tsquery('english', ?)) as rank", req.Query).
			Order("rank DESC")
	}

	// Filter by type
	if req.Type != "" {
		query = query.Where("type = ?", req.Type)
	}

	// Count total results
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count search results: %w", err)
	}

	// Apply pagination
	limit := req.Limit
	if limit <= 0 {
		limit = 20
	}
	offset := req.Offset
	if offset < 0 {
		offset = 0
	}

	var searchIndexes []domain.SearchIndex
	if err := query.Limit(limit).Offset(offset).Find(&searchIndexes).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to search: %w", err)
	}

	// Convert to search results
	for _, index := range searchIndexes {
		media := r.searchIndexToMedia(&index)
		result := &domain.SearchResult{
			Media: media,
			Score: 1.0, // PostgreSQL rank would be available in a raw query
		}
		results = append(results, result)
	}

	return results, total, nil
}

// Suggest provides search suggestions
func (r *PostgresSearchRepository) Suggest(ctx context.Context, req *domain.SuggestRequest) ([]*domain.Suggestion, error) {
	var suggestions []*domain.Suggestion

	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}

	// Get suggestions from titles
	query := `
		SELECT title as suggestion, COUNT(*) as count FROM search_index 
		WHERE title ILIKE ? 
		GROUP BY title 
		ORDER BY count DESC 
		LIMIT ?`

	likePattern := "%" + req.Query + "%"

	rows, err := r.conn.DB.WithContext(ctx).Raw(query, likePattern, limit).Rows()
	if err != nil {
		return nil, fmt.Errorf("failed to get suggestions: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var suggestion domain.Suggestion
		if err := rows.Scan(&suggestion.Text, &suggestion.Count); err != nil {
			continue
		}
		suggestions = append(suggestions, &suggestion)
	}

	return suggestions, nil
}

// IndexMedia adds or updates media in search index
func (r *PostgresSearchRepository) IndexMedia(ctx context.Context, media *domain.Media) error {
	// Create searchable content by combining title and description
	content := media.Title + " " + media.Description

	searchIndex := &domain.SearchIndex{
		ID:          media.ID,
		MediaID:     media.ID,
		Title:       media.Title,
		Description: media.Description,
		Content:     content,
		Type:        media.Type,
	}

	// Use ON CONFLICT to handle updates
	if err := r.conn.DB.WithContext(ctx).Save(searchIndex).Error; err != nil {
		return fmt.Errorf("failed to index media: %w", err)
	}

	return nil
}

// RemoveFromIndex removes media from search index
func (r *PostgresSearchRepository) RemoveFromIndex(ctx context.Context, mediaID string) error {
	if err := r.conn.DB.WithContext(ctx).Delete(&domain.SearchIndex{}, "media_id = ?", mediaID).Error; err != nil {
		return fmt.Errorf("failed to remove from index: %w", err)
	}
	return nil
}

// ReindexAll rebuilds the entire search index
func (r *PostgresSearchRepository) ReindexAll(ctx context.Context, mediaList []*domain.Media) error {
	// Start transaction
	tx := r.conn.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Clear existing index
	if err := tx.Exec("TRUNCATE search_index").Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to clear search index: %w", err)
	}

	// Index all media
	for _, media := range mediaList {
		content := media.Title + " " + media.Description

		searchIndex := &domain.SearchIndex{
			ID:          media.ID,
			MediaID:     media.ID,
			Title:       media.Title,
			Description: media.Description,
			Content:     content,
			Type:        media.Type,
		}

		if err := tx.Create(searchIndex).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to index media %s: %w", media.ID, err)
		}
	}

	return tx.Commit().Error
}

// searchIndexToMedia converts SearchIndex back to Media
func (r *PostgresSearchRepository) searchIndexToMedia(index *domain.SearchIndex) *domain.Media {
	return &domain.Media{
		ID:          index.MediaID,
		Title:       index.Title,
		Description: index.Description,
		Type:        index.Type,
		Status:      domain.StatusReady, // Search results are ready
	}
}
