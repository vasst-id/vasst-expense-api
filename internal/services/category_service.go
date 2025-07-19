package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/vasst-id/vasst-expense-api/internal/entities"
	"github.com/vasst-id/vasst-expense-api/internal/repositories"
	errorsutil "github.com/vasst-id/vasst-expense-api/internal/utils/errors"
)

//go:generate mockgen -source=category_service.go -package=mock -destination=mock/category_service_mock.go
type (
	CategoryService interface {
		// System categories
		CreateSystemCategory(ctx context.Context, input *entities.CreateCategoryInput) (*entities.Category, error)
		GetSystemCategories(ctx context.Context, pagination *entities.PaginationRequest) ([]*entities.Category, int, error)
		GetSystemCategoryByID(ctx context.Context, categoryID uuid.UUID) (*entities.Category, error)

		// User categories
		CreateUserCategory(ctx context.Context, userID uuid.UUID, input *entities.CreateUserCategoryInput) (*entities.UserCategory, error)
		UpdateUserCategory(ctx context.Context, userID uuid.UUID, userCategoryID uuid.UUID, input *entities.UpdateUserCategoryInput) (*entities.UserCategory, error)
		DeleteUserCategory(ctx context.Context, userID uuid.UUID, userCategoryID uuid.UUID) error
		GetUserCategoryByID(ctx context.Context, userID uuid.UUID, userCategoryID uuid.UUID) (*entities.UserCategory, error)
		GetUserCategories(ctx context.Context, userID uuid.UUID, pagination *entities.PaginationRequest, filter *entities.FilterRequest) ([]*entities.UserCategory, int, error)
		GetActiveUserCategories(ctx context.Context, userID uuid.UUID) ([]*entities.UserCategory, error)

		// Category operations
		AddSystemCategoryToUser(ctx context.Context, userID, categoryID uuid.UUID) error
		GetCategoriesWithTransactionCount(ctx context.Context, userID uuid.UUID) ([]map[string]interface{}, error)
	}

	categoryService struct {
		categoryRepo repositories.CategoryRepository
	}
)

// NewCategoryService creates a new category service
func NewCategoryService(categoryRepo repositories.CategoryRepository) CategoryService {
	return &categoryService{
		categoryRepo: categoryRepo,
	}
}

// CreateSystemCategory creates a new system category
func (s *categoryService) CreateSystemCategory(ctx context.Context, input *entities.CreateCategoryInput) (*entities.Category, error) {
	// Validate required fields
	if input.Name == "" {
		return nil, errors.New("category name is required")
	}

	// Create new category
	category := &entities.Category{
		CategoryID:       uuid.New(),
		Name:             input.Name,
		Description:      input.Description,
		Icon:             input.Icon,
		ParentCategoryID: input.ParentCategoryID,
		IsSystemCategory: input.IsSystemCategory,
	}

	// Create the category - the repository will populate the struct with the actual data from DB
	createdCategory, err := s.categoryRepo.CreateSystemCategory(ctx, category)
	if err != nil {
		return nil, err
	}

	// Return the category with data populated from the database
	return &createdCategory, nil
}

// GetSystemCategories returns system categories with pagination
func (s *categoryService) GetSystemCategories(ctx context.Context, pagination *entities.PaginationRequest) ([]*entities.Category, int, error) {
	return s.categoryRepo.GetSystemCategories(ctx, pagination)
}

// GetSystemCategoryByID returns a system category by ID
func (s *categoryService) GetSystemCategoryByID(ctx context.Context, categoryID uuid.UUID) (*entities.Category, error) {
	category, err := s.categoryRepo.FindSystemCategoryByID(ctx, categoryID)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, errorsutil.New(404, "category not found")
	}
	return category, nil
}

// CreateUserCategory creates a new user category
func (s *categoryService) CreateUserCategory(ctx context.Context, userID uuid.UUID, input *entities.CreateUserCategoryInput) (*entities.UserCategory, error) {
	// Validate required fields
	if input.Name == "" {
		return nil, errors.New("category name is required")
	}

	// Check if category with the same name already exists for this user
	exists, err := s.categoryRepo.ExistsByName(ctx, userID, input.Name)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errorsutil.New(409, "category with this name already exists")
	}

	// Create new user category
	userCategory := &entities.UserCategory{
		UserCategoryID: uuid.New(),
		UserID:         userID,
		CategoryID:     input.CategoryID,
		Name:           input.Name,
		Description:    input.Description,
		Icon:           input.Icon,
		IsCustom:       input.IsCustom,
		IsActive:       true,
	}

	// Create the user category - the repository will populate the struct with the actual data from DB
	createdUserCategory, err := s.categoryRepo.CreateUserCategory(ctx, userCategory)
	if err != nil {
		return nil, err
	}

	// Return the user category with data populated from the database
	return &createdUserCategory, nil
}

// UpdateUserCategory updates an existing user category
func (s *categoryService) UpdateUserCategory(ctx context.Context, userID uuid.UUID, userCategoryID uuid.UUID, input *entities.UpdateUserCategoryInput) (*entities.UserCategory, error) {
	// Get existing user category and verify ownership
	existingUserCategory, err := s.categoryRepo.FindUserCategoryByID(ctx, userCategoryID)
	if err != nil {
		return nil, err
	}
	if existingUserCategory == nil {
		return nil, errorsutil.New(404, "user category not found")
	}
	if existingUserCategory.UserID != userID {
		return nil, errorsutil.New(403, "access denied")
	}

	// Validate required fields
	if input.Name == "" {
		return nil, errors.New("category name is required")
	}

	// Check for category name uniqueness if it's being changed
	if input.Name != existingUserCategory.Name {
		exists, err := s.categoryRepo.ExistsByName(ctx, userID, input.Name)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errorsutil.New(409, "category name already in use")
		}
	}

	// Update fields
	existingUserCategory.Name = input.Name
	existingUserCategory.Description = input.Description
	existingUserCategory.Icon = input.Icon
	existingUserCategory.IsActive = input.IsActive

	// Update the user category - the repository will populate the struct with the actual data from DB
	updatedUserCategory, err := s.categoryRepo.UpdateUserCategory(ctx, existingUserCategory)
	if err != nil {
		return nil, err
	}

	// Return the user category with data populated from the database
	return &updatedUserCategory, nil
}

// DeleteUserCategory deletes a user category (soft delete)
func (s *categoryService) DeleteUserCategory(ctx context.Context, userID uuid.UUID, userCategoryID uuid.UUID) error {
	// Get existing user category and verify ownership
	existingUserCategory, err := s.categoryRepo.FindUserCategoryByID(ctx, userCategoryID)
	if err != nil {
		return err
	}
	if existingUserCategory == nil {
		return errorsutil.New(404, "user category not found")
	}
	if existingUserCategory.UserID != userID {
		return errorsutil.New(403, "access denied")
	}

	// Check if category is being used
	usage, err := s.categoryRepo.GetCategoryUsage(ctx, userCategoryID)
	if err != nil {
		return err
	}
	if usage > 0 {
		return errorsutil.New(409, "cannot delete category with existing transactions")
	}

	return s.categoryRepo.DeleteUserCategory(ctx, userCategoryID)
}

// GetUserCategoryByID returns a user category by ID (with user ownership verification)
func (s *categoryService) GetUserCategoryByID(ctx context.Context, userID uuid.UUID, userCategoryID uuid.UUID) (*entities.UserCategory, error) {
	userCategory, err := s.categoryRepo.FindUserCategoryByID(ctx, userCategoryID)
	if err != nil {
		return nil, err
	}
	if userCategory == nil {
		return nil, errorsutil.New(404, "user category not found")
	}
	if userCategory.UserID != userID {
		return nil, errorsutil.New(403, "access denied")
	}
	return userCategory, nil
}

// GetUserCategories returns user categories with pagination and filtering
func (s *categoryService) GetUserCategories(ctx context.Context, userID uuid.UUID, pagination *entities.PaginationRequest, filter *entities.FilterRequest) ([]*entities.UserCategory, int, error) {
	return s.categoryRepo.FindUserCategories(ctx, userID, pagination, filter)
}

// GetActiveUserCategories returns all active user categories
func (s *categoryService) GetActiveUserCategories(ctx context.Context, userID uuid.UUID) ([]*entities.UserCategory, error) {
	return s.categoryRepo.FindActiveUserCategories(ctx, userID)
}

// AddSystemCategoryToUser adds a system category to a user's categories
func (s *categoryService) AddSystemCategoryToUser(ctx context.Context, userID, categoryID uuid.UUID) error {
	return s.categoryRepo.AddSystemCategoryToUser(ctx, userID, categoryID)
}

// GetCategoriesWithTransactionCount gets user categories with their transaction counts
func (s *categoryService) GetCategoriesWithTransactionCount(ctx context.Context, userID uuid.UUID) ([]map[string]interface{}, error) {
	return s.categoryRepo.GetCategoriesWithTransactionCount(ctx, userID)
}
