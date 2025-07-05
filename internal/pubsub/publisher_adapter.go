package pubsub

import (
	"context"

	"github.com/vasst-id/vasst-expense-api/internal/entities"
	"github.com/vasst-id/vasst-expense-api/internal/events"
)

// PublisherAdapter adapts the pubsub Client to implement events.Publisher
type PublisherAdapter struct {
	client Client
}

// NewPublisherAdapter creates a new publisher adapter
func NewPublisherAdapter(client Client) events.Publisher {
	return &PublisherAdapter{
		client: client,
	}
}

// Publish publishes a message to the specified topic
func (p *PublisherAdapter) Publish(ctx context.Context, topic string, message *entities.PubSubMessage) error {
	return p.client.Publish(ctx, topic, message.Data)
}

// Close closes the underlying client
func (p *PublisherAdapter) Close() error {
	return p.client.Close()
}
