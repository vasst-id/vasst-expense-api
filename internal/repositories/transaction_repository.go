package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/vasst-id/vasst-expense-api/internal/entities"
	"github.com/vasst-id/vasst-expense-api/internal/utils/postgres"
)

type (
	transactionRepository struct {
		*postgres.Postgres
	}

	TransactionRepository interface {
		Create(ctx context.Context, transaction *entities.Transaction) (entities.Transaction, error)
		Update(ctx context.Context, transaction *entities.Transaction) (entities.Transaction, error)
		Delete(ctx context.Context, transactionID uuid.UUID) error
		FindByID(ctx context.Context, transactionID uuid.UUID) (*entities.Transaction, error)
		FindByWorkspace(ctx context.Context, workspaceID uuid.UUID, params *entities.TransactionListParams, limit, offset int) ([]*entities.Transaction, error)
		FindByAccountID(ctx context.Context, accountID uuid.UUID, limit, offset int) ([]*entities.Transaction, error)
		FindByCategoryID(ctx context.Context, categoryID uuid.UUID, limit, offset int) ([]*entities.Transaction, error)
		CountByWorkspace(ctx context.Context, workspaceID uuid.UUID, params *entities.TransactionListParams) (int64, error)
	}
)

// NewTransactionRepository creates a new TransactionRepository
func NewTransactionRepository(pg *postgres.Postgres) TransactionRepository {
	return &transactionRepository{pg}
}

// Create creates a new transaction
func (r *transactionRepository) Create(ctx context.Context, transaction *entities.Transaction) (entities.Transaction, error) {
	query := `
		INSERT INTO "vasst_expense".transactions 
		(transaction_id, workspace_id, account_id, category_id, description, amount, 
		 transaction_type, transaction_date, merchant_name, location, 
		 notes, receipt_url, is_recurring, recurrence_interval, recurrence_end_date, 
		 parent_transaction_id, ai_confidence_score, ai_categorized, credit_status, 
		 created_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING transaction_id, workspace_id, account_id, category_id, description, amount,
		          transaction_type, transaction_date, merchant_name, location,
		          notes, receipt_url, is_recurring, recurrence_interval, recurrence_end_date,
		          parent_transaction_id, ai_confidence_score, ai_categorized, credit_status,
		          created_by, created_at, updated_at
	`

	var createdTransaction entities.Transaction
	err := r.DB.QueryRowContext(ctx, query,
		transaction.TransactionID, transaction.WorkspaceID, transaction.AccountID, transaction.CategoryID,
		transaction.Description, transaction.Amount, transaction.TransactionType,
		transaction.TransactionDate, transaction.MerchantName, transaction.Location, transaction.Notes,
		transaction.ReceiptURL, transaction.IsRecurring, transaction.RecurrenceInterval, transaction.RecurrenceEndDate,
		transaction.ParentTransactionID, transaction.AIConfidenceScore, transaction.AICategorized, transaction.CreditStatus,
		transaction.CreatedBy,
	).Scan(
		&createdTransaction.TransactionID, &createdTransaction.WorkspaceID, &createdTransaction.AccountID, &createdTransaction.CategoryID,
		&createdTransaction.Description, &createdTransaction.Amount, &createdTransaction.TransactionType,
		&createdTransaction.TransactionDate, &createdTransaction.MerchantName, &createdTransaction.Location, &createdTransaction.Notes,
		&createdTransaction.ReceiptURL, &createdTransaction.IsRecurring, &createdTransaction.RecurrenceInterval, &createdTransaction.RecurrenceEndDate,
		&createdTransaction.ParentTransactionID, &createdTransaction.AIConfidenceScore, &createdTransaction.AICategorized, &createdTransaction.CreditStatus,
		&createdTransaction.CreatedBy, &createdTransaction.CreatedAt, &createdTransaction.UpdatedAt,
	)

	return createdTransaction, err
}

// Update updates a transaction
func (r *transactionRepository) Update(ctx context.Context, transaction *entities.Transaction) (entities.Transaction, error) {
	query := `
		UPDATE "vasst_expense".transactions 
		SET account_id = $2, category_id = $3, description = $4, amount = $5,
		    transaction_type = $6, transaction_date = $8,
		    merchant_name = $9, location = $10, notes = $11, receipt_url = $12,
		    is_recurring = $13, recurrence_interval = $14, recurrence_end_date = $15,
		    parent_transaction_id = $16, ai_confidence_score = $17, ai_categorized = $18,
		    credit_status = $19, updated_at = CURRENT_TIMESTAMP
		WHERE transaction_id = $1
		RETURNING transaction_id, workspace_id, account_id, category_id, description, amount,
		          transaction_type, transaction_date, merchant_name, location,
		          notes, receipt_url, is_recurring, recurrence_interval, recurrence_end_date,
		          parent_transaction_id, ai_confidence_score, ai_categorized, credit_status,
		          created_by, created_at, updated_at
	`

	var updatedTransaction entities.Transaction
	err := r.DB.QueryRowContext(ctx, query,
		transaction.TransactionID, transaction.AccountID, transaction.CategoryID, transaction.Description,
		transaction.Amount, transaction.TransactionType, transaction.TransactionDate,
		transaction.MerchantName, transaction.Location, transaction.Notes, transaction.ReceiptURL,
		transaction.IsRecurring, transaction.RecurrenceInterval, transaction.RecurrenceEndDate,
		transaction.ParentTransactionID, transaction.AIConfidenceScore, transaction.AICategorized, transaction.CreditStatus,
	).Scan(
		&updatedTransaction.TransactionID, &updatedTransaction.WorkspaceID, &updatedTransaction.AccountID, &updatedTransaction.CategoryID,
		&updatedTransaction.Description, &updatedTransaction.Amount, &updatedTransaction.TransactionType,
		&updatedTransaction.TransactionDate, &updatedTransaction.MerchantName, &updatedTransaction.Location, &updatedTransaction.Notes,
		&updatedTransaction.ReceiptURL, &updatedTransaction.IsRecurring, &updatedTransaction.RecurrenceInterval, &updatedTransaction.RecurrenceEndDate,
		&updatedTransaction.ParentTransactionID, &updatedTransaction.AIConfidenceScore, &updatedTransaction.AICategorized, &updatedTransaction.CreditStatus,
		&updatedTransaction.CreatedBy, &updatedTransaction.CreatedAt, &updatedTransaction.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return entities.Transaction{}, sql.ErrNoRows // Transaction not found
		}
		return entities.Transaction{}, err
	}

	return updatedTransaction, nil
}

// Delete soft deletes a transaction (marks as inactive)
func (r *transactionRepository) Delete(ctx context.Context, transactionID uuid.UUID) error {
	query := `
		DELETE FROM "vasst_expense".transactions 
		WHERE transaction_id = $1
	`

	result, err := r.DB.ExecContext(ctx, query, transactionID)
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

// FindByID finds a transaction by ID
func (r *transactionRepository) FindByID(ctx context.Context, transactionID uuid.UUID) (*entities.Transaction, error) {
	query := `
		SELECT transaction_id, workspace_id, account_id, category_id, description, amount,
		       transaction_type, transaction_date, merchant_name, location,
		       notes, receipt_url, is_recurring, recurrence_interval, recurrence_end_date,
		       parent_transaction_id, ai_confidence_score, ai_categorized, credit_status,
		       created_by, created_at, updated_at
		FROM "vasst_expense".transactions 
		WHERE transaction_id = $1
	`

	var transaction entities.Transaction
	err := r.DB.QueryRowContext(ctx, query, transactionID).Scan(
		&transaction.TransactionID, &transaction.WorkspaceID, &transaction.AccountID, &transaction.CategoryID,
		&transaction.Description, &transaction.Amount, &transaction.TransactionType,
		&transaction.TransactionDate, &transaction.MerchantName, &transaction.Location, &transaction.Notes,
		&transaction.ReceiptURL, &transaction.IsRecurring, &transaction.RecurrenceInterval, &transaction.RecurrenceEndDate,
		&transaction.ParentTransactionID, &transaction.AIConfidenceScore, &transaction.AICategorized, &transaction.CreditStatus,
		&transaction.CreatedBy, &transaction.CreatedAt, &transaction.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &transaction, nil
}

// FindByWorkspace finds transactions by workspace with filtering and pagination
func (r *transactionRepository) FindByWorkspace(ctx context.Context, workspaceID uuid.UUID, params *entities.TransactionListParams, limit, offset int) ([]*entities.Transaction, error) {
	query := `
		SELECT transaction_id, workspace_id, account_id, category_id, description, amount,
		       transaction_type, transaction_date, merchant_name, location,
		       notes, receipt_url, is_recurring, recurrence_interval, recurrence_end_date,
		       parent_transaction_id, ai_confidence_score, ai_categorized, credit_status,
		       created_by, created_at, updated_at
		FROM "vasst_expense".transactions 
		WHERE workspace_id = $1
	`

	args := []interface{}{workspaceID}
	argIndex := 2

	// Add filters from params
	if params != nil {
		if params.AccountID != nil {
			query += fmt.Sprintf(" AND account_id = $%d", argIndex)
			args = append(args, *params.AccountID)
			argIndex++
		}
		if params.CategoryID != nil {
			query += fmt.Sprintf(" AND category_id = $%d", argIndex)
			args = append(args, *params.CategoryID)
			argIndex++
		}
		if params.StartDate != nil {
			query += fmt.Sprintf(" AND transaction_date >= $%d", argIndex)
			args = append(args, *params.StartDate)
			argIndex++
		}
		if params.EndDate != nil {
			query += fmt.Sprintf(" AND transaction_date <= $%d", argIndex)
			args = append(args, *params.EndDate)
			argIndex++
		}
		if params.Description != nil {
			query += fmt.Sprintf(" AND description ILIKE $%d", argIndex)
			args = append(args, "%"+*params.Description+"%")
			argIndex++
		}
		if params.MerchantName != nil {
			query += fmt.Sprintf(" AND merchant_name ILIKE $%d", argIndex)
			args = append(args, "%"+*params.MerchantName+"%")
			argIndex++
		}
		if params.Amount != nil {
			query += fmt.Sprintf(" AND amount = $%d", argIndex)
			args = append(args, *params.Amount)
			argIndex++
		}
		if params.IsRecurring != nil {
			query += fmt.Sprintf(" AND is_recurring = $%d", argIndex)
			args = append(args, *params.IsRecurring)
			argIndex++
		}
		if params.CreditStatus != nil {
			query += fmt.Sprintf(" AND credit_status = $%d", argIndex)
			args = append(args, *params.CreditStatus)
			argIndex++
		}
	}

	query += " ORDER BY transaction_date DESC, created_at DESC"
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, limit, offset)

	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*entities.Transaction
	for rows.Next() {
		var transaction entities.Transaction
		err := rows.Scan(
			&transaction.TransactionID, &transaction.WorkspaceID, &transaction.AccountID, &transaction.CategoryID,
			&transaction.Description, &transaction.Amount, &transaction.TransactionType,
			&transaction.TransactionDate, &transaction.MerchantName, &transaction.Location, &transaction.Notes,
			&transaction.ReceiptURL, &transaction.IsRecurring, &transaction.RecurrenceInterval, &transaction.RecurrenceEndDate,
			&transaction.ParentTransactionID, &transaction.AIConfidenceScore, &transaction.AICategorized, &transaction.CreditStatus,
			&transaction.CreatedBy, &transaction.CreatedAt, &transaction.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, &transaction)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}

// FindByAccountID finds transactions by account ID with pagination
func (r *transactionRepository) FindByAccountID(ctx context.Context, accountID uuid.UUID, limit, offset int) ([]*entities.Transaction, error) {
	query := `
		SELECT transaction_id, workspace_id, account_id, category_id, description, amount,
		       transaction_type, transaction_date, merchant_name, location,
		       notes, receipt_url, is_recurring, recurrence_interval, recurrence_end_date,
		       parent_transaction_id, ai_confidence_score, ai_categorized, credit_status,
		       created_by, created_at, updated_at
		FROM "vasst_expense".transactions 
		WHERE account_id = $1
		ORDER BY transaction_date DESC, created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.DB.QueryContext(ctx, query, accountID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*entities.Transaction
	for rows.Next() {
		var transaction entities.Transaction
		err := rows.Scan(
			&transaction.TransactionID, &transaction.WorkspaceID, &transaction.AccountID, &transaction.CategoryID,
			&transaction.Description, &transaction.Amount, &transaction.TransactionType,
			&transaction.TransactionDate, &transaction.MerchantName, &transaction.Location, &transaction.Notes,
			&transaction.ReceiptURL, &transaction.IsRecurring, &transaction.RecurrenceInterval, &transaction.RecurrenceEndDate,
			&transaction.ParentTransactionID, &transaction.AIConfidenceScore, &transaction.AICategorized, &transaction.CreditStatus,
			&transaction.CreatedBy, &transaction.CreatedAt, &transaction.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, &transaction)
	}

	return transactions, nil
}

// FindByCategoryID finds transactions by category ID with pagination
func (r *transactionRepository) FindByCategoryID(ctx context.Context, categoryID uuid.UUID, limit, offset int) ([]*entities.Transaction, error) {
	query := `
		SELECT transaction_id, workspace_id, account_id, category_id, description, amount,
		       transaction_type, transaction_date, merchant_name, location,
		       notes, receipt_url, is_recurring, recurrence_interval, recurrence_end_date,
		       parent_transaction_id, ai_confidence_score, ai_categorized, credit_status,
		       created_by, created_at, updated_at
		FROM "vasst_expense".transactions 
		WHERE category_id = $1
		ORDER BY transaction_date DESC, created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.DB.QueryContext(ctx, query, categoryID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*entities.Transaction
	for rows.Next() {
		var transaction entities.Transaction
		err := rows.Scan(
			&transaction.TransactionID, &transaction.WorkspaceID, &transaction.AccountID, &transaction.CategoryID,
			&transaction.Description, &transaction.Amount, &transaction.TransactionType,
			&transaction.TransactionDate, &transaction.MerchantName, &transaction.Location, &transaction.Notes,
			&transaction.ReceiptURL, &transaction.IsRecurring, &transaction.RecurrenceInterval, &transaction.RecurrenceEndDate,
			&transaction.ParentTransactionID, &transaction.AIConfidenceScore, &transaction.AICategorized, &transaction.CreditStatus,
			&transaction.CreatedBy, &transaction.CreatedAt, &transaction.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, &transaction)
	}

	return transactions, nil
}

// CountByWorkspace counts transactions by workspace with filtering
func (r *transactionRepository) CountByWorkspace(ctx context.Context, workspaceID uuid.UUID, params *entities.TransactionListParams) (int64, error) {
	query := `
		SELECT COUNT(*)
		FROM "vasst_expense".transactions 
		WHERE workspace_id = $1
	`

	args := []interface{}{workspaceID}
	argIndex := 2

	// Add same filters as FindByWorkspace
	if params != nil {
		if params.AccountID != nil {
			query += fmt.Sprintf(" AND account_id = $%d", argIndex)
			args = append(args, *params.AccountID)
			argIndex++
		}
		if params.CategoryID != nil {
			query += fmt.Sprintf(" AND category_id = $%d", argIndex)
			args = append(args, *params.CategoryID)
			argIndex++
		}
		if params.StartDate != nil {
			query += fmt.Sprintf(" AND transaction_date >= $%d", argIndex)
			args = append(args, *params.StartDate)
			argIndex++
		}
		if params.EndDate != nil {
			query += fmt.Sprintf(" AND transaction_date <= $%d", argIndex)
			args = append(args, *params.EndDate)
			argIndex++
		}
		if params.Description != nil {
			query += fmt.Sprintf(" AND description ILIKE $%d", argIndex)
			args = append(args, "%"+*params.Description+"%")
			argIndex++
		}
		if params.MerchantName != nil {
			query += fmt.Sprintf(" AND merchant_name ILIKE $%d", argIndex)
			args = append(args, "%"+*params.MerchantName+"%")
			argIndex++
		}
		if params.Amount != nil {
			query += fmt.Sprintf(" AND amount = $%d", argIndex)
			args = append(args, *params.Amount)
			argIndex++
		}
		if params.IsRecurring != nil {
			query += fmt.Sprintf(" AND is_recurring = $%d", argIndex)
			args = append(args, *params.IsRecurring)
			argIndex++
		}
		if params.CreditStatus != nil {
			query += fmt.Sprintf(" AND credit_status = $%d", argIndex)
			args = append(args, *params.CreditStatus)
			argIndex++
		}
	}

	var count int64
	err := r.DB.QueryRowContext(ctx, query, args...).Scan(&count)
	return count, err
}
