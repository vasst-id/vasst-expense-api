package repositories

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/vasst-id/vasst-expense-api/internal/entities"
	"github.com/vasst-id/vasst-expense-api/internal/utils/postgres"
)

type (
	accountRepository struct {
		*postgres.Postgres
	}

	AccountRepository interface {
		Create(ctx context.Context, account *entities.Account) (entities.Account, error)
		Update(ctx context.Context, account *entities.Account) (entities.Account, error)
		Delete(ctx context.Context, accountID uuid.UUID) error
		FindByID(ctx context.Context, accountID uuid.UUID) (*entities.Account, error)
		FindByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entities.Account, error)
		FindActiveByUserID(ctx context.Context, userID uuid.UUID) ([]*entities.Account, error)
		FindByNameAndUserID(ctx context.Context, userID uuid.UUID, accountName string) (*entities.Account, error)
	}
)

// NewAccountRepository creates a new AccountRepository
func NewAccountRepository(pg *postgres.Postgres) AccountRepository {
	return &accountRepository{pg}
}

// Create creates a new account
func (r *accountRepository) Create(ctx context.Context, account *entities.Account) (entities.Account, error) {
	query := `
		INSERT INTO "vasst_expense".accounts 
		(account_id, user_id, account_name, account_type, bank_id, account_number, 
		 current_balance, credit_limit, due_date, currency_id, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING account_id, user_id, account_name, account_type, bank_id, account_number,
		          current_balance, credit_limit, due_date, currency_id, is_active, created_at, updated_at
	`

	var createdAccount entities.Account
	err := r.DB.QueryRowContext(ctx, query,
		account.AccountID, account.UserID, account.AccountName, account.AccountType,
		account.BankID, account.AccountNumber, account.CurrentBalance, account.CreditLimit,
		account.DueDate, account.CurrencyID, account.IsActive,
	).Scan(
		&createdAccount.AccountID, &createdAccount.UserID, &createdAccount.AccountName, &createdAccount.AccountType,
		&createdAccount.BankID, &createdAccount.AccountNumber, &createdAccount.CurrentBalance, &createdAccount.CreditLimit,
		&createdAccount.DueDate, &createdAccount.CurrencyID, &createdAccount.IsActive, &createdAccount.CreatedAt, &createdAccount.UpdatedAt,
	)

	return createdAccount, err
}

// Update updates an account
func (r *accountRepository) Update(ctx context.Context, account *entities.Account) (entities.Account, error) {
	query := `
		UPDATE "vasst_expense".accounts 
		SET account_name = $2, account_type = $3, bank_id = $4, account_number = $5,
		    current_balance = $6, credit_limit = $7, due_date = $8, currency_id = $9,
		    is_active = $10, updated_at = CURRENT_TIMESTAMP
		WHERE account_id = $1
		RETURNING account_id, user_id, account_name, account_type, bank_id, account_number,
		          current_balance, credit_limit, due_date, currency_id, is_active, created_at, updated_at
	`

	var updatedAccount entities.Account
	err := r.DB.QueryRowContext(ctx, query,
		account.AccountID, account.AccountName, account.AccountType, account.BankID,
		account.AccountNumber, account.CurrentBalance, account.CreditLimit, account.DueDate,
		account.CurrencyID, account.IsActive,
	).Scan(
		&updatedAccount.AccountID, &updatedAccount.UserID, &updatedAccount.AccountName, &updatedAccount.AccountType,
		&updatedAccount.BankID, &updatedAccount.AccountNumber, &updatedAccount.CurrentBalance, &updatedAccount.CreditLimit,
		&updatedAccount.DueDate, &updatedAccount.CurrencyID, &updatedAccount.IsActive, &updatedAccount.CreatedAt, &updatedAccount.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return entities.Account{}, sql.ErrNoRows // Account not found
		}
		return entities.Account{}, err
	}

	return updatedAccount, nil
}

// Delete soft deletes an account (sets is_active to false)
func (r *accountRepository) Delete(ctx context.Context, accountID uuid.UUID) error {
	query := `
		UPDATE "vasst_expense".accounts 
		SET is_active = false, updated_at = CURRENT_TIMESTAMP
		WHERE account_id = $1
	`

	result, err := r.DB.ExecContext(ctx, query, accountID)
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

// FindByID finds an account by ID
func (r *accountRepository) FindByID(ctx context.Context, accountID uuid.UUID) (*entities.Account, error) {
	query := `
		SELECT account_id, user_id, account_name, account_type, bank_id, account_number,
		       current_balance, credit_limit, due_date, currency_id, is_active, created_at, updated_at
		FROM "vasst_expense".accounts 
		WHERE account_id = $1 AND is_active = true
	`

	var account entities.Account
	err := r.DB.QueryRowContext(ctx, query, accountID).Scan(
		&account.AccountID, &account.UserID, &account.AccountName, &account.AccountType,
		&account.BankID, &account.AccountNumber, &account.CurrentBalance, &account.CreditLimit,
		&account.DueDate, &account.CurrencyID, &account.IsActive, &account.CreatedAt, &account.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &account, nil
}

// FindByUserID finds accounts by user ID with pagination
func (r *accountRepository) FindByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entities.Account, error) {
	query := `
		SELECT account_id, user_id, account_name, account_type, bank_id, account_number,
		       current_balance, credit_limit, due_date, currency_id, is_active, created_at, updated_at
		FROM "vasst_expense".accounts 
		WHERE user_id = $1 AND is_active = true
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.DB.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []*entities.Account
	for rows.Next() {
		var account entities.Account
		err := rows.Scan(
			&account.AccountID, &account.UserID, &account.AccountName, &account.AccountType,
			&account.BankID, &account.AccountNumber, &account.CurrentBalance, &account.CreditLimit,
			&account.DueDate, &account.CurrencyID, &account.IsActive, &account.CreatedAt, &account.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, &account)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return accounts, nil
}

// FindActiveByUserID finds all active accounts by user ID (no pagination)
func (r *accountRepository) FindActiveByUserID(ctx context.Context, userID uuid.UUID) ([]*entities.Account, error) {
	query := `
		SELECT account_id, user_id, account_name, account_type, bank_id, account_number,
		       current_balance, credit_limit, due_date, currency_id, is_active, created_at, updated_at
		FROM "vasst_expense".accounts 
		WHERE user_id = $1 AND is_active = true
		ORDER BY account_name ASC
	`

	rows, err := r.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []*entities.Account
	for rows.Next() {
		var account entities.Account
		err := rows.Scan(
			&account.AccountID, &account.UserID, &account.AccountName, &account.AccountType,
			&account.BankID, &account.AccountNumber, &account.CurrentBalance, &account.CreditLimit,
			&account.DueDate, &account.CurrencyID, &account.IsActive, &account.CreatedAt, &account.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, &account)
	}

	return accounts, nil
}

// FindByNameAndUserID finds an account by name and user ID
func (r *accountRepository) FindByNameAndUserID(ctx context.Context, userID uuid.UUID, accountName string) (*entities.Account, error) {
	query := `
		SELECT account_id, user_id, account_name, account_type, bank_id, account_number,
		       current_balance, credit_limit, due_date, currency_id, is_active, created_at, updated_at
		FROM "vasst_expense".accounts 
		WHERE user_id = $1 AND account_name = $2 AND is_active = true
	`

	var account entities.Account
	err := r.DB.QueryRowContext(ctx, query, userID, accountName).Scan(
		&account.AccountID, &account.UserID, &account.AccountName, &account.AccountType,
		&account.BankID, &account.AccountNumber, &account.CurrentBalance, &account.CreditLimit,
		&account.DueDate, &account.CurrencyID, &account.IsActive, &account.CreatedAt, &account.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &account, nil
}
