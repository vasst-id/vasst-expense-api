package repositories

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/vasst-id/vasst-expense-api/internal/entities"
	"github.com/vasst-id/vasst-expense-api/internal/utils/postgres"
)

type (
	transactionTagsRepository struct {
		*postgres.Postgres
	}

	TransactionTagsRepository interface {
		Create(ctx context.Context, transactionTag *entities.TransactionTag) (entities.TransactionTag, error)
		CreateMultiple(ctx context.Context, transactionTags []*entities.TransactionTag) ([]entities.TransactionTag, error)
		Delete(ctx context.Context, transactionTagID uuid.UUID) error
		DeleteByTransactionID(ctx context.Context, transactionID uuid.UUID) error
		DeleteByUserTagID(ctx context.Context, userTagID uuid.UUID) error
		FindByID(ctx context.Context, transactionTagID uuid.UUID) (*entities.TransactionTag, error)
		FindByTransactionID(ctx context.Context, transactionID uuid.UUID) ([]*entities.TransactionTag, error)
		FindByUserTagID(ctx context.Context, userTagID uuid.UUID, limit, offset int) ([]*entities.TransactionTag, error)
		FindTransactionsByUserTagID(ctx context.Context, userTagID uuid.UUID, limit, offset int) ([]*entities.TransactionWithTags, error)
		GetTaggedTransactionsSummary(ctx context.Context, userID uuid.UUID) ([]*entities.TaggedTransactionSummary, error)
	}
)

// NewTransactionTagsRepository creates a new TransactionTagsRepository
func NewTransactionTagsRepository(pg *postgres.Postgres) TransactionTagsRepository {
	return &transactionTagsRepository{pg}
}

// Create creates a new transaction tag
func (r *transactionTagsRepository) Create(ctx context.Context, transactionTag *entities.TransactionTag) (entities.TransactionTag, error) {
	query := `
		INSERT INTO "vasst_expense".transaction_tags 
		(transaction_tag_id, transaction_id, user_tag_id, applied_by, applied_at)
		VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP)
		RETURNING transaction_tag_id, transaction_id, user_tag_id, applied_by, applied_at
	`

	var createdTransactionTag entities.TransactionTag
	err := r.DB.QueryRowContext(ctx, query,
		transactionTag.TransactionTagID, transactionTag.TransactionID,
		transactionTag.UserTagID, transactionTag.AppliedBy,
	).Scan(
		&createdTransactionTag.TransactionTagID, &createdTransactionTag.TransactionID,
		&createdTransactionTag.UserTagID, &createdTransactionTag.AppliedBy,
		&createdTransactionTag.AppliedAt,
	)

	return createdTransactionTag, err
}

// CreateMultiple creates multiple transaction tags in a single transaction
func (r *transactionTagsRepository) CreateMultiple(ctx context.Context, transactionTags []*entities.TransactionTag) ([]entities.TransactionTag, error) {
	if len(transactionTags) == 0 {
		return []entities.TransactionTag{}, nil
	}

	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO "vasst_expense".transaction_tags 
		(transaction_tag_id, transaction_id, user_tag_id, applied_by, applied_at)
		VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP)
		RETURNING transaction_tag_id, transaction_id, user_tag_id, applied_by, applied_at
	`

	var createdTransactionTags []entities.TransactionTag
	for _, transactionTag := range transactionTags {
		var createdTransactionTag entities.TransactionTag
		err := tx.QueryRowContext(ctx, query,
			transactionTag.TransactionTagID, transactionTag.TransactionID,
			transactionTag.UserTagID, transactionTag.AppliedBy,
		).Scan(
			&createdTransactionTag.TransactionTagID, &createdTransactionTag.TransactionID,
			&createdTransactionTag.UserTagID, &createdTransactionTag.AppliedBy,
			&createdTransactionTag.AppliedAt,
		)
		if err != nil {
			return nil, err
		}
		createdTransactionTags = append(createdTransactionTags, createdTransactionTag)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return createdTransactionTags, nil
}

// Delete deletes a transaction tag
func (r *transactionTagsRepository) Delete(ctx context.Context, transactionTagID uuid.UUID) error {
	query := `DELETE FROM "vasst_expense".transaction_tags WHERE transaction_tag_id = $1`

	result, err := r.DB.ExecContext(ctx, query, transactionTagID)
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

// DeleteByTransactionID deletes all transaction tags for a specific transaction
func (r *transactionTagsRepository) DeleteByTransactionID(ctx context.Context, transactionID uuid.UUID) error {
	query := `DELETE FROM "vasst_expense".transaction_tags WHERE transaction_id = $1`

	_, err := r.DB.ExecContext(ctx, query, transactionID)
	return err
}

// DeleteByUserTagID deletes all transaction tags for a specific user tag
func (r *transactionTagsRepository) DeleteByUserTagID(ctx context.Context, userTagID uuid.UUID) error {
	query := `DELETE FROM "vasst_expense".transaction_tags WHERE user_tag_id = $1`

	_, err := r.DB.ExecContext(ctx, query, userTagID)
	return err
}

// FindByID finds a transaction tag by ID
func (r *transactionTagsRepository) FindByID(ctx context.Context, transactionTagID uuid.UUID) (*entities.TransactionTag, error) {
	query := `
		SELECT transaction_tag_id, transaction_id, user_tag_id, applied_by, applied_at
		FROM "vasst_expense".transaction_tags 
		WHERE transaction_tag_id = $1
	`

	var transactionTag entities.TransactionTag
	err := r.DB.QueryRowContext(ctx, query, transactionTagID).Scan(
		&transactionTag.TransactionTagID, &transactionTag.TransactionID,
		&transactionTag.UserTagID, &transactionTag.AppliedBy,
		&transactionTag.AppliedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &transactionTag, nil
}

// FindByTransactionID finds all transaction tags for a specific transaction
func (r *transactionTagsRepository) FindByTransactionID(ctx context.Context, transactionID uuid.UUID) ([]*entities.TransactionTag, error) {
	query := `
		SELECT transaction_tag_id, transaction_id, user_tag_id, applied_by, applied_at
		FROM "vasst_expense".transaction_tags 
		WHERE transaction_id = $1
		ORDER BY applied_at DESC
	`

	rows, err := r.DB.QueryContext(ctx, query, transactionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactionTags []*entities.TransactionTag
	for rows.Next() {
		var transactionTag entities.TransactionTag
		err := rows.Scan(
			&transactionTag.TransactionTagID, &transactionTag.TransactionID,
			&transactionTag.UserTagID, &transactionTag.AppliedBy,
			&transactionTag.AppliedAt,
		)
		if err != nil {
			return nil, err
		}
		transactionTags = append(transactionTags, &transactionTag)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return transactionTags, nil
}

// FindByUserTagID finds all transaction tags for a specific user tag with pagination
func (r *transactionTagsRepository) FindByUserTagID(ctx context.Context, userTagID uuid.UUID, limit, offset int) ([]*entities.TransactionTag, error) {
	query := `
		SELECT transaction_tag_id, transaction_id, user_tag_id, applied_by, applied_at
		FROM "vasst_expense".transaction_tags 
		WHERE user_tag_id = $1
		ORDER BY applied_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.DB.QueryContext(ctx, query, userTagID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactionTags []*entities.TransactionTag
	for rows.Next() {
		var transactionTag entities.TransactionTag
		err := rows.Scan(
			&transactionTag.TransactionTagID, &transactionTag.TransactionID,
			&transactionTag.UserTagID, &transactionTag.AppliedBy,
			&transactionTag.AppliedAt,
		)
		if err != nil {
			return nil, err
		}
		transactionTags = append(transactionTags, &transactionTag)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return transactionTags, nil
}

// FindTransactionsByUserTagID finds transactions with their tags for a specific user tag
func (r *transactionTagsRepository) FindTransactionsByUserTagID(ctx context.Context, userTagID uuid.UUID, limit, offset int) ([]*entities.TransactionWithTags, error) {
	query := `
		SELECT DISTINCT t.transaction_id
		FROM "vasst_expense".transactions t
		INNER JOIN "vasst_expense".transaction_tags tt ON t.transaction_id = tt.transaction_id
		WHERE tt.user_tag_id = $1
		ORDER BY t.transaction_date DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.DB.QueryContext(ctx, query, userTagID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactionsWithTags []*entities.TransactionWithTags
	for rows.Next() {
		var transactionID uuid.UUID
		err := rows.Scan(&transactionID)
		if err != nil {
			return nil, err
		}

		// Get tags for this transaction
		tags, err := r.FindByTransactionID(ctx, transactionID)
		if err != nil {
			return nil, err
		}

		// Convert to simple format
		var tagSimples []entities.TransactionTagSimple
		for _, tag := range tags {
			tagSimples = append(tagSimples, entities.TransactionTagSimple{
				TransactionTagID: tag.TransactionTagID,
				TransactionID:    tag.TransactionID,
				UserTagID:        tag.UserTagID,
				AppliedBy:        tag.AppliedBy,
				AppliedAt:        tag.AppliedAt,
			})
		}

		transactionWithTags := &entities.TransactionWithTags{
			TransactionID: transactionID,
			Tags:          tagSimples,
		}
		transactionsWithTags = append(transactionsWithTags, transactionWithTags)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return transactionsWithTags, nil
}

// GetTaggedTransactionsSummary gets summary of tagged transactions by user
func (r *transactionTagsRepository) GetTaggedTransactionsSummary(ctx context.Context, userID uuid.UUID) ([]*entities.TaggedTransactionSummary, error) {
	query := `
		SELECT 
			ut.user_tag_id,
			ut.name as tag_name,
			COUNT(DISTINCT tt.transaction_id) as transaction_count,
			COALESCE(SUM(t.amount), 0) as total_amount,
			COALESCE(MAX(tt.applied_at), ut.created_at) as last_used_at
		FROM "vasst_expense".user_tags ut
		LEFT JOIN "vasst_expense".transaction_tags tt ON ut.user_tag_id = tt.user_tag_id
		LEFT JOIN "vasst_expense".transactions t ON tt.transaction_id = t.transaction_id
		WHERE ut.user_id = $1 AND ut.is_active = true
		GROUP BY ut.user_tag_id, ut.name, ut.created_at
		HAVING COUNT(DISTINCT tt.transaction_id) > 0
		ORDER BY transaction_count DESC, total_amount DESC
	`

	rows, err := r.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var summaries []*entities.TaggedTransactionSummary
	for rows.Next() {
		var summary entities.TaggedTransactionSummary
		err := rows.Scan(
			&summary.UserTagID, &summary.TagName, &summary.TransactionCount,
			&summary.TotalAmount, &summary.LastUsedAt,
		)
		if err != nil {
			return nil, err
		}
		summaries = append(summaries, &summary)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return summaries, nil
}
