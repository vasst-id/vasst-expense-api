package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/vasst-id/vasst-expense-api/internal/entities"
	"github.com/vasst-id/vasst-expense-api/internal/repositories"
	errorsutil "github.com/vasst-id/vasst-expense-api/internal/utils/errors"
)

//go:generate mockgen -source=budget_service.go -package=mock -destination=mock/budget_service_mock.go
type (
	BudgetService interface {
		CreateBudget(ctx context.Context, workspaceID uuid.UUID, userID uuid.UUID, input *entities.CreateBudgetRequest) (*entities.Budget, error)
		UpdateBudget(ctx context.Context, budgetID uuid.UUID, workspaceID uuid.UUID, input *entities.UpdateBudgetRequest) (*entities.Budget, error)
		DeleteBudget(ctx context.Context, budgetID uuid.UUID, workspaceID uuid.UUID) error
		GetAllBudgets(ctx context.Context, workspaceID uuid.UUID, limit, offset int) ([]*entities.BudgetSimple, error)
		GetBudgetByID(ctx context.Context, budgetID uuid.UUID, workspaceID uuid.UUID) (*entities.BudgetSimple, error)
	}

	budgetService struct {
		budgetRepo repositories.BudgetRepository
	}
)

// NewBudgetService creates a new budget service
func NewBudgetService(budgetRepo repositories.BudgetRepository) BudgetService {
	return &budgetService{
		budgetRepo: budgetRepo,
	}
}

// CreateBudget creates a new budget
func (s *budgetService) CreateBudget(ctx context.Context, workspaceID uuid.UUID, userID uuid.UUID, input *entities.CreateBudgetRequest) (*entities.Budget, error) {
	// Validate required fields
	if input.Name == "" {
		return nil, errors.New("budget name is required")
	}
	if input.BudgetedAmount <= 0 {
		return nil, errors.New("budgeted amount must be greater than 0")
	}
	if input.PeriodType < 1 || input.PeriodType > 4 {
		return nil, errors.New("invalid period type")
	}
	if input.PeriodStart.IsZero() {
		return nil, errors.New("period start is required")
	}
	if input.PeriodEnd.IsZero() {
		return nil, errors.New("period end is required")
	}
	if input.PeriodEnd.Before(input.PeriodStart) {
		return nil, errors.New("period end must be after period start")
	}
	if input.UserCategoryID == uuid.Nil {
		return nil, errors.New("user category ID is required")
	}

	budget := &entities.Budget{
		BudgetID:       uuid.New(),
		WorkspaceID:    workspaceID,
		UserCategoryID: input.UserCategoryID,
		Name:           input.Name,
		BudgetedAmount: input.BudgetedAmount,
		PeriodType:     input.PeriodType,
		PeriodStart:    input.PeriodStart,
		PeriodEnd:      input.PeriodEnd,
		SpentAmount:    0, // Default to 0
		IsActive:       true,
		CreatedBy:      userID,
	}

	// Create the budget - the repository will populate the struct with the actual data from DB
	createdBudget, err := s.budgetRepo.Create(ctx, budget)
	if err != nil {
		return nil, err
	}

	// Return the budget with data populated from the database
	return &createdBudget, nil
}

// UpdateBudget updates an existing budget
func (s *budgetService) UpdateBudget(ctx context.Context, budgetID uuid.UUID, workspaceID uuid.UUID, input *entities.UpdateBudgetRequest) (*entities.Budget, error) {
	// Validate required fields
	if input.Name == "" {
		return nil, errors.New("budget name is required")
	}
	if input.BudgetedAmount <= 0 {
		return nil, errors.New("budgeted amount must be greater than 0")
	}
	if input.PeriodType < 1 || input.PeriodType > 4 {
		return nil, errors.New("invalid period type")
	}
	if input.PeriodStart.IsZero() {
		return nil, errors.New("period start is required")
	}
	if input.PeriodEnd.IsZero() {
		return nil, errors.New("period end is required")
	}
	if input.PeriodEnd.Before(input.PeriodStart) {
		return nil, errors.New("period end must be after period start")
	}
	if input.UserCategoryID == uuid.Nil {
		return nil, errors.New("user category ID is required")
	}

	// Check if budget exists and belongs to workspace
	existingBudget, err := s.budgetRepo.FindByID(ctx, budgetID)
	if err != nil {
		return nil, err
	}
	if existingBudget == nil {
		return nil, errorsutil.New(404, "budget not found")
	}
	if existingBudget.WorkspaceID != workspaceID {
		return nil, errorsutil.New(404, "budget not found")
	}

	// Update fields
	existingBudget.UserCategoryID = input.UserCategoryID
	existingBudget.Name = input.Name
	existingBudget.BudgetedAmount = input.BudgetedAmount
	existingBudget.PeriodType = input.PeriodType
	existingBudget.PeriodStart = input.PeriodStart
	existingBudget.PeriodEnd = input.PeriodEnd
	existingBudget.SpentAmount = input.SpentAmount
	existingBudget.IsActive = input.IsActive

	// Update the budget - the repository will populate the struct with the actual data from DB
	updatedBudget, err := s.budgetRepo.Update(ctx, existingBudget)
	if err != nil {
		return nil, err
	}

	// Return the budget with data populated from the database
	return &updatedBudget, nil
}

// DeleteBudget deletes a budget
func (s *budgetService) DeleteBudget(ctx context.Context, budgetID uuid.UUID, workspaceID uuid.UUID) error {
	// Check if budget exists and belongs to workspace
	existingBudget, err := s.budgetRepo.FindByID(ctx, budgetID)
	if err != nil {
		return err
	}
	if existingBudget == nil {
		return errorsutil.New(404, "budget not found")
	}
	if existingBudget.WorkspaceID != workspaceID {
		return errorsutil.New(404, "budget not found")
	}

	return s.budgetRepo.Delete(ctx, budgetID)
}

// GetAllBudgets returns all budgets for a workspace with pagination
func (s *budgetService) GetAllBudgets(ctx context.Context, workspaceID uuid.UUID, limit, offset int) ([]*entities.BudgetSimple, error) {
	return s.budgetRepo.FindByWorkspace(ctx, workspaceID, limit, offset)
}

// GetBudgetByID returns a budget by ID within a workspace
func (s *budgetService) GetBudgetByID(ctx context.Context, budgetID uuid.UUID, workspaceID uuid.UUID) (*entities.BudgetSimple, error) {
	budget, err := s.budgetRepo.FindByIDWithWorkspace(ctx, budgetID, workspaceID)
	if err != nil {
		return nil, err
	}
	if budget == nil {
		return nil, errorsutil.New(404, "budget not found")
	}
	return budget, nil
}
