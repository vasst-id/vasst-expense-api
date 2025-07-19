package repositories

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/vasst-id/vasst-expense-api/internal/entities"
	"github.com/vasst-id/vasst-expense-api/internal/utils/postgres"
)

type (
	budgetRepository struct {
		*postgres.Postgres
	}

	BudgetRepository interface {
		Create(ctx context.Context, budget *entities.Budget) (entities.Budget, error)
		Update(ctx context.Context, budget *entities.Budget) (entities.Budget, error)
		Delete(ctx context.Context, budgetID uuid.UUID) error
		FindByID(ctx context.Context, budgetID uuid.UUID) (*entities.Budget, error)
		FindByIDWithWorkspace(ctx context.Context, budgetID uuid.UUID, workspaceID uuid.UUID) (*entities.BudgetSimple, error)
		FindByWorkspace(ctx context.Context, workspaceID uuid.UUID, limit, offset int) ([]*entities.BudgetSimple, error)
	}
)

// NewBudgetRepository creates a new BudgetRepository
func NewBudgetRepository(pg *postgres.Postgres) BudgetRepository {
	return &budgetRepository{pg}
}

// Create creates a new budget
func (r *budgetRepository) Create(ctx context.Context, budget *entities.Budget) (entities.Budget, error) {
	query := `
		INSERT INTO "vasst_expense".budgets (
			budget_id, workspace_id, user_category_id, name, budgeted_amount, 
			period_type, period_start, period_end, spent_amount, is_active, 
			created_by, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING budget_id, workspace_id, user_category_id, name, budgeted_amount, 
		          period_type, period_start, period_end, spent_amount, is_active, 
		          created_by, created_at, updated_at
	`

	var createdBudget entities.Budget

	err := r.DB.QueryRowContext(ctx, query,
		budget.BudgetID,
		budget.WorkspaceID,
		budget.UserCategoryID,
		budget.Name,
		budget.BudgetedAmount,
		budget.PeriodType,
		budget.PeriodStart,
		budget.PeriodEnd,
		budget.SpentAmount,
		budget.IsActive,
		budget.CreatedBy,
	).Scan(
		&createdBudget.BudgetID,
		&createdBudget.WorkspaceID,
		&createdBudget.UserCategoryID,
		&createdBudget.Name,
		&createdBudget.BudgetedAmount,
		&createdBudget.PeriodType,
		&createdBudget.PeriodStart,
		&createdBudget.PeriodEnd,
		&createdBudget.SpentAmount,
		&createdBudget.IsActive,
		&createdBudget.CreatedBy,
		&createdBudget.CreatedAt,
		&createdBudget.UpdatedAt,
	)

	return createdBudget, err
}

// Update updates a budget
func (r *budgetRepository) Update(ctx context.Context, budget *entities.Budget) (entities.Budget, error) {
	query := `
		UPDATE "vasst_expense".budgets
		SET user_category_id = $2,
			name = $3,
			budgeted_amount = $4,
			period_type = $5,
			period_start = $6,
			period_end = $7,
			spent_amount = $8,
			is_active = $9,
			updated_at = CURRENT_TIMESTAMP
		WHERE budget_id = $1
		RETURNING budget_id, workspace_id, user_category_id, name, budgeted_amount, 
		          period_type, period_start, period_end, spent_amount, is_active, 
		          created_by, created_at, updated_at
	`

	var updatedBudget entities.Budget
	err := r.DB.QueryRowContext(ctx, query,
		budget.BudgetID,
		budget.UserCategoryID,
		budget.Name,
		budget.BudgetedAmount,
		budget.PeriodType,
		budget.PeriodStart,
		budget.PeriodEnd,
		budget.SpentAmount,
		budget.IsActive,
	).Scan(
		&updatedBudget.BudgetID,
		&updatedBudget.WorkspaceID,
		&updatedBudget.UserCategoryID,
		&updatedBudget.Name,
		&updatedBudget.BudgetedAmount,
		&updatedBudget.PeriodType,
		&updatedBudget.PeriodStart,
		&updatedBudget.PeriodEnd,
		&updatedBudget.SpentAmount,
		&updatedBudget.IsActive,
		&updatedBudget.CreatedBy,
		&updatedBudget.CreatedAt,
		&updatedBudget.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return entities.Budget{}, sql.ErrNoRows // Budget not found
		}
		return entities.Budget{}, err
	}

	return updatedBudget, nil
}

// Delete deletes a budget
func (r *budgetRepository) Delete(ctx context.Context, budgetID uuid.UUID) error {
	query := `
		DELETE FROM "vasst_expense".budgets
		WHERE budget_id = $1
	`

	result, err := r.DB.ExecContext(ctx, query, budgetID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// FindByID finds a budget by ID (for internal operations)
func (r *budgetRepository) FindByID(ctx context.Context, budgetID uuid.UUID) (*entities.Budget, error) {
	query := `
		SELECT budget_id, workspace_id, user_category_id, name, budgeted_amount, 
			   period_type, period_start, period_end, spent_amount, is_active, 
			   created_by, created_at, updated_at
		FROM "vasst_expense".budgets
		WHERE budget_id = $1
	`

	var budget entities.Budget
	err := r.DB.QueryRowContext(ctx, query, budgetID).Scan(
		&budget.BudgetID,
		&budget.WorkspaceID,
		&budget.UserCategoryID,
		&budget.Name,
		&budget.BudgetedAmount,
		&budget.PeriodType,
		&budget.PeriodStart,
		&budget.PeriodEnd,
		&budget.SpentAmount,
		&budget.IsActive,
		&budget.CreatedBy,
		&budget.CreatedAt,
		&budget.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &budget, nil
}

// FindByIDWithWorkspace finds a budget by ID within a specific workspace (returns BudgetSimple)
func (r *budgetRepository) FindByIDWithWorkspace(ctx context.Context, budgetID uuid.UUID, workspaceID uuid.UUID) (*entities.BudgetSimple, error) {
	query := `
		SELECT 
			b.budget_id,
			COALESCE(uc.name, 'Uncategorized') as user_category_name,
			b.name,
			b.budgeted_amount,
			CASE 
				WHEN b.period_type = 1 THEN 'Weekly'
				WHEN b.period_type = 2 THEN 'Monthly'
				WHEN b.period_type = 3 THEN 'Yearly'
				WHEN b.period_type = 4 THEN 'Event'
				ELSE 'Unknown'
			END as period_type_label,
			b.period_start,
			b.period_end,
			b.spent_amount,
			(b.budgeted_amount - b.spent_amount) as remaining_amount,
			CASE 
				WHEN b.budgeted_amount > 0 THEN (b.spent_amount / b.budgeted_amount * 100)
				ELSE 0
			END as percentage_used,
			GREATEST(0, (b.period_end - CURRENT_DATE)) as days_remaining,
			(b.spent_amount > b.budgeted_amount) as is_overspent
		FROM "vasst_expense".budgets b
		LEFT JOIN "vasst_expense".user_categories uc ON b.user_category_id = uc.user_category_id
		WHERE b.budget_id = $1 AND b.workspace_id = $2 AND b.is_active = true
	`

	var budget entities.BudgetSimple
	err := r.DB.QueryRowContext(ctx, query, budgetID, workspaceID).Scan(
		&budget.BudgetID,
		&budget.UserCategoryName,
		&budget.Name,
		&budget.BudgetedAmount,
		&budget.PeriodTypeLabel,
		&budget.PeriodStart,
		&budget.PeriodEnd,
		&budget.SpentAmount,
		&budget.RemainingAmount,
		&budget.PercentageUsed,
		&budget.DaysRemaining,
		&budget.IsOverspent,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &budget, nil
}

// FindByWorkspace finds all budgets by workspace with pagination (returns BudgetSimple)
func (r *budgetRepository) FindByWorkspace(ctx context.Context, workspaceID uuid.UUID, limit, offset int) ([]*entities.BudgetSimple, error) {
	query := `
		SELECT 
			b.budget_id,
			COALESCE(uc.name, 'Uncategorized') as user_category_name,
			b.name,
			b.budgeted_amount,
			CASE 
				WHEN b.period_type = 1 THEN 'Weekly'
				WHEN b.period_type = 2 THEN 'Monthly'
				WHEN b.period_type = 3 THEN 'Yearly'
				WHEN b.period_type = 4 THEN 'Event'
				ELSE 'Unknown'
			END as period_type_label,
			b.period_start,
			b.period_end,
			b.spent_amount,
			(b.budgeted_amount - b.spent_amount) as remaining_amount,
			CASE 
				WHEN b.budgeted_amount > 0 THEN (b.spent_amount / b.budgeted_amount * 100)
				ELSE 0
			END as percentage_used,
			GREATEST(0, (b.period_end - CURRENT_DATE)) as days_remaining,
			(b.spent_amount > b.budgeted_amount) as is_overspent
		FROM "vasst_expense".budgets b
		LEFT JOIN "vasst_expense".user_categories uc ON b.user_category_id = uc.user_category_id
		WHERE b.workspace_id = $1 AND b.is_active = true
		ORDER BY b.created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.DB.QueryContext(ctx, query, workspaceID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var budgets []*entities.BudgetSimple
	for rows.Next() {
		var budget entities.BudgetSimple
		err := rows.Scan(
			&budget.BudgetID,
			&budget.UserCategoryName,
			&budget.Name,
			&budget.BudgetedAmount,
			&budget.PeriodTypeLabel,
			&budget.PeriodStart,
			&budget.PeriodEnd,
			&budget.SpentAmount,
			&budget.RemainingAmount,
			&budget.PercentageUsed,
			&budget.DaysRemaining,
			&budget.IsOverspent,
		)
		if err != nil {
			return nil, err
		}

		budgets = append(budgets, &budget)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return budgets, nil
}
