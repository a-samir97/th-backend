package service

import (
	"context"
	"log"

	"thamaniyah/internal/domain"
	"thamaniyah/internal/repository"
)

// MediaEventHandler handles media events for search indexing
type MediaEventHandler struct {
	searchRepo repository.SearchRepository
}

// NewMediaEventHandler creates a new media event handler
func NewMediaEventHandler(searchRepo repository.SearchRepository) *MediaEventHandler {
	return &MediaEventHandler{
		searchRepo: searchRepo,
	}
}

// HandleMediaCreated handles media creation events
func (h *MediaEventHandler) HandleMediaCreated(ctx context.Context, media *domain.Media) error {
	log.Printf("Indexing newly created media: %s", media.ID)
	return h.searchRepo.IndexMedia(ctx, media)
}

// HandleMediaUpdated handles media update events
func (h *MediaEventHandler) HandleMediaUpdated(ctx context.Context, media *domain.Media) error {
	log.Printf("Reindexing updated media: %s", media.ID)
	return h.searchRepo.IndexMedia(ctx, media)
}

// HandleMediaDeleted handles media deletion events
func (h *MediaEventHandler) HandleMediaDeleted(ctx context.Context, mediaID string) error {
	log.Printf("Removing deleted media from index: %s", mediaID)
	return h.searchRepo.RemoveFromIndex(ctx, mediaID)
}
