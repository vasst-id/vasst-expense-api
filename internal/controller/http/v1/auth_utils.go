package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/vasst-id/vasst-expense-api/internal/entities"
)

// GetAuthenticatedUserID extracts and validates the authenticated user ID from the Gin context
// Returns the user ID and a boolean indicating success
// If extraction fails, it automatically sends an error response and returns false
func GetAuthenticatedUserID(c *gin.Context) (uuid.UUID, bool) {
	userIDInterface, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, &entities.ApiResponse{
			Success: false,
			Error:   "user authentication required",
		})
		return uuid.Nil, false
	}

	// Handle both uuid.UUID and string types
	var userID uuid.UUID
	switch v := userIDInterface.(type) {
	case uuid.UUID:
		userID = v
	case string:
		parsedID, err := uuid.Parse(v)
		if err != nil {
			c.JSON(http.StatusInternalServerError, &entities.ApiResponse{
				Success: false,
				Error:   "invalid user ID format",
			})
			return uuid.Nil, false
		}
		userID = parsedID
	default:
		c.JSON(http.StatusInternalServerError, &entities.ApiResponse{
			Success: false,
			Error:   "invalid user ID format in context",
		})
		return uuid.Nil, false
	}

	return userID, true
}

// GetOptionalAuthenticatedUserID extracts the authenticated user ID from context without sending error responses
// Returns the user ID and a boolean indicating if extraction was successful
// Use this when user authentication is optional
func GetOptionalAuthenticatedUserID(c *gin.Context) (uuid.UUID, bool) {
	userIDInterface, ok := c.Get("user_id")
	if !ok {
		return uuid.Nil, false
	}

	// Handle both uuid.UUID and string types
	var userID uuid.UUID
	switch v := userIDInterface.(type) {
	case uuid.UUID:
		userID = v
	case string:
		parsedID, err := uuid.Parse(v)
		if err != nil {
			return uuid.Nil, false
		}
		userID = parsedID
	default:
		return uuid.Nil, false
	}

	return userID, true
}
