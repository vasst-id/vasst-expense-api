package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/vasst-id/vasst-expense-api/internal/entities"
	"github.com/vasst-id/vasst-expense-api/internal/events"
	"github.com/vasst-id/vasst-expense-api/internal/events/handlers/processors"
)

// WebhookEventHandler handles webhooks from all platforms
type WebhookEventHandler struct {
	publisher events.Publisher
}

// NewWebhookEventHandler creates a new unified webhook event handler
func NewWebhookEventHandler(pub events.Publisher) *WebhookEventHandler {
	return &WebhookEventHandler{
		publisher: pub,
	}
}

// HandleWebhook processes webhook payloads from any platform
func (h *WebhookEventHandler) HandleWebhook(ctx context.Context, platform string, payload map[string]interface{}) error {
	log.Info().
		Str("platform", platform).
		Msg("Processing webhook event")

	// Get platform-specific processor
	processor := processors.GetMessageProcessor(platform)
	if processor == nil {
		return fmt.Errorf("unsupported platform: %s", platform)
	}

	// Extract messages using platform processor
	messages, err := processor.ExtractMessages(payload)
	if err != nil {
		log.Error().
			Err(err).
			Str("platform", platform).
			Msg("Failed to extract messages from webhook payload")
		return fmt.Errorf("failed to extract messages: %w", err)
	}

	if len(messages) == 0 {
		log.Warn().
			Str("platform", platform).
			Msg("No messages found in webhook payload")
		return nil
	}

	// Create webhook event for publishing
	webhookEvent := &entities.WebhookReceivedEvent{
		EventID:        uuid.New(),
		Platform:       platform,
		OrganizationID: organizationID,
		MediumID:       processor.GetMediumID(),
		Messages:       h.convertToWebhookMessages(messages),
		RawPayload:     payload,
		ReceivedAt:     time.Now(),
	}

	// Publish webhook event to Pub/Sub
	if err := h.publishWebhookEvent(ctx, webhookEvent); err != nil {
		log.Error().
			Err(err).
			Str("platform", platform).
			Str("organization_id", organizationID.String()).
			Msg("Failed to publish webhook event")
		return fmt.Errorf("failed to publish webhook event: %w", err)
	}

	log.Info().
		Str("platform", platform).
		Str("organization_id", organizationID.String()).
		Int("message_count", len(messages)).
		Msg("Successfully processed webhook event")

	return nil
}

// convertToWebhookMessages converts processor messages to webhook event messages
func (h *WebhookEventHandler) convertToWebhookMessages(messages []*processors.MessageInfo) []entities.WebhookMessage {
	webhookMessages := make([]entities.WebhookMessage, len(messages))

	for i, msg := range messages {
		whatsappMessageID := ""

		if msg.MessageOriginId != "" {
			whatsappMessageID = msg.MessageOriginId
		}

		webhookMessage := entities.WebhookMessage{
			PhoneNumber:       msg.PhoneNumber,
			Content:           msg.Content,
			MediaURL:          msg.MediaURL,
			MessageType:       msg.MessageType,
			Metadata:          msg.Metadata,
			WhatsAppMessageID: whatsappMessageID,
		}

		fmt.Printf("DEBUG: Created webhookMessage with WhatsAppMessageID = '%s'\n", webhookMessage.WhatsAppMessageID)
		webhookMessages[i] = webhookMessage
	}

	return webhookMessages
}

// publishWebhookEvent publishes the webhook event to Pub/Sub
func (h *WebhookEventHandler) publishWebhookEvent(ctx context.Context, event *entities.WebhookReceivedEvent) error {
	eventData, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal webhook event: %w", err)
	}

	message := &entities.PubSubMessage{
		Data: eventData,
		Attributes: map[string]string{
			"event_type":      "webhook_received",
			"platform":        event.Platform,
			"organization_id": event.OrganizationID.String(),
			"medium_id":       fmt.Sprintf("%d", event.MediumID),
		},
	}

	return h.publisher.Publish(ctx, entities.TopicWebhookReceived, message)
}
