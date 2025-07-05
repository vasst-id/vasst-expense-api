package pubsub

import (
	"context"
	"fmt"

	"cloud.google.com/go/pubsub"
	"google.golang.org/api/option"
)

type Client interface {
	Publish(ctx context.Context, topic string, data []byte) error
	Subscribe(ctx context.Context, subscription string, handler func([]byte) error) error
	CreateTopicIfNotExists(ctx context.Context, topicID string) error
	CreateSubscriptionIfNotExists(ctx context.Context, subscriptionID, topicID string) error
	Close() error
}

type googlePubSubClient struct {
	client    *pubsub.Client
	projectID string
}

func NewGooglePubSubClient(projectID string) (Client, error) {
	return NewGooglePubSubClientWithCredentials(projectID, "")
}

func NewGooglePubSubClientWithCredentials(projectID, credentialsFile string) (Client, error) {
	ctx := context.Background()
	
	var client *pubsub.Client
	var err error
	
	if credentialsFile != "" {
		// Use specific credentials file
		client, err = pubsub.NewClient(ctx, projectID, option.WithCredentialsFile(credentialsFile))
	} else {
		// Use default credentials (environment variable or service account)
		client, err = pubsub.NewClient(ctx, projectID)
	}
	
	if err != nil {
		return nil, fmt.Errorf("failed to create pubsub client: %w", err)
	}

	return &googlePubSubClient{
		client:    client,
		projectID: projectID,
	}, nil
}

func (g *googlePubSubClient) Publish(ctx context.Context, topicID string, data []byte) error {
	topic := g.client.Topic(topicID)
	result := topic.Publish(ctx, &pubsub.Message{Data: data})
	_, err := result.Get(ctx)
	if err != nil {
		return fmt.Errorf("failed to publish message to topic %s: %w", topicID, err)
	}
	return nil
}

func (g *googlePubSubClient) Subscribe(ctx context.Context, subscriptionID string, handler func([]byte) error) error {
	sub := g.client.Subscription(subscriptionID)
	return sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		if err := handler(msg.Data); err != nil {
			msg.Nack()
			return
		}
		msg.Ack()
	})
}

func (g *googlePubSubClient) CreateTopicIfNotExists(ctx context.Context, topicID string) error {
	topic := g.client.Topic(topicID)
	exists, err := topic.Exists(ctx)
	if err != nil {
		return fmt.Errorf("failed to check if topic %s exists: %w", topicID, err)
	}

	if !exists {
		_, err = g.client.CreateTopic(ctx, topicID)
		if err != nil {
			return fmt.Errorf("failed to create topic %s: %w", topicID, err)
		}
	}
	return nil
}

func (g *googlePubSubClient) CreateSubscriptionIfNotExists(ctx context.Context, subscriptionID, topicID string) error {
	sub := g.client.Subscription(subscriptionID)
	exists, err := sub.Exists(ctx)
	if err != nil {
		return fmt.Errorf("failed to check if subscription %s exists: %w", subscriptionID, err)
	}

	if !exists {
		topic := g.client.Topic(topicID)
		_, err = g.client.CreateSubscription(ctx, subscriptionID, pubsub.SubscriptionConfig{
			Topic: topic,
		})
		if err != nil {
			return fmt.Errorf("failed to create subscription %s: %w", subscriptionID, err)
		}
	}
	return nil
}

func (g *googlePubSubClient) Close() error {
	return g.client.Close()
}
