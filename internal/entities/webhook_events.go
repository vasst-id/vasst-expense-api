package entities

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type WebhookEvent struct {
	EventID        uuid.UUID       `json:"event_id" db:"event_id"`
	OrganizationID uuid.UUID       `json:"organization_id" db:"organization_id"`
	MediumID       int             `json:"medium_id" db:"medium_id"`
	Platform       string          `json:"platform" db:"platform"`
	Payload        json.RawMessage `json:"payload" db:"payload"`
	Status         int             `json:"status" db:"status"`
	ErrorMessage   string          `json:"error_message" db:"error_message"`
	RetryCount     int             `json:"retry_count" db:"retry_count"`
	ProcessedAt    time.Time       `json:"processed_at" db:"processed_at"`
	CreatedAt      time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at" db:"updated_at"`
}

// WebhookMessage represents a single message extracted from webhook payload
type WebhookMessage struct {
	PhoneNumber       string                 `json:"phone_number"`
	Content           string                 `json:"content"`
	MediaURL          string                 `json:"media_url"`
	MediaID           string                 `json:"media_id"`
	WhatsAppMessageID string                 `json:"whatsapp_message_id"`
	MessageType       int                    `json:"message_type"`
	Metadata          map[string]interface{} `json:"metadata"`
}

// WebhookReceivedEvent represents a webhook event for Pub/Sub
type WebhookReceivedEvent struct {
	EventID        uuid.UUID              `json:"event_id"`
	Platform       string                 `json:"platform"`
	OrganizationID uuid.UUID              `json:"organization_id"`
	MediumID       int                    `json:"medium_id"`
	Messages       []WebhookMessage       `json:"messages"`
	RawPayload     map[string]interface{} `json:"raw_payload"`
	ReceivedAt     time.Time              `json:"received_at"`
}

// MessageDeliveryEvent represents an event to trigger message delivery
type MessageDeliveryEvent struct {
	EventID           uuid.UUID `json:"event_id"`
	MessageID         uuid.UUID `json:"message_id"`
	ConversationID    uuid.UUID `json:"conversation_id"`
	OrganizationID    uuid.UUID `json:"organization_id"`
	ContactID         uuid.UUID `json:"contact_id"`
	Medium            string    `json:"medium"`
	WhatsAppMessageID string    `json:"whatsapp_message_id"`
	CreatedAt         time.Time `json:"created_at"`
}

// PubSubMessage represents a message structure for Pub/Sub
type PubSubMessage struct {
	Data       []byte            `json:"data"`
	Attributes map[string]string `json:"attributes"`
	MessageID  string            `json:"message_id,omitempty"`
}
