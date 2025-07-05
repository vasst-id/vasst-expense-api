package entities

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// ConversationStatus represents the status of a conversation
type ConversationStatus int

const (
	ConversationStatusOpen     ConversationStatus = 0
	ConversationStatusClosed   ConversationStatus = 1
	ConversationStatusPending  ConversationStatus = 2
	ConversationStatusResolved ConversationStatus = 3
)

// ConversationPriority represents the priority of a conversation
type ConversationPriority int

const (
	ConversationPriorityLow    ConversationPriority = 0
	ConversationPriorityMedium ConversationPriority = 1
	ConversationPriorityHigh   ConversationPriority = 2
	ConversationPriorityUrgent ConversationPriority = 3
)

// AIConfig represents AI-specific settings for a conversation
type AIConfig struct {
	SystemPrompt string            `json:"system_prompt,omitempty"`
	Model        string            `json:"model,omitempty"`
	Temperature  float64           `json:"temperature,omitempty"`
	MaxTokens    int               `json:"max_tokens,omitempty"`
	Knowledge    map[string]string `json:"knowledge,omitempty"`
	// Add more AI configuration fields as needed
}

// ConversationMetadata represents metadata for conversations
type ConversationMetadata struct {
	Tags       []string          `json:"tags,omitempty"`
	Categories []string          `json:"categories,omitempty"`
	Properties map[string]string `json:"properties,omitempty"`
	// Add more metadata fields as needed
}

// Conversation represents a conversation in the system
type Conversation struct {
	ConversationID      uuid.UUID       `json:"conversation_id" db:"conversation_id"`
	OrganizationID      uuid.UUID       `json:"organization_id" db:"organization_id"`
	UserID              uuid.UUID       `json:"user_id" db:"user_id"`
	ContactID           uuid.UUID       `json:"contact_id" db:"contact_id"`
	MediumID            int             `json:"medium_id" db:"medium_id"`
	IsActive            bool            `json:"is_active" db:"is_active"`
	IsArchived          bool            `json:"is_archived" db:"is_archived"`
	IsDeleted           bool            `json:"is_deleted" db:"is_deleted"`
	Status              int             `json:"status" db:"status"`
	Priority            int             `json:"priority" db:"priority"`
	AIEnabled           bool            `json:"ai_enabled" db:"ai_enabled"`
	AIConfig            json.RawMessage `json:"ai_config" db:"ai_config"`
	Metadata            json.RawMessage `json:"metadata" db:"metadata"`
	LastMessageAt       *time.Time      `json:"last_message_at" db:"last_message_at"`
	LastHumanMessageAt  *time.Time      `json:"last_human_message_at" db:"last_human_message_at"`
	LastAIMessageAt     *time.Time      `json:"last_ai_message_at" db:"last_ai_message_at"`
	LastMessageByID     *uuid.UUID      `json:"last_message_by_id" db:"last_message_by_id"`
	LastMessageByType   int             `json:"last_message_by_type" db:"last_message_by_type"`
	LastMessageByName   string          `json:"last_message_by_name" db:"last_message_by_name"`
	LastMessageContent  string          `json:"last_message_content" db:"last_message_content"`
	LastMessageTypeID   int             `json:"last_message_type_id" db:"last_message_type_id"`
	LastMessageMediaURL string          `json:"last_message_media_url" db:"last_message_media_url"`
	CreatedAt           time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time       `json:"updated_at" db:"updated_at"`
}

// CreateConversationInput is used for creating a new conversation
type CreateConversationInput struct {
	OrganizationID      uuid.UUID             `json:"organization_id" binding:"required"`
	UserID              uuid.UUID             `json:"user_id" binding:"required"`
	ContactID           uuid.UUID             `json:"contact_id" binding:"required"`
	MediumID            int                   `json:"medium_id" binding:"required"`
	Status              int                   `json:"status"`
	Priority            int                   `json:"priority"`
	AIEnabled           bool                  `json:"ai_enabled"`
	AIConfig            *AIConfig             `json:"ai_config"`
	Metadata            *ConversationMetadata `json:"metadata"`
	LastMessageByID     *uuid.UUID            `json:"last_message_by_id"`
	LastMessageByType   int                   `json:"last_message_by_type"`
	LastMessageByName   string                `json:"last_message_by_name"`
	LastMessageContent  string                `json:"last_message_content"`
	LastMessageTypeID   int                   `json:"last_message_type_id"`
	LastMessageMediaURL string                `json:"last_message_media_url"`
}

// UpdateConversationInput is used for updating an existing conversation
type UpdateConversationInput struct {
	UserID     uuid.UUID             `json:"user_id"`
	MediumID   int                   `json:"medium_id"`
	Status     *int                  `json:"status"`
	Priority   *int                  `json:"priority"`
	AIEnabled  *bool                 `json:"ai_enabled"`
	AIConfig   *AIConfig             `json:"ai_config"`
	Metadata   *ConversationMetadata `json:"metadata"`
	IsActive   *bool                 `json:"is_active"`
	IsArchived *bool                 `json:"is_archived"`
}

// UpdateConversationLastMessageInput is used for updating last message information
type UpdateConversationLastMessageInput struct {
	LastMessageAt       time.Time `json:"last_message_at"`
	LastMessageByID     uuid.UUID `json:"last_message_by_id"`
	LastMessageByType   int       `json:"last_message_by_type"`
	LastMessageByName   string    `json:"last_message_by_name"`
	LastMessageContent  string    `json:"last_message_content"`
	LastMessageTypeID   int       `json:"last_message_type_id"`
	LastMessageMediaURL string    `json:"last_message_media_url"`
}

// ConversationListResponse represents a conversation in list view (limited fields)
type ConversationListResponse struct {
	ConversationID     uuid.UUID  `json:"conversation_id"`
	ContactID          uuid.UUID  `json:"contact_id"`
	ContactName        string     `json:"contact_name"`
	MediumID           int        `json:"medium_id"`
	MediumName         string     `json:"medium_name"`
	Status             int        `json:"status"`
	Priority           int        `json:"priority"`
	AIEnabled          bool       `json:"ai_enabled"`
	LastMessageAt      *time.Time `json:"last_message_at"`
	LastMessageContent string     `json:"last_message_content"`
	LastMessageByName  string     `json:"last_message_by_name"`
	UnreadCount        int        `json:"unread_count"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

// ConversationDetailResponse represents a conversation with messages
type ConversationDetailResponse struct {
	Conversation    *Conversation `json:"conversation"`
	Messages        []*Message    `json:"messages"`
	TotalMessages   int           `json:"total_messages"`
	HasMoreMessages bool          `json:"has_more_messages"`
}
