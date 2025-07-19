package repositories

import (
	"context"
	"database/sql"

	"github.com/vasst-id/vasst-expense-api/internal/entities"
	"github.com/vasst-id/vasst-expense-api/internal/utils/postgres"
)

type (
	taxonomyRepository struct {
		*postgres.Postgres
	}

	TaxonomyRepository interface {
		Create(ctx context.Context, taxonomy *entities.Taxonomy) (entities.Taxonomy, error)
		Update(ctx context.Context, taxonomy *entities.Taxonomy) (entities.Taxonomy, error)
		Delete(ctx context.Context, taxonomyID int) error
		FindByID(ctx context.Context, taxonomyID int) (*entities.Taxonomy, error)
		FindByType(ctx context.Context, taxonomyType string, limit, offset int) ([]*entities.Taxonomy, error)
		FindByTypeAndValue(ctx context.Context, taxonomyType, value string) (*entities.Taxonomy, error)
		FindActive(ctx context.Context, limit, offset int) ([]*entities.Taxonomy, error)
		FindAll(ctx context.Context, limit, offset int) ([]*entities.Taxonomy, error)
	}
)

// NewTaxonomyRepository creates a new TaxonomyRepository
func NewTaxonomyRepository(pg *postgres.Postgres) TaxonomyRepository {
	return &taxonomyRepository{pg}
}

// Create creates a new taxonomy
func (r *taxonomyRepository) Create(ctx context.Context, taxonomy *entities.Taxonomy) (entities.Taxonomy, error) {
	query := `
		INSERT INTO "vasst_expense".taxonomy 
		(label, value, type, type_label, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING taxonomy_id, label, value, type, type_label, status, created_at, updated_at
	`

	var createdTaxonomy entities.Taxonomy
	err := r.DB.QueryRowContext(ctx, query,
		taxonomy.Label, taxonomy.Value, taxonomy.Type, taxonomy.TypeLabel, taxonomy.Status,
	).Scan(
		&createdTaxonomy.TaxonomyID, &createdTaxonomy.Label, &createdTaxonomy.Value,
		&createdTaxonomy.Type, &createdTaxonomy.TypeLabel, &createdTaxonomy.Status,
		&createdTaxonomy.CreatedAt, &createdTaxonomy.UpdatedAt,
	)

	return createdTaxonomy, err
}

// Update updates a taxonomy
func (r *taxonomyRepository) Update(ctx context.Context, taxonomy *entities.Taxonomy) (entities.Taxonomy, error) {
	query := `
		UPDATE "vasst_expense".taxonomy 
		SET label = $2, value = $3, type = $4, type_label = $5, status = $6, updated_at = CURRENT_TIMESTAMP
		WHERE taxonomy_id = $1
		RETURNING taxonomy_id, label, value, type, type_label, status, created_at, updated_at
	`

	var updatedTaxonomy entities.Taxonomy
	err := r.DB.QueryRowContext(ctx, query,
		taxonomy.TaxonomyID, taxonomy.Label, taxonomy.Value, taxonomy.Type,
		taxonomy.TypeLabel, taxonomy.Status,
	).Scan(
		&updatedTaxonomy.TaxonomyID, &updatedTaxonomy.Label, &updatedTaxonomy.Value,
		&updatedTaxonomy.Type, &updatedTaxonomy.TypeLabel, &updatedTaxonomy.Status,
		&updatedTaxonomy.CreatedAt, &updatedTaxonomy.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return entities.Taxonomy{}, sql.ErrNoRows // Taxonomy not found
		}
		return entities.Taxonomy{}, err
	}

	return updatedTaxonomy, nil
}

// Delete soft deletes a taxonomy (sets status to inactive)
func (r *taxonomyRepository) Delete(ctx context.Context, taxonomyID int) error {
	query := `
		UPDATE "vasst_expense".taxonomy 
		SET status = $2, updated_at = CURRENT_TIMESTAMP
		WHERE taxonomy_id = $1
	`

	result, err := r.DB.ExecContext(ctx, query, taxonomyID, entities.TaxonomyStatusInactive)
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

// FindByID finds a taxonomy by ID
func (r *taxonomyRepository) FindByID(ctx context.Context, taxonomyID int) (*entities.Taxonomy, error) {
	query := `
		SELECT taxonomy_id, label, value, type, type_label, status, created_at, updated_at
		FROM "vasst_expense".taxonomy 
		WHERE taxonomy_id = $1
	`

	var taxonomy entities.Taxonomy
	err := r.DB.QueryRowContext(ctx, query, taxonomyID).Scan(
		&taxonomy.TaxonomyID, &taxonomy.Label, &taxonomy.Value,
		&taxonomy.Type, &taxonomy.TypeLabel, &taxonomy.Status,
		&taxonomy.CreatedAt, &taxonomy.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &taxonomy, nil
}

// FindByType finds taxonomies by type with pagination
func (r *taxonomyRepository) FindByType(ctx context.Context, taxonomyType string, limit, offset int) ([]*entities.Taxonomy, error) {
	query := `
		SELECT taxonomy_id, label, value, type, type_label, status, created_at, updated_at
		FROM "vasst_expense".taxonomy 
		WHERE type = $1 AND status = $2
		ORDER BY label ASC
		LIMIT $3 OFFSET $4
	`

	rows, err := r.DB.QueryContext(ctx, query, taxonomyType, entities.TaxonomyStatusActive, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var taxonomies []*entities.Taxonomy
	for rows.Next() {
		var taxonomy entities.Taxonomy
		err := rows.Scan(
			&taxonomy.TaxonomyID, &taxonomy.Label, &taxonomy.Value,
			&taxonomy.Type, &taxonomy.TypeLabel, &taxonomy.Status,
			&taxonomy.CreatedAt, &taxonomy.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		taxonomies = append(taxonomies, &taxonomy)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return taxonomies, nil
}

// FindByTypeAndValue finds a taxonomy by type and value
func (r *taxonomyRepository) FindByTypeAndValue(ctx context.Context, taxonomyType, value string) (*entities.Taxonomy, error) {
	query := `
		SELECT taxonomy_id, label, value, type, type_label, status, created_at, updated_at
		FROM "vasst_expense".taxonomy 
		WHERE type = $1 AND value = $2 AND status = $3
	`

	var taxonomy entities.Taxonomy
	err := r.DB.QueryRowContext(ctx, query, taxonomyType, value, entities.TaxonomyStatusActive).Scan(
		&taxonomy.TaxonomyID, &taxonomy.Label, &taxonomy.Value,
		&taxonomy.Type, &taxonomy.TypeLabel, &taxonomy.Status,
		&taxonomy.CreatedAt, &taxonomy.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &taxonomy, nil
}

// FindActive finds all active taxonomies with pagination
func (r *taxonomyRepository) FindActive(ctx context.Context, limit, offset int) ([]*entities.Taxonomy, error) {
	query := `
		SELECT taxonomy_id, label, value, type, type_label, status, created_at, updated_at
		FROM "vasst_expense".taxonomy 
		WHERE status = $1
		ORDER BY type ASC, label ASC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.DB.QueryContext(ctx, query, entities.TaxonomyStatusActive, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var taxonomies []*entities.Taxonomy
	for rows.Next() {
		var taxonomy entities.Taxonomy
		err := rows.Scan(
			&taxonomy.TaxonomyID, &taxonomy.Label, &taxonomy.Value,
			&taxonomy.Type, &taxonomy.TypeLabel, &taxonomy.Status,
			&taxonomy.CreatedAt, &taxonomy.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		taxonomies = append(taxonomies, &taxonomy)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return taxonomies, nil
}

// FindAll finds all taxonomies with pagination
func (r *taxonomyRepository) FindAll(ctx context.Context, limit, offset int) ([]*entities.Taxonomy, error) {
	query := `
		SELECT taxonomy_id, label, value, type, type_label, status, created_at, updated_at
		FROM "vasst_expense".taxonomy 
		ORDER BY type ASC, label ASC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.DB.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var taxonomies []*entities.Taxonomy
	for rows.Next() {
		var taxonomy entities.Taxonomy
		err := rows.Scan(
			&taxonomy.TaxonomyID, &taxonomy.Label, &taxonomy.Value,
			&taxonomy.Type, &taxonomy.TypeLabel, &taxonomy.Status,
			&taxonomy.CreatedAt, &taxonomy.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		taxonomies = append(taxonomies, &taxonomy)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return taxonomies, nil
}
