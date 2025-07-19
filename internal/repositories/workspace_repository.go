package repositories

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/vasst-id/vasst-expense-api/internal/entities"
	"github.com/vasst-id/vasst-expense-api/internal/utils/postgres"
)

type (
	workspaceRepository struct {
		*postgres.Postgres
	}

	// WorkspaceRepository defines methods for interacting with workspaces in the database
	WorkspaceRepository interface {
		// Global methods
		Create(ctx context.Context, workspace *entities.Workspace) (entities.Workspace, error)
		Update(ctx context.Context, workspace *entities.Workspace) (entities.Workspace, error)
		Delete(ctx context.Context, workspaceID uuid.UUID) error
		ListAll(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entities.Workspace, error)
		FindByID(ctx context.Context, workspaceID uuid.UUID) (*entities.Workspace, error)
		FindByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entities.Workspace, error)
		FindByName(ctx context.Context, name string) (*entities.Workspace, error)
	}
)

// NewWorkspaceRepository creates a new WorkspaceRepository
func NewWorkspaceRepository(pg *postgres.Postgres) WorkspaceRepository {
	return &workspaceRepository{pg}
}

// Create creates a new workspace
func (r *workspaceRepository) Create(ctx context.Context, workspace *entities.Workspace) (entities.Workspace, error) {
	query := `
		INSERT INTO "vasst_expense".workspaces (
			workspace_id, name, description, workspace_type, icon, color_code,
			currency_id, timezone, settings, is_active, created_by, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING workspace_id, name, description, workspace_type, icon, color_code,
		          currency_id, timezone, settings, is_active, created_by, created_at, updated_at
	`

	var createdWorkspace entities.Workspace
	err := r.DB.QueryRowContext(ctx, query,
		workspace.WorkspaceID,
		workspace.Name,
		workspace.Description,
		workspace.WorkspaceType,
		workspace.Icon,
		workspace.ColorCode,
		workspace.CurrencyID,
		workspace.Timezone,
		workspace.Settings,
		workspace.IsActive,
		workspace.CreatedBy,
	).Scan(
		&createdWorkspace.WorkspaceID,
		&createdWorkspace.Name,
		&createdWorkspace.Description,
		&createdWorkspace.WorkspaceType,
		&createdWorkspace.Icon,
		&createdWorkspace.ColorCode,
		&createdWorkspace.CurrencyID,
		&createdWorkspace.Timezone,
		&createdWorkspace.Settings,
		&createdWorkspace.IsActive,
		&createdWorkspace.CreatedBy,
		&createdWorkspace.CreatedAt,
		&createdWorkspace.UpdatedAt,
	)

	return createdWorkspace, err
}

// Update updates a workspace
func (r *workspaceRepository) Update(ctx context.Context, workspace *entities.Workspace) (entities.Workspace, error) {
	query := `
		UPDATE "vasst_expense".workspaces
		SET name = $2,
			description = $3,
			workspace_type = $4,
			icon = $5,
			color_code = $6,
			currency_id = $7,
			timezone = $8,
			settings = $9,
			is_active = $10,
			updated_at = CURRENT_TIMESTAMP
		WHERE workspace_id = $1
		RETURNING workspace_id, name, description, workspace_type, icon, color_code,
		          currency_id, timezone, settings, is_active, created_by, created_at, updated_at
	`

	var updatedWorkspace entities.Workspace
	err := r.DB.QueryRowContext(ctx, query,
		workspace.WorkspaceID,
		workspace.Name,
		workspace.Description,
		workspace.WorkspaceType,
		workspace.Icon,
		workspace.ColorCode,
		workspace.CurrencyID,
		workspace.Timezone,
		workspace.Settings,
		workspace.IsActive,
	).Scan(
		&updatedWorkspace.WorkspaceID,
		&updatedWorkspace.Name,
		&updatedWorkspace.Description,
		&updatedWorkspace.WorkspaceType,
		&updatedWorkspace.Icon,
		&updatedWorkspace.ColorCode,
		&updatedWorkspace.CurrencyID,
		&updatedWorkspace.Timezone,
		&updatedWorkspace.Settings,
		&updatedWorkspace.IsActive,
		&updatedWorkspace.CreatedBy,
		&updatedWorkspace.CreatedAt,
		&updatedWorkspace.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return entities.Workspace{}, sql.ErrNoRows
		}
		return entities.Workspace{}, err
	}

	return updatedWorkspace, nil
}

// Delete deletes a workspace
func (r *workspaceRepository) Delete(ctx context.Context, workspaceID uuid.UUID) error {
	query := `
		DELETE FROM "vasst_expense".workspaces
		WHERE workspace_id = $1
	`

	result, err := r.DB.ExecContext(ctx, query, workspaceID)
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

// ListAll returns all workspaces (global, superadmin)
func (r *workspaceRepository) ListAll(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entities.Workspace, error) {
	query := `
		SELECT workspace_id, name, description, workspace_type, icon, color_code,
			   currency_id, timezone, settings, is_active, created_at, updated_at
		FROM "vasst_expense".workspaces
		WHERE created_by = $3
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.DB.QueryContext(ctx, query, limit, offset, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var workspaces []*entities.Workspace
	for rows.Next() {
		var workspace entities.Workspace

		err := rows.Scan(
			&workspace.WorkspaceID,
			&workspace.Name,
			&workspace.Description,
			&workspace.WorkspaceType,
			&workspace.Icon,
			&workspace.ColorCode,
			&workspace.CurrencyID,
			&workspace.Timezone,
			&workspace.Settings,
			&workspace.IsActive,
			&workspace.CreatedAt,
			&workspace.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		workspaces = append(workspaces, &workspace)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return workspaces, nil
}

// FindByID returns a workspace by ID
func (r *workspaceRepository) FindByID(ctx context.Context, workspaceID uuid.UUID) (*entities.Workspace, error) {
	query := `
		SELECT workspace_id, name, description, workspace_type, icon, color_code,
			   currency_id, timezone, settings, is_active, created_by, created_at, updated_at
		FROM "vasst_expense".workspaces
		WHERE workspace_id = $1
	`

	var workspace entities.Workspace

	err := r.DB.QueryRowContext(ctx, query, workspaceID).Scan(
		&workspace.WorkspaceID,
		&workspace.Name,
		&workspace.Description,
		&workspace.WorkspaceType,
		&workspace.Icon,
		&workspace.ColorCode,
		&workspace.CurrencyID,
		&workspace.Timezone,
		&workspace.Settings,
		&workspace.IsActive,
		&workspace.CreatedBy,
		&workspace.CreatedAt,
		&workspace.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &workspace, nil
}

// FindByUserID returns workspaces created by a specific user
func (r *workspaceRepository) FindByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entities.Workspace, error) {
	query := `
		SELECT workspace_id, name, description, workspace_type, icon, color_code,
			   currency_id, timezone, settings, is_active, created_by, created_at, updated_at
		FROM "vasst_expense".workspaces
		WHERE created_by = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.DB.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var workspaces []*entities.Workspace
	for rows.Next() {
		var workspace entities.Workspace

		err := rows.Scan(
			&workspace.WorkspaceID,
			&workspace.Name,
			&workspace.Description,
			&workspace.WorkspaceType,
			&workspace.Icon,
			&workspace.ColorCode,
			&workspace.CurrencyID,
			&workspace.Timezone,
			&workspace.Settings,
			&workspace.IsActive,
			&workspace.CreatedBy,
			&workspace.CreatedAt,
			&workspace.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		workspaces = append(workspaces, &workspace)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return workspaces, nil
}

// FindByName returns a workspace by name
func (r *workspaceRepository) FindByName(ctx context.Context, name string) (*entities.Workspace, error) {
	query := `
		SELECT workspace_id, name, description, workspace_type, icon, color_code,
			   currency_id, timezone, settings, is_active, created_by, created_at, updated_at
		FROM "vasst_expense".workspaces
		WHERE name = $1
	`

	var workspace entities.Workspace

	err := r.DB.QueryRowContext(ctx, query, name).Scan(
		&workspace.WorkspaceID,
		&workspace.Name,
		&workspace.Description,
		&workspace.WorkspaceType,
		&workspace.Icon,
		&workspace.ColorCode,
		&workspace.CurrencyID,
		&workspace.Timezone,
		&workspace.Settings,
		&workspace.IsActive,
		&workspace.CreatedBy,
		&workspace.CreatedAt,
		&workspace.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &workspace, nil
}
