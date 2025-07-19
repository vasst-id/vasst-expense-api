package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vasst-id/vasst-expense-api/internal/entities"
	"github.com/vasst-id/vasst-expense-api/internal/middleware"
	"github.com/vasst-id/vasst-expense-api/internal/services"
)

type subscriptionPlanRoutes struct {
	subscriptionPlanService services.SubscriptionPlanService
	auth                    *middleware.AuthMiddleware
}

func newSubscriptionPlanRoutes(handler *gin.RouterGroup, subscriptionPlanService services.SubscriptionPlanService, auth *middleware.AuthMiddleware) {
	r := &subscriptionPlanRoutes{
		subscriptionPlanService: subscriptionPlanService,
		auth:                    auth,
	}

	// Subscription Plan endpoints
	subscriptionPlans := handler.Group("/subscription-plans")
	{
		// Public endpoint - no authentication required
		subscriptionPlans.GET("", r.GetAllSubscriptionPlans)
	}
}

// @Summary Get all subscription plans
// @Description Get a list of all active subscription plans (public endpoint)
// @Tags subscription-plans
// @Accept json
// @Produce json
// @Success 200 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /subscription-plans [get]
func (r *subscriptionPlanRoutes) GetAllSubscriptionPlans(c *gin.Context) {
	subscriptionPlans, err := r.subscriptionPlanService.GetAllSubscriptionPlans(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &entities.ApiResponse{
		Success: true,
		Data:    subscriptionPlans,
	})
}
