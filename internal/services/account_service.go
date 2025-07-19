package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/vasst-id/vasst-expense-api/internal/entities"
	"github.com/vasst-id/vasst-expense-api/internal/repositories"
	errorsutil "github.com/vasst-id/vasst-expense-api/internal/utils/errors"
)

//go:generate mockgen -source=account_service.go -package=mock -destination=mock/account_service_mock.go
type (
	AccountService interface {
		CreateAccount(ctx context.Context, userID uuid.UUID, input *entities.CreateAccountRequest) (*entities.Account, error)
		UpdateAccount(ctx context.Context, userID uuid.UUID, accountID uuid.UUID, input *entities.UpdateAccountRequest) (*entities.Account, error)
		DeleteAccount(ctx context.Context, userID uuid.UUID, accountID uuid.UUID) error
		GetAccountsByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entities.Account, error)
		GetActiveAccountsByUserID(ctx context.Context, userID uuid.UUID) ([]*entities.Account, error)
		GetAccountByID(ctx context.Context, userID uuid.UUID, accountID uuid.UUID) (*entities.Account, error)
	}

	accountService struct {
		accountRepo repositories.AccountRepository
	}
)

// NewAccountService creates a new account service
func NewAccountService(accountRepo repositories.AccountRepository) AccountService {
	return &accountService{
		accountRepo: accountRepo,
	}
}

// CreateAccount creates a new account
func (s *accountService) CreateAccount(ctx context.Context, userID uuid.UUID, input *entities.CreateAccountRequest) (*entities.Account, error) {
	// Validate required fields
	if input.AccountName == "" {
		return nil, errors.New("account name is required")
	}
	if input.AccountType == 0 {
		return nil, errors.New("account type is required")
	}
	if input.CurrencyID == 0 {
		return nil, errors.New("currency ID is required")
	}

	// Check if account with the same name already exists for this user
	existingAccount, err := s.accountRepo.FindByNameAndUserID(ctx, userID, input.AccountName)
	if err != nil {
		return nil, err
	}
	if existingAccount != nil {
		return nil, errorsutil.New(409, "account with this name already exists")
	}

	// Create new account
	account := &entities.Account{
		AccountID:      uuid.New(),
		UserID:         userID,
		AccountName:    input.AccountName,
		AccountType:    input.AccountType,
		BankID:         input.BankID,
		AccountNumber:  input.AccountNumber,
		CurrentBalance: input.CurrentBalance,
		CreditLimit:    input.CreditLimit,
		DueDate:        input.DueDate,
		CurrencyID:     input.CurrencyID,
		IsActive:       true,
	}

	// Create the account - the repository will populate the struct with the actual data from DB
	createdAccount, err := s.accountRepo.Create(ctx, account)
	if err != nil {
		return nil, err
	}

	// Return the account with data populated from the database
	return &createdAccount, nil
}

// UpdateAccount updates an existing account
func (s *accountService) UpdateAccount(ctx context.Context, userID uuid.UUID, accountID uuid.UUID, input *entities.UpdateAccountRequest) (*entities.Account, error) {
	// Get existing account and verify ownership
	existingAccount, err := s.accountRepo.FindByID(ctx, accountID)
	if err != nil {
		return nil, err
	}
	if existingAccount == nil {
		return nil, errorsutil.New(404, "account not found")
	}
	if existingAccount.UserID != userID {
		return nil, errorsutil.New(403, "access denied")
	}

	// Validate required fields
	if input.AccountName == "" {
		return nil, errors.New("account name is required")
	}
	if input.AccountType == 0 {
		return nil, errors.New("account type is required")
	}
	if input.CurrencyID == 0 {
		return nil, errors.New("currency ID is required")
	}

	// Check for account name uniqueness if it's being changed
	if input.AccountName != existingAccount.AccountName {
		accountWithName, err := s.accountRepo.FindByNameAndUserID(ctx, userID, input.AccountName)
		if err != nil {
			return nil, err
		}
		if accountWithName != nil && accountWithName.AccountID != accountID {
			return nil, errorsutil.New(409, "account name already in use")
		}
	}

	// Update fields
	existingAccount.AccountName = input.AccountName
	existingAccount.AccountType = input.AccountType
	existingAccount.BankID = input.BankID
	existingAccount.AccountNumber = input.AccountNumber
	existingAccount.CurrentBalance = input.CurrentBalance
	existingAccount.CreditLimit = input.CreditLimit
	existingAccount.DueDate = input.DueDate
	existingAccount.CurrencyID = input.CurrencyID

	// Update the account - the repository will populate the struct with the actual data from DB
	updatedAccount, err := s.accountRepo.Update(ctx, existingAccount)
	if err != nil {
		return nil, err
	}

	// Return the account with data populated from the database
	return &updatedAccount, nil
}

// DeleteAccount deletes an account (soft delete)
func (s *accountService) DeleteAccount(ctx context.Context, userID uuid.UUID, accountID uuid.UUID) error {
	// Get existing account and verify ownership
	existingAccount, err := s.accountRepo.FindByID(ctx, accountID)
	if err != nil {
		return err
	}
	if existingAccount == nil {
		return errorsutil.New(404, "account not found")
	}
	if existingAccount.UserID != userID {
		return errorsutil.New(403, "access denied")
	}

	return s.accountRepo.Delete(ctx, accountID)
}

// GetAccountsByUserID returns accounts for a user with pagination
func (s *accountService) GetAccountsByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entities.Account, error) {
	return s.accountRepo.FindByUserID(ctx, userID, limit, offset)
}

// GetActiveAccountsByUserID returns all active accounts for a user
func (s *accountService) GetActiveAccountsByUserID(ctx context.Context, userID uuid.UUID) ([]*entities.Account, error) {
	return s.accountRepo.FindActiveByUserID(ctx, userID)
}

// GetAccountByID returns an account by ID (with user ownership verification)
func (s *accountService) GetAccountByID(ctx context.Context, userID uuid.UUID, accountID uuid.UUID) (*entities.Account, error) {
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
	return account, nil
}
