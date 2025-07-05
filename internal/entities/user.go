package entities

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user in the system
type User struct {
	UserID         uuid.UUID `json:"user_id" db:"user_id"`
	OrganizationID uuid.UUID `json:"organization_id" db:"organization_id"`
	RoleID         int       `json:"role_id" db:"role_id"`
	UserFullName   string    `json:"full_name" db:"user_fullname"`
	PhoneNumber    string    `json:"phone_number" db:"phone_number"`
	Username       string    `json:"username" db:"username"`
	Password       string    `json:"-" db:"password"` // Password is not exposed in JSON
	IsActive       bool      `json:"is_active" db:"is_active"`
	AccessToken    string    `json:"access_token,omitempty" db:"access_token"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

// CreateUserInput is used for creating a new user
type CreateUserInput struct {
	OrganizationID uuid.UUID `json:"organization_id" binding:"required"`
	RoleID         int       `json:"role_id" binding:"required"`
	UserFullName   string    `json:"user_fullname"`
	PhoneNumber    string    `json:"phone_number" binding:"required"`
	Username       string    `json:"username" binding:"required"`
	Password       string    `json:"password" binding:"required"`
	IsActive       bool      `json:"is_active"`
}

// UpdateUserInput is used for updating an existing user
type UpdateUserInput struct {
	RoleID       int    `json:"role_id"`
	UserFullName string `json:"user_fullname"`
	PhoneNumber  string `json:"phone_number"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	IsActive     *bool  `json:"is_active"`
}

type ResetPasswordInput struct {
	UserID      uuid.UUID `json:"user_id" binding:"required"`
	OldPassword string    `json:"old_password" binding:"required"`
	NewPassword string    `json:"new_password" binding:"required"`
}

type GenerateUserPasswordInput struct {
	Password string `json:"password" binding:"required"`
}

type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	AccessToken    string    `json:"access_token"`
	UserID         uuid.UUID `json:"user_id"`
	OrganizationID uuid.UUID `json:"organization_id"`
	Username       string    `json:"username"`
	RoleID         int       `json:"role_id"`
	ExpiresAt      time.Time `json:"expires_at"`
}
