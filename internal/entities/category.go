package entities

import (
	"time"

	"github.com/google/uuid"
)

// Category represents an expense/income category
type Category struct {
	CategoryID       uuid.UUID  `json:"category_id" db:"category_id"`
	Name             string     `json:"name" db:"name"`
	Description      *string    `json:"description" db:"description"`
	Icon             *string    `json:"icon" db:"icon"`
	ParentCategoryID *uuid.UUID `json:"parent_category_id" db:"parent_category_id"`
	IsSystemCategory bool       `json:"is_system_category" db:"is_system_category"`
	CreatedAt        time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at" db:"updated_at"`
}

type CreateCategoryInput struct {
	Name             string     `json:"name" binding:"required"`
	Description      *string    `json:"description"`
	Icon             *string    `json:"icon"`
	ParentCategoryID *uuid.UUID `json:"parent_category_id" db:"parent_category_id"`
	IsSystemCategory bool       `json:"is_system_category" db:"is_system_category"`
}

type UpdateCategoryInput struct {
	Name             string     `json:"name" binding:"required"`
	Description      *string    `json:"description"`
	Icon             *string    `json:"icon"`
	ParentCategoryID *uuid.UUID `json:"parent_category_id" db:"parent_category_id"`
	IsSystemCategory bool       `json:"is_system_category" db:"is_system_category"`
}
