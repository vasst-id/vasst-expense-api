package entities

import (
	"time"

	"github.com/google/uuid"
)

// Conversation represents a chat thread between the system and a user
type Conversation struct {
	ConversationID uuid.UUID `json:"conversation_id" db:"conversation_id"`
	UserID         uuid.UUID `json:"user_id" db:"user_id"`
	Channel        string    `json:"channel" db:"channel"`
	IsActive       bool      `json:"is_active" db:"is_active"`
	Context        *string   `json:"context" db:"context"`
	Metadata       *string   `json:"metadata" db:"metadata"` // JSON string
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

// CreateConversationRequest represents the create conversation request
type CreateConversationRequest struct {
	UserID   uuid.UUID `json:"user_id" binding:"required"`
	Channel  string    `json:"channel" binding:"required"`
	Context  *string   `json:"context"`
	Metadata *string   `json:"metadata"`
}

// UpdateConversationRequest represents the update conversation request
type UpdateConversationRequest struct {
	Channel  string  `json:"channel"`
	Context  *string `json:"context"`
	Metadata *string `json:"metadata"`
	IsActive bool    `json:"is_active"`
}

// ConversationSimple represents a simplified conversation for listing
type ConversationSimple struct {
	ConversationID uuid.UUID `json:"conversation_id" db:"conversation_id"`
	Channel        string    `json:"channel" db:"channel"`
	IsActive       bool      `json:"is_active" db:"is_active"`
	LastMessage    *string   `json:"last_message" db:"last_message"`
	MessageCount   int       `json:"message_count" db:"message_count"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

// Constants for conversation channels
const (
	ChannelWhatsApp = "whatsapp"
	ChannelWeb      = "web"
	ChannelAPI      = "api"
)
