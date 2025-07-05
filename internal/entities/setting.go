package entities

import "time"

// Setting represents application settings in the system
type Setting struct {
	SettingID            int       `json:"setting_id" db:"setting_id"`
	BusinessName         string    `json:"business_name" db:"business_name"`
	BusinessAddress      string    `json:"business_address" db:"business_address"`
	BusinessPhone        string    `json:"business_phone" db:"business_phone"`
	BusinessEmail        string    `json:"business_email" db:"business_email"`
	TaxRate              float64   `json:"tax_rate" db:"tax_rate"`
	DefaultDeliveryFee   float64   `json:"default_delivery_fee" db:"default_delivery_fee"`
	Logo                 string    `json:"logo" db:"logo"`
	InvoiceFooterText    string    `json:"invoice_footer_text" db:"invoice_footer_text"`
	OrderNumberPrefix    string    `json:"order_number_prefix" db:"order_number_prefix"`
	BatchNumberPrefix    string    `json:"batch_number_prefix" db:"batch_number_prefix"`
	CurrencySymbol       string    `json:"currency_symbol" db:"currency_symbol"`
	DateFormat           string    `json:"date_format" db:"date_format"`
	TimeZone             string    `json:"time_zone" db:"time_zone"`
	InventoryAlertsEmail string    `json:"inventory_alerts_email" db:"inventory_alerts_email"`
	CreatedAt            time.Time `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time `json:"updated_at" db:"updated_at"`
}
