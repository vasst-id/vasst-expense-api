package entities

import (
	"time"
)

// OrganizationCategory represents an organization category in the system
type OrganizationCategory struct {
	CategoryID  int       `json:"category_id" db:"organization_category_id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	ImageURL    string    `json:"image_url" db:"image_url"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// CreateOrganizationCategoryInput is used for creating a new organization category
type CreateOrganizationCategoryInput struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
	IsActive    bool   `json:"is_active"`
}

// UpdateOrganizationCategoryInput is used for updating an existing organization category
type UpdateOrganizationCategoryInput struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
	IsActive    *bool  `json:"is_active"`
}
