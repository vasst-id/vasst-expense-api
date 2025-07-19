package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/vasst-id/vasst-expense-api/internal/entities"
	"github.com/vasst-id/vasst-expense-api/internal/repositories"
	errorsutil "github.com/vasst-id/vasst-expense-api/internal/utils/errors"
)

//go:generate mockgen -source=transaction_tags_service.go -package=mock -destination=mock/transaction_tags_service_mock.go
type (
	TransactionTagsService interface {
		CreateTransactionTag(ctx context.Context, userID uuid.UUID, input *entities.CreateTransactionTagRequest) (*entities.TransactionTag, error)
		CreateMultipleTransactionTags(ctx context.Context, userID uuid.UUID, input *entities.CreateMultipleTransactionTagsRequest) ([]entities.TransactionTag, error)
		DeleteTransactionTag(ctx context.Context, userID uuid.UUID, transactionTagID uuid.UUID) error
		DeleteTransactionTagsByTransaction(ctx context.Context, userID uuid.UUID, transactionID uuid.UUID) error
		GetTransactionTagByID(ctx context.Context, userID uuid.UUID, transactionTagID uuid.UUID) (*entities.TransactionTag, error)
		GetTransactionTagsByTransaction(ctx context.Context, userID uuid.UUID, transactionID uuid.UUID) ([]*entities.TransactionTag, error)
		GetTransactionTagsByUserTag(ctx context.Context, userID uuid.UUID, userTagID uuid.UUID, limit, offset int) ([]*entities.TransactionTag, error)
		GetTransactionsByUserTag(ctx context.Context, userID uuid.UUID, userTagID uuid.UUID, limit, offset int) ([]*entities.TransactionWithTags, error)
		GetTaggedTransactionsSummary(ctx context.Context, userID uuid.UUID) ([]*entities.TaggedTransactionSummary, error)
	}

	transactionTagsService struct {
		transactionTagsRepo repositories.TransactionTagsRepository
		userTagsRepo        repositories.UserTagsRepository
	}
)

// NewTransactionTagsService creates a new transaction tags service
func NewTransactionTagsService(transactionTagsRepo repositories.TransactionTagsRepository, userTagsRepo repositories.UserTagsRepository) TransactionTagsService {
	return &transactionTagsService{
		transactionTagsRepo: transactionTagsRepo,
		userTagsRepo:        userTagsRepo,
	}
}

// CreateTransactionTag creates a new transaction tag
func (s *transactionTagsService) CreateTransactionTag(ctx context.Context, userID uuid.UUID, input *entities.CreateTransactionTagRequest) (*entities.TransactionTag, error) {
	// Validate required fields
	if input.TransactionID == uuid.Nil {
		return nil, errors.New("transaction ID is required")
	}
	if input.UserTagID == uuid.Nil {
		return nil, errors.New("user tag ID is required")
	}

	// Verify user owns the tag
	userTag, err := s.userTagsRepo.FindByID(ctx, input.UserTagID)
	if err != nil {
		return nil, err
	}
	if userTag == nil {
		return nil, errorsutil.New(404, "user tag not found")
	}
	if userTag.UserID != userID {
		return nil, errorsutil.New(403, "access denied")
	}

	// Create new transaction tag
	transactionTag := &entities.TransactionTag{
		TransactionTagID: uuid.New(),
		TransactionID:    input.TransactionID,
		UserTagID:        input.UserTagID,
		AppliedBy:        userID,
	}

	// Create the transaction tag - the repository will populate the struct with the actual data from DB
	createdTransactionTag, err := s.transactionTagsRepo.Create(ctx, transactionTag)
	if err != nil {
		return nil, err
	}

	// Return the transaction tag with data populated from the database
	return &createdTransactionTag, nil
}

// CreateMultipleTransactionTags creates multiple transaction tags for a single transaction
func (s *transactionTagsService) CreateMultipleTransactionTags(ctx context.Context, userID uuid.UUID, input *entities.CreateMultipleTransactionTagsRequest) ([]entities.TransactionTag, error) {
	// Validate required fields
	if input.TransactionID == uuid.Nil {
		return nil, errors.New("transaction ID is required")
	}
	if len(input.UserTagIDs) == 0 {
		return nil, errors.New("at least one user tag ID is required")
	}

	// Verify user owns all tags
	var transactionTags []*entities.TransactionTag
	for _, userTagID := range input.UserTagIDs {
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

		transactionTag := &entities.TransactionTag{
			TransactionTagID: uuid.New(),
			TransactionID:    input.TransactionID,
			UserTagID:        userTagID,
			AppliedBy:        userID,
		}
		transactionTags = append(transactionTags, transactionTag)
	}

	// Create the transaction tags - the repository will populate the structs with the actual data from DB
	createdTransactionTags, err := s.transactionTagsRepo.CreateMultiple(ctx, transactionTags)
	if err != nil {
		return nil, err
	}

	// Return the transaction tags with data populated from the database
	return createdTransactionTags, nil
}

// DeleteTransactionTag deletes a transaction tag
func (s *transactionTagsService) DeleteTransactionTag(ctx context.Context, userID uuid.UUID, transactionTagID uuid.UUID) error {
	// Get existing transaction tag
	existingTransactionTag, err := s.transactionTagsRepo.FindByID(ctx, transactionTagID)
	if err != nil {
		return err
	}
	if existingTransactionTag == nil {
		return errorsutil.New(404, "transaction tag not found")
	}

	// Verify user owns the tag
	userTag, err := s.userTagsRepo.FindByID(ctx, existingTransactionTag.UserTagID)
	if err != nil {
		return err
	}
	if userTag == nil {
		return errorsutil.New(404, "user tag not found")
	}
	if userTag.UserID != userID {
		return errorsutil.New(403, "access denied")
	}

	return s.transactionTagsRepo.Delete(ctx, transactionTagID)
}

// DeleteTransactionTagsByTransaction deletes all transaction tags for a specific transaction
func (s *transactionTagsService) DeleteTransactionTagsByTransaction(ctx context.Context, userID uuid.UUID, transactionID uuid.UUID) error {
	// Get existing transaction tags
	existingTransactionTags, err := s.transactionTagsRepo.FindByTransactionID(ctx, transactionID)
	if err != nil {
		return err
	}

	// Verify user owns all tags
	for _, transactionTag := range existingTransactionTags {
		userTag, err := s.userTagsRepo.FindByID(ctx, transactionTag.UserTagID)
		if err != nil {
			return err
		}
		if userTag == nil {
			return errorsutil.New(404, "user tag not found")
		}
		if userTag.UserID != userID {
			return errorsutil.New(403, "access denied")
		}
	}

	return s.transactionTagsRepo.DeleteByTransactionID(ctx, transactionID)
}

// GetTransactionTagByID returns a transaction tag by ID (with user ownership verification)
func (s *transactionTagsService) GetTransactionTagByID(ctx context.Context, userID uuid.UUID, transactionTagID uuid.UUID) (*entities.TransactionTag, error) {
	transactionTag, err := s.transactionTagsRepo.FindByID(ctx, transactionTagID)
	if err != nil {
		return nil, err
	}
	if transactionTag == nil {
		return nil, errorsutil.New(404, "transaction tag not found")
	}

	// Verify user owns the tag
	userTag, err := s.userTagsRepo.FindByID(ctx, transactionTag.UserTagID)
	if err != nil {
		return nil, err
	}
	if userTag == nil {
		return nil, errorsutil.New(404, "user tag not found")
	}
	if userTag.UserID != userID {
		return nil, errorsutil.New(403, "access denied")
	}

	return transactionTag, nil
}

// GetTransactionTagsByTransaction returns all transaction tags for a specific transaction
func (s *transactionTagsService) GetTransactionTagsByTransaction(ctx context.Context, userID uuid.UUID, transactionID uuid.UUID) ([]*entities.TransactionTag, error) {
	transactionTags, err := s.transactionTagsRepo.FindByTransactionID(ctx, transactionID)
	if err != nil {
		return nil, err
	}

	// Verify user owns all tags
	var userOwnedTags []*entities.TransactionTag
	for _, transactionTag := range transactionTags {
		userTag, err := s.userTagsRepo.FindByID(ctx, transactionTag.UserTagID)
		if err != nil {
			return nil, err
		}
		if userTag != nil && userTag.UserID == userID {
			userOwnedTags = append(userOwnedTags, transactionTag)
		}
	}

	return userOwnedTags, nil
}

// GetTransactionTagsByUserTag returns transaction tags for a specific user tag with pagination
func (s *transactionTagsService) GetTransactionTagsByUserTag(ctx context.Context, userID uuid.UUID, userTagID uuid.UUID, limit, offset int) ([]*entities.TransactionTag, error) {
	// Verify user owns the tag
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

	return s.transactionTagsRepo.FindByUserTagID(ctx, userTagID, limit, offset)
}

// GetTransactionsByUserTag returns transactions with their tags for a specific user tag
func (s *transactionTagsService) GetTransactionsByUserTag(ctx context.Context, userID uuid.UUID, userTagID uuid.UUID, limit, offset int) ([]*entities.TransactionWithTags, error) {
	// Verify user owns the tag
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

	return s.transactionTagsRepo.FindTransactionsByUserTagID(ctx, userTagID, limit, offset)
}

// GetTaggedTransactionsSummary returns summary of tagged transactions by user
func (s *transactionTagsService) GetTaggedTransactionsSummary(ctx context.Context, userID uuid.UUID) ([]*entities.TaggedTransactionSummary, error) {
	return s.transactionTagsRepo.GetTaggedTransactionsSummary(ctx, userID)
}
