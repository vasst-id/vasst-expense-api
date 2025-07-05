package entities

import (
	"time"

	"github.com/google/uuid"
)

// ContactTags represents the relationship between contacts and tags
type ContactTags struct {
	ContactID uuid.UUID `json:"contact_id" db:"contact_id"`
	TagID     uuid.UUID `json:"tag_id" db:"tag_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// CreateContactTagsInput is used for creating a new contact-tag relationship
type CreateContactTagsInput struct {
	ContactID uuid.UUID `json:"contact_id" binding:"required"`
	TagID     uuid.UUID `json:"tag_id" binding:"required"`
}

// UpdateContactTagsInput is used for updating an existing contact-tag relationship
type UpdateContactTagsInput struct {
	TagID uuid.UUID `json:"tag_id" binding:"required"`
}
