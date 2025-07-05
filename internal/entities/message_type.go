package entities

// MessageType represents a type of message in the system
type MessageType struct {
	MessageTypeID int    `json:"message_type_id" db:"message_type_id"`
	Name          string `json:"name" db:"name"`
	Description   string `json:"description" db:"description"`
	IsActive      bool   `json:"is_active" db:"is_active"`
}

// CreateMessageTypeInput is used for creating a new message type
type CreateMessageTypeInput struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active"`
}

// UpdateMessageTypeInput is used for updating an existing message type
type UpdateMessageTypeInput struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	IsActive    *bool  `json:"is_active"`
}
