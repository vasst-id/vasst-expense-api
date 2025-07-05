package entities

import (
	"time"

	"github.com/google/uuid"
)

// AIResponseEvent represents an AI response event
type AIResponseEvent struct {
	EventID           uuid.UUID `json:"event_id"`
	MessageID         uuid.UUID `json:"message_id"`
	ConversationID    uuid.UUID `json:"conversation_id"`
	OrganizationID    uuid.UUID `json:"organization_id"`
	ContactID         uuid.UUID `json:"contact_id"`
	Response          string    `json:"response"`
	Model             string    `json:"model"`
	ConfidenceScore   float64   `json:"confidence_score"`
	ProcessingTime    int64     `json:"processing_time_ms"`
	WhatsAppMessageID string    `json:"whatsapp_message_id"`
	CreatedAt         time.Time `json:"created_at"`
}
