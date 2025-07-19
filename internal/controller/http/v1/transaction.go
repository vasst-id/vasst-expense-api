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

// bindTransactionRequest binds the request with custom date parsing for both create and update
func bindTransactionRequest(c *gin.Context, isUpdate bool) (interface{}, error) {
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
		var input entities.UpdateTransactionRequest

		// Parse transaction date manually
		if transactionDateStr, ok := rawData["transaction_date"].(string); ok {
			if t, err := parseDateOnly(transactionDateStr); err == nil {
				input.TransactionDate = t
			} else {
				return nil, fmt.Errorf("invalid transaction_date: %v", err)
			}
		} else {
			return nil, fmt.Errorf("transaction_date is required")
		}

		// Parse other fields
		if description, ok := rawData["description"].(string); ok {
			input.Description = description
		} else {
			return nil, fmt.Errorf("description is required")
		}

		if amount, ok := rawData["amount"].(float64); ok {
			input.Amount = amount
		} else {
			return nil, fmt.Errorf("amount is required")
		}

		if transactionType, ok := rawData["transaction_type"].(float64); ok {
			input.TransactionType = int(transactionType)
		} else {
			return nil, fmt.Errorf("transaction_type is required")
		}

		// if paymentMethod, ok := rawData["payment_method"].(float64); ok {
		// 	input.PaymentMethod = int(paymentMethod)
		// } else {
		// 	return nil, fmt.Errorf("payment_method is required")
		// }

		// Optional fields
		if accountIDStr, ok := rawData["account_id"].(string); ok {
			if id, err := uuid.Parse(accountIDStr); err == nil {
				input.AccountID = &id
			} else {
				return nil, fmt.Errorf("invalid account_id: %v", err)
			}
		}

		if categoryIDStr, ok := rawData["category_id"].(string); ok {
			if id, err := uuid.Parse(categoryIDStr); err == nil {
				input.CategoryID = &id
			} else {
				return nil, fmt.Errorf("invalid category_id: %v", err)
			}
		}

		if merchantName, ok := rawData["merchant_name"].(string); ok {
			input.MerchantName = &merchantName
		}

		if location, ok := rawData["location"].(string); ok {
			input.Location = &location
		}

		if notes, ok := rawData["notes"].(string); ok {
			input.Notes = &notes
		}

		if isRecurring, ok := rawData["is_recurring"].(bool); ok {
			input.IsRecurring = &isRecurring
		}

		if recurrenceInterval, ok := rawData["recurrence_interval"].(float64); ok {
			interval := int(recurrenceInterval)
			input.RecurrenceInterval = &interval
		}

		if recurrenceEndDateStr, ok := rawData["recurrence_end_date"].(string); ok {
			if t, err := parseDateOnly(recurrenceEndDateStr); err == nil {
				input.RecurrenceEndDate = &t
			} else {
				return nil, fmt.Errorf("invalid recurrence_end_date: %v", err)
			}
		}

		return &input, nil
	} else {
		var input entities.CreateTransactionRequest

		// Parse transaction date manually
		if transactionDateStr, ok := rawData["transaction_date"].(string); ok {
			if t, err := parseDateOnly(transactionDateStr); err == nil {
				input.TransactionDate = t
			} else {
				return nil, fmt.Errorf("invalid transaction_date: %v", err)
			}
		} else {
			return nil, fmt.Errorf("transaction_date is required")
		}

		// Parse required fields
		if description, ok := rawData["description"].(string); ok {
			input.Description = description
		} else {
			return nil, fmt.Errorf("description is required")
		}

		if amount, ok := rawData["amount"].(float64); ok {
			input.Amount = amount
		} else {
			return nil, fmt.Errorf("amount is required")
		}

		if transactionType, ok := rawData["transaction_type"].(float64); ok {
			input.TransactionType = int(transactionType)
		} else {
			return nil, fmt.Errorf("transaction_type is required")
		}

		// if paymentMethod, ok := rawData["payment_method"].(float64); ok {
		// 	input.PaymentMethod = int(paymentMethod)
		// } else {
		// 	return nil, fmt.Errorf("payment_method is required")
		// }

		// Parse workspace_id
		if workspaceIDStr, ok := rawData["workspace_id"].(string); ok {
			if id, err := uuid.Parse(workspaceIDStr); err == nil {
				input.WorkspaceID = id
			} else {
				return nil, fmt.Errorf("invalid workspace_id: %v", err)
			}
		} else {
			return nil, fmt.Errorf("workspace_id is required")
		}

		// Parse account_id
		if accountIDStr, ok := rawData["account_id"].(string); ok {
			if id, err := uuid.Parse(accountIDStr); err == nil {
				input.AccountID = id
			} else {
				return nil, fmt.Errorf("invalid account_id: %v", err)
			}
		} else {
			return nil, fmt.Errorf("account_id is required")
		}

		// Optional fields
		if categoryIDStr, ok := rawData["category_id"].(string); ok {
			if id, err := uuid.Parse(categoryIDStr); err == nil {
				input.CategoryID = &id
			} else {
				return nil, fmt.Errorf("invalid category_id: %v", err)
			}
		}

		if merchantName, ok := rawData["merchant_name"].(string); ok {
			input.MerchantName = &merchantName
		}

		if location, ok := rawData["location"].(string); ok {
			input.Location = &location
		}

		if notes, ok := rawData["notes"].(string); ok {
			input.Notes = &notes
		}

		if isRecurring, ok := rawData["is_recurring"].(bool); ok {
			input.IsRecurring = &isRecurring
		}

		if recurrenceInterval, ok := rawData["recurrence_interval"].(float64); ok {
			input.RecurrenceInterval = int(recurrenceInterval)
		}

		if recurrenceEndDateStr, ok := rawData["recurrence_end_date"].(string); ok {
			if t, err := parseDateOnly(recurrenceEndDateStr); err == nil {
				input.RecurrenceEndDate = &t
			} else {
				return nil, fmt.Errorf("invalid recurrence_end_date: %v", err)
			}
		}

		return &input, nil
	}
}

type transactionRoutes struct {
	transactionService services.TransactionService
	auth               *middleware.AuthMiddleware
}

func newTransactionRoutes(handler *gin.RouterGroup, transactionService services.TransactionService, auth *middleware.AuthMiddleware) {
	r := &transactionRoutes{
		transactionService: transactionService,
		auth:               auth,
	}

	// Transaction endpoints - all require authentication
	transactions := handler.Group("/transactions")
	transactions.Use(auth.AuthRequired())
	{
		transactions.GET("", r.GetTransactionsByWorkspace)
		transactions.POST("", r.CreateTransaction)
		transactions.GET("/:id", r.GetTransactionByID)
		transactions.PUT("/:id", r.UpdateTransaction)
		transactions.DELETE("/:id", r.DeleteTransaction)
	}
}

// @Summary Get transactions by workspace
// @Description Get a list of transactions for a workspace with filtering and pagination
// @Tags transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param workspace_id query string true "Workspace ID"
// @Param account_id query string false "Filter by account ID"
// @Param category_id query string false "Filter by category ID"
// @Param start_date query string false "Start date filter (YYYY-MM-DD)"
// @Param end_date query string false "End date filter (YYYY-MM-DD)"
// @Param payment_method query int false "Filter by payment method"
// @Param description query string false "Filter by description (partial match)"
// @Param merchant_name query string false "Filter by merchant name (partial match)"
// @Param amount query number false "Filter by exact amount"
// @Param is_recurring query boolean false "Filter by recurring status"
// @Param credit_status query int false "Filter by credit status"
// @Param limit query int false "Limit for pagination (default: 10)"
// @Param offset query int false "Offset for pagination (default: 0)"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 401 {object} entities.ApiResponse
// @Failure 403 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /transactions [get]
func (r *transactionRoutes) GetTransactionsByWorkspace(c *gin.Context) {
	userID, ok := GetAuthenticatedUserID(c)
	if !ok {
		return
	}

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

	// Parse filter parameters
	params := &entities.TransactionListParams{}

	if accountIDStr := c.Query("account_id"); accountIDStr != "" {
		if accountID, err := uuid.Parse(accountIDStr); err == nil {
			params.AccountID = &accountID
		}
	}

	if categoryIDStr := c.Query("category_id"); categoryIDStr != "" {
		if categoryID, err := uuid.Parse(categoryIDStr); err == nil {
			params.CategoryID = &categoryID
		}
	}

	if startDateStr := c.Query("start_date"); startDateStr != "" {
		if startDate, err := time.Parse("2006-01-02", startDateStr); err == nil {
			params.StartDate = &startDate
		}
	}

	if endDateStr := c.Query("end_date"); endDateStr != "" {
		if endDate, err := time.Parse("2006-01-02", endDateStr); err == nil {
			params.EndDate = &endDate
		}
	}

	if paymentMethodStr := c.Query("payment_method"); paymentMethodStr != "" {
		if paymentMethod, err := strconv.Atoi(paymentMethodStr); err == nil {
			params.PaymentMethod = &paymentMethod
		}
	}

	if description := c.Query("description"); description != "" {
		params.Description = &description
	}

	if merchantName := c.Query("merchant_name"); merchantName != "" {
		params.MerchantName = &merchantName
	}

	if amountStr := c.Query("amount"); amountStr != "" {
		if amount, err := strconv.ParseFloat(amountStr, 64); err == nil {
			params.Amount = &amount
		}
	}

	if isRecurringStr := c.Query("is_recurring"); isRecurringStr != "" {
		if isRecurring, err := strconv.ParseBool(isRecurringStr); err == nil {
			params.IsRecurring = &isRecurring
		}
	}

	if creditStatusStr := c.Query("credit_status"); creditStatusStr != "" {
		if creditStatus, err := strconv.Atoi(creditStatusStr); err == nil {
			params.CreditStatus = &creditStatus
		}
	}

	transactions, totalCount, err := r.transactionService.GetTransactionsByWorkspace(c.Request.Context(), userID, workspaceID, params, limit, offset)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "workspace not found" {
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
			"transactions": transactions,
			"total":        totalCount,
			"limit":        limit,
			"offset":       offset,
		},
	})
}

// @Summary Create a new transaction
// @Description Create a new transaction
// @Tags transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body entities.CreateTransactionRequest true "Transaction details"
// @Success 201 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 401 {object} entities.ApiResponse
// @Failure 403 {object} entities.ApiResponse
// @Failure 404 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /transactions [post]
func (r *transactionRoutes) CreateTransaction(c *gin.Context) {
	userID, ok := GetAuthenticatedUserID(c)
	if !ok {
		return
	}

	input, err := bindTransactionRequest(c, false)
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	createInput := input.(*entities.CreateTransactionRequest)
	transaction, err := r.transactionService.CreateTransaction(c.Request.Context(), userID, createInput)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "description is required" ||
			err.Error() == "amount is required" ||
			err.Error() == "transaction type is required" ||
			err.Error() == "payment method is required" ||
			err.Error() == "transaction date is required" {
			status = http.StatusBadRequest
		} else if err.Error() == "workspace not found" ||
			err.Error() == "account not found" {
			status = http.StatusNotFound
		} else if err.Error() == "access denied to workspace" ||
			err.Error() == "access denied to account" {
			status = http.StatusForbidden
		}
		c.JSON(status, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, &entities.ApiResponse{
		Success: true,
		Data:    transaction,
	})
}

// @Summary Get transaction by ID
// @Description Get a transaction by its ID
// @Tags transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Transaction ID"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 401 {object} entities.ApiResponse
// @Failure 403 {object} entities.ApiResponse
// @Failure 404 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /transactions/{id} [get]
func (r *transactionRoutes) GetTransactionByID(c *gin.Context) {
	userID, ok := GetAuthenticatedUserID(c)
	if !ok {
		return
	}

	transactionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "invalid transaction ID format",
		})
		return
	}

	transaction, err := r.transactionService.GetTransactionByID(c.Request.Context(), userID, transactionID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "transaction not found" {
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
		Data:    transaction,
	})
}

// @Summary Update a transaction
// @Description Update an existing transaction's details
// @Tags transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Transaction ID"
// @Param input body entities.UpdateTransactionRequest true "Updated transaction details"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 401 {object} entities.ApiResponse
// @Failure 403 {object} entities.ApiResponse
// @Failure 404 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /transactions/{id} [put]
func (r *transactionRoutes) UpdateTransaction(c *gin.Context) {
	userID, ok := GetAuthenticatedUserID(c)
	if !ok {
		return
	}

	transactionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "invalid transaction ID format",
		})
		return
	}

	input, err := bindTransactionRequest(c, true)
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	updateInput := input.(*entities.UpdateTransactionRequest)
	transaction, err := r.transactionService.UpdateTransaction(c.Request.Context(), userID, transactionID, updateInput)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "description is required" ||
			err.Error() == "amount is required" ||
			err.Error() == "transaction type is required" ||
			err.Error() == "payment method is required" ||
			err.Error() == "transaction date is required" {
			status = http.StatusBadRequest
		} else if err.Error() == "transaction not found" ||
			err.Error() == "account not found" {
			status = http.StatusNotFound
		} else if err.Error() == "access denied" ||
			err.Error() == "access denied to account" {
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
		Data:    transaction,
	})
}

// @Summary Delete a transaction
// @Description Delete a transaction by its ID
// @Tags transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Transaction ID"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 401 {object} entities.ApiResponse
// @Failure 403 {object} entities.ApiResponse
// @Failure 404 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /transactions/{id} [delete]
func (r *transactionRoutes) DeleteTransaction(c *gin.Context) {
	userID, ok := GetAuthenticatedUserID(c)
	if !ok {
		return
	}

	transactionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "invalid transaction ID format",
		})
		return
	}

	err = r.transactionService.DeleteTransaction(c.Request.Context(), userID, transactionID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "transaction not found" {
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
		Message: "Transaction deleted successfully",
	})
}
