package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vasst-id/vasst-expense-api/internal/entities"
	"github.com/vasst-id/vasst-expense-api/internal/middleware"
	"github.com/vasst-id/vasst-expense-api/internal/services"
)

type bankRoutes struct {
	bankService services.BankService
	auth        *middleware.AuthMiddleware
}

func newBankRoutes(handler *gin.RouterGroup, bankService services.BankService, auth *middleware.AuthMiddleware) {
	r := &bankRoutes{
		bankService: bankService,
		auth:        auth,
	}

	// Bank endpoints
	banks := handler.Group("/banks")
	{
		// Public endpoint - no authentication required
		banks.GET("", r.GetAllBanks)
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
func (r *bankRoutes) GetAllBanks(c *gin.Context) {
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
