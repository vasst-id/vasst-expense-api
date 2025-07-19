package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/vasst-id/vasst-expense-api/internal/entities"
	"github.com/vasst-id/vasst-expense-api/internal/repositories"
	errorsutil "github.com/vasst-id/vasst-expense-api/internal/utils/errors"
)

//go:generate mockgen -source=workspace_service.go -package=mock -destination=mock/workspace_service_mock.go
type (
	WorkspaceService interface {
		// Global (superadmin) methods
		CreateWorkspace(ctx context.Context, userID uuid.UUID, input *entities.CreateWorkspaceInput) (*entities.Workspace, error)
		UpdateWorkspace(ctx context.Context, workspaceID uuid.UUID, input *entities.UpdateWorkspaceInput) (*entities.Workspace, error)
		DeleteWorkspace(ctx context.Context, workspaceID uuid.UUID) error
		ListAllWorkspaces(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entities.Workspace, error)
		GetWorkspaceByID(ctx context.Context, workspaceID uuid.UUID) (*entities.Workspace, error)
	}

	workspaceService struct {
		workspaceRepo repositories.WorkspaceRepository
	}
)

// NewWorkspaceService creates a new workspace service
func NewWorkspaceService(workspaceRepo repositories.WorkspaceRepository) WorkspaceService {
	return &workspaceService{
		workspaceRepo: workspaceRepo,
	}
}

// CreateWorkspace creates a new workspace
func (s *workspaceService) CreateWorkspace(ctx context.Context, userID uuid.UUID, input *entities.CreateWorkspaceInput) (*entities.Workspace, error) {
	// Validate required fields
	if input.Name == "" {
		return nil, errors.New("workspace name is required")
	}
	if input.WorkspaceType == 0 {
		return nil, errors.New("workspace type is required")
	}
	if input.CurrencyID == 0 {
		return nil, errors.New("currency ID is required")
	}

	// Set default timezone if not provided
	timezone := input.Timezone
	if timezone == "" {
		timezone = "Asia/Jakarta"
	}

	// Check if workspace with same name already exists
	existingWorkspace, err := s.workspaceRepo.FindByName(ctx, input.Name)
	if err != nil {
		return nil, err
	}
	if existingWorkspace != nil {
		return nil, errorsutil.New(409, "workspace with this name already exists")
	}

	workspace := &entities.Workspace{
		WorkspaceID:   uuid.New(),
		Name:          input.Name,
		Description:   input.Description,
		WorkspaceType: input.WorkspaceType,
		Icon:          input.Icon,
		ColorCode:     input.ColorCode,
		CurrencyID:    input.CurrencyID,
		Timezone:      timezone,
		Settings:      "{}",   // Default empty JSON
		IsActive:      true,   // Default to active
		CreatedBy:     userID, // Use the authenticated user ID
	}

	// Create the workspace - the repository will populate the struct with the actual data from DB
	createdWorkspace, err := s.workspaceRepo.Create(ctx, workspace)
	if err != nil {
		return nil, err
	}

	// Return the workspace with data populated from the database
	return &createdWorkspace, nil
}

// UpdateWorkspace updates an existing workspace
func (s *workspaceService) UpdateWorkspace(ctx context.Context, workspaceID uuid.UUID, input *entities.UpdateWorkspaceInput) (*entities.Workspace, error) {
	existingWorkspace, err := s.workspaceRepo.FindByID(ctx, workspaceID)
	if err != nil {
		return nil, err
	}
	if existingWorkspace == nil {
		return nil, errorsutil.New(404, "workspace not found")
	}

	// Check for name uniqueness if name is being changed
	if input.Name != "" && input.Name != existingWorkspace.Name {
		workspaceWithName, err := s.workspaceRepo.FindByName(ctx, input.Name)
		if err != nil {
			return nil, err
		}
		if workspaceWithName != nil && workspaceWithName.WorkspaceID != workspaceID {
			return nil, errorsutil.New(409, "workspace name already in use")
		}
	}

	// Update fields
	if input.Name != "" {
		existingWorkspace.Name = input.Name
	}
	if input.Description != "" {
		existingWorkspace.Description = input.Description
	}
	if input.WorkspaceType != 0 {
		existingWorkspace.WorkspaceType = input.WorkspaceType
	}
	if input.Icon != "" {
		existingWorkspace.Icon = input.Icon
	}
	if input.ColorCode != "" {
		existingWorkspace.ColorCode = input.ColorCode
	}
	if input.CurrencyID != 0 {
		existingWorkspace.CurrencyID = input.CurrencyID
	}
	if input.Timezone != "" {
		existingWorkspace.Timezone = input.Timezone
	}
	// Update IsActive field
	existingWorkspace.IsActive = input.IsActive

	// Update the workspace - the repository will populate the struct with the actual data from DB
	updatedWorkspace, err := s.workspaceRepo.Update(ctx, existingWorkspace)
	if err != nil {
		return nil, err
	}

	// Return the workspace with data populated from the database
	return &updatedWorkspace, nil
}

// DeleteWorkspace deletes a workspace
func (s *workspaceService) DeleteWorkspace(ctx context.Context, workspaceID uuid.UUID) error {
	existingWorkspace, err := s.workspaceRepo.FindByID(ctx, workspaceID)
	if err != nil {
		return err
	}
	if existingWorkspace == nil {
		return errorsutil.New(404, "workspace not found")
	}
	return s.workspaceRepo.Delete(ctx, workspaceID)
}

// ListAllWorkspaces returns all workspaces with pagination
func (s *workspaceService) ListAllWorkspaces(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entities.Workspace, error) {
	return s.workspaceRepo.ListAll(ctx, userID, limit, offset)
}

// GetWorkspaceByID returns a workspace by ID
func (s *workspaceService) GetWorkspaceByID(ctx context.Context, workspaceID uuid.UUID) (*entities.Workspace, error) {
	workspace, err := s.workspaceRepo.FindByID(ctx, workspaceID)
	if err != nil {
		return nil, err
	}
	if workspace == nil {
		return nil, errorsutil.New(404, "workspace not found")
	}
	return workspace, nil
}
