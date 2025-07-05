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
	conversationRepository struct {
		*postgres.Postgres
	}

	// ConversationRepository defines methods for interacting with conversations in the database
	ConversationRepository interface {
		// Create creates a new conversation
		Create(ctx context.Context, conversation *entities.Conversation) error

		// Update updates a conversation
		Update(ctx context.Context, conversation *entities.Conversation) error

		// UpdateLastMessage updates last message information
		UpdateLastMessage(ctx context.Context, conversationID uuid.UUID, input *entities.UpdateConversationLastMessageInput) error

		// Delete deletes a conversation
		Delete(ctx context.Context, conversationID uuid.UUID) error

		// FindByID returns a conversation by ID (organization-scoped)
		FindByID(ctx context.Context, conversationID, organizationID uuid.UUID) (*entities.Conversation, error)

		// FindActiveByOrganizationUserContactMedium finds the active conversation for specific parameters
		FindActiveByOrganizationUserContactMedium(ctx context.Context, organizationID, userID, contactID uuid.UUID, mediumID int) (*entities.Conversation, error)

		// ListByOrganization returns conversations for an organization (organization-scoped)
		ListByOrganization(ctx context.Context, organizationID uuid.UUID, limit, offset int) ([]*entities.ConversationListResponse, error)

		// ListByOrganizationWithFilters returns conversations with filters
		ListByOrganizationWithFilters(ctx context.Context, organizationID uuid.UUID, status, priority *int, isActive *bool, limit, offset int) ([]*entities.ConversationListResponse, error)

		// FindByUserID returns conversations for a user (organization-scoped)
		FindByUserID(ctx context.Context, organizationID, userID uuid.UUID, limit, offset int) ([]*entities.ConversationListResponse, error)

		// FindByContactID returns conversations for a contact (organization-scoped)
		FindByContactID(ctx context.Context, organizationID, contactID uuid.UUID, limit, offset int) ([]*entities.ConversationListResponse, error)

		// FindByStatus returns conversations by status (organization-scoped)
		FindByStatus(ctx context.Context, organizationID uuid.UUID, status int, limit, offset int) ([]*entities.ConversationListResponse, error)

		// FindByPriority returns conversations by priority (organization-scoped)
		FindByPriority(ctx context.Context, organizationID uuid.UUID, priority int, limit, offset int) ([]*entities.ConversationListResponse, error)

		// CountByOrganization returns total count of conversations for an organization
		CountByOrganization(ctx context.Context, organizationID uuid.UUID) (int, error)

		// DeactivateOtherConversations deactivates other conversations for the same parameters
		DeactivateOtherConversations(ctx context.Context, organizationID, userID, contactID uuid.UUID, mediumID int, excludeConversationID uuid.UUID) error
	}
)

// NewConversationRepository creates a new ConversationRepository
func NewConversationRepository(pg *postgres.Postgres) ConversationRepository {
	return &conversationRepository{pg}
}

// Create creates a new conversation
func (r *conversationRepository) Create(ctx context.Context, conversation *entities.Conversation) error {
	query := `
		INSERT INTO "vasst_ca".conversation (
			organization_id, user_id, contact_id, medium_id, status, priority, 
			ai_enabled, ai_config, metadata, last_message_by_id, last_message_by_type, 
			last_message_by_name, last_message_content, last_message_type_id, last_message_media_url
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		RETURNING conversation_id, created_at, updated_at
	`

	// Convert AI config to JSON
	aiConfigJSON, err := json.Marshal(conversation.AIConfig)
	if err != nil {
		return err
	}

	// Convert metadata to JSON
	metadataJSON, err := json.Marshal(conversation.Metadata)
	if err != nil {
		return err
	}

	err = r.DB.QueryRowContext(ctx, query,
		conversation.OrganizationID,
		conversation.UserID,
		conversation.ContactID,
		conversation.MediumID,
		conversation.Status,
		conversation.Priority,
		conversation.AIEnabled,
		aiConfigJSON,
		metadataJSON,
		conversation.LastMessageByID,
		conversation.LastMessageByType,
		conversation.LastMessageByName,
		conversation.LastMessageContent,
		conversation.LastMessageTypeID,
		conversation.LastMessageMediaURL,
	).Scan(&conversation.ConversationID, &conversation.CreatedAt, &conversation.UpdatedAt)

	return err
}

// Update updates a conversation
func (r *conversationRepository) Update(ctx context.Context, conversation *entities.Conversation) error {
	query := `
		UPDATE "vasst_ca".conversation
		SET user_id = $1,
			medium_id = $2,
			status = $3,
			priority = $4,
			ai_enabled = $5,
			ai_config = $6,
			metadata = $7,
			is_active = $8,
			is_archived = $9,
			updated_at = CURRENT_TIMESTAMP
		WHERE conversation_id = $10 AND organization_id = $11
	`

	// Convert AI config to JSON
	aiConfigJSON, err := json.Marshal(conversation.AIConfig)
	if err != nil {
		return err
	}

	// Convert metadata to JSON
	metadataJSON, err := json.Marshal(conversation.Metadata)
	if err != nil {
		return err
	}

	result, err := r.DB.ExecContext(ctx, query,
		conversation.UserID,
		conversation.MediumID,
		conversation.Status,
		conversation.Priority,
		conversation.AIEnabled,
		aiConfigJSON,
		metadataJSON,
		conversation.IsActive,
		conversation.IsArchived,
		conversation.ConversationID,
		conversation.OrganizationID,
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

// UpdateLastMessage updates last message information
func (r *conversationRepository) UpdateLastMessage(ctx context.Context, conversationID uuid.UUID, input *entities.UpdateConversationLastMessageInput) error {
	query := `
		UPDATE "vasst_ca".conversation
		SET last_message_at = $1,
			last_message_by_id = $2,
			last_message_by_type = $3,
			last_message_by_name = $4,
			last_message_content = $5,
			last_message_type_id = $6,
			last_message_media_url = $7,
			updated_at = CURRENT_TIMESTAMP
		WHERE conversation_id = $8
	`

	result, err := r.DB.ExecContext(ctx, query,
		input.LastMessageAt,
		input.LastMessageByID,
		input.LastMessageByType,
		input.LastMessageByName,
		input.LastMessageContent,
		input.LastMessageTypeID,
		input.LastMessageMediaURL,
		conversationID,
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

// Delete deletes a conversation
func (r *conversationRepository) Delete(ctx context.Context, conversationID uuid.UUID) error {
	query := `
		DELETE FROM "vasst_ca".conversation
		WHERE conversation_id = $1
	`

	result, err := r.DB.ExecContext(ctx, query, conversationID)
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

// FindByID returns a conversation by ID (organization-scoped)
func (r *conversationRepository) FindByID(ctx context.Context, conversationID, organizationID uuid.UUID) (*entities.Conversation, error) {
	query := `
		SELECT conversation_id, organization_id, user_id, contact_id, medium_id, is_active, 
		       is_archived, is_deleted, status, priority, ai_enabled, ai_config, metadata,
		       last_message_at, last_human_message_at, last_ai_message_at, last_message_by_id,
		       last_message_by_type, last_message_by_name, last_message_content, last_message_type_id,
		       last_message_media_url, created_at, updated_at
		FROM "vasst_ca".conversation
		WHERE conversation_id = $1 AND organization_id = $2 AND is_deleted = false
	`

	var conversation entities.Conversation
	err := r.DB.QueryRowContext(ctx, query, conversationID, organizationID).Scan(
		&conversation.ConversationID,
		&conversation.OrganizationID,
		&conversation.UserID,
		&conversation.ContactID,
		&conversation.MediumID,
		&conversation.IsActive,
		&conversation.IsArchived,
		&conversation.IsDeleted,
		&conversation.Status,
		&conversation.Priority,
		&conversation.AIEnabled,
		&conversation.AIConfig,
		&conversation.Metadata,
		&conversation.LastMessageAt,
		&conversation.LastHumanMessageAt,
		&conversation.LastAIMessageAt,
		&conversation.LastMessageByID,
		&conversation.LastMessageByType,
		&conversation.LastMessageByName,
		&conversation.LastMessageContent,
		&conversation.LastMessageTypeID,
		&conversation.LastMessageMediaURL,
		&conversation.CreatedAt,
		&conversation.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &conversation, nil
}

// FindActiveByOrganizationUserContactMedium finds the active conversation for specific parameters
func (r *conversationRepository) FindActiveByOrganizationUserContactMedium(ctx context.Context, organizationID, userID, contactID uuid.UUID, mediumID int) (*entities.Conversation, error) {
	query := `
		SELECT conversation_id, organization_id, user_id, contact_id, medium_id, is_active, 
		       is_archived, is_deleted, status, priority, ai_enabled, ai_config, metadata,
		       last_message_at, last_human_message_at, last_ai_message_at, last_message_by_id,
		       last_message_by_type, last_message_by_name, last_message_content, last_message_type_id,
		       last_message_media_url, created_at, updated_at
		FROM "vasst_ca".conversation
		WHERE organization_id = $1 AND user_id = $2 AND contact_id = $3 AND medium_id = $4 
		      AND is_active = true AND is_deleted = false
		ORDER BY created_at DESC
		LIMIT 1
	`

	var conversation entities.Conversation
	err := r.DB.QueryRowContext(ctx, query, organizationID, userID, contactID, mediumID).Scan(
		&conversation.ConversationID,
		&conversation.OrganizationID,
		&conversation.UserID,
		&conversation.ContactID,
		&conversation.MediumID,
		&conversation.IsActive,
		&conversation.IsArchived,
		&conversation.IsDeleted,
		&conversation.Status,
		&conversation.Priority,
		&conversation.AIEnabled,
		&conversation.AIConfig,
		&conversation.Metadata,
		&conversation.LastMessageAt,
		&conversation.LastHumanMessageAt,
		&conversation.LastAIMessageAt,
		&conversation.LastMessageByID,
		&conversation.LastMessageByType,
		&conversation.LastMessageByName,
		&conversation.LastMessageContent,
		&conversation.LastMessageTypeID,
		&conversation.LastMessageMediaURL,
		&conversation.CreatedAt,
		&conversation.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &conversation, nil
}

// ListByOrganization returns conversations for an organization (organization-scoped)
func (r *conversationRepository) ListByOrganization(ctx context.Context, organizationID uuid.UUID, limit, offset int) ([]*entities.ConversationListResponse, error) {
	query := `
		SELECT 
			c.conversation_id,
			c.contact_id,
			co.contact_name,
			c.medium_id,
			m.name as medium_name,
			c.status,
			c.priority,
			c.ai_enabled,
			c.last_message_at,
			c.last_message_content,
			c.last_message_by_name,
			c.created_at,
			c.updated_at,
			COALESCE(unread.unread_count, 0) as unread_count
		FROM "vasst_ca".conversation c
		LEFT JOIN "vasst_ca".contact co ON c.contact_id = co.contact_id
		LEFT JOIN "vasst_ca".medium m ON c.medium_id = m.medium_id
		LEFT JOIN (
			SELECT conversation_id, COUNT(*) as unread_count
			FROM "vasst_ca".messages
			WHERE status = 0 AND direction = 'i'
			GROUP BY conversation_id
		) unread ON c.conversation_id = unread.conversation_id
		WHERE c.organization_id = $1 AND c.is_deleted = false
		ORDER BY c.last_message_at DESC NULLS LAST, c.created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.DB.QueryContext(ctx, query, organizationID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var conversations []*entities.ConversationListResponse
	for rows.Next() {
		var conv entities.ConversationListResponse
		err := rows.Scan(
			&conv.ConversationID,
			&conv.ContactID,
			&conv.ContactName,
			&conv.MediumID,
			&conv.MediumName,
			&conv.Status,
			&conv.Priority,
			&conv.AIEnabled,
			&conv.LastMessageAt,
			&conv.LastMessageContent,
			&conv.LastMessageByName,
			&conv.CreatedAt,
			&conv.UpdatedAt,
			&conv.UnreadCount,
		)
		if err != nil {
			return nil, err
		}
		conversations = append(conversations, &conv)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return conversations, nil
}

// ListByOrganizationWithFilters returns conversations with filters
func (r *conversationRepository) ListByOrganizationWithFilters(ctx context.Context, organizationID uuid.UUID, status, priority *int, isActive *bool, limit, offset int) ([]*entities.ConversationListResponse, error) {
	query := `
		SELECT 
			c.conversation_id,
			c.contact_id,
			co.contact_name,
			c.medium_id,
			m.name as medium_name,
			c.status,
			c.priority,
			c.ai_enabled,
			c.last_message_at,
			c.last_message_content,
			c.last_message_by_name,
			c.created_at,
			c.updated_at,
			COALESCE(unread.unread_count, 0) as unread_count
		FROM "vasst_ca".conversation c
		LEFT JOIN "vasst_ca".contact co ON c.contact_id = co.contact_id
		LEFT JOIN "vasst_ca".medium m ON c.medium_id = m.medium_id
		LEFT JOIN (
			SELECT conversation_id, COUNT(*) as unread_count
			FROM "vasst_ca".messages
			WHERE status = 0 AND direction = 'i'
			GROUP BY conversation_id
		) unread ON c.conversation_id = unread.conversation_id
		WHERE c.organization_id = $1 AND c.is_deleted = false
	`

	args := []interface{}{organizationID}
	argCount := 1

	if status != nil {
		argCount++
		query += ` AND c.status = $` + string(rune(argCount+'0'))
		args = append(args, *status)
	}

	if priority != nil {
		argCount++
		query += ` AND c.priority = $` + string(rune(argCount+'0'))
		args = append(args, *priority)
	}

	if isActive != nil {
		argCount++
		query += ` AND c.is_active = $` + string(rune(argCount+'0'))
		args = append(args, *isActive)
	}

	query += ` ORDER BY c.last_message_at DESC NULLS LAST, c.created_at DESC LIMIT $` + string(rune(argCount+1+'0')) + ` OFFSET $` + string(rune(argCount+2+'0'))
	args = append(args, limit, offset)

	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var conversations []*entities.ConversationListResponse
	for rows.Next() {
		var conv entities.ConversationListResponse
		err := rows.Scan(
			&conv.ConversationID,
			&conv.ContactID,
			&conv.ContactName,
			&conv.MediumID,
			&conv.MediumName,
			&conv.Status,
			&conv.Priority,
			&conv.AIEnabled,
			&conv.LastMessageAt,
			&conv.LastMessageContent,
			&conv.LastMessageByName,
			&conv.CreatedAt,
			&conv.UpdatedAt,
			&conv.UnreadCount,
		)
		if err != nil {
			return nil, err
		}
		conversations = append(conversations, &conv)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return conversations, nil
}

// FindByUserID returns conversations for a user (organization-scoped)
func (r *conversationRepository) FindByUserID(ctx context.Context, organizationID, userID uuid.UUID, limit, offset int) ([]*entities.ConversationListResponse, error) {
	query := `
		SELECT 
			c.conversation_id,
			c.contact_id,
			co.contact_name,
			c.medium_id,
			m.name as medium_name,
			c.status,
			c.priority,
			c.ai_enabled,
			c.last_message_at,
			c.last_message_content,
			c.last_message_by_name,
			c.created_at,
			c.updated_at,
			COALESCE(unread.unread_count, 0) as unread_count
		FROM "vasst_ca".conversation c
		LEFT JOIN "vasst_ca".contact co ON c.contact_id = co.contact_id
		LEFT JOIN "vasst_ca".medium m ON c.medium_id = m.medium_id
		LEFT JOIN (
			SELECT conversation_id, COUNT(*) as unread_count
			FROM "vasst_ca".messages
			WHERE status = 0 AND direction = 'i'
			GROUP BY conversation_id
		) unread ON c.conversation_id = unread.conversation_id
		WHERE c.organization_id = $1 AND c.user_id = $2 AND c.is_deleted = false
		ORDER BY c.last_message_at DESC NULLS LAST, c.created_at DESC
		LIMIT $3 OFFSET $4
	`

	rows, err := r.DB.QueryContext(ctx, query, organizationID, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var conversations []*entities.ConversationListResponse
	for rows.Next() {
		var conv entities.ConversationListResponse
		err := rows.Scan(
			&conv.ConversationID,
			&conv.ContactID,
			&conv.ContactName,
			&conv.MediumID,
			&conv.MediumName,
			&conv.Status,
			&conv.Priority,
			&conv.AIEnabled,
			&conv.LastMessageAt,
			&conv.LastMessageContent,
			&conv.LastMessageByName,
			&conv.CreatedAt,
			&conv.UpdatedAt,
			&conv.UnreadCount,
		)
		if err != nil {
			return nil, err
		}
		conversations = append(conversations, &conv)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return conversations, nil
}

// FindByContactID returns conversations for a contact (organization-scoped)
func (r *conversationRepository) FindByContactID(ctx context.Context, organizationID, contactID uuid.UUID, limit, offset int) ([]*entities.ConversationListResponse, error) {
	query := `
		SELECT 
			c.conversation_id,
			c.contact_id,
			co.contact_name,
			c.medium_id,
			m.name as medium_name,
			c.status,
			c.priority,
			c.ai_enabled,
			c.last_message_at,
			c.last_message_content,
			c.last_message_by_name,
			c.created_at,
			c.updated_at,
			COALESCE(unread.unread_count, 0) as unread_count
		FROM "vasst_ca".conversation c
		LEFT JOIN "vasst_ca".contact co ON c.contact_id = co.contact_id
		LEFT JOIN "vasst_ca".medium m ON c.medium_id = m.medium_id
		LEFT JOIN (
			SELECT conversation_id, COUNT(*) as unread_count
			FROM "vasst_ca".messages
			WHERE status = 0 AND direction = 'i'
			GROUP BY conversation_id
		) unread ON c.conversation_id = unread.conversation_id
		WHERE c.organization_id = $1 AND c.contact_id = $2 AND c.is_deleted = false
		ORDER BY c.last_message_at DESC NULLS LAST, c.created_at DESC
		LIMIT $3 OFFSET $4
	`

	rows, err := r.DB.QueryContext(ctx, query, organizationID, contactID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var conversations []*entities.ConversationListResponse
	for rows.Next() {
		var conv entities.ConversationListResponse
		err := rows.Scan(
			&conv.ConversationID,
			&conv.ContactID,
			&conv.ContactName,
			&conv.MediumID,
			&conv.MediumName,
			&conv.Status,
			&conv.Priority,
			&conv.AIEnabled,
			&conv.LastMessageAt,
			&conv.LastMessageContent,
			&conv.LastMessageByName,
			&conv.CreatedAt,
			&conv.UpdatedAt,
			&conv.UnreadCount,
		)
		if err != nil {
			return nil, err
		}
		conversations = append(conversations, &conv)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return conversations, nil
}

// FindByStatus returns conversations by status (organization-scoped)
func (r *conversationRepository) FindByStatus(ctx context.Context, organizationID uuid.UUID, status int, limit, offset int) ([]*entities.ConversationListResponse, error) {
	statusPtr := &status
	return r.ListByOrganizationWithFilters(ctx, organizationID, statusPtr, nil, nil, limit, offset)
}

// FindByPriority returns conversations by priority (organization-scoped)
func (r *conversationRepository) FindByPriority(ctx context.Context, organizationID uuid.UUID, priority int, limit, offset int) ([]*entities.ConversationListResponse, error) {
	priorityPtr := &priority
	return r.ListByOrganizationWithFilters(ctx, organizationID, nil, priorityPtr, nil, limit, offset)
}

// CountByOrganization returns total count of conversations for an organization
func (r *conversationRepository) CountByOrganization(ctx context.Context, organizationID uuid.UUID) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM "vasst_ca".conversation
		WHERE organization_id = $1 AND is_deleted = false
	`

	var count int
	err := r.DB.QueryRowContext(ctx, query, organizationID).Scan(&count)
	return count, err
}

// DeactivateOtherConversations deactivates other conversations for the same parameters
func (r *conversationRepository) DeactivateOtherConversations(ctx context.Context, organizationID, userID, contactID uuid.UUID, mediumID int, excludeConversationID uuid.UUID) error {
	query := `
		UPDATE "vasst_ca".conversation
		SET is_active = false, updated_at = CURRENT_TIMESTAMP
		WHERE organization_id = $1 AND user_id = $2 AND contact_id = $3 AND medium_id = $4 
		      AND conversation_id != $5 AND is_active = true
	`

	_, err := r.DB.ExecContext(ctx, query, organizationID, userID, contactID, mediumID, excludeConversationID)
	return err
}
