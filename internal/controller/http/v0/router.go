package v0

import (
	"github.com/gin-gonic/gin"

	"github.com/vasst-id/vasst-expense-api/config"
	"github.com/vasst-id/vasst-expense-api/internal/middleware"
	"github.com/vasst-id/vasst-expense-api/internal/services"
	"github.com/vasst-id/vasst-expense-api/internal/utils"
)

type Services struct {
	Cfg *config.Config

	UserService             services.UserService
	BankService             services.BankService
	CurrencyService         services.CurrencyService
	SubscriptionPlanService services.SubscriptionPlanService
	TaxonomyService         services.TaxonomyService

	// ConversationService services.ConversationService
	// MessageService      services.MessageService
	// OpenAIService       services.OpenAIService
	// GeminiService       services.GeminiService
	// WhatsAppService     services.WhatsAppService

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

	// API Routers
	h := handler.Group("v0")
	{
		newUserAdminRoutes(h, s.UserService, s.AuthMiddleware)                         // User management routes
		newBankAdminRoutes(h, s.BankService, s.AuthMiddleware)                         // Bank management routes
		newCurrencyAdminRoutes(h, s.CurrencyService, s.AuthMiddleware)                 // Currency management routes
		newSubscriptionPlanAdminRoutes(h, s.SubscriptionPlanService, s.AuthMiddleware) // Subscription plan management routes
		newTaxonomyAdminRoutes(h, s.TaxonomyService, s.AuthMiddleware)                 // Taxonomy management routes
	}
}
