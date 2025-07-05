package entities

import (
	"time"

	"github.com/google/uuid"
)

// OrganizationIntegration represents the relationship between organizations and integrations
type OrganizationIntegration struct {
	OrganizationIntegrationID uuid.UUID  `json:"organization_integration_id" db:"organization_integration_id"`
	OrganizationID            uuid.UUID  `json:"organization_id" db:"organization_id"`
	IntegrationID             uuid.UUID  `json:"integration_id" db:"integration_id"`
	Token                     string     `json:"token" db:"token"`
	LastUsedAt                *time.Time `json:"last_used_at" db:"last_used_at"`
	IsActive                  bool       `json:"is_active" db:"is_active"`
	IsAiEnabled               bool       `json:"is_ai_enabled" db:"is_ai_enabled"`
	CreatedAt                 time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt                 time.Time  `json:"updated_at" db:"updated_at"`
}

// CreateOrganizationIntegrationInput is used for creating a new organization-integration relationship
type CreateOrganizationIntegrationInput struct {
	OrganizationID uuid.UUID `json:"organization_id" binding:"required"`
	IntegrationID  uuid.UUID `json:"integration_id" binding:"required"`
	Token          string    `json:"token" binding:"required"`
	IsActive       bool      `json:"is_active"`
	IsAiEnabled    bool      `json:"is_ai_enabled"`
}

// UpdateOrganizationIntegrationInput is used for updating an existing organization-integration relationship
type UpdateOrganizationIntegrationInput struct {
	Token       string     `json:"token"`
	IsAiEnabled *bool      `json:"is_ai_enabled"`
	LastUsedAt  *time.Time `json:"last_used_at"`
	IsActive    *bool      `json:"is_active"`
}
