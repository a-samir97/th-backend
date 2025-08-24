package repository

import (
	"context"
	"fmt"
	"log"

	"thamaniyah/internal/domain"
	"thamaniyah/pkg/elasticsearch"
)

// ElasticsearchSearchRepository implements SearchRepository using Elasticsearch
type ElasticsearchSearchRepository struct {
	client *elasticsearch.Client
}

// NewElasticsearchSearchRepository creates a new Elasticsearch search repository
func NewElasticsearchSearchRepository(client *elasticsearch.Client) SearchRepository {
	return &ElasticsearchSearchRepository{
		client: client,
	}
}

// Search performs full-text search using Elasticsearch
func (r *ElasticsearchSearchRepository) Search(ctx context.Context, req *domain.SearchRequest) ([]*domain.SearchResult, int64, error) {
	// Build Elasticsearch query
	query := r.buildSearchQuery(req)

	// Execute search
	searchResp, err := r.client.Search(ctx, query)
	if err != nil {
		return nil, 0, fmt.Errorf("elasticsearch search failed: %w", err)
	}

	// Convert results
	results := make([]*domain.SearchResult, 0, len(searchResp.Hits.Hits))
	for _, hit := range searchResp.Hits.Hits {
		media := r.hitToMedia(hit.Source)
		result := &domain.SearchResult{
			Media: media,
			Score: hit.Score,
		}
		results = append(results, result)
	}

	return results, searchResp.Hits.Total.Value, nil
}

// Suggest provides search suggestions using Elasticsearch
func (r *ElasticsearchSearchRepository) Suggest(ctx context.Context, req *domain.SuggestRequest) ([]*domain.Suggestion, error) {
	// Build suggestion query using match_phrase_prefix
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match_phrase_prefix": map[string]interface{}{
				"title": req.Query,
			},
		},
		"size":    req.Limit,
		"_source": []string{"title"},
		"aggs": map[string]interface{}{
			"suggestions": map[string]interface{}{
				"terms": map[string]interface{}{
					"field": "title.keyword",
					"size":  req.Limit,
				},
			},
		},
	}

	searchResp, err := r.client.Search(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("elasticsearch suggest failed: %w", err)
	}

	// For simplicity, extract suggestions from search hits
	suggestions := make([]*domain.Suggestion, 0)
	seen := make(map[string]bool)

	for _, hit := range searchResp.Hits.Hits {
		if titleInterface, ok := hit.Source["title"]; ok {
			if title, ok := titleInterface.(string); ok && !seen[title] {
				suggestions = append(suggestions, &domain.Suggestion{
					Text:  title,
					Count: 1, // In a real implementation, you'd count occurrences
				})
				seen[title] = true
			}
		}

		if len(suggestions) >= req.Limit {
			break
		}
	}

	return suggestions, nil
}

// IndexMedia adds or updates media in Elasticsearch index
func (r *ElasticsearchSearchRepository) IndexMedia(ctx context.Context, media *domain.Media) error {
	// Create searchable document
	doc := r.mediaToDocument(media)

	// Index the document
	if err := r.client.IndexDocument(ctx, media.ID, doc); err != nil {
		return fmt.Errorf("failed to index media: %w", err)
	}

	log.Printf("Successfully indexed media: %s", media.ID)
	return nil
}

// RemoveFromIndex removes media from Elasticsearch index
func (r *ElasticsearchSearchRepository) RemoveFromIndex(ctx context.Context, mediaID string) error {
	if err := r.client.DeleteDocument(ctx, mediaID); err != nil {
		return fmt.Errorf("failed to remove media from index: %w", err)
	}

	log.Printf("Successfully removed media from index: %s", mediaID)
	return nil
}

// ReindexAll rebuilds the entire Elasticsearch index
func (r *ElasticsearchSearchRepository) ReindexAll(ctx context.Context, mediaList []*domain.Media) error {
	// Clear existing index
	if err := r.client.ClearIndex(ctx); err != nil {
		return fmt.Errorf("failed to clear index: %w", err)
	}

	if len(mediaList) == 0 {
		log.Println("No media to index")
		return nil
	}

	// Prepare bulk documents
	documents := make([]elasticsearch.BulkDocument, 0, len(mediaList))
	for _, media := range mediaList {
		doc := r.mediaToDocument(media)
		documents = append(documents, elasticsearch.BulkDocument{
			ID:     media.ID,
			Source: doc,
		})
	}

	// Bulk index
	if err := r.client.BulkIndex(ctx, documents); err != nil {
		return fmt.Errorf("failed to bulk index media: %w", err)
	}

	log.Printf("Successfully reindexed %d media items", len(mediaList))
	return nil
}

// Helper methods

// buildSearchQuery constructs Elasticsearch query from SearchRequest
func (r *ElasticsearchSearchRepository) buildSearchQuery(req *domain.SearchRequest) map[string]interface{} {
	query := map[string]interface{}{
		"size": req.Limit,
		"from": req.Offset,
	}

	// Build bool query
	boolQuery := map[string]interface{}{
		"must": []interface{}{},
	}

	// Add text search if query provided
	if req.Query != "" {
		textQuery := map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  req.Query,
				"fields": []string{"title^2", "description", "content"}, // Boost title matches
				"type":   "best_fields",
			},
		}
		boolQuery["must"] = append(boolQuery["must"].([]interface{}), textQuery)
	} else {
		// If no query, match all
		boolQuery["must"] = append(boolQuery["must"].([]interface{}), map[string]interface{}{
			"match_all": map[string]interface{}{},
		})
	}

	// Add type filter if specified
	if req.Type != "" {
		typeFilter := map[string]interface{}{
			"term": map[string]interface{}{
				"type": req.Type,
			},
		}
		if boolQuery["filter"] == nil {
			boolQuery["filter"] = []interface{}{}
		}
		boolQuery["filter"] = append(boolQuery["filter"].([]interface{}), typeFilter)
	}

	query["query"] = map[string]interface{}{
		"bool": boolQuery,
	}

	// Add sorting
	query["sort"] = []map[string]interface{}{
		{"_score": map[string]string{"order": "desc"}},
		{"created_at": map[string]string{"order": "desc"}},
	}

	return query
}

// mediaToDocument converts Media to Elasticsearch document
func (r *ElasticsearchSearchRepository) mediaToDocument(media *domain.Media) map[string]interface{} {
	// Create searchable content
	content := media.Title + " " + media.Description

	return map[string]interface{}{
		"id":          media.ID,
		"title":       media.Title,
		"description": media.Description,
		"content":     content,
		"type":        media.Type,
		"status":      media.Status,
		"file_path":   media.FilePath,
		"file_size":   media.FileSize,
		"duration":    media.Duration,
		"format":      media.Format,
		"created_at":  media.CreatedAt,
		"updated_at":  media.UpdatedAt,
	}
}

// hitToMedia converts Elasticsearch hit to Media domain model
func (r *ElasticsearchSearchRepository) hitToMedia(source map[string]interface{}) *domain.Media {
	media := &domain.Media{}

	if id, ok := source["id"].(string); ok {
		media.ID = id
	}
	if title, ok := source["title"].(string); ok {
		media.Title = title
	}
	if description, ok := source["description"].(string); ok {
		media.Description = description
	}
	if mediaType, ok := source["type"].(string); ok {
		media.Type = domain.MediaType(mediaType)
	}
	if status, ok := source["status"].(string); ok {
		media.Status = domain.MediaStatus(status)
	}
	if filePath, ok := source["file_path"].(string); ok {
		media.FilePath = filePath
	}
	if fileSize, ok := source["file_size"].(float64); ok {
		media.FileSize = int64(fileSize)
	}
	if duration, ok := source["duration"].(float64); ok {
		media.Duration = int(duration)
	}
	if format, ok := source["format"].(string); ok {
		media.Format = format
	}

	// For search results, we set status as ready since we only index ready content
	if media.Status == "" {
		media.Status = domain.StatusReady
	}

	return media
}
