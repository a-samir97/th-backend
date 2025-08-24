package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"

	"thamaniyah/internal/config"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

// Client wraps the Elasticsearch client with additional functionality
type Client struct {
	es    *elasticsearch.Client
	index string
}

// NewClient creates a new Elasticsearch client
func NewClient(cfg *config.Config) (*Client, error) {
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{cfg.Elasticsearch.URL},
		// Add authentication if needed
		// Username: "elastic",
		// Password: "password",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create elasticsearch client: %w", err)
	}

	client := &Client{
		es:    es,
		index: cfg.Elasticsearch.Index,
	}

	// Check connection
	if err := client.ping(context.Background()); err != nil {
		return nil, fmt.Errorf("elasticsearch connection failed: %w", err)
	}

	// Create index if it doesn't exist
	if err := client.createIndexIfNotExists(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to create index: %w", err)
	}

	return client, nil
}

// ping checks if Elasticsearch is reachable
func (c *Client) ping(ctx context.Context) error {
	res, err := c.es.Info()
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("elasticsearch ping failed: %s", res.Status())
	}

	return nil
}

// createIndexIfNotExists creates the media index if it doesn't exist
func (c *Client) createIndexIfNotExists(ctx context.Context) error {
	// Check if index exists
	res, err := c.es.Indices.Exists([]string{c.index})
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// If index exists (200), return
	if res.StatusCode == 200 {
		log.Printf("Elasticsearch index '%s' already exists", c.index)
		return nil
	}

	// Create index with mapping
	mapping := `{
		"mappings": {
			"properties": {
				"id": {
					"type": "keyword"
				},
				"title": {
					"type": "text",
					"analyzer": "standard",
					"fields": {
						"keyword": {
							"type": "keyword"
						}
					}
				},
				"description": {
					"type": "text",
					"analyzer": "standard"
				},
				"content": {
					"type": "text",
					"analyzer": "standard"
				},
				"type": {
					"type": "keyword"
				},
				"status": {
					"type": "keyword"
				},
				"file_path": {
					"type": "keyword"
				},
				"file_size": {
					"type": "long"
				},
				"duration": {
					"type": "integer"
				},
				"format": {
					"type": "keyword"
				},
				"created_at": {
					"type": "date"
				},
				"updated_at": {
					"type": "date"
				}
			}
		},
		"settings": {
			"number_of_shards": 1,
			"number_of_replicas": 0,
			"analysis": {
				"analyzer": {
					"standard": {
						"type": "standard"
					}
				}
			}
		}
	}`

	res, err = c.es.Indices.Create(
		c.index,
		c.es.Indices.Create.WithBody(strings.NewReader(mapping)),
		c.es.Indices.Create.WithContext(ctx),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("failed to create index: %s", res.Status())
	}

	log.Printf("Elasticsearch index '%s' created successfully", c.index)
	return nil
}

// IndexDocument indexes a document
func (c *Client) IndexDocument(ctx context.Context, docID string, doc interface{}) error {
	docBytes, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("failed to marshal document: %w", err)
	}

	req := esapi.IndexRequest{
		Index:      c.index,
		DocumentID: docID,
		Body:       bytes.NewReader(docBytes),
		Refresh:    "true",
	}

	res, err := req.Do(ctx, c.es)
	if err != nil {
		return fmt.Errorf("failed to index document: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("index request failed: %s", res.Status())
	}

	return nil
}

// DeleteDocument deletes a document
func (c *Client) DeleteDocument(ctx context.Context, docID string) error {
	req := esapi.DeleteRequest{
		Index:      c.index,
		DocumentID: docID,
		Refresh:    "true",
	}

	res, err := req.Do(ctx, c.es)
	if err != nil {
		return fmt.Errorf("failed to delete document: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() && res.StatusCode != 404 {
		return fmt.Errorf("delete request failed: %s", res.Status())
	}

	return nil
}

// Search performs a search query
func (c *Client) Search(ctx context.Context, query map[string]interface{}) (*SearchResponse, error) {
	queryBytes, err := json.Marshal(query)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal query: %w", err)
	}

	req := esapi.SearchRequest{
		Index: []string{c.index},
		Body:  bytes.NewReader(queryBytes),
	}

	res, err := req.Do(ctx, c.es)
	if err != nil {
		return nil, fmt.Errorf("search request failed: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("search failed: %s", res.Status())
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var searchResp SearchResponse
	if err := json.Unmarshal(body, &searchResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal search response: %w", err)
	}

	return &searchResp, nil
}

// BulkIndex indexes multiple documents in batch
func (c *Client) BulkIndex(ctx context.Context, documents []BulkDocument) error {
	if len(documents) == 0 {
		return nil
	}

	var body strings.Builder

	for _, doc := range documents {
		// Index action
		indexAction := map[string]interface{}{
			"index": map[string]interface{}{
				"_index": c.index,
				"_id":    doc.ID,
			},
		}

		actionBytes, err := json.Marshal(indexAction)
		if err != nil {
			return fmt.Errorf("failed to marshal index action: %w", err)
		}

		body.Write(actionBytes)
		body.WriteString("\n")

		// Document
		docBytes, err := json.Marshal(doc.Source)
		if err != nil {
			return fmt.Errorf("failed to marshal document: %w", err)
		}

		body.Write(docBytes)
		body.WriteString("\n")
	}

	req := esapi.BulkRequest{
		Index:   c.index,
		Body:    strings.NewReader(body.String()),
		Refresh: "true",
	}

	res, err := req.Do(ctx, c.es)
	if err != nil {
		return fmt.Errorf("bulk request failed: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("bulk index failed: %s", res.Status())
	}

	return nil
}

// ClearIndex removes all documents from the index
func (c *Client) ClearIndex(ctx context.Context) error {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match_all": map[string]interface{}{},
		},
	}

	queryBytes, err := json.Marshal(query)
	if err != nil {
		return fmt.Errorf("failed to marshal delete query: %w", err)
	}

	req := esapi.DeleteByQueryRequest{
		Index: []string{c.index},
		Body:  bytes.NewReader(queryBytes),
	}

	res, err := req.Do(ctx, c.es)
	if err != nil {
		return fmt.Errorf("delete by query failed: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("clear index failed: %s", res.Status())
	}

	return nil
}

// Close closes the client connection
func (c *Client) Close() error {
	// The go-elasticsearch client doesn't require explicit closing
	return nil
}

// Response structures

// SearchResponse represents Elasticsearch search response
type SearchResponse struct {
	Hits struct {
		Total struct {
			Value int64 `json:"value"`
		} `json:"total"`
		Hits []struct {
			ID     string                 `json:"_id"`
			Score  float64                `json:"_score"`
			Source map[string]interface{} `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

// BulkDocument represents a document for bulk indexing
type BulkDocument struct {
	ID     string
	Source interface{}
}
