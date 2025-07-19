package repositories

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/vasst-id/vasst-expense-api/internal/entities"
	"github.com/vasst-id/vasst-expense-api/internal/utils/postgres"
)

type (
	userTagsRepository struct {
		*postgres.Postgres
	}

	UserTagsRepository interface {
		Create(ctx context.Context, userTag *entities.UserTag) (entities.UserTag, error)
		Update(ctx context.Context, userTag *entities.UserTag) (entities.UserTag, error)
		Delete(ctx context.Context, userTagID uuid.UUID) error
		FindByID(ctx context.Context, userTagID uuid.UUID) (*entities.UserTag, error)
		FindByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entities.UserTag, error)
		FindActiveByUserID(ctx context.Context, userID uuid.UUID) ([]*entities.UserTag, error)
		FindByNameAndUserID(ctx context.Context, userID uuid.UUID, name string) (*entities.UserTag, error)
		GetUserTagsWithUsage(ctx context.Context, userID uuid.UUID) ([]*entities.UserTagWithUsage, error)
	}
)

// NewUserTagsRepository creates a new UserTagsRepository
func NewUserTagsRepository(pg *postgres.Postgres) UserTagsRepository {
	return &userTagsRepository{pg}
}

// Create creates a new user tag
func (r *userTagsRepository) Create(ctx context.Context, userTag *entities.UserTag) (entities.UserTag, error) {
	query := `
		INSERT INTO "vasst_expense".user_tags 
		(user_tag_id, user_id, name, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING user_tag_id, user_id, name, is_active, created_at, updated_at
	`

	var createdUserTag entities.UserTag
	err := r.DB.QueryRowContext(ctx, query,
		userTag.UserTagID, userTag.UserID, userTag.Name, userTag.IsActive,
	).Scan(
		&createdUserTag.UserTagID, &createdUserTag.UserID, &createdUserTag.Name,
		&createdUserTag.IsActive, &createdUserTag.CreatedAt, &createdUserTag.UpdatedAt,
	)

	return createdUserTag, err
}

// Update updates a user tag
func (r *userTagsRepository) Update(ctx context.Context, userTag *entities.UserTag) (entities.UserTag, error) {
	query := `
		UPDATE "vasst_expense".user_tags 
		SET name = $2, is_active = $3, updated_at = CURRENT_TIMESTAMP
		WHERE user_tag_id = $1
		RETURNING user_tag_id, user_id, name, is_active, created_at, updated_at
	`

	var updatedUserTag entities.UserTag
	err := r.DB.QueryRowContext(ctx, query,
		userTag.UserTagID, userTag.Name, userTag.IsActive,
	).Scan(
		&updatedUserTag.UserTagID, &updatedUserTag.UserID, &updatedUserTag.Name,
		&updatedUserTag.IsActive, &updatedUserTag.CreatedAt, &updatedUserTag.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return entities.UserTag{}, sql.ErrNoRows // User tag not found
		}
		return entities.UserTag{}, err
	}

	return updatedUserTag, nil
}

// Delete soft deletes a user tag (sets is_active to false)
func (r *userTagsRepository) Delete(ctx context.Context, userTagID uuid.UUID) error {
	query := `
		UPDATE "vasst_expense".user_tags 
		SET is_active = false, updated_at = CURRENT_TIMESTAMP
		WHERE user_tag_id = $1
	`

	result, err := r.DB.ExecContext(ctx, query, userTagID)
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

// FindByID finds a user tag by ID
func (r *userTagsRepository) FindByID(ctx context.Context, userTagID uuid.UUID) (*entities.UserTag, error) {
	query := `
		SELECT user_tag_id, user_id, name, is_active, created_at, updated_at
		FROM "vasst_expense".user_tags 
		WHERE user_tag_id = $1 AND is_active = true
	`

	var userTag entities.UserTag
	err := r.DB.QueryRowContext(ctx, query, userTagID).Scan(
		&userTag.UserTagID, &userTag.UserID, &userTag.Name,
		&userTag.IsActive, &userTag.CreatedAt, &userTag.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &userTag, nil
}

// FindByUserID finds user tags by user ID with pagination
func (r *userTagsRepository) FindByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entities.UserTag, error) {
	query := `
		SELECT user_tag_id, user_id, name, is_active, created_at, updated_at
		FROM "vasst_expense".user_tags 
		WHERE user_id = $1 AND is_active = true
		ORDER BY name ASC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.DB.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userTags []*entities.UserTag
	for rows.Next() {
		var userTag entities.UserTag
		err := rows.Scan(
			&userTag.UserTagID, &userTag.UserID, &userTag.Name,
			&userTag.IsActive, &userTag.CreatedAt, &userTag.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		userTags = append(userTags, &userTag)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return userTags, nil
}

// FindActiveByUserID finds all active user tags by user ID (no pagination)
func (r *userTagsRepository) FindActiveByUserID(ctx context.Context, userID uuid.UUID) ([]*entities.UserTag, error) {
	query := `
		SELECT user_tag_id, user_id, name, is_active, created_at, updated_at
		FROM "vasst_expense".user_tags 
		WHERE user_id = $1 AND is_active = true
		ORDER BY name ASC
	`

	rows, err := r.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userTags []*entities.UserTag
	for rows.Next() {
		var userTag entities.UserTag
		err := rows.Scan(
			&userTag.UserTagID, &userTag.UserID, &userTag.Name,
			&userTag.IsActive, &userTag.CreatedAt, &userTag.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		userTags = append(userTags, &userTag)
	}

	return userTags, nil
}

// FindByNameAndUserID finds a user tag by name and user ID
func (r *userTagsRepository) FindByNameAndUserID(ctx context.Context, userID uuid.UUID, name string) (*entities.UserTag, error) {
	query := `
		SELECT user_tag_id, user_id, name, is_active, created_at, updated_at
		FROM "vasst_expense".user_tags 
		WHERE user_id = $1 AND name = $2 AND is_active = true
	`

	var userTag entities.UserTag
	err := r.DB.QueryRowContext(ctx, query, userID, name).Scan(
		&userTag.UserTagID, &userTag.UserID, &userTag.Name,
		&userTag.IsActive, &userTag.CreatedAt, &userTag.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &userTag, nil
}

// GetUserTagsWithUsage gets user tags with their usage statistics
func (r *userTagsRepository) GetUserTagsWithUsage(ctx context.Context, userID uuid.UUID) ([]*entities.UserTagWithUsage, error) {
	query := `
		SELECT 
			ut.user_tag_id,
			ut.name,
			ut.is_active,
			ut.created_at,
			ut.updated_at,
			COUNT(tt.transaction_tag_id) as usage_count,
			COALESCE(MAX(tt.applied_at), ut.created_at) as last_used_at
		FROM "vasst_expense".user_tags ut
		LEFT JOIN "vasst_expense".transaction_tags tt ON ut.user_tag_id = tt.user_tag_id
		WHERE ut.user_id = $1 AND ut.is_active = true
		GROUP BY ut.user_tag_id, ut.name, ut.is_active, ut.created_at, ut.updated_at
		ORDER BY usage_count DESC, ut.name ASC
	`

	rows, err := r.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userTagsWithUsage []*entities.UserTagWithUsage
	for rows.Next() {
		var userTagWithUsage entities.UserTagWithUsage
		err := rows.Scan(
			&userTagWithUsage.UserTagID, &userTagWithUsage.Name,
			&userTagWithUsage.IsActive, &userTagWithUsage.CreatedAt, &userTagWithUsage.UpdatedAt,
			&userTagWithUsage.UsageCount, &userTagWithUsage.LastUsedAt,
		)
		if err != nil {
			return nil, err
		}
		userTagsWithUsage = append(userTagsWithUsage, &userTagWithUsage)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return userTagsWithUsage, nil
}
