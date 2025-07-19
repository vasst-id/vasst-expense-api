package entities

import (
	"time"

	"github.com/google/uuid"
)

// UserTag represents a user's custom tag for categorizing transactions
type UserTag struct {
	UserTagID uuid.UUID `json:"user_tag_id" db:"user_tag_id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	Name      string    `json:"name" db:"name"`
	IsActive  bool      `json:"is_active" db:"is_active"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// CreateUserTagRequest represents the create user tag request
type CreateUserTagRequest struct {
	Name string `json:"name" binding:"required"`
}

// UpdateUserTagRequest represents the update user tag request
type UpdateUserTagRequest struct {
	Name     string `json:"name" binding:"required"`
	IsActive bool   `json:"is_active" binding:"required"`
}

// UserTagSimple represents a simplified user tag for dropdown/selection purposes
type UserTagSimple struct {
	UserTagID uuid.UUID `json:"user_tag_id" db:"user_tag_id"`
	Name      string    `json:"name" db:"name"`
	IsActive  bool      `json:"is_active" db:"is_active"`
}

// UserTagWithUsage represents a user tag with usage statistics
type UserTagWithUsage struct {
	UserTagID  uuid.UUID  `json:"user_tag_id" db:"user_tag_id"`
	Name       string     `json:"name" db:"name"`
	IsActive   bool       `json:"is_active" db:"is_active"`
	UsageCount int        `json:"usage_count" db:"usage_count"`
	LastUsedAt *time.Time `json:"last_used_at" db:"last_used_at"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at" db:"updated_at"`
}
