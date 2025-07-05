package entities

import (
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

// JWTClaims represents the JWT claims structure
type JWTClaims struct {
	UserID         uuid.UUID `json:"user_id"`
	OrganizationID uuid.UUID `json:"organization_id"`
	Username       string    `json:"username"`
	Role           string    `json:"role"`
	jwt.StandardClaims
}
