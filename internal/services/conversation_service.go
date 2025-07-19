package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/vasst-id/vasst-expense-api/internal/entities"
	"github.com/vasst-id/vasst-expense-api/internal/repositories"
	errorsutil "github.com/vasst-id/vasst-expense-api/internal/utils/errors"
)

//go:generate mockgen -source=conversation_service.go -package=mock -destination=mock/conversation_service_mock.go
type (
	ConversationService interface {
		CreateConversation(ctx context.Context, userID uuid.UUID, input *entities.CreateConversationRequest) (*entities.Conversation, error)
		UpdateConversation(ctx context.Context, userID uuid.UUID, conversationID uuid.UUID, input *entities.UpdateConversationRequest) (*entities.Conversation, error)
		DeleteConversation(ctx context.Context, userID uuid.UUID, conversationID uuid.UUID) error
		GetConversationsByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entities.Conversation, error)
		GetActiveConversationsByUserID(ctx context.Context, userID uuid.UUID) ([]*entities.Conversation, error)
		GetConversationByID(ctx context.Context, userID uuid.UUID, conversationID uuid.UUID) (*entities.Conversation, error)
		GetSimpleConversationsByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entities.ConversationSimple, int64, error)
		GetOrCreateConversationByChannel(ctx context.Context, userID uuid.UUID, channel string) (*entities.Conversation, error)
	}

	conversationService struct {
		conversationRepo repositories.ConversationRepository
		userRepo         repositories.UserRepository
	}
)

// NewConversationService creates a new conversation service
func NewConversationService(
	conversationRepo repositories.ConversationRepository,
	userRepo repositories.UserRepository,
) ConversationService {
	return &conversationService{
		conversationRepo: conversationRepo,
		userRepo:         userRepo,
	}
}

// CreateConversation creates a new conversation
func (s *conversationService) CreateConversation(ctx context.Context, userID uuid.UUID, input *entities.CreateConversationRequest) (*entities.Conversation, error) {
	// Validate required fields
	if input.Channel == "" {
		return nil, errors.New("channel is required")
	}

	// Verify user exists
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errorsutil.New(404, "user not found")
	}

	// Create new conversation
	conversation := &entities.Conversation{
		ConversationID: uuid.New(),
		UserID:         userID,
		Channel:        input.Channel,
		IsActive:       true,
		Context:        input.Context,
		Metadata:       input.Metadata,
	}

	// Create the conversation - the repository will populate the struct with the actual data from DB
	createdConversation, err := s.conversationRepo.Create(ctx, conversation)
	if err != nil {
		return nil, err
	}

	// Return the conversation with data populated from the database
	return &createdConversation, nil
}

// UpdateConversation updates an existing conversation
func (s *conversationService) UpdateConversation(ctx context.Context, userID uuid.UUID, conversationID uuid.UUID, input *entities.UpdateConversationRequest) (*entities.Conversation, error) {
	// Get existing conversation and verify ownership
	existingConversation, err := s.conversationRepo.FindByID(ctx, conversationID)
	if err != nil {
		return nil, err
	}
	if existingConversation == nil {
		return nil, errorsutil.New(404, "conversation not found")
	}
	if existingConversation.UserID != userID {
		return nil, errorsutil.New(403, "access denied")
	}

	// Validate required fields
	if input.Channel == "" {
		return nil, errors.New("channel is required")
	}

	// Update fields
	existingConversation.Channel = input.Channel
	existingConversation.Context = input.Context
	existingConversation.Metadata = input.Metadata
	existingConversation.IsActive = input.IsActive

	// Update the conversation - the repository will populate the struct with the actual data from DB
	updatedConversation, err := s.conversationRepo.Update(ctx, existingConversation)
	if err != nil {
		return nil, err
	}

	// Return the conversation with data populated from the database
	return &updatedConversation, nil
}

// DeleteConversation deletes a conversation (soft delete)
func (s *conversationService) DeleteConversation(ctx context.Context, userID uuid.UUID, conversationID uuid.UUID) error {
	// Get existing conversation and verify ownership
	existingConversation, err := s.conversationRepo.FindByID(ctx, conversationID)
	if err != nil {
		return err
	}
	if existingConversation == nil {
		return errorsutil.New(404, "conversation not found")
	}
	if existingConversation.UserID != userID {
		return errorsutil.New(403, "access denied")
	}

	return s.conversationRepo.Delete(ctx, conversationID)
}

// GetConversationsByUserID returns conversations for a user with pagination
func (s *conversationService) GetConversationsByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entities.Conversation, error) {
	// Verify user exists
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errorsutil.New(404, "user not found")
	}

	return s.conversationRepo.FindByUserID(ctx, userID, limit, offset)
}

// GetActiveConversationsByUserID returns all active conversations for a user
func (s *conversationService) GetActiveConversationsByUserID(ctx context.Context, userID uuid.UUID) ([]*entities.Conversation, error) {
	// Verify user exists
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errorsutil.New(404, "user not found")
	}

	return s.conversationRepo.FindActiveByUserID(ctx, userID)
}

// GetConversationByID returns a conversation by ID (with user ownership verification)
func (s *conversationService) GetConversationByID(ctx context.Context, userID uuid.UUID, conversationID uuid.UUID) (*entities.Conversation, error) {
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
	return conversation, nil
}

// GetSimpleConversationsByUserID returns simplified conversations with pagination and total count
func (s *conversationService) GetSimpleConversationsByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entities.ConversationSimple, int64, error) {
	// Verify user exists
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, 0, err
	}
	if user == nil {
		return nil, 0, errorsutil.New(404, "user not found")
	}

	// Get conversations
	conversations, err := s.conversationRepo.FindSimpleByUserID(ctx, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	// Get total count
	totalCount, err := s.conversationRepo.CountByUserID(ctx, userID)
	if err != nil {
		return nil, 0, err
	}

	return conversations, totalCount, nil
}

// GetOrCreateConversationByChannel gets an existing conversation by channel or creates a new one
func (s *conversationService) GetOrCreateConversationByChannel(ctx context.Context, userID uuid.UUID, channel string) (*entities.Conversation, error) {
	// Validate required fields
	if channel == "" {
		return nil, errors.New("channel is required")
	}

	// Verify user exists
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errorsutil.New(404, "user not found")
	}

	// Try to find existing conversation by channel
	existingConversation, err := s.conversationRepo.FindByUserIDAndChannel(ctx, userID, channel)
	if err != nil {
		return nil, err
	}
	if existingConversation != nil {
		return existingConversation, nil
	}

	// Create new conversation if none exists
	input := &entities.CreateConversationRequest{
		UserID:  userID,
		Channel: channel,
	}

	return s.CreateConversation(ctx, userID, input)
}
