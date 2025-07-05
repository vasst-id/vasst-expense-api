package entities

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Contact represents a contact in the system
type Contact struct {
	ContactID          uuid.UUID       `json:"contact_id" db:"contact_id"`
	OrganizationID     uuid.UUID       `json:"organization_id" db:"organization_id"`
	Name               string          `json:"name" db:"contact_name"`
	PhoneNumber        string          `json:"phone_number" db:"phone_number"`
	Email              string          `json:"email" db:"email"`
	Salutation         string          `json:"salutation" db:"salutation"`
	Notes              string          `json:"notes" db:"notes"`
	CustomSystemPrompt string          `json:"custom_system_prompt" db:"custom_system_prompt"`
	Context            json.RawMessage `json:"context" db:"context"`
	CreatedAt          time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time       `json:"updated_at" db:"updated_at"`
}

// CreateContactInput is used for creating a new contact
type CreateContactInput struct {
	OrganizationID     uuid.UUID       `json:"organization_id"`
	Name               string          `json:"name" binding:"required"`
	PhoneNumber        string          `json:"phone_number" binding:"required"`
	Email              string          `json:"email"`
	Salutation         string          `json:"salutation"`
	Notes              string          `json:"notes"`
	CustomSystemPrompt string          `json:"custom_system_prompt"`
	Context            json.RawMessage `json:"context"`
}

// UpdateContactInput is used for updating an existing contact
type UpdateContactInput struct {
	Name               string          `json:"name"`
	PhoneNumber        string          `json:"phone_number"`
	Email              string          `json:"email"`
	Salutation         string          `json:"salutation"`
	Notes              string          `json:"notes"`
	CustomSystemPrompt string          `json:"custom_system_prompt"`
	Context            json.RawMessage `json:"context"`
}
