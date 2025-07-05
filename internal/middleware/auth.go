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
	UserID         uuid.UUID `json:"user_id"`
	OrganizationID uuid.UUID `json:"organization_id"`
	Username       string    `json:"username"`
	RoleID         int64     `json:"role_id"`
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

func (m *AuthMiddleware) GenerateToken(userID uuid.UUID, organizationID uuid.UUID, username string, roleID int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, JWTClaims{
		UserID:         userID,
		OrganizationID: organizationID,
		Username:       username,
		RoleID:         roleID,
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
		c.Set("username", claims.Username)
		c.Set("role_id", claims.RoleID)
		c.Set("organization_id", claims.OrganizationID)
		c.Next()
	}
}

func (m *AuthMiddleware) RoleRequired(roleIDs ...int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleID, exists := c.Get("role_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "role not found in context"})
			c.Abort()
			return
		}

		roleIDInt64, ok := roleID.(int64)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid role type"})
			c.Abort()
			return
		}

		for _, requiredRoleID := range roleIDs {
			if roleIDInt64 == requiredRoleID {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
		c.Abort()
	}
}

// SuperAdminRequired middleware that checks if user has SuperAdmin role
func (m *AuthMiddleware) SuperAdminRequired() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// First check if user is authenticated
		m.AuthRequired()(c)
		if c.IsAborted() {
			return
		}

		// Then check if user has SuperAdmin role
		roleID, exists := c.Get("role_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "role not found in context"})
			c.Abort()
			return
		}

		roleIDInt64, ok := roleID.(int64)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid role type"})
			c.Abort()
			return
		}

		if roleIDInt64 != entities.RoleSuperAdmin {
			c.JSON(http.StatusForbidden, gin.H{"error": "superadmin access required"})
			c.Abort()
			return
		}

		c.Next()
	})
}
