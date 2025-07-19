package entities

import (
	"time"

	"github.com/google/uuid"
)

// Message represents a conversation message
type Message struct {
	MessageID            uuid.UUID  `json:"message_id" db:"message_id"`
	ConversationID       uuid.UUID  `json:"conversation_id" db:"conversation_id"`
	UserID               *uuid.UUID `json:"user_id" db:"user_id"`
	SenderType           int        `json:"sender_type" db:"sender_type"`
	Direction            string     `json:"direction" db:"direction"`
	MessageType          int        `json:"message_type" db:"message_type"`
	Content              *string    `json:"content" db:"content"`
	MediaURL             *string    `json:"media_url" db:"media_url"`
	Attachments          *string    `json:"attachments" db:"attachments"` // JSON string
	MediaMimeType        *string    `json:"media_mime_type" db:"media_mime_type"`
	Transcription        *string    `json:"transcription" db:"transcription"`
	AIProcessed          bool       `json:"ai_processed" db:"ai_processed"`
	AIModel              *string    `json:"ai_model" db:"ai_model"`
	AIConfidenceScore    *float64   `json:"ai_confidence_score" db:"ai_confidence_score"`
	RelatedTransactionID *uuid.UUID `json:"related_transaction_id" db:"related_transaction_id"`
	ScheduledTaskID      *uuid.UUID `json:"scheduled_task_id" db:"scheduled_task_id"`
	CreatedAt            time.Time  `json:"created_at" db:"created_at"`
}

// CreateMessageRequest represents the create message request
type CreateMessageRequest struct {
	ConversationID       uuid.UUID  `json:"conversation_id" binding:"required"`
	UserID               *uuid.UUID `json:"user_id"`
	SenderType           int        `json:"sender_type" binding:"required"`
	Direction            string     `json:"direction" binding:"required"`
	MessageType          int        `json:"message_type" binding:"required"`
	Content              *string    `json:"content"`
	MediaURL             *string    `json:"media_url"`
	Attachments          *string    `json:"attachments"`
	MediaMimeType        *string    `json:"media_mime_type"`
	RelatedTransactionID *uuid.UUID `json:"related_transaction_id"`
}

// UpdateMessageRequest represents the update message request
type UpdateMessageRequest struct {
	Content              *string    `json:"content"`
	MediaURL             *string    `json:"media_url"`
	Attachments          *string    `json:"attachments"`
	MediaMimeType        *string    `json:"media_mime_type"`
	Transcription        *string    `json:"transcription"`
	AIProcessed          *bool      `json:"ai_processed"`
	AIModel              *string    `json:"ai_model"`
	AIConfidenceScore    *float64   `json:"ai_confidence_score"`
	RelatedTransactionID *uuid.UUID `json:"related_transaction_id"`
}

// MessageSimple represents a simplified message for listing
type MessageSimple struct {
	MessageID        uuid.UUID `json:"message_id" db:"message_id"`
	SenderType       int       `json:"sender_type" db:"sender_type"`
	SenderTypeLabel  string    `json:"sender_type_label" db:"sender_type_label"`
	Direction        string    `json:"direction" db:"direction"`
	MessageType      int       `json:"message_type" db:"message_type"`
	MessageTypeLabel string    `json:"message_type_label" db:"message_type_label"`
	Content          *string   `json:"content" db:"content"`
	MediaURL         *string   `json:"media_url" db:"media_url"`
	AIProcessed      bool      `json:"ai_processed" db:"ai_processed"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
}

// MessageListParams represents parameters for filtering messages
type MessageListParams struct {
	ConversationID *uuid.UUID `json:"conversation_id"`
	SenderType     *int       `json:"sender_type"`
	Direction      *string    `json:"direction"`
	MessageType    *int       `json:"message_type"`
	AIProcessed    *bool      `json:"ai_processed"`
	StartDate      *time.Time `json:"start_date"`
	EndDate        *time.Time `json:"end_date"`
}

// Constants for sender types
const (
	SenderTypeUser      = 1
	SenderTypeAI        = 2
	SenderTypeSystem    = 3
	SenderTypeScheduler = 4
)

// Constants for message directions
const (
	DirectionInbound  = "i"
	DirectionOutbound = "o"
)

// Constants for message types (these would typically reference taxonomy)
const (
	MessageTypeText     = 1
	MessageTypeImage    = 2
	MessageTypeDocument = 3
	MessageTypeAudio    = 4
	MessageTypeVideo    = 5
)
