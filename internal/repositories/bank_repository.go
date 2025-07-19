package repositories

import (
	"context"
	"database/sql"

	"github.com/vasst-id/vasst-expense-api/internal/entities"
	"github.com/vasst-id/vasst-expense-api/internal/utils/postgres"
)

type (
	bankRepository struct {
		*postgres.Postgres
	}

	BankRepository interface {
		Create(ctx context.Context, bank *entities.Bank) (entities.Bank, error)
		Update(ctx context.Context, bank *entities.Bank) (entities.Bank, error)
		Delete(ctx context.Context, bankID int) error
		FindAll(ctx context.Context) ([]*entities.BankSimple, error)
		FindByID(ctx context.Context, bankID int) (*entities.Bank, error)
		FindByCode(ctx context.Context, bankCode string) (*entities.Bank, error)
	}
)

// NewBankRepository creates a new BankRepository
func NewBankRepository(pg *postgres.Postgres) BankRepository {
	return &bankRepository{pg}
}

// Create creates a new bank
func (r *bankRepository) Create(ctx context.Context, bank *entities.Bank) (entities.Bank, error) {
	query := `
		INSERT INTO "vasst_expense".banks (
			bank_name, bank_code, bank_logo_url, status, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING bank_id, bank_name, bank_code, bank_logo_url, status, created_at, updated_at
	`

	var createdBank entities.Bank
	err := r.DB.QueryRowContext(ctx, query,
		bank.BankName,
		bank.BankCode,
		bank.BankLogoURL,
		bank.Status,
	).Scan(
		&createdBank.BankID,
		&createdBank.BankName,
		&createdBank.BankCode,
		&createdBank.BankLogoURL,
		&createdBank.Status,
		&createdBank.CreatedAt,
		&createdBank.UpdatedAt,
	)

	return createdBank, err
}

// Update updates a bank
func (r *bankRepository) Update(ctx context.Context, bank *entities.Bank) (entities.Bank, error) {
	query := `
		UPDATE "vasst_expense".banks
		SET bank_name = $2,
			bank_code = $3,
			bank_logo_url = $4,
			status = $5,
			updated_at = CURRENT_TIMESTAMP
		WHERE bank_id = $1
		RETURNING bank_id, bank_name, bank_code, bank_logo_url, status, created_at, updated_at
	`

	var updatedBank entities.Bank
	err := r.DB.QueryRowContext(ctx, query,
		bank.BankID,
		bank.BankName,
		bank.BankCode,
		bank.BankLogoURL,
		bank.Status,
	).Scan(
		&updatedBank.BankID,
		&updatedBank.BankName,
		&updatedBank.BankCode,
		&updatedBank.BankLogoURL,
		&updatedBank.Status,
		&updatedBank.CreatedAt,
		&updatedBank.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return entities.Bank{}, sql.ErrNoRows
		}
		return entities.Bank{}, err
	}

	return updatedBank, nil
}

// Delete deletes a bank
func (r *bankRepository) Delete(ctx context.Context, bankID int) error {
	query := `
		DELETE FROM "vasst_expense".banks
		WHERE bank_id = $1
	`

	result, err := r.DB.ExecContext(ctx, query, bankID)
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

// FindAll returns all active banks in simple format
func (r *bankRepository) FindAll(ctx context.Context) ([]*entities.BankSimple, error) {
	query := `
		SELECT bank_id, bank_name, bank_code, bank_logo_url
		FROM "vasst_expense".banks
		WHERE status = 1
		ORDER BY bank_name ASC
	`

	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var banks []*entities.BankSimple
	for rows.Next() {
		var bank entities.BankSimple

		err := rows.Scan(
			&bank.BankID,
			&bank.BankName,
			&bank.BankCode,
			&bank.BankLogoURL,
		)
		if err != nil {
			return nil, err
		}

		banks = append(banks, &bank)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return banks, nil
}

// FindByID returns a bank by ID
func (r *bankRepository) FindByID(ctx context.Context, bankID int) (*entities.Bank, error) {
	query := `
		SELECT bank_id, bank_name, bank_code, bank_logo_url, status, created_at, updated_at
		FROM "vasst_expense".banks
		WHERE bank_id = $1
	`

	var bank entities.Bank

	err := r.DB.QueryRowContext(ctx, query, bankID).Scan(
		&bank.BankID,
		&bank.BankName,
		&bank.BankCode,
		&bank.BankLogoURL,
		&bank.Status,
		&bank.CreatedAt,
		&bank.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &bank, nil
}

// FindByCode returns a bank by bank code
func (r *bankRepository) FindByCode(ctx context.Context, bankCode string) (*entities.Bank, error) {
	query := `
		SELECT bank_id, bank_name, bank_code, bank_logo_url, status, created_at, updated_at
		FROM "vasst_expense".banks
		WHERE bank_code = $1
	`

	var bank entities.Bank

	err := r.DB.QueryRowContext(ctx, query, bankCode).Scan(
		&bank.BankID,
		&bank.BankName,
		&bank.BankCode,
		&bank.BankLogoURL,
		&bank.Status,
		&bank.CreatedAt,
		&bank.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &bank, nil
}
