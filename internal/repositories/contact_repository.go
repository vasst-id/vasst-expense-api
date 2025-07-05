package repositories

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/vasst-id/vasst-expense-api/internal/entities"
	"github.com/vasst-id/vasst-expense-api/internal/utils/postgres"
)

type (
	contactRepository struct {
		*postgres.Postgres
	}

	// ContactRepository defines methods for interacting with contacts in the database
	ContactRepository interface {
		// Create creates a new contact
		Create(ctx context.Context, contact *entities.Contact) error

		// Update updates a contact
		Update(ctx context.Context, contact *entities.Contact) error

		// Delete deletes a contact
		Delete(ctx context.Context, contactID uuid.UUID) error

		// List returns all contacts with optional filtering
		List(ctx context.Context, limit, offset int) ([]*entities.Contact, error)

		// FindByID returns a contact by ID
		FindByID(ctx context.Context, contactID uuid.UUID) (*entities.Contact, error)

		// FindByPhoneNumber returns a contact by phone number
		FindByPhoneNumber(ctx context.Context, phoneNumber string) (*entities.Contact, error)

		// ListByOrganizationID returns all contacts for an organization
		ListByOrganizationID(ctx context.Context, organizationID uuid.UUID, limit, offset int) ([]*entities.Contact, error)
	}
)

// NewContactRepository creates a new ContactRepository
func NewContactRepository(pg *postgres.Postgres) ContactRepository {
	return &contactRepository{pg}
}

// Create creates a new contact
func (r *contactRepository) Create(ctx context.Context, contact *entities.Contact) error {

	// Get default context structure from organization setting
	organizationSetting, err := r.GetOrganizationSettingContactInfoStructure(ctx, contact.OrganizationID)
	if err != nil {
		return err
	}

	if organizationSetting.ContactInfoStructure != nil {
		contact.Context = organizationSetting.ContactInfoStructure
	}

	query := `
		INSERT INTO "vasst_ca".contact (
			contact_id, organization_id, contact_name, phone_number, email, 
			salutation, notes, custom_system_prompt, context, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`

	_, err = r.DB.ExecContext(ctx, query,
		contact.ContactID,
		contact.OrganizationID,
		contact.Name,
		contact.PhoneNumber,
		contact.Email,
		contact.Salutation,
		contact.Notes,
		contact.CustomSystemPrompt,
		contact.Context,
	)

	return err
}

// Update updates a contact
func (r *contactRepository) Update(ctx context.Context, contact *entities.Contact) error {
	query := `
		UPDATE "vasst_ca".contact
		SET contact_name = $1,
			phone_number = $2,
			email = $3,
			salutation = $4,
			notes = $5,
			custom_system_prompt = $6,
			context = $7,
			updated_at = CURRENT_TIMESTAMP
		WHERE contact_id = $8
	`

	result, err := r.DB.ExecContext(ctx, query,
		contact.Name,
		contact.PhoneNumber,
		contact.Email,
		contact.Salutation,
		contact.Notes,
		contact.CustomSystemPrompt,
		contact.Context,
		contact.ContactID,
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

// Delete deletes a contact
func (r *contactRepository) Delete(ctx context.Context, contactID uuid.UUID) error {
	query := `
		DELETE FROM "vasst_ca".contact
		WHERE contact_id = $1
	`

	result, err := r.DB.ExecContext(ctx, query, contactID)
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

// List returns all contacts with optional filtering
func (r *contactRepository) List(ctx context.Context, limit, offset int) ([]*entities.Contact, error) {
	query := `
		SELECT contact_id, organization_id, contact_name, phone_number, email, 
			   salutation, notes, custom_system_prompt, context, created_at, updated_at
		FROM "vasst_ca".contact
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.DB.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contacts []*entities.Contact
	for rows.Next() {
		var contact entities.Contact
		var customSystemPrompt sql.NullString
		var notes sql.NullString
		var salutation sql.NullString
		var context sql.NullString

		err := rows.Scan(
			&contact.ContactID,
			&contact.OrganizationID,
			&contact.Name,
			&contact.PhoneNumber,
			&contact.Email,
			&salutation,
			&notes,
			&customSystemPrompt,
			&context,
			&contact.CreatedAt,
			&contact.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Handle nullable fields
		if customSystemPrompt.Valid {
			contact.CustomSystemPrompt = customSystemPrompt.String
		}
		if notes.Valid {
			contact.Notes = notes.String
		}
		if salutation.Valid {
			contact.Salutation = salutation.String
		}
		if context.Valid {
			contact.Context = json.RawMessage(context.String)
		}

		contacts = append(contacts, &contact)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return contacts, nil
}

// FindByID returns a contact by ID
func (r *contactRepository) FindByID(ctx context.Context, contactID uuid.UUID) (*entities.Contact, error) {
	query := `
		SELECT contact_id, organization_id, contact_name, phone_number, email, 
			   salutation, notes, custom_system_prompt, context, created_at, updated_at
		FROM "vasst_ca".contact
		WHERE contact_id = $1
	`

	var contact entities.Contact
	var customSystemPrompt sql.NullString
	var notes sql.NullString
	var salutation sql.NullString
	var context sql.NullString

	err := r.DB.QueryRowContext(ctx, query, contactID).Scan(
		&contact.ContactID,
		&contact.OrganizationID,
		&contact.Name,
		&contact.PhoneNumber,
		&contact.Email,
		&salutation,
		&notes,
		&customSystemPrompt,
		&context,
		&contact.CreatedAt,
		&contact.UpdatedAt,
	)

	// Handle nullable fields
	if customSystemPrompt.Valid {
		contact.CustomSystemPrompt = customSystemPrompt.String
	}
	if notes.Valid {
		contact.Notes = notes.String
	}
	if salutation.Valid {
		contact.Salutation = salutation.String
	}
	if context.Valid {
		contact.Context = json.RawMessage(context.String)
	}
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &contact, nil
}

// FindByPhoneNumber returns a contact by phone number
func (r *contactRepository) FindByPhoneNumber(ctx context.Context, phoneNumber string) (*entities.Contact, error) {
	query := `
		SELECT contact_id, organization_id, contact_name, phone_number, email, 
			   salutation, notes, custom_system_prompt, context, created_at, updated_at
		FROM "vasst_ca".contact
		WHERE phone_number = $1
	`

	var contact entities.Contact
	var customSystemPrompt sql.NullString
	var notes sql.NullString
	var salutation sql.NullString
	var context sql.NullString

	err := r.DB.QueryRowContext(ctx, query, phoneNumber).Scan(
		&contact.ContactID,
		&contact.OrganizationID,
		&contact.Name,
		&contact.PhoneNumber,
		&contact.Email,
		&salutation,
		&notes,
		&customSystemPrompt,
		&context,
		&contact.CreatedAt,
		&contact.UpdatedAt,
	)

	// Handle nullable fields
	if customSystemPrompt.Valid {
		contact.CustomSystemPrompt = customSystemPrompt.String
	}
	if notes.Valid {
		contact.Notes = notes.String
	}
	if salutation.Valid {
		contact.Salutation = salutation.String
	}
	if context.Valid {
		contact.Context = json.RawMessage(context.String)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &contact, nil
}

// ListByOrganizationID returns all contacts for an organization
func (r *contactRepository) ListByOrganizationID(ctx context.Context, organizationID uuid.UUID, limit, offset int) ([]*entities.Contact, error) {
	query := `
		SELECT contact_id, organization_id, contact_name, phone_number, email, 
			   salutation, notes, custom_system_prompt, context, created_at, updated_at
		FROM "vasst_ca".contact
		WHERE organization_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.DB.QueryContext(ctx, query, organizationID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contacts []*entities.Contact
	for rows.Next() {
		var contact entities.Contact
		var customSystemPrompt sql.NullString
		var notes sql.NullString
		var salutation sql.NullString
		var context sql.NullString

		err := rows.Scan(
			&contact.ContactID,
			&contact.OrganizationID,
			&contact.Name,
			&contact.PhoneNumber,
			&contact.Email,
			&salutation,
			&notes,
			&customSystemPrompt,
			&context,
			&contact.CreatedAt,
			&contact.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Handle nullable fields
		if customSystemPrompt.Valid {
			contact.CustomSystemPrompt = customSystemPrompt.String
		}
		if notes.Valid {
			contact.Notes = notes.String
		}
		if salutation.Valid {
			contact.Salutation = salutation.String
		}
		if context.Valid {
			contact.Context = json.RawMessage(context.String)
		}
		if err != nil {
			return nil, err
		}
		contacts = append(contacts, &contact)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return contacts, nil
}

func (r *contactRepository) GetOrganizationSettingContactInfoStructure(ctx context.Context, organizationID uuid.UUID) (*entities.OrganizationSetting, error) {
	query := `
			SELECT contact_info_structure
			FROM "vasst_ca".organization_setting
			WHERE organization_id = $1
		`

	var organizationSetting entities.OrganizationSetting
	err := r.DB.QueryRowContext(ctx, query, organizationID).Scan(
		&organizationSetting.ContactInfoStructure,
	)

	if err != nil {
		return nil, err
	}

	return &organizationSetting, nil
}
