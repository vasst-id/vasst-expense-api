package v0

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vasst-id/vasst-expense-api/internal/entities"
	"github.com/vasst-id/vasst-expense-api/internal/middleware"
	"github.com/vasst-id/vasst-expense-api/internal/services"
)

type subscriptionPlanAdminRoutes struct {
	subscriptionPlanService services.SubscriptionPlanService
	auth                    *middleware.AuthMiddleware
}

func newSubscriptionPlanAdminRoutes(handler *gin.RouterGroup, subscriptionPlanService services.SubscriptionPlanService, auth *middleware.AuthMiddleware) {
	r := &subscriptionPlanAdminRoutes{
		subscriptionPlanService: subscriptionPlanService,
		auth:                    auth,
	}

	// Plan endpoints
	subscriptionPlans := handler.Group("/subscription-plans")
	{
		subscriptionPlans.GET("", auth.AuthRequired(), r.GetAllSubscriptionPlans)
		subscriptionPlans.POST("", auth.AuthRequired(), r.CreateSubscriptionPlan)
		subscriptionPlans.GET("/:id", auth.AuthRequired(), r.GetSubscriptionPlanByID)
		subscriptionPlans.PUT("/:id", auth.AuthRequired(), r.UpdateSubscriptionPlan)
		subscriptionPlans.DELETE("/:id", auth.AuthRequired(), r.DeleteSubscriptionPlan)
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
func (r *subscriptionPlanAdminRoutes) GetAllSubscriptionPlans(c *gin.Context) {
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

// @Summary Get subscription plan by ID
// @Description Get a subscription plan by their ID
// @Tags subscription-plans
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Subscription Plan ID"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 401 {object} entities.ApiResponse
// @Failure 404 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /subscription-plans/{id} [get]
func (r *subscriptionPlanAdminRoutes) GetSubscriptionPlanByID(c *gin.Context) {
	subscriptionPlanID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "invalid subscription plan ID format",
		})
		return
	}

	subscriptionPlan, err := r.subscriptionPlanService.GetSubscriptionPlanByID(c.Request.Context(), subscriptionPlanID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "subscription plan not found" {
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
		Data:    subscriptionPlan,
	})
}

// @Summary Create a new subscription plan
// @Description Create a new subscription plan with the provided details
// @Tags subscription-plans
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body entities.CreateSubscriptionPlanInput true "Subscription plan details"
// @Success 201 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 401 {object} entities.ApiResponse
// @Failure 409 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /subscription-plans [post]
func (r *subscriptionPlanAdminRoutes) CreateSubscriptionPlan(c *gin.Context) {
	var input entities.CreateSubscriptionPlanInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	subscriptionPlan, err := r.subscriptionPlanService.CreateSubscriptionPlan(c.Request.Context(), &input)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "subscription plan name is required" ||
			err.Error() == "subscription plan price is required" ||
			err.Error() == "subscription plan currency ID is required" {
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
		Data:    subscriptionPlan,
	})
}

// @Summary Update a subscription plan
// @Description Update an existing subscription plan's details
// @Tags subscription-plans
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Subscription Plan ID"
// @Param input body entities.UpdateSubscriptionPlanInput true "Updated subscription plan details"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 401 {object} entities.ApiResponse
// @Failure 404 {object} entities.ApiResponse
// @Failure 409 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /subscription-plans/{id} [put]
func (r *subscriptionPlanAdminRoutes) UpdateSubscriptionPlan(c *gin.Context) {
	subscriptionPlanID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "invalid subscription plan ID format",
		})
		return
	}

	var input entities.UpdateSubscriptionPlanInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	subscriptionPlan, err := r.subscriptionPlanService.UpdateSubscriptionPlan(c.Request.Context(), subscriptionPlanID, &input)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "subscription plan not found" {
			status = http.StatusNotFound
		} else if err.Error() == "subscription plan name already in use" {
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
		Data:    subscriptionPlan,
	})
}

// @Summary Delete a subscription plan
// @Description Delete a subscription plan by their ID
// @Tags subscription-plans
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Subscription Plan ID"
// @Success 200 {object} entities.ApiResponse
// @Failure 400 {object} entities.ApiResponse
// @Failure 401 {object} entities.ApiResponse
// @Failure 404 {object} entities.ApiResponse
// @Failure 500 {object} entities.ApiResponse
// @Router /subscription-plans/{id} [delete]
func (r *subscriptionPlanAdminRoutes) DeleteSubscriptionPlan(c *gin.Context) {
	subscriptionPlanID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "invalid subscription plan ID format",
		})
		return
	}

	err = r.subscriptionPlanService.DeleteSubscriptionPlan(c.Request.Context(), subscriptionPlanID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "subscription plan not found" {
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
		Message: "Subscription plan deleted successfully",
	})
}
