package entities

import (
	"time"

	"github.com/google/uuid"
)

// Model represents an AI model in the system
type Model struct {
	ModelID     uuid.UUID `json:"model_id" db:"model_id"`
	Name        string    `json:"name" db:"name"`
	ModelAPIKey string    `json:"model_api_key" db:"model_api_key"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// CreateModelInput is used for creating a new model
type CreateModelInput struct {
	Name        string `json:"name" binding:"required"`
	ModelAPIKey string `json:"model_api_key" binding:"required"`
	IsActive    bool   `json:"is_active"`
}

// UpdateModelInput is used for updating an existing model
type UpdateModelInput struct {
	Name        string `json:"name"`
	ModelAPIKey string `json:"model_api_key"`
	IsActive    *bool  `json:"is_active"`
}
