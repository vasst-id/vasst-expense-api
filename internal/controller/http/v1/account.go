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

type accountRoutes struct {
	accountService services.AccountService
	auth           *middleware.AuthMiddleware
}

func newAccountRoutes(handler *gin.RouterGroup, accountService services.AccountService, auth *middleware.AuthMiddleware) {
	r := &accountRoutes{
		accountService: accountService,
		auth:           auth,
	}

	// Account endpoints - all require authentication
	accounts := handler.Group("/accounts")
	accounts.Use(auth.AuthRequired())
	{
		accounts.GET("", r.GetAccountsByUserID)
		accounts.POST("", r.CreateAccount)
		accounts.GET("/:id", r.GetAccountByID)
		accounts.PUT("/:id", r.UpdateAccount)
		accounts.DELETE("/:id", r.DeleteAccount)
	}
}

// @Summary Get accounts by user
// @Description Get a list of accounts for the authenticated user with pagination
// @Tags accounts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Limit for pagination"
// @Param offset query int false "Offset for pagination"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 401 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /accounts [get]
func (r *accountRoutes) GetAccountsByUserID(c *gin.Context) {
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

	accounts, err := r.accountService.GetAccountsByUserID(c.Request.Context(), userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &entities.ApiResponse{
		Success: true,
		Data:    accounts,
	})
}

// @Summary Get active accounts
// @Description Get all active accounts for the authenticated user
// @Tags accounts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} entities.ApiResponse
// @Failure 401 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /accounts/active [get]
func (r *accountRoutes) GetActiveAccountsByUserID(c *gin.Context) {
	userID, ok := GetAuthenticatedUserID(c)
	if !ok {
		return
	}

	accounts, err := r.accountService.GetActiveAccountsByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &entities.ApiResponse{
		Success: true,
		Data:    accounts,
	})
}

// @Summary Get account by ID
// @Description Get an account by its ID
// @Tags accounts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Account ID"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 401 {object} entities.ApiResponse
// @Failure 403 {object} entities.ApiResponse
// @Failure 404 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /accounts/{id} [get]
func (r *accountRoutes) GetAccountByID(c *gin.Context) {
	userID, ok := GetAuthenticatedUserID(c)
	if !ok {
		return
	}

	accountID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "invalid account ID format",
		})
		return
	}

	account, err := r.accountService.GetAccountByID(c.Request.Context(), userID, accountID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "account not found" {
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
		Data:    account,
	})
}

// @Summary Create a new account
// @Description Create a new account for the authenticated user
// @Tags accounts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body entities.CreateAccountRequest true "Account details"
// @Success 201 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 401 {object} entities.ApiResponse
// @Failure 409 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /accounts [post]
func (r *accountRoutes) CreateAccount(c *gin.Context) {
	userID, ok := GetAuthenticatedUserID(c)
	if !ok {
		return
	}

	var input entities.CreateAccountRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	input.UserID = userID

	account, err := r.accountService.CreateAccount(c.Request.Context(), userID, &input)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "account with this name already exists" {
			status = http.StatusConflict
		} else if err.Error() == "account name is required" ||
			err.Error() == "account type is required" ||
			err.Error() == "currency ID is required" {
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
		Data:    account,
	})
}

// @Summary Update an account
// @Description Update an existing account's details
// @Tags accounts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Account ID"
// @Param input body entities.UpdateAccountRequest true "Updated account details"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 401 {object} entities.ApiResponse
// @Failure 403 {object} entities.ApiResponse
// @Failure 404 {object} entities.ApiResponse
// @Failure 409 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /accounts/{id} [put]
func (r *accountRoutes) UpdateAccount(c *gin.Context) {
	userID, ok := GetAuthenticatedUserID(c)
	if !ok {
		return
	}

	accountID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "invalid account ID format",
		})
		return
	}

	var input entities.UpdateAccountRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	account, err := r.accountService.UpdateAccount(c.Request.Context(), userID, accountID, &input)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "account not found" {
			status = http.StatusNotFound
		} else if err.Error() == "access denied" {
			status = http.StatusForbidden
		} else if err.Error() == "account name already in use" {
			status = http.StatusConflict
		} else if err.Error() == "account name is required" ||
			err.Error() == "account type is required" ||
			err.Error() == "currency ID is required" {
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
		Data:    account,
	})
}

// @Summary Delete an account
// @Description Delete an account by its ID (soft delete)
// @Tags accounts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Account ID"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 401 {object} entities.ApiResponse
// @Failure 403 {object} entities.ApiResponse
// @Failure 404 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /accounts/{id} [delete]
func (r *accountRoutes) DeleteAccount(c *gin.Context) {
	userID, ok := GetAuthenticatedUserID(c)
	if !ok {
		return
	}

	accountID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "invalid account ID format",
		})
		return
	}

	err = r.accountService.DeleteAccount(c.Request.Context(), userID, accountID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "account not found" {
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
		Message: "Account deleted successfully",
	})
}
