package entities

import "time"

// Bank represents a financial institution
type Bank struct {
	BankID      int       `json:"bank_id" db:"bank_id"`
	BankName    string    `json:"bank_name" db:"bank_name"`
	BankCode    string    `json:"bank_code" db:"bank_code"`
	BankLogoURL string    `json:"bank_logo_url" db:"bank_logo_url"`
	Status      int       `json:"status" db:"status"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type CreateBankInput struct {
	BankName    string `json:"bank_name" db:"bank_name"`
	BankCode    string `json:"bank_code" db:"bank_code"`
	BankLogoURL string `json:"bank_logo_url" db:"bank_logo_url"`
	Status      int    `json:"status" db:"status"`
}

type UpdateBankInput struct {
	BankName    string `json:"bank_name" db:"bank_name"`
	BankCode    string `json:"bank_code" db:"bank_code"`
	BankLogoURL string `json:"bank_logo_url" db:"bank_logo_url"`
	Status      int    `json:"status" db:"status"`
}

type BankSimple struct {
	BankID      int    `json:"bank_id" db:"bank_id"`
	BankName    string `json:"bank_name" db:"bank_name"`
	BankCode    string `json:"bank_code" db:"bank_code"`
	BankLogoURL string `json:"bank_logo_url" db:"bank_logo_url"`
}
