package entities

import (
	"time"

	"github.com/google/uuid"
)

// Integration represents an integration in the system
type Integration struct {
	IntegrationID   uuid.UUID `json:"integration_id" db:"integration_id"`
	IntegrationName string    `json:"integration_name" db:"integration_name"`
	GlobalToken     string    `json:"global_token" db:"global_token"`
	IsActive        bool      `json:"is_active" db:"is_active"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

// CreateIntegrationInput is used for creating a new integration
type CreateIntegrationInput struct {
	IntegrationName string `json:"integration_name" binding:"required"`
	GlobalToken     string `json:"global_token" binding:"required"`
	IsActive        bool   `json:"is_active"`
}

// UpdateIntegrationInput is used for updating an existing integration
type UpdateIntegrationInput struct {
	IntegrationName string `json:"integration_name"`
	GlobalToken     string `json:"global_token"`
	IsActive        *bool  `json:"is_active"`
}
