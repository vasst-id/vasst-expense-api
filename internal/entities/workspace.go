package entities

import (
	"time"

	"github.com/google/uuid"
)

// Workspace represents a workspace in the expense tracking system
type Workspace struct {
	WorkspaceID   uuid.UUID `json:"workspace_id" db:"workspace_id"`
	Name          string    `json:"name" db:"name"`
	Description   string    `json:"description" db:"description"`
	WorkspaceType int       `json:"workspace_type" db:"workspace_type"`
	Icon          string    `json:"icon" db:"icon"`
	ColorCode     string    `json:"color_code" db:"color_code"`
	CurrencyID    int       `json:"currency_id" db:"currency_id"`
	Timezone      string    `json:"timezone" db:"timezone"`
	Settings      string    `json:"settings" db:"settings"` // JSON string
	IsActive      bool      `json:"is_active" db:"is_active"`
	CreatedBy     uuid.UUID `json:"created_by" db:"created_by"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

type CreateWorkspaceInput struct {
	Name          string `json:"name" binding:"required"`
	Description   string `json:"description"`
	WorkspaceType int    `json:"workspace_type" binding:"required"`
	Icon          string `json:"icon"`
	ColorCode     string `json:"color_code"`
	CurrencyID    int    `json:"currency_id" binding:"required"`
	Timezone      string `json:"timezone"`
}

type UpdateWorkspaceInput struct {
	Name          string `json:"name"`
	Description   string `json:"description"`
	WorkspaceType int    `json:"workspace_type"`
	Icon          string `json:"icon"`
	ColorCode     string `json:"color_code"`
	CurrencyID    int    `json:"currency_id"`
	Timezone      string `json:"timezone"`
	IsActive      bool   `json:"is_active"`
}

// Constants for workspace types
const (
	WorkspaceTypePersonal = 1
	WorkspaceTypeBusiness = 2
	WorkspaceTypeEvent    = 3
	WorkspaceTypeTravel   = 4
	WorkspaceTypeProject  = 5
	WorkspaceTypeShared   = 6
)
