package entities

import (
	"time"
)

// Taxonomy represents a taxonomy entry for categorizing various system values
type Taxonomy struct {
	TaxonomyID int       `json:"taxonomy_id" db:"taxonomy_id"`
	Label      string    `json:"label" db:"label"`
	Value      string    `json:"value" db:"value"`
	Type       string    `json:"type" db:"type"`
	TypeLabel  string    `json:"type_label" db:"type_label"`
	Status     int       `json:"status" db:"status"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

// CreateTaxonomyRequest represents the create taxonomy request
type CreateTaxonomyRequest struct {
	Label     string `json:"label" binding:"required"`
	Value     string `json:"value" binding:"required"`
	Type      string `json:"type" binding:"required"`
	TypeLabel string `json:"type_label" binding:"required"`
	Status    int    `json:"status"`
}

// UpdateTaxonomyRequest represents the update taxonomy request
type UpdateTaxonomyRequest struct {
	Label     string `json:"label" binding:"required"`
	Value     string `json:"value" binding:"required"`
	Type      string `json:"type" binding:"required"`
	TypeLabel string `json:"type_label" binding:"required"`
	Status    int    `json:"status" binding:"required"`
}

// TaxonomySimple represents a simplified taxonomy for dropdown/selection purposes
type TaxonomySimple struct {
	TaxonomyID int    `json:"taxonomy_id" db:"taxonomy_id"`
	Label      string `json:"label" db:"label"`
	Value      string `json:"value" db:"value"`
	Type       string `json:"type" db:"type"`
	TypeLabel  string `json:"type_label" db:"type_label"`
}

// Constants for taxonomy status
const (
	TaxonomyStatusActive   = 1
	TaxonomyStatusInactive = 0
)

// Common taxonomy types
const (
	TaxonomyTypeAccountType      = "account_type"
	TaxonomyTypeTransactionType  = "transaction_type"
	TaxonomyTypePaymentMethod    = "payment_method"
	TaxonomyTypeWorkspaceType    = "workspace_type"
	TaxonomyTypePeriodType       = "period_type"
	TaxonomyTypeRecurrenceType   = "recurrence_interval"
	TaxonomyTypeCreditStatus     = "credit_status"
	TaxonomyTypeDocumentType     = "document_type"
	TaxonomyTypeProcessingStatus = "processing_status"
)
