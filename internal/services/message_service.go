package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"

	"github.com/google/uuid"
	"github.com/vasst-id/vasst-expense-api/internal/entities"
	"github.com/vasst-id/vasst-expense-api/internal/repositories"
	errorsutil "github.com/vasst-id/vasst-expense-api/internal/utils/errors"
)

//go:generate mockgen -source=message_service.go -package=mock -destination=mock/message_service_mock.go
type (
	MessageService interface {
		CreateMessage(ctx context.Context, input *entities.CreateMessageInput, userID uuid.UUID) (*entities.Message, error)
		CreateMessageWithMedia(ctx context.Context, input *entities.CreateMessageInput, file *multipart.FileHeader) (*entities.Message, error)
		UpdateMessage(ctx context.Context, messageID, organizationID uuid.UUID, input *entities.UpdateMessageInput) (*entities.Message, error)
		UpdateMessageStatus(ctx context.Context, messageID, organizationID uuid.UUID, input *entities.UpdateMessageStatusInput) error
		DeleteMessage(ctx context.Context, messageID, organizationID uuid.UUID) error
		GetMessageByID(ctx context.Context, messageID, organizationID uuid.UUID) (*entities.Message, error)
		ListMessagesByConversation(ctx context.Context, conversationID, organizationID uuid.UUID, limit, offset int) ([]*entities.Message, error)
		ListMessagesByOrganization(ctx context.Context, organizationID uuid.UUID, limit, offset int) ([]*entities.Message, error)
		GetPendingMessages(ctx context.Context, organizationID uuid.UUID, limit, offset int) ([]*entities.Message, error)
		GetMessagesByStatus(ctx context.Context, organizationID uuid.UUID, status int, limit, offset int) ([]*entities.Message, error)
		GetMessagesBySenderType(ctx context.Context, organizationID uuid.UUID, senderTypeID int, limit, offset int) ([]*entities.Message, error)
		GetMessagesBySenderID(ctx context.Context, organizationID uuid.UUID, senderID uuid.UUID, limit, offset int) ([]*entities.Message, error)
	}

	messageService struct {
		messageRepo      repositories.MessageRepository
		conversationRepo repositories.ConversationRepository
		conversationSvc  ConversationService
		storageSvc       GoogleStorageService
		organizationRepo repositories.OrganizationRepository
	}
)

// NewMessageService creates a new message service
func NewMessageService(messageRepo repositories.MessageRepository, conversationRepo repositories.ConversationRepository, conversationSvc ConversationService, storageSvc GoogleStorageService) MessageService {
	return &messageService{
		messageRepo:      messageRepo,
		conversationRepo: conversationRepo,
		conversationSvc:  conversationSvc,
		storageSvc:       storageSvc,
	}
}

// CreateMessage creates a new message
func (s *messageService) CreateMessage(ctx context.Context, input *entities.CreateMessageInput, userID uuid.UUID) (*entities.Message, error) {
	// Validate input

	if input.OrganizationID == uuid.Nil {
		return nil, errors.New("organization ID is required")
	}

	fmt.Println("input", input)

	switch input.MessageTypeID {
	case entities.MessageTypeText:
		if input.Content == "" {
			return nil, errors.New("message content is required")
		}
	case entities.MessageTypeImage, entities.MessageTypeVideo, entities.MessageTypeAudio, entities.MessageTypeDocument, entities.MessageTypeLocation, entities.MessageTypeContact, entities.MessageTypeSticker:
		if input.MediaURL == "" {
			return nil, errors.New("media URL is required")
		}
	default:
		return nil, errors.New("invalid message type")
	}

	// Set default status to pending if not provided
	if input.Status == 0 {
		input.Status = int(entities.MessageStatusPending)
	}

	// Handle conversation logic
	var conversationID uuid.UUID
	var err error

	fmt.Println("input", input)

	if input.ConversationID != uuid.Nil {
		// Conversation ID provided, verify it exists and belongs to organization
		conversation, err := s.conversationRepo.FindByID(ctx, input.ConversationID, input.OrganizationID)
		if err != nil {
			return nil, err
		}
		if conversation == nil {
			return nil, errorsutil.New(404, "conversation not found")
		}
		conversationID = input.ConversationID
	} else {
		// No conversation ID provided, need to find or create conversation
		// Extract user ID from context (set by auth middleware)

		if userID == uuid.Nil {
			return nil, errors.New("user ID not found in context")
		}

		// Extract contact ID from input or context
		contactID := input.ContactID
		if contactID == uuid.Nil {
			return nil, errors.New("contact ID is required when conversation ID is not provided")
		}

		// Default medium ID to WhatsApp (1) if not provided
		mediumID := input.MediumID
		if mediumID == 0 {
			mediumID = 1 // WhatsApp default
		}

		// Try to find existing active conversation
		existingConversation, err := s.conversationRepo.FindActiveByOrganizationUserContactMedium(ctx, input.OrganizationID, userID, contactID, mediumID)
		if err != nil {
			return nil, err
		}

		if existingConversation != nil {
			conversationID = existingConversation.ConversationID
		} else {
			// Create new conversation
			conversationInput := &entities.CreateConversationInput{
				OrganizationID:      input.OrganizationID,
				UserID:              userID,
				ContactID:           contactID,
				MediumID:            mediumID,
				Status:              int(entities.ConversationStatusOpen),
				Priority:            int(entities.ConversationPriorityLow),
				AIEnabled:           true, // Default to true
				LastMessageByID:     &userID,
				LastMessageByType:   int(entities.SenderTypeCustomer),
				LastMessageByName:   "Customer",
				LastMessageContent:  input.Content,
				LastMessageTypeID:   input.MessageTypeID,
				LastMessageMediaURL: input.MediaURL,
			}

			newConversation, err := s.conversationSvc.CreateConversation(ctx, conversationInput)
			if err != nil {
				return nil, err
			}
			conversationID = newConversation.ConversationID
		}
	}

	// Prepare attachments and metadata
	var attachmentsJSON json.RawMessage
	if input.Attachments != nil {
		attachmentsJSON, err = json.Marshal(input.Attachments)
		if err != nil {
			return nil, err
		}
	} else {
		attachmentsJSON = json.RawMessage("[]")
	}

	var metadataJSON json.RawMessage
	if input.Metadata != nil {
		metadataJSON, err = json.Marshal(input.Metadata)
		if err != nil {
			return nil, err
		}
	} else {
		metadataJSON = json.RawMessage("{}")
	}

	// Create new message
	message := &entities.Message{
		ConversationID:    conversationID,
		OrganizationID:    input.OrganizationID,
		SenderTypeID:      input.SenderTypeID,
		SenderID:          input.SenderID,
		Direction:         input.Direction,
		MessageTypeID:     input.MessageTypeID,
		Content:           input.Content,
		MediaURL:          input.MediaURL,
		Attachments:       attachmentsJSON,
		IsBroadcast:       input.IsBroadcast,
		IsOrderMessage:    input.IsOrderMessage,
		Metadata:          metadataJSON,
		AIGenerated:       input.AIGenerated,
		AIConfidenceScore: input.AIConfidenceScore,
		Status:            input.Status,
	}

	if err := s.messageRepo.Create(ctx, message); err != nil {
		return nil, err
	}

	// Update conversation's last message information
	lastMessageInput := &entities.UpdateConversationLastMessageInput{
		LastMessageAt:       message.CreatedAt,
		LastMessageByType:   message.SenderTypeID,
		LastMessageContent:  message.Content,
		LastMessageTypeID:   message.MessageTypeID,
		LastMessageMediaURL: message.MediaURL,
	}

	// Handle nullable SenderID
	if message.SenderID != nil {
		lastMessageInput.LastMessageByID = *message.SenderID
	}

	lastMessageInput.LastMessageByName = s.getSenderName(ctx, message.SenderTypeID, message.SenderID)

	// Update conversation last message (don't fail if this fails)
	if err := s.conversationRepo.UpdateLastMessage(ctx, conversationID, lastMessageInput); err != nil {
		// Log error but don't fail the message creation
		// TODO: Add proper logging
	}

	return message, nil
}

// CreateMessageWithMedia creates a new message with media file
func (s *messageService) CreateMessageWithMedia(ctx context.Context, input *entities.CreateMessageInput, file *multipart.FileHeader) (*entities.Message, error) {
	// Get user ID from context
	userID, exists := ctx.Value("user_id").(uuid.UUID)
	if !exists {
		return nil, errors.New("user ID not found in context")
	}

	// Validate that message type supports media
	if !s.isMediaMessageType(input.MessageTypeID) {
		return nil, errors.New("message type does not support media")
	}

	// First create the message to get the conversation ID
	message, err := s.CreateMessage(ctx, input, userID)
	if err != nil {
		return nil, err
	}

	// Get organization code
	organization, err := s.organizationRepo.FindOrganizationByID(ctx, input.OrganizationID)
	if err != nil {
		return nil, err
	}
	if organization == nil {
		return nil, errorsutil.New(404, "organization not found")
	}

	// Upload file to Google Cloud Storage using conversation bucket
	uploadResult, err := s.storageSvc.UploadFile(ctx, organization.OrganizationCode, message.ConversationID, file)
	if err != nil {
		return nil, fmt.Errorf("failed to upload media: %w", err)
	}

	// Update the message with the media URL
	message.MediaURL = uploadResult.FileURL

	// Add file metadata to attachments
	attachment := map[string]interface{}{
		"file_name":    uploadResult.FileName,
		"file_url":     uploadResult.FileURL,
		"file_size":    uploadResult.FileSize,
		"content_type": uploadResult.ContentType,
		"uploaded_at":  uploadResult.UploadedAt,
		"bucket_name":  uploadResult.BucketName,
		"object_name":  uploadResult.ObjectName,
	}

	// Parse existing attachments or create new array
	var attachments []map[string]interface{}
	if message.Attachments != nil {
		if err := json.Unmarshal(message.Attachments, &attachments); err != nil {
			attachments = []map[string]interface{}{}
		}
	}

	// Add new attachment
	attachments = append(attachments, attachment)

	// Update attachments JSON
	attachmentsJSON, err := json.Marshal(attachments)
	if err != nil {
		return nil, err
	}
	message.Attachments = attachmentsJSON

	// Update the message in the database
	if err := s.messageRepo.Update(ctx, message); err != nil {
		return nil, err
	}

	return message, nil
}

// UpdateMessage updates an existing message
func (s *messageService) UpdateMessage(ctx context.Context, messageID, organizationID uuid.UUID, input *entities.UpdateMessageInput) (*entities.Message, error) {
	// Check if message exists and belongs to organization
	existingMessage, err := s.messageRepo.FindByID(ctx, messageID)
	if err != nil {
		return nil, err
	}

	if existingMessage == nil {
		return nil, errorsutil.New(404, "message not found")
	}

	if existingMessage.OrganizationID != organizationID {
		return nil, errorsutil.New(403, "access denied")
	}

	// Update message fields
	if input.SenderTypeID != nil {
		existingMessage.SenderTypeID = *input.SenderTypeID
	}
	if input.SenderID != nil {
		existingMessage.SenderID = input.SenderID
	}
	if input.Direction != "" {
		existingMessage.Direction = input.Direction
	}
	if input.MessageTypeID != 0 {
		existingMessage.MessageTypeID = input.MessageTypeID
	}
	if input.Content != "" {
		existingMessage.Content = input.Content
	}
	if input.MediaURL != "" {
		existingMessage.MediaURL = input.MediaURL
	}
	if input.Attachments != nil {
		attachmentsJSON, err := json.Marshal(input.Attachments)
		if err != nil {
			return nil, err
		}
		existingMessage.Attachments = attachmentsJSON
	}
	if input.IsBroadcast != nil {
		existingMessage.IsBroadcast = *input.IsBroadcast
	}
	if input.IsOrderMessage != nil {
		existingMessage.IsOrderMessage = *input.IsOrderMessage
	}
	if input.Metadata != nil {
		metadataJSON, err := json.Marshal(input.Metadata)
		if err != nil {
			return nil, err
		}
		existingMessage.Metadata = metadataJSON
	}
	if input.Status != nil {
		existingMessage.Status = *input.Status
	}
	if input.AIGenerated != nil {
		existingMessage.AIGenerated = *input.AIGenerated
	}
	if input.AIConfidenceScore != nil {
		existingMessage.AIConfidenceScore = input.AIConfidenceScore
	}
	if input.FailureReason != nil {
		existingMessage.FailureReason = input.FailureReason
	}

	if err := s.messageRepo.Update(ctx, existingMessage); err != nil {
		return nil, err
	}

	return existingMessage, nil
}

// UpdateMessageStatus updates only the status of a message
func (s *messageService) UpdateMessageStatus(ctx context.Context, messageID, organizationID uuid.UUID, input *entities.UpdateMessageStatusInput) error {
	// Check if message exists and belongs to organization
	existingMessage, err := s.messageRepo.FindByID(ctx, messageID)
	if err != nil {
		return err
	}

	if existingMessage == nil {
		return errorsutil.New(404, "message not found")
	}

	if existingMessage.OrganizationID != organizationID {
		return errorsutil.New(403, "access denied")
	}

	return s.messageRepo.UpdateStatus(ctx, messageID, input.Status, input.FailureReason)
}

// DeleteMessage deletes a message by its ID
func (s *messageService) DeleteMessage(ctx context.Context, messageID, organizationID uuid.UUID) error {
	// Check if message exists and belongs to organization
	existingMessage, err := s.messageRepo.FindByID(ctx, messageID)
	if err != nil {
		return err
	}

	if existingMessage == nil {
		return errorsutil.New(404, "message not found")
	}

	if existingMessage.OrganizationID != organizationID {
		return errorsutil.New(403, "access denied")
	}

	return s.messageRepo.Delete(ctx, messageID)
}

// GetMessageByID returns a message by its ID (organization-scoped)
func (s *messageService) GetMessageByID(ctx context.Context, messageID, organizationID uuid.UUID) (*entities.Message, error) {
	message, err := s.messageRepo.FindByID(ctx, messageID)
	if err != nil {
		return nil, err
	}

	if message == nil {
		return nil, errorsutil.New(404, "message not found")
	}

	if message.OrganizationID != organizationID {
		return nil, errorsutil.New(403, "access denied")
	}

	return message, nil
}

// ListMessagesByConversation returns messages for a conversation (organization-scoped)
func (s *messageService) ListMessagesByConversation(ctx context.Context, conversationID, organizationID uuid.UUID, limit, offset int) ([]*entities.Message, error) {
	return s.messageRepo.FindByConversationAndOrganization(ctx, conversationID, organizationID, limit, offset)
}

// ListMessagesByOrganization returns messages for an organization
func (s *messageService) ListMessagesByOrganization(ctx context.Context, organizationID uuid.UUID, limit, offset int) ([]*entities.Message, error) {
	return s.messageRepo.FindByOrganizationID(ctx, organizationID, limit, offset)
}

// GetPendingMessages returns pending messages for an organization
func (s *messageService) GetPendingMessages(ctx context.Context, organizationID uuid.UUID, limit, offset int) ([]*entities.Message, error) {
	return s.messageRepo.FindPendingMessages(ctx, limit, offset)
}

// GetMessagesByStatus returns messages by status for an organization
func (s *messageService) GetMessagesByStatus(ctx context.Context, organizationID uuid.UUID, status int, limit, offset int) ([]*entities.Message, error) {
	return s.messageRepo.FindByStatus(ctx, status, limit, offset)
}

// GetMessagesBySenderType returns messages by sender type for an organization
func (s *messageService) GetMessagesBySenderType(ctx context.Context, organizationID uuid.UUID, senderTypeID int, limit, offset int) ([]*entities.Message, error) {
	return s.messageRepo.FindBySenderType(ctx, senderTypeID, limit, offset)
}

// GetMessagesBySenderID returns messages by sender ID for an organization
func (s *messageService) GetMessagesBySenderID(ctx context.Context, organizationID uuid.UUID, senderID uuid.UUID, limit, offset int) ([]*entities.Message, error) {
	return s.messageRepo.FindBySenderID(ctx, senderID, limit, offset)
}

// Helper methods

// isValidMessageType checks if the message type is valid
func (s *messageService) isValidMessageType(messageTypeID int) bool {
	// TODO: Implement validation against message_type table
	// For now, accept any positive integer
	return messageTypeID > 0
}

// isMediaMessageType checks if the message type supports media
func (s *messageService) isMediaMessageType(messageTypeID int) bool {
	mediaTypes := []int{
		entities.MessageTypeImage,
		entities.MessageTypeVideo,
		entities.MessageTypeAudio,
		entities.MessageTypeDocument,
		entities.MessageTypeSticker,
	}

	for _, mediaType := range mediaTypes {
		if messageTypeID == mediaType {
			return true
		}
	}
	return false
}

// uploadMediaToStorage uploads a media file to Google Cloud Storage
// This method is deprecated and replaced by Google Cloud Storage service integration
func (s *messageService) uploadMediaToStorage(ctx context.Context, file *multipart.FileHeader) (string, error) {
	// This method is deprecated - use Google Cloud Storage service instead
	return "", errors.New("uploadMediaToStorage is deprecated - use Google Cloud Storage service")
}

// getSenderName gets the sender name based on sender type and ID
func (s *messageService) getSenderName(ctx context.Context, senderTypeID int, senderID *uuid.UUID) string {
	// TODO: Implement logic to get sender name from user/contact tables
	// For now, return a placeholder
	switch senderTypeID {
	case int(entities.SenderTypeCustomer):
		return "Customer"
	case int(entities.SenderTypeAgent):
		return "Agent"
	case int(entities.SenderTypeAI):
		return "AI Assistant"
	case int(entities.SenderTypeSystem):
		return "System"
	default:
		return "Unknown"
	}
}
