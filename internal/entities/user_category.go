package entities

import (
	"time"

	"github.com/google/uuid"
)

// UserCategory represents user's custom categories
type UserCategory struct {
	UserCategoryID uuid.UUID `json:"user_category_id" db:"user_category_id"`
	UserID         uuid.UUID `json:"user_id" db:"user_id"`
	CategoryID     uuid.UUID `json:"category_id" db:"category_id"`
	Name           string    `json:"name" db:"name"`
	Description    *string   `json:"description" db:"description"`
	Icon           *string   `json:"icon" db:"icon"`
	IsCustom       bool      `json:"is_custom" db:"is_custom"`
	IsActive       bool      `json:"is_active" db:"is_active"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

type CreateUserCategoryInput struct {
	CategoryID  uuid.UUID `json:"category_id" db:"category_id"`
	UserID      uuid.UUID `json:"user_id" db:"user_id"`
	Name        string    `json:"name" db:"name"`
	Description *string   `json:"description" db:"description"`
	Icon        *string   `json:"icon" db:"icon"`
	IsCustom    bool      `json:"is_custom" db:"is_custom"`
	IsActive    bool      `json:"is_active" db:"is_active"`
}

type UpdateUserCategoryInput struct {
	CategoryID  uuid.UUID `json:"category_id" db:"category_id"`
	UserID      uuid.UUID `json:"user_id" db:"user_id"`
	Name        string    `json:"name" db:"name"`
	Description *string   `json:"description" db:"description"`
	Icon        *string   `json:"icon" db:"icon"`
	IsCustom    bool      `json:"is_custom" db:"is_custom"`
	IsActive    bool      `json:"is_active" db:"is_active"`
}
