package entities

import (
	"time"
)

// Customer represents a customer in the system
type Customer struct {
	CustomerID  int       `json:"customer_id" db:"customer_id"`
	UserID      int       `json:"user_id" db:"user_id"`
	FullName    string    `json:"full_name" db:"full_name"`
	PhoneNumber string    `json:"phone_number" db:"phone_number"`
	Address     string    `json:"address" db:"address"`
	Email       string    `json:"email" db:"email"`
	Area        string    `json:"area" db:"area"`
	Notes       string    `json:"notes" db:"notes"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// CustomerListParams is used for filtering and pagination when listing customers
type CustomerListParams struct {
	Offset   int64  `json:"offset"`
	Limit    int64  `json:"limit"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	IsActive bool   `json:"is_active"`
}

// CreateCustomerInput is used for creating a new customer
type CreateCustomerInput struct {
	UserID      int64  `json:"user_id" binding:"required"`
	FullName    string `json:"full_name" binding:"required"`
	PhoneNumber string `json:"phone_number" binding:"required"`
	Address     string `json:"address"`
	Email       string `json:"email"`
	Area        string `json:"area"`
	Notes       string `json:"notes"`
	IsActive    bool   `json:"is_active"`
}

// UpdateCustomerInput is used for updating an existing customer
type UpdateCustomerInput struct {
	FullName    string `json:"full_name"`
	PhoneNumber string `json:"phone_number"`
	Address     string `json:"address"`
	Email       string `json:"email"`
	Area        string `json:"area"`
	Notes       string `json:"notes"`
	IsActive    *bool  `json:"is_active"`
}
