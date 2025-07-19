package v0

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vasst-id/vasst-expense-api/internal/entities"
	"github.com/vasst-id/vasst-expense-api/internal/middleware"
	"github.com/vasst-id/vasst-expense-api/internal/services"
)

type currencyAdminRoutes struct {
	currencyService services.CurrencyService
	auth            *middleware.AuthMiddleware
}

func newCurrencyAdminRoutes(handler *gin.RouterGroup, currencyService services.CurrencyService, auth *middleware.AuthMiddleware) {
	r := &currencyAdminRoutes{
		currencyService: currencyService,
		auth:            auth,
	}

	// Currency endpoints
	currencies := handler.Group("/currencies")
	{
		currencies.GET("", auth.AuthRequired(), r.GetAllCurrencies)
		currencies.POST("", auth.AuthRequired(), r.CreateCurrency)
		currencies.GET("/:id", auth.AuthRequired(), r.GetCurrencyByID)
		currencies.PUT("/:id", auth.AuthRequired(), r.UpdateCurrency)
		currencies.DELETE("/:id", auth.AuthRequired(), r.DeleteCurrency)
	}
}

// @Summary Get all currencies
// @Description Get a list of all active currencies (public endpoint)
// @Tags currencies
// @Accept json
// @Produce json
// @Success 200 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /currencies [get]
func (r *currencyAdminRoutes) GetAllCurrencies(c *gin.Context) {
	currencies, err := r.currencyService.GetAllCurrencies(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &entities.ApiResponse{
		Success: true,
		Data:    currencies,
	})
}

// @Summary Get currency by ID
// @Description Get a currency by their ID
// @Tags currencies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Currency ID"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 401 {object} entities.ApiResponse
// @Failure 404 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /currencies/{id} [get]
func (r *currencyAdminRoutes) GetCurrencyByID(c *gin.Context) {
	currencyID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "invalid currency ID format",
		})
		return
	}

	currency, err := r.currencyService.GetCurrencyByID(c.Request.Context(), currencyID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "currency not found" {
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
		Data:    currency,
	})
}

// @Summary Create a new currency
// @Description Create a new currency with the provided details
// @Tags currencies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body entities.CreateCurrencyInput true "Currency details"
// @Success 201 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 401 {object} entities.ApiResponse
// @Failure 409 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /currencies [post]
func (r *currencyAdminRoutes) CreateCurrency(c *gin.Context) {
	var input entities.CreateCurrencyInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	currency, err := r.currencyService.CreateCurrency(c.Request.Context(), &input)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "currency with this code already exists" {
			status = http.StatusConflict
		} else if err.Error() == "currency code is required" ||
			err.Error() == "currency name is required" ||
			err.Error() == "currency symbol is required" {
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
		Data:    currency,
	})
}

// @Summary Update a currency
// @Description Update an existing currency's details
// @Tags currencies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Currency ID"
// @Param input body entities.UpdateCurrencyInput true "Updated currency details"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 401 {object} entities.ApiResponse
// @Failure 404 {object} entities.ApiResponse
// @Failure 409 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /currencies/{id} [put]
func (r *currencyAdminRoutes) UpdateCurrency(c *gin.Context) {
	currencyID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "invalid currency ID format",
		})
		return
	}

	var input entities.UpdateCurrencyInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	currency, err := r.currencyService.UpdateCurrency(c.Request.Context(), currencyID, &input)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "currency not found" {
			status = http.StatusNotFound
		} else if err.Error() == "currency code already in use" {
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
		Data:    currency,
	})
}

// @Summary Delete a currency
// @Description Delete a currency by their ID
// @Tags currencies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Currency ID"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 401 {object} entities.ApiResponse
// @Failure 404 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /currencies/{id} [delete]
func (r *currencyAdminRoutes) DeleteCurrency(c *gin.Context) {
	currencyID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "invalid currency ID format",
		})
		return
	}

	err = r.currencyService.DeleteCurrency(c.Request.Context(), currencyID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "currency not found" {
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
		Message: "Currency deleted successfully",
	})
}

// @Summary Get currency by code
// @Description Get a currency by their currency code
// @Tags currencies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param code path string true "Currency Code"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 401 {object} entities.ApiResponse
// @Failure 404 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /currencies/code/{code} [get]
func (r *currencyAdminRoutes) GetCurrencyByCode(c *gin.Context) {
	currencyCode := c.Param("code")
	if currencyCode == "" {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "currency code is required",
		})
		return
	}

	currency, err := r.currencyService.GetCurrencyByCode(c.Request.Context(), currencyCode)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "currency not found" {
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
		Data:    currency,
	})
}
