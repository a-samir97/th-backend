package messagequeue

import (
	"context"
)

// MessageQueue defines the interface for message queue operations
type MessageQueue interface {
	// Publish sends a message to a topic/queue
	Publish(ctx context.Context, topic string, message []byte) error

	// Subscribe listens for messages on a topic/queue
	Subscribe(ctx context.Context, topic string, handler MessageHandler) error

	// Close closes the message queue connection
	Close() error
}

// MessageHandler defines the function signature for message handlers
type MessageHandler func(ctx context.Context, message []byte) error

// MediaIndexEvent represents a media indexing event
type MediaIndexEvent struct {
	EventType string      `json:"event_type"` // "created", "updated", "deleted"
	MediaID   string      `json:"media_id"`
	Media     interface{} `json:"media,omitempty"`
}
