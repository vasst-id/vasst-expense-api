package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/vasst-id/vasst-expense-api/config"
	"github.com/vasst-id/vasst-expense-api/internal/entities"
	"github.com/vasst-id/vasst-expense-api/internal/pubsub"
	"github.com/vasst-id/vasst-expense-api/internal/services"
	logs "github.com/vasst-id/vasst-expense-api/internal/utils/logger"
)

// CachedPrompt represents cached prompt data for an organization
type CachedPrompt struct {
	SystemPrompt string
	Knowledge    string
	FullPrompt   string
	CachedAt     time.Time
}

// CachedContext represents cached context data for a contact conversation
type CachedContext struct {
	OrganizationPrompt  string
	ContactContext      string
	ConversationHistory string
	FullContext         string
	MessageCount        int
	CachedAt            time.Time
}

type AIEventHandler struct {
	pubsubClient        pubsub.Client
	geminiService       services.GeminiService
	messageService      services.MessageService
	contactService      services.ContactService
	organizationService services.OrganizationService
	whatsappService     services.WhatsAppService
	config              *config.Config
	logger              *logs.Logger
	// Cache for organization prompts (orgID -> cached prompt)
	promptCache map[uuid.UUID]*CachedPrompt
	// Cache for full context (contactID:conversationID -> cached context)
	contextCache map[string]*CachedContext
	cacheMutex   sync.RWMutex
	cacheExpiry  time.Duration
}

func NewAIEventHandler(
	pubsubClient pubsub.Client,
	geminiService services.GeminiService,
	messageService services.MessageService,
	contactService services.ContactService,
	organizationService services.OrganizationService,
	whatsappService services.WhatsAppService,
	config *config.Config,
	logger *logs.Logger,
) *AIEventHandler {
	return &AIEventHandler{
		pubsubClient:        pubsubClient,
		geminiService:       geminiService,
		messageService:      messageService,
		contactService:      contactService,
		organizationService: organizationService,
		whatsappService:     whatsappService,
		config:              config,
		logger:              logger,
		promptCache:         make(map[uuid.UUID]*CachedPrompt),
		contextCache:        make(map[string]*CachedContext),
		cacheExpiry:         15 * time.Minute, // Cache expires after 15 minutes
	}
}

func (h *AIEventHandler) PublishAIResponseReceived(ctx context.Context, event *entities.AIResponseEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		h.logger.Error().Err(err).Str("event_id", event.EventID.String()).Msg("Failed to marshal AI response event")
		return fmt.Errorf("failed to marshal AI response event: %w", err)
	}

	if err := h.pubsubClient.Publish(ctx, entities.TopicAIResponseReceived, data); err != nil {
		h.logger.Error().Err(err).Str("event_id", event.EventID.String()).Msg("Failed to publish AI response event")
		return fmt.Errorf("failed to publish AI response event: %w", err)
	}

	h.logger.Info().Str("event_id", event.EventID.String()).Str("message_id", event.MessageID.String()).Msg("AI response event published")
	return nil
}

func (h *AIEventHandler) PublishMessageDelivery(ctx context.Context, event *entities.MessageDeliveryEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		h.logger.Error().Err(err).Str("event_id", event.EventID.String()).Msg("Failed to marshal message delivery event")
		return fmt.Errorf("failed to marshal message delivery event: %w", err)
	}

	if err := h.pubsubClient.Publish(ctx, entities.TopicMessageDelivery, data); err != nil {
		h.logger.Error().Err(err).Str("event_id", event.EventID.String()).Msg("Failed to publish message delivery event")
		return fmt.Errorf("failed to publish message delivery event: %w", err)
	}

	h.logger.Info().Str("event_id", event.EventID.String()).Str("message_id", event.MessageID.String()).Msg("Message delivery event published")
	return nil
}

// HandleMessageCreated processes message and generates AI response
func (h *AIEventHandler) HandleMessageCreated(ctx context.Context, data []byte) error {
	var messageEvent entities.MessageCreatedEvent
	if err := json.Unmarshal(data, &messageEvent); err != nil {
		h.logger.Error().Err(err).Msg("Failed to unmarshal message created event")
		return fmt.Errorf("failed to unmarshal message created event: %w", err)
	}

	h.logger.Info().Str("event_id", messageEvent.EventID.String()).Str("message_id", messageEvent.MessageID.String()).Msg("Processing message for AI response")

	// Only process incoming customer messages
	if messageEvent.Direction != string(entities.MessageDirectionIncoming) || messageEvent.SenderTypeID != int(entities.SenderTypeCustomer) {
		h.logger.Info().Str("message_id", messageEvent.MessageID.String()).Msg("Skipping AI processing for non-customer message")
		return nil
	}

	// Send typing indicator
	if h.config.EnableTypingIndicators && messageEvent.WhatsAppMessageID != "" {
		if err := h.whatsappService.SendTypingIndicator(ctx, messageEvent.WhatsAppMessageID); err != nil {
			h.logger.Warn().Err(err).
				Str("whatsapp_message_id", messageEvent.WhatsAppMessageID).
				Msg("Failed to send typing indicator before AI processing")
			// Continue with AI processing even if typing indicator fails
		} else {
			h.logger.Info().
				Str("whatsapp_message_id", messageEvent.WhatsAppMessageID).
				Msg("Typing indicator sent before AI processing")
		}
	}

	// Prepare full context including organization prompt, knowledge, contact context, and conversation history
	fullPrompt, err := h.prepareContext(ctx, messageEvent.OrganizationID, messageEvent.ContactID, messageEvent.ConversationID, messageEvent.MessageID)
	if err != nil {
		h.logger.Error().Err(err).
			Str("organization_id", messageEvent.OrganizationID.String()).
			Str("contact_id", messageEvent.ContactID.String()).
			Str("conversation_id", messageEvent.ConversationID.String()).
			Msg("Failed to prepare context")
		return fmt.Errorf("failed to prepare context: %w", err)
	}

	// Debug: Log the prepared context (first 500 chars to avoid spam)
	contextPreview := fullPrompt
	if len(contextPreview) > 500 {
		contextPreview = contextPreview[:500] + "... [TRUNCATED]"
	}
	h.logger.Info().
		Str("message_id", messageEvent.MessageID.String()).
		Str("contact_id", messageEvent.ContactID.String()).
		Int("full_context_length", len(fullPrompt)).
		Str("context_preview", contextPreview).
		Msg("🧠 CONTEXT PREPARED - Ready for AI processing")

	// Update contact context after processing (async to not block AI response)
	go func() {
		// Use a separate context with timeout for the update
		updateCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		
		isCustomerMessage := messageEvent.SenderTypeID == int(entities.SenderTypeCustomer)
		if err := h.updateContactContext(updateCtx, messageEvent.ContactID, messageEvent.OrganizationID, messageEvent.ConversationID, messageEvent.Content, isCustomerMessage); err != nil {
			h.logger.Warn().Err(err).
				Str("contact_id", messageEvent.ContactID.String()).
				Msg("Failed to update contact context (non-blocking)")
		}
	}()

	// Process with Gemini AI using organization-specific prompt
	startTime := time.Now()
	var response string

	// Check if this is a media message and use appropriate processing
	if h.isMediaMessage(messageEvent.MessageTypeID) && messageEvent.MediaURL != "" {
		response, err = h.geminiService.ProcessMediaMessage(ctx, messageEvent.Content, messageEvent.MessageTypeID, messageEvent.MediaURL, fullPrompt)
	} else {
		response, err = h.geminiService.ProcessCustomerMessageWithPrompt(ctx, messageEvent.Content, fullPrompt)
	}

	if err != nil {
		h.logger.Error().Err(err).Str("message_id", messageEvent.MessageID.String()).Msg("Failed to process message with Gemini")
		return fmt.Errorf("failed to process message with Gemini: %w", err)
	}
	processingTime := time.Since(startTime).Milliseconds()

	// Publish AI response event
	aiEvent := &entities.AIResponseEvent{
		EventID:           uuid.New(),
		MessageID:         messageEvent.MessageID,
		ConversationID:    messageEvent.ConversationID,
		OrganizationID:    messageEvent.OrganizationID,
		ContactID:         messageEvent.ContactID,
		Response:          response,
		Model:             "gemini-2.5-flash",
		ConfidenceScore:   0.85, // You might want to get this from Gemini response
		ProcessingTime:    processingTime,
		WhatsAppMessageID: messageEvent.WhatsAppMessageID,
		CreatedAt:         time.Now(),
	}

	return h.PublishAIResponseReceived(ctx, aiEvent)
}

// HandleAIResponseReceived processes AI response and creates outgoing message
func (h *AIEventHandler) HandleAIResponseReceived(ctx context.Context, data []byte) error {
	var aiEvent entities.AIResponseEvent
	if err := json.Unmarshal(data, &aiEvent); err != nil {
		h.logger.Error().Err(err).Msg("Failed to unmarshal AI response event")
		return fmt.Errorf("failed to unmarshal AI response event: %w", err)
	}

	h.logger.Info().Str("event_id", aiEvent.EventID.String()).Str("message_id", aiEvent.MessageID.String()).Msg("Processing AI response to create outgoing message")

	// Create outgoing message with AI response
	messageInput := &entities.CreateMessageInput{
		ConversationID:    aiEvent.ConversationID,
		OrganizationID:    aiEvent.OrganizationID,
		ContactID:         aiEvent.ContactID,
		MediumID:          1, // WhatsApp default
		SenderTypeID:      int(entities.SenderTypeAI),
		Direction:         string(entities.MessageDirectionOutgoing),
		MessageTypeID:     entities.MessageTypeText,
		Content:           aiEvent.Response,
		AIGenerated:       true,
		AIConfidenceScore: &aiEvent.ConfidenceScore,
		Status:            int(entities.MessageStatusPending),
	}

	// Use system user ID for AI messages
	systemUserID := uuid.New() // This should be a real system user ID
	message, err := h.messageService.CreateMessage(ctx, messageInput, systemUserID)
	if err != nil {
		h.logger.Error().Err(err).Str("conversation_id", aiEvent.ConversationID.String()).Msg("Failed to create AI response message")
		return fmt.Errorf("failed to create AI response message: %w", err)
	}

	h.logger.Info().Str("message_id", message.MessageID.String()).Str("ai_event_id", aiEvent.EventID.String()).Msg("AI response message created successfully")

	// Trigger message delivery
	deliveryEvent := &entities.MessageDeliveryEvent{
		EventID:           uuid.New(),
		MessageID:         message.MessageID,
		ConversationID:    message.ConversationID,
		OrganizationID:    message.OrganizationID,
		ContactID:         aiEvent.ContactID,
		Medium:            "whatsapp",                // TODO: Get from conversation/message medium
		WhatsAppMessageID: aiEvent.WhatsAppMessageID, // Pass the original WhatsApp message ID for typing indicators
		CreatedAt:         time.Now(),
	}

	if err := h.PublishMessageDelivery(ctx, deliveryEvent); err != nil {
		h.logger.Error().Err(err).Str("message_id", message.MessageID.String()).Msg("Failed to publish message delivery event")
		return fmt.Errorf("failed to publish message delivery event: %w", err)
	}

	h.logger.Info().Str("message_id", message.MessageID.String()).Str("delivery_event_id", deliveryEvent.EventID.String()).Msg("Message delivery event published")
	return nil
}

// getOrganizationPrompt gets cached organization prompt or fetches from database
func (h *AIEventHandler) getOrganizationPrompt(ctx context.Context, orgID uuid.UUID) (string, error) {
	h.cacheMutex.RLock()
	if cached, exists := h.promptCache[orgID]; exists && time.Since(cached.CachedAt) < h.cacheExpiry {
		h.cacheMutex.RUnlock()
		h.logger.Debug().Str("organization_id", orgID.String()).Msg("Using cached organization prompt")
		return cached.FullPrompt, nil
	}
	h.cacheMutex.RUnlock()

	// Cache miss or expired - fetch from database
	h.logger.Debug().Str("organization_id", orgID.String()).Msg("Fetching organization prompt from database")
	return h.refreshOrganizationPrompt(ctx, orgID)
}

// refreshOrganizationPrompt fetches fresh prompt data from database and caches it
func (h *AIEventHandler) refreshOrganizationPrompt(ctx context.Context, orgID uuid.UUID) (string, error) {
	// Get system prompt from organization settings
	var systemPrompt string
	orgSettings, err := h.organizationService.GetSettingByOrgID(ctx, orgID)
	if err != nil || orgSettings == nil || orgSettings.SystemPrompt == nil {
		// Fallback to default prompt from file
		systemPrompt = h.getDefaultSystemPrompt()
		h.logger.Info().Str("organization_id", orgID.String()).Msg("Using default system prompt")
	} else {
		systemPrompt = *orgSettings.SystemPrompt
		h.logger.Info().Str("organization_id", orgID.String()).Msg("Using organization-specific system prompt")
	}

	// Get knowledge base from organization knowledge
	knowledgeList, err := h.organizationService.ListKnowledgeByOrgID(ctx, orgID)
	if err != nil {
		h.logger.Error().Err(err).Str("organization_id", orgID.String()).Msg("Failed to get organization knowledge, using empty knowledge")
		knowledgeList = []*entities.OrganizationKnowledge{}
	}

	// Combine knowledge entries including SourceURL
	var knowledgeBuilder strings.Builder
	knowledgeBuilder.WriteString("# Organization Knowledge Base\n\n")

	for _, kb := range knowledgeList {
		if !kb.IsActive {
			continue // Skip inactive knowledge entries
		}

		// Add title if available
		if kb.Title != nil && *kb.Title != "" {
			knowledgeBuilder.WriteString(fmt.Sprintf("## %s\n\n", *kb.Title))
		}

		// Add content
		knowledgeBuilder.WriteString(kb.Content)
		knowledgeBuilder.WriteString("\n\n")

		// Add source URL if available
		if kb.SourceURL != nil && *kb.SourceURL != "" {
			knowledgeBuilder.WriteString(fmt.Sprintf("Source: %s\n\n", *kb.SourceURL))
		}

		// Add description if available
		if kb.Description != nil && *kb.Description != "" {
			knowledgeBuilder.WriteString(fmt.Sprintf("Description: %s\n\n", *kb.Description))
		}

		knowledgeBuilder.WriteString("---\n\n")
	}

	knowledge := knowledgeBuilder.String()

	// Create full prompt
	fullPrompt := systemPrompt + "\n\n" + knowledge

	// Cache the result
	h.cacheMutex.Lock()
	h.promptCache[orgID] = &CachedPrompt{
		SystemPrompt: systemPrompt,
		Knowledge:    knowledge,
		FullPrompt:   fullPrompt,
		CachedAt:     time.Now(),
	}
	h.cacheMutex.Unlock()

	h.logger.Info().
		Str("organization_id", orgID.String()).
		Int("knowledge_entries", len(knowledgeList)).
		Int("prompt_length", len(fullPrompt)).
		Msg("Organization prompt cached successfully")

	return fullPrompt, nil
}

// getDefaultSystemPrompt loads the default system prompt from file
func (h *AIEventHandler) getDefaultSystemPrompt() string {
	// Load system prompt from docs/tge_agent_prompt.md
	systemPromptPath := filepath.Join("internal", "docs", "tge_agent_prompt.md")
	systemPromptBytes, err := os.ReadFile(systemPromptPath)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to read default system prompt file")
		return "You are a helpful AI assistant. Please assist the customer with their inquiries."
	}

	systemPrompt := string(systemPromptBytes)

	// Load knowledge base from docs/knowledge.md
	knowledgePath := filepath.Join("internal", "docs", "knowledge.md")
	knowledgeBytes, err := os.ReadFile(knowledgePath)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to read default knowledge file")
		return systemPrompt
	}

	knowledge := string(knowledgeBytes)

	// Combine system prompt and knowledge base
	return systemPrompt + "\n\n# Default Knowledge Base\n" + knowledge
}

// InvalidateOrganizationPromptCache invalidates cache for specific organization (useful for updates)
func (h *AIEventHandler) InvalidateOrganizationPromptCache(orgID uuid.UUID) {
	h.cacheMutex.Lock()
	delete(h.promptCache, orgID)
	h.cacheMutex.Unlock()
	h.logger.Info().Str("organization_id", orgID.String()).Msg("Organization prompt cache invalidated")
}

// isMediaMessage checks if message type supports media
func (h *AIEventHandler) isMediaMessage(messageType int) bool {
	mediaTypes := []int{
		entities.MessageTypeImage,
		entities.MessageTypeVideo,
		entities.MessageTypeAudio,
		entities.MessageTypeDocument,
		entities.MessageTypeLocation,
		entities.MessageTypeContact,
		entities.MessageTypeSticker,
	}

	for _, mediaType := range mediaTypes {
		if messageType == mediaType {
			return true
		}
	}
	return false
}

// getConversationHistory fetches and formats recent conversation history as context string
func (h *AIEventHandler) getConversationHistory(ctx context.Context, conversationID uuid.UUID, organizationID uuid.UUID, excludeMessageID uuid.UUID, limit int) (string, error) {
	// Get recent messages from conversation (excluding current message)
	messages, err := h.messageService.ListMessagesByConversation(ctx, conversationID, organizationID, limit*2, 0) // Get more to filter out excluded
	if err != nil {
		h.logger.Error().Err(err).
			Str("conversation_id", conversationID.String()).
			Msg("Failed to fetch conversation history")
		return "", fmt.Errorf("failed to fetch conversation history: %w", err)
	}

	if len(messages) == 0 {
		h.logger.Debug().
			Str("conversation_id", conversationID.String()).
			Msg("No conversation history found")
		return "", nil
	}

	var historyBuilder strings.Builder
	historyBuilder.WriteString("=== CONVERSATION HISTORY ===\n")
	
	messageCount := 0
	// Process messages in reverse order (newest first, but we want recent history)
	for i := len(messages) - 1; i >= 0 && messageCount < limit; i-- {
		message := messages[i]
		
		// Skip the current message being processed
		if message.MessageID == excludeMessageID {
			continue
		}
		
		// Format timestamp
		timestamp := message.CreatedAt.Format("2006-01-02 15:04")
		
		// Determine sender type
		var senderLabel string
		switch entities.SenderType(message.SenderTypeID) {
		case entities.SenderTypeCustomer:
			senderLabel = "Customer"
		case entities.SenderTypeAI:
			senderLabel = "AI"
		case entities.SenderTypeAgent:
			senderLabel = "Agent"
		default:
			senderLabel = "System"
		}
		
		// Format message content (truncate if too long)
		content := message.Content
		if len(content) > 200 {
			content = content[:200] + "..."
		}
		
		// Add to history
		historyBuilder.WriteString(fmt.Sprintf("[%s] %s: \"%s\"\n", timestamp, senderLabel, content))
		messageCount++
	}
	
	historyBuilder.WriteString("=== END HISTORY ===\n")
	
	h.logger.Debug().
		Str("conversation_id", conversationID.String()).
		Int("messages_included", messageCount).
		Msg("Conversation history formatted successfully")
	
	return historyBuilder.String(), nil
}

// formatContactContext extracts and formats contact context from the Contact.Context JSON field
func (h *AIEventHandler) formatContactContext(contact *entities.Contact) string {
	if contact == nil || len(contact.Context) == 0 {
		return ""
	}

	// Parse contact context JSON
	var contextData map[string]interface{}
	if err := json.Unmarshal(contact.Context, &contextData); err != nil {
		h.logger.Warn().Err(err).
			Str("contact_id", contact.ContactID.String()).
			Msg("Failed to parse contact context JSON")
		return ""
	}

	var contextBuilder strings.Builder
	contextBuilder.WriteString("=== CONTACT CONTEXT ===\n")

	// Basic contact information
	if contact.Name != "" {
		contextBuilder.WriteString(fmt.Sprintf("Name: %s\n", contact.Name))
	}
	if contact.Salutation != "" {
		contextBuilder.WriteString(fmt.Sprintf("Salutation: %s\n", contact.Salutation))
	}

	// Customer information from context
	if customerInfo, ok := contextData["customer_info"].(map[string]interface{}); ok {
		if customerType, ok := customerInfo["type"].(string); ok {
			contextBuilder.WriteString(fmt.Sprintf("Customer Type: %s\n", customerType))
		}
		if ordersCount, ok := customerInfo["orders_count"].(float64); ok {
			contextBuilder.WriteString(fmt.Sprintf("Orders Count: %.0f\n", ordersCount))
		}
		if responseStyle, ok := customerInfo["response_style"].(string); ok {
			contextBuilder.WriteString(fmt.Sprintf("Preferred Style: %s\n", responseStyle))
		}
		if favoriteProducts, ok := customerInfo["favorite_product"].([]interface{}); ok && len(favoriteProducts) > 0 {
			contextBuilder.WriteString("Favorite Products: ")
			for i, product := range favoriteProducts {
				if i > 0 {
					contextBuilder.WriteString(", ")
				}
				contextBuilder.WriteString(fmt.Sprintf("%v", product))
			}
			contextBuilder.WriteString("\n")
		}
		if tags, ok := customerInfo["tags"].([]interface{}); ok && len(tags) > 0 {
			contextBuilder.WriteString("Customer Tags: ")
			for i, tag := range tags {
				if i > 0 {
					contextBuilder.WriteString(", ")
				}
				contextBuilder.WriteString(fmt.Sprintf("%v", tag))
			}
			contextBuilder.WriteString("\n")
		}
	}

	// Memory - Important facts
	if memory, ok := contextData["memory"].(map[string]interface{}); ok {
		if importantFacts, ok := memory["important_facts"].([]interface{}); ok && len(importantFacts) > 0 {
			contextBuilder.WriteString("Important Facts:\n")
			for _, fact := range importantFacts {
				contextBuilder.WriteString(fmt.Sprintf("- %v\n", fact))
			}
		}
		if previousIssues, ok := memory["previous_issues"].([]interface{}); ok && len(previousIssues) > 0 {
			contextBuilder.WriteString("Previous Issues:\n")
			for _, issue := range previousIssues {
				contextBuilder.WriteString(fmt.Sprintf("- %v\n", issue))
			}
		}
	}

	// Active context
	if activeContext, ok := contextData["active_context"].(map[string]interface{}); ok {
		if currentTopic, ok := activeContext["current_topic"].(string); ok && currentTopic != "" {
			contextBuilder.WriteString(fmt.Sprintf("Current Topic: %s\n", currentTopic))
		}
		if description, ok := activeContext["description"].(string); ok && description != "" {
			contextBuilder.WriteString(fmt.Sprintf("Context Description: %s\n", description))
		}
	}

	// Session summary
	if sessionSummary, ok := contextData["session_summary"].(map[string]interface{}); ok {
		if summary, ok := sessionSummary["summary"].(string); ok && summary != "" {
			contextBuilder.WriteString(fmt.Sprintf("Session Summary: %s\n", summary))
		}
		if sentiment, ok := sessionSummary["sentiment"].(string); ok && sentiment != "" {
			contextBuilder.WriteString(fmt.Sprintf("Current Sentiment: %s\n", sentiment))
		}
		if needsHuman, ok := sessionSummary["needs_human"].(bool); ok && needsHuman {
			contextBuilder.WriteString("⚠️ ATTENTION: Customer needs human agent intervention\n")
		}
		if lastQuestion, ok := sessionSummary["last_question"].(string); ok && lastQuestion != "" {
			contextBuilder.WriteString(fmt.Sprintf("Last Question: %s\n", lastQuestion))
		}
	}

	contextBuilder.WriteString("=== END CONTACT CONTEXT ===\n")

	h.logger.Debug().
		Str("contact_id", contact.ContactID.String()).
		Msg("Contact context formatted successfully")

	return contextBuilder.String()
}

// prepareContext combines all context sources: organization prompt + knowledge + contact context + conversation history
func (h *AIEventHandler) prepareContext(ctx context.Context, orgID uuid.UUID, contactID uuid.UUID, conversationID uuid.UUID, excludeMessageID uuid.UUID) (string, error) {
	// Create cache key
	cacheKey := fmt.Sprintf("%s:%s", contactID.String(), conversationID.String())
	
	// Check cache first
	h.cacheMutex.RLock()
	if cached, exists := h.contextCache[cacheKey]; exists && time.Since(cached.CachedAt) < h.cacheExpiry {
		// Check if message count is still valid (no new messages since cache)
		messages, err := h.messageService.ListMessagesByConversation(ctx, conversationID, orgID, 1, 0)
		if err == nil && len(messages) == cached.MessageCount {
			h.cacheMutex.RUnlock()
			h.logger.Debug().
				Str("cache_key", cacheKey).
				Int("cached_context_length", len(cached.FullContext)).
				Msg("Using cached context")
			return cached.FullContext, nil
		}
	}
	h.cacheMutex.RUnlock()

	// Cache miss or expired - build context fresh
	var contextBuilder strings.Builder
	
	// 1. Get organization system prompt and knowledge (existing functionality)
	organizationPrompt, err := h.getOrganizationPrompt(ctx, orgID)
	if err != nil {
		h.logger.Error().Err(err).
			Str("organization_id", orgID.String()).
			Msg("Failed to get organization prompt")
		return "", fmt.Errorf("failed to get organization prompt: %w", err)
	}
	
	// Add organization prompt and knowledge
	contextBuilder.WriteString(organizationPrompt)
	contextBuilder.WriteString("\n\n")

	// 2. Get and format contact context
	var contactContext string
	contact, err := h.contactService.GetContactByIDAndOrganization(ctx, contactID, orgID)
	if err != nil {
		h.logger.Warn().Err(err).
			Str("contact_id", contactID.String()).
			Str("organization_id", orgID.String()).
			Msg("Failed to get contact for context - continuing without contact context")
	} else {
		contactContext = h.formatContactContext(contact)
		if contactContext != "" {
			contextBuilder.WriteString(contactContext)
			contextBuilder.WriteString("\n")
		}
	}

	// 3. Get and format conversation history
	var conversationHistory string
	conversationHistory, err = h.getConversationHistory(ctx, conversationID, orgID, excludeMessageID, 10) // Default limit of 10 messages
	if err != nil {
		h.logger.Warn().Err(err).
			Str("conversation_id", conversationID.String()).
			Msg("Failed to get conversation history - continuing without history")
	} else if conversationHistory != "" {
		contextBuilder.WriteString(conversationHistory)
		contextBuilder.WriteString("\n")
	}

	// 4. Add current message section header
	contextBuilder.WriteString("=== CURRENT MESSAGE ===\n")
	contextBuilder.WriteString("Please respond to the customer's message below, taking into account all the context provided above.\n\n")

	fullContext := contextBuilder.String()

	// Get current message count for cache validation
	messages, _ := h.messageService.ListMessagesByConversation(ctx, conversationID, orgID, 100, 0)
	messageCount := len(messages)

	// Cache the result
	h.cacheMutex.Lock()
	h.contextCache[cacheKey] = &CachedContext{
		OrganizationPrompt:  organizationPrompt,
		ContactContext:      contactContext,
		ConversationHistory: conversationHistory,
		FullContext:         fullContext,
		MessageCount:        messageCount,
		CachedAt:            time.Now(),
	}
	h.cacheMutex.Unlock()
	
	h.logger.Debug().
		Str("organization_id", orgID.String()).
		Str("contact_id", contactID.String()).
		Str("conversation_id", conversationID.String()).
		Str("cache_key", cacheKey).
		Int("context_length", len(fullContext)).
		Int("message_count", messageCount).
		Msg("Full context prepared and cached successfully")

	return fullContext, nil
}

// updateContactContext updates the contact's context with session summary after a threshold of messages
func (h *AIEventHandler) updateContactContext(ctx context.Context, contactID uuid.UUID, organizationID uuid.UUID, conversationID uuid.UUID, currentMessage string, isCustomerMessage bool) error {
	h.logger.Debug().
		Str("contact_id", contactID.String()).
		Str("conversation_id", conversationID.String()).
		Bool("is_customer_message", isCustomerMessage).
		Msg("🔄 CONTEXT UPDATE - updateContactContext called")

	// Only update context for customer messages to avoid updating on every AI response
	if !isCustomerMessage {
		h.logger.Debug().
			Str("contact_id", contactID.String()).
			Msg("⏭️ CONTEXT UPDATE - Skipping context update (not a customer message)")
		return nil
	}

	// Count total messages in conversation to determine if we should update
	messages, err := h.messageService.ListMessagesByConversation(ctx, conversationID, organizationID, 100, 0) // Get recent messages
	if err != nil {
		h.logger.Warn().Err(err).
			Str("conversation_id", conversationID.String()).
			Msg("❌ CONTEXT UPDATE - Failed to count messages for context update")
		return nil // Don't fail the main flow
	}

	h.logger.Info().
		Str("contact_id", contactID.String()).
		Str("conversation_id", conversationID.String()).
		Int("message_count", len(messages)).
		Msg("📊 CONTEXT UPDATE - Message count retrieved")

	// Check if we've reached the threshold (5+ messages) for context update
	threshold := 5
	if len(messages) < threshold || len(messages)%threshold != 0 {
		// Only update every N messages to avoid excessive updates
		h.logger.Info().
			Str("conversation_id", conversationID.String()).
			Int("message_count", len(messages)).
			Int("threshold", threshold).
			Bool("is_multiple_of_threshold", len(messages)%threshold == 0).
			Msg("⏸️ CONTEXT UPDATE - Threshold not met for context update")
		return nil
	}

	h.logger.Info().
		Str("contact_id", contactID.String()).
		Str("conversation_id", conversationID.String()).
		Int("message_count", len(messages)).
		Msg("✅ CONTEXT UPDATE - Threshold met, proceeding with context update")

	// Get current contact
	contact, err := h.contactService.GetContactByIDAndOrganization(ctx, contactID, organizationID)
	if err != nil {
		h.logger.Error().Err(err).
			Str("contact_id", contactID.String()).
			Msg("❌ CONTEXT UPDATE - Failed to get contact for context update")
		return fmt.Errorf("failed to get contact for context update: %w", err)
	}

	h.logger.Info().
		Str("contact_id", contactID.String()).
		Str("contact_name", contact.Name).
		Int("existing_context_length", len(contact.Context)).
		Msg("📋 CONTEXT UPDATE - Contact retrieved successfully")

	// Parse existing context or create new one
	var contextData map[string]interface{}
	if len(contact.Context) > 0 {
		if err := json.Unmarshal(contact.Context, &contextData); err != nil {
			h.logger.Warn().Err(err).
				Str("contact_id", contactID.String()).
				Msg("Failed to parse existing contact context, creating new")
			contextData = make(map[string]interface{})
		}
	} else {
		contextData = make(map[string]interface{})
	}

	// Update session summary
	h.updateSessionSummary(contextData, messages, currentMessage)
	
	// Update system fields
	h.updateSystemFields(contextData)

	// Marshal updated context
	updatedContext, err := json.Marshal(contextData)
	if err != nil {
		h.logger.Error().Err(err).
			Str("contact_id", contactID.String()).
			Msg("Failed to marshal updated context")
		return fmt.Errorf("failed to marshal updated context: %w", err)
	}

	// Update contact context in database
	updateInput := &entities.UpdateContactInput{
		Context: updatedContext,
	}

	h.logger.Info().
		Str("contact_id", contactID.String()).
		Int("new_context_length", len(updatedContext)).
		Msg("💾 CONTEXT UPDATE - Attempting to update contact context in database")

	updatedContact, err := h.contactService.UpdateContact(ctx, contactID, updateInput)
	if err != nil {
		h.logger.Error().Err(err).
			Str("contact_id", contactID.String()).
			Msg("❌ CONTEXT UPDATE - Failed to update contact context in database")
		return fmt.Errorf("failed to update contact context: %w", err)
	}

	h.logger.Info().
		Str("contact_id", contactID.String()).
		Int("updated_context_length", len(updatedContact.Context)).
		Msg("✅ CONTEXT UPDATE - Contact context updated successfully in database")

	// Invalidate context cache for this contact/conversation
	h.invalidateContextCache(contactID, conversationID)

	h.logger.Info().
		Str("contact_id", contactID.String()).
		Str("conversation_id", conversationID.String()).
		Int("message_count", len(messages)).
		Msg("🎉 CONTEXT UPDATE - Contact context update completed successfully")

	return nil
}

// updateSessionSummary updates the session summary based on recent conversation
func (h *AIEventHandler) updateSessionSummary(contextData map[string]interface{}, messages []*entities.Message, currentMessage string) {
	// Ensure session_summary exists
	sessionSummary, ok := contextData["session_summary"].(map[string]interface{})
	if !ok {
		sessionSummary = make(map[string]interface{})
		contextData["session_summary"] = sessionSummary
	}

	// Update message count
	sessionSummary["messages_count"] = len(messages)
	sessionSummary["last_message_at"] = time.Now().Format(time.RFC3339)

	// Generate summary from recent messages
	var summaryBuilder strings.Builder
	customerMessageCount := 0
	aiMessageCount := 0
	
	// Analyze last 10 messages for summary
	startIdx := len(messages) - 10
	if startIdx < 0 {
		startIdx = 0
	}

	for i := startIdx; i < len(messages); i++ {
		message := messages[i]
		if message.SenderTypeID == int(entities.SenderTypeCustomer) {
			customerMessageCount++
		} else if message.SenderTypeID == int(entities.SenderTypeAI) {
			aiMessageCount++
		}
	}

	// Create intelligent summary
	if customerMessageCount > 0 {
		summaryBuilder.WriteString(fmt.Sprintf("Customer telah mengirim %d pesan", customerMessageCount))
		if aiMessageCount > 0 {
			summaryBuilder.WriteString(fmt.Sprintf(" dan menerima %d respon dari AI", aiMessageCount))
		}
		
		// Analyze sentiment from most recent customer message
		sentiment := h.analyzeSentiment(currentMessage)
		sessionSummary["sentiment"] = sentiment
		
		// Store current message as last question if it's a question
		if strings.Contains(currentMessage, "?") || 
		   strings.Contains(strings.ToLower(currentMessage), "berapa") ||
		   strings.Contains(strings.ToLower(currentMessage), "kapan") ||
		   strings.Contains(strings.ToLower(currentMessage), "bagaimana") {
			sessionSummary["last_question"] = currentMessage
		}
	}

	sessionSummary["summary"] = summaryBuilder.String()
	
	// Check if human intervention is needed
	needsHuman := strings.Contains(strings.ToLower(currentMessage), "manusia") ||
		strings.Contains(strings.ToLower(currentMessage), "admin") ||
		strings.Contains(strings.ToLower(currentMessage), "manajer") ||
		strings.Contains(strings.ToLower(currentMessage), "komplain")
	sessionSummary["needs_human"] = needsHuman
}

// updateSystemFields updates system tracking fields
func (h *AIEventHandler) updateSystemFields(contextData map[string]interface{}) {
	systemFields, ok := contextData["system_fields"].(map[string]interface{})
	if !ok {
		systemFields = make(map[string]interface{})
		contextData["system_fields"] = systemFields
	}

	now := time.Now()
	systemFields["updated_at"] = now.Format(time.RFC3339)
	systemFields["last_message_at"] = now.Format(time.RFC3339)
	systemFields["version"] = "1.0"
	systemFields["context_health"] = "active"

	// Calculate approximate token count
	if fullContext, err := json.Marshal(contextData); err == nil {
		systemFields["token_count"] = len(string(fullContext)) / 4 // Rough token estimation
	}
}

// analyzeSentiment performs basic sentiment analysis on message content
func (h *AIEventHandler) analyzeSentiment(message string) string {
	lowerMessage := strings.ToLower(message)
	
	// Positive indicators
	positiveWords := []string{"terima kasih", "bagus", "suka", "senang", "puas", "mantap", "oke", "baik"}
	negativeWords := []string{"kecewa", "buruk", "tidak puas", "jelek", "lambat", "lama", "mahal", "marah"}
	
	positiveCount := 0
	negativeCount := 0
	
	for _, word := range positiveWords {
		if strings.Contains(lowerMessage, word) {
			positiveCount++
		}
	}
	
	for _, word := range negativeWords {
		if strings.Contains(lowerMessage, word) {
			negativeCount++
		}
	}
	
	if positiveCount > negativeCount {
		return "positif"
	} else if negativeCount > positiveCount {
		return "negatif"
	}
	
	return "netral"
}

// invalidateContextCache removes cached context for a specific contact/conversation
func (h *AIEventHandler) invalidateContextCache(contactID uuid.UUID, conversationID uuid.UUID) {
	cacheKey := fmt.Sprintf("%s:%s", contactID.String(), conversationID.String())
	h.cacheMutex.Lock()
	delete(h.contextCache, cacheKey)
	h.cacheMutex.Unlock()
	
	h.logger.Debug().
		Str("contact_id", contactID.String()).
		Str("conversation_id", conversationID.String()).
		Str("cache_key", cacheKey).
		Msg("Context cache invalidated")
}

// InvalidateContactCache invalidates all cached contexts for a specific contact (useful when contact data changes)
func (h *AIEventHandler) InvalidateContactCache(contactID uuid.UUID) {
	h.cacheMutex.Lock()
	defer h.cacheMutex.Unlock()
	
	contactIDStr := contactID.String()
	keysToDelete := make([]string, 0)
	
	for key := range h.contextCache {
		if strings.HasPrefix(key, contactIDStr+":") {
			keysToDelete = append(keysToDelete, key)
		}
	}
	
	for _, key := range keysToDelete {
		delete(h.contextCache, key)
	}
	
	h.logger.Debug().
		Str("contact_id", contactID.String()).
		Int("invalidated_keys", len(keysToDelete)).
		Msg("All contact context caches invalidated")
}
