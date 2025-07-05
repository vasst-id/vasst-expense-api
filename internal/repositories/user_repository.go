package repositories

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/vasst-id/vasst-expense-api/internal/entities"
	"github.com/vasst-id/vasst-expense-api/internal/utils/postgres"
)

type (
	userRepository struct {
		*postgres.Postgres
	}

	// UserRepository defines methods for interacting with users in the database
	UserRepository interface {
		// Global (superadmin) methods
		Create(ctx context.Context, user *entities.User) error
		Update(ctx context.Context, user *entities.User) error
		Delete(ctx context.Context, userID uuid.UUID) error
		ListAll(ctx context.Context, limit, offset int) ([]*entities.User, error)
		FindByID(ctx context.Context, userID uuid.UUID) (*entities.User, error)
		FindByPhoneNumber(ctx context.Context, phoneNumber string) (*entities.User, error)
		FindByUsername(ctx context.Context, username string) (*entities.User, error)

		// Organization-scoped methods
		ListByOrganization(ctx context.Context, organizationID uuid.UUID, limit, offset int) ([]*entities.User, error)
		FindByIDAndOrganization(ctx context.Context, userID, organizationID uuid.UUID) (*entities.User, error)
		FindByPhoneNumberAndOrganization(ctx context.Context, phoneNumber string, organizationID uuid.UUID) (*entities.User, error)
		FindByUsernameAndOrganization(ctx context.Context, username string, organizationID uuid.UUID) (*entities.User, error)
	}
)

// NewUserRepository creates a new UserRepository
func NewUserRepository(pg *postgres.Postgres) UserRepository {
	return &userRepository{pg}
}

// Create creates a new user
func (r *userRepository) Create(ctx context.Context, user *entities.User) error {
	query := `
		INSERT INTO "vasst_ca".user (user_id, organization_id, role_id, user_fullname, phone_number, username, password, is_active, access_token, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`

	// Handle access_token - if empty, set to NULL
	var accessToken sql.NullString
	if user.AccessToken != "" {
		accessToken.String = user.AccessToken
		accessToken.Valid = true
	}

	_, err := r.DB.ExecContext(ctx, query,
		user.UserID,
		user.OrganizationID,
		user.RoleID,
		user.UserFullName,
		user.PhoneNumber,
		user.Username,
		user.Password,
		user.IsActive,
		accessToken,
	)

	return err
}

// Update updates a user
func (r *userRepository) Update(ctx context.Context, user *entities.User) error {
	query := `
		UPDATE "vasst_ca".user
		SET organization_id = $1,
			role_id = $2,
			user_fullname = $3,
			phone_number = $4,
			username = $5,
			password = $6,
			is_active = $7,
			access_token = $8,
			updated_at = CURRENT_TIMESTAMP
		WHERE user_id = $9
	`

	// Handle access_token - if empty, set to NULL
	var accessToken sql.NullString
	if user.AccessToken != "" {
		accessToken.String = user.AccessToken
		accessToken.Valid = true
	}

	result, err := r.DB.ExecContext(ctx, query,
		user.OrganizationID,
		user.RoleID,
		user.UserFullName,
		user.PhoneNumber,
		user.Username,
		user.Password,
		user.IsActive,
		accessToken,
		user.UserID,
	)
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

// Delete deletes a user
func (r *userRepository) Delete(ctx context.Context, userID uuid.UUID) error {
	query := `
		DELETE FROM "vasst_ca".user
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
		SELECT user_id, organization_id, role_id, user_fullname, phone_number, username, password, is_active, access_token, created_at, updated_at
		FROM "vasst_ca".user
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
		var accessToken sql.NullString
		err := rows.Scan(
			&user.UserID,
			&user.OrganizationID,
			&user.RoleID,
			&user.UserFullName,
			&user.PhoneNumber,
			&user.Username,
			&user.Password,
			&user.IsActive,
			&accessToken,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		if accessToken.Valid {
			user.AccessToken = accessToken.String
		}
		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// ListByOrganization returns all users for a specific organization
func (r *userRepository) ListByOrganization(ctx context.Context, organizationID uuid.UUID, limit, offset int) ([]*entities.User, error) {
	query := `
		SELECT user_id, organization_id, role_id, user_fullname, phone_number, username, password, is_active, access_token, created_at, updated_at
		FROM "vasst_ca".user
		WHERE organization_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.DB.QueryContext(ctx, query, organizationID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*entities.User
	for rows.Next() {
		var user entities.User
		var accessToken sql.NullString
		err := rows.Scan(
			&user.UserID,
			&user.OrganizationID,
			&user.RoleID,
			&user.UserFullName,
			&user.PhoneNumber,
			&user.Username,
			&user.Password,
			&user.IsActive,
			&accessToken,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		if accessToken.Valid {
			user.AccessToken = accessToken.String
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
		SELECT user_id, organization_id, role_id, user_fullname, phone_number, username, password, is_active, access_token, created_at, updated_at
		FROM "vasst_ca".user
		WHERE user_id = $1
	`

	var user entities.User
	var accessToken sql.NullString
	err := r.DB.QueryRowContext(ctx, query, userID).Scan(
		&user.UserID,
		&user.OrganizationID,
		&user.RoleID,
		&user.UserFullName,
		&user.PhoneNumber,
		&user.Username,
		&user.Password,
		&user.IsActive,
		&accessToken,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if accessToken.Valid {
		user.AccessToken = accessToken.String
	}

	return &user, nil
}

// FindByIDAndOrganization returns a user by ID and organization
func (r *userRepository) FindByIDAndOrganization(ctx context.Context, userID, organizationID uuid.UUID) (*entities.User, error) {
	query := `
		SELECT user_id, organization_id, role_id, user_fullname, phone_number, username, password, is_active, access_token, created_at, updated_at
		FROM "vasst_ca".user
		WHERE user_id = $1 AND organization_id = $2
	`

	var user entities.User
	var accessToken sql.NullString
	err := r.DB.QueryRowContext(ctx, query, userID, organizationID).Scan(
		&user.UserID,
		&user.OrganizationID,
		&user.RoleID,
		&user.UserFullName,
		&user.PhoneNumber,
		&user.Username,
		&user.Password,
		&user.IsActive,
		&accessToken,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if accessToken.Valid {
		user.AccessToken = accessToken.String
	}

	return &user, nil
}

// FindByPhoneNumber returns a user by phone number (global)
func (r *userRepository) FindByPhoneNumber(ctx context.Context, phoneNumber string) (*entities.User, error) {
	query := `
		SELECT user_id, organization_id, role_id, user_fullname, phone_number, username, password, is_active, access_token, created_at, updated_at
		FROM "vasst_ca".user
		WHERE phone_number = $1
	`

	var user entities.User
	var accessToken sql.NullString
	err := r.DB.QueryRowContext(ctx, query, phoneNumber).Scan(
		&user.UserID,
		&user.OrganizationID,
		&user.RoleID,
		&user.UserFullName,
		&user.PhoneNumber,
		&user.Username,
		&user.Password,
		&user.IsActive,
		&accessToken,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if accessToken.Valid {
		user.AccessToken = accessToken.String
	}

	return &user, nil
}

// FindByPhoneNumberAndOrganization returns a user by phone number and organization
func (r *userRepository) FindByPhoneNumberAndOrganization(ctx context.Context, phoneNumber string, organizationID uuid.UUID) (*entities.User, error) {
	query := `
		SELECT user_id, organization_id, role_id, user_fullname, phone_number, username, password, is_active, access_token, created_at, updated_at
		FROM "vasst_ca".user
		WHERE phone_number = $1 AND organization_id = $2
	`

	var user entities.User
	var accessToken sql.NullString
	err := r.DB.QueryRowContext(ctx, query, phoneNumber, organizationID).Scan(
		&user.UserID,
		&user.OrganizationID,
		&user.RoleID,
		&user.UserFullName,
		&user.PhoneNumber,
		&user.Username,
		&user.Password,
		&user.IsActive,
		&accessToken,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if accessToken.Valid {
		user.AccessToken = accessToken.String
	}

	return &user, nil
}

// FindByUsername returns a user by username (global)
func (r *userRepository) FindByUsername(ctx context.Context, username string) (*entities.User, error) {
	query := `
		SELECT user_id, organization_id, role_id, user_fullname, phone_number, username, password, is_active, access_token, created_at, updated_at
		FROM "vasst_ca".user
		WHERE username = $1
	`

	var user entities.User
	var accessToken sql.NullString
	err := r.DB.QueryRowContext(ctx, query, username).Scan(
		&user.UserID,
		&user.OrganizationID,
		&user.RoleID,
		&user.UserFullName,
		&user.PhoneNumber,
		&user.Username,
		&user.Password,
		&user.IsActive,
		&accessToken,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if accessToken.Valid {
		user.AccessToken = accessToken.String
	}

	return &user, nil
}

// FindByUsernameAndOrganization returns a user by username and organization
func (r *userRepository) FindByUsernameAndOrganization(ctx context.Context, username string, organizationID uuid.UUID) (*entities.User, error) {
	query := `
		SELECT user_id, organization_id, role_id, user_fullname, phone_number, username, password, is_active, access_token, created_at, updated_at
		FROM "vasst_ca".user
		WHERE username = $1 AND organization_id = $2
	`

	var user entities.User
	var accessToken sql.NullString
	err := r.DB.QueryRowContext(ctx, query, username, organizationID).Scan(
		&user.UserID,
		&user.OrganizationID,
		&user.RoleID,
		&user.UserFullName,
		&user.PhoneNumber,
		&user.Username,
		&user.Password,
		&user.IsActive,
		&accessToken,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if accessToken.Valid {
		user.AccessToken = accessToken.String
	}

	return &user, nil
}
