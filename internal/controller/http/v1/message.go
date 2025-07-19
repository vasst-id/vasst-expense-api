package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/vasst-id/vasst-expense-api/internal/entities"
	"github.com/vasst-id/vasst-expense-api/internal/middleware"
	"github.com/vasst-id/vasst-expense-api/internal/services"
)

type messageRoutes struct {
	messageService services.MessageService
	auth           *middleware.AuthMiddleware
}

func newMessageRoutes(handler *gin.RouterGroup, messageService services.MessageService, auth *middleware.AuthMiddleware) {
	r := &messageRoutes{
		messageService: messageService,
		auth:           auth,
	}

	// Message endpoints - all require authentication
	messages := handler.Group("/messages")
	messages.Use(auth.AuthRequired())
	{
		messages.GET("/conversation/:conversation_id", r.GetMessagesByConversationID)
	}
}

// @Summary Get messages by conversation
// @Description Get a list of messages for a specific conversation with pagination
// @Tags messages
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param conversation_id path string true "Conversation ID"
// @Param limit query int false "Limit for pagination (default: 10)"
// @Param offset query int false "Offset for pagination (default: 0)"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 401 {object} entities.ApiResponse
// @Failure 403 {object} entities.ApiResponse
// @Failure 404 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /messages/conversation/{conversation_id} [get]
func (r *messageRoutes) GetMessagesByConversationID(c *gin.Context) {
	userID, ok := GetAuthenticatedUserID(c)
	if !ok {
		return
	}

	conversationID, err := uuid.Parse(c.Param("conversation_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "invalid conversation ID format",
		})
		return
	}

	// Parse pagination parameters
	limit := 10
	offset := 0

	if limitStr := c.Query("limit"); limitStr != "" {
		if val, err := strconv.Atoi(limitStr); err == nil && val > 0 {
			limit = val
		}
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if val, err := strconv.Atoi(offsetStr); err == nil && val >= 0 {
			offset = val
		}
	}

	messages, totalCount, err := r.messageService.GetMessagesByConversationID(c.Request.Context(), userID, conversationID, limit, offset)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "conversation not found" {
			status = http.StatusNotFound
		} else if err.Error() == "access denied" {
			status = http.StatusForbidden
		}
		c.JSON(status, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &entities.ApiResponse{
		Success: true,
		Data: map[string]interface{}{
			"messages": messages,
			"total":    totalCount,
			"limit":    limit,
			"offset":   offset,
		},
	})
}
