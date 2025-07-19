package entities

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user in the expense tracking system
type User struct {
	UserID             uuid.UUID  `json:"user_id" db:"user_id"`
	Email              string     `json:"email" db:"email"`
	PhoneNumber        string     `json:"phone_number" db:"phone_number"`
	PasswordHash       string     `json:"-" db:"password_hash"`
	FirstName          string     `json:"first_name" db:"first_name"`
	LastName           string     `json:"last_name" db:"last_name"`
	Timezone           string     `json:"timezone" db:"timezone"`
	CurrencyID         int        `json:"currency_id" db:"currency_id"`
	SubscriptionPlanID int        `json:"subscription_plan_id" db:"subscription_plan_id"`
	EmailVerifiedAt    *time.Time `json:"email_verified_at" db:"email_verified_at"`
	PhoneVerifiedAt    *time.Time `json:"phone_verified_at" db:"phone_verified_at"`
	Status             int        `json:"status" db:"status"`
	CreatedAt          time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at" db:"updated_at"`

	// Virtual fields for response
	AccessToken  string `json:"access_token,omitempty" db:"-"`
	RefreshToken string `json:"refresh_token,omitempty" db:"-"`
}

type CreateUserInput struct {
	PhoneNumber        string `json:"phone_number" binding:"required"`
	Password           string `json:"password" binding:"required,len=6,numeric"`
	FirstName          string `json:"first_name" binding:"required"`
	LastName           string `json:"last_name" binding:"required"`
	Timezone           string `json:"timezone,omitempty"`
	CurrencyID         int    `json:"currency_id,omitempty"`
	SubscriptionPlanID int    `json:"subscription_plan_id,omitempty"`
}

type UpdateUserInput struct {
	Email              string     `json:"email" binding:"email"`
	PhoneNumber        string     `json:"phone_number" binding:"required"`
	FirstName          string     `json:"first_name" binding:"required"`
	LastName           string     `json:"last_name" binding:"required"`
	Timezone           string     `json:"timezone,omitempty"`
	CurrencyID         int        `json:"currency_id" binding:"required"`
	SubscriptionPlanID int        `json:"subscription_plan_id" binding:"required"`
	Status             int        `json:"status" binding:"required"`
	EmailVerifiedAt    *time.Time `json:"email_verified_at" db:"email_verified_at"`
	PhoneVerifiedAt    *time.Time `json:"phone_verified_at" db:"phone_verified_at"`
}

// RefreshToken represents a refresh token for JWT authentication
type RefreshToken struct {
	ID        uuid.UUID `json:"id" db:"id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	Token     string    `json:"token" db:"token"`
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	IsRevoked bool      `json:"is_revoked" db:"is_revoked"`
	IPAddress string    `json:"ip_address" db:"ip_address"`
	UserAgent string    `json:"user_agent" db:"user_agent"`
}

// JWT Claims for expense system
type JWTClaims struct {
	UserID             uuid.UUID `json:"user_id"`
	DefaultWorkspaceID uuid.UUID `json:"default_workspace_id"`
	Email              string    `json:"email"`
	StandardClaims
}

// StandardClaims represents the standard JWT claims
type StandardClaims struct {
	Issuer    string `json:"iss"`
	Subject   string `json:"sub"`
	Audience  string `json:"aud"`
	ExpiresAt int64  `json:"exp"`
	NotBefore int64  `json:"nbf"`
	IssuedAt  int64  `json:"iat"`
	ID        string `json:"jti"`
}

// LoginRequest represents the login request
type LoginRequest struct {
	PhoneNumber string `json:"phone_number,omitempty" binding:"required"`
	Password    string `json:"password" binding:"required,len=6,numeric"`
}

// LoginResponse represents the login response
type LoginResponse struct {
	AccessToken  string     `json:"access_token"`
	RefreshToken string     `json:"refresh_token"`
	TokenType    string     `json:"token_type"`
	ExpiresIn    int        `json:"expires_in"`
	User         *User      `json:"user"`
	Workspace    *Workspace `json:"workspace,omitempty"`
}

// RefreshTokenRequest represents the refresh token request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// RefreshTokenResponse represents the refresh token response
type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

// ForgotPasswordRequest represents the forgot password request
type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// ResetPasswordRequest represents the reset password request
type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

// ChangePasswordRequest represents the change password request
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=8"`
}

// VerifyPhoneRequest represents the phone verification request
type VerifyPhoneRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
	Code        string `json:"code" binding:"required"`
}

// ResendVerificationCodeRequest represents the resend verification code request
type ResendVerificationCodeRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
}

// VerifyEmailRequest represents the email verification request
type VerifyEmailRequest struct {
	Email string `json:"email" binding:"required,email"`
	Token string `json:"token" binding:"required"`
}

// ResendVerificationEmailRequest represents the resend verification email request
type ResendVerificationEmailRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// ResetPasswordInput represents the reset password input for existing users
type ResetPasswordInput struct {
	UserID      uuid.UUID `json:"user_id" binding:"required"`
	OldPassword string    `json:"old_password" binding:"required"`
	NewPassword string    `json:"new_password" binding:"required,min=8"`
}

// Constants for user status
const (
	UserStatusActive   = 1
	UserStatusInactive = 0
)

// Constants for token types
const (
	TokenTypeAccess  = "access"
	TokenTypeRefresh = "refresh"
)

// ApiResponse represents a generic API response
type ApiResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
}
