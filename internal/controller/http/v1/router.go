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

	UserService         services.UserService
	ContactService      services.ContactService
	OrganizationService services.OrganizationService
	ConversationService services.ConversationService
	MessageService      services.MessageService
	OpenAIService       services.OpenAIService
	GeminiService       services.GeminiService
	WhatsAppService     services.WhatsAppService

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
	h := handler.Group("v1")
	{
		newOrganizationRoutes(h, s.OrganizationService, s.AuthMiddleware) // To get organization
		newUserRoutes(h, s.UserService, s.AuthMiddleware)                 // To get organization user
		newContactRoutes(h, s.ContactService, s.AuthMiddleware)           // To get organization contact
		newConversationRoutes(h, s.ConversationService, s.AuthMiddleware) // To get organization conversation
		newMessageRoutes(h, s.MessageService, s.AuthMiddleware)           // To get organization message
	}
}
