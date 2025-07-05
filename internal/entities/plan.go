package entities

import (
	"encoding/json"
	"time"
)

// Plan represents a subscription plan in the system
type Plan struct {
	PlanID       int             `json:"plan_id" db:"plan_id"`
	Name         string          `json:"name" db:"name"`
	Description  string          `json:"description" db:"description"`
	Price        string          `json:"price" db:"price"`
	Currency     string          `json:"currency" db:"currency"`
	Duration     int             `json:"duration" db:"duration"`
	DurationType string          `json:"duration_type" db:"duration_type"`
	Features     json.RawMessage `json:"features" db:"features"`
	IsActive     bool            `json:"is_active" db:"is_active"`
	CreatedAt    time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at" db:"updated_at"`
}

// CreatePlanInput is used for creating a new plan
type CreatePlanInput struct {
	Name         string          `json:"name" binding:"required"`
	Description  string          `json:"description"`
	Price        string          `json:"price" binding:"required"`
	Currency     string          `json:"currency"`
	Duration     int             `json:"duration" binding:"required"`
	DurationType string          `json:"duration_type" binding:"required"`
	Features     json.RawMessage `json:"features" binding:"required"`
	IsActive     bool            `json:"is_active"`
}

// UpdatePlanInput is used for updating an existing plan
type UpdatePlanInput struct {
	Name         string          `json:"name"`
	Description  string          `json:"description"`
	Price        string          `json:"price"`
	Currency     string          `json:"currency"`
	Duration     int             `json:"duration"`
	DurationType string          `json:"duration_type"`
	Features     json.RawMessage `json:"features"`
	IsActive     *bool           `json:"is_active"`
}
