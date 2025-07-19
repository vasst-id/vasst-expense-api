package services

import (
	"context"
	"errors"

	"github.com/vasst-id/vasst-expense-api/internal/entities"
	"github.com/vasst-id/vasst-expense-api/internal/repositories"
	errorsutil "github.com/vasst-id/vasst-expense-api/internal/utils/errors"
)

//go:generate mockgen -source=taxonomy_service.go -package=mock -destination=mock/taxonomy_service_mock.go
type (
	TaxonomyService interface {
		CreateTaxonomy(ctx context.Context, input *entities.CreateTaxonomyRequest) (*entities.Taxonomy, error)
		UpdateTaxonomy(ctx context.Context, taxonomyID int, input *entities.UpdateTaxonomyRequest) (*entities.Taxonomy, error)
		DeleteTaxonomy(ctx context.Context, taxonomyID int) error
		GetTaxonomyByID(ctx context.Context, taxonomyID int) (*entities.Taxonomy, error)
		GetTaxonomiesByType(ctx context.Context, taxonomyType string, limit, offset int) ([]*entities.Taxonomy, error)
		GetTaxonomyByTypeAndValue(ctx context.Context, taxonomyType, value string) (*entities.Taxonomy, error)
		GetActiveTaxonomies(ctx context.Context, limit, offset int) ([]*entities.Taxonomy, error)
		GetAllTaxonomies(ctx context.Context, limit, offset int) ([]*entities.Taxonomy, error)
	}

	taxonomyService struct {
		taxonomyRepo repositories.TaxonomyRepository
	}
)

// NewTaxonomyService creates a new taxonomy service
func NewTaxonomyService(taxonomyRepo repositories.TaxonomyRepository) TaxonomyService {
	return &taxonomyService{
		taxonomyRepo: taxonomyRepo,
	}
}

// CreateTaxonomy creates a new taxonomy
func (s *taxonomyService) CreateTaxonomy(ctx context.Context, input *entities.CreateTaxonomyRequest) (*entities.Taxonomy, error) {
	// Validate required fields
	if input.Label == "" {
		return nil, errors.New("label is required")
	}
	if input.Value == "" {
		return nil, errors.New("value is required")
	}
	if input.Type == "" {
		return nil, errors.New("type is required")
	}
	if input.TypeLabel == "" {
		return nil, errors.New("type label is required")
	}

	// Check if taxonomy with the same type and value already exists
	existingTaxonomy, err := s.taxonomyRepo.FindByTypeAndValue(ctx, input.Type, input.Value)
	if err != nil {
		return nil, err
	}
	if existingTaxonomy != nil {
		return nil, errorsutil.New(409, "taxonomy with this type and value already exists")
	}

	// Create new taxonomy
	taxonomy := &entities.Taxonomy{
		Label:     input.Label,
		Value:     input.Value,
		Type:      input.Type,
		TypeLabel: input.TypeLabel,
		Status:    entities.TaxonomyStatusActive,
	}

	// Create the taxonomy - the repository will populate the struct with the actual data from DB
	createdTaxonomy, err := s.taxonomyRepo.Create(ctx, taxonomy)
	if err != nil {
		return nil, err
	}

	// Return the taxonomy with data populated from the database
	return &createdTaxonomy, nil
}

// UpdateTaxonomy updates an existing taxonomy
func (s *taxonomyService) UpdateTaxonomy(ctx context.Context, taxonomyID int, input *entities.UpdateTaxonomyRequest) (*entities.Taxonomy, error) {
	// Get existing taxonomy
	existingTaxonomy, err := s.taxonomyRepo.FindByID(ctx, taxonomyID)
	if err != nil {
		return nil, err
	}
	if existingTaxonomy == nil {
		return nil, errorsutil.New(404, "taxonomy not found")
	}

	// Validate required fields
	if input.Label == "" {
		return nil, errors.New("label is required")
	}
	if input.Value == "" {
		return nil, errors.New("value is required")
	}
	if input.Type == "" {
		return nil, errors.New("type is required")
	}
	if input.TypeLabel == "" {
		return nil, errors.New("type label is required")
	}

	// Check for uniqueness if type or value is being changed
	if input.Type != existingTaxonomy.Type || input.Value != existingTaxonomy.Value {
		taxonomyWithTypeAndValue, err := s.taxonomyRepo.FindByTypeAndValue(ctx, input.Type, input.Value)
		if err != nil {
			return nil, err
		}
		if taxonomyWithTypeAndValue != nil && taxonomyWithTypeAndValue.TaxonomyID != taxonomyID {
			return nil, errorsutil.New(409, "taxonomy with this type and value already exists")
		}
	}

	// Update fields
	existingTaxonomy.Label = input.Label
	existingTaxonomy.Value = input.Value
	existingTaxonomy.Type = input.Type
	existingTaxonomy.TypeLabel = input.TypeLabel
	existingTaxonomy.Status = input.Status

	// Update the taxonomy - the repository will populate the struct with the actual data from DB
	updatedTaxonomy, err := s.taxonomyRepo.Update(ctx, existingTaxonomy)
	if err != nil {
		return nil, err
	}

	// Return the taxonomy with data populated from the database
	return &updatedTaxonomy, nil
}

// DeleteTaxonomy deletes a taxonomy (soft delete)
func (s *taxonomyService) DeleteTaxonomy(ctx context.Context, taxonomyID int) error {
	// Get existing taxonomy
	existingTaxonomy, err := s.taxonomyRepo.FindByID(ctx, taxonomyID)
	if err != nil {
		return err
	}
	if existingTaxonomy == nil {
		return errorsutil.New(404, "taxonomy not found")
	}

	return s.taxonomyRepo.Delete(ctx, taxonomyID)
}

// GetTaxonomyByID returns a taxonomy by ID
func (s *taxonomyService) GetTaxonomyByID(ctx context.Context, taxonomyID int) (*entities.Taxonomy, error) {
	taxonomy, err := s.taxonomyRepo.FindByID(ctx, taxonomyID)
	if err != nil {
		return nil, err
	}
	if taxonomy == nil {
		return nil, errorsutil.New(404, "taxonomy not found")
	}
	return taxonomy, nil
}

// GetTaxonomiesByType returns taxonomies by type with pagination
func (s *taxonomyService) GetTaxonomiesByType(ctx context.Context, taxonomyType string, limit, offset int) ([]*entities.Taxonomy, error) {
	if taxonomyType == "" {
		return nil, errors.New("taxonomy type is required")
	}
	return s.taxonomyRepo.FindByType(ctx, taxonomyType, limit, offset)
}

// GetTaxonomyByTypeAndValue returns a taxonomy by type and value
func (s *taxonomyService) GetTaxonomyByTypeAndValue(ctx context.Context, taxonomyType, value string) (*entities.Taxonomy, error) {
	if taxonomyType == "" {
		return nil, errors.New("taxonomy type is required")
	}
	if value == "" {
		return nil, errors.New("value is required")
	}

	taxonomy, err := s.taxonomyRepo.FindByTypeAndValue(ctx, taxonomyType, value)
	if err != nil {
		return nil, err
	}
	if taxonomy == nil {
		return nil, errorsutil.New(404, "taxonomy not found")
	}
	return taxonomy, nil
}

// GetActiveTaxonomies returns all active taxonomies with pagination
func (s *taxonomyService) GetActiveTaxonomies(ctx context.Context, limit, offset int) ([]*entities.Taxonomy, error) {
	return s.taxonomyRepo.FindActive(ctx, limit, offset)
}

// GetAllTaxonomies returns all taxonomies with pagination
func (s *taxonomyService) GetAllTaxonomies(ctx context.Context, limit, offset int) ([]*entities.Taxonomy, error) {
	return s.taxonomyRepo.FindAll(ctx, limit, offset)
}
