package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/vasst-id/vasst-expense-api/internal/entities"
	"github.com/vasst-id/vasst-expense-api/internal/repositories"
	errorsutil "github.com/vasst-id/vasst-expense-api/internal/utils/errors"
)

//go:generate mockgen -source=transaction_service.go -package=mock -destination=mock/transaction_service_mock.go
type (
	TransactionService interface {
		CreateTransaction(ctx context.Context, userID uuid.UUID, input *entities.CreateTransactionRequest) (*entities.Transaction, error)
		UpdateTransaction(ctx context.Context, userID uuid.UUID, transactionID uuid.UUID, input *entities.UpdateTransactionRequest) (*entities.Transaction, error)
		DeleteTransaction(ctx context.Context, userID uuid.UUID, transactionID uuid.UUID) error
		GetTransactionsByWorkspace(ctx context.Context, userID uuid.UUID, workspaceID uuid.UUID, params *entities.TransactionListParams, limit, offset int) ([]*entities.Transaction, int64, error)
		GetTransactionsByAccount(ctx context.Context, userID uuid.UUID, accountID uuid.UUID, limit, offset int) ([]*entities.Transaction, error)
		GetTransactionsByCategory(ctx context.Context, userID uuid.UUID, categoryID uuid.UUID, limit, offset int) ([]*entities.Transaction, error)
		GetTransactionByID(ctx context.Context, userID uuid.UUID, transactionID uuid.UUID) (*entities.Transaction, error)
	}

	transactionService struct {
		transactionRepo repositories.TransactionRepository
		workspaceRepo   repositories.WorkspaceRepository
		accountRepo     repositories.AccountRepository
	}
)

// NewTransactionService creates a new transaction service
func NewTransactionService(
	transactionRepo repositories.TransactionRepository,
	workspaceRepo repositories.WorkspaceRepository,
	accountRepo repositories.AccountRepository,
) TransactionService {
	return &transactionService{
		transactionRepo: transactionRepo,
		workspaceRepo:   workspaceRepo,
		accountRepo:     accountRepo,
	}
}

// CreateTransaction creates a new transaction
func (s *transactionService) CreateTransaction(ctx context.Context, userID uuid.UUID, input *entities.CreateTransactionRequest) (*entities.Transaction, error) {
	// Validate required fields
	if input.Description == "" {
		return nil, errors.New("description is required")
	}
	if input.Amount == 0 {
		return nil, errors.New("amount is required")
	}
	if input.TransactionType == 0 {
		return nil, errors.New("transaction type is required")
	}
	// if input.PaymentMethod == 0 {
	// 	return nil, errors.New("payment method is required")
	// }
	if input.TransactionDate.IsZero() {
		return nil, errors.New("transaction date is required")
	}

	// Verify workspace ownership
	workspace, err := s.workspaceRepo.FindByID(ctx, input.WorkspaceID)
	if err != nil {
		return nil, err
	}
	if workspace == nil {
		return nil, errorsutil.New(404, "workspace not found")
	}
	if workspace.CreatedBy != userID {
		return nil, errorsutil.New(403, "access denied to workspace")
	}

	// Verify account ownership if account is specified
	if input.AccountID != uuid.Nil {
		account, err := s.accountRepo.FindByID(ctx, input.AccountID)
		if err != nil {
			return nil, err
		}
		if account == nil {
			return nil, errorsutil.New(404, "account not found")
		}
		if account.UserID != userID {
			return nil, errorsutil.New(403, "access denied to account")
		}
	}

	// Create new transaction
	transaction := &entities.Transaction{
		TransactionID:   uuid.New(),
		WorkspaceID:     &input.WorkspaceID,
		AccountID:       &input.AccountID,
		CategoryID:      input.CategoryID,
		Description:     input.Description,
		Amount:          input.Amount,
		TransactionType: input.TransactionType,
		// PaymentMethod:      input.PaymentMethod,
		TransactionDate:    input.TransactionDate,
		MerchantName:       input.MerchantName,
		Location:           input.Location,
		Notes:              input.Notes,
		IsRecurring:        input.IsRecurring != nil && *input.IsRecurring,
		RecurrenceInterval: input.RecurrenceInterval,
		RecurrenceEndDate:  input.RecurrenceEndDate,
		CreatedBy:          &userID,
	}

	// Create the transaction - the repository will populate the struct with the actual data from DB
	createdTransaction, err := s.transactionRepo.Create(ctx, transaction)
	if err != nil {
		return nil, err
	}

	// Return the transaction with data populated from the database
	return &createdTransaction, nil
}

// UpdateTransaction updates an existing transaction
func (s *transactionService) UpdateTransaction(ctx context.Context, userID uuid.UUID, transactionID uuid.UUID, input *entities.UpdateTransactionRequest) (*entities.Transaction, error) {
	// Get existing transaction and verify ownership
	existingTransaction, err := s.transactionRepo.FindByID(ctx, transactionID)
	if err != nil {
		return nil, err
	}
	if existingTransaction == nil {
		return nil, errorsutil.New(404, "transaction not found")
	}

	// Verify workspace ownership
	if existingTransaction.WorkspaceID != nil {
		workspace, err := s.workspaceRepo.FindByID(ctx, *existingTransaction.WorkspaceID)
		if err != nil {
			return nil, err
		}
		if workspace == nil || workspace.CreatedBy != userID {
			return nil, errorsutil.New(403, "access denied")
		}
	}

	// Validate required fields
	if input.Description == "" {
		return nil, errors.New("description is required")
	}
	if input.Amount == 0 {
		return nil, errors.New("amount is required")
	}
	if input.TransactionType == 0 {
		return nil, errors.New("transaction type is required")
	}
	// if input.PaymentMethod == 0 {
	// 	return nil, errors.New("payment method is required")
	// }
	if input.TransactionDate.IsZero() {
		return nil, errors.New("transaction date is required")
	}

	// Verify account ownership if account is being changed
	if input.AccountID != nil && *input.AccountID != uuid.Nil {
		account, err := s.accountRepo.FindByID(ctx, *input.AccountID)
		if err != nil {
			return nil, err
		}
		if account == nil {
			return nil, errorsutil.New(404, "account not found")
		}
		if account.UserID != userID {
			return nil, errorsutil.New(403, "access denied to account")
		}
	}

	// Update fields
	existingTransaction.AccountID = input.AccountID
	existingTransaction.CategoryID = input.CategoryID
	existingTransaction.Description = input.Description
	existingTransaction.Amount = input.Amount
	existingTransaction.TransactionType = input.TransactionType
	// existingTransaction.PaymentMethod = input.PaymentMethod
	existingTransaction.TransactionDate = input.TransactionDate
	existingTransaction.MerchantName = input.MerchantName
	existingTransaction.Location = input.Location
	existingTransaction.Notes = input.Notes
	if input.IsRecurring != nil {
		existingTransaction.IsRecurring = *input.IsRecurring
	}
	if input.RecurrenceInterval != nil {
		existingTransaction.RecurrenceInterval = *input.RecurrenceInterval
	}
	existingTransaction.RecurrenceEndDate = input.RecurrenceEndDate

	// Update the transaction - the repository will populate the struct with the actual data from DB
	updatedTransaction, err := s.transactionRepo.Update(ctx, existingTransaction)
	if err != nil {
		return nil, err
	}

	// Return the transaction with data populated from the database
	return &updatedTransaction, nil
}

// DeleteTransaction deletes a transaction
func (s *transactionService) DeleteTransaction(ctx context.Context, userID uuid.UUID, transactionID uuid.UUID) error {
	// Get existing transaction and verify ownership
	existingTransaction, err := s.transactionRepo.FindByID(ctx, transactionID)
	if err != nil {
		return err
	}
	if existingTransaction == nil {
		return errorsutil.New(404, "transaction not found")
	}

	// Verify workspace ownership
	if existingTransaction.WorkspaceID != nil {
		workspace, err := s.workspaceRepo.FindByID(ctx, *existingTransaction.WorkspaceID)
		if err != nil {
			return err
		}
		if workspace == nil || workspace.CreatedBy != userID {
			return errorsutil.New(403, "access denied")
		}
	}

	return s.transactionRepo.Delete(ctx, transactionID)
}

// GetTransactionsByWorkspace returns transactions for a workspace with pagination and filtering
func (s *transactionService) GetTransactionsByWorkspace(ctx context.Context, userID uuid.UUID, workspaceID uuid.UUID, params *entities.TransactionListParams, limit, offset int) ([]*entities.Transaction, int64, error) {
	// Verify workspace ownership
	workspace, err := s.workspaceRepo.FindByID(ctx, workspaceID)
	if err != nil {
		return nil, 0, err
	}
	if workspace == nil {
		return nil, 0, errorsutil.New(404, "workspace not found")
	}
	if workspace.CreatedBy != userID {
		return nil, 0, errorsutil.New(403, "access denied")
	}

	// Get transactions
	transactions, err := s.transactionRepo.FindByWorkspace(ctx, workspaceID, params, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	// Get total count
	totalCount, err := s.transactionRepo.CountByWorkspace(ctx, workspaceID, params)
	if err != nil {
		return nil, 0, err
	}

	return transactions, totalCount, nil
}

// GetTransactionsByAccount returns transactions for an account with pagination
func (s *transactionService) GetTransactionsByAccount(ctx context.Context, userID uuid.UUID, accountID uuid.UUID, limit, offset int) ([]*entities.Transaction, error) {
	// Verify account ownership
	account, err := s.accountRepo.FindByID(ctx, accountID)
	if err != nil {
		return nil, err
	}
	if account == nil {
		return nil, errorsutil.New(404, "account not found")
	}
	if account.UserID != userID {
		return nil, errorsutil.New(403, "access denied")
	}

	return s.transactionRepo.FindByAccountID(ctx, accountID, limit, offset)
}

// GetTransactionsByCategory returns transactions for a category with pagination
func (s *transactionService) GetTransactionsByCategory(ctx context.Context, userID uuid.UUID, categoryID uuid.UUID, limit, offset int) ([]*entities.Transaction, error) {
	// Note: Category ownership verification would require additional repository method
	// For now, we'll trust that the user has access to transactions they can see
	return s.transactionRepo.FindByCategoryID(ctx, categoryID, limit, offset)
}

// GetTransactionByID returns a transaction by ID (with user ownership verification)
func (s *transactionService) GetTransactionByID(ctx context.Context, userID uuid.UUID, transactionID uuid.UUID) (*entities.Transaction, error) {
	transaction, err := s.transactionRepo.FindByID(ctx, transactionID)
	if err != nil {
		return nil, err
	}
	if transaction == nil {
		return nil, errorsutil.New(404, "transaction not found")
	}

	// Verify workspace ownership
	if transaction.WorkspaceID != nil {
		workspace, err := s.workspaceRepo.FindByID(ctx, *transaction.WorkspaceID)
		if err != nil {
			return nil, err
		}
		if workspace == nil || workspace.CreatedBy != userID {
			return nil, errorsutil.New(403, "access denied")
		}
	}

	return transaction, nil
}
