package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/vasst-id/vasst-expense-api/internal/entities"
	"github.com/vasst-id/vasst-expense-api/internal/utils/postgres"
)

type (
	organizationRepository struct {
		*postgres.Postgres
	}

	// OrganizationRepository defines methods for interacting with organizations and related entities
	OrganizationRepository interface {
		// Organization methods
		CreateOrganization(ctx context.Context, org *entities.Organization) error
		UpdateOrganization(ctx context.Context, org *entities.Organization) error
		DeleteOrganization(ctx context.Context, orgID uuid.UUID) error
		ListOrganizations(ctx context.Context, limit, offset int) ([]*entities.Organization, error)
		FindOrganizationByID(ctx context.Context, orgID uuid.UUID) (*entities.Organization, error)
		FindOrganizationByCode(ctx context.Context, code string) (*entities.Organization, error)
		FindOrganizationByKey(ctx context.Context, key string) (*entities.Organization, error)
		FindDefaultUserByOrgID(ctx context.Context, orgID uuid.UUID) (*entities.User, error)

		// OrganizationCategory methods
		CreateCategory(ctx context.Context, category *entities.OrganizationCategory) error
		UpdateCategory(ctx context.Context, category *entities.OrganizationCategory) error
		DeleteCategory(ctx context.Context, categoryID int) error
		ListCategories(ctx context.Context) ([]*entities.OrganizationCategory, error)
		FindCategoryByID(ctx context.Context, categoryID int) (*entities.OrganizationCategory, error)

		// OrganizationSetting methods
		UpdateSetting(ctx context.Context, setting *entities.OrganizationSetting) error
		FindSettingByOrgID(ctx context.Context, orgID uuid.UUID) (*entities.OrganizationSetting, error)

		// OrganizationKnowledge methods
		CreateKnowledge(ctx context.Context, knowledge *entities.OrganizationKnowledge) error
		UpdateKnowledge(ctx context.Context, knowledge *entities.OrganizationKnowledge) error
		DeleteKnowledge(ctx context.Context, knowledgeID uuid.UUID) error
		ListKnowledgeByOrgID(ctx context.Context, orgID uuid.UUID) ([]*entities.OrganizationKnowledge, error)
		FindKnowledgeByID(ctx context.Context, knowledgeID uuid.UUID) (*entities.OrganizationKnowledge, error)

		// OrganizationModel methods
		CreateModel(ctx context.Context, model *entities.OrganizationModel) error
		DeleteModel(ctx context.Context, modelID uuid.UUID) error
		ListModelsByOrgID(ctx context.Context, orgID uuid.UUID) ([]*entities.OrganizationModel, error)

		// OrganizationIntegration methods
		CreateIntegration(ctx context.Context, integration *entities.OrganizationIntegration) error
		UpdateIntegration(ctx context.Context, integration *entities.OrganizationIntegration) error
		DeleteIntegration(ctx context.Context, integrationID uuid.UUID) error
		ListIntegrationsByOrgID(ctx context.Context, orgID uuid.UUID) ([]*entities.OrganizationIntegration, error)
		FindIntegrationByID(ctx context.Context, integrationID uuid.UUID) (*entities.OrganizationIntegration, error)
		FindIntegrationTokenByOrgIDAndType(ctx context.Context, orgID uuid.UUID, integrationType string) (string, error)
		FindIntegrationByOrgIDAndType(ctx context.Context, orgID uuid.UUID, integrationType string) (*entities.OrganizationIntegration, error)
	}
)

// NewOrganizationRepository creates a new OrganizationRepository
func NewOrganizationRepository(pg *postgres.Postgres) OrganizationRepository {
	return &organizationRepository{pg}
}

// Organization methods
func (r *organizationRepository) CreateOrganization(ctx context.Context, org *entities.Organization) error {
	query := `
		INSERT INTO "vasst_ca".organization (
			organization_id, organization_code, name, contact_name, phone_number,
			email, whatsapp_number, organization_type, organization_category_id,
			address, city, province, postal_code, country, status, api_key,
			created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17,
			CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`

	_, err := r.DB.ExecContext(ctx, query,
		org.OrganizationID,
		org.OrganizationCode,
		org.Name,
		org.ContactName,
		org.PhoneNumber,
		org.Email,
		org.WhatsappNumber,
		org.OrganizationType,
		org.CategoryID,
		org.Address,
		org.City,
		org.Province,
		org.PostalCode,
		org.Country,
		org.Status,
		org.APIKey,
	)

	// Create organization setting
	planID := entities.PlanStarter
	planStartDate := time.Now()
	planEndDate := time.Now().AddDate(1, 0, 0)
	planStatus := 1
	planAmount := 0
	planCurrency := "IDR"
	maxContacts := 100
	maxMessages := 1000
	maxBroadcasts := 10
	maxUsers := 100
	maxTags := 10
	maxOrders := 100
	currentContacts := 0
	currentMessages := 0
	currentBroadcasts := 0
	currentUsers := 0
	currentTags := 0
	currentOrders := 0

	query = `
		INSERT INTO "vasst_ca".organization_setting (
			organization_id, plan_id, plan_start_date, plan_end_date,
			plan_status, plan_amount, plan_currency,
			max_contacts, max_messages, max_broadcasts, max_users,
			max_tags, max_orders,
			current_contacts, current_messages, current_broadcasts,
			current_users, current_tags, current_orders,
			created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15,
			$16, $17, $18, $19, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`

	_, err = r.DB.ExecContext(ctx, query,
		org.OrganizationID,
		planID,
		planStartDate,
		planEndDate,
		planStatus,
		planAmount,
		planCurrency,
		maxContacts,
		maxMessages,
		maxBroadcasts,
		maxUsers,
		maxTags,
		maxOrders,
		currentContacts,
		currentMessages,
		currentBroadcasts,
		currentUsers,
		currentTags,
		currentOrders,
	)

	return err
}

func (r *organizationRepository) UpdateOrganization(ctx context.Context, org *entities.Organization) error {
	query := `
		UPDATE "vasst_ca".organization
		SET name = $1,
			contact_name = $2,
			phone_number = $3,
			email = $4,
			whatsapp_number = $5,
			organization_type = $6,
			organization_category_id = $7,
			address = $8,
			city = $9,
			province = $10,
			postal_code = $11,
			country = $12,
			status = $13,
			updated_at = CURRENT_TIMESTAMP
		WHERE organization_id = $14
	`

	result, err := r.DB.ExecContext(ctx, query,
		org.Name,
		org.ContactName,
		org.PhoneNumber,
		org.Email,
		org.WhatsappNumber,
		org.OrganizationType,
		org.CategoryID,
		org.Address,
		org.City,
		org.Province,
		org.PostalCode,
		org.Country,
		org.Status,
		org.OrganizationID,
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

func (r *organizationRepository) DeleteOrganization(ctx context.Context, orgID uuid.UUID) error {
	query := `
		DELETE FROM "vasst_ca".organization
		WHERE organization_id = $1
	`

	result, err := r.DB.ExecContext(ctx, query, orgID)
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

func (r *organizationRepository) ListOrganizations(ctx context.Context, limit, offset int) ([]*entities.Organization, error) {
	query := `
		SELECT organization_id, organization_code, name, contact_name, phone_number,
			email, whatsapp_number, organization_type, organization_category_id,
			address, city, province, postal_code, country, status,
			created_at, updated_at
		FROM "vasst_ca".organization
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.DB.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var organizations []*entities.Organization
	for rows.Next() {
		var org entities.Organization
		err := rows.Scan(
			&org.OrganizationID,
			&org.OrganizationCode,
			&org.Name,
			&org.ContactName,
			&org.PhoneNumber,
			&org.Email,
			&org.WhatsappNumber,
			&org.OrganizationType,
			&org.CategoryID,
			&org.Address,
			&org.City,
			&org.Province,
			&org.PostalCode,
			&org.Country,
			&org.Status,
			&org.CreatedAt,
			&org.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		organizations = append(organizations, &org)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return organizations, nil
}

func (r *organizationRepository) FindOrganizationByID(ctx context.Context, orgID uuid.UUID) (*entities.Organization, error) {
	query := `
		SELECT organization_id, organization_code, name, contact_name, phone_number,
			email, whatsapp_number, organization_type, organization_category_id,
			address, city, province, postal_code, country, status,
			created_at, updated_at
		FROM "vasst_ca".organization
		WHERE organization_id = $1
	`

	var org entities.Organization
	err := r.DB.QueryRowContext(ctx, query, orgID).Scan(
		&org.OrganizationID,
		&org.OrganizationCode,
		&org.Name,
		&org.ContactName,
		&org.PhoneNumber,
		&org.Email,
		&org.WhatsappNumber,
		&org.OrganizationType,
		&org.CategoryID,
		&org.Address,
		&org.City,
		&org.Province,
		&org.PostalCode,
		&org.Country,
		&org.Status,
		&org.CreatedAt,
		&org.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &org, nil
}

func (r *organizationRepository) FindOrganizationByCode(ctx context.Context, code string) (*entities.Organization, error) {
	query := `
		SELECT organization_id, organization_code, name, contact_name, phone_number,
			email, whatsapp_number, organization_type, organization_category_id,
			address, city, province, postal_code, country, status,
			created_at, updated_at
		FROM "vasst_ca".organization
		WHERE organization_code = $1
	`

	var org entities.Organization
	err := r.DB.QueryRowContext(ctx, query, code).Scan(
		&org.OrganizationID,
		&org.OrganizationCode,
		&org.Name,
		&org.ContactName,
		&org.PhoneNumber,
		&org.Email,
		&org.WhatsappNumber,
		&org.OrganizationType,
		&org.CategoryID,
		&org.Address,
		&org.City,
		&org.Province,
		&org.PostalCode,
		&org.Country,
		&org.Status,
		&org.CreatedAt,
		&org.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &org, nil
}

// OrganizationCategory methods
func (r *organizationRepository) CreateCategory(ctx context.Context, category *entities.OrganizationCategory) error {
	query := `
		INSERT INTO "vasst_ca".organization_category (
			organization_category_id, name, description, image_url, is_active,
			created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`

	_, err := r.DB.ExecContext(ctx, query,
		category.CategoryID,
		category.Name,
		category.Description,
		category.ImageURL,
		category.IsActive,
	)

	return err
}

func (r *organizationRepository) UpdateCategory(ctx context.Context, category *entities.OrganizationCategory) error {
	query := `
		UPDATE "vasst_ca".organization_category
		SET name = $1,
			description = $2,
			image_url = $3,
			is_active = $4,
			updated_at = CURRENT_TIMESTAMP
		WHERE organization_category_id = $5
	`

	result, err := r.DB.ExecContext(ctx, query,
		category.Name,
		category.Description,
		category.ImageURL,
		category.IsActive,
		category.CategoryID,
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

func (r *organizationRepository) DeleteCategory(ctx context.Context, categoryID int) error {
	query := `
		DELETE FROM "vasst_ca".organization_category
		WHERE organization_category_id = $1
	`

	result, err := r.DB.ExecContext(ctx, query, categoryID)
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

func (r *organizationRepository) ListCategories(ctx context.Context) ([]*entities.OrganizationCategory, error) {
	query := `
		SELECT organization_category_id, name, description, image_url, is_active,
			created_at, updated_at
		FROM "vasst_ca".organization_category
		ORDER BY name ASC
	`

	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*entities.OrganizationCategory
	for rows.Next() {
		var category entities.OrganizationCategory
		err := rows.Scan(
			&category.CategoryID,
			&category.Name,
			&category.Description,
			&category.ImageURL,
			&category.IsActive,
			&category.CreatedAt,
			&category.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		categories = append(categories, &category)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *organizationRepository) FindCategoryByID(ctx context.Context, categoryID int) (*entities.OrganizationCategory, error) {
	query := `
		SELECT organization_category_id, name, description, image_url, is_active,
			created_at, updated_at
		FROM "vasst_ca".organization_category
		WHERE organization_category_id = $1
	`

	var category entities.OrganizationCategory
	err := r.DB.QueryRowContext(ctx, query, categoryID).Scan(
		&category.CategoryID,
		&category.Name,
		&category.Description,
		&category.ImageURL,
		&category.IsActive,
		&category.CreatedAt,
		&category.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &category, nil
}

// OrganizationSetting methods

func (r *organizationRepository) UpdateSetting(ctx context.Context, setting *entities.OrganizationSetting) error {
	query := `
		UPDATE "vasst_ca".organization_setting
		SET plan_id = $1,
			plan_start_date = $2,
			plan_end_date = $3,
			plan_status = $4,
			plan_amount = $5,
			plan_currency = $6,
			max_contacts = $7,
			max_messages = $8,
			max_broadcasts = $9,
			max_users = $10,
			max_tags = $11,
			max_orders = $12,
			current_contacts = $13,
			current_messages = $14,
			current_broadcasts = $15,
			current_users = $16,
			current_tags = $17,
			current_orders = $18,
			system_prompt = $19,
			ai_assistant_name = $20,
			ai_communication_style = $21,
			ai_communication_language = $22,
			updated_at = CURRENT_TIMESTAMP
		WHERE organization_id = $23
	`

	result, err := r.DB.ExecContext(ctx, query,
		setting.PlanID,
		setting.PlanStartDate,
		setting.PlanEndDate,
		setting.PlanStatus,
		setting.PlanAmount,
		setting.PlanCurrency,
		setting.MaxContacts,
		setting.MaxMessages,
		setting.MaxBroadcasts,
		setting.MaxUsers,
		setting.MaxTags,
		setting.MaxOrders,
		setting.CurrentContacts,
		setting.CurrentMessages,
		setting.CurrentBroadcasts,
		setting.CurrentUsers,
		setting.CurrentTags,
		setting.CurrentOrders,
		setting.SystemPrompt,
		setting.AIAssistantName,
		setting.AICommunicationStyle,
		setting.AICommunicationLanguage,
		setting.OrganizationID,
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

func (r *organizationRepository) FindSettingByOrgID(ctx context.Context, orgID uuid.UUID) (*entities.OrganizationSetting, error) {
	query := `
		SELECT organization_id, os.plan_id, p.name,
			os.plan_start_date, os.plan_end_date,
			os.plan_status, os.plan_amount, os.plan_currency,
			os.max_contacts, os.max_messages, os.max_broadcasts, os.max_users,
			os.max_tags, os.max_orders,
			os.current_contacts, os.current_messages, os.current_broadcasts,
			os.current_users, os.current_tags, os.current_orders,
			os.system_prompt, os.ai_assistant_name, os.ai_communication_style, os.ai_communication_language,
			os.created_at, os.updated_at
		FROM "vasst_ca".organization_setting os
		JOIN "vasst_ca".plan p ON os.plan_id = p.plan_id
		WHERE organization_id = $1
	`

	var setting entities.OrganizationSetting
	err := r.DB.QueryRowContext(ctx, query, orgID).Scan(
		&setting.OrganizationID,
		&setting.PlanID,
		&setting.PlanName,
		&setting.PlanStartDate,
		&setting.PlanEndDate,
		&setting.PlanStatus,
		&setting.PlanAmount,
		&setting.PlanCurrency,
		&setting.MaxContacts,
		&setting.MaxMessages,
		&setting.MaxBroadcasts,
		&setting.MaxUsers,
		&setting.MaxTags,
		&setting.MaxOrders,
		&setting.CurrentContacts,
		&setting.CurrentMessages,
		&setting.CurrentBroadcasts,
		&setting.CurrentUsers,
		&setting.CurrentTags,
		&setting.CurrentOrders,
		&setting.SystemPrompt,
		&setting.AIAssistantName,
		&setting.AICommunicationStyle,
		&setting.AICommunicationLanguage,
		&setting.CreatedAt,
		&setting.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &setting, nil
}

// OrganizationKnowledge methods
func (r *organizationRepository) CreateKnowledge(ctx context.Context, knowledge *entities.OrganizationKnowledge) error {
	query := `
		INSERT INTO "vasst_ca".organization_knowledge (
			organization_id, knowledge_type, title, knowledge_content, description, 
			source_url, metadata, file_name, file_size, content_type, 
			bucket_name, object_name, is_active
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING organization_knowledge_id, created_at, updated_at
	`

	// Convert metadata to JSON
	var metadataJSON json.RawMessage
	if knowledge.Metadata != nil {
		metadataJSON = knowledge.Metadata
	} else {
		metadataJSON = json.RawMessage("{}")
	}

	err := r.DB.QueryRowContext(ctx, query,
		knowledge.OrganizationID,
		knowledge.KnowledgeType,
		knowledge.Title,
		knowledge.Content,
		knowledge.Description,
		knowledge.SourceURL,
		metadataJSON,
		knowledge.FileName,
		knowledge.FileSize,
		knowledge.ContentType,
		knowledge.BucketName,
		knowledge.ObjectName,
		knowledge.IsActive,
	).Scan(&knowledge.KnowledgeID, &knowledge.CreatedAt, &knowledge.UpdatedAt)

	return err
}

func (r *organizationRepository) UpdateKnowledge(ctx context.Context, knowledge *entities.OrganizationKnowledge) error {
	query := `
		UPDATE "vasst_ca".organization_knowledge
		SET knowledge_type = $1,
			title = $2,
			knowledge_content = $3,
			description = $4,
			source_url = $5,
			metadata = $6,
			file_name = $7,
			file_size = $8,
			content_type = $9,
			bucket_name = $10,
			object_name = $11,
			is_active = $12,
			updated_at = CURRENT_TIMESTAMP
		WHERE organization_knowledge_id = $13
	`

	// Convert metadata to JSON
	var metadataJSON json.RawMessage
	if knowledge.Metadata != nil {
		metadataJSON = knowledge.Metadata
	} else {
		metadataJSON = json.RawMessage("{}")
	}

	result, err := r.DB.ExecContext(ctx, query,
		knowledge.KnowledgeType,
		knowledge.Title,
		knowledge.Content,
		knowledge.Description,
		knowledge.SourceURL,
		metadataJSON,
		knowledge.FileName,
		knowledge.FileSize,
		knowledge.ContentType,
		knowledge.BucketName,
		knowledge.ObjectName,
		knowledge.IsActive,
		knowledge.KnowledgeID,
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

func (r *organizationRepository) DeleteKnowledge(ctx context.Context, knowledgeID uuid.UUID) error {
	query := `
		DELETE FROM "vasst_ca".organization_knowledge
		WHERE organization_knowledge_id = $1
	`

	result, err := r.DB.ExecContext(ctx, query, knowledgeID)
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

func (r *organizationRepository) ListKnowledgeByOrgID(ctx context.Context, orgID uuid.UUID) ([]*entities.OrganizationKnowledge, error) {
	query := `
		SELECT organization_knowledge_id, organization_id, knowledge_type, title,
			knowledge_content, description, source_url, metadata, file_name, 
			file_size, content_type, bucket_name, object_name, is_active, 
			created_at, updated_at
		FROM "vasst_ca".organization_knowledge
		WHERE organization_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.DB.QueryContext(ctx, query, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var knowledgeList []*entities.OrganizationKnowledge
	for rows.Next() {
		var knowledge entities.OrganizationKnowledge
		err := rows.Scan(
			&knowledge.KnowledgeID,
			&knowledge.OrganizationID,
			&knowledge.KnowledgeType,
			&knowledge.Title,
			&knowledge.Content,
			&knowledge.Description,
			&knowledge.SourceURL,
			&knowledge.Metadata,
			&knowledge.FileName,
			&knowledge.FileSize,
			&knowledge.ContentType,
			&knowledge.BucketName,
			&knowledge.ObjectName,
			&knowledge.IsActive,
			&knowledge.CreatedAt,
			&knowledge.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		knowledgeList = append(knowledgeList, &knowledge)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return knowledgeList, nil
}

func (r *organizationRepository) FindKnowledgeByID(ctx context.Context, knowledgeID uuid.UUID) (*entities.OrganizationKnowledge, error) {
	query := `
		SELECT organization_knowledge_id, organization_id, knowledge_type, title,
			knowledge_content, description, source_url, metadata, file_name, 
			file_size, content_type, bucket_name, object_name, is_active, 
			created_at, updated_at
		FROM "vasst_ca".organization_knowledge
		WHERE organization_knowledge_id = $1
	`

	var knowledge entities.OrganizationKnowledge
	err := r.DB.QueryRowContext(ctx, query, knowledgeID).Scan(
		&knowledge.KnowledgeID,
		&knowledge.OrganizationID,
		&knowledge.KnowledgeType,
		&knowledge.Title,
		&knowledge.Content,
		&knowledge.Description,
		&knowledge.SourceURL,
		&knowledge.Metadata,
		&knowledge.FileName,
		&knowledge.FileSize,
		&knowledge.ContentType,
		&knowledge.BucketName,
		&knowledge.ObjectName,
		&knowledge.IsActive,
		&knowledge.CreatedAt,
		&knowledge.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &knowledge, nil
}

// OrganizationModel methods
func (r *organizationRepository) CreateModel(ctx context.Context, model *entities.OrganizationModel) error {
	query := `
		INSERT INTO "vasst_ca".organization_model (
			organization_model_id, organization_id, model_id,
			is_active, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`

	_, err := r.DB.ExecContext(ctx, query,
		model.OrganizationModelID,
		model.OrganizationID,
		model.ModelID,
		model.IsActive,
	)

	return err
}

func (r *organizationRepository) DeleteModel(ctx context.Context, modelID uuid.UUID) error {
	query := `
		DELETE FROM "vasst_ca".organization_model
		WHERE organization_model_id = $1
	`

	result, err := r.DB.ExecContext(ctx, query, modelID)
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

func (r *organizationRepository) ListModelsByOrgID(ctx context.Context, orgID uuid.UUID) ([]*entities.OrganizationModel, error) {
	query := `
		SELECT organization_model_id, organization_id, model_id,
			is_active, created_at, updated_at
		FROM "vasst_ca".organization_model
		WHERE organization_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.DB.QueryContext(ctx, query, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var models []*entities.OrganizationModel
	for rows.Next() {
		var model entities.OrganizationModel
		err := rows.Scan(
			&model.OrganizationModelID,
			&model.OrganizationID,
			&model.ModelID,
			&model.IsActive,
			&model.CreatedAt,
			&model.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		models = append(models, &model)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return models, nil
}

// OrganizationIntegration methods
func (r *organizationRepository) CreateIntegration(ctx context.Context, integration *entities.OrganizationIntegration) error {
	query := `
		INSERT INTO "vasst_ca".organization_integration (
			organization_integration_id, organization_id, integration_id,
			token, last_used_at, is_active, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`

	_, err := r.DB.ExecContext(ctx, query,
		integration.OrganizationIntegrationID,
		integration.OrganizationID,
		integration.IntegrationID,
		integration.Token,
		integration.LastUsedAt,
		integration.IsActive,
	)

	return err
}

func (r *organizationRepository) UpdateIntegration(ctx context.Context, integration *entities.OrganizationIntegration) error {
	query := `
		UPDATE "vasst_ca".organization_integration
		SET token = $1,
			last_used_at = $2,
			is_active = $3,
			updated_at = CURRENT_TIMESTAMP
		WHERE organization_integration_id = $4
	`

	result, err := r.DB.ExecContext(ctx, query,
		integration.Token,
		integration.LastUsedAt,
		integration.IsActive,
		integration.OrganizationIntegrationID,
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

func (r *organizationRepository) DeleteIntegration(ctx context.Context, integrationID uuid.UUID) error {
	query := `
		DELETE FROM "vasst_ca".organization_integration
		WHERE organization_integration_id = $1
	`

	result, err := r.DB.ExecContext(ctx, query, integrationID)
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

func (r *organizationRepository) ListIntegrationsByOrgID(ctx context.Context, orgID uuid.UUID) ([]*entities.OrganizationIntegration, error) {
	query := `
		SELECT organization_integration_id, organization_id, integration_id,
			token, last_used_at, is_active, created_at, updated_at
		FROM "vasst_ca".organization_integration
		WHERE organization_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.DB.QueryContext(ctx, query, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var integrations []*entities.OrganizationIntegration
	for rows.Next() {
		var integration entities.OrganizationIntegration
		err := rows.Scan(
			&integration.OrganizationIntegrationID,
			&integration.OrganizationID,
			&integration.IntegrationID,
			&integration.Token,
			&integration.LastUsedAt,
			&integration.IsActive,
			&integration.CreatedAt,
			&integration.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		integrations = append(integrations, &integration)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return integrations, nil
}

func (r *organizationRepository) FindIntegrationByID(ctx context.Context, integrationID uuid.UUID) (*entities.OrganizationIntegration, error) {
	query := `
		SELECT organization_integration_id, organization_id, integration_id,
			token, last_used_at, is_active, created_at, updated_at
		FROM "vasst_ca".organization_integration
		WHERE organization_integration_id = $1
	`

	var integration entities.OrganizationIntegration
	err := r.DB.QueryRowContext(ctx, query, integrationID).Scan(
		&integration.OrganizationIntegrationID,
		&integration.OrganizationID,
		&integration.IntegrationID,
		&integration.Token,
		&integration.LastUsedAt,
		&integration.IsActive,
		&integration.CreatedAt,
		&integration.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &integration, nil
}

func (r *organizationRepository) FindOrganizationByKey(ctx context.Context, key string) (*entities.Organization, error) {
	query := `
		SELECT organization_id, organization_code, name, contact_name, phone_number,
			email, whatsapp_number, organization_type, organization_category_id,
			address, city, province, postal_code, country, status,
			created_at, updated_at
		FROM "vasst_ca".organization
		WHERE api_key = $1
	`

	var org entities.Organization
	err := r.DB.QueryRowContext(ctx, query, key).Scan(
		&org.OrganizationID,
		&org.OrganizationCode,
		&org.Name,
		&org.ContactName,
		&org.PhoneNumber,
		&org.Email,
		&org.WhatsappNumber,
		&org.OrganizationType,
		&org.CategoryID,
		&org.Address,
		&org.City,
		&org.Province,
		&org.PostalCode,
		&org.Country,
		&org.Status,
		&org.CreatedAt,
		&org.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &org, nil
}

func (r *organizationRepository) FindIntegrationTokenByOrgIDAndType(ctx context.Context, orgID uuid.UUID, integrationType string) (string, error) {
	query := `
		SELECT token, last_used_at
		FROM "vasst_ca".organization_integration oi
		JOIN "vasst_ca".integration i
		ON oi.integration_id = i.integration_id
		WHERE oi.organization_id = $1 AND i.integration_name = $2
	`

	var token string
	var lastUsedAt time.Time
	err := r.DB.QueryRowContext(ctx, query, orgID, integrationType).Scan(
		&token,
		&lastUsedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", err
	}

	// update last used at
	_, err = r.DB.ExecContext(ctx, `
		UPDATE "vasst_ca".organization_integration
		SET last_used_at = CURRENT_TIMESTAMP
		WHERE organization_integration_id = $1
	`, orgID)

	return token, nil
}

func (r *organizationRepository) FindIntegrationByOrgIDAndType(ctx context.Context, orgID uuid.UUID, integrationType string) (*entities.OrganizationIntegration, error) {
	query := `
		SELECT oi.organization_integration_id, oi.last_used_at, oi.is_active, oi.is_ai_enabled, oi.updated_at
		FROM "vasst_ca".organization_integration oi
		JOIN "vasst_ca".integration i
		ON oi.integration_id = i.integration_id
		WHERE oi.organization_id = $1 AND i.integration_name = $2
	`

	var integration entities.OrganizationIntegration
	err := r.DB.QueryRowContext(ctx, query, orgID, integrationType).Scan(
		&integration.OrganizationIntegrationID,
		&integration.LastUsedAt,
		&integration.IsActive,
		&integration.IsAiEnabled,
		&integration.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &integration, nil
}

func (r *organizationRepository) FindDefaultUserByOrgID(ctx context.Context, orgID uuid.UUID) (*entities.User, error) {
	query := `
		SELECT user_id, user_fullname, phone_number, username, updated_at
		FROM "vasst_ca".user
		WHERE organization_id = $1
	`

	fmt.Println("query", query)
	fmt.Println("orgID", orgID)

	var user entities.User
	err := r.DB.QueryRowContext(ctx, query, orgID).Scan(
		&user.UserID,
		&user.UserFullName,
		&user.PhoneNumber,
		&user.Username,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}
