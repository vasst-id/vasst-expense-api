package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/vasst-id/vasst-expense-api/internal/entities"
	"github.com/vasst-id/vasst-expense-api/internal/utils/postgres"
)

type (
	categoryRepository struct {
		*postgres.Postgres
	}

	// CategoryRepository defines methods for interacting with categories
	CategoryRepository interface {
		// System categories
		CreateSystemCategory(ctx context.Context, category *entities.Category) (entities.Category, error)
		GetSystemCategories(ctx context.Context, pagination *entities.PaginationRequest) ([]*entities.Category, int, error)
		FindSystemCategoryByID(ctx context.Context, categoryID uuid.UUID) (*entities.Category, error)

		// User categories
		CreateUserCategory(ctx context.Context, userCategory *entities.UserCategory) (entities.UserCategory, error)
		UpdateUserCategory(ctx context.Context, userCategory *entities.UserCategory) (entities.UserCategory, error)
		DeleteUserCategory(ctx context.Context, userCategoryID uuid.UUID) error
		FindUserCategoryByID(ctx context.Context, userCategoryID uuid.UUID) (*entities.UserCategory, error)
		FindUserCategories(ctx context.Context, userID uuid.UUID, pagination *entities.PaginationRequest, filter *entities.FilterRequest) ([]*entities.UserCategory, int, error)
		FindActiveUserCategories(ctx context.Context, userID uuid.UUID) ([]*entities.UserCategory, error)

		// Category operations
		ExistsByName(ctx context.Context, userID uuid.UUID, name string) (bool, error)
		GetCategoryUsage(ctx context.Context, userCategoryID uuid.UUID) (int, error)
		GetCategoriesWithTransactionCount(ctx context.Context, userID uuid.UUID) ([]map[string]interface{}, error)

		// User category management
		AddSystemCategoryToUser(ctx context.Context, userID, categoryID uuid.UUID) error
		RemoveUserCategory(ctx context.Context, userCategoryID uuid.UUID) error
		GetUserCategoryHierarchy(ctx context.Context, userID uuid.UUID) ([]*entities.UserCategory, error)
	}
)

// NewCategoryRepository creates a new CategoryRepository
func NewCategoryRepository(pg *postgres.Postgres) CategoryRepository {
	return &categoryRepository{pg}
}

// CreateSystemCategory creates a new system category
func (r *categoryRepository) CreateSystemCategory(ctx context.Context, category *entities.Category) (entities.Category, error) {
	query := `
		INSERT INTO "vasst_expense".categories 
		(category_id, name, description, icon, parent_category_id, is_system_category, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING category_id, name, description, icon, parent_category_id, is_system_category, created_at, updated_at
	`

	var createdCategory entities.Category
	err := r.DB.QueryRowContext(ctx, query,
		category.CategoryID, category.Name, category.Description, category.Icon,
		category.ParentCategoryID, category.IsSystemCategory,
	).Scan(
		&createdCategory.CategoryID, &createdCategory.Name, &createdCategory.Description, &createdCategory.Icon,
		&createdCategory.ParentCategoryID, &createdCategory.IsSystemCategory, &createdCategory.CreatedAt, &createdCategory.UpdatedAt,
	)

	if err != nil {
		return entities.Category{}, fmt.Errorf("failed to create system category: %w", err)
	}

	return createdCategory, nil
}

// GetSystemCategories gets system categories with pagination
func (r *categoryRepository) GetSystemCategories(ctx context.Context, pagination *entities.PaginationRequest) ([]*entities.Category, int, error) {
	// Set defaults
	if pagination == nil {
		pagination = &entities.PaginationRequest{Page: 1, PageSize: 20}
	}
	pagination.SetDefaults()

	// Count total records
	countQuery := `
		SELECT COUNT(*) 
		FROM "vasst_expense".categories 
		WHERE is_system_category = true
	`

	var totalCount int
	err := r.DB.QueryRowContext(ctx, countQuery).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count system categories: %w", err)
	}

	// Get paginated results
	query := `
		SELECT category_id, name, description, icon, parent_category_id, is_system_category, created_at, updated_at
		FROM "vasst_expense".categories
		WHERE is_system_category = true
		ORDER BY name ASC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.DB.QueryContext(ctx, query, pagination.PageSize, pagination.GetOffset())
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query system categories: %w", err)
	}
	defer rows.Close()

	var categories []*entities.Category
	for rows.Next() {
		category := &entities.Category{}
		err := rows.Scan(
			&category.CategoryID, &category.Name, &category.Description, &category.Icon,
			&category.ParentCategoryID, &category.IsSystemCategory, &category.CreatedAt, &category.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan system category: %w", err)
		}
		categories = append(categories, category)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("rows iteration error: %w", err)
	}

	return categories, totalCount, nil
}

// FindSystemCategoryByID finds a system category by ID
func (r *categoryRepository) FindSystemCategoryByID(ctx context.Context, categoryID uuid.UUID) (*entities.Category, error) {
	query := `
		SELECT category_id, name, description, icon, parent_category_id, is_system_category, created_at, updated_at
		FROM "vasst_expense".categories 
		WHERE category_id = $1 AND is_system_category = true
	`

	category := &entities.Category{}
	err := r.DB.QueryRowContext(ctx, query, categoryID).Scan(
		&category.CategoryID, &category.Name, &category.Description, &category.Icon,
		&category.ParentCategoryID, &category.IsSystemCategory, &category.CreatedAt, &category.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, entities.ErrNotFound
		}
		return nil, fmt.Errorf("failed to find system category: %w", err)
	}

	return category, nil
}

// CreateUserCategory creates a new user category
func (r *categoryRepository) CreateUserCategory(ctx context.Context, userCategory *entities.UserCategory) (entities.UserCategory, error) {
	query := `
		INSERT INTO "vasst_expense".user_categories 
		(user_category_id, user_id, category_id, name, description, icon, is_custom, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING user_category_id, user_id, category_id, name, description, icon, is_custom, is_active, created_at, updated_at
	`

	var createdUserCategory entities.UserCategory
	err := r.DB.QueryRowContext(ctx, query,
		userCategory.UserCategoryID, userCategory.UserID, userCategory.CategoryID,
		userCategory.Name, userCategory.Description, userCategory.Icon,
		userCategory.IsCustom, userCategory.IsActive,
	).Scan(
		&createdUserCategory.UserCategoryID, &createdUserCategory.UserID, &createdUserCategory.CategoryID,
		&createdUserCategory.Name, &createdUserCategory.Description, &createdUserCategory.Icon,
		&createdUserCategory.IsCustom, &createdUserCategory.IsActive, &createdUserCategory.CreatedAt, &createdUserCategory.UpdatedAt,
	)

	if err != nil {
		return entities.UserCategory{}, fmt.Errorf("failed to create user category: %w", err)
	}

	return createdUserCategory, nil
}

// UpdateUserCategory updates a user category
func (r *categoryRepository) UpdateUserCategory(ctx context.Context, userCategory *entities.UserCategory) (entities.UserCategory, error) {
	query := `
		UPDATE "vasst_expense".user_categories 
		SET name = $2, description = $3, icon = $4, is_active = $5, updated_at = CURRENT_TIMESTAMP
		WHERE user_category_id = $1
		RETURNING user_category_id, user_id, category_id, name, description, icon, is_custom, is_active, created_at, updated_at
	`

	var updatedUserCategory entities.UserCategory
	err := r.DB.QueryRowContext(ctx, query,
		userCategory.UserCategoryID, userCategory.Name, userCategory.Description,
		userCategory.Icon, userCategory.IsActive,
	).Scan(
		&updatedUserCategory.UserCategoryID, &updatedUserCategory.UserID, &updatedUserCategory.CategoryID,
		&updatedUserCategory.Name, &updatedUserCategory.Description, &updatedUserCategory.Icon,
		&updatedUserCategory.IsCustom, &updatedUserCategory.IsActive, &updatedUserCategory.CreatedAt, &updatedUserCategory.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return entities.UserCategory{}, entities.ErrNotFound
		}
		return entities.UserCategory{}, fmt.Errorf("failed to update user category: %w", err)
	}

	return updatedUserCategory, nil
}

// DeleteUserCategory soft deletes a user category (sets is_active to false)
func (r *categoryRepository) DeleteUserCategory(ctx context.Context, userCategoryID uuid.UUID) error {
	query := `
		UPDATE "vasst_expense".user_categories 
		SET is_active = false, updated_at = CURRENT_TIMESTAMP
		WHERE user_category_id = $1
	`

	result, err := r.DB.ExecContext(ctx, query, userCategoryID)
	if err != nil {
		return fmt.Errorf("failed to delete user category: %w", err)
	}

	if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
		return entities.ErrNotFound
	}

	return nil
}

// FindUserCategoryByID finds a user category by ID
func (r *categoryRepository) FindUserCategoryByID(ctx context.Context, userCategoryID uuid.UUID) (*entities.UserCategory, error) {
	query := `
		SELECT user_category_id, user_id, category_id, name, description, icon, is_custom, is_active, created_at, updated_at
		FROM "vasst_expense".user_categories 
		WHERE user_category_id = $1 AND is_active = true
	`

	userCategory := &entities.UserCategory{}
	err := r.DB.QueryRowContext(ctx, query, userCategoryID).Scan(
		&userCategory.UserCategoryID, &userCategory.UserID, &userCategory.CategoryID,
		&userCategory.Name, &userCategory.Description, &userCategory.Icon,
		&userCategory.IsCustom, &userCategory.IsActive, &userCategory.CreatedAt, &userCategory.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, entities.ErrNotFound
		}
		return nil, fmt.Errorf("failed to find user category: %w", err)
	}

	return userCategory, nil
}

// FindUserCategories finds user categories with pagination and filtering
func (r *categoryRepository) FindUserCategories(ctx context.Context, userID uuid.UUID, pagination *entities.PaginationRequest, filter *entities.FilterRequest) ([]*entities.UserCategory, int, error) {
	// Set defaults
	if pagination == nil {
		pagination = &entities.PaginationRequest{Page: 1, PageSize: 20}
	}
	pagination.SetDefaults()

	if filter == nil {
		filter = &entities.FilterRequest{}
	}
	filter.SetDefaults()

	// Build WHERE clause
	whereConditions := []string{"uc.user_id = $1", "uc.is_active = true"}
	args := []interface{}{userID}
	argIndex := 2

	if filter.Search != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("uc.name ILIKE $%d", argIndex))
		args = append(args, "%"+filter.Search+"%")
		argIndex++
	}

	whereClause := strings.Join(whereConditions, " AND ")

	// Build ORDER BY clause
	orderBy := "uc.created_at"
	if filter.SortBy != "" {
		switch filter.SortBy {
		case "name", "created_at", "updated_at":
			orderBy = "uc." + filter.SortBy
		}
	}
	orderClause := fmt.Sprintf("ORDER BY %s %s", orderBy, strings.ToUpper(filter.SortOrder))

	// Count total records
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*) 
		FROM "vasst_expense".user_categories uc 
		WHERE %s
	`, whereClause)

	var totalCount int
	err := r.DB.QueryRowContext(ctx, countQuery, args[:argIndex-1]...).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count user categories: %w", err)
	}

	// Get paginated results
	query := fmt.Sprintf(`
		SELECT uc.user_category_id, uc.user_id, uc.category_id, uc.name, uc.description, 
		       uc.icon, uc.is_custom, uc.is_active, uc.created_at, uc.updated_at
		FROM "vasst_expense".user_categories uc
		WHERE %s
		%s
		LIMIT $%d OFFSET $%d
	`, whereClause, orderClause, argIndex, argIndex+1)

	args = append(args, pagination.PageSize, pagination.GetOffset())

	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query user categories: %w", err)
	}
	defer rows.Close()

	var userCategories []*entities.UserCategory
	for rows.Next() {
		userCategory := &entities.UserCategory{}
		err := rows.Scan(
			&userCategory.UserCategoryID, &userCategory.UserID, &userCategory.CategoryID,
			&userCategory.Name, &userCategory.Description, &userCategory.Icon,
			&userCategory.IsCustom, &userCategory.IsActive, &userCategory.CreatedAt, &userCategory.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan user category: %w", err)
		}
		userCategories = append(userCategories, userCategory)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("rows iteration error: %w", err)
	}

	return userCategories, totalCount, nil
}

// FindActiveUserCategories finds all active user categories (no pagination)
func (r *categoryRepository) FindActiveUserCategories(ctx context.Context, userID uuid.UUID) ([]*entities.UserCategory, error) {
	query := `
		SELECT user_category_id, user_id, category_id, name, description, icon, is_custom, is_active, created_at, updated_at
		FROM "vasst_expense".user_categories 
		WHERE user_id = $1 AND is_active = true
		ORDER BY name ASC
	`

	rows, err := r.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query active user categories: %w", err)
	}
	defer rows.Close()

	var userCategories []*entities.UserCategory
	for rows.Next() {
		userCategory := &entities.UserCategory{}
		err := rows.Scan(
			&userCategory.UserCategoryID, &userCategory.UserID, &userCategory.CategoryID,
			&userCategory.Name, &userCategory.Description, &userCategory.Icon,
			&userCategory.IsCustom, &userCategory.IsActive, &userCategory.CreatedAt, &userCategory.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user category: %w", err)
		}
		userCategories = append(userCategories, userCategory)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return userCategories, nil
}

// ExistsByName checks if a user category with the given name exists for the user
func (r *categoryRepository) ExistsByName(ctx context.Context, userID uuid.UUID, name string) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM "vasst_expense".user_categories 
			WHERE user_id = $1 AND name = $2 AND is_active = true
		)
	`

	var exists bool
	err := r.DB.QueryRowContext(ctx, query, userID, name).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check category name existence: %w", err)
	}

	return exists, nil
}

// GetCategoryUsage gets the number of transactions using this category
func (r *categoryRepository) GetCategoryUsage(ctx context.Context, userCategoryID uuid.UUID) (int, error) {
	query := `
		SELECT COUNT(*) 
		FROM "vasst_expense".transactions 
		WHERE category_id = $1
	`

	var count int
	err := r.DB.QueryRowContext(ctx, query, userCategoryID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get category usage: %w", err)
	}

	return count, nil
}

// GetCategoriesWithTransactionCount gets user categories with their transaction counts
func (r *categoryRepository) GetCategoriesWithTransactionCount(ctx context.Context, userID uuid.UUID) ([]map[string]interface{}, error) {
	query := `
		SELECT 
			uc.user_category_id,
			uc.name,
			uc.icon,
			uc.is_custom,
			COUNT(t.transaction_id) as transaction_count,
			COALESCE(SUM(CASE WHEN t.transaction_type = $1 THEN t.amount ELSE 0 END), 0) as total_spent
		FROM "vasst_expense".user_categories uc
		LEFT JOIN "vasst_expense".transactions t ON uc.user_category_id = t.category_id
		WHERE uc.user_id = $2 AND uc.is_active = true
		GROUP BY uc.user_category_id, uc.name, uc.icon, uc.is_custom
		ORDER BY total_spent DESC
	`

	rows, err := r.DB.QueryContext(ctx, query, entities.TransactionTypeExpense, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get categories with transaction count: %w", err)
	}
	defer rows.Close()

	var result []map[string]interface{}
	for rows.Next() {
		var userCategoryID uuid.UUID
		var name, icon string
		var isCustom bool
		var transactionCount int
		var totalSpent float64

		err := rows.Scan(&userCategoryID, &name, &icon, &isCustom, &transactionCount, &totalSpent)
		if err != nil {
			return nil, fmt.Errorf("failed to scan category with transaction count: %w", err)
		}

		result = append(result, map[string]interface{}{
			"user_category_id":  userCategoryID,
			"name":              name,
			"icon":              icon,
			"is_custom":         isCustom,
			"transaction_count": transactionCount,
			"total_spent":       totalSpent,
		})
	}

	return result, nil
}

// AddSystemCategoryToUser adds a system category to a user's categories
func (r *categoryRepository) AddSystemCategoryToUser(ctx context.Context, userID, categoryID uuid.UUID) error {
	// First check if the system category exists
	var exists bool
	checkQuery := `
		SELECT EXISTS(
			SELECT 1 FROM "vasst_expense".categories 
			WHERE category_id = $1 AND is_system_category = true
		)
	`
	err := r.DB.QueryRowContext(ctx, checkQuery, categoryID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check system category existence: %w", err)
	}

	if !exists {
		return entities.ErrNotFound
	}

	// Check if user already has this category
	checkUserQuery := `
		SELECT EXISTS(
			SELECT 1 FROM "vasst_expense".user_categories 
			WHERE user_id = $1 AND category_id = $2
		)
	`
	err = r.DB.QueryRowContext(ctx, checkUserQuery, userID, categoryID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check existing user category: %w", err)
	}

	if exists {
		return fmt.Errorf("user already has this category")
	}

	// Get system category details
	var categoryName, categoryIcon string
	var categoryDescription sql.NullString
	getQuery := `
		SELECT name, description, icon 
		FROM "vasst_expense".categories 
		WHERE category_id = $1
	`
	err = r.DB.QueryRowContext(ctx, getQuery, categoryID).Scan(&categoryName, &categoryDescription, &categoryIcon)
	if err != nil {
		return fmt.Errorf("failed to get system category details: %w", err)
	}

	// Create user category
	userCategory := &entities.UserCategory{
		UserCategoryID: uuid.New(),
		UserID:         userID,
		CategoryID:     categoryID,
		Name:           categoryName,
		Icon:           &categoryIcon,
		IsCustom:       false,
		IsActive:       true,
	}

	if categoryDescription.Valid {
		userCategory.Description = &categoryDescription.String
	}

	_, err = r.CreateUserCategory(ctx, userCategory)
	return err
}

// RemoveUserCategory removes a user category (hard delete for custom categories, soft delete for system categories)
func (r *categoryRepository) RemoveUserCategory(ctx context.Context, userCategoryID uuid.UUID) error {
	// Check if it's a custom category
	var isCustom bool
	checkQuery := `
		SELECT is_custom 
		FROM "vasst_expense".user_categories 
		WHERE user_category_id = $1
	`
	err := r.DB.QueryRowContext(ctx, checkQuery, userCategoryID).Scan(&isCustom)
	if err != nil {
		if err == sql.ErrNoRows {
			return entities.ErrNotFound
		}
		return fmt.Errorf("failed to check category type: %w", err)
	}

	if isCustom {
		// Hard delete for custom categories (but only if no transactions use it)
		var transactionCount int
		countQuery := `
			SELECT COUNT(*) 
			FROM "vasst_expense".transactions 
			WHERE category_id = $1
		`
		err = r.DB.QueryRowContext(ctx, countQuery, userCategoryID).Scan(&transactionCount)
		if err != nil {
			return fmt.Errorf("failed to count transactions: %w", err)
		}

		if transactionCount > 0 {
			return fmt.Errorf("cannot delete category with existing transactions")
		}

		deleteQuery := `DELETE FROM "vasst_expense".user_categories WHERE user_category_id = $1`
		result, err := r.DB.ExecContext(ctx, deleteQuery, userCategoryID)
		if err != nil {
			return fmt.Errorf("failed to delete custom category: %w", err)
		}

		if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
			return entities.ErrNotFound
		}
	} else {
		// Soft delete for system categories
		return r.DeleteUserCategory(ctx, userCategoryID)
	}

	return nil
}

// GetUserCategoryHierarchy gets user categories organized by hierarchy
func (r *categoryRepository) GetUserCategoryHierarchy(ctx context.Context, userID uuid.UUID) ([]*entities.UserCategory, error) {
	query := `
		SELECT 
			uc.user_category_id, uc.user_id, uc.category_id, uc.name, uc.description, 
			uc.icon, uc.is_custom, uc.is_active, uc.created_at, uc.updated_at
		FROM "vasst_expense".user_categories uc
		LEFT JOIN "vasst_expense".categories c ON uc.category_id = c.category_id
		WHERE uc.user_id = $1 AND uc.is_active = true
		ORDER BY c.parent_category_id NULLS FIRST, uc.name ASC
	`

	rows, err := r.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query user category hierarchy: %w", err)
	}
	defer rows.Close()

	var userCategories []*entities.UserCategory
	for rows.Next() {
		userCategory := &entities.UserCategory{}
		err := rows.Scan(
			&userCategory.UserCategoryID, &userCategory.UserID, &userCategory.CategoryID,
			&userCategory.Name, &userCategory.Description, &userCategory.Icon,
			&userCategory.IsCustom, &userCategory.IsActive, &userCategory.CreatedAt, &userCategory.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user category: %w", err)
		}
		userCategories = append(userCategories, userCategory)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return userCategories, nil
}
