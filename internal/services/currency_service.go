package services

import (
	"context"
	"errors"

	"github.com/vasst-id/vasst-expense-api/internal/entities"
	"github.com/vasst-id/vasst-expense-api/internal/repositories"
	errorsutil "github.com/vasst-id/vasst-expense-api/internal/utils/errors"
)

//go:generate mockgen -source=currency_service.go -package=mock -destination=mock/currency_service_mock.go
type (
	CurrencyService interface {
		CreateCurrency(ctx context.Context, input *entities.CreateCurrencyInput) (*entities.Currency, error)
		UpdateCurrency(ctx context.Context, currencyID int, input *entities.UpdateCurrencyInput) (*entities.Currency, error)
		DeleteCurrency(ctx context.Context, currencyID int) error
		GetAllCurrencies(ctx context.Context) ([]*entities.CurrencySimple, error)
		GetCurrencyByID(ctx context.Context, currencyID int) (*entities.Currency, error)
		GetCurrencyByCode(ctx context.Context, currencyCode string) (*entities.Currency, error)
	}

	currencyService struct {
		currencyRepo repositories.CurrencyRepository
	}
)

// NewCurrencyService creates a new currency service
func NewCurrencyService(currencyRepo repositories.CurrencyRepository) CurrencyService {
	return &currencyService{
		currencyRepo: currencyRepo,
	}
}

// CreateCurrency creates a new currency
func (s *currencyService) CreateCurrency(ctx context.Context, input *entities.CreateCurrencyInput) (*entities.Currency, error) {
	// Validate required fields
	if input.CurrencyCode == "" {
		return nil, errors.New("currency code is required")
	}
	if input.CurrencyName == "" {
		return nil, errors.New("currency name is required")
	}
	if input.CurrencySymbol == "" {
		return nil, errors.New("currency symbol is required")
	}

	// Check if currency with the same code already exists
	existingCurrency, err := s.currencyRepo.FindByCode(ctx, input.CurrencyCode)
	if err != nil {
		return nil, err
	}
	if existingCurrency != nil {
		return nil, errorsutil.New(409, "currency with this code already exists")
	}

	// Set default status if not provided
	status := input.CurrencyStatus
	if status == 0 {
		status = 1 // Active by default
	}

	// Set default decimal places if not provided
	decimalPlaces := input.CurrencyDecimalPlaces
	if decimalPlaces == 0 {
		decimalPlaces = 2 // Default to 2 decimal places
	}

	currency := &entities.Currency{
		CurrencyCode:          input.CurrencyCode,
		CurrencyName:          input.CurrencyName,
		CurrencySymbol:        input.CurrencySymbol,
		CurrencyDecimalPlaces: decimalPlaces,
		CurrencyStatus:        status,
	}

	// Create the currency - the repository will populate the struct with the actual data from DB
	createdCurrency, err := s.currencyRepo.Create(ctx, currency)
	if err != nil {
		return nil, err
	}

	// Return the currency with data populated from the database
	return &createdCurrency, nil
}

// UpdateCurrency updates an existing currency
func (s *currencyService) UpdateCurrency(ctx context.Context, currencyID int, input *entities.UpdateCurrencyInput) (*entities.Currency, error) {
	existingCurrency, err := s.currencyRepo.FindByID(ctx, currencyID)
	if err != nil {
		return nil, err
	}
	if existingCurrency == nil {
		return nil, errorsutil.New(404, "currency not found")
	}

	// Check for currency code uniqueness if it's being changed
	if input.CurrencyCode != "" && input.CurrencyCode != existingCurrency.CurrencyCode {
		currencyWithCode, err := s.currencyRepo.FindByCode(ctx, input.CurrencyCode)
		if err != nil {
			return nil, err
		}
		if currencyWithCode != nil && currencyWithCode.CurrencyID != currencyID {
			return nil, errorsutil.New(409, "currency code already in use")
		}
	}

	// Update fields
	if input.CurrencyCode != "" {
		existingCurrency.CurrencyCode = input.CurrencyCode
	}
	if input.CurrencyName != "" {
		existingCurrency.CurrencyName = input.CurrencyName
	}
	if input.CurrencySymbol != "" {
		existingCurrency.CurrencySymbol = input.CurrencySymbol
	}
	if input.CurrencyDecimalPlaces != 0 {
		existingCurrency.CurrencyDecimalPlaces = input.CurrencyDecimalPlaces
	}
	if input.CurrencyStatus != 0 {
		existingCurrency.CurrencyStatus = input.CurrencyStatus
	}

	// Update the currency - the repository will populate the struct with the actual data from DB
	updatedCurrency, err := s.currencyRepo.Update(ctx, existingCurrency)
	if err != nil {
		return nil, err
	}

	// Return the currency with data populated from the database
	return &updatedCurrency, nil
}

// DeleteCurrency deletes a currency
func (s *currencyService) DeleteCurrency(ctx context.Context, currencyID int) error {
	existingCurrency, err := s.currencyRepo.FindByID(ctx, currencyID)
	if err != nil {
		return err
	}
	if existingCurrency == nil {
		return errorsutil.New(404, "currency not found")
	}
	return s.currencyRepo.Delete(ctx, currencyID)
}

// GetAllCurrencies returns all active currencies in simple format
func (s *currencyService) GetAllCurrencies(ctx context.Context) ([]*entities.CurrencySimple, error) {
	return s.currencyRepo.FindAll(ctx)
}

// GetCurrencyByID returns a currency by ID
func (s *currencyService) GetCurrencyByID(ctx context.Context, currencyID int) (*entities.Currency, error) {
	currency, err := s.currencyRepo.FindByID(ctx, currencyID)
	if err != nil {
		return nil, err
	}
	if currency == nil {
		return nil, errorsutil.New(404, "currency not found")
	}
	return currency, nil
}

// GetCurrencyByCode returns a currency by code
func (s *currencyService) GetCurrencyByCode(ctx context.Context, currencyCode string) (*entities.Currency, error) {
	currency, err := s.currencyRepo.FindByCode(ctx, currencyCode)
	if err != nil {
		return nil, err
	}
	if currency == nil {
		return nil, errorsutil.New(404, "currency not found")
	}
	return currency, nil
}
