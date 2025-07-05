package services

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"github.com/vasst-id/vasst-expense-api/internal/entities"
	"github.com/vasst-id/vasst-expense-api/internal/repositories"
	errorsutil "github.com/vasst-id/vasst-expense-api/internal/utils/errors"
)

//go:generate mockgen -source=conversation_service.go -package=mock -destination=mock/conversation_service_mock.go
type (
	ConversationService interface {
		CreateConversation(ctx context.Context, input *entities.CreateConversationInput) (*entities.Conversation, error)
		UpdateConversation(ctx context.Context, conversationID, organizationID uuid.UUID, input *entities.UpdateConversationInput) (*entities.Conversation, error)
		UpdateConversationLastMessage(ctx context.Context, conversationID uuid.UUID, input *entities.UpdateConversationLastMessageInput) error
		DeleteConversation(ctx context.Context, conversationID, organizationID uuid.UUID) error

		// Organization-scoped operations
		GetConversationByID(ctx context.Context, conversationID, organizationID uuid.UUID) (*entities.Conversation, error)
		GetConversationDetail(ctx context.Context, conversationID, organizationID uuid.UUID, messageLimit, messageOffset int) (*entities.ConversationDetailResponse, error)
		ListConversationsByOrganization(ctx context.Context, organizationID uuid.UUID, limit, offset int) ([]*entities.ConversationListResponse, error)
		ListConversationsByOrganizationWithFilters(ctx context.Context, organizationID uuid.UUID, status, priority *int, isActive *bool, limit, offset int) ([]*entities.ConversationListResponse, error)

		// Find active conversation for specific parameters
		GetActiveConversation(ctx context.Context, organizationID, userID, contactID uuid.UUID, mediumID int) (*entities.Conversation, error)

		// Filtered queries (organization-scoped)
		GetConversationsByUserID(ctx context.Context, organizationID, userID uuid.UUID, limit, offset int) ([]*entities.ConversationListResponse, error)
		GetConversationsByContactID(ctx context.Context, organizationID, contactID uuid.UUID, limit, offset int) ([]*entities.ConversationListResponse, error)
		GetConversationsByStatus(ctx context.Context, organizationID uuid.UUID, status int, limit, offset int) ([]*entities.ConversationListResponse, error)
		GetConversationsByPriority(ctx context.Context, organizationID uuid.UUID, priority int, limit, offset int) ([]*entities.ConversationListResponse, error)

		// Count operations
		GetConversationCountByOrganization(ctx context.Context, organizationID uuid.UUID) (int, error)
	}

	conversationService struct {
		conversationRepo repositories.ConversationRepository
		messageRepo      repositories.MessageRepository
	}
)

// NewConversationService creates a new conversation service
func NewConversationService(conversationRepo repositories.ConversationRepository, messageRepo repositories.MessageRepository) ConversationService {
	return &conversationService{
		conversationRepo: conversationRepo,
		messageRepo:      messageRepo,
	}
}

// CreateConversation creates a new conversation
func (s *conversationService) CreateConversation(ctx context.Context, input *entities.CreateConversationInput) (*entities.Conversation, error) {
	// Validate input
	if input.OrganizationID == uuid.Nil {
		return nil, errors.New("organization ID is required")
	}

	if input.UserID == uuid.Nil {
		return nil, errors.New("user ID is required")
	}

	if input.ContactID == uuid.Nil {
		return nil, errors.New("contact ID is required")
	}

	if input.MediumID == 0 {
		return nil, errors.New("medium ID is required")
	}

	// Check if there's already an active conversation for these parameters
	existingConversation, err := s.conversationRepo.FindActiveByOrganizationUserContactMedium(ctx, input.OrganizationID, input.UserID, input.ContactID, input.MediumID)
	if err != nil {
		return nil, err
	}

	if existingConversation != nil {
		return existingConversation, nil // Return existing conversation instead of creating new one
	}

	// Prepare AI config
	var aiConfigJSON json.RawMessage
	if input.AIConfig != nil {
		aiConfigJSON, err = json.Marshal(input.AIConfig)
		if err != nil {
			return nil, err
		}
	} else {
		aiConfigJSON = json.RawMessage("{}")
	}

	// Prepare metadata
	var metadataJSON json.RawMessage
	if input.Metadata != nil {
		metadataJSON, err = json.Marshal(input.Metadata)
		if err != nil {
			return nil, err
		}
	} else {
		metadataJSON = json.RawMessage("{}")
	}

	// Set default values
	status := input.Status
	if status == 0 {
		status = int(entities.ConversationStatusOpen)
	}

	priority := input.Priority
	if priority == 0 {
		priority = int(entities.ConversationPriorityLow)
	}

	aiEnabled := input.AIEnabled
	if !aiEnabled {
		aiEnabled = true // Default to true
	}

	// Create new conversation
	conversation := &entities.Conversation{
		ConversationID:      uuid.New(),
		OrganizationID:      input.OrganizationID,
		UserID:              input.UserID,
		ContactID:           input.ContactID,
		MediumID:            input.MediumID,
		Status:              status,
		Priority:            priority,
		AIEnabled:           aiEnabled,
		AIConfig:            aiConfigJSON,
		Metadata:            metadataJSON,
		LastMessageByID:     input.LastMessageByID,
		LastMessageByType:   input.LastMessageByType,
		LastMessageByName:   input.LastMessageByName,
		LastMessageContent:  input.LastMessageContent,
		LastMessageTypeID:   input.LastMessageTypeID,
		LastMessageMediaURL: input.LastMessageMediaURL,
		IsActive:            true,
		IsArchived:          false,
		IsDeleted:           false,
	}

	if err := s.conversationRepo.Create(ctx, conversation); err != nil {
		return nil, err
	}

	// Deactivate other conversations for the same parameters
	if err := s.conversationRepo.DeactivateOtherConversations(ctx, input.OrganizationID, input.UserID, input.ContactID, input.MediumID, conversation.ConversationID); err != nil {
		// Log error but don't fail the creation
		// TODO: Add proper logging
	}

	return conversation, nil
}

// UpdateConversation updates an existing conversation
func (s *conversationService) UpdateConversation(ctx context.Context, conversationID, organizationID uuid.UUID, input *entities.UpdateConversationInput) (*entities.Conversation, error) {
	// Check if conversation exists and belongs to organization
	existingConversation, err := s.conversationRepo.FindByID(ctx, conversationID, organizationID)
	if err != nil {
		return nil, err
	}

	if existingConversation == nil {
		return nil, errorsutil.New(404, "conversation not found")
	}

	// Update conversation fields
	if input.UserID != uuid.Nil {
		existingConversation.UserID = input.UserID
	}
	if input.MediumID != 0 {
		existingConversation.MediumID = input.MediumID
	}
	if input.Status != nil {
		existingConversation.Status = *input.Status
	}
	if input.Priority != nil {
		existingConversation.Priority = *input.Priority
	}
	if input.AIEnabled != nil {
		existingConversation.AIEnabled = *input.AIEnabled
	}
	if input.AIConfig != nil {
		aiConfigJSON, err := json.Marshal(input.AIConfig)
		if err != nil {
			return nil, err
		}
		existingConversation.AIConfig = aiConfigJSON
	}
	if input.Metadata != nil {
		metadataJSON, err := json.Marshal(input.Metadata)
		if err != nil {
			return nil, err
		}
		existingConversation.Metadata = metadataJSON
	}
	if input.IsActive != nil {
		existingConversation.IsActive = *input.IsActive
	}
	if input.IsArchived != nil {
		existingConversation.IsArchived = *input.IsArchived
	}

	if err := s.conversationRepo.Update(ctx, existingConversation); err != nil {
		return nil, err
	}

	return existingConversation, nil
}

// UpdateConversationLastMessage updates last message information
func (s *conversationService) UpdateConversationLastMessage(ctx context.Context, conversationID uuid.UUID, input *entities.UpdateConversationLastMessageInput) error {
	return s.conversationRepo.UpdateLastMessage(ctx, conversationID, input)
}

// DeleteConversation deletes a conversation by its ID
func (s *conversationService) DeleteConversation(ctx context.Context, conversationID, organizationID uuid.UUID) error {
	// Check if conversation exists and belongs to organization
	existingConversation, err := s.conversationRepo.FindByID(ctx, conversationID, organizationID)
	if err != nil {
		return err
	}

	if existingConversation == nil {
		return errorsutil.New(404, "conversation not found")
	}

	return s.conversationRepo.Delete(ctx, conversationID)
}

// GetConversationByID returns a conversation by its ID (organization-scoped)
func (s *conversationService) GetConversationByID(ctx context.Context, conversationID, organizationID uuid.UUID) (*entities.Conversation, error) {
	conversation, err := s.conversationRepo.FindByID(ctx, conversationID, organizationID)
	if err != nil {
		return nil, err
	}

	if conversation == nil {
		return nil, errorsutil.New(404, "conversation not found")
	}

	return conversation, nil
}

// GetConversationDetail returns a conversation with messages
func (s *conversationService) GetConversationDetail(ctx context.Context, conversationID, organizationID uuid.UUID, messageLimit, messageOffset int) (*entities.ConversationDetailResponse, error) {
	// Get conversation
	conversation, err := s.conversationRepo.FindByID(ctx, conversationID, organizationID)
	if err != nil {
		return nil, err
	}

	if conversation == nil {
		return nil, errorsutil.New(404, "conversation not found")
	}

	// Set default message limit if not provided
	if messageLimit == 0 {
		messageLimit = 50
	}

	// Get messages for this conversation
	messages, err := s.messageRepo.FindByConversationAndOrganization(ctx, conversationID, organizationID, messageLimit, messageOffset)
	if err != nil {
		return nil, err
	}

	// Check if there are more messages
	hasMoreMessages := false
	if len(messages) == messageLimit {
		// Check if there's one more message
		nextMessages, err := s.messageRepo.FindByConversationAndOrganization(ctx, conversationID, organizationID, 1, messageOffset+messageLimit)
		if err == nil && len(nextMessages) > 0 {
			hasMoreMessages = true
		}
	}

	// Count total messages
	totalMessages := len(messages) + messageOffset
	if hasMoreMessages {
		// Get total count (this is a simplified approach)
		allMessages, err := s.messageRepo.FindByConversationAndOrganization(ctx, conversationID, organizationID, 1000, 0)
		if err == nil {
			totalMessages = len(allMessages)
		}
	}

	return &entities.ConversationDetailResponse{
		Conversation:    conversation,
		Messages:        messages,
		TotalMessages:   totalMessages,
		HasMoreMessages: hasMoreMessages,
	}, nil
}

// ListConversationsByOrganization returns conversations for an organization
func (s *conversationService) ListConversationsByOrganization(ctx context.Context, organizationID uuid.UUID, limit, offset int) ([]*entities.ConversationListResponse, error) {
	return s.conversationRepo.ListByOrganization(ctx, organizationID, limit, offset)
}

// ListConversationsByOrganizationWithFilters returns conversations with filters
func (s *conversationService) ListConversationsByOrganizationWithFilters(ctx context.Context, organizationID uuid.UUID, status, priority *int, isActive *bool, limit, offset int) ([]*entities.ConversationListResponse, error) {
	return s.conversationRepo.ListByOrganizationWithFilters(ctx, organizationID, status, priority, isActive, limit, offset)
}

// GetActiveConversation finds the active conversation for specific parameters
func (s *conversationService) GetActiveConversation(ctx context.Context, organizationID, userID, contactID uuid.UUID, mediumID int) (*entities.Conversation, error) {
	return s.conversationRepo.FindActiveByOrganizationUserContactMedium(ctx, organizationID, userID, contactID, mediumID)
}

// GetConversationsByUserID returns conversations for a user (organization-scoped)
func (s *conversationService) GetConversationsByUserID(ctx context.Context, organizationID, userID uuid.UUID, limit, offset int) ([]*entities.ConversationListResponse, error) {
	return s.conversationRepo.FindByUserID(ctx, organizationID, userID, limit, offset)
}

// GetConversationsByContactID returns conversations for a contact (organization-scoped)
func (s *conversationService) GetConversationsByContactID(ctx context.Context, organizationID, contactID uuid.UUID, limit, offset int) ([]*entities.ConversationListResponse, error) {
	return s.conversationRepo.FindByContactID(ctx, organizationID, contactID, limit, offset)
}

// GetConversationsByStatus returns conversations by status (organization-scoped)
func (s *conversationService) GetConversationsByStatus(ctx context.Context, organizationID uuid.UUID, status int, limit, offset int) ([]*entities.ConversationListResponse, error) {
	return s.conversationRepo.FindByStatus(ctx, organizationID, status, limit, offset)
}

// GetConversationsByPriority returns conversations by priority (organization-scoped)
func (s *conversationService) GetConversationsByPriority(ctx context.Context, organizationID uuid.UUID, priority int, limit, offset int) ([]*entities.ConversationListResponse, error) {
	return s.conversationRepo.FindByPriority(ctx, organizationID, priority, limit, offset)
}

// GetConversationCountByOrganization returns total count of conversations for an organization
func (s *conversationService) GetConversationCountByOrganization(ctx context.Context, organizationID uuid.UUID) (int, error) {
	return s.conversationRepo.CountByOrganization(ctx, organizationID)
}
