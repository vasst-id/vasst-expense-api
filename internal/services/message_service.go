package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/vasst-id/vasst-expense-api/internal/entities"
	"github.com/vasst-id/vasst-expense-api/internal/repositories"
	errorsutil "github.com/vasst-id/vasst-expense-api/internal/utils/errors"
)

//go:generate mockgen -source=message_service.go -package=mock -destination=mock/message_service_mock.go
type (
	MessageService interface {
		CreateMessage(ctx context.Context, userID uuid.UUID, input *entities.CreateMessageRequest) (*entities.Message, error)
		UpdateMessage(ctx context.Context, userID uuid.UUID, messageID uuid.UUID, input *entities.UpdateMessageRequest) (*entities.Message, error)
		DeleteMessage(ctx context.Context, userID uuid.UUID, messageID uuid.UUID) error
		GetMessagesByConversationID(ctx context.Context, userID uuid.UUID, conversationID uuid.UUID, limit, offset int) ([]*entities.Message, int64, error)
		GetMessagesByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entities.Message, error)
		GetMessageByID(ctx context.Context, userID uuid.UUID, messageID uuid.UUID) (*entities.Message, error)
		GetMessagesWithFilters(ctx context.Context, userID uuid.UUID, params *entities.MessageListParams, limit, offset int) ([]*entities.Message, int64, error)
		GetSimpleMessagesByConversationID(ctx context.Context, userID uuid.UUID, conversationID uuid.UUID, limit, offset int) ([]*entities.MessageSimple, error)
		MarkMessageAsProcessed(ctx context.Context, userID uuid.UUID, messageID uuid.UUID, aiModel string, confidenceScore *float64) error
	}

	messageService struct {
		messageRepo      repositories.MessageRepository
		conversationRepo repositories.ConversationRepository
		userRepo         repositories.UserRepository
	}
)

// NewMessageService creates a new message service
func NewMessageService(
	messageRepo repositories.MessageRepository,
	conversationRepo repositories.ConversationRepository,
	userRepo repositories.UserRepository,
) MessageService {
	return &messageService{
		messageRepo:      messageRepo,
		conversationRepo: conversationRepo,
		userRepo:         userRepo,
	}
}

// CreateMessage creates a new message
func (s *messageService) CreateMessage(ctx context.Context, userID uuid.UUID, input *entities.CreateMessageRequest) (*entities.Message, error) {
	// Validate required fields
	if input.SenderType == 0 {
		return nil, errors.New("sender type is required")
	}
	if input.Direction == "" {
		return nil, errors.New("direction is required")
	}
	if input.MessageType == 0 {
		return nil, errors.New("message type is required")
	}

	// Verify conversation exists and user has access
	conversation, err := s.conversationRepo.FindByID(ctx, input.ConversationID)
	if err != nil {
		return nil, err
	}
	if conversation == nil {
		return nil, errorsutil.New(404, "conversation not found")
	}
	if conversation.UserID != userID {
		return nil, errorsutil.New(403, "access denied to conversation")
	}

	// Create new message
	message := &entities.Message{
		MessageID:            uuid.New(),
		ConversationID:       input.ConversationID,
		UserID:               input.UserID,
		SenderType:           input.SenderType,
		Direction:            input.Direction,
		MessageType:          input.MessageType,
		Content:              input.Content,
		MediaURL:             input.MediaURL,
		Attachments:          input.Attachments,
		MediaMimeType:        input.MediaMimeType,
		AIProcessed:          false,
		RelatedTransactionID: input.RelatedTransactionID,
	}

	// Create the message - the repository will populate the struct with the actual data from DB
	createdMessage, err := s.messageRepo.Create(ctx, message)
	if err != nil {
		return nil, err
	}

	// Return the message with data populated from the database
	return &createdMessage, nil
}

// UpdateMessage updates an existing message
func (s *messageService) UpdateMessage(ctx context.Context, userID uuid.UUID, messageID uuid.UUID, input *entities.UpdateMessageRequest) (*entities.Message, error) {
	// Get existing message and verify ownership
	existingMessage, err := s.messageRepo.FindByID(ctx, messageID)
	if err != nil {
		return nil, err
	}
	if existingMessage == nil {
		return nil, errorsutil.New(404, "message not found")
	}

	// Verify conversation ownership
	conversation, err := s.conversationRepo.FindByID(ctx, existingMessage.ConversationID)
	if err != nil {
		return nil, err
	}
	if conversation == nil || conversation.UserID != userID {
		return nil, errorsutil.New(403, "access denied")
	}

	// Update fields
	if input.Content != nil {
		existingMessage.Content = input.Content
	}
	if input.MediaURL != nil {
		existingMessage.MediaURL = input.MediaURL
	}
	if input.Attachments != nil {
		existingMessage.Attachments = input.Attachments
	}
	if input.MediaMimeType != nil {
		existingMessage.MediaMimeType = input.MediaMimeType
	}
	if input.Transcription != nil {
		existingMessage.Transcription = input.Transcription
	}
	if input.AIProcessed != nil {
		existingMessage.AIProcessed = *input.AIProcessed
	}
	if input.AIModel != nil {
		existingMessage.AIModel = input.AIModel
	}
	if input.AIConfidenceScore != nil {
		existingMessage.AIConfidenceScore = input.AIConfidenceScore
	}
	if input.RelatedTransactionID != nil {
		existingMessage.RelatedTransactionID = input.RelatedTransactionID
	}

	// Update the message - the repository will populate the struct with the actual data from DB
	updatedMessage, err := s.messageRepo.Update(ctx, existingMessage)
	if err != nil {
		return nil, err
	}

	// Return the message with data populated from the database
	return &updatedMessage, nil
}

// DeleteMessage deletes a message
func (s *messageService) DeleteMessage(ctx context.Context, userID uuid.UUID, messageID uuid.UUID) error {
	// Get existing message and verify ownership
	existingMessage, err := s.messageRepo.FindByID(ctx, messageID)
	if err != nil {
		return err
	}
	if existingMessage == nil {
		return errorsutil.New(404, "message not found")
	}

	// Verify conversation ownership
	conversation, err := s.conversationRepo.FindByID(ctx, existingMessage.ConversationID)
	if err != nil {
		return err
	}
	if conversation == nil || conversation.UserID != userID {
		return errorsutil.New(403, "access denied")
	}

	return s.messageRepo.Delete(ctx, messageID)
}

// GetMessagesByConversationID returns messages for a conversation with pagination
func (s *messageService) GetMessagesByConversationID(ctx context.Context, userID uuid.UUID, conversationID uuid.UUID, limit, offset int) ([]*entities.Message, int64, error) {
	// Verify conversation exists and user has access
	conversation, err := s.conversationRepo.FindByID(ctx, conversationID)
	if err != nil {
		return nil, 0, err
	}
	if conversation == nil {
		return nil, 0, errorsutil.New(404, "conversation not found")
	}
	if conversation.UserID != userID {
		return nil, 0, errorsutil.New(403, "access denied")
	}

	// Get messages
	messages, err := s.messageRepo.FindByConversationID(ctx, conversationID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	// Get total count
	totalCount, err := s.messageRepo.CountByConversationID(ctx, conversationID)
	if err != nil {
		return nil, 0, err
	}

	return messages, totalCount, nil
}

// GetMessagesByUserID returns messages for a user with pagination
func (s *messageService) GetMessagesByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entities.Message, error) {
	// Verify user exists
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errorsutil.New(404, "user not found")
	}

	return s.messageRepo.FindByUserID(ctx, userID, limit, offset)
}

// GetMessageByID returns a message by ID (with user ownership verification)
func (s *messageService) GetMessageByID(ctx context.Context, userID uuid.UUID, messageID uuid.UUID) (*entities.Message, error) {
	message, err := s.messageRepo.FindByID(ctx, messageID)
	if err != nil {
		return nil, err
	}
	if message == nil {
		return nil, errorsutil.New(404, "message not found")
	}

	// Verify conversation ownership
	conversation, err := s.conversationRepo.FindByID(ctx, message.ConversationID)
	if err != nil {
		return nil, err
	}
	if conversation == nil || conversation.UserID != userID {
		return nil, errorsutil.New(403, "access denied")
	}

	return message, nil
}

// GetMessagesWithFilters returns messages with filtering and pagination
func (s *messageService) GetMessagesWithFilters(ctx context.Context, userID uuid.UUID, params *entities.MessageListParams, limit, offset int) ([]*entities.Message, int64, error) {
	// If filtering by conversation, verify user has access
	if params != nil && params.ConversationID != nil {
		conversation, err := s.conversationRepo.FindByID(ctx, *params.ConversationID)
		if err != nil {
			return nil, 0, err
		}
		if conversation == nil {
			return nil, 0, errorsutil.New(404, "conversation not found")
		}
		if conversation.UserID != userID {
			return nil, 0, errorsutil.New(403, "access denied")
		}
	}

	// Get messages
	messages, err := s.messageRepo.FindWithFilters(ctx, params, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	// Filter messages to only include those from user's conversations
	var filteredMessages []*entities.Message
	for _, message := range messages {
		conversation, err := s.conversationRepo.FindByID(ctx, message.ConversationID)
		if err != nil {
			continue // Skip messages from conversations we can't access
		}
		if conversation != nil && conversation.UserID == userID {
			filteredMessages = append(filteredMessages, message)
		}
	}

	// Get total count
	totalCount, err := s.messageRepo.CountWithFilters(ctx, params)
	if err != nil {
		return nil, 0, err
	}

	return filteredMessages, totalCount, nil
}

// GetSimpleMessagesByConversationID returns simplified messages for a conversation
func (s *messageService) GetSimpleMessagesByConversationID(ctx context.Context, userID uuid.UUID, conversationID uuid.UUID, limit, offset int) ([]*entities.MessageSimple, error) {
	// Verify conversation exists and user has access
	conversation, err := s.conversationRepo.FindByID(ctx, conversationID)
	if err != nil {
		return nil, err
	}
	if conversation == nil {
		return nil, errorsutil.New(404, "conversation not found")
	}
	if conversation.UserID != userID {
		return nil, errorsutil.New(403, "access denied")
	}

	return s.messageRepo.FindSimpleByConversationID(ctx, conversationID, limit, offset)
}

// MarkMessageAsProcessed marks a message as AI processed
func (s *messageService) MarkMessageAsProcessed(ctx context.Context, userID uuid.UUID, messageID uuid.UUID, aiModel string, confidenceScore *float64) error {
	// Get existing message and verify ownership
	existingMessage, err := s.messageRepo.FindByID(ctx, messageID)
	if err != nil {
		return err
	}
	if existingMessage == nil {
		return errorsutil.New(404, "message not found")
	}

	// Verify conversation ownership
	conversation, err := s.conversationRepo.FindByID(ctx, existingMessage.ConversationID)
	if err != nil {
		return err
	}
	if conversation == nil || conversation.UserID != userID {
		return errorsutil.New(403, "access denied")
	}

	return s.messageRepo.MarkAsProcessed(ctx, messageID, aiModel, confidenceScore)
}
