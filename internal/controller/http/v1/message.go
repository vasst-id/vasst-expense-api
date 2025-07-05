package v1

import (
	"fmt"
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

	h := handler.Group("/messages")
	{
		// Protected routes (auth required)
		authGroup := h.Group("")
		authGroup.Use(r.auth.AuthRequired())
		{
			authGroup.GET("", r.GetMessagesByOrganizationID)
			authGroup.GET("/:id", r.GetMessageByID)
			authGroup.POST("", r.CreateMessage)
			authGroup.POST("/with-media", r.CreateMessageWithMedia)
			authGroup.PUT("/:id", r.UpdateMessage)
			authGroup.PUT("/:id/status", r.UpdateMessageStatus)
			authGroup.DELETE("/:id", r.DeleteMessage)
			authGroup.GET("/conversation/:conversation_id", r.GetMessagesByConversationID)
			authGroup.GET("/pending", r.GetPendingMessages)
			authGroup.GET("/status/:status", r.GetMessagesByStatus)
		}
	}
}

// @Summary Get message by ID
// @Description Get a message by its ID (organization-scoped)
// @Tags messages
// @Accept json
// @Produce json
// @Param id path string true "Message ID"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 401 {object} entities.ApiResponse
// @Failure 403 {object} entities.ApiResponse
// @Failure 404 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /messages/{id} [get]
func (r *messageRoutes) GetMessageByID(c *gin.Context) {
	messageID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "invalid message ID format",
		})
		return
	}

	// Get organization ID from context
	organizationID, exists := c.Get("organization_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, &entities.ApiResponse{
			Success: false,
			Error:   "organization ID not found in context",
		})
		return
	}

	orgID, ok := organizationID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "invalid organization ID format",
		})
		return
	}

	message, err := r.messageService.GetMessageByID(c.Request.Context(), messageID, orgID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "message not found" {
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
		Data:    message,
	})
}

// @Summary Create a new message
// @Description Create a new message with the provided details (will create conversation if needed)
// @Tags messages
// @Accept json
// @Produce json
// @Param input body entities.CreateMessageInput true "Message details"
// @Success 201 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 401 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /messages [post]
func (r *messageRoutes) CreateMessage(c *gin.Context) {
	var input entities.CreateMessageInput

	// Get organization ID from context
	organizationID, exists := c.Get("organization_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, &entities.ApiResponse{
			Success: false,
			Error:   "organization ID not found in context",
		})
		return
	}

	orgID, ok := organizationID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "invalid organization ID format",
		})
		return
	}
	input.OrganizationID = orgID

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, &entities.ApiResponse{
			Success: false,
			Error:   "user ID not found in context",
		})
		return
	}
	fmt.Println("userID", userID)

	userIDUUID, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "invalid user ID format",
		})
		return
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	message, err := r.messageService.CreateMessage(c.Request.Context(), &input, userIDUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, &entities.ApiResponse{
		Success: true,
		Data:    message,
	})
}

// @Summary Create a new message with media
// @Description Create a new message with media file upload
// @Tags messages
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Media file"
// @Param conversation_id formData string true "Conversation ID"
// @Param organization_id formData string true "Organization ID"
// @Param direction formData string true "Message direction (i/o)"
// @Param message_type_id formData int true "Message type ID"
// @Param content formData string false "Message content"
// @Param is_broadcast formData bool false "Is broadcast message"
// @Param is_order_message formData bool false "Is order message"
// @Success 201 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /messages/with-media [post]
func (r *messageRoutes) CreateMessageWithMedia(c *gin.Context) {
	// Get the uploaded file
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "file is required",
		})
		return
	}

	// Parse form data
	conversationID, err := uuid.Parse(c.PostForm("conversation_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "invalid conversation ID format",
		})
		return
	}

	organizationID, err := uuid.Parse(c.PostForm("organization_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "invalid organization ID format",
		})
		return
	}

	direction := c.PostForm("direction")
	if direction == "" {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "direction is required",
		})
		return
	}

	messageTypeID, err := strconv.Atoi(c.PostForm("message_type_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "invalid message type ID",
		})
		return
	}

	content := c.PostForm("content")
	isBroadcast := c.PostForm("is_broadcast") == "true"
	isOrderMessage := c.PostForm("is_order_message") == "true"

	input := &entities.CreateMessageInput{
		ConversationID: conversationID,
		OrganizationID: organizationID,
		Direction:      direction,
		MessageTypeID:  messageTypeID,
		Content:        content,
		IsBroadcast:    isBroadcast,
		IsOrderMessage: isOrderMessage,
	}

	message, err := r.messageService.CreateMessageWithMedia(c.Request.Context(), input, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, &entities.ApiResponse{
		Success: true,
		Data:    message,
	})
}

// @Summary Update a message
// @Description Update an existing message's details
// @Tags messages
// @Accept json
// @Produce json
// @Param id path string true "Message ID"
// @Param input body entities.UpdateMessageInput true "Updated message details"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 401 {object} entities.ApiResponse
// @Failure 403 {object} entities.ApiResponse
// @Failure 404 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /messages/{id} [put]
func (r *messageRoutes) UpdateMessage(c *gin.Context) {
	messageID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "invalid message ID format",
		})
		return
	}

	// Get organization ID from context
	organizationID, exists := c.Get("organization_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, &entities.ApiResponse{
			Success: false,
			Error:   "organization ID not found in context",
		})
		return
	}

	orgID, ok := organizationID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "invalid organization ID format",
		})
		return
	}

	var input entities.UpdateMessageInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	message, err := r.messageService.UpdateMessage(c.Request.Context(), messageID, orgID, &input)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "message not found" {
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
		Data:    message,
	})
}

// @Summary Update message status
// @Description Update only the status of a message
// @Tags messages
// @Accept json
// @Produce json
// @Param id path string true "Message ID"
// @Param input body entities.UpdateMessageStatusInput true "Status update details"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 401 {object} entities.ApiResponse
// @Failure 403 {object} entities.ApiResponse
// @Failure 404 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /messages/{id}/status [put]
func (r *messageRoutes) UpdateMessageStatus(c *gin.Context) {
	messageID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "invalid message ID format",
		})
		return
	}

	// Get organization ID from context
	organizationID, exists := c.Get("organization_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, &entities.ApiResponse{
			Success: false,
			Error:   "organization ID not found in context",
		})
		return
	}

	orgID, ok := organizationID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "invalid organization ID format",
		})
		return
	}

	var input entities.UpdateMessageStatusInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	err = r.messageService.UpdateMessageStatus(c.Request.Context(), messageID, orgID, &input)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "message not found" {
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
		Data:    "Message status updated successfully",
	})
}

// @Summary Delete a message
// @Description Delete a message by its ID
// @Tags messages
// @Accept json
// @Produce json
// @Param id path string true "Message ID"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 401 {object} entities.ApiResponse
// @Failure 403 {object} entities.ApiResponse
// @Failure 404 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /messages/{id} [delete]
func (r *messageRoutes) DeleteMessage(c *gin.Context) {
	messageID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "invalid message ID format",
		})
		return
	}

	// Get organization ID from context
	organizationID, exists := c.Get("organization_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, &entities.ApiResponse{
			Success: false,
			Error:   "organization ID not found in context",
		})
		return
	}

	orgID, ok := organizationID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "invalid organization ID format",
		})
		return
	}

	err = r.messageService.DeleteMessage(c.Request.Context(), messageID, orgID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "message not found" {
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
		Data:    "Message deleted successfully",
	})
}

// @Summary Get messages by conversation ID
// @Description Get all messages for a conversation (organization-scoped)
// @Tags messages
// @Accept json
// @Produce json
// @Param conversation_id path string true "Conversation ID"
// @Param limit query int false "Limit for pagination (default: 50)"
// @Param offset query int false "Offset for pagination (default: 0)"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 401 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /messages/conversation/{conversation_id} [get]
func (r *messageRoutes) GetMessagesByConversationID(c *gin.Context) {
	conversationID, err := uuid.Parse(c.Param("conversation_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "invalid conversation ID format",
		})
		return
	}

	// Get organization ID from context
	organizationID, exists := c.Get("organization_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, &entities.ApiResponse{
			Success: false,
			Error:   "organization ID not found in context",
		})
		return
	}

	orgID, ok := organizationID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "invalid organization ID format",
		})
		return
	}

	limit := 50
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

	messages, err := r.messageService.ListMessagesByConversation(c.Request.Context(), conversationID, orgID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &entities.ApiResponse{
		Success: true,
		Data:    messages,
	})
}

// @Summary Get messages by organization ID
// @Description Get all messages for a specific organization
// @Tags messages
// @Accept json
// @Produce json
// @Param org_id path string true "Organization ID"
// @Param limit query int false "Limit for pagination (default: 50)"
// @Param offset query int false "Offset for pagination (default: 0)"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 401 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /messages/organization/{org_id} [get]
func (r *messageRoutes) GetMessagesByOrganizationID(c *gin.Context) {
	// Get organization ID from context
	organizationID, exists := c.Get("organization_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, &entities.ApiResponse{
			Success: false,
			Error:   "organization ID not found in context",
		})
		return
	}

	orgID, ok := organizationID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "invalid organization ID format",
		})
		return
	}

	limit := 50
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

	messages, err := r.messageService.ListMessagesByOrganization(c.Request.Context(), orgID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &entities.ApiResponse{
		Success: true,
		Data:    messages,
	})
}

// @Summary Get pending messages
// @Description Get all pending messages for an organization
// @Tags messages
// @Accept json
// @Produce json
// @Param limit query int false "Limit for pagination (default: 50)"
// @Param offset query int false "Offset for pagination (default: 0)"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 401 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /messages/pending [get]
func (r *messageRoutes) GetPendingMessages(c *gin.Context) {
	// Get organization ID from context
	organizationID, exists := c.Get("organization_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, &entities.ApiResponse{
			Success: false,
			Error:   "organization ID not found in context",
		})
		return
	}

	orgID, ok := organizationID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "invalid organization ID format",
		})
		return
	}

	limit := 50
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

	messages, err := r.messageService.GetPendingMessages(c.Request.Context(), orgID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &entities.ApiResponse{
		Success: true,
		Data:    messages,
	})
}

// @Summary Get messages by status
// @Description Get messages by status for an organization
// @Tags messages
// @Accept json
// @Produce json
// @Param status path int true "Message status (0: Pending, 1: Sent, 2: Delivered, 3: Read, 4: Failed)"
// @Param limit query int false "Limit for pagination (default: 50)"
// @Param offset query int false "Offset for pagination (default: 0)"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 401 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /messages/status/{status} [get]
func (r *messageRoutes) GetMessagesByStatus(c *gin.Context) {
	status, err := strconv.Atoi(c.Param("status"))
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "invalid status format",
		})
		return
	}

	// Get organization ID from context
	organizationID, exists := c.Get("organization_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, &entities.ApiResponse{
			Success: false,
			Error:   "organization ID not found in context",
		})
		return
	}

	orgID, ok := organizationID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "invalid organization ID format",
		})
		return
	}

	limit := 50
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

	messages, err := r.messageService.GetMessagesByStatus(c.Request.Context(), orgID, status, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &entities.ApiResponse{
		Success: true,
		Data:    messages,
	})
}
