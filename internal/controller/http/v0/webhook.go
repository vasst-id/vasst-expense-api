package v0

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/vasst-id/vasst-expense-api/internal/entities"
	"github.com/vasst-id/vasst-expense-api/internal/events/handlers"
	"github.com/vasst-id/vasst-expense-api/internal/services"
)

type webhookRoutes struct {
	webhookHandler      *handlers.WebhookEventHandler
	organizationService services.OrganizationService
}

func newWebhookRoutes(handler *gin.RouterGroup, webhookHandler *handlers.WebhookEventHandler, organizationService services.OrganizationService) {
	r := &webhookRoutes{
		webhookHandler:      webhookHandler,
		organizationService: organizationService,
	}

	webhook := handler.Group("/webhook")
	{
		webhook.POST("/whatsapp", r.handleWebhook(entities.PlatformWhatsApp))
		webhook.POST("/instagram", r.handleWebhook(entities.PlatformInstagram))
		webhook.POST("/facebook", r.handleWebhook(entities.PlatformFacebook))
		webhook.POST("/email", r.handleWebhook(entities.PlatformEmail))

		// WhatsApp verification endpoint (required by WhatsApp)
		webhook.GET("/whatsapp", r.verifyWhatsAppWebhook)
	}
}

// handleWebhook returns a handler function for the specified platform
func (r *webhookRoutes) handleWebhook(platform string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Log incoming webhook request
		fmt.Printf("Received webhook for platform: %s, Content-Length: %d, Content-Type: %s\n", 
			platform, c.Request.ContentLength, c.GetHeader("Content-Type"))

		// Check if request has body content
		if c.Request.ContentLength == 0 {
			fmt.Printf("ERR: Request body is empty for platform: %s\n", platform)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Request body is empty"})
			return
		}

		var payload map[string]interface{}
		if err := c.ShouldBindJSON(&payload); err != nil {
			fmt.Printf("ERR: Failed to parse JSON for platform %s: %v\n", platform, err)
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid JSON payload: %v", err)})
			return
		}

		// Check if payload is empty
		if len(payload) == 0 {
			fmt.Printf("ERR: Empty JSON payload for platform: %s\n", platform)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Empty JSON payload"})
			return
		}

		fmt.Printf("Successfully parsed JSON payload for platform: %s with %d fields\n", platform, len(payload))

		// Extract organization ID based on platform requirements
		organizationID, err := r.extractOrganizationID(c, platform)
		if err != nil {
			fmt.Printf("ERR: Failed to extract organization ID for platform %s: %v\n", platform, err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Organization ID required: " + err.Error()})
			return
		}

		fmt.Printf("Processing webhook for organization: %s, platform: %s\n", organizationID.String(), platform)

		// Handle webhook using unified handler
		if err := r.webhookHandler.HandleWebhook(c.Request.Context(), platform, organizationID, payload); err != nil {
			fmt.Printf("ERR: Failed to process webhook event for platform %s: %v\n", platform, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process webhook: " + err.Error()})
			return
		}

		fmt.Printf("Successfully processed webhook for platform: %s\n", platform)
		c.JSON(http.StatusOK, gin.H{"status": "received"})
	}
}

// @Summary Verify WhatsApp webhook
// @Description Verify WhatsApp webhook subscription (required by WhatsApp)
// @Tags webhooks
// @Produce json
// @Param hub.mode query string true "Verification mode"
// @Param hub.verify_token query string true "Verification token"
// @Param hub.challenge query string true "Challenge string"
// @Param key query string true "Organization key"
// @Success 200 {string} string "Challenge response"
// @Failure 403 {object} map[string]interface{} "Forbidden"
// @Router /webhook/whatsapp [get]
func (r *webhookRoutes) verifyWhatsAppWebhook(c *gin.Context) {
	mode := c.Query("hub.mode")
	token := c.Query("hub.verify_token")
	challenge := c.Query("hub.challenge")
	orgKey := c.Query("key")

	if orgKey == "" {
		c.JSON(http.StatusBadRequest, &entities.ApiResponse{
			Success: false,
			Error:   "Organization key is required",
		})
		return
	}

	// Get Organization by key
	org, err := r.organizationService.GetOrganizationByKey(c.Request.Context(), orgKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Get Organization WhatsApp integration token
	integrationToken, err := r.organizationService.GetIntegrationTokenByOrgIDAndType(c.Request.Context(), org.OrganizationID, "WhatsApp")
	if err != nil {
		c.JSON(http.StatusInternalServerError, &entities.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	if integrationToken == "" {
		c.JSON(http.StatusForbidden, &entities.ApiResponse{
			Success: false,
			Error:   "WhatsApp integration token not found",
		})
		return
	}

	if mode == "subscribe" && strings.TrimSpace(token) == strings.TrimSpace(integrationToken) {
		c.String(http.StatusOK, challenge)
		return
	}

	c.JSON(http.StatusForbidden, &entities.ApiResponse{
		Success: false,
		Error:   "Verification failed",
	})
}

// extractOrganizationID extracts organization ID based on platform-specific requirements
func (r *webhookRoutes) extractOrganizationID(c *gin.Context, platform string) (uuid.UUID, error) {
	switch platform {
	case "whatsapp":
		// WhatsApp uses organization key instead of direct ID
		orgKey := c.Query("key")
		if orgKey == "" {
			return uuid.Nil, fmt.Errorf("organization key is required for WhatsApp webhooks")
		}

		org, err := r.organizationService.GetOrganizationByKey(c.Request.Context(), orgKey)
		if err != nil {
			return uuid.Nil, fmt.Errorf("failed to get organization by key: %w", err)
		}

		return org.OrganizationID, nil

	case "instagram", "facebook", "email":
		// Other platforms use organization_id directly as UUID string
		if orgIDStr := c.Query("organization_id"); orgIDStr != "" {
			orgID, err := uuid.Parse(orgIDStr)
			if err != nil {
				return uuid.Nil, fmt.Errorf("invalid organization_id format: %w", err)
			}
			return orgID, nil
		}
		return uuid.Nil, fmt.Errorf("organization_id is required in query parameters")

	default:
		return uuid.Nil, fmt.Errorf("unsupported platform: %s", platform)
	}
}
