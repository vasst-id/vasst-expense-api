package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/vasst-id/vasst-expense-api/internal/entities"
	"github.com/vasst-id/vasst-expense-api/internal/repositories"
	errorsutil "github.com/vasst-id/vasst-expense-api/internal/utils/errors"
)

//go:generate mockgen -source=user_tags_service.go -package=mock -destination=mock/user_tags_service_mock.go
type (
	UserTagsService interface {
		CreateUserTag(ctx context.Context, userID uuid.UUID, input *entities.CreateUserTagRequest) (*entities.UserTag, error)
		UpdateUserTag(ctx context.Context, userID uuid.UUID, userTagID uuid.UUID, input *entities.UpdateUserTagRequest) (*entities.UserTag, error)
		DeleteUserTag(ctx context.Context, userID uuid.UUID, userTagID uuid.UUID) error
		GetUserTagByID(ctx context.Context, userID uuid.UUID, userTagID uuid.UUID) (*entities.UserTag, error)
		GetUserTags(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entities.UserTag, error)
		GetActiveUserTags(ctx context.Context, userID uuid.UUID) ([]*entities.UserTag, error)
		GetUserTagsWithUsage(ctx context.Context, userID uuid.UUID) ([]*entities.UserTagWithUsage, error)
	}

	userTagsService struct {
		userTagsRepo repositories.UserTagsRepository
	}
)

// NewUserTagsService creates a new user tags service
func NewUserTagsService(userTagsRepo repositories.UserTagsRepository) UserTagsService {
	return &userTagsService{
		userTagsRepo: userTagsRepo,
	}
}

// CreateUserTag creates a new user tag
func (s *userTagsService) CreateUserTag(ctx context.Context, userID uuid.UUID, input *entities.CreateUserTagRequest) (*entities.UserTag, error) {
	// Validate required fields
	if input.Name == "" {
		return nil, errors.New("tag name is required")
	}

	// Check if user tag with the same name already exists for this user
	existingUserTag, err := s.userTagsRepo.FindByNameAndUserID(ctx, userID, input.Name)
	if err != nil {
		return nil, err
	}
	if existingUserTag != nil {
		return nil, errorsutil.New(409, "tag with this name already exists")
	}

	// Create new user tag
	userTag := &entities.UserTag{
		UserTagID: uuid.New(),
		UserID:    userID,
		Name:      input.Name,
		IsActive:  true,
	}

	// Create the user tag - the repository will populate the struct with the actual data from DB
	createdUserTag, err := s.userTagsRepo.Create(ctx, userTag)
	if err != nil {
		return nil, err
	}

	// Return the user tag with data populated from the database
	return &createdUserTag, nil
}

// UpdateUserTag updates an existing user tag
func (s *userTagsService) UpdateUserTag(ctx context.Context, userID uuid.UUID, userTagID uuid.UUID, input *entities.UpdateUserTagRequest) (*entities.UserTag, error) {
	// Get existing user tag and verify ownership
	existingUserTag, err := s.userTagsRepo.FindByID(ctx, userTagID)
	if err != nil {
		return nil, err
	}
	if existingUserTag == nil {
		return nil, errorsutil.New(404, "user tag not found")
	}
	if existingUserTag.UserID != userID {
		return nil, errorsutil.New(403, "access denied")
	}

	// Validate required fields
	if input.Name == "" {
		return nil, errors.New("tag name is required")
	}

	// Check for tag name uniqueness if it's being changed
	if input.Name != existingUserTag.Name {
		userTagWithName, err := s.userTagsRepo.FindByNameAndUserID(ctx, userID, input.Name)
		if err != nil {
			return nil, err
		}
		if userTagWithName != nil && userTagWithName.UserTagID != userTagID {
			return nil, errorsutil.New(409, "tag name already in use")
		}
	}

	// Update fields
	existingUserTag.Name = input.Name
	existingUserTag.IsActive = input.IsActive

	// Update the user tag - the repository will populate the struct with the actual data from DB
	updatedUserTag, err := s.userTagsRepo.Update(ctx, existingUserTag)
	if err != nil {
		return nil, err
	}

	// Return the user tag with data populated from the database
	return &updatedUserTag, nil
}

// DeleteUserTag deletes a user tag (soft delete)
func (s *userTagsService) DeleteUserTag(ctx context.Context, userID uuid.UUID, userTagID uuid.UUID) error {
	// Get existing user tag and verify ownership
	existingUserTag, err := s.userTagsRepo.FindByID(ctx, userTagID)
	if err != nil {
		return err
	}
	if existingUserTag == nil {
		return errorsutil.New(404, "user tag not found")
	}
	if existingUserTag.UserID != userID {
		return errorsutil.New(403, "access denied")
	}

	return s.userTagsRepo.Delete(ctx, userTagID)
}

// GetUserTagByID returns a user tag by ID (with user ownership verification)
func (s *userTagsService) GetUserTagByID(ctx context.Context, userID uuid.UUID, userTagID uuid.UUID) (*entities.UserTag, error) {
	userTag, err := s.userTagsRepo.FindByID(ctx, userTagID)
	if err != nil {
		return nil, err
	}
	if userTag == nil {
		return nil, errorsutil.New(404, "user tag not found")
	}
	if userTag.UserID != userID {
		return nil, errorsutil.New(403, "access denied")
	}
	return userTag, nil
}

// GetUserTags returns user tags with pagination
func (s *userTagsService) GetUserTags(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entities.UserTag, error) {
	return s.userTagsRepo.FindByUserID(ctx, userID, limit, offset)
}

// GetActiveUserTags returns all active user tags
func (s *userTagsService) GetActiveUserTags(ctx context.Context, userID uuid.UUID) ([]*entities.UserTag, error) {
	return s.userTagsRepo.FindActiveByUserID(ctx, userID)
}

// GetUserTagsWithUsage returns user tags with their usage statistics
func (s *userTagsService) GetUserTagsWithUsage(ctx context.Context, userID uuid.UUID) ([]*entities.UserTagWithUsage, error) {
	return s.userTagsRepo.GetUserTagsWithUsage(ctx, userID)
}
