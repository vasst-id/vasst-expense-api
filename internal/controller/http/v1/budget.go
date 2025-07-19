package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/vasst-id/vasst-expense-api/internal/entities"
	"github.com/vasst-id/vasst-expense-api/internal/middleware"
	"github.com/vasst-id/vasst-expense-api/internal/services"
)

// parseDateOnly parses a date string in YYYY-MM-DD format
func parseDateOnly(dateStr string) (time.Time, error) {
	return time.Parse("2006-01-02", dateStr)
}

// bindBudgetRequest binds the request with custom date parsing for both create and update
func bindBudgetRequest(c *gin.Context, isUpdate bool) (interface{}, error) {
	// Read the raw body first
	body, err := c.GetRawData()
	if err != nil {
		return nil, err
	}

	// If body is empty, return error
	if len(body) == 0 {
		return nil, fmt.Errorf("empty request body")
	}

	var rawData map[string]interface{}
	if err := json.Unmarshal(body, &rawData); err != nil {
		return nil, fmt.Errorf("invalid JSON: %v", err)
	}

	if isUpdate {
		var input entities.UpdateBudgetRequest

		// Parse dates manually
		if periodStartStr, ok := rawData["period_start"].(string); ok {
			if t, err := parseDateOnly(periodStartStr); err == nil {
				input.PeriodStart = t
			} else {
				return nil, fmt.Errorf("invalid period_start date: %v", err)
			}
		} else {
			return nil, fmt.Errorf("period_start is required")
		}

		if periodEndStr, ok := rawData["period_end"].(string); ok {
			if t, err := parseDateOnly(periodEndStr); err == nil {
				input.PeriodEnd = t
			} else {
				return nil, fmt.Errorf("invalid period_end date: %v", err)
			}
		} else {
			return nil, fmt.Errorf("period_end is required")
		}

		// Parse other fields
		if userCategoryIDStr, ok := rawData["user_category_id"].(string); ok {
			if id, err := uuid.Parse(userCategoryIDStr); err == nil {
				input.UserCategoryID = id
			} else {
				return nil, fmt.Errorf("invalid user_category_id: %v", err)
			}
		} else {
			return nil, fmt.Errorf("user_category_id is required")
		}

		if name, ok := rawData["name"].(string); ok {
			input.Name = name
		} else {
			return nil, fmt.Errorf("name is required")
		}

		if budgetedAmount, ok := rawData["budgeted_amount"].(float64); ok {
			input.BudgetedAmount = budgetedAmount
		} else {
			return nil, fmt.Errorf("budgeted_amount is required")
		}

		if periodType, ok := rawData["period_type"].(float64); ok {
			input.PeriodType = int(periodType)
		} else {
			return nil, fmt.Errorf("period_type is required")
		}

		if spentAmount, ok := rawData["spent_amount"].(float64); ok {
			input.SpentAmount = spentAmount
		} else {
			return nil, fmt.Errorf("spent_amount is required")
		}

		if isActive, ok := rawData["is_active"].(bool); ok {
			input.IsActive = isActive
		} else {
			return nil, fmt.Errorf("is_active is required")
		}

		return &input, nil
	} else {
		var input entities.CreateBudgetRequest

		// For create requests, try to get workspace_id from JSON body first, then query parameter as fallback
		var workspaceID uuid.UUID
		var workspaceIDErr error

		// Try to get from JSON body first
		if workspaceIDStr, ok := rawData["workspace_id"].(string); ok {
			workspaceID, workspaceIDErr = uuid.Parse(workspaceIDStr)
			if workspaceIDErr != nil {
				return nil, fmt.Errorf("invalid workspace_id in request body: %v", workspaceIDErr)
			}
		}
		input.WorkspaceID = workspaceID

		// Parse dates manually
		if periodStartStr, ok := rawData["period_start"].(string); ok {
			if t, err := parseDateOnly(periodStartStr); err == nil {
				input.PeriodStart = t
			} else {
				return nil, fmt.Errorf("invalid period_start date: %v", err)
			}
		} else {
			return nil, fmt.Errorf("period_start is required")
		}

		if periodEndStr, ok := rawData["period_end"].(string); ok {
			if t, err := parseDateOnly(periodEndStr); err == nil {
				input.PeriodEnd = t
			} else {
				return nil, fmt.Errorf("invalid period_end date: %v", err)
			}
		} else {
			return nil, fmt.Errorf("period_end is required")
		}

		// Parse other fields
		if userCategoryIDStr, ok := rawData["user_category_id"].(string); ok {
			if id, err := uuid.Parse(userCategoryIDStr); err == nil {
				input.UserCategoryID = id
			} else {
				return nil, fmt.Errorf("invalid user_category_id: %v", err)
			}
		} else {
			return nil, fmt.Errorf("user_category_id is required")
		}

		if name, ok := rawData["name"].(string); ok {
			input.Name = name
		} else {
			return nil, fmt.Errorf("name is required")
		}

		if budgetedAmount, ok := rawData["budgeted_amount"].(float64); ok {
			input.BudgetedAmount = budgetedAmount
		} else {
			return nil, fmt.Errorf("budgeted_amount is required")
		}

		if periodType, ok := rawData["period_type"].(float64); ok {
			input.PeriodType = int(periodType)
		} else {
			return nil, fmt.Errorf("period_type is required")
		}

		return &input, nil
	}
}

type budgetRoutes struct {
	budgetService services.BudgetService
	auth          *middleware.AuthMiddleware
}

func newBudgetRoutes(handler *gin.RouterGroup, budgetService services.BudgetService, auth *middleware.AuthMiddleware) {
	r := &budgetRoutes{
		budgetService: budgetService,
		auth:          auth,
	}

	// All budget endpoints require authentication
	budgets := handler.Group("/budgets").Use(auth.AuthRequired())
	{
		budgets.GET("", r.GetAllBudgets)
		budgets.POST("", r.CreateBudget)
		budgets.GET("/:id", r.GetBudgetByID)
		budgets.PUT("/:id", r.UpdateBudget)
		budgets.DELETE("/:id", r.DeleteBudget)
	}
}

// @Summary Get all budgets
// @Description Get a list of budgets for the authenticated user's workspace with optional pagination
// @Tags budgets
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param workspace_id query string true "Workspace ID"
// @Param limit query int false "Limit for pagination"
// @Param offset query int false "Offset for pagination"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 401 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /budgets [get]
func (r *budgetRoutes) GetAllBudgets(c *gin.Context) {
	// Get authenticated user ID
	userID, ok := GetAuthenticatedUserID(c)
	if !ok {
		return
	}
	_ = userID // We have the userID for future workspace validation

	// Get workspace ID from query parameter
	workspaceIDStr := c.Query("workspace_id")
	if workspaceIDStr == "" {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "workspace_id is required",
		})
		return
	}

	workspaceID, err := uuid.Parse(workspaceIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "invalid workspace_id format",
		})
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

	budgets, err := r.budgetService.GetAllBudgets(c.Request.Context(), workspaceID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &entities.ApiResponse{
		Success: true,
		Data:    budgets,
	})
}

// @Summary Get budget by ID
// @Description Get a budget by its ID within the authenticated user's workspace
// @Tags budgets
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Budget ID"
// @Param workspace_id query string true "Workspace ID"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 401 {object} entities.ApiResponse
// @Failure 404 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /budgets/{id} [get]
func (r *budgetRoutes) GetBudgetByID(c *gin.Context) {
	// Get authenticated user ID
	userID, ok := GetAuthenticatedUserID(c)
	if !ok {
		return
	}
	_ = userID // We have the userID for future workspace validation

	budgetID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "invalid budget ID format",
		})
		return
	}

	// Get workspace ID from query parameter
	workspaceIDStr := c.Query("workspace_id")
	if workspaceIDStr == "" {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "workspace_id is required",
		})
		return
	}

	workspaceID, err := uuid.Parse(workspaceIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "invalid workspace_id format",
		})
		return
	}

	budget, err := r.budgetService.GetBudgetByID(c.Request.Context(), budgetID, workspaceID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "budget not found" {
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
		Data:    budget,
	})
}

// @Summary Create a new budget
// @Description Create a new budget in the authenticated user's workspace
// @Tags budgets
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param workspace_id query string true "Workspace ID"
// @Param input body entities.CreateBudgetRequest true "Budget details"
// @Success 201 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 401 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /budgets [post]
func (r *budgetRoutes) CreateBudget(c *gin.Context) {
	// Get authenticated user ID
	userID, ok := GetAuthenticatedUserID(c)
	if !ok {
		return
	}

	input, err := bindBudgetRequest(c, false)
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	createInput := input.(*entities.CreateBudgetRequest)
	budget, err := r.budgetService.CreateBudget(c.Request.Context(), createInput.WorkspaceID, userID, createInput)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "budget name is required" ||
			err.Error() == "budgeted amount must be greater than 0" ||
			err.Error() == "invalid period type" ||
			err.Error() == "period start is required" ||
			err.Error() == "period end is required" ||
			err.Error() == "period end must be after period start" ||
			err.Error() == "user category ID is required" {
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
		Data:    budget,
	})
}

// @Summary Update a budget
// @Description Update an existing budget in the authenticated user's workspace
// @Tags budgets
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Budget ID"
// @Param workspace_id query string true "Workspace ID"
// @Param input body entities.UpdateBudgetRequest true "Updated budget details"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 401 {object} entities.ApiResponse
// @Failure 404 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /budgets/{id} [put]
func (r *budgetRoutes) UpdateBudget(c *gin.Context) {
	// Get authenticated user ID
	userID, ok := GetAuthenticatedUserID(c)
	if !ok {
		return
	}
	_ = userID // We have the userID for future workspace validation

	budgetID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "invalid budget ID format",
		})
		return
	}

	// Get workspace ID from query parameter
	workspaceIDStr := c.Query("workspace_id")
	if workspaceIDStr == "" {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "workspace_id is required",
		})
		return
	}

	workspaceID, err := uuid.Parse(workspaceIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "invalid workspace_id format",
		})
		return
	}

	input, err := bindBudgetRequest(c, true)
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	updateInput := input.(*entities.UpdateBudgetRequest)
	budget, err := r.budgetService.UpdateBudget(c.Request.Context(), budgetID, workspaceID, updateInput)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "budget not found" {
			status = http.StatusNotFound
		} else if err.Error() == "budget name is required" ||
			err.Error() == "budgeted amount must be greater than 0" ||
			err.Error() == "invalid period type" ||
			err.Error() == "period start is required" ||
			err.Error() == "period end is required" ||
			err.Error() == "period end must be after period start" ||
			err.Error() == "user category ID is required" {
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
		Data:    budget,
	})
}

// @Summary Delete a budget
// @Description Delete a budget by its ID within the authenticated user's workspace
// @Tags budgets
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Budget ID"
// @Param workspace_id query string true "Workspace ID"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 401 {object} entities.ApiResponse
// @Failure 404 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /budgets/{id} [delete]
func (r *budgetRoutes) DeleteBudget(c *gin.Context) {
	// Get authenticated user ID
	userID, ok := GetAuthenticatedUserID(c)
	if !ok {
		return
	}
	_ = userID // We have the userID for future workspace validation

	budgetID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "invalid budget ID format",
		})
		return
	}

	// Get workspace ID from query parameter
	workspaceIDStr := c.Query("workspace_id")
	if workspaceIDStr == "" {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "workspace_id is required",
		})
		return
	}

	workspaceID, err := uuid.Parse(workspaceIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "invalid workspace_id format",
		})
		return
	}

	err = r.budgetService.DeleteBudget(c.Request.Context(), budgetID, workspaceID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "budget not found" {
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
		Message: "Budget deleted successfully",
	})
}
