package services

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/vasst-id/vasst-expense-api/config"
	"github.com/vasst-id/vasst-expense-api/internal/entities"
	"github.com/vasst-id/vasst-expense-api/internal/utils"
	errorsutil "github.com/vasst-id/vasst-expense-api/internal/utils/errors"
)

//go:generate mockgen -source=whatsapp_service.go -package=mock -destination=mock/whatsapp_service_mock.go
type (
	WhatsAppService interface {
		HandleWebhook(ctx context.Context, webhook *entities.WhatsAppWebhook, orgID uuid.UUID) error
		SendTemplateMessage(ctx context.Context, to, templateName, languageCode string) error
		SendTextMessage(ctx context.Context, to, message string) error
		SendTypingIndicator(ctx context.Context, messageId string) error
		SendReadStatus(ctx context.Context, phoneNumber, messageId string) error
		HandleMediaMessage(ctx context.Context, message *entities.WhatsAppMessage, contact *entities.Contact, user *entities.User, orgID uuid.UUID) error
	}

	whatsAppService struct {
		phoneNumberID       string
		accessToken         string
		baseURL             string
		httpClient          *http.Client
		openAIService       OpenAIService
		geminiService       GeminiService
		userService         UserService
		organizationService OrganizationService
		messageService      MessageService
		contactService      ContactService
		storageSvc          GoogleStorageService
	}
)

// NewWhatsAppService creates a new WhatsApp service
func NewWhatsAppService(cfg *config.Config, userService UserService, organizationService OrganizationService, messageService MessageService, contactService ContactService, openAIService OpenAIService, geminiService GeminiService, storageSvc GoogleStorageService) (WhatsAppService, error) {
	phoneNumberID := cfg.WhatsAppPhoneNumberID
	if phoneNumberID == "" {
		return nil, errors.New("WHATSAPP_PHONE_NUMBER_ID environment variable is required")
	}

	accessToken := cfg.WhatsAppAccessToken
	if accessToken == "" {
		return nil, errors.New("WHATSAPP_ACCESS_TOKEN environment variable is required")
	}

	return &whatsAppService{
		phoneNumberID:       phoneNumberID,
		accessToken:         accessToken,
		baseURL:             "https://graph.facebook.com/v22.0",
		httpClient:          &http.Client{},
		openAIService:       openAIService,
		geminiService:       geminiService,
		userService:         userService,
		organizationService: organizationService,
		messageService:      messageService,
		contactService:      contactService,
		storageSvc:          storageSvc,
	}, nil
}

// HandleWebhook processes incoming WhatsApp webhooks
func (s *whatsAppService) HandleWebhook(ctx context.Context, webhook *entities.WhatsAppWebhook, orgID uuid.UUID) error {
	if webhook == nil {
		return errors.New("webhook is required")
	}

	if webhook.Object != "whatsapp_business_account" {
		return errors.New("invalid webhook object type")
	}

	// get organization integration
	integration, err := s.organizationService.GetIntegrationByOrgIDAndType(ctx, orgID, "WhatsApp")
	if err != nil {
		return fmt.Errorf("failed to get organization integration: %w", err)
	}

	if !integration.IsAiEnabled {
		return nil
	}

	fmt.Printf("webhook: %+v\n", webhook)

	for _, entry := range webhook.Entry {
		for _, change := range entry.Changes {
			if change.Field != "messages" {
				continue
			}

			// Process each message in the webhook
			for _, message := range change.Value.Messages {
				phoneNumber := utils.SanitizePhoneNumber(message.From)
				fmt.Printf("phoneNumber: %s\n", phoneNumber)

				// Get or create contact
				contact, err := s.contactService.GetContactByPhoneNumber(ctx, phoneNumber)
				if err != nil {
					// if contact doesn't exist, create a new one
					contact, err = s.contactService.CreateContact(ctx, &entities.CreateContactInput{
						PhoneNumber:    message.From,
						OrganizationID: orgID,
					})
					if err != nil {
						return fmt.Errorf("failed to create contact: %w", err)
					}
				}

				fmt.Printf("contact: %+v\n", contact)

				// Get the default user Id from Organization
				user, err := s.organizationService.GetDefaultUserByOrgID(ctx, orgID)
				fmt.Println("user", user)
				if err != nil {
					return fmt.Errorf("failed to get default user: %w", err)
				}

				// Handle different message types
				if message.Type == "text" && message.Text != nil {
					// Handle text message
					err = s.handleTextMessage(ctx, message, contact, user, orgID, integration)
				} else if message.Type == "image" || message.Type == "video" || message.Type == "audio" || message.Type == "document" || message.Type == "sticker" {
					// Handle media message
					err = s.HandleMediaMessage(ctx, &message, contact, user, orgID)
				} else {
					fmt.Printf("Unsupported message type: %s\n", message.Type)
					continue
				}

				if err != nil {
					return fmt.Errorf("failed to process message: %w", err)
				}
			}
		}
	}

	return nil
}

// SendTemplateMessage sends a template message via WhatsApp
func (s *whatsAppService) SendTemplateMessage(ctx context.Context, to, templateName, languageCode string) error {
	if to == "" {
		return errors.New("recipient phone number is required")
	}

	if templateName == "" {
		return errors.New("template name is required")
	}

	if languageCode == "" {
		languageCode = "en_US" // Default to English
	}

	request := entities.WhatsAppMessageRequest{
		MessagingProduct: "whatsapp",
		To:               to,
		Type:             "template",
		Template: entities.WhatsAppTemplate{
			Name: templateName,
			Language: struct {
				Code string `json:"code"`
			}{
				Code: languageCode,
			},
		},
	}

	return s.sendMessage(ctx, request)
}

// SendTextMessage sends a text message via WhatsApp
func (s *whatsAppService) SendTextMessage(ctx context.Context, to, message string) error {
	if to == "" {
		return errors.New("recipient phone number is required")
	}
	if message == "" {
		return errors.New("message content is required")
	}

	// Build the correct WhatsApp Cloud API text message payload
	request := map[string]interface{}{
		"messaging_product": "whatsapp",
		"recipient_type":    "individual",
		"to":                to,
		"type":              "text",
		"text": map[string]interface{}{
			"preview_url": true,
			"body":        message,
		},
	}

	return s.sendMessage(ctx, request)
}

// sendMessage is a helper function to send messages via WhatsApp API
func (s *whatsAppService) sendMessage(ctx context.Context, request interface{}) error {
	url := fmt.Sprintf("%s/%s/messages", s.baseURL, s.phoneNumberID)
	fmt.Printf("Sending request to URL: %s\n", url)

	jsonData, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	fmt.Printf("Request body: %s\n", string(jsonData))

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.accessToken))
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return errorsutil.New(resp.StatusCode, fmt.Sprintf("failed to send WhatsApp message: %s, response: %s", resp.Status, string(body)))
	}

	return nil
}

// HandleMediaMessage processes media messages from WhatsApp webhooks
func (s *whatsAppService) HandleMediaMessage(ctx context.Context, message *entities.WhatsAppMessage, contact *entities.Contact, user *entities.User, orgID uuid.UUID) error {
	// Determine message type based on media type
	var messageTypeID int
	var mediaURL string
	var content string

	if message.Image != nil {
		messageTypeID = entities.MessageTypeImage
		mediaURL = message.Image.ID
		content = message.Image.Caption
	} else if message.Video != nil {
		messageTypeID = entities.MessageTypeVideo
		mediaURL = message.Video.ID
		content = message.Video.Caption
	} else if message.Audio != nil {
		messageTypeID = entities.MessageTypeAudio
		mediaURL = message.Audio.ID
		content = ""
	} else if message.Document != nil {
		messageTypeID = entities.MessageTypeDocument
		mediaURL = message.Document.ID
		content = message.Document.Caption
	} else if message.Sticker != nil {
		messageTypeID = entities.MessageTypeSticker
		mediaURL = message.Sticker.ID
		content = ""
	} else {
		return errors.New("unsupported media type")
	}

	// Download media from WhatsApp
	mediaBytes, contentType, err := s.downloadMediaFromWhatsApp(ctx, mediaURL)
	if err != nil {
		return fmt.Errorf("failed to download media: %w", err)
	}

	// Create message first to get conversation ID
	messageInput := &entities.CreateMessageInput{
		ConversationID:    uuid.Nil,
		ContactID:         contact.ContactID,
		OrganizationID:    orgID,
		MediumID:          entities.MediumWhatsApp,
		SenderTypeID:      int(entities.SenderTypeCustomer),
		SenderID:          &user.UserID,
		Direction:         "i",
		MessageTypeID:     messageTypeID,
		Content:           content,
		IsBroadcast:       false,
		IsOrderMessage:    false,
		AIGenerated:       false,
		AIConfidenceScore: nil,
	}

	messageResult, err := s.messageService.CreateMessage(ctx, messageInput, user.UserID)
	if err != nil {
		return fmt.Errorf("failed to create message: %w", err)
	}

	// Generate filename for the media
	filename := fmt.Sprintf("whatsapp_%s_%s", mediaURL, time.Now().Format("20060102-150405"))

	// Get organization code
	organization, err := s.organizationService.GetOrganizationByID(ctx, orgID)
	if err != nil {
		return fmt.Errorf("failed to get organization: %w", err)
	}

	if organization == nil {
		return fmt.Errorf("organization not found")
	}

	// Upload media to Google Cloud Storage
	uploadResult, err := s.storageSvc.UploadFileFromBytes(ctx, organization.OrganizationCode, messageResult.ConversationID, filename, mediaBytes, contentType)
	if err != nil {
		return fmt.Errorf("failed to upload media to storage: %w", err)
	}

	// Update message with media URL
	messageResult.MediaURL = uploadResult.FileURL

	// Add file metadata to attachments
	attachment := entities.Attachment{
		ID:       uploadResult.ObjectName,
		Type:     getAttachmentType(uploadResult.ContentType),
		URL:      uploadResult.FileURL,
		Filename: uploadResult.FileName,
		Size:     uploadResult.FileSize,
		MimeType: uploadResult.ContentType,
	}

	// Parse existing attachments or create new array
	var attachments []entities.Attachment
	if messageResult.Attachments != nil {
		if err := json.Unmarshal(messageResult.Attachments, &attachments); err != nil {
			attachments = []entities.Attachment{}
		}
	}

	// Add new attachment
	attachments = append(attachments, attachment)

	// Update the message in the database
	_, err = s.messageService.UpdateMessage(ctx, messageResult.MessageID, orgID, &entities.UpdateMessageInput{
		MediaURL:    uploadResult.FileURL,
		Attachments: attachments,
	})
	if err != nil {
		return fmt.Errorf("failed to update message with media: %w", err)
	}

	return nil
}

// downloadMediaFromWhatsApp downloads media from WhatsApp API
func (s *whatsAppService) downloadMediaFromWhatsApp(ctx context.Context, mediaID string) ([]byte, string, error) {
	// First, get media URL from WhatsApp
	url := fmt.Sprintf("%s/%s", s.baseURL, mediaID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, "", err
	}

	req.Header.Set("Authorization", "Bearer "+s.accessToken)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("failed to get media URL: %s", resp.Status)
	}

	var mediaResponse struct {
		URL string `json:"url"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&mediaResponse); err != nil {
		return nil, "", err
	}

	// Download the actual media file
	mediaReq, err := http.NewRequestWithContext(ctx, "GET", mediaResponse.URL, nil)
	if err != nil {
		return nil, "", err
	}

	mediaReq.Header.Set("Authorization", "Bearer "+s.accessToken)

	mediaResp, err := s.httpClient.Do(mediaReq)
	if err != nil {
		return nil, "", err
	}
	defer mediaResp.Body.Close()

	if mediaResp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("failed to download media: %s", mediaResp.Status)
	}

	mediaBytes, err := io.ReadAll(mediaResp.Body)
	if err != nil {
		return nil, "", err
	}

	contentType := mediaResp.Header.Get("Content-Type")

	return mediaBytes, contentType, nil
}

// handleTextMessage processes text messages from WhatsApp webhooks
func (s *whatsAppService) handleTextMessage(ctx context.Context, message entities.WhatsAppMessage, contact *entities.Contact, user *entities.User, orgID uuid.UUID, integration *entities.OrganizationIntegration) error {
	// Create a new message
	messageResult, err := s.messageService.CreateMessage(ctx, &entities.CreateMessageInput{
		ConversationID:    uuid.Nil,
		ContactID:         contact.ContactID,
		OrganizationID:    orgID,
		MediumID:          entities.MediumWhatsApp,
		SenderTypeID:      int(entities.SenderTypeCustomer),
		SenderID:          &user.UserID,
		Direction:         "i",
		MessageTypeID:     entities.MessageTypeText,
		Content:           message.Text.Body,
		IsBroadcast:       false,
		IsOrderMessage:    false,
		AIGenerated:       false,
		AIConfidenceScore: nil,
	}, user.UserID)
	if err != nil {
		return fmt.Errorf("failed to create message: %w", err)
	}

	// Check if AI is enabled for the organization
	if !integration.IsAiEnabled {
		return nil
	}

	// Process message with selected AI Model (Right now default to Gemini)
	response, err := s.geminiService.ProcessCustomerMessage(ctx, messageResult.Content)
	if err != nil {
		return fmt.Errorf("failed to process message with Gemini: %w", err)
	}

	// Send the response back to the user
	phoneNumber := utils.SanitizePhoneNumber(message.From)
	err = s.SendTextMessage(ctx, phoneNumber, response)
	if err != nil {
		return fmt.Errorf("failed to send response: %w", err)
	}

	return nil
}

// SendTypingIndicator sends a typing indicator to show the bot is typing
func (s *whatsAppService) SendTypingIndicator(ctx context.Context, messageId string) error {
	// WhatsApp typing indicators require a message_id and specific format
	if messageId == "" {
		// If no message ID, we can't send a typing indicator via the official API
		// This is a limitation of WhatsApp Cloud API - typing indicators must reference a specific message
		return nil // Silently skip rather than error
	}

	request := map[string]interface{}{
		"messaging_product": "whatsapp",
		"status":            "read",
		"message_id":        messageId,
		"typing_indicator": map[string]interface{}{
			"type": "text",
		},
	}

	return s.sendMessage(ctx, request)
}

// SendReadStatus marks a message as read
func (s *whatsAppService) SendReadStatus(ctx context.Context, phoneNumber, messageId string) error {
	if messageId == "" {
		return fmt.Errorf("messageId is required for read status")
	}

	request := map[string]interface{}{
		"messaging_product": "whatsapp",
		"status":            "read",
		"message_id":        messageId,
	}

	return s.sendMessage(ctx, request)
}

// getAttachmentType determines the attachment type based on content type
func getAttachmentType(contentType string) string {
	switch {
	case strings.HasPrefix(contentType, "image/"):
		return "image"
	case strings.HasPrefix(contentType, "video/"):
		return "video"
	case strings.HasPrefix(contentType, "audio/"):
		return "audio"
	case strings.HasPrefix(contentType, "application/pdf"):
		return "pdf"
	case strings.HasPrefix(contentType, "application/"):
		return "document"
	default:
		return "document"
	}
}
