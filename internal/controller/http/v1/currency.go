package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vasst-id/vasst-expense-api/internal/entities"
	"github.com/vasst-id/vasst-expense-api/internal/middleware"
	"github.com/vasst-id/vasst-expense-api/internal/services"
)

type currencyRoutes struct {
	currencyService services.CurrencyService
	auth            *middleware.AuthMiddleware
}

func newCurrencyRoutes(handler *gin.RouterGroup, currencyService services.CurrencyService, auth *middleware.AuthMiddleware) {
	r := &currencyRoutes{
		currencyService: currencyService,
		auth:            auth,
	}

	// Currency endpoints
	currencies := handler.Group("/currencies")
	{
		// Public endpoint - no authentication required
		currencies.GET("", r.GetAllCurrencies)
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
func (r *currencyRoutes) GetAllCurrencies(c *gin.Context) {
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
