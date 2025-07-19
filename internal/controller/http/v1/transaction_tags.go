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

type transactionTagsRoutes struct {
	transactionTagsService services.TransactionTagsService
	auth                   *middleware.AuthMiddleware
}

func newTransactionTagsRoutes(handler *gin.RouterGroup, transactionTagsService services.TransactionTagsService, auth *middleware.AuthMiddleware) {
	r := &transactionTagsRoutes{
		transactionTagsService: transactionTagsService,
		auth:                   auth,
	}

	// Transaction tags endpoints - all require authentication
	transactionTags := handler.Group("/transaction-tags")
	transactionTags.Use(auth.AuthRequired())
	{
		transactionTags.POST("", r.CreateTransactionTag)
		transactionTags.DELETE("/:id", r.DeleteTransactionTag)

	}
}

// @Summary Create transaction tag
// @Description Create a new transaction tag
// @Tags transaction-tags
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body entities.CreateTransactionTagRequest true "Transaction tag details"
// @Success 201 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 401 {object} entities.ApiResponse
// @Failure 403 {object} entities.ApiResponse
// @Failure 404 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /transaction-tags [post]
func (r *transactionTagsRoutes) CreateTransactionTag(c *gin.Context) {
	userID, ok := GetAuthenticatedUserID(c)
	if !ok {
		return
	}

	var input entities.CreateTransactionTagRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	transactionTag, err := r.transactionTagsService.CreateTransactionTag(c.Request.Context(), userID, &input)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "user tag not found" {
			status = http.StatusNotFound
		} else if err.Error() == "access denied" {
			status = http.StatusForbidden
		} else if err.Error() == "transaction ID is required" || err.Error() == "user tag ID is required" {
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
		Data:    transactionTag,
	})
}

// @Summary Create multiple transaction tags
// @Description Create multiple transaction tags for a single transaction
// @Tags transaction-tags
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body entities.CreateMultipleTransactionTagsRequest true "Multiple transaction tags details"
// @Success 201 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 401 {object} entities.ApiResponse
// @Failure 403 {object} entities.ApiResponse
// @Failure 404 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /transaction-tags/multiple [post]
func (r *transactionTagsRoutes) CreateMultipleTransactionTags(c *gin.Context) {
	userID, ok := GetAuthenticatedUserID(c)
	if !ok {
		return
	}

	var input entities.CreateMultipleTransactionTagsRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	transactionTags, err := r.transactionTagsService.CreateMultipleTransactionTags(c.Request.Context(), userID, &input)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "user tag not found" {
			status = http.StatusNotFound
		} else if err.Error() == "access denied" {
			status = http.StatusForbidden
		} else if err.Error() == "transaction ID is required" || err.Error() == "at least one user tag ID is required" {
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
		Data:    transactionTags,
	})
}

// @Summary Get transaction tag by ID
// @Description Get a transaction tag by its ID
// @Tags transaction-tags
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Transaction Tag ID"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 401 {object} entities.ApiResponse
// @Failure 403 {object} entities.ApiResponse
// @Failure 404 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /transaction-tags/{id} [get]
func (r *transactionTagsRoutes) GetTransactionTagByID(c *gin.Context) {
	userID, ok := GetAuthenticatedUserID(c)
	if !ok {
		return
	}

	transactionTagID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "invalid transaction tag ID format",
		})
		return
	}

	transactionTag, err := r.transactionTagsService.GetTransactionTagByID(c.Request.Context(), userID, transactionTagID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "transaction tag not found" || err.Error() == "user tag not found" {
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
		Data:    transactionTag,
	})
}

// @Summary Delete transaction tag
// @Description Delete a transaction tag by its ID
// @Tags transaction-tags
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Transaction Tag ID"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 401 {object} entities.ApiResponse
// @Failure 403 {object} entities.ApiResponse
// @Failure 404 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /transaction-tags/{id} [delete]
func (r *transactionTagsRoutes) DeleteTransactionTag(c *gin.Context) {
	userID, ok := GetAuthenticatedUserID(c)
	if !ok {
		return
	}

	transactionTagID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "invalid transaction tag ID format",
		})
		return
	}

	err = r.transactionTagsService.DeleteTransactionTag(c.Request.Context(), userID, transactionTagID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "transaction tag not found" || err.Error() == "user tag not found" {
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
		Message: "Transaction tag deleted successfully",
	})
}

// @Summary Get transaction tags by transaction
// @Description Get all transaction tags for a specific transaction
// @Tags transaction-tags
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param transaction_id path string true "Transaction ID"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 401 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /transaction-tags/transaction/{transaction_id} [get]
func (r *transactionTagsRoutes) GetTransactionTagsByTransaction(c *gin.Context) {
	userID, ok := GetAuthenticatedUserID(c)
	if !ok {
		return
	}

	transactionID, err := uuid.Parse(c.Param("transaction_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "invalid transaction ID format",
		})
		return
	}

	transactionTags, err := r.transactionTagsService.GetTransactionTagsByTransaction(c.Request.Context(), userID, transactionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &entities.ApiResponse{
		Success: true,
		Data:    transactionTags,
	})
}

// @Summary Delete transaction tags by transaction
// @Description Delete all transaction tags for a specific transaction
// @Tags transaction-tags
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param transaction_id path string true "Transaction ID"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 401 {object} entities.ApiResponse
// @Failure 403 {object} entities.ApiResponse
// @Failure 404 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /transaction-tags/transaction/{transaction_id} [delete]
func (r *transactionTagsRoutes) DeleteTransactionTagsByTransaction(c *gin.Context) {
	userID, ok := GetAuthenticatedUserID(c)
	if !ok {
		return
	}

	transactionID, err := uuid.Parse(c.Param("transaction_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "invalid transaction ID format",
		})
		return
	}

	err = r.transactionTagsService.DeleteTransactionTagsByTransaction(c.Request.Context(), userID, transactionID)
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
		Message: "Transaction tags deleted successfully",
	})
}

// @Summary Get transaction tags by user tag
// @Description Get transaction tags for a specific user tag with pagination
// @Tags transaction-tags
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param user_tag_id path string true "User Tag ID"
// @Param limit query int false "Limit for pagination"
// @Param offset query int false "Offset for pagination"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 401 {object} entities.ApiResponse
// @Failure 403 {object} entities.ApiResponse
// @Failure 404 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /transaction-tags/user-tag/{user_tag_id} [get]
func (r *transactionTagsRoutes) GetTransactionTagsByUserTag(c *gin.Context) {
	userID, ok := GetAuthenticatedUserID(c)
	if !ok {
		return
	}

	userTagID, err := uuid.Parse(c.Param("user_tag_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "invalid user tag ID format",
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

	transactionTags, err := r.transactionTagsService.GetTransactionTagsByUserTag(c.Request.Context(), userID, userTagID, limit, offset)
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
		Data:    transactionTags,
	})
}

// @Summary Get transactions by user tag
// @Description Get transactions with their tags for a specific user tag
// @Tags transaction-tags
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param user_tag_id path string true "User Tag ID"
// @Param limit query int false "Limit for pagination"
// @Param offset query int false "Offset for pagination"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 401 {object} entities.ApiResponse
// @Failure 403 {object} entities.ApiResponse
// @Failure 404 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /transaction-tags/user-tag/{user_tag_id}/transactions [get]
func (r *transactionTagsRoutes) GetTransactionsByUserTag(c *gin.Context) {
	userID, ok := GetAuthenticatedUserID(c)
	if !ok {
		return
	}

	userTagID, err := uuid.Parse(c.Param("user_tag_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "invalid user tag ID format",
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

	transactions, err := r.transactionTagsService.GetTransactionsByUserTag(c.Request.Context(), userID, userTagID, limit, offset)
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
		Data:    transactions,
	})
}

// @Summary Get tagged transactions summary
// @Description Get summary of tagged transactions by user
// @Tags transaction-tags
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} entities.ApiResponse
// @Failure 401 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /transaction-tags/summary [get]
func (r *transactionTagsRoutes) GetTaggedTransactionsSummary(c *gin.Context) {
	userID, ok := GetAuthenticatedUserID(c)
	if !ok {
		return
	}

	summary, err := r.transactionTagsService.GetTaggedTransactionsSummary(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &entities.ApiResponse{
		Success: true,
		Data:    summary,
	})
}
