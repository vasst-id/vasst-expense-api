package v0

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vasst-id/vasst-expense-api/internal/entities"
	"github.com/vasst-id/vasst-expense-api/internal/middleware"
	"github.com/vasst-id/vasst-expense-api/internal/services"
)

type bankAdminRoutes struct {
	bankService services.BankService
	auth        *middleware.AuthMiddleware
}

func newBankAdminRoutes(handler *gin.RouterGroup, bankService services.BankService, auth *middleware.AuthMiddleware) {
	r := &bankAdminRoutes{
		bankService: bankService,
		auth:        auth,
	}

	// Bank endpoints
	banks := handler.Group("/banks")
	{
		// Protected endpoints - authentication required
		banks.GET("", auth.AuthRequired(), r.GetAllBanks)
		banks.POST("", auth.AuthRequired(), r.CreateBank)
		banks.GET("/:id", auth.AuthRequired(), r.GetBankByID)
		banks.PUT("/:id", auth.AuthRequired(), r.UpdateBank)
		banks.DELETE("/:id", auth.AuthRequired(), r.DeleteBank)
		banks.GET("/code/:code", auth.AuthRequired(), r.GetBankByCode)
	}
}

// @Summary Get all banks
// @Description Get a list of all active banks (public endpoint)
// @Tags banks
// @Accept json
// @Produce json
// @Success 200 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /banks [get]
func (r *bankAdminRoutes) GetAllBanks(c *gin.Context) {
	banks, err := r.bankService.GetAllBanks(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &entities.ApiResponse{
		Success: true,
		Data:    banks,
	})
}

// @Summary Get bank by ID
// @Description Get a bank by their ID
// @Tags banks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Bank ID"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 401 {object} entities.ApiResponse
// @Failure 404 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /banks/{id} [get]
func (r *bankAdminRoutes) GetBankByID(c *gin.Context) {
	bankID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "invalid bank ID format",
		})
		return
	}

	bank, err := r.bankService.GetBankByID(c.Request.Context(), bankID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "bank not found" {
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
		Data:    bank,
	})
}

// @Summary Create a new bank
// @Description Create a new bank with the provided details
// @Tags banks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body entities.CreateBankInput true "Bank details"
// @Success 201 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 401 {object} entities.ApiResponse
// @Failure 409 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /banks [post]
func (r *bankAdminRoutes) CreateBank(c *gin.Context) {
	var input entities.CreateBankInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	bank, err := r.bankService.CreateBank(c.Request.Context(), &input)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "bank with this code already exists" {
			status = http.StatusConflict
		} else if err.Error() == "bank name is required" ||
			err.Error() == "bank code is required" {
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
		Data:    bank,
	})
}

// @Summary Update a bank
// @Description Update an existing bank's details
// @Tags banks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Bank ID"
// @Param input body entities.UpdateBankInput true "Updated bank details"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 401 {object} entities.ApiResponse
// @Failure 404 {object} entities.ApiResponse
// @Failure 409 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /banks/{id} [put]
func (r *bankAdminRoutes) UpdateBank(c *gin.Context) {
	bankID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "invalid bank ID format",
		})
		return
	}

	var input entities.UpdateBankInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	bank, err := r.bankService.UpdateBank(c.Request.Context(), bankID, &input)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "bank not found" {
			status = http.StatusNotFound
		} else if err.Error() == "bank code already in use" {
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
		Data:    bank,
	})
}

// @Summary Delete a bank
// @Description Delete a bank by their ID
// @Tags banks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Bank ID"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 401 {object} entities.ApiResponse
// @Failure 404 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /banks/{id} [delete]
func (r *bankAdminRoutes) DeleteBank(c *gin.Context) {
	bankID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "invalid bank ID format",
		})
		return
	}

	err = r.bankService.DeleteBank(c.Request.Context(), bankID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "bank not found" {
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
		Message: "Bank deleted successfully",
	})
}

// @Summary Get bank by code
// @Description Get a bank by their bank code
// @Tags banks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param code path string true "Bank Code"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 401 {object} entities.ApiResponse
// @Failure 404 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /banks/code/{code} [get]
func (r *bankAdminRoutes) GetBankByCode(c *gin.Context) {
	bankCode := c.Param("code")
	if bankCode == "" {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "bank code is required",
		})
		return
	}

	bank, err := r.bankService.GetBankByCode(c.Request.Context(), bankCode)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "bank not found" {
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
		Data:    bank,
	})
}
