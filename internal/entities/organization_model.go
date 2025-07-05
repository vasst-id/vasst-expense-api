package entities

import (
	"time"

	"github.com/google/uuid"
)

// OrganizationModel represents the relationship between organizations and models
type OrganizationModel struct {
	OrganizationModelID uuid.UUID `json:"organization_model_id" db:"organization_model_id"`
	OrganizationID      uuid.UUID `json:"organization_id" db:"organization_id"`
	ModelID             uuid.UUID `json:"model_id" db:"model_id"`
	IsActive            bool      `json:"is_active" db:"is_active"`
	CreatedAt           time.Time `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time `json:"updated_at" db:"updated_at"`
}

// CreateOrganizationModelInput is used for creating a new organization-model relationship
type CreateOrganizationModelInput struct {
	OrganizationID uuid.UUID `json:"organization_id" binding:"required"`
	ModelID        uuid.UUID `json:"model_id" binding:"required"`
	IsActive       bool      `json:"is_active"`
}

// UpdateOrganizationModelInput is used for updating an existing organization-model relationship
type UpdateOrganizationModelInput struct {
	IsActive *bool `json:"is_active"`
}
