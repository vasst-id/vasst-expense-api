package events

import (
	"context"

	"github.com/vasst-id/vasst-expense-api/internal/entities"
)

// Publisher defines the interface for publishing events to Pub/Sub
type Publisher interface {
	Publish(ctx context.Context, topic string, message *entities.PubSubMessage) error
	Close() error
}
