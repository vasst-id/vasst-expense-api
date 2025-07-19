package entities

import (
	"time"

	"github.com/google/uuid"
)

type VerificationCode struct {
	VerificationCodeID uuid.UUID `json:"verification_code_id" db:"verification_code_id"`
	PhoneNumber        string    `json:"phone_number" db:"phone_number"`
	Code               string    `json:"code" db:"code"`
	CodeType           string    `json:"code_type" db:"code_type"`
	ExpiresAt          time.Time `json:"expires_at" db:"expires_at"`
	IsUsed             bool      `json:"is_used" db:"is_used"`
	AttemptsCount      int       `json:"attempts_count" db:"attempts_count"`
	MaxAttempts        int       `json:"max_attempts" db:"max_attempts"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`
}

type CreateVerificationCodeRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
	CodeType    string `json:"code_type" binding:"required"`
}

type VerifyVerificationCodeRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
	Code        string `json:"code" binding:"required"`
}
