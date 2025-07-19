package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/vasst-id/vasst-expense-api/internal/entities"
	"github.com/vasst-id/vasst-expense-api/internal/utils/postgres"
)

type (
	messageRepository struct {
		*postgres.Postgres
	}

	MessageRepository interface {
		Create(ctx context.Context, message *entities.Message) (entities.Message, error)
		Update(ctx context.Context, message *entities.Message) (entities.Message, error)
		Delete(ctx context.Context, messageID uuid.UUID) error
		FindByID(ctx context.Context, messageID uuid.UUID) (*entities.Message, error)
		FindByConversationID(ctx context.Context, conversationID uuid.UUID, limit, offset int) ([]*entities.Message, error)
		FindByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entities.Message, error)
		FindWithFilters(ctx context.Context, params *entities.MessageListParams, limit, offset int) ([]*entities.Message, error)
		FindSimpleByConversationID(ctx context.Context, conversationID uuid.UUID, limit, offset int) ([]*entities.MessageSimple, error)
		CountByConversationID(ctx context.Context, conversationID uuid.UUID) (int64, error)
		CountWithFilters(ctx context.Context, params *entities.MessageListParams) (int64, error)
		MarkAsProcessed(ctx context.Context, messageID uuid.UUID, aiModel string, confidenceScore *float64) error
	}
)

// NewMessageRepository creates a new MessageRepository
func NewMessageRepository(pg *postgres.Postgres) MessageRepository {
	return &messageRepository{pg}
}

// Create creates a new message
func (r *messageRepository) Create(ctx context.Context, message *entities.Message) (entities.Message, error) {
	query := `
		INSERT INTO "vasst_expense".messages 
		(message_id, conversation_id, user_id, sender_type, direction, message_type, 
		 content, media_url, attachments, media_mime_type, transcription, ai_processed, 
		 ai_model, ai_confidence_score, related_transaction_id, scheduled_task_id, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, CURRENT_TIMESTAMP)
		RETURNING message_id, conversation_id, user_id, sender_type, direction, message_type,
		          content, media_url, attachments, media_mime_type, transcription, ai_processed,
		          ai_model, ai_confidence_score, related_transaction_id, scheduled_task_id, created_at
	`

	var createdMessage entities.Message
	err := r.DB.QueryRowContext(ctx, query,
		message.MessageID, message.ConversationID, message.UserID, message.SenderType,
		message.Direction, message.MessageType, message.Content, message.MediaURL,
		message.Attachments, message.MediaMimeType, message.Transcription, message.AIProcessed,
		message.AIModel, message.AIConfidenceScore, message.RelatedTransactionID, message.ScheduledTaskID,
	).Scan(
		&createdMessage.MessageID, &createdMessage.ConversationID, &createdMessage.UserID, &createdMessage.SenderType,
		&createdMessage.Direction, &createdMessage.MessageType, &createdMessage.Content, &createdMessage.MediaURL,
		&createdMessage.Attachments, &createdMessage.MediaMimeType, &createdMessage.Transcription, &createdMessage.AIProcessed,
		&createdMessage.AIModel, &createdMessage.AIConfidenceScore, &createdMessage.RelatedTransactionID, &createdMessage.ScheduledTaskID,
		&createdMessage.CreatedAt,
	)

	return createdMessage, err
}

// Update updates a message
func (r *messageRepository) Update(ctx context.Context, message *entities.Message) (entities.Message, error) {
	query := `
		UPDATE "vasst_expense".messages 
		SET content = $2, media_url = $3, attachments = $4, media_mime_type = $5,
		    transcription = $6, ai_processed = $7, ai_model = $8, ai_confidence_score = $9,
		    related_transaction_id = $10
		WHERE message_id = $1
		RETURNING message_id, conversation_id, user_id, sender_type, direction, message_type,
		          content, media_url, attachments, media_mime_type, transcription, ai_processed,
		          ai_model, ai_confidence_score, related_transaction_id, scheduled_task_id, created_at
	`

	var updatedMessage entities.Message
	err := r.DB.QueryRowContext(ctx, query,
		message.MessageID, message.Content, message.MediaURL, message.Attachments,
		message.MediaMimeType, message.Transcription, message.AIProcessed, message.AIModel,
		message.AIConfidenceScore, message.RelatedTransactionID,
	).Scan(
		&updatedMessage.MessageID, &updatedMessage.ConversationID, &updatedMessage.UserID, &updatedMessage.SenderType,
		&updatedMessage.Direction, &updatedMessage.MessageType, &updatedMessage.Content, &updatedMessage.MediaURL,
		&updatedMessage.Attachments, &updatedMessage.MediaMimeType, &updatedMessage.Transcription, &updatedMessage.AIProcessed,
		&updatedMessage.AIModel, &updatedMessage.AIConfidenceScore, &updatedMessage.RelatedTransactionID, &updatedMessage.ScheduledTaskID,
		&updatedMessage.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return entities.Message{}, sql.ErrNoRows // Message not found
		}
		return entities.Message{}, err
	}

	return updatedMessage, nil
}

// Delete deletes a message (hard delete)
func (r *messageRepository) Delete(ctx context.Context, messageID uuid.UUID) error {
	query := `
		DELETE FROM "vasst_expense".messages 
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

// FindByID finds a message by ID
func (r *messageRepository) FindByID(ctx context.Context, messageID uuid.UUID) (*entities.Message, error) {
	query := `
		SELECT message_id, conversation_id, user_id, sender_type, direction, message_type,
		       content, media_url, attachments, media_mime_type, transcription, ai_processed,
		       ai_model, ai_confidence_score, related_transaction_id, scheduled_task_id, created_at
		FROM "vasst_expense".messages 
		WHERE message_id = $1
	`

	var message entities.Message
	err := r.DB.QueryRowContext(ctx, query, messageID).Scan(
		&message.MessageID, &message.ConversationID, &message.UserID, &message.SenderType,
		&message.Direction, &message.MessageType, &message.Content, &message.MediaURL,
		&message.Attachments, &message.MediaMimeType, &message.Transcription, &message.AIProcessed,
		&message.AIModel, &message.AIConfidenceScore, &message.RelatedTransactionID, &message.ScheduledTaskID,
		&message.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &message, nil
}

// FindByConversationID finds messages by conversation ID with pagination
func (r *messageRepository) FindByConversationID(ctx context.Context, conversationID uuid.UUID, limit, offset int) ([]*entities.Message, error) {
	query := `
		SELECT message_id, conversation_id, user_id, sender_type, direction, message_type,
		       content, media_url, attachments, media_mime_type, transcription, ai_processed,
		       ai_model, ai_confidence_score, related_transaction_id, scheduled_task_id, created_at
		FROM "vasst_expense".messages 
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
		var message entities.Message
		err := rows.Scan(
			&message.MessageID, &message.ConversationID, &message.UserID, &message.SenderType,
			&message.Direction, &message.MessageType, &message.Content, &message.MediaURL,
			&message.Attachments, &message.MediaMimeType, &message.Transcription, &message.AIProcessed,
			&message.AIModel, &message.AIConfidenceScore, &message.RelatedTransactionID, &message.ScheduledTaskID,
			&message.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		messages = append(messages, &message)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}

// FindByUserID finds messages by user ID with pagination
func (r *messageRepository) FindByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entities.Message, error) {
	query := `
		SELECT message_id, conversation_id, user_id, sender_type, direction, message_type,
		       content, media_url, attachments, media_mime_type, transcription, ai_processed,
		       ai_model, ai_confidence_score, related_transaction_id, scheduled_task_id, created_at
		FROM "vasst_expense".messages 
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.DB.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*entities.Message
	for rows.Next() {
		var message entities.Message
		err := rows.Scan(
			&message.MessageID, &message.ConversationID, &message.UserID, &message.SenderType,
			&message.Direction, &message.MessageType, &message.Content, &message.MediaURL,
			&message.Attachments, &message.MediaMimeType, &message.Transcription, &message.AIProcessed,
			&message.AIModel, &message.AIConfidenceScore, &message.RelatedTransactionID, &message.ScheduledTaskID,
			&message.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		messages = append(messages, &message)
	}

	return messages, nil
}

// FindWithFilters finds messages with filtering and pagination
func (r *messageRepository) FindWithFilters(ctx context.Context, params *entities.MessageListParams, limit, offset int) ([]*entities.Message, error) {
	query := `
		SELECT message_id, conversation_id, user_id, sender_type, direction, message_type,
		       content, media_url, attachments, media_mime_type, transcription, ai_processed,
		       ai_model, ai_confidence_score, related_transaction_id, scheduled_task_id, created_at
		FROM "vasst_expense".messages 
		WHERE 1=1
	`

	args := []interface{}{}
	argIndex := 1

	// Add filters from params
	if params != nil {
		if params.ConversationID != nil {
			query += fmt.Sprintf(" AND conversation_id = $%d", argIndex)
			args = append(args, *params.ConversationID)
			argIndex++
		}
		if params.SenderType != nil {
			query += fmt.Sprintf(" AND sender_type = $%d", argIndex)
			args = append(args, *params.SenderType)
			argIndex++
		}
		if params.Direction != nil {
			query += fmt.Sprintf(" AND direction = $%d", argIndex)
			args = append(args, *params.Direction)
			argIndex++
		}
		if params.MessageType != nil {
			query += fmt.Sprintf(" AND message_type = $%d", argIndex)
			args = append(args, *params.MessageType)
			argIndex++
		}
		if params.AIProcessed != nil {
			query += fmt.Sprintf(" AND ai_processed = $%d", argIndex)
			args = append(args, *params.AIProcessed)
			argIndex++
		}
		if params.StartDate != nil {
			query += fmt.Sprintf(" AND created_at >= $%d", argIndex)
			args = append(args, *params.StartDate)
			argIndex++
		}
		if params.EndDate != nil {
			query += fmt.Sprintf(" AND created_at <= $%d", argIndex)
			args = append(args, *params.EndDate)
			argIndex++
		}
	}

	query += " ORDER BY created_at DESC"
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, limit, offset)

	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*entities.Message
	for rows.Next() {
		var message entities.Message
		err := rows.Scan(
			&message.MessageID, &message.ConversationID, &message.UserID, &message.SenderType,
			&message.Direction, &message.MessageType, &message.Content, &message.MediaURL,
			&message.Attachments, &message.MediaMimeType, &message.Transcription, &message.AIProcessed,
			&message.AIModel, &message.AIConfidenceScore, &message.RelatedTransactionID, &message.ScheduledTaskID,
			&message.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		messages = append(messages, &message)
	}

	return messages, nil
}

// FindSimpleByConversationID finds simplified messages with taxonomy labels
func (r *messageRepository) FindSimpleByConversationID(ctx context.Context, conversationID uuid.UUID, limit, offset int) ([]*entities.MessageSimple, error) {
	query := `
		SELECT 
			m.message_id,
			m.sender_type,
			COALESCE(st.label, 'Unknown') as sender_type_label,
			m.direction,
			m.message_type,
			COALESCE(mt.label, 'Unknown') as message_type_label,
			m.content,
			m.media_url,
			m.ai_processed,
			m.created_at
		FROM "vasst_expense".messages m
		LEFT JOIN "vasst_expense".taxonomy st ON m.sender_type = st.taxonomy_id
		LEFT JOIN "vasst_expense".taxonomy mt ON m.message_type = mt.taxonomy_id
		WHERE m.conversation_id = $1
		ORDER BY m.created_at ASC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.DB.QueryContext(ctx, query, conversationID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*entities.MessageSimple
	for rows.Next() {
		var message entities.MessageSimple
		err := rows.Scan(
			&message.MessageID, &message.SenderType, &message.SenderTypeLabel,
			&message.Direction, &message.MessageType, &message.MessageTypeLabel,
			&message.Content, &message.MediaURL, &message.AIProcessed, &message.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		messages = append(messages, &message)
	}

	return messages, nil
}

// CountByConversationID counts messages by conversation ID
func (r *messageRepository) CountByConversationID(ctx context.Context, conversationID uuid.UUID) (int64, error) {
	query := `
		SELECT COUNT(*)
		FROM "vasst_expense".messages 
		WHERE conversation_id = $1
	`

	var count int64
	err := r.DB.QueryRowContext(ctx, query, conversationID).Scan(&count)
	return count, err
}

// CountWithFilters counts messages with filtering
func (r *messageRepository) CountWithFilters(ctx context.Context, params *entities.MessageListParams) (int64, error) {
	query := `
		SELECT COUNT(*)
		FROM "vasst_expense".messages 
		WHERE 1=1
	`

	args := []interface{}{}
	argIndex := 1

	// Add same filters as FindWithFilters
	if params != nil {
		if params.ConversationID != nil {
			query += fmt.Sprintf(" AND conversation_id = $%d", argIndex)
			args = append(args, *params.ConversationID)
			argIndex++
		}
		if params.SenderType != nil {
			query += fmt.Sprintf(" AND sender_type = $%d", argIndex)
			args = append(args, *params.SenderType)
			argIndex++
		}
		if params.Direction != nil {
			query += fmt.Sprintf(" AND direction = $%d", argIndex)
			args = append(args, *params.Direction)
			argIndex++
		}
		if params.MessageType != nil {
			query += fmt.Sprintf(" AND message_type = $%d", argIndex)
			args = append(args, *params.MessageType)
			argIndex++
		}
		if params.AIProcessed != nil {
			query += fmt.Sprintf(" AND ai_processed = $%d", argIndex)
			args = append(args, *params.AIProcessed)
			argIndex++
		}
		if params.StartDate != nil {
			query += fmt.Sprintf(" AND created_at >= $%d", argIndex)
			args = append(args, *params.StartDate)
			argIndex++
		}
		if params.EndDate != nil {
			query += fmt.Sprintf(" AND created_at <= $%d", argIndex)
			args = append(args, *params.EndDate)
			argIndex++
		}
	}

	var count int64
	err := r.DB.QueryRowContext(ctx, query, args...).Scan(&count)
	return count, err
}

// MarkAsProcessed marks a message as AI processed
func (r *messageRepository) MarkAsProcessed(ctx context.Context, messageID uuid.UUID, aiModel string, confidenceScore *float64) error {
	query := `
		UPDATE "vasst_expense".messages 
		SET ai_processed = true, ai_model = $2, ai_confidence_score = $3
		WHERE message_id = $1
	`

	result, err := r.DB.ExecContext(ctx, query, messageID, aiModel, confidenceScore)
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
