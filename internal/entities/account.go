package entities

import (
	"time"

	"github.com/google/uuid"
)

// Account represents a financial account
type Account struct {
	AccountID      uuid.UUID `json:"account_id" db:"account_id"`
	UserID         uuid.UUID `json:"user_id" db:"user_id"`
	AccountName    string    `json:"account_name" db:"account_name"`
	AccountType    int       `json:"account_type" db:"account_type"`
	BankID         *int      `json:"bank_id" db:"bank_id"`
	AccountNumber  *string   `json:"account_number" db:"account_number"`
	CurrentBalance float64   `json:"current_balance" db:"current_balance"`
	CreditLimit    *float64  `json:"credit_limit" db:"credit_limit"`
	DueDate        *int      `json:"due_date" db:"due_date"`
	CurrencyID     int       `json:"currency_id" db:"currency_id"`
	IsActive       bool      `json:"is_active" db:"is_active"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

// CreateAccountRequest represents the create account request
type CreateAccountRequest struct {
	UserID         uuid.UUID `json:"user_id"`
	AccountName    string    `json:"account_name" binding:"required"`
	AccountType    int       `json:"account_type" binding:"required"`
	BankID         *int      `json:"bank_id"`
	AccountNumber  *string   `json:"account_number"`
	CurrentBalance float64   `json:"current_balance"`
	CreditLimit    *float64  `json:"credit_limit"`
	DueDate        *int      `json:"due_date"`
	CurrencyID     int       `json:"currency_id" binding:"required"`
}

type UpdateAccountRequest struct {
	AccountID      uuid.UUID `json:"account_id" binding:"required"`
	AccountName    string    `json:"account_name" binding:"required"`
	AccountType    int       `json:"account_type" binding:"required"`
	BankID         *int      `json:"bank_id"`
	AccountNumber  *string   `json:"account_number"`
	CurrentBalance float64   `json:"current_balance"`
	CreditLimit    *float64  `json:"credit_limit"`
	DueDate        *int      `json:"due_date"`
	CurrencyID     int       `json:"currency_id" binding:"required"`
}

// Constants for account types
const (
	AccountTypeDebit   = 1
	AccountTypeCredit  = 2
	AccountTypeSavings = 3
	AccountTypeCash    = 4
	AccountTypeShared  = 5
)

type AccountSimple struct {
	AccountID        uuid.UUID `json:"account_id" db:"account_id"`
	AccountName      string    `json:"account_name" db:"account_name"`
	AccountTypeLabel string    `json:"account_type_label" db:"account_type_label"`
	BankName         string    `json:"bank_name" db:"bank_name"`
	AccountNumber    *string   `json:"account_number" db:"account_number"`
	CurrentBalance   float64   `json:"current_balance" db:"current_balance"`
	CurrencySymbol   string    `json:"currency_symbol" db:"currency_symbol"`
	CreditLimit      *float64  `json:"credit_limit" db:"credit_limit"`
	DueDate          *int      `json:"due_date" db:"due_date"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}
