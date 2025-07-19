package entities

import (
	"time"

	"github.com/google/uuid"
)

// TransactionTag represents the many-to-many relationship between transactions and user tags
type TransactionTag struct {
	TransactionTagID uuid.UUID `json:"transaction_tag_id" db:"transaction_tag_id"`
	TransactionID    uuid.UUID `json:"transaction_id" db:"transaction_id"`
	UserTagID        uuid.UUID `json:"user_tag_id" db:"user_tag_id"`
	AppliedBy        uuid.UUID `json:"applied_by" db:"applied_by"`
	AppliedAt        time.Time `json:"applied_at" db:"applied_at"`
}

// CreateTransactionTagRequest represents the create transaction tag request
type CreateTransactionTagRequest struct {
	TransactionID uuid.UUID `json:"transaction_id" binding:"required"`
	UserTagID     uuid.UUID `json:"user_tag_id" binding:"required"`
}

// CreateMultipleTransactionTagsRequest represents the request to add multiple tags to a transaction
type CreateMultipleTransactionTagsRequest struct {
	TransactionID uuid.UUID   `json:"transaction_id" binding:"required"`
	UserTagIDs    []uuid.UUID `json:"user_tag_ids" binding:"required"`
}

// TransactionTagSimple represents a simplified transaction tag for API responses
type TransactionTagSimple struct {
	TransactionTagID uuid.UUID `json:"transaction_tag_id" db:"transaction_tag_id"`
	TransactionID    uuid.UUID `json:"transaction_id" db:"transaction_id"`
	UserTagID        uuid.UUID `json:"user_tag_id" db:"user_tag_id"`
	TagName          string    `json:"tag_name" db:"tag_name"`
	AppliedBy        uuid.UUID `json:"applied_by" db:"applied_by"`
	AppliedAt        time.Time `json:"applied_at" db:"applied_at"`
}

// TransactionWithTags represents a transaction with its associated tags
type TransactionWithTags struct {
	TransactionID uuid.UUID              `json:"transaction_id" db:"transaction_id"`
	Tags          []TransactionTagSimple `json:"tags"`
}

// TaggedTransactionSummary represents a summary of transactions for a specific tag
type TaggedTransactionSummary struct {
	UserTagID        uuid.UUID `json:"user_tag_id" db:"user_tag_id"`
	TagName          string    `json:"tag_name" db:"tag_name"`
	TransactionCount int       `json:"transaction_count" db:"transaction_count"`
	TotalAmount      float64   `json:"total_amount" db:"total_amount"`
	LastUsedAt       time.Time `json:"last_used_at" db:"last_used_at"`
}
