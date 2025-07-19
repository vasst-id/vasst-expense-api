package repositories

import (
	"context"
	"database/sql"

	"github.com/vasst-id/vasst-expense-api/internal/entities"
	"github.com/vasst-id/vasst-expense-api/internal/utils/postgres"
)

type (
	currencyRepository struct {
		*postgres.Postgres
	}

	CurrencyRepository interface {
		Create(ctx context.Context, currency *entities.Currency) (entities.Currency, error)
		Update(ctx context.Context, currency *entities.Currency) (entities.Currency, error)
		Delete(ctx context.Context, currencyID int) error
		FindAll(ctx context.Context) ([]*entities.CurrencySimple, error)
		FindByID(ctx context.Context, currencyID int) (*entities.Currency, error)
		FindByCode(ctx context.Context, currencyCode string) (*entities.Currency, error)
	}
)

// NewCurrencyRepository creates a new CurrencyRepository
func NewCurrencyRepository(pg *postgres.Postgres) CurrencyRepository {
	return &currencyRepository{pg}
}

// Create creates a new currency
func (r *currencyRepository) Create(ctx context.Context, currency *entities.Currency) (entities.Currency, error) {
	query := `
		INSERT INTO "vasst_expense".currency (
			currency_code, currency_name, currency_symbol, currency_decimal_places, 
			currency_status, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING currency_id, currency_code, currency_name, currency_symbol, currency_decimal_places, 
		          currency_status, created_at, updated_at
	`

	var createdCurrency entities.Currency
	err := r.DB.QueryRowContext(ctx, query,
		currency.CurrencyCode,
		currency.CurrencyName,
		currency.CurrencySymbol,
		currency.CurrencyDecimalPlaces,
		currency.CurrencyStatus,
	).Scan(
		&createdCurrency.CurrencyID,
		&createdCurrency.CurrencyCode,
		&createdCurrency.CurrencyName,
		&createdCurrency.CurrencySymbol,
		&createdCurrency.CurrencyDecimalPlaces,
		&createdCurrency.CurrencyStatus,
		&createdCurrency.CreatedAt,
		&createdCurrency.UpdatedAt,
	)

	return createdCurrency, err
}

// Update updates a currency
func (r *currencyRepository) Update(ctx context.Context, currency *entities.Currency) (entities.Currency, error) {
	query := `
		UPDATE "vasst_expense".currency
		SET currency_code = $2,
			currency_name = $3,
			currency_symbol = $4,
			currency_decimal_places = $5,
			currency_status = $6,
			updated_at = CURRENT_TIMESTAMP
		WHERE currency_id = $1
		RETURNING currency_id, currency_code, currency_name, currency_symbol, currency_decimal_places, 
		          currency_status, created_at, updated_at
	`

	var updatedCurrency entities.Currency
	err := r.DB.QueryRowContext(ctx, query,
		currency.CurrencyID,
		currency.CurrencyCode,
		currency.CurrencyName,
		currency.CurrencySymbol,
		currency.CurrencyDecimalPlaces,
		currency.CurrencyStatus,
	).Scan(
		&updatedCurrency.CurrencyID,
		&updatedCurrency.CurrencyCode,
		&updatedCurrency.CurrencyName,
		&updatedCurrency.CurrencySymbol,
		&updatedCurrency.CurrencyDecimalPlaces,
		&updatedCurrency.CurrencyStatus,
		&updatedCurrency.CreatedAt,
		&updatedCurrency.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return entities.Currency{}, sql.ErrNoRows
		}
		return entities.Currency{}, err
	}

	return updatedCurrency, nil
}

// Delete deletes a currency
func (r *currencyRepository) Delete(ctx context.Context, currencyID int) error {
	query := `
		DELETE FROM "vasst_expense".currency
		WHERE currency_id = $1
	`

	result, err := r.DB.ExecContext(ctx, query, currencyID)
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

// FindAll returns all active currencies in simple format
func (r *currencyRepository) FindAll(ctx context.Context) ([]*entities.CurrencySimple, error) {
	query := `
		SELECT currency_id, currency_code, currency_name, currency_symbol, currency_decimal_places
		FROM "vasst_expense".currency
		WHERE currency_status = 1
		ORDER BY currency_name ASC
	`

	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var currencies []*entities.CurrencySimple
	for rows.Next() {
		var currency entities.CurrencySimple

		err := rows.Scan(
			&currency.CurrencyID,
			&currency.CurrencyCode,
			&currency.CurrencyName,
			&currency.CurrencySymbol,
			&currency.CurrencyDecimalPlaces,
		)
		if err != nil {
			return nil, err
		}

		currencies = append(currencies, &currency)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return currencies, nil
}

// FindByID returns a currency by ID
func (r *currencyRepository) FindByID(ctx context.Context, currencyID int) (*entities.Currency, error) {
	query := `
		SELECT currency_id, currency_code, currency_name, currency_symbol, 
			   currency_decimal_places, currency_status, created_at, updated_at
		FROM "vasst_expense".currency
		WHERE currency_id = $1
	`

	var currency entities.Currency

	err := r.DB.QueryRowContext(ctx, query, currencyID).Scan(
		&currency.CurrencyID,
		&currency.CurrencyCode,
		&currency.CurrencyName,
		&currency.CurrencySymbol,
		&currency.CurrencyDecimalPlaces,
		&currency.CurrencyStatus,
		&currency.CreatedAt,
		&currency.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &currency, nil
}

// FindByCode returns a currency by currency code
func (r *currencyRepository) FindByCode(ctx context.Context, currencyCode string) (*entities.Currency, error) {
	query := `
		SELECT currency_id, currency_code, currency_name, currency_symbol, 
			   currency_decimal_places, currency_status, created_at, updated_at
		FROM "vasst_expense".currency
		WHERE currency_code = $1
	`

	var currency entities.Currency

	err := r.DB.QueryRowContext(ctx, query, currencyCode).Scan(
		&currency.CurrencyID,
		&currency.CurrencyCode,
		&currency.CurrencyName,
		&currency.CurrencySymbol,
		&currency.CurrencyDecimalPlaces,
		&currency.CurrencyStatus,
		&currency.CreatedAt,
		&currency.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &currency, nil
}
