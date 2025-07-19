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
	userRepository struct {
		*postgres.Postgres
	}

	UserRepository interface {
		Create(ctx context.Context, user *entities.User) (entities.User, error)
		Update(ctx context.Context, user *entities.User) (entities.User, error)
		Delete(ctx context.Context, userID uuid.UUID) error
		ListAll(ctx context.Context, limit, offset int) ([]*entities.User, error)
		FindByID(ctx context.Context, userID uuid.UUID) (*entities.User, error)
		FindByEmail(ctx context.Context, email string) (*entities.User, error)
		FindByPhoneNumber(ctx context.Context, phoneNumber string) (*entities.User, error)
	}
)

// NewUserRepository creates a new UserRepository
func NewUserRepository(pg *postgres.Postgres) UserRepository {
	return &userRepository{pg}
}

// Create creates a new user
func (r *userRepository) Create(ctx context.Context, user *entities.User) (entities.User, error) {
	query := `
		INSERT INTO "vasst_expense".users (
			user_id, email, phone_number, password_hash, first_name, last_name, 
			timezone, currency_id, subscription_plan_id, email_verified_at, 
			phone_verified_at, status, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING user_id, email, phone_number, password_hash, first_name, last_name, 
		          timezone, currency_id, subscription_plan_id, email_verified_at, 
		          phone_verified_at, status, created_at, updated_at
	`

	// Handle nullable timestamp fields
	var emailVerifiedAt, phoneVerifiedAt sql.NullTime
	if user.EmailVerifiedAt != nil {
		emailVerifiedAt.Time = *user.EmailVerifiedAt
		emailVerifiedAt.Valid = true
	}
	if user.PhoneVerifiedAt != nil {
		phoneVerifiedAt.Time = *user.PhoneVerifiedAt
		phoneVerifiedAt.Valid = true
	}

	// Handle email field - convert empty string to NULL
	var email sql.NullString
	if user.Email != "" {
		email.String = user.Email
		email.Valid = true
	}

	var createdUser entities.User
	var returnedEmailVerifiedAt, returnedPhoneVerifiedAt sql.NullTime
	var returnedEmail sql.NullString

	err := r.DB.QueryRowContext(ctx, query,
		user.UserID,
		email,
		user.PhoneNumber,
		user.PasswordHash,
		user.FirstName,
		user.LastName,
		user.Timezone,
		user.CurrencyID,
		user.SubscriptionPlanID,
		emailVerifiedAt,
		phoneVerifiedAt,
		user.Status,
	).Scan(
		&createdUser.UserID,
		&returnedEmail,
		&createdUser.PhoneNumber,
		&createdUser.PasswordHash,
		&createdUser.FirstName,
		&createdUser.LastName,
		&createdUser.Timezone,
		&createdUser.CurrencyID,
		&createdUser.SubscriptionPlanID,
		&returnedEmailVerifiedAt,
		&returnedPhoneVerifiedAt,
		&createdUser.Status,
		&createdUser.CreatedAt,
		&createdUser.UpdatedAt,
	)

	if err != nil {
		return entities.User{}, err
	}

	// Handle nullable fields
	if returnedEmail.Valid {
		createdUser.Email = returnedEmail.String
	}
	if returnedEmailVerifiedAt.Valid {
		createdUser.EmailVerifiedAt = &returnedEmailVerifiedAt.Time
	}
	if returnedPhoneVerifiedAt.Valid {
		createdUser.PhoneVerifiedAt = &returnedPhoneVerifiedAt.Time
	}

	// Create a workspace for the user
	workspace := &entities.Workspace{
		Name:          "My Personal Workspace",
		Description:   "This is my personal workspace for managing my expenses",
		WorkspaceType: entities.WorkspaceTypePersonal,
		CurrencyID:    createdUser.CurrencyID,
		Timezone:      createdUser.Timezone,
		CreatedBy:     createdUser.UserID,
		IsActive:      true,
	}

	queryWorkspace := `
		INSERT INTO "vasst_expense".workspaces (
			name, description, workspace_type, currency_id, timezone, is_active, created_by
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	fmt.Println("workspace", workspace)

	_, err = r.DB.ExecContext(ctx, queryWorkspace,
		workspace.Name,
		workspace.Description,
		workspace.WorkspaceType,
		workspace.CurrencyID,
		workspace.Timezone,
		workspace.IsActive,
		workspace.CreatedBy,
	)
	if err != nil {
		return entities.User{}, err
	}

	// // Create user categories
	// // Get all system categories
	// systemCategoriesQuery := `
	// 	SELECT category_id, name, description, icon, parent_category_id
	// 	FROM "vasst_expense".categories
	// 	WHERE is_system_category = true
	// `
	// rows, err := r.DB.QueryContext(ctx, systemCategoriesQuery)
	// if err != nil {
	// 	return entities.User{}, err
	// }
	// defer rows.Close()

	// type sysCat struct {
	// 	CategoryID       string
	// 	Name             string
	// 	Description      sql.NullString
	// 	Icon             sql.NullString
	// 	ParentCategoryID sql.NullString
	// }
	// var systemCategories []sysCat

	// for rows.Next() {
	// 	var cat sysCat
	// 	if err := rows.Scan(&cat.CategoryID, &cat.Name, &cat.Description, &cat.Icon, &cat.ParentCategoryID); err != nil {
	// 		return entities.User{}, err
	// 	}
	// 	systemCategories = append(systemCategories, cat)
	// }
	// if err := rows.Err(); err != nil {
	// 	return entities.User{}, err
	// }

	// if len(systemCategories) > 0 {
	// 	// Prepare batch insert for user_categories
	// 	valueStrings := make([]string, 0, len(systemCategories))
	// 	valueArgs := make([]interface{}, 0, len(systemCategories)*8)
	// 	for i, cat := range systemCategories {
	// 		userCategoryID := uuid.New()
	// 		valueStrings = append(valueStrings, fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d)",
	// 			i*8+1, i*8+2, i*8+3, i*8+4, i*8+5, i*8+6, i*8+7, i*8+8))
	// 		valueArgs = append(valueArgs,
	// 			userCategoryID,     // user_category_id
	// 			createdUser.UserID, // user_id
	// 			cat.CategoryID,     // category_id
	// 			cat.Name,           // name
	// 			cat.Description,    // description
	// 			cat.Icon,           // icon
	// 			false,              // is_custom
	// 			true,               // is_active
	// 		)
	// 	}
	// 	insertQuery := `
	// 		INSERT INTO "vasst_expense".user_categories
	// 		(user_category_id, user_id, category_id, name, description, icon, is_custom, is_active)
	// 		VALUES ` + strings.Join(valueStrings, ",")
	// 	_, err = r.DB.ExecContext(ctx, insertQuery, valueArgs...)
	// 	if err != nil {
	// 		return entities.User{}, err
	// 	}
	// }

	return createdUser, err
}

// Update updates a user
func (r *userRepository) Update(ctx context.Context, user *entities.User) (entities.User, error) {
	query := `
		UPDATE "vasst_expense".users
		SET email = $2,
			phone_number = $3,
			password_hash = $4,
			first_name = $5,
			last_name = $6,
			timezone = $7,
			currency_id = $8,
			subscription_plan_id = $9,
			email_verified_at = $10,
			phone_verified_at = $11,
			status = $12,
			updated_at = CURRENT_TIMESTAMP
		WHERE user_id = $1
		RETURNING user_id, email, phone_number, password_hash, first_name, last_name, 
		          timezone, currency_id, subscription_plan_id, email_verified_at, 
		          phone_verified_at, status, created_at, updated_at
	`

	// Handle nullable timestamp fields
	var emailVerifiedAt, phoneVerifiedAt sql.NullTime
	if user.EmailVerifiedAt != nil {
		emailVerifiedAt.Time = *user.EmailVerifiedAt
		emailVerifiedAt.Valid = true
	}
	if user.PhoneVerifiedAt != nil {
		phoneVerifiedAt.Time = *user.PhoneVerifiedAt
		phoneVerifiedAt.Valid = true
	}

	// Handle email field - convert empty string to NULL
	var email sql.NullString
	if user.Email != "" {
		email.String = user.Email
		email.Valid = true
	}

	var updatedUser entities.User
	var returnedEmailVerifiedAt, returnedPhoneVerifiedAt sql.NullTime

	err := r.DB.QueryRowContext(ctx, query,
		user.UserID,
		email,
		user.PhoneNumber,
		user.PasswordHash,
		user.FirstName,
		user.LastName,
		user.Timezone,
		user.CurrencyID,
		user.SubscriptionPlanID,
		emailVerifiedAt,
		phoneVerifiedAt,
		user.Status,
	).Scan(
		&updatedUser.UserID,
		&updatedUser.Email,
		&updatedUser.PhoneNumber,
		&updatedUser.PasswordHash,
		&updatedUser.FirstName,
		&updatedUser.LastName,
		&updatedUser.Timezone,
		&updatedUser.CurrencyID,
		&updatedUser.SubscriptionPlanID,
		&returnedEmailVerifiedAt,
		&returnedPhoneVerifiedAt,
		&updatedUser.Status,
		&updatedUser.CreatedAt,
		&updatedUser.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return entities.User{}, sql.ErrNoRows
		}
		return entities.User{}, err
	}

	// Handle nullable timestamp fields
	if returnedEmailVerifiedAt.Valid {
		updatedUser.EmailVerifiedAt = &returnedEmailVerifiedAt.Time
	}
	if returnedPhoneVerifiedAt.Valid {
		updatedUser.PhoneVerifiedAt = &returnedPhoneVerifiedAt.Time
	}

	return updatedUser, nil
}

// Delete deletes a user
func (r *userRepository) Delete(ctx context.Context, userID uuid.UUID) error {
	query := `
		DELETE FROM "vasst_expense".users
		WHERE user_id = $1
	`

	result, err := r.DB.ExecContext(ctx, query, userID)
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

// ListAll returns all users (global, superadmin)
func (r *userRepository) ListAll(ctx context.Context, limit, offset int) ([]*entities.User, error) {
	query := `
		SELECT user_id, email, phone_number, first_name, last_name, 
			   timezone, currency_id, subscription_plan_id, email_verified_at, 
			   phone_verified_at, status, created_at, updated_at
		FROM "vasst_expense".users
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.DB.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*entities.User
	for rows.Next() {
		var user entities.User
		var emailVerifiedAt, phoneVerifiedAt sql.NullTime
		var email sql.NullString

		err := rows.Scan(
			&user.UserID,
			&email,
			&user.PhoneNumber,
			&user.FirstName,
			&user.LastName,
			&user.Timezone,
			&user.CurrencyID,
			&user.SubscriptionPlanID,
			&emailVerifiedAt,
			&phoneVerifiedAt,
			&user.Status,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Handle nullable fields
		if email.Valid {
			user.Email = email.String
		}
		if emailVerifiedAt.Valid {
			user.EmailVerifiedAt = &emailVerifiedAt.Time
		}
		if phoneVerifiedAt.Valid {
			user.PhoneVerifiedAt = &phoneVerifiedAt.Time
		}

		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// FindByID returns a user by ID (global)
func (r *userRepository) FindByID(ctx context.Context, userID uuid.UUID) (*entities.User, error) {
	query := `
		SELECT user_id, email, phone_number, password_hash, first_name, last_name, 
		       timezone, currency_id, subscription_plan_id, email_verified_at, 
		       phone_verified_at, status, created_at, updated_at
		FROM "vasst_expense".users
		WHERE user_id = $1
	`

	var user entities.User
	var emailVerifiedAt, phoneVerifiedAt sql.NullTime
	var email sql.NullString

	err := r.DB.QueryRowContext(ctx, query, userID).Scan(
		&user.UserID,
		&email,
		&user.PhoneNumber,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.Timezone,
		&user.CurrencyID,
		&user.SubscriptionPlanID,
		&emailVerifiedAt,
		&phoneVerifiedAt,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// Handle nullable fields
	if email.Valid {
		user.Email = email.String
	}
	if emailVerifiedAt.Valid {
		user.EmailVerifiedAt = &emailVerifiedAt.Time
	}
	if phoneVerifiedAt.Valid {
		user.PhoneVerifiedAt = &phoneVerifiedAt.Time
	}

	return &user, nil
}

// FindByEmail returns a user by email (global)
func (r *userRepository) FindByEmail(ctx context.Context, email string) (*entities.User, error) {
	query := `
		SELECT user_id, email, phone_number, password_hash, first_name, last_name, 
		       timezone, currency_id, subscription_plan_id, email_verified_at, 
		       phone_verified_at, status, created_at, updated_at
		FROM "vasst_expense".users
		WHERE email = $1
	`

	var user entities.User
	var emailVerifiedAt, phoneVerifiedAt sql.NullTime
	var returnedEmail sql.NullString

	err := r.DB.QueryRowContext(ctx, query, email).Scan(
		&user.UserID,
		&returnedEmail,
		&user.PhoneNumber,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.Timezone,
		&user.CurrencyID,
		&user.SubscriptionPlanID,
		&emailVerifiedAt,
		&phoneVerifiedAt,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// Handle nullable fields
	if returnedEmail.Valid {
		user.Email = returnedEmail.String
	}
	if emailVerifiedAt.Valid {
		user.EmailVerifiedAt = &emailVerifiedAt.Time
	}
	if phoneVerifiedAt.Valid {
		user.PhoneVerifiedAt = &phoneVerifiedAt.Time
	}

	return &user, nil
}

// FindByPhoneNumber returns a user by phone number (global)
func (r *userRepository) FindByPhoneNumber(ctx context.Context, phoneNumber string) (*entities.User, error) {
	query := `
		SELECT user_id, email, phone_number, password_hash, first_name, last_name, 
		       timezone, currency_id, subscription_plan_id, email_verified_at, 
		       phone_verified_at, status, created_at, updated_at
		FROM "vasst_expense".users
		WHERE phone_number = $1
	`

	var user entities.User
	var emailVerifiedAt, phoneVerifiedAt sql.NullTime
	var email sql.NullString

	err := r.DB.QueryRowContext(ctx, query, phoneNumber).Scan(
		&user.UserID,
		&email,
		&user.PhoneNumber,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.Timezone,
		&user.CurrencyID,
		&user.SubscriptionPlanID,
		&emailVerifiedAt,
		&phoneVerifiedAt,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// Handle nullable fields
	if email.Valid {
		user.Email = email.String
	}
	if emailVerifiedAt.Valid {
		user.EmailVerifiedAt = &emailVerifiedAt.Time
	}
	if phoneVerifiedAt.Valid {
		user.PhoneVerifiedAt = &phoneVerifiedAt.Time
	}

	return &user, nil
}
