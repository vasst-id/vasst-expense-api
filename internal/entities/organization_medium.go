package entities

import (
	"time"

	"github.com/google/uuid"
)

// OrganizationMedium represents a medium configuration for an organization
type OrganizationMedium struct {
	OrganizationMediumID uuid.UUID `json:"organization_medium_id" db:"organization_medium_id"`
	OrganizationID       uuid.UUID `json:"organization_id" db:"organization_id"`
	MediumID             int       `json:"medium_id" db:"medium_id"`
	MediumSystemPrompt   string    `json:"medium_system_prompt" db:"medium_system_prompt"`
	IsActive             bool      `json:"is_active" db:"is_active"`
	CreatedAt            time.Time `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time `json:"updated_at" db:"updated_at"`
}

// CreateOrganizationMediumInput is used for creating a new organization medium
type CreateOrganizationMediumInput struct {
	OrganizationID     uuid.UUID `json:"organization_id" binding:"required"`
	MediumID           int       `json:"medium_id" binding:"required"`
	MediumSystemPrompt string    `json:"medium_system_prompt"`
	IsActive           bool      `json:"is_active"`
}

// UpdateOrganizationMediumInput is used for updating an existing organization medium
type UpdateOrganizationMediumInput struct {
	MediumSystemPrompt string `json:"medium_system_prompt"`
	IsActive           *bool  `json:"is_active"`
}
