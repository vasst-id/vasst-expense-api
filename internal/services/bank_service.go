package services

import (
	"context"
	"errors"

	"github.com/vasst-id/vasst-expense-api/internal/entities"
	"github.com/vasst-id/vasst-expense-api/internal/repositories"
	errorsutil "github.com/vasst-id/vasst-expense-api/internal/utils/errors"
)

//go:generate mockgen -source=bank_service.go -package=mock -destination=mock/bank_service_mock.go
type (
	BankService interface {
		CreateBank(ctx context.Context, input *entities.CreateBankInput) (*entities.Bank, error)
		UpdateBank(ctx context.Context, bankID int, input *entities.UpdateBankInput) (*entities.Bank, error)
		DeleteBank(ctx context.Context, bankID int) error
		GetAllBanks(ctx context.Context) ([]*entities.BankSimple, error)
		GetBankByID(ctx context.Context, bankID int) (*entities.Bank, error)
		GetBankByCode(ctx context.Context, bankCode string) (*entities.Bank, error)
	}

	bankService struct {
		bankRepo repositories.BankRepository
	}
)

// NewBankService creates a new bank service
func NewBankService(bankRepo repositories.BankRepository) BankService {
	return &bankService{
		bankRepo: bankRepo,
	}
}

// CreateBank creates a new bank
func (s *bankService) CreateBank(ctx context.Context, input *entities.CreateBankInput) (*entities.Bank, error) {
	// Validate required fields
	if input.BankName == "" {
		return nil, errors.New("bank name is required")
	}
	if input.BankCode == "" {
		return nil, errors.New("bank code is required")
	}

	// Check if bank with the same code already exists
	existingBank, err := s.bankRepo.FindByCode(ctx, input.BankCode)
	if err != nil {
		return nil, err
	}
	if existingBank != nil {
		return nil, errorsutil.New(409, "bank with this code already exists")
	}

	// Set default status if not provided
	status := input.Status
	if status == 0 {
		status = 1 // Active by default
	}

	bank := &entities.Bank{
		BankName:    input.BankName,
		BankCode:    input.BankCode,
		BankLogoURL: input.BankLogoURL,
		Status:      status,
	}

	// Create the bank - the repository will populate the struct with the actual data from DB
	createdBank, err := s.bankRepo.Create(ctx, bank)
	if err != nil {
		return nil, err
	}

	// Return the bank with data populated from the database
	return &createdBank, nil
}

// UpdateBank updates an existing bank
func (s *bankService) UpdateBank(ctx context.Context, bankID int, input *entities.UpdateBankInput) (*entities.Bank, error) {
	existingBank, err := s.bankRepo.FindByID(ctx, bankID)
	if err != nil {
		return nil, err
	}
	if existingBank == nil {
		return nil, errorsutil.New(404, "bank not found")
	}

	// Check for bank code uniqueness if it's being changed
	if input.BankCode != "" && input.BankCode != existingBank.BankCode {
		bankWithCode, err := s.bankRepo.FindByCode(ctx, input.BankCode)
		if err != nil {
			return nil, err
		}
		if bankWithCode != nil && bankWithCode.BankID != bankID {
			return nil, errorsutil.New(409, "bank code already in use")
		}
	}

	// Update fields
	if input.BankName != "" {
		existingBank.BankName = input.BankName
	}
	if input.BankCode != "" {
		existingBank.BankCode = input.BankCode
	}
	if input.BankLogoURL != "" {
		existingBank.BankLogoURL = input.BankLogoURL
	}
	if input.Status != 0 {
		existingBank.Status = input.Status
	}

	// Update the bank - the repository will populate the struct with the actual data from DB
	updatedBank, err := s.bankRepo.Update(ctx, existingBank)
	if err != nil {
		return nil, err
	}

	// Return the bank with data populated from the database
	return &updatedBank, nil
}

// DeleteBank deletes a bank
func (s *bankService) DeleteBank(ctx context.Context, bankID int) error {
	existingBank, err := s.bankRepo.FindByID(ctx, bankID)
	if err != nil {
		return err
	}
	if existingBank == nil {
		return errorsutil.New(404, "bank not found")
	}
	return s.bankRepo.Delete(ctx, bankID)
}

// GetAllBanks returns all active banks in simple format
func (s *bankService) GetAllBanks(ctx context.Context) ([]*entities.BankSimple, error) {
	return s.bankRepo.FindAll(ctx)
}

// GetBankByID returns a bank by ID
func (s *bankService) GetBankByID(ctx context.Context, bankID int) (*entities.Bank, error) {
	bank, err := s.bankRepo.FindByID(ctx, bankID)
	if err != nil {
		return nil, err
	}
	if bank == nil {
		return nil, errorsutil.New(404, "bank not found")
	}
	return bank, nil
}

// GetBankByCode returns a bank by code
func (s *bankService) GetBankByCode(ctx context.Context, bankCode string) (*entities.Bank, error) {
	bank, err := s.bankRepo.FindByCode(ctx, bankCode)
	if err != nil {
		return nil, err
	}
	if bank == nil {
		return nil, errorsutil.New(404, "bank not found")
	}
	return bank, nil
}
