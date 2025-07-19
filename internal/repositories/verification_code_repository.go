package repositories

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/vasst-id/vasst-expense-api/internal/entities"
	"github.com/vasst-id/vasst-expense-api/internal/utils/postgres"
)

type (
	verificationCodeRepository struct {
		*postgres.Postgres
	}

	VerificationCodeRepository interface {
		Create(ctx context.Context, verificationCode *entities.VerificationCode) (entities.VerificationCode, error)
		Update(ctx context.Context, verificationCode *entities.VerificationCode) (entities.VerificationCode, error)
		Delete(ctx context.Context, verificationCodeID uuid.UUID) error
		FindByID(ctx context.Context, verificationCodeID uuid.UUID) (*entities.VerificationCode, error)
		FindByPhoneNumberAndType(ctx context.Context, phoneNumber, codeType string) (*entities.VerificationCode, error)
		FindActiveByPhoneNumberAndType(ctx context.Context, phoneNumber, codeType string) (*entities.VerificationCode, error)
		CleanupExpiredCodes(ctx context.Context) error
		IncrementAttempts(ctx context.Context, verificationCodeID uuid.UUID) error
		MarkAsUsed(ctx context.Context, verificationCodeID uuid.UUID) error
	}
)

// NewVerificationCodeRepository creates a new VerificationCodeRepository
func NewVerificationCodeRepository(pg *postgres.Postgres) VerificationCodeRepository {
	return &verificationCodeRepository{pg}
}

// Create creates a new verification code
func (r *verificationCodeRepository) Create(ctx context.Context, verificationCode *entities.VerificationCode) (entities.VerificationCode, error) {
	query := `
		INSERT INTO "vasst_expense".verification_codes (
			verification_code_id, phone_number, code, code_type, expires_at, 
			is_used, attempts_count, max_attempts, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING verification_code_id, phone_number, code, code_type, expires_at, 
		          is_used, attempts_count, max_attempts, created_at, updated_at
	`

	var createdVerificationCode entities.VerificationCode
	err := r.DB.QueryRowContext(ctx, query,
		verificationCode.VerificationCodeID,
		verificationCode.PhoneNumber,
		verificationCode.Code,
		verificationCode.CodeType,
		verificationCode.ExpiresAt,
		verificationCode.IsUsed,
		verificationCode.AttemptsCount,
		verificationCode.MaxAttempts,
	).Scan(
		&createdVerificationCode.VerificationCodeID,
		&createdVerificationCode.PhoneNumber,
		&createdVerificationCode.Code,
		&createdVerificationCode.CodeType,
		&createdVerificationCode.ExpiresAt,
		&createdVerificationCode.IsUsed,
		&createdVerificationCode.AttemptsCount,
		&createdVerificationCode.MaxAttempts,
		&createdVerificationCode.CreatedAt,
		&createdVerificationCode.UpdatedAt,
	)

	return createdVerificationCode, err
}

// Update updates a verification code
func (r *verificationCodeRepository) Update(ctx context.Context, verificationCode *entities.VerificationCode) (entities.VerificationCode, error) {
	query := `
		UPDATE "vasst_expense".verification_codes
		SET phone_number = $2,
			code = $3,
			code_type = $4,
			expires_at = $5,
			is_used = $6,
			attempts_count = $7,
			max_attempts = $8,
			updated_at = CURRENT_TIMESTAMP
		WHERE verification_code_id = $1
		RETURNING verification_code_id, phone_number, code, code_type, expires_at, 
		          is_used, attempts_count, max_attempts, created_at, updated_at
	`

	var updatedVerificationCode entities.VerificationCode
	err := r.DB.QueryRowContext(ctx, query,
		verificationCode.VerificationCodeID,
		verificationCode.PhoneNumber,
		verificationCode.Code,
		verificationCode.CodeType,
		verificationCode.ExpiresAt,
		verificationCode.IsUsed,
		verificationCode.AttemptsCount,
		verificationCode.MaxAttempts,
	).Scan(
		&updatedVerificationCode.VerificationCodeID,
		&updatedVerificationCode.PhoneNumber,
		&updatedVerificationCode.Code,
		&updatedVerificationCode.CodeType,
		&updatedVerificationCode.ExpiresAt,
		&updatedVerificationCode.IsUsed,
		&updatedVerificationCode.AttemptsCount,
		&updatedVerificationCode.MaxAttempts,
		&updatedVerificationCode.CreatedAt,
		&updatedVerificationCode.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return entities.VerificationCode{}, sql.ErrNoRows
		}
		return entities.VerificationCode{}, err
	}

	return updatedVerificationCode, nil
}

// Delete deletes a verification code
func (r *verificationCodeRepository) Delete(ctx context.Context, verificationCodeID uuid.UUID) error {
	query := `
		DELETE FROM "vasst_expense".verification_codes
		WHERE verification_code_id = $1
	`

	result, err := r.DB.ExecContext(ctx, query, verificationCodeID)
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

// FindByID returns a verification code by ID
func (r *verificationCodeRepository) FindByID(ctx context.Context, verificationCodeID uuid.UUID) (*entities.VerificationCode, error) {
	query := `
		SELECT verification_code_id, phone_number, code, code_type, expires_at, 
		       is_used, attempts_count, max_attempts, created_at, updated_at
		FROM "vasst_expense".verification_codes
		WHERE verification_code_id = $1
	`

	var verificationCode entities.VerificationCode
	err := r.DB.QueryRowContext(ctx, query, verificationCodeID).Scan(
		&verificationCode.VerificationCodeID,
		&verificationCode.PhoneNumber,
		&verificationCode.Code,
		&verificationCode.CodeType,
		&verificationCode.ExpiresAt,
		&verificationCode.IsUsed,
		&verificationCode.AttemptsCount,
		&verificationCode.MaxAttempts,
		&verificationCode.CreatedAt,
		&verificationCode.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &verificationCode, nil
}

// FindByPhoneNumberAndType returns a verification code by phone number and type
func (r *verificationCodeRepository) FindByPhoneNumberAndType(ctx context.Context, phoneNumber, codeType string) (*entities.VerificationCode, error) {
	query := `
		SELECT verification_code_id, phone_number, code, code_type, expires_at, 
		       is_used, attempts_count, max_attempts, created_at, updated_at
		FROM "vasst_expense".verification_codes
		WHERE phone_number = $1 AND code_type = $2
		ORDER BY created_at DESC
		LIMIT 1
	`

	var verificationCode entities.VerificationCode
	err := r.DB.QueryRowContext(ctx, query, phoneNumber, codeType).Scan(
		&verificationCode.VerificationCodeID,
		&verificationCode.PhoneNumber,
		&verificationCode.Code,
		&verificationCode.CodeType,
		&verificationCode.ExpiresAt,
		&verificationCode.IsUsed,
		&verificationCode.AttemptsCount,
		&verificationCode.MaxAttempts,
		&verificationCode.CreatedAt,
		&verificationCode.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &verificationCode, nil
}

// FindActiveByPhoneNumberAndType returns an active verification code by phone number and type
func (r *verificationCodeRepository) FindActiveByPhoneNumberAndType(ctx context.Context, phoneNumber, codeType string) (*entities.VerificationCode, error) {
	query := `
		SELECT verification_code_id, phone_number, code, code_type, expires_at, 
		       is_used, attempts_count, max_attempts, created_at, updated_at
		FROM "vasst_expense".verification_codes
		WHERE phone_number = $1 AND code_type = $2 AND is_used = false AND expires_at > CURRENT_TIMESTAMP
		ORDER BY created_at DESC
		LIMIT 1
	`

	var verificationCode entities.VerificationCode
	err := r.DB.QueryRowContext(ctx, query, phoneNumber, codeType).Scan(
		&verificationCode.VerificationCodeID,
		&verificationCode.PhoneNumber,
		&verificationCode.Code,
		&verificationCode.CodeType,
		&verificationCode.ExpiresAt,
		&verificationCode.IsUsed,
		&verificationCode.AttemptsCount,
		&verificationCode.MaxAttempts,
		&verificationCode.CreatedAt,
		&verificationCode.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &verificationCode, nil
}

// CleanupExpiredCodes removes expired verification codes
func (r *verificationCodeRepository) CleanupExpiredCodes(ctx context.Context) error {
	query := `
		DELETE FROM "vasst_expense".verification_codes
		WHERE expires_at < CURRENT_TIMESTAMP
	`

	_, err := r.DB.ExecContext(ctx, query)
	return err
}

// IncrementAttempts increments the attempts count for a verification code
func (r *verificationCodeRepository) IncrementAttempts(ctx context.Context, verificationCodeID uuid.UUID) error {
	query := `
		UPDATE "vasst_expense".verification_codes
		SET attempts_count = attempts_count + 1, updated_at = CURRENT_TIMESTAMP
		WHERE verification_code_id = $1
	`

	result, err := r.DB.ExecContext(ctx, query, verificationCodeID)
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

// MarkAsUsed marks a verification code as used
func (r *verificationCodeRepository) MarkAsUsed(ctx context.Context, verificationCodeID uuid.UUID) error {
	query := `
		UPDATE "vasst_expense".verification_codes
		SET is_used = true, updated_at = CURRENT_TIMESTAMP
		WHERE verification_code_id = $1
	`

	result, err := r.DB.ExecContext(ctx, query, verificationCodeID)
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
