package entities

import "time"

// Currency represents a currency
type Currency struct {
	CurrencyID            int       `json:"currency_id" db:"currency_id"`
	CurrencyCode          string    `json:"currency_code" db:"currency_code"`
	CurrencyName          string    `json:"currency_name" db:"currency_name"`
	CurrencySymbol        string    `json:"currency_symbol" db:"currency_symbol"`
	CurrencyDecimalPlaces int       `json:"currency_decimal_places" db:"currency_decimal_places"`
	CurrencyStatus        int       `json:"currency_status" db:"currency_status"`
	CreatedAt             time.Time `json:"created_at" db:"created_at"`
	UpdatedAt             time.Time `json:"updated_at" db:"updated_at"`
}

type CreateCurrencyInput struct {
	CurrencyCode          string `json:"currency_code" db:"currency_code"`
	CurrencyName          string `json:"currency_name" db:"currency_name"`
	CurrencySymbol        string `json:"currency_symbol" db:"currency_symbol"`
	CurrencyDecimalPlaces int    `json:"currency_decimal_places" db:"currency_decimal_places"`
	CurrencyStatus        int    `json:"currency_status" db:"currency_status"`
}

type UpdateCurrencyInput struct {
	CurrencyCode          string `json:"currency_code" db:"currency_code"`
	CurrencyName          string `json:"currency_name" db:"currency_name"`
	CurrencySymbol        string `json:"currency_symbol" db:"currency_symbol"`
	CurrencyDecimalPlaces int    `json:"currency_decimal_places" db:"currency_decimal_places"`
	CurrencyStatus        int    `json:"currency_status" db:"currency_status"`
}

type CurrencySimple struct {
	CurrencyID            int    `json:"currency_id" db:"currency_id"`
	CurrencyCode          string `json:"currency_code" db:"currency_code"`
	CurrencyName          string `json:"currency_name" db:"currency_name"`
	CurrencySymbol        string `json:"currency_symbol" db:"currency_symbol"`
	CurrencyDecimalPlaces int    `json:"currency_decimal_places" db:"currency_decimal_places"`
}
