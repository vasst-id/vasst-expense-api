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

type userTagsRoutes struct {
	userTagsService services.UserTagsService
	auth            *middleware.AuthMiddleware
}

func newUserTagsRoutes(handler *gin.RouterGroup, userTagsService services.UserTagsService, auth *middleware.AuthMiddleware) {
	r := &userTagsRoutes{
		userTagsService: userTagsService,
		auth:            auth,
	}

	// User tags endpoints - all require authentication
	userTags := handler.Group("/user-tags")
	userTags.Use(auth.AuthRequired())
	{
		userTags.GET("", r.GetUserTags)
		userTags.POST("", r.CreateUserTag)
		userTags.GET("/active", r.GetActiveUserTags)
		userTags.GET("/:id", r.GetUserTagByID)
		userTags.PUT("/:id", r.UpdateUserTag)
		userTags.DELETE("/:id", r.DeleteUserTag)
	}
}

// @Summary Get user tags
// @Description Get a list of user tags with pagination
// @Tags user-tags
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Limit for pagination"
// @Param offset query int false "Offset for pagination"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 401 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /user-tags [get]
func (r *userTagsRoutes) GetUserTags(c *gin.Context) {
	userID, ok := GetAuthenticatedUserID(c)
	if !ok {
		return
	}

	limit := 10
	offset := 0

	if limitStr := c.Query("limit"); limitStr != "" {
		if val, err := strconv.Atoi(limitStr); err == nil {
			limit = val
		}
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if val, err := strconv.Atoi(offsetStr); err == nil {
			offset = val
		}
	}

	userTags, err := r.userTagsService.GetUserTags(c.Request.Context(), userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &entities.ApiResponse{
		Success: true,
		Data:    userTags,
	})
}

// @Summary Create user tag
// @Description Create a new user tag
// @Tags user-tags
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body entities.CreateUserTagRequest true "User tag details"
// @Success 201 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 401 {object} entities.ApiResponse
// @Failure 409 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /user-tags [post]
func (r *userTagsRoutes) CreateUserTag(c *gin.Context) {
	userID, ok := GetAuthenticatedUserID(c)
	if !ok {
		return
	}

	var input entities.CreateUserTagRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	userTag, err := r.userTagsService.CreateUserTag(c.Request.Context(), userID, &input)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "tag with this name already exists" {
			status = http.StatusConflict
		} else if err.Error() == "tag name is required" {
			status = http.StatusBadRequest
		}
		c.JSON(status, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, &entities.ApiResponse{
		Success: true,
		Data:    userTag,
	})
}

// @Summary Get active user tags
// @Description Get all active user tags for the authenticated user
// @Tags user-tags
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} entities.ApiResponse
// @Failure 401 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /user-tags/active [get]
func (r *userTagsRoutes) GetActiveUserTags(c *gin.Context) {
	userID, ok := GetAuthenticatedUserID(c)
	if !ok {
		return
	}

	userTags, err := r.userTagsService.GetActiveUserTags(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &entities.ApiResponse{
		Success: true,
		Data:    userTags,
	})
}

// @Summary Get user tag by ID
// @Description Get a user tag by its ID
// @Tags user-tags
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User Tag ID"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 401 {object} entities.ApiResponse
// @Failure 403 {object} entities.ApiResponse
// @Failure 404 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /user-tags/{id} [get]
func (r *userTagsRoutes) GetUserTagByID(c *gin.Context) {
	userID, ok := GetAuthenticatedUserID(c)
	if !ok {
		return
	}

	userTagID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "invalid user tag ID format",
		})
		return
	}

	userTag, err := r.userTagsService.GetUserTagByID(c.Request.Context(), userID, userTagID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "user tag not found" {
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
		Data:    userTag,
	})
}

// @Summary Update user tag
// @Description Update an existing user tag
// @Tags user-tags
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User Tag ID"
// @Param input body entities.UpdateUserTagRequest true "Updated user tag details"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 401 {object} entities.ApiResponse
// @Failure 403 {object} entities.ApiResponse
// @Failure 404 {object} entities.ApiResponse
// @Failure 409 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /user-tags/{id} [put]
func (r *userTagsRoutes) UpdateUserTag(c *gin.Context) {
	userID, ok := GetAuthenticatedUserID(c)
	if !ok {
		return
	}

	userTagID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "invalid user tag ID format",
		})
		return
	}

	var input entities.UpdateUserTagRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	userTag, err := r.userTagsService.UpdateUserTag(c.Request.Context(), userID, userTagID, &input)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "user tag not found" {
			status = http.StatusNotFound
		} else if err.Error() == "access denied" {
			status = http.StatusForbidden
		} else if err.Error() == "tag name already in use" {
			status = http.StatusConflict
		} else if err.Error() == "tag name is required" {
			status = http.StatusBadRequest
		}
		c.JSON(status, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &entities.ApiResponse{
		Success: true,
		Data:    userTag,
	})
}

// @Summary Delete user tag
// @Description Delete a user tag by its ID (soft delete)
// @Tags user-tags
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User Tag ID"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 401 {object} entities.ApiResponse
// @Failure 403 {object} entities.ApiResponse
// @Failure 404 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /user-tags/{id} [delete]
func (r *userTagsRoutes) DeleteUserTag(c *gin.Context) {
	userID, ok := GetAuthenticatedUserID(c)
	if !ok {
		return
	}

	userTagID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "invalid user tag ID format",
		})
		return
	}

	err = r.userTagsService.DeleteUserTag(c.Request.Context(), userID, userTagID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "user tag not found" {
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
		Message: "User tag deleted successfully",
	})
}
