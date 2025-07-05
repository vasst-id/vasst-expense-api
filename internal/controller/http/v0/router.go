package v0

import (
	"github.com/gin-gonic/gin"

	"github.com/vasst-id/vasst-expense-api/config"
	"github.com/vasst-id/vasst-expense-api/internal/events/handlers"
	"github.com/vasst-id/vasst-expense-api/internal/middleware"
	"github.com/vasst-id/vasst-expense-api/internal/services"
	"github.com/vasst-id/vasst-expense-api/internal/utils"
)

type Services struct {
	Cfg *config.Config

	UserService         services.UserService
	OrganizationService services.OrganizationService
	LeadsService        services.LeadsService

	// Event handlers
	WebhookEventHandler *handlers.WebhookEventHandler

	AuthMiddleware *middleware.AuthMiddleware
}

func (s Services) Initialized() error {
	return utils.ValidateStruct(s)
}

func NewRouter(handler *gin.Engine, s Services) {

	// panic if any of the services field is not initialized
	if err := s.Initialized(); err != nil {
		panic(err)
	}

	// Organization Routers
	h := handler.Group("v0")
	{
		newUserRoutes(h, s.UserService, s.AuthMiddleware)
		newOrganizationRoutes(h, s.OrganizationService, s.AuthMiddleware)
		newLeadsRoutes(h, s.LeadsService, s.AuthMiddleware)

		// Webhook routes (no auth required for webhooks)
		newWebhookRoutes(h, s.WebhookEventHandler, s.OrganizationService)
	}
}
