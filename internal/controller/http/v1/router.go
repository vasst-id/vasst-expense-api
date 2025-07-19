package v1

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
	WorkspaceService        services.WorkspaceService
	AccountService          services.AccountService
	CategoryService         services.CategoryService
	BankService             services.BankService
	CurrencyService         services.CurrencyService
	SubscriptionPlanService services.SubscriptionPlanService
	BudgetService           services.BudgetService
	TransactionService      services.TransactionService
	ConversationService     services.ConversationService
	MessageService          services.MessageService
	TaxonomyService         services.TaxonomyService
	UserTagsService         services.UserTagsService
	TransactionTagsService  services.TransactionTagsService
	VerificationCodeService services.VerificationCodeService
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
	h := handler.Group("v1")
	{
		newUserRoutes(h, s.UserService, s.AuthMiddleware)                         // User management routes
		newWorkspaceRoutes(h, s.WorkspaceService, s.AuthMiddleware)               // Workspace management routes
		newAccountRoutes(h, s.AccountService, s.AuthMiddleware)                   // Account management routes
		newCategoryRoutes(h, s.CategoryService, s.AuthMiddleware)                 // Category management routes
		newBankRoutes(h, s.BankService, s.AuthMiddleware)                         // Bank management routes
		newCurrencyRoutes(h, s.CurrencyService, s.AuthMiddleware)                 // Currency management routes
		newSubscriptionPlanRoutes(h, s.SubscriptionPlanService, s.AuthMiddleware) // Subscription plan management routes
		newBudgetRoutes(h, s.BudgetService, s.AuthMiddleware)                     // Budget management routes
		newTransactionRoutes(h, s.TransactionService, s.AuthMiddleware)           // Transaction management routes
		newConversationRoutes(h, s.ConversationService, s.AuthMiddleware)         // Conversation management routes
		newMessageRoutes(h, s.MessageService, s.AuthMiddleware)                   // Message management routes
		newTaxonomyRoutes(h, s.TaxonomyService, s.AuthMiddleware)                 // Taxonomy management routes
		newUserTagsRoutes(h, s.UserTagsService, s.AuthMiddleware)                 // User tags management routes
		newTransactionTagsRoutes(h, s.TransactionTagsService, s.AuthMiddleware)   // Transaction tags management routes
		newVerificationCodeRoutes(h, s.VerificationCodeService, s.AuthMiddleware) // Verification code management routes
	}
}
