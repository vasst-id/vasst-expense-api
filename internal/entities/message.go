package entities

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// MessageDirection represents the direction of a message
type MessageDirection string

const (
	MessageDirectionIncoming MessageDirection = "i"
	MessageDirectionOutgoing MessageDirection = "o"
)

// MessageStatus represents the status of a message
type MessageStatus int

const (
	MessageStatusPending   MessageStatus = 0
	MessageStatusSent      MessageStatus = 1
	MessageStatusDelivered MessageStatus = 2
	MessageStatusRead      MessageStatus = 3
	MessageStatusFailed    MessageStatus = 4
)

// SenderType represents the type of message sender
type SenderType int

const (
	SenderTypeCustomer SenderType = 1
	SenderTypeAgent    SenderType = 2
	SenderTypeAI       SenderType = 3
	SenderTypeSystem   SenderType = 4
)

// Attachment represents a media attachment
type Attachment struct {
	ID       string `json:"id"`
	Type     string `json:"type"` // pdf, image, video, audio, document
	URL      string `json:"url"`
	Filename string `json:"filename"`
	Size     int64  `json:"size"`
	MimeType string `json:"mime_type"`
}

// MessageMetadata represents metadata for messages (especially images)
type MessageMetadata struct {
	Width    int    `json:"width,omitempty"`
	Height   int    `json:"height,omitempty"`
	Format   string `json:"format,omitempty"`
	Duration int    `json:"duration,omitempty"` // for videos/audio
	Size     int64  `json:"size,omitempty"`
	// Add more metadata fields as needed
}

// Message represents a message in the system
type Message struct {
	MessageID         uuid.UUID       `json:"message_id" db:"message_id"`
	ConversationID    uuid.UUID       `json:"conversation_id" db:"conversation_id"`
	OrganizationID    uuid.UUID       `json:"organization_id" db:"organization_id"`
	SenderTypeID      int             `json:"sender_type_id" db:"sender_type_id"`
	SenderID          *uuid.UUID      `json:"sender_id" db:"sender_id"`
	Direction         string          `json:"direction" db:"direction"`
	MessageTypeID     int             `json:"message_type_id" db:"message_type_id"`
	Content           string          `json:"content" db:"content"`
	MediaURL          string          `json:"media_url" db:"media_url"`
	Attachments       json.RawMessage `json:"attachments" db:"attachments"`
	IsBroadcast       bool            `json:"is_broadcast" db:"is_broadcast"`
	IsOrderMessage    bool            `json:"is_order_message" db:"is_order_message"`
	Metadata          json.RawMessage `json:"metadata" db:"metadata"`
	ReadAt            *time.Time      `json:"read_at" db:"read_at"`
	DeliveredAt       *time.Time      `json:"delivered_at" db:"delivered_at"`
	FailedAt          *time.Time      `json:"failed_at" db:"failed_at"`
	FailureReason     *string         `json:"failure_reason" db:"failure_reason"`
	AIGenerated       bool            `json:"ai_generated" db:"ai_generated"`
	AIConfidenceScore *float64        `json:"ai_confidence_score" db:"ai_confidence_score"`
	Status            int             `json:"status" db:"status"`
	CreatedAt         time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at" db:"updated_at"`
}

// CreateMessageInput is used for creating a new message
type CreateMessageInput struct {
	ConversationID    uuid.UUID       `json:"conversation_id"`
	OrganizationID    uuid.UUID       `json:"organization_id"`
	ContactID         uuid.UUID       `json:"contact_id"`
	MediumID          int             `json:"medium_id"`
	SenderTypeID      int             `json:"sender_type_id" binding:"required"`
	SenderID          *uuid.UUID      `json:"sender_id"`
	Direction         string          `json:"direction" binding:"required"`
	MessageTypeID     int             `json:"message_type_id"`
	Content           string          `json:"content"`
	MediaURL          string          `json:"media_url"`
	Attachments       []Attachment    `json:"attachments"`
	IsBroadcast       bool            `json:"is_broadcast"`
	IsOrderMessage    bool            `json:"is_order_message"`
	Metadata          json.RawMessage `json:"metadata"`
	AIGenerated       bool            `json:"ai_generated"`
	AIConfidenceScore *float64        `json:"ai_confidence_score"`
	Status            int             `json:"status"`
}

// UpdateMessageInput is used for updating an existing message
type UpdateMessageInput struct {
	SenderTypeID      *int            `json:"sender_type_id"`
	SenderID          *uuid.UUID      `json:"sender_id"`
	Direction         string          `json:"direction"`
	MessageTypeID     int             `json:"message_type_id"`
	Content           string          `json:"content"`
	MediaURL          string          `json:"media_url"`
	Attachments       []Attachment    `json:"attachments"`
	IsBroadcast       *bool           `json:"is_broadcast"`
	IsOrderMessage    *bool           `json:"is_order_message"`
	Metadata          json.RawMessage `json:"metadata"`
	Status            *int            `json:"status"`
	AIGenerated       *bool           `json:"ai_generated"`
	AIConfidenceScore *float64        `json:"ai_confidence_score"`
	FailureReason     *string         `json:"failure_reason"`
}

// UpdateMessageStatusInput is used for updating only the status of a message
type UpdateMessageStatusInput struct {
	Status        int     `json:"status" binding:"required"`
	FailureReason *string `json:"failure_reason"`
}

// MessageCreatedEvent represents a message creation event
type MessageCreatedEvent struct {
	EventID           uuid.UUID `json:"event_id"`
	MessageID         uuid.UUID `json:"message_id"`
	ConversationID    uuid.UUID `json:"conversation_id"`
	OrganizationID    uuid.UUID `json:"organization_id"`
	ContactID         uuid.UUID `json:"contact_id"`
	Content           string    `json:"content"`
	SenderTypeID      int       `json:"sender_type_id"`
	Direction         string    `json:"direction"`
	MessageTypeID     int       `json:"message_type_id"`
	MediaURL          string    `json:"media_url"`
	WhatsAppMessageID string    `json:"whatsapp_message_id"`
	CreatedAt         time.Time `json:"created_at"`
}
