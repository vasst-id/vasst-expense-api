package entities

import (
	"time"

	"github.com/google/uuid"
)

type Lead struct {
	LeadID              uuid.UUID `json:"lead_id" db:"lead_id"`
	Name                string    `json:"name" db:"name"`
	PhoneNumber         string    `json:"phone_number" db:"phone_number"`
	Email               string    `json:"email" db:"email"`
	BusinessName        string    `json:"business_name" db:"business_name"`
	BusinessAddress     string    `json:"business_address" db:"business_address"`
	BusinessPhoneNumber string    `json:"business_phone_number" db:"business_phone_number"`
	BusinessEmail       string    `json:"business_email" db:"business_email"`
	BusinessWebsite     string    `json:"business_website" db:"business_website"`
	BusinessIndustry    string    `json:"business_industry" db:"business_industry"`
	BusinessSize        string    `json:"business_size" db:"business_size"`
	CreatedAt           time.Time `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time `json:"updated_at" db:"updated_at"`
}

type CreateLeadInput struct {
	Name                string `json:"name" binding:"required"`
	PhoneNumber         string `json:"phone_number" binding:"required"`
	Email               string `json:"email"`
	BusinessName        string `json:"business_name"`
	BusinessAddress     string `json:"business_address"`
	BusinessPhoneNumber string `json:"business_phone_number"`
	BusinessEmail       string `json:"business_email"`
	BusinessWebsite     string `json:"business_website"`
	BusinessIndustry    string `json:"business_industry"`
	BusinessSize        string `json:"business_size"`
}

type UpdateLeadInput struct {
	Name                string `json:"name"`
	PhoneNumber         string `json:"phone_number"`
	Email               string `json:"email"`
	BusinessName        string `json:"business_name"`
	BusinessAddress     string `json:"business_address"`
	BusinessPhoneNumber string `json:"business_phone_number"`
	BusinessEmail       string `json:"business_email"`
	BusinessWebsite     string `json:"business_website"`
	BusinessIndustry    string `json:"business_industry"`
	BusinessSize        string `json:"business_size"`
}
