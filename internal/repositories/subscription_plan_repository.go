package repositories

import (
	"context"
	"database/sql"

	"github.com/vasst-id/vasst-expense-api/internal/entities"
	"github.com/vasst-id/vasst-expense-api/internal/utils/postgres"
)

type (
	subscriptionPlanRepository struct {
		*postgres.Postgres
	}

	SubscriptionPlanRepository interface {
		Create(ctx context.Context, plan *entities.SubscriptionPlan) (entities.SubscriptionPlan, error)
		Update(ctx context.Context, plan *entities.SubscriptionPlan) (entities.SubscriptionPlan, error)
		Delete(ctx context.Context, subscriptionPlanID int) error
		FindAll(ctx context.Context) ([]*entities.SubscriptionPlanSimple, error)
		FindByID(ctx context.Context, subscriptionPlanID int) (*entities.SubscriptionPlan, error)
	}
)

// NewSubscriptionPlanRepository creates a new SubscriptionPlanRepository
func NewSubscriptionPlanRepository(pg *postgres.Postgres) SubscriptionPlanRepository {
	return &subscriptionPlanRepository{pg}
}

// Create creates a new plan
func (r *subscriptionPlanRepository) Create(ctx context.Context, plan *entities.SubscriptionPlan) (entities.SubscriptionPlan, error) {
	query := `
		INSERT INTO "vasst_expense".subscription_plan (
			subscription_plan_name, subscription_plan_description, subscription_plan_price, subscription_plan_currency_id, subscription_plan_features, subscription_plan_status, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING subscription_plan_id, subscription_plan_name, subscription_plan_description, subscription_plan_price, subscription_plan_currency_id, subscription_plan_features, subscription_plan_status, created_at, updated_at
	`

	var createdPlan entities.SubscriptionPlan
	err := r.DB.QueryRowContext(ctx, query,
		plan.SubscriptionPlanName,
		plan.SubscriptionPlanDescription,
		plan.SubscriptionPlanPrice,
		plan.SubscriptionPlanCurrencyID,
		plan.SubscriptionPlanFeatures,
		plan.SubscriptionPlanStatus,
	).Scan(
		&createdPlan.SubscriptionPlanID,
		&createdPlan.SubscriptionPlanName,
		&createdPlan.SubscriptionPlanDescription,
		&createdPlan.SubscriptionPlanPrice,
		&createdPlan.SubscriptionPlanCurrencyID,
		&createdPlan.SubscriptionPlanFeatures,
		&createdPlan.SubscriptionPlanStatus,
		&createdPlan.CreatedAt,
		&createdPlan.UpdatedAt,
	)

	return createdPlan, err
}

// Update updates a plan
func (r *subscriptionPlanRepository) Update(ctx context.Context, plan *entities.SubscriptionPlan) (entities.SubscriptionPlan, error) {
	query := `
		UPDATE "vasst_expense".subscription_plan
		SET subscription_plan_name = $2,
			subscription_plan_description = $3,
			subscription_plan_price = $4,
			subscription_plan_currency_id = $5,
			subscription_plan_features = $6,
			subscription_plan_status = $7,
			updated_at = CURRENT_TIMESTAMP
		WHERE subscription_plan_id = $1
		RETURNING subscription_plan_id, subscription_plan_name, subscription_plan_description, subscription_plan_price, subscription_plan_currency_id, subscription_plan_features, subscription_plan_status, created_at, updated_at
	`

	var updatedPlan entities.SubscriptionPlan
	err := r.DB.QueryRowContext(ctx, query,
		plan.SubscriptionPlanID,
		plan.SubscriptionPlanName,
		plan.SubscriptionPlanDescription,
		plan.SubscriptionPlanPrice,
		plan.SubscriptionPlanCurrencyID,
		plan.SubscriptionPlanFeatures,
		plan.SubscriptionPlanStatus,
	).Scan(
		&updatedPlan.SubscriptionPlanID,
		&updatedPlan.SubscriptionPlanName,
		&updatedPlan.SubscriptionPlanDescription,
		&updatedPlan.SubscriptionPlanPrice,
		&updatedPlan.SubscriptionPlanCurrencyID,
		&updatedPlan.SubscriptionPlanFeatures,
		&updatedPlan.SubscriptionPlanStatus,
		&updatedPlan.CreatedAt,
		&updatedPlan.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return entities.SubscriptionPlan{}, sql.ErrNoRows
		}
		return entities.SubscriptionPlan{}, err
	}

	return updatedPlan, nil
}

// Delete deletes a plan
func (r *subscriptionPlanRepository) Delete(ctx context.Context, subscriptionPlanID int) error {
	query := `
		DELETE FROM "vasst_expense".subscription_plan
		WHERE subscription_plan_id = $1
	`

	result, err := r.DB.ExecContext(ctx, query, subscriptionPlanID)
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

// FindAll returns all active plans in simple format
func (r *subscriptionPlanRepository) FindAll(ctx context.Context) ([]*entities.SubscriptionPlanSimple, error) {
	query := `
		SELECT subscription_plan_id, subscription_plan_name, subscription_plan_description, subscription_plan_price, subscription_plan_currency_id, subscription_plan_features, subscription_plan_status
		FROM "vasst_expense".subscription_plan
		WHERE subscription_plan_status = 1
		ORDER BY subscription_plan_name ASC
	`

	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var plans []*entities.SubscriptionPlanSimple
	for rows.Next() {
		var plan entities.SubscriptionPlanSimple

		err := rows.Scan(
			&plan.SubscriptionPlanID,
			&plan.SubscriptionPlanName,
			&plan.SubscriptionPlanDescription,
			&plan.SubscriptionPlanPrice,
			&plan.SubscriptionPlanCurrencyID,
			&plan.SubscriptionPlanFeatures,
			&plan.SubscriptionPlanStatus,
		)
		if err != nil {
			return nil, err
		}

		plans = append(plans, &plan)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return plans, nil
}

// FindByID returns a plan by ID
func (r *subscriptionPlanRepository) FindByID(ctx context.Context, planID int) (*entities.SubscriptionPlan, error) {
	query := `
		SELECT subscription_plan_id, subscription_plan_name, subscription_plan_description, subscription_plan_price, subscription_plan_currency_id, subscription_plan_features, subscription_plan_status, created_at, updated_at
		FROM "vasst_expense".subscription_plan
		WHERE subscription_plan_id = $1
	`

	var plan entities.SubscriptionPlan

	err := r.DB.QueryRowContext(ctx, query, planID).Scan(
		&plan.SubscriptionPlanID,
		&plan.SubscriptionPlanName,
		&plan.SubscriptionPlanDescription,
		&plan.SubscriptionPlanPrice,
		&plan.SubscriptionPlanCurrencyID,
		&plan.SubscriptionPlanFeatures,
		&plan.SubscriptionPlanStatus,
		&plan.CreatedAt,
		&plan.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &plan, nil
}
