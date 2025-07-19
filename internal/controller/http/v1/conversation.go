package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vasst-id/vasst-expense-api/internal/entities"
	"github.com/vasst-id/vasst-expense-api/internal/middleware"
	"github.com/vasst-id/vasst-expense-api/internal/services"
)

type conversationRoutes struct {
	conversationService services.ConversationService
	auth                *middleware.AuthMiddleware
}

func newConversationRoutes(handler *gin.RouterGroup, conversationService services.ConversationService, auth *middleware.AuthMiddleware) {
	r := &conversationRoutes{
		conversationService: conversationService,
		auth:                auth,
	}

	// Conversation endpoints - all require authentication
	conversations := handler.Group("/conversations")
	conversations.Use(auth.AuthRequired())
	{
		conversations.GET("/active", r.GetActiveConversationsByUserID)
	}
}

// @Summary Get active conversations
// @Description Get all active conversations for the authenticated user
// @Tags conversations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} entities.ApiResponse
// @Failure 401 {object} entities.ApiResponse
// @Failure 404 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /conversations/active [get]
func (r *conversationRoutes) GetActiveConversationsByUserID(c *gin.Context) {
	userID, ok := GetAuthenticatedUserID(c)
	if !ok {
		return
	}

	conversations, err := r.conversationService.GetActiveConversationsByUserID(c.Request.Context(), userID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "user not found" {
			status = http.StatusNotFound
		}
		c.JSON(status, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &entities.ApiResponse{
		Success: true,
		Data:    conversations,
	})
}
