package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/vasst-id/vasst-expense-api/config"
	"github.com/vasst-id/vasst-expense-api/internal/entities"
	"github.com/vasst-id/vasst-expense-api/internal/pubsub"
	"github.com/vasst-id/vasst-expense-api/internal/services"
	"github.com/vasst-id/vasst-expense-api/internal/utils"
	logs "github.com/vasst-id/vasst-expense-api/internal/utils/logger"
)

type MessageEventHandler struct {
	pubsubClient        pubsub.Client
	messageService      services.MessageService
	conversationService services.ConversationService
	contactService      services.ContactService
	userService         services.UserService
	whatsappMediaSvc    services.WhatsAppMediaService
	whatsappSvc         services.WhatsAppService
	storageSvc          services.GoogleStorageService
	organizationSvc     services.OrganizationService
	config              *config.Config
	logger              *logs.Logger
}

func NewMessageEventHandler(
	pubsubClient pubsub.Client,
	messageService services.MessageService,
	conversationService services.ConversationService,
	contactService services.ContactService,
	userService services.UserService,
	whatsappMediaSvc services.WhatsAppMediaService,
	whatsappSvc services.WhatsAppService,
	storageSvc services.GoogleStorageService,
	organizationSvc services.OrganizationService,
	config *config.Config,
	logger *logs.Logger,
) *MessageEventHandler {
	return &MessageEventHandler{
		pubsubClient:        pubsubClient,
		messageService:      messageService,
		conversationService: conversationService,
		contactService:      contactService,
		userService:         userService,
		whatsappMediaSvc:    whatsappMediaSvc,
		whatsappSvc:         whatsappSvc,
		storageSvc:          storageSvc,
		organizationSvc:     organizationSvc,
		config:              config,
		logger:              logger,
	}
}

func (h *MessageEventHandler) PublishMessageCreated(ctx context.Context, event *entities.MessageCreatedEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		h.logger.Error().Err(err).Str("event_id", event.EventID.String()).Msg("Failed to marshal message created event")
		return fmt.Errorf("failed to marshal message created event: %w", err)
	}

	if err := h.pubsubClient.Publish(ctx, entities.TopicMessageCreated, data); err != nil {
		h.logger.Error().Err(err).Str("event_id", event.EventID.String()).Msg("Failed to publish message created event")
		return fmt.Errorf("failed to publish message created event: %w", err)
	}

	h.logger.Info().Str("event_id", event.EventID.String()).Str("message_id", event.MessageID.String()).Msg("Message created event published")
	return nil
}

// HandleWebhookReceived processes webhook and creates message
func (h *MessageEventHandler) HandleWebhookReceived(ctx context.Context, data []byte) error {
	var webhookEvent entities.WebhookReceivedEvent
	if err := json.Unmarshal(data, &webhookEvent); err != nil {
		h.logger.Error().Err(err).Msg("Failed to unmarshal webhook received event")
		return fmt.Errorf("failed to unmarshal webhook received event: %w", err)
	}

	h.logger.Info().Str("event_id", webhookEvent.EventID.String()).Str("platform", webhookEvent.Platform).Msg("Processing webhook received event")

	// Process all messages from the webhook event
	for _, webhookMessage := range webhookEvent.Messages {

		fmt.Printf("DEBUG MESSAGE_HANDLER: Processing webhookMessage = %+v\n", webhookMessage)
		fmt.Printf("DEBUG MESSAGE_HANDLER: webhookMessage.WhatsAppMessageID = '%s'\n", webhookMessage.WhatsAppMessageID)
		fmt.Printf("DEBUG MESSAGE_HANDLER: webhookMessage.Metadata = %+v\n", webhookMessage.Metadata)

		// Validate message length to prevent abuse
		if h.isMessageTooLong(webhookMessage.Content) {
			h.logger.Warn().
				Str("phone", webhookMessage.PhoneNumber).
				Int("word_count", utils.CountWords(webhookMessage.Content)).
				Int("max_allowed", h.config.MaxMessageWordCount).
				Msg("Message too long, sending abuse response")

			// Send abuse response immediately
			if err := h.sendAbuseResponse(ctx, webhookMessage.PhoneNumber); err != nil {
				h.logger.Error().Err(err).Str("phone", webhookMessage.PhoneNumber).Msg("Failed to send abuse response")
			}
			continue // Skip AI processing for this message
		}

		// Get or create contact
		contact, err := h.getOrCreateContactByPhone(ctx, webhookEvent.OrganizationID, webhookMessage.PhoneNumber)
		if err != nil {
			h.logger.Error().Err(err).Str("phone", webhookMessage.PhoneNumber).Msg("Failed to get or create contact")
			return fmt.Errorf("failed to get or create contact: %w", err)
		}

		// Get or create conversation
		conversation, err := h.getOrCreateConversation(ctx, webhookEvent.OrganizationID, contact.ContactID, webhookEvent.MediumID, webhookMessage)
		if err != nil {
			h.logger.Error().Err(err).Str("contact_id", contact.ContactID.String()).Msg("Failed to get or create conversation")
			return fmt.Errorf("failed to get or create conversation: %w", err)
		}

		// Store media ID for later processing after conversation is created
		var mediaID string
		var attachments []entities.Attachment

		if webhookMessage.MediaID != "" && h.isMediaMessage(webhookMessage.MessageType) {
			mediaID = webhookMessage.MediaID
		}

		// Create message
		messageInput := &entities.CreateMessageInput{
			ConversationID: conversation.ConversationID,
			OrganizationID: webhookEvent.OrganizationID,
			ContactID:      contact.ContactID,
			MediumID:       webhookEvent.MediumID,
			SenderTypeID:   int(entities.SenderTypeCustomer),
			Direction:      string(entities.MessageDirectionIncoming),
			MessageTypeID:  webhookMessage.MessageType,
			Content:        webhookMessage.Content,
			MediaURL:       webhookMessage.MediaURL,
			Attachments:    attachments,
			Status:         int(entities.MessageStatusDelivered),
		}

		// Get a system user for webhook messages
		systemUserID := conversation.UserID

		message, err := h.messageService.CreateMessage(ctx, messageInput, systemUserID)
		if err != nil {
			h.logger.Error().Err(err).Str("conversation_id", conversation.ConversationID.String()).Msg("Failed to create message")
			return fmt.Errorf("failed to create message: %w", err)
		}

		// Process media after message creation if media ID is available
		if mediaID != "" {
			storedMediaURL, attachment, err := h.processMediaMessage(ctx, webhookEvent.OrganizationID, conversation.ConversationID, mediaID, webhookMessage.MessageType)
			if err != nil {
				h.logger.Error().Err(err).Str("media_id", mediaID).Msg("Failed to process media message")
				// Continue without updating media URL - message already created
			} else {
				// Update message with media URL and attachments
				updateInput := &entities.UpdateMessageInput{
					MediaURL: storedMediaURL,
				}
				if attachment != nil {
					updateInput.Attachments = []entities.Attachment{*attachment}
				}

				_, err = h.messageService.UpdateMessage(ctx, message.MessageID, webhookEvent.OrganizationID, updateInput)
				if err != nil {
					h.logger.Error().Err(err).Str("message_id", message.MessageID.String()).Msg("Failed to update message with media")
				} else {
					message.MediaURL = storedMediaURL
					h.logger.Info().Str("message_id", message.MessageID.String()).Str("media_url", storedMediaURL).Msg("Updated message with media")
				}
			}
		}

		whatsappMessageID := ""
		if msgID, ok := webhookMessage.Metadata["whatsapp_message_id"]; ok {
			whatsappMessageID = msgID.(string)
		}

		// Publish message created event
		messageEvent := &entities.MessageCreatedEvent{
			EventID:           uuid.New(),
			MessageID:         message.MessageID,
			ConversationID:    conversation.ConversationID,
			OrganizationID:    webhookEvent.OrganizationID,
			ContactID:         contact.ContactID,
			Content:           message.Content,
			SenderTypeID:      message.SenderTypeID,
			Direction:         message.Direction,
			MessageTypeID:     message.MessageTypeID,
			MediaURL:          message.MediaURL,
			WhatsAppMessageID: whatsappMessageID,
			CreatedAt:         time.Now(),
		}

		if err := h.PublishMessageCreated(ctx, messageEvent); err != nil {
			h.logger.Error().Err(err).Str("message_id", message.MessageID.String()).Msg("Failed to publish message created event")
			return fmt.Errorf("failed to publish message created event: %w", err)
		}

		h.logger.Info().Str("message_id", message.MessageID.String()).Str("phone", webhookMessage.PhoneNumber).Msg("Successfully processed webhook message")
	}

	h.logger.Info().Str("event_id", webhookEvent.EventID.String()).Int("message_count", len(webhookEvent.Messages)).Msg("Successfully processed all webhook messages")
	return nil
}

// getOrCreateContactByPhone gets an existing contact or creates a new one
func (h *MessageEventHandler) getOrCreateContactByPhone(ctx context.Context, organizationID uuid.UUID, phoneNumber string) (*entities.Contact, error) {
	// Try to get existing contact
	contact, err := h.contactService.GetContactByPhoneNumber(ctx, phoneNumber)
	if err == nil {
		return contact, nil
	}

	// Contact doesn't exist, create new one
	createInput := &entities.CreateContactInput{
		OrganizationID: organizationID,
		Name:           phoneNumber, // Use phone number as name initially
		PhoneNumber:    phoneNumber,
	}

	return h.contactService.CreateContact(ctx, createInput)
}

// getOrCreateConversation gets an existing active conversation or creates a new one
func (h *MessageEventHandler) getOrCreateConversation(ctx context.Context, organizationID, contactID uuid.UUID, mediumID int, webhookMessage entities.WebhookMessage) (*entities.Conversation, error) {
	// Get a system user for webhook conversations
	systemUserID, err := h.getSystemUserID(ctx, organizationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get system user: %w", err)
	}

	// Get the

	conversation, err := h.conversationService.GetActiveConversation(ctx, organizationID, systemUserID, contactID, mediumID)
	if err == nil && conversation != nil {
		return conversation, nil
	}

	// Conversation doesn't exist, create new one
	createInput := &entities.CreateConversationInput{
		OrganizationID:      organizationID,
		UserID:              systemUserID,
		ContactID:           contactID,
		MediumID:            mediumID,
		Status:              int(entities.ConversationStatusOpen),
		Priority:            int(entities.ConversationPriorityLow),
		LastMessageByID:     &systemUserID,
		LastMessageByName:   "System",
		LastMessageContent:  webhookMessage.Content,
		LastMessageTypeID:   int(webhookMessage.MessageType),
		LastMessageMediaURL: "",
		AIEnabled:           true,
	}

	return h.conversationService.CreateConversation(ctx, createInput)
}

// getSystemUserID gets a system user for the organization (for webhook operations)
func (h *MessageEventHandler) getSystemUserID(ctx context.Context, organizationID uuid.UUID) (uuid.UUID, error) {
	// Try to get any admin user from the organization to use as system user
	// This is a simple approach - in production you might want a dedicated system user
	users, err := h.userService.ListUsersByOrganization(ctx, organizationID, 1, 0)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to get users for organization: %w", err)
	}

	if len(users) == 0 {
		return uuid.Nil, fmt.Errorf("no users found for organization %s", organizationID.String())
	}

	// Use the first available user as system user for webhook operations
	return users[0].UserID, nil
}

// isMediaMessage checks if message type supports media
func (h *MessageEventHandler) isMediaMessage(messageType int) bool {
	mediaTypes := []int{
		entities.MessageTypeImage,
		entities.MessageTypeVideo,
		entities.MessageTypeAudio,
		entities.MessageTypeDocument,
		entities.MessageTypeSticker,
	}

	for _, mediaType := range mediaTypes {
		if messageType == mediaType {
			return true
		}
	}
	return false
}

// processMediaMessage downloads media from WhatsApp and stores it in Google Cloud Storage
func (h *MessageEventHandler) processMediaMessage(ctx context.Context, organizationID, conversationID uuid.UUID, mediaID string, messageType int) (string, *entities.Attachment, error) {
	// Get organization info for storage bucket
	organization, err := h.organizationSvc.GetOrganizationByID(ctx, organizationID)
	if err != nil {
		return "", nil, fmt.Errorf("failed to get organization: %w", err)
	}

	// Get media URL from WhatsApp
	mediaURL, err := h.whatsappMediaSvc.GetMediaURL(ctx, mediaID)
	if err != nil {
		return "", nil, fmt.Errorf("failed to get media URL: %w", err)
	}

	// Download media content
	mediaData, contentType, err := h.whatsappMediaSvc.DownloadMedia(ctx, mediaURL)
	if err != nil {
		return "", nil, fmt.Errorf("failed to download media: %w", err)
	}

	// Generate file name based on media type and ID
	fileName := h.generateMediaFileName(mediaID, contentType, messageType)

	// Upload to Google Cloud Storage using conversation-specific bucket
	uploadResult, err := h.storageSvc.UploadFileFromBytes(ctx, organization.OrganizationCode, conversationID, fileName, mediaData, contentType)
	if err != nil {
		return "", nil, fmt.Errorf("failed to upload media: %w", err)
	}

	// Create attachment metadata
	attachment := &entities.Attachment{
		ID:       mediaID,
		Type:     h.getAttachmentType(messageType),
		URL:      uploadResult.FileURL,
		Filename: fileName,
		Size:     int64(len(mediaData)),
		MimeType: contentType,
	}

	return uploadResult.FileURL, attachment, nil
}

// generateMediaFileName generates a unique file name for media
func (h *MessageEventHandler) generateMediaFileName(mediaID, contentType string, messageType int) string {
	// Get file extension from content type
	ext := h.getFileExtensionFromContentType(contentType)

	// Generate filename with media type prefix
	typePrefix := h.getMediaTypePrefix(messageType)

	return fmt.Sprintf("%s_%s%s", typePrefix, mediaID, ext)
}

// getFileExtensionFromContentType returns file extension based on content type
func (h *MessageEventHandler) getFileExtensionFromContentType(contentType string) string {
	switch {
	case strings.Contains(contentType, "image/jpeg"):
		return ".jpg"
	case strings.Contains(contentType, "image/png"):
		return ".png"
	case strings.Contains(contentType, "image/gif"):
		return ".gif"
	case strings.Contains(contentType, "image/webp"):
		return ".webp"
	case strings.Contains(contentType, "video/mp4"):
		return ".mp4"
	case strings.Contains(contentType, "video/3gpp"):
		return ".3gp"
	case strings.Contains(contentType, "audio/aac"):
		return ".aac"
	case strings.Contains(contentType, "audio/mp4"):
		return ".m4a"
	case strings.Contains(contentType, "audio/mpeg"):
		return ".mp3"
	case strings.Contains(contentType, "audio/ogg"):
		return ".ogg"
	case strings.Contains(contentType, "application/pdf"):
		return ".pdf"
	case strings.Contains(contentType, "application/msword"):
		return ".doc"
	case strings.Contains(contentType, "application/vnd.openxmlformats-officedocument.wordprocessingml.document"):
		return ".docx"
	case strings.Contains(contentType, "application/vnd.ms-excel"):
		return ".xls"
	case strings.Contains(contentType, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"):
		return ".xlsx"
	default:
		return ""
	}
}

// getMediaTypePrefix returns prefix based on message type
func (h *MessageEventHandler) getMediaTypePrefix(messageType int) string {
	switch messageType {
	case entities.MessageTypeImage:
		return "img"
	case entities.MessageTypeVideo:
		return "vid"
	case entities.MessageTypeAudio:
		return "aud"
	case entities.MessageTypeDocument:
		return "doc"
	case entities.MessageTypeSticker:
		return "stk"
	default:
		return "media"
	}
}

// getAttachmentType returns attachment type string based on message type
func (h *MessageEventHandler) getAttachmentType(messageType int) string {
	switch messageType {
	case entities.MessageTypeImage:
		return "image"
	case entities.MessageTypeVideo:
		return "video"
	case entities.MessageTypeAudio:
		return "audio"
	case entities.MessageTypeDocument:
		return "document"
	case entities.MessageTypeSticker:
		return "sticker"
	default:
		return "unknown"
	}
}

// isMessageTooLong checks if message exceeds configured word limit
func (h *MessageEventHandler) isMessageTooLong(content string) bool {
	if content == "" {
		return false
	}

	wordCount := utils.CountWords(content)
	return wordCount > h.config.MaxMessageWordCount
}

// sendAbuseResponse sends an automated response for messages that are too long
func (h *MessageEventHandler) sendAbuseResponse(ctx context.Context, phoneNumber string) error {
	abuseMessage := "We can't process your message as it's too long. Our team will review and get back to you."

	return h.whatsappSvc.SendTextMessage(ctx, phoneNumber, abuseMessage)
}
