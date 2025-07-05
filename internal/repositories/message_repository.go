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
	messageRepository struct {
		*postgres.Postgres
	}

	// MessageRepository defines methods for interacting with messages in the database
	MessageRepository interface {
		// Create creates a new message
		Create(ctx context.Context, message *entities.Message) error

		// Update updates a message
		Update(ctx context.Context, message *entities.Message) error

		// UpdateStatus updates only the status of a message with proper timestamp handling
		UpdateStatus(ctx context.Context, messageID uuid.UUID, status int, failureReason *string) error

		// Delete deletes a message
		Delete(ctx context.Context, messageID uuid.UUID) error

		// FindByConversationAndOrganization returns messages for a specific conversation and organization
		FindByConversationAndOrganization(ctx context.Context, conversationID, organizationID uuid.UUID, limit, offset int) ([]*entities.Message, error)

		// FindByID returns a message by ID
		FindByID(ctx context.Context, messageID uuid.UUID) (*entities.Message, error)

		// FindByConversationID returns all messages for a conversation
		FindByConversationID(ctx context.Context, conversationID uuid.UUID, limit, offset int) ([]*entities.Message, error)

		// FindByOrganizationID returns all messages for an organization
		FindByOrganizationID(ctx context.Context, organizationID uuid.UUID, limit, offset int) ([]*entities.Message, error)

		// FindPendingMessages returns all pending messages
		FindPendingMessages(ctx context.Context, limit, offset int) ([]*entities.Message, error)

		// FindByStatus returns messages by status
		FindByStatus(ctx context.Context, status int, limit, offset int) ([]*entities.Message, error)

		// FindBySenderType returns messages by sender type
		FindBySenderType(ctx context.Context, senderTypeID int, limit, offset int) ([]*entities.Message, error)

		// FindBySenderID returns messages by sender ID
		FindBySenderID(ctx context.Context, senderID uuid.UUID, limit, offset int) ([]*entities.Message, error)
	}
)

// NewMessageRepository creates a new MessageRepository
func NewMessageRepository(pg *postgres.Postgres) MessageRepository {
	return &messageRepository{pg}
}

// Create creates a new message (message_id is auto-generated)
func (r *messageRepository) Create(ctx context.Context, message *entities.Message) error {
	query := `
		INSERT INTO "vasst_ca".messages (
			conversation_id, organization_id, sender_type_id, sender_id, direction, 
			message_type_id, content, media_url, attachments, is_broadcast, 
			is_order_message, metadata, ai_generated, ai_confidence_score, status
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		RETURNING message_id, created_at, updated_at
	`

	// Convert attachments to JSON
	attachmentsJSON, err := json.Marshal(message.Attachments)
	if err != nil {
		return err
	}

	// Convert metadata to JSON
	metadataJSON, err := json.Marshal(message.Metadata)
	if err != nil {
		return err
	}

	err = r.DB.QueryRowContext(ctx, query,
		message.ConversationID,
		message.OrganizationID,
		message.SenderTypeID,
		message.SenderID,
		message.Direction,
		message.MessageTypeID,
		message.Content,
		message.MediaURL,
		attachmentsJSON,
		message.IsBroadcast,
		message.IsOrderMessage,
		metadataJSON,
		message.AIGenerated,
		message.AIConfidenceScore,
		message.Status,
	).Scan(&message.MessageID, &message.CreatedAt, &message.UpdatedAt)

	return err
}

// Update updates a message
func (r *messageRepository) Update(ctx context.Context, message *entities.Message) error {
	query := `
		UPDATE "vasst_ca".messages
		SET sender_type_id = $1,
			sender_id = $2,
			direction = $3,
			message_type_id = $4,
			content = $5,
			media_url = $6,
			attachments = $7,
			is_broadcast = $8,
			is_order_message = $9,
			metadata = $10,
			ai_generated = $11,
			ai_confidence_score = $12,
			failure_reason = $13,
			updated_at = CURRENT_TIMESTAMP
		WHERE message_id = $14
	`

	// Convert attachments to JSON
	attachmentsJSON, err := json.Marshal(message.Attachments)
	if err != nil {
		return err
	}

	// Convert metadata to JSON
	metadataJSON, err := json.Marshal(message.Metadata)
	if err != nil {
		return err
	}

	result, err := r.DB.ExecContext(ctx, query,
		message.SenderTypeID,
		message.SenderID,
		message.Direction,
		message.MessageTypeID,
		message.Content,
		message.MediaURL,
		attachmentsJSON,
		message.IsBroadcast,
		message.IsOrderMessage,
		metadataJSON,
		message.AIGenerated,
		message.AIConfidenceScore,
		message.FailureReason,
		message.MessageID,
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

// UpdateStatus updates only the status of a message with proper timestamp handling
func (r *messageRepository) UpdateStatus(ctx context.Context, messageID uuid.UUID, status int, failureReason *string) error {
	// Build dynamic query based on status
	var query string
	var args []interface{}

	switch status {
	case int(entities.MessageStatusDelivered):
		query = `
			UPDATE "vasst_ca".messages
			SET status = $1, delivered_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
			WHERE message_id = $2
		`
		args = []interface{}{status, messageID}
	case int(entities.MessageStatusRead):
		query = `
			UPDATE "vasst_ca".messages
			SET status = $1, read_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
			WHERE message_id = $2
		`
		args = []interface{}{status, messageID}
	case int(entities.MessageStatusFailed):
		query = `
			UPDATE "vasst_ca".messages
			SET status = $1, failed_at = CURRENT_TIMESTAMP, failure_reason = $2, updated_at = CURRENT_TIMESTAMP
			WHERE message_id = $3
		`
		args = []interface{}{status, failureReason, messageID}
	default:
		query = `
			UPDATE "vasst_ca".messages
			SET status = $1, updated_at = CURRENT_TIMESTAMP
			WHERE message_id = $2
		`
		args = []interface{}{status, messageID}
	}

	result, err := r.DB.ExecContext(ctx, query, args...)
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

// Delete deletes a message
func (r *messageRepository) Delete(ctx context.Context, messageID uuid.UUID) error {
	query := `
		DELETE FROM "vasst_ca".messages
		WHERE message_id = $1
	`

	result, err := r.DB.ExecContext(ctx, query, messageID)
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

// FindByConversationAndOrganization returns messages for a specific conversation and organization
func (r *messageRepository) FindByConversationAndOrganization(ctx context.Context, conversationID, organizationID uuid.UUID, limit, offset int) ([]*entities.Message, error) {
	query := `
		SELECT message_id, conversation_id, organization_id, sender_type_id, sender_id, direction, 
		       message_type_id, content, media_url, attachments, is_broadcast, is_order_message, 
		       metadata, read_at, delivered_at, failed_at, failure_reason, ai_generated, 
		       ai_confidence_score, status, created_at, updated_at
		FROM "vasst_ca".messages
		WHERE conversation_id = $1 AND organization_id = $2
		ORDER BY created_at ASC
		LIMIT $3 OFFSET $4
	`

	rows, err := r.DB.QueryContext(ctx, query, conversationID, organizationID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*entities.Message
	for rows.Next() {
		message, err := r.scanMessage(rows)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}

// FindByID returns a message by ID
func (r *messageRepository) FindByID(ctx context.Context, messageID uuid.UUID) (*entities.Message, error) {
	query := `
		SELECT message_id, conversation_id, organization_id, sender_type_id, sender_id, direction, 
		       message_type_id, content, media_url, attachments, is_broadcast, is_order_message, 
		       metadata, read_at, delivered_at, failed_at, failure_reason, ai_generated, 
		       ai_confidence_score, status, created_at, updated_at
		FROM "vasst_ca".messages
		WHERE message_id = $1
	`

	var message entities.Message
	err := r.DB.QueryRowContext(ctx, query, messageID).Scan(
		&message.MessageID,
		&message.ConversationID,
		&message.OrganizationID,
		&message.SenderTypeID,
		&message.SenderID,
		&message.Direction,
		&message.MessageTypeID,
		&message.Content,
		&message.MediaURL,
		&message.Attachments,
		&message.IsBroadcast,
		&message.IsOrderMessage,
		&message.Metadata,
		&message.ReadAt,
		&message.DeliveredAt,
		&message.FailedAt,
		&message.FailureReason,
		&message.AIGenerated,
		&message.AIConfidenceScore,
		&message.Status,
		&message.CreatedAt,
		&message.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &message, nil
}

// FindByConversationID returns all messages for a conversation
func (r *messageRepository) FindByConversationID(ctx context.Context, conversationID uuid.UUID, limit, offset int) ([]*entities.Message, error) {
	query := `
		SELECT message_id, conversation_id, organization_id, sender_type_id, sender_id, direction, 
		       message_type_id, content, media_url, attachments, is_broadcast, is_order_message, 
		       metadata, read_at, delivered_at, failed_at, failure_reason, ai_generated, 
		       ai_confidence_score, status, created_at, updated_at
		FROM "vasst_ca".messages
		WHERE conversation_id = $1
		ORDER BY created_at ASC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.DB.QueryContext(ctx, query, conversationID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*entities.Message
	for rows.Next() {
		message, err := r.scanMessage(rows)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}

// FindByOrganizationID returns all messages for an organization
func (r *messageRepository) FindByOrganizationID(ctx context.Context, organizationID uuid.UUID, limit, offset int) ([]*entities.Message, error) {
	query := `
		SELECT message_id, conversation_id, organization_id, sender_type_id, sender_id, direction, 
		       message_type_id, content, media_url, attachments, is_broadcast, is_order_message, 
		       metadata, read_at, delivered_at, failed_at, failure_reason, ai_generated, 
		       ai_confidence_score, status, created_at, updated_at
		FROM "vasst_ca".messages
		WHERE organization_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.DB.QueryContext(ctx, query, organizationID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*entities.Message
	for rows.Next() {
		message, err := r.scanMessage(rows)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}

// FindPendingMessages returns all pending messages
func (r *messageRepository) FindPendingMessages(ctx context.Context, limit, offset int) ([]*entities.Message, error) {
	query := `
		SELECT message_id, conversation_id, organization_id, sender_type_id, sender_id, direction, 
		       message_type_id, content, media_url, attachments, is_broadcast, is_order_message, 
		       metadata, read_at, delivered_at, failed_at, failure_reason, ai_generated, 
		       ai_confidence_score, status, created_at, updated_at
		FROM "vasst_ca".messages
		WHERE status = 0
		ORDER BY created_at ASC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.DB.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*entities.Message
	for rows.Next() {
		message, err := r.scanMessage(rows)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}

// FindByStatus returns messages by status
func (r *messageRepository) FindByStatus(ctx context.Context, status int, limit, offset int) ([]*entities.Message, error) {
	query := `
		SELECT message_id, conversation_id, organization_id, sender_type_id, sender_id, direction, 
		       message_type_id, content, media_url, attachments, is_broadcast, is_order_message, 
		       metadata, read_at, delivered_at, failed_at, failure_reason, ai_generated, 
		       ai_confidence_score, status, created_at, updated_at
		FROM "vasst_ca".messages
		WHERE status = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.DB.QueryContext(ctx, query, status, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*entities.Message
	for rows.Next() {
		message, err := r.scanMessage(rows)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}

// FindBySenderType returns messages by sender type
func (r *messageRepository) FindBySenderType(ctx context.Context, senderTypeID int, limit, offset int) ([]*entities.Message, error) {
	query := `
		SELECT message_id, conversation_id, organization_id, sender_type_id, sender_id, direction, 
		       message_type_id, content, media_url, attachments, is_broadcast, is_order_message, 
		       metadata, read_at, delivered_at, failed_at, failure_reason, ai_generated, 
		       ai_confidence_score, status, created_at, updated_at
		FROM "vasst_ca".messages
		WHERE sender_type_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.DB.QueryContext(ctx, query, senderTypeID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*entities.Message
	for rows.Next() {
		message, err := r.scanMessage(rows)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}

// FindBySenderID returns messages by sender ID
func (r *messageRepository) FindBySenderID(ctx context.Context, senderID uuid.UUID, limit, offset int) ([]*entities.Message, error) {
	query := `
		SELECT message_id, conversation_id, organization_id, sender_type_id, sender_id, direction, 
		       message_type_id, content, media_url, attachments, is_broadcast, is_order_message, 
		       metadata, read_at, delivered_at, failed_at, failure_reason, ai_generated, 
		       ai_confidence_score, status, created_at, updated_at
		FROM "vasst_ca".messages
		WHERE sender_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.DB.QueryContext(ctx, query, senderID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*entities.Message
	for rows.Next() {
		message, err := r.scanMessage(rows)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}

// scanMessage is a helper function to scan message data from database rows
func (r *messageRepository) scanMessage(rows *sql.Rows) (*entities.Message, error) {
	var message entities.Message
	err := rows.Scan(
		&message.MessageID,
		&message.ConversationID,
		&message.OrganizationID,
		&message.SenderTypeID,
		&message.SenderID,
		&message.Direction,
		&message.MessageTypeID,
		&message.Content,
		&message.MediaURL,
		&message.Attachments,
		&message.IsBroadcast,
		&message.IsOrderMessage,
		&message.Metadata,
		&message.ReadAt,
		&message.DeliveredAt,
		&message.FailedAt,
		&message.FailureReason,
		&message.AIGenerated,
		&message.AIConfidenceScore,
		&message.Status,
		&message.CreatedAt,
		&message.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &message, nil
}
