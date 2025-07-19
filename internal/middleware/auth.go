package middleware

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/vasst-id/vasst-expense-api/internal/entities"
)

const (
	AuthorizationHeader = "Authorization"
	BearerPrefix        = "Bearer "
)

type JWTClaims struct {
	UserID      uuid.UUID `json:"user_id"`
	PhoneNumber string    `json:"phone_number"`
	Email       string    `json:"email"`
	Status      int64     `json:"status"`
	jwt.RegisteredClaims
}

type AuthMiddleware struct {
	secretKey []byte
}

func NewAuthMiddleware(secretKey string) *AuthMiddleware {
	return &AuthMiddleware{
		secretKey: []byte(secretKey),
	}
}

func (m *AuthMiddleware) GenerateToken(user *entities.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, JWTClaims{
		UserID:      user.UserID,
		PhoneNumber: user.PhoneNumber,
		Email:       user.Email,
		Status:      int64(user.Status),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * 7 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	})

	tokenString, err := token.SignedString(m.secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (m *AuthMiddleware) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return m.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func (m *AuthMiddleware) AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(AuthorizationHeader)
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header is required"})
			c.Abort()
			return
		}

		if !strings.HasPrefix(authHeader, BearerPrefix) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, BearerPrefix)
		claims, err := m.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("phone_number", claims.PhoneNumber)
		c.Set("email", claims.Email)
		c.Set("status", claims.Status)
		c.Next()
	}
}
