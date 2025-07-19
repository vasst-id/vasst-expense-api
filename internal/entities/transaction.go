package entities

import (
	"time"

	"github.com/google/uuid"
)

// Transaction represents a financial transaction
type Transaction struct {
	TransactionID       uuid.UUID  `json:"transaction_id" db:"transaction_id"`
	WorkspaceID         *uuid.UUID `json:"workspace_id" db:"workspace_id"`
	AccountID           *uuid.UUID `json:"account_id" db:"account_id"`
	CategoryID          *uuid.UUID `json:"category_id" db:"category_id"`
	Description         string     `json:"description" db:"description"`
	Amount              float64    `json:"amount" db:"amount"`
	TransactionType     int        `json:"transaction_type" db:"transaction_type"`
	PaymentMethod       int        `json:"payment_method" db:"payment_method"`
	TransactionDate     time.Time  `json:"transaction_date" db:"transaction_date"`
	MerchantName        *string    `json:"merchant_name" db:"merchant_name"`
	Location            *string    `json:"location" db:"location"`
	Notes               *string    `json:"notes" db:"notes"`
	ReceiptURL          *string    `json:"receipt_url" db:"receipt_url"`
	IsRecurring         bool       `json:"is_recurring" db:"is_recurring"`
	RecurrenceInterval  int        `json:"recurrence_interval" db:"recurrence_interval"`
	RecurrenceEndDate   *time.Time `json:"recurrence_end_date" db:"recurrence_end_date"`
	ParentTransactionID *uuid.UUID `json:"parent_transaction_id" db:"parent_transaction_id"`
	AIConfidenceScore   *float64   `json:"ai_confidence_score" db:"ai_confidence_score"`
	AICategorized       bool       `json:"ai_categorized" db:"ai_categorized"`
	CreditStatus        *int       `json:"credit_status" db:"credit_status"`
	CreatedBy           *uuid.UUID `json:"created_by" db:"created_by"`
	CreatedAt           time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at" db:"updated_at"`
}

type TransactionListParams struct {
	WorkspaceID   *uuid.UUID `json:"workspace_id"`
	AccountID     *uuid.UUID `json:"account_id"`
	CategoryID    *uuid.UUID `json:"category_id"`
	StartDate     *time.Time `json:"start_date"`
	EndDate       *time.Time `json:"end_date"`
	PaymentMethod *int       `json:"payment_method"`
	Description   *string    `json:"description"`
	MerchantName  *string    `json:"merchant_name"`
	Amount        *float64   `json:"amount"`
	IsRecurring   *bool      `json:"is_recurring"`
	CreditStatus  *int       `json:"credit_status"`
}

// CreateTransactionRequest represents the create transaction request
type CreateTransactionRequest struct {
	WorkspaceID        uuid.UUID  `json:"workspace_id"`
	AccountID          uuid.UUID  `json:"account_id"`
	CategoryID         *uuid.UUID `json:"category_id"`
	Description        string     `json:"description" binding:"required"`
	Amount             float64    `json:"amount" binding:"required"`
	TransactionType    int        `json:"transaction_type" binding:"required"`
	PaymentMethod      int        `json:"payment_method"`
	TransactionDate    time.Time  `json:"transaction_date" binding:"required"`
	MerchantName       *string    `json:"merchant_name"`
	Location           *string    `json:"location"`
	Notes              *string    `json:"notes"`
	IsRecurring        *bool      `json:"is_recurring"`
	RecurrenceInterval int        `json:"recurrence_interval"`
	RecurrenceEndDate  *time.Time `json:"recurrence_end_date"`
}

// UpdateTransactionRequest represents the update transaction request
type UpdateTransactionRequest struct {
	AccountID          *uuid.UUID `json:"account_id"`
	CategoryID         *uuid.UUID `json:"category_id"`
	Description        string     `json:"description"`
	Amount             float64    `json:"amount"`
	TransactionType    int        `json:"transaction_type"`
	PaymentMethod      int        `json:"payment_method"`
	TransactionDate    time.Time  `json:"transaction_date"`
	MerchantName       *string    `json:"merchant_name"`
	Location           *string    `json:"location"`
	Notes              *string    `json:"notes"`
	IsRecurring        *bool      `json:"is_recurring"`
	RecurrenceInterval *int       `json:"recurrence_interval"`
	RecurrenceEndDate  *time.Time `json:"recurrence_end_date"`
}

type TransationSimple struct {
	TransactionID   uuid.UUID `json:"transaction_id"`
	Description     string    `json:"description"`
	Amount          float64   `json:"amount"`
	TransactionType int       `json:"transaction_type"`
	PaymentMethod   int       `json:"payment_method"`
	TransactionDate time.Time `json:"transaction_date"`
	CreditStatus    int       `json:"credit_status"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// Constants for transaction types
const (
	TransactionTypeIncome  = 1
	TransactionTypeExpense = 2
)

// Constants for payment methods
const (
	PaymentMethodDebitQRIS = 1
	PaymentMethodCredit    = 2
	PaymentMethodCash      = 3
	PaymentMethodTransfer  = 4
)

// Constants for recurrence intervals
const (
	RecurrenceIntervalDaily   = 1
	RecurrenceIntervalWeekly  = 2
	RecurrenceIntervalMonthly = 3
	RecurrenceIntervalYearly  = 4
)

// Constants for credit status
const (
	CreditStatusPaid   = 1
	CreditStatusUnpaid = 2
)
