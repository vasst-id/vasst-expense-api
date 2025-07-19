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

type categoryRoutes struct {
	categoryService services.CategoryService
	auth            *middleware.AuthMiddleware
}

func newCategoryRoutes(handler *gin.RouterGroup, categoryService services.CategoryService, auth *middleware.AuthMiddleware) {
	r := &categoryRoutes{
		categoryService: categoryService,
		auth:            auth,
	}

	// System category endpoints - all require authentication
	systemCategories := handler.Group("/system-categories")
	systemCategories.Use(auth.AuthRequired())
	{
		systemCategories.GET("", r.GetSystemCategories)
		systemCategories.POST("", r.CreateSystemCategory)
		systemCategories.GET("/:id", r.GetSystemCategoryByID)
		systemCategories.POST("/:id/add-to-user", r.AddSystemCategoryToUser)
	}

	// User category endpoints - all require authentication
	userCategories := handler.Group("/user-categories")
	userCategories.Use(auth.AuthRequired())
	{
		userCategories.GET("", r.GetUserCategories)
		userCategories.POST("", r.CreateUserCategory)
		userCategories.GET("/active", r.GetActiveUserCategories)
		userCategories.GET("/with-transaction-count", r.GetCategoriesWithTransactionCount)
		userCategories.GET("/:id", r.GetUserCategoryByID)
		userCategories.PUT("/:id", r.UpdateUserCategory)
		userCategories.DELETE("/:id", r.DeleteUserCategory)
	}
}

// @Summary Get system categories
// @Description Get a list of system categories with pagination
// @Tags categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 401 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /system-categories [get]
func (r *categoryRoutes) GetSystemCategories(c *gin.Context) {
	page := 1
	pageSize := 20

	if pageStr := c.Query("page"); pageStr != "" {
		if val, err := strconv.Atoi(pageStr); err == nil {
			page = val
		}
	}

	if pageSizeStr := c.Query("page_size"); pageSizeStr != "" {
		if val, err := strconv.Atoi(pageSizeStr); err == nil {
			pageSize = val
		}
	}

	pagination := &entities.PaginationRequest{
		Page:     page,
		PageSize: pageSize,
	}

	categories, _, err := r.categoryService.GetSystemCategories(c.Request.Context(), pagination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &entities.ApiResponse{
		Success: true,
		Data:    categories,
	})
}

// @Summary Create system category
// @Description Create a new system category
// @Tags categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body entities.CreateCategoryInput true "Category details"
// @Success 201 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 401 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /system-categories [post]
func (r *categoryRoutes) CreateSystemCategory(c *gin.Context) {
	var input entities.CreateCategoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	category, err := r.categoryService.CreateSystemCategory(c.Request.Context(), &input)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "category name is required" {
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
		Data:    category,
	})
}

// @Summary Get system category by ID
// @Description Get a system category by its ID
// @Tags categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Category ID"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 401 {object} entities.ApiResponse
// @Failure 404 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /system-categories/{id} [get]
func (r *categoryRoutes) GetSystemCategoryByID(c *gin.Context) {
	categoryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "invalid category ID format",
		})
		return
	}

	category, err := r.categoryService.GetSystemCategoryByID(c.Request.Context(), categoryID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "category not found" {
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
		Data:    category,
	})
}

// @Summary Add system category to user
// @Description Add a system category to the authenticated user's categories
// @Tags categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Category ID"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 401 {object} entities.ApiResponse
// @Failure 404 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /system-categories/{id}/add-to-user [post]
func (r *categoryRoutes) AddSystemCategoryToUser(c *gin.Context) {
	userID, ok := GetAuthenticatedUserID(c)
	if !ok {
		return
	}

	categoryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "invalid category ID format",
		})
		return
	}

	err = r.categoryService.AddSystemCategoryToUser(c.Request.Context(), userID, categoryID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "category not found" {
			status = http.StatusNotFound
		} else if err.Error() == "user already has this category" {
			status = http.StatusConflict
		}
		c.JSON(status, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &entities.ApiResponse{
		Success: true,
		Message: "Category added to user successfully",
	})
}

// @Summary Get user categories
// @Description Get a list of user categories with pagination and filtering
// @Tags categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Param search query string false "Search term"
// @Param sort_by query string false "Sort by field"
// @Param sort_order query string false "Sort order (asc/desc)"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 401 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /user-categories [get]
func (r *categoryRoutes) GetUserCategories(c *gin.Context) {
	userID, ok := GetAuthenticatedUserID(c)
	if !ok {
		return
	}

	page := 1
	pageSize := 20

	if pageStr := c.Query("page"); pageStr != "" {
		if val, err := strconv.Atoi(pageStr); err == nil {
			page = val
		}
	}

	if pageSizeStr := c.Query("page_size"); pageSizeStr != "" {
		if val, err := strconv.Atoi(pageSizeStr); err == nil {
			pageSize = val
		}
	}

	pagination := &entities.PaginationRequest{
		Page:     page,
		PageSize: pageSize,
	}

	filter := &entities.FilterRequest{
		Search:    c.Query("search"),
		SortBy:    c.Query("sort_by"),
		SortOrder: c.Query("sort_order"),
	}

	categories, _, err := r.categoryService.GetUserCategories(c.Request.Context(), userID, pagination, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &entities.ApiResponse{
		Success: true,
		Data:    categories,
	})
}

// @Summary Create user category
// @Description Create a new user category
// @Tags categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body entities.CreateUserCategoryInput true "Category details"
// @Success 201 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 401 {object} entities.ApiResponse
// @Failure 409 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /user-categories [post]
func (r *categoryRoutes) CreateUserCategory(c *gin.Context) {
	userID, ok := GetAuthenticatedUserID(c)
	if !ok {
		return
	}

	var input entities.CreateUserCategoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	category, err := r.categoryService.CreateUserCategory(c.Request.Context(), userID, &input)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "category with this name already exists" {
			status = http.StatusConflict
		} else if err.Error() == "category name is required" {
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
		Data:    category,
	})
}

// @Summary Get active user categories
// @Description Get all active user categories for the authenticated user
// @Tags categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} entities.ApiResponse
// @Failure 401 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /user-categories/active [get]
func (r *categoryRoutes) GetActiveUserCategories(c *gin.Context) {
	userID, ok := GetAuthenticatedUserID(c)
	if !ok {
		return
	}

	categories, err := r.categoryService.GetActiveUserCategories(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &entities.ApiResponse{
		Success: true,
		Data:    categories,
	})
}

// @Summary Get categories with transaction count
// @Description Get user categories with their transaction counts
// @Tags categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} entities.ApiResponse
// @Failure 401 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /user-categories/with-transaction-count [get]
func (r *categoryRoutes) GetCategoriesWithTransactionCount(c *gin.Context) {
	userID, ok := GetAuthenticatedUserID(c)
	if !ok {
		return
	}

	categories, err := r.categoryService.GetCategoriesWithTransactionCount(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &entities.ApiResponse{
		Success: true,
		Data:    categories,
	})
}

// @Summary Get user category by ID
// @Description Get a user category by its ID
// @Tags categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Category ID"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 401 {object} entities.ApiResponse
// @Failure 403 {object} entities.ApiResponse
// @Failure 404 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /user-categories/{id} [get]
func (r *categoryRoutes) GetUserCategoryByID(c *gin.Context) {
	userID, ok := GetAuthenticatedUserID(c)
	if !ok {
		return
	}

	categoryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "invalid category ID format",
		})
		return
	}

	category, err := r.categoryService.GetUserCategoryByID(c.Request.Context(), userID, categoryID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "user category not found" {
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
		Data:    category,
	})
}

// @Summary Update user category
// @Description Update an existing user category
// @Tags categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Category ID"
// @Param input body entities.UpdateUserCategoryInput true "Updated category details"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 401 {object} entities.ApiResponse
// @Failure 403 {object} entities.ApiResponse
// @Failure 404 {object} entities.ApiResponse
// @Failure 409 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /user-categories/{id} [put]
func (r *categoryRoutes) UpdateUserCategory(c *gin.Context) {
	userID, ok := GetAuthenticatedUserID(c)
	if !ok {
		return
	}

	categoryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "invalid category ID format",
		})
		return
	}

	var input entities.UpdateUserCategoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	category, err := r.categoryService.UpdateUserCategory(c.Request.Context(), userID, categoryID, &input)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "user category not found" {
			status = http.StatusNotFound
		} else if err.Error() == "access denied" {
			status = http.StatusForbidden
		} else if err.Error() == "category name already in use" {
			status = http.StatusConflict
		} else if err.Error() == "category name is required" {
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
		Data:    category,
	})
}

// @Summary Delete user category
// @Description Delete a user category by its ID (soft delete)
// @Tags categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Category ID"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 401 {object} entities.ApiResponse
// @Failure 403 {object} entities.ApiResponse
// @Failure 404 {object} entities.ApiResponse
// @Failure 409 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /user-categories/{id} [delete]
func (r *categoryRoutes) DeleteUserCategory(c *gin.Context) {
	userID, ok := GetAuthenticatedUserID(c)
	if !ok {
		return
	}

	categoryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "invalid category ID format",
		})
		return
	}

	err = r.categoryService.DeleteUserCategory(c.Request.Context(), userID, categoryID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "user category not found" {
			status = http.StatusNotFound
		} else if err.Error() == "access denied" {
			status = http.StatusForbidden
		} else if err.Error() == "cannot delete category with existing transactions" {
			status = http.StatusConflict
		}
		c.JSON(status, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &entities.ApiResponse{
		Success: true,
		Message: "Category deleted successfully",
	})
}
