package repositories

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/vasst-id/vasst-expense-api/internal/entities"
	"github.com/vasst-id/vasst-expense-api/internal/utils/postgres"
)

type (
	leadsRepository struct {
		*postgres.Postgres
	}

	// LeadsRepository defines methods for interacting with leads in the database
	LeadsRepository interface {
		Create(ctx context.Context, lead *entities.Lead) error
		Update(ctx context.Context, lead *entities.Lead) error
		FindByID(ctx context.Context, leadID uuid.UUID) (*entities.Lead, error)
		ListAll(ctx context.Context, limit, offset int) ([]*entities.Lead, error)
		FindByPhoneNumber(ctx context.Context, phoneNumber string) (*entities.Lead, error)
		FindByEmail(ctx context.Context, email string) (*entities.Lead, error)
	}
)

// NewLeadsRepository creates a new LeadsRepository
func NewLeadsRepository(pg *postgres.Postgres) LeadsRepository {
	return &leadsRepository{pg}
}

// Create creates a new lead
func (r *leadsRepository) Create(ctx context.Context, lead *entities.Lead) error {
	query := `
		INSERT INTO "vasst_ca".leads (lead_id, name, phone_number, email, business_name, business_address, business_phone_number, business_email, business_website, business_industry, business_size, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`

	_, err := r.DB.ExecContext(ctx, query,
		lead.LeadID,
		lead.Name,
		lead.PhoneNumber,
		lead.Email,
		lead.BusinessName,
		lead.BusinessAddress,
		lead.BusinessPhoneNumber,
		lead.BusinessEmail,
		lead.BusinessWebsite,
		lead.BusinessIndustry,
		lead.BusinessSize,
	)

	return err
}

// Update updates a lead
func (r *leadsRepository) Update(ctx context.Context, lead *entities.Lead) error {
	query := `
		UPDATE "vasst_ca".leads
		SET name = $1,
			phone_number = $2,
			email = $3,
			business_name = $4,
			business_address = $5,
			business_phone_number = $6,
			business_email = $7,
			business_website = $8,
			business_industry = $9,
			business_size = $10,
			updated_at = CURRENT_TIMESTAMP
		WHERE lead_id = $11
	`

	result, err := r.DB.ExecContext(ctx, query,
		lead.Name,
		lead.PhoneNumber,
		lead.Email,
		lead.BusinessName,
		lead.BusinessAddress,
		lead.BusinessPhoneNumber,
		lead.BusinessEmail,
		lead.BusinessWebsite,
		lead.BusinessIndustry,
		lead.BusinessSize,
		lead.LeadID,
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

// FindByID returns a lead by ID
func (r *leadsRepository) FindByID(ctx context.Context, leadID uuid.UUID) (*entities.Lead, error) {
	query := `
		SELECT lead_id, name, phone_number, email, business_name, business_address, business_phone_number, business_email, business_website, business_industry, business_size, created_at, updated_at
		FROM "vasst_ca".leads
		WHERE lead_id = $1
	`

	var lead entities.Lead
	err := r.DB.QueryRowContext(ctx, query, leadID).Scan(
		&lead.LeadID,
		&lead.Name,
		&lead.PhoneNumber,
		&lead.Email,
		&lead.BusinessName,
		&lead.BusinessAddress,
		&lead.BusinessPhoneNumber,
		&lead.BusinessEmail,
		&lead.BusinessWebsite,
		&lead.BusinessIndustry,
		&lead.BusinessSize,
		&lead.CreatedAt,
		&lead.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &lead, nil
}

// ListAll returns all leads with pagination
func (r *leadsRepository) ListAll(ctx context.Context, limit, offset int) ([]*entities.Lead, error) {
	query := `
		SELECT lead_id, name, phone_number, email, business_name, business_address, business_phone_number, business_email, business_website, business_industry, business_size, created_at, updated_at
		FROM "vasst_ca".leads
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.DB.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var leads []*entities.Lead
	for rows.Next() {
		var lead entities.Lead
		err := rows.Scan(
			&lead.LeadID,
			&lead.Name,
			&lead.PhoneNumber,
			&lead.Email,
			&lead.BusinessName,
			&lead.BusinessAddress,
			&lead.BusinessPhoneNumber,
			&lead.BusinessEmail,
			&lead.BusinessWebsite,
			&lead.BusinessIndustry,
			&lead.BusinessSize,
			&lead.CreatedAt,
			&lead.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		leads = append(leads, &lead)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return leads, nil
}

// FindByPhoneNumber returns a lead by phone number
func (r *leadsRepository) FindByPhoneNumber(ctx context.Context, phoneNumber string) (*entities.Lead, error) {
	query := `
		SELECT lead_id, name, phone_number, email, business_name, business_address, business_phone_number, business_email, business_website, business_industry, business_size, created_at, updated_at
		FROM "vasst_ca".leads
		WHERE phone_number = $1
	`

	var lead entities.Lead
	err := r.DB.QueryRowContext(ctx, query, phoneNumber).Scan(
		&lead.LeadID,
		&lead.Name,
		&lead.PhoneNumber,
		&lead.Email,
		&lead.BusinessName,
		&lead.BusinessAddress,
		&lead.BusinessPhoneNumber,
		&lead.BusinessEmail,
		&lead.BusinessWebsite,
		&lead.BusinessIndustry,
		&lead.BusinessSize,
		&lead.CreatedAt,
		&lead.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &lead, nil
}

// FindByEmail returns a lead by email
func (r *leadsRepository) FindByEmail(ctx context.Context, email string) (*entities.Lead, error) {
	query := `
		SELECT lead_id, name, phone_number, email, business_name, business_address, business_phone_number, business_email, business_website, business_industry, business_size, created_at, updated_at
		FROM "vasst_ca".leads
		WHERE email = $1
	`

	var lead entities.Lead
	err := r.DB.QueryRowContext(ctx, query, email).Scan(
		&lead.LeadID,
		&lead.Name,
		&lead.PhoneNumber,
		&lead.Email,
		&lead.BusinessName,
		&lead.BusinessAddress,
		&lead.BusinessPhoneNumber,
		&lead.BusinessEmail,
		&lead.BusinessWebsite,
		&lead.BusinessIndustry,
		&lead.BusinessSize,
		&lead.CreatedAt,
		&lead.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &lead, nil
}
