package v1

import (
	"net/http"
	"strconv"

	"github.com/vasst-id/vasst-expense-api/internal/entities"
	"github.com/vasst-id/vasst-expense-api/internal/middleware"
	"github.com/vasst-id/vasst-expense-api/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

	h := handler.Group("/conversations")
	{
		// Protected routes (auth required)
		authGroup := h.Group("")
		authGroup.Use(r.auth.AuthRequired())
		{
			// Organization-scoped routes
			authGroup.GET("", r.ListConversationsByOrganization)
			authGroup.GET("/count", r.GetConversationCountByOrganization)
			authGroup.GET("/:id", r.GetConversationByID)
			authGroup.GET("/:id/detail", r.GetConversationDetail)
			authGroup.POST("", r.CreateConversation)
			authGroup.PUT("/:id", r.UpdateConversation)
			authGroup.DELETE("/:id", r.DeleteConversation)

			// Filtered routes
			authGroup.GET("/user/:user_id", r.GetConversationsByUserID)
			authGroup.GET("/contact/:contact_id", r.GetConversationsByContactID)
			authGroup.GET("/status/:status", r.GetConversationsByStatus)
			authGroup.GET("/priority/:priority", r.GetConversationsByPriority)

			// Active conversation route
			authGroup.GET("/active/:user_id/:contact_id/:medium_id", r.GetActiveConversation)
		}
	}
}

// @Summary Get conversations by organization
// @Description Get a list of conversations for the authenticated organization
// @Tags conversations
// @Accept json
// @Produce json
// @Param limit query int false "Limit for pagination (default: 10)"
// @Param offset query int false "Offset for pagination (default: 0)"
// @Param status query int false "Filter by status"
// @Param priority query int false "Filter by priority"
// @Param is_active query bool false "Filter by active status"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 401 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /conversations [get]
func (r *conversationRoutes) ListConversationsByOrganization(c *gin.Context) {
	// Get organization ID from context (set by auth middleware)
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

	// Check for filters
	var status, priority *int
	var isActive *bool

	if statusStr := c.Query("status"); statusStr != "" {
		if val, err := strconv.Atoi(statusStr); err == nil {
			status = &val
		}
	}

	if priorityStr := c.Query("priority"); priorityStr != "" {
		if val, err := strconv.Atoi(priorityStr); err == nil {
			priority = &val
		}
	}

	if isActiveStr := c.Query("is_active"); isActiveStr != "" {
		if val, err := strconv.ParseBool(isActiveStr); err == nil {
			isActive = &val
		}
	}

	var conversations interface{}
	var err error

	if status != nil || priority != nil || isActive != nil {
		conversations, err = r.conversationService.ListConversationsByOrganizationWithFilters(c.Request.Context(), orgID, status, priority, isActive, limit, offset)
	} else {
		conversations, err = r.conversationService.ListConversationsByOrganization(c.Request.Context(), orgID, limit, offset)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, &entities.ApiResponse{
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

// @Summary Get conversation count by organization
// @Description Get total count of conversations for the authenticated organization
// @Tags conversations
// @Accept json
// @Produce json
// @Success 200 {object} entities.ApiResponse
// @Failure 401 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /conversations/count [get]
func (r *conversationRoutes) GetConversationCountByOrganization(c *gin.Context) {
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

	count, err := r.conversationService.GetConversationCountByOrganization(c.Request.Context(), orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &entities.ApiResponse{
		Success: true,
		Data: map[string]interface{}{
			"count": count,
		},
	})
}

// @Summary Get conversation by ID
// @Description Get a conversation by its ID (organization-scoped)
// @Tags conversations
// @Accept json
// @Produce json
// @Param id path string true "Conversation ID"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 401 {object} entities.ApiResponse
// @Failure 404 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /conversations/{id} [get]
func (r *conversationRoutes) GetConversationByID(c *gin.Context) {
	conversationID, err := uuid.Parse(c.Param("id"))
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

	conversation, err := r.conversationService.GetConversationByID(c.Request.Context(), conversationID, orgID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "conversation not found" {
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
		Data:    conversation,
	})
}

// @Summary Get conversation detail with messages
// @Description Get a conversation with its messages (organization-scoped)
// @Tags conversations
// @Accept json
// @Produce json
// @Param id path string true "Conversation ID"
// @Param message_limit query int false "Message limit (default: 50)"
// @Param message_offset query int false "Message offset (default: 0)"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 401 {object} entities.ApiResponse
// @Failure 404 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /conversations/{id}/detail [get]
func (r *conversationRoutes) GetConversationDetail(c *gin.Context) {
	conversationID, err := uuid.Parse(c.Param("id"))
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

	messageLimit := 50
	messageOffset := 0

	if limitStr := c.Query("message_limit"); limitStr != "" {
		if val, err := strconv.Atoi(limitStr); err == nil && val > 0 {
			messageLimit = val
		}
	}

	if offsetStr := c.Query("message_offset"); offsetStr != "" {
		if val, err := strconv.Atoi(offsetStr); err == nil && val >= 0 {
			messageOffset = val
		}
	}

	conversationDetail, err := r.conversationService.GetConversationDetail(c.Request.Context(), conversationID, orgID, messageLimit, messageOffset)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "conversation not found" {
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
		Data:    conversationDetail,
	})
}

// @Summary Create a new conversation
// @Description Create a new conversation with the provided details
// @Tags conversations
// @Accept json
// @Produce json
// @Param input body entities.CreateConversationInput true "Conversation details"
// @Success 201 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 401 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /conversations [post]
func (r *conversationRoutes) CreateConversation(c *gin.Context) {
	var input entities.CreateConversationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Get organization ID from context and override input
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

	conversation, err := r.conversationService.CreateConversation(c.Request.Context(), &input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, &entities.ApiResponse{
		Success: true,
		Data:    conversation,
	})
}

// @Summary Update a conversation
// @Description Update an existing conversation's details
// @Tags conversations
// @Accept json
// @Produce json
// @Param id path string true "Conversation ID"
// @Param input body entities.UpdateConversationInput true "Updated conversation details"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 401 {object} entities.ApiResponse
// @Failure 404 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /conversations/{id} [put]
func (r *conversationRoutes) UpdateConversation(c *gin.Context) {
	conversationID, err := uuid.Parse(c.Param("id"))
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

	var input entities.UpdateConversationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	conversation, err := r.conversationService.UpdateConversation(c.Request.Context(), conversationID, orgID, &input)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "conversation not found" {
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
		Data:    conversation,
	})
}

// @Summary Delete a conversation
// @Description Delete a conversation by its ID
// @Tags conversations
// @Accept json
// @Produce json
// @Param id path string true "Conversation ID"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 401 {object} entities.ApiResponse
// @Failure 404 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /conversations/{id} [delete]
func (r *conversationRoutes) DeleteConversation(c *gin.Context) {
	conversationID, err := uuid.Parse(c.Param("id"))
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

	err = r.conversationService.DeleteConversation(c.Request.Context(), conversationID, orgID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "conversation not found" {
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
		Data:    "Conversation deleted successfully",
	})
}

// @Summary Get conversations by user ID
// @Description Get all conversations for a specific user (organization-scoped)
// @Tags conversations
// @Accept json
// @Produce json
// @Param user_id path string true "User ID"
// @Param limit query int false "Limit for pagination (default: 10)"
// @Param offset query int false "Offset for pagination (default: 0)"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 401 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /conversations/user/{user_id} [get]
func (r *conversationRoutes) GetConversationsByUserID(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "invalid user ID format",
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

	conversations, err := r.conversationService.GetConversationsByUserID(c.Request.Context(), orgID, userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &entities.ApiResponse{
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

// @Summary Get conversations by contact ID
// @Description Get all conversations for a specific contact (organization-scoped)
// @Tags conversations
// @Accept json
// @Produce json
// @Param contact_id path string true "Contact ID"
// @Param limit query int false "Limit for pagination (default: 10)"
// @Param offset query int false "Offset for pagination (default: 0)"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 401 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /conversations/contact/{contact_id} [get]
func (r *conversationRoutes) GetConversationsByContactID(c *gin.Context) {
	contactID, err := uuid.Parse(c.Param("contact_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "invalid contact ID format",
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

	conversations, err := r.conversationService.GetConversationsByContactID(c.Request.Context(), orgID, contactID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &entities.ApiResponse{
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

// @Summary Get conversations by status
// @Description Get conversations by status (organization-scoped)
// @Tags conversations
// @Accept json
// @Produce json
// @Param status path int true "Status (0: Open, 1: Closed, 2: Pending, 3: Resolved)"
// @Param limit query int false "Limit for pagination (default: 10)"
// @Param offset query int false "Offset for pagination (default: 0)"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 401 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /conversations/status/{status} [get]
func (r *conversationRoutes) GetConversationsByStatus(c *gin.Context) {
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

	conversations, err := r.conversationService.GetConversationsByStatus(c.Request.Context(), orgID, status, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &entities.ApiResponse{
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

// @Summary Get conversations by priority
// @Description Get conversations by priority (organization-scoped)
// @Tags conversations
// @Accept json
// @Produce json
// @Param priority path int true "Priority (0: Low, 1: Medium, 2: High, 3: Urgent)"
// @Param limit query int false "Limit for pagination (default: 10)"
// @Param offset query int false "Offset for pagination (default: 0)"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 401 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /conversations/priority/{priority} [get]
func (r *conversationRoutes) GetConversationsByPriority(c *gin.Context) {
	priority, err := strconv.Atoi(c.Param("priority"))
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "invalid priority format",
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

	conversations, err := r.conversationService.GetConversationsByPriority(c.Request.Context(), orgID, priority, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &entities.ApiResponse{
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

// @Summary Get active conversation
// @Description Get the active conversation for specific user, contact, and medium
// @Tags conversations
// @Accept json
// @Produce json
// @Param user_id path string true "User ID"
// @Param contact_id path string true "Contact ID"
// @Param medium_id path int true "Medium ID"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 401 {object} entities.ApiResponse
// @Failure 404 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /conversations/active/{user_id}/{contact_id}/{medium_id} [get]
func (r *conversationRoutes) GetActiveConversation(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "invalid user ID format",
		})
		return
	}

	contactID, err := uuid.Parse(c.Param("contact_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "invalid contact ID format",
		})
		return
	}

	mediumID, err := strconv.Atoi(c.Param("medium_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "invalid medium ID format",
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

	conversation, err := r.conversationService.GetActiveConversation(c.Request.Context(), orgID, userID, contactID, mediumID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	if conversation == nil {
		c.JSON(http.StatusNotFound, &entities.ApiResponse{
			Success: false,
			Error:   "no active conversation found",
		})
		return
	}

	c.JSON(http.StatusOK, &entities.ApiResponse{
		Success: true,
		Data:    conversation,
	})
}
