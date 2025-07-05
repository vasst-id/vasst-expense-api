package entities

import (
	"time"

	"github.com/google/uuid"
)

// Organization represents an organization in the system
type Organization struct {
	OrganizationID   uuid.UUID `json:"organization_id" db:"organization_id"`
	OrganizationCode string    `json:"organization_code" db:"organization_code"`
	Name             string    `json:"name" db:"name"`
	ContactName      string    `json:"contact_name" db:"contact_name"`
	PhoneNumber      string    `json:"phone_number" db:"phone_number"`
	Email            string    `json:"email" db:"email"`
	WhatsappNumber   string    `json:"whatsapp_number" db:"whatsapp_number"`
	OrganizationType int       `json:"organization_type" db:"organization_type"`
	CategoryID       int       `json:"category_id" db:"organization_category_id"`
	Address          string    `json:"address" db:"address"`
	City             string    `json:"city" db:"city"`
	Province         string    `json:"province" db:"province"`
	PostalCode       string    `json:"postal_code" db:"postal_code"`
	Country          string    `json:"country" db:"country"`
	Status           int       `json:"status" db:"status"`
	APIKey           string    `json:"api_key" db:"api_key"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}

// CreateOrganizationInput is used for creating a new organization
type CreateOrganizationInput struct {
	OrganizationCode string `json:"organization_code" binding:"required"`
	Name             string `json:"name" binding:"required"`
	ContactName      string `json:"contact_name" binding:"required"`
	PhoneNumber      string `json:"phone_number" binding:"required"`
	Email            string `json:"email"`
	WhatsappNumber   string `json:"whatsapp_number"`
	OrganizationType int    `json:"organization_type" binding:"required"`
	CategoryID       int    `json:"category_id" binding:"required"`
	Address          string `json:"address"`
	City             string `json:"city"`
	Province         string `json:"province"`
	PostalCode       string `json:"postal_code"`
	Country          string `json:"country"`
	Status           int    `json:"status"`
}

// UpdateOrganizationInput is used for updating an existing organization
type UpdateOrganizationInput struct {
	Name             string `json:"name"`
	ContactName      string `json:"contact_name"`
	PhoneNumber      string `json:"phone_number"`
	Email            string `json:"email"`
	WhatsappNumber   string `json:"whatsapp_number"`
	OrganizationType int    `json:"organization_type"`
	CategoryID       int    `json:"category_id"`
	Address          string `json:"address"`
	City             string `json:"city"`
	Province         string `json:"province"`
	PostalCode       string `json:"postal_code"`
	Country          string `json:"country"`
	Status           *int   `json:"status"`
}
