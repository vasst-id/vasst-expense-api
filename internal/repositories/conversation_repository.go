package repositories

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/vasst-id/vasst-expense-api/internal/entities"
	"github.com/vasst-id/vasst-expense-api/internal/utils/postgres"
)

type (
	conversationRepository struct {
		*postgres.Postgres
	}

	ConversationRepository interface {
		Create(ctx context.Context, conversation *entities.Conversation) (entities.Conversation, error)
		Update(ctx context.Context, conversation *entities.Conversation) (entities.Conversation, error)
		Delete(ctx context.Context, conversationID uuid.UUID) error
		FindByID(ctx context.Context, conversationID uuid.UUID) (*entities.Conversation, error)
		FindByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entities.Conversation, error)
		FindActiveByUserID(ctx context.Context, userID uuid.UUID) ([]*entities.Conversation, error)
		FindByUserIDAndChannel(ctx context.Context, userID uuid.UUID, channel string) (*entities.Conversation, error)
		FindSimpleByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entities.ConversationSimple, error)
		CountByUserID(ctx context.Context, userID uuid.UUID) (int64, error)
	}
)

// NewConversationRepository creates a new ConversationRepository
func NewConversationRepository(pg *postgres.Postgres) ConversationRepository {
	return &conversationRepository{pg}
}

// Create creates a new conversation
func (r *conversationRepository) Create(ctx context.Context, conversation *entities.Conversation) (entities.Conversation, error) {
	query := `
		INSERT INTO "vasst_expense".conversations 
		(conversation_id, user_id, channel, is_active, context, metadata, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING conversation_id, user_id, channel, is_active, context, metadata, created_at, updated_at
	`

	var createdConversation entities.Conversation
	err := r.DB.QueryRowContext(ctx, query,
		conversation.ConversationID, conversation.UserID, conversation.Channel,
		conversation.IsActive, conversation.Context, conversation.Metadata,
	).Scan(
		&createdConversation.ConversationID, &createdConversation.UserID, &createdConversation.Channel,
		&createdConversation.IsActive, &createdConversation.Context, &createdConversation.Metadata,
		&createdConversation.CreatedAt, &createdConversation.UpdatedAt,
	)

	return createdConversation, err
}

// Update updates a conversation
func (r *conversationRepository) Update(ctx context.Context, conversation *entities.Conversation) (entities.Conversation, error) {
	query := `
		UPDATE "vasst_expense".conversations 
		SET channel = $2, is_active = $3, context = $4, metadata = $5, updated_at = CURRENT_TIMESTAMP
		WHERE conversation_id = $1
		RETURNING conversation_id, user_id, channel, is_active, context, metadata, created_at, updated_at
	`

	var updatedConversation entities.Conversation
	err := r.DB.QueryRowContext(ctx, query,
		conversation.ConversationID, conversation.Channel, conversation.IsActive,
		conversation.Context, conversation.Metadata,
	).Scan(
		&updatedConversation.ConversationID, &updatedConversation.UserID, &updatedConversation.Channel,
		&updatedConversation.IsActive, &updatedConversation.Context, &updatedConversation.Metadata,
		&updatedConversation.CreatedAt, &updatedConversation.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return entities.Conversation{}, sql.ErrNoRows // Conversation not found
		}
		return entities.Conversation{}, err
	}

	return updatedConversation, nil
}

// Delete soft deletes a conversation (sets is_active to false)
func (r *conversationRepository) Delete(ctx context.Context, conversationID uuid.UUID) error {
	query := `
		UPDATE "vasst_expense".conversations 
		SET is_active = false, updated_at = CURRENT_TIMESTAMP
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

// FindByID finds a conversation by ID
func (r *conversationRepository) FindByID(ctx context.Context, conversationID uuid.UUID) (*entities.Conversation, error) {
	query := `
		SELECT conversation_id, user_id, channel, is_active, context, metadata, created_at, updated_at
		FROM "vasst_expense".conversations 
		WHERE conversation_id = $1
	`

	var conversation entities.Conversation
	err := r.DB.QueryRowContext(ctx, query, conversationID).Scan(
		&conversation.ConversationID, &conversation.UserID, &conversation.Channel,
		&conversation.IsActive, &conversation.Context, &conversation.Metadata,
		&conversation.CreatedAt, &conversation.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &conversation, nil
}

// FindByUserID finds conversations by user ID with pagination
func (r *conversationRepository) FindByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entities.Conversation, error) {
	query := `
		SELECT conversation_id, user_id, channel, is_active, context, metadata, created_at, updated_at
		FROM "vasst_expense".conversations 
		WHERE user_id = $1
		ORDER BY updated_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.DB.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var conversations []*entities.Conversation
	for rows.Next() {
		var conversation entities.Conversation
		err := rows.Scan(
			&conversation.ConversationID, &conversation.UserID, &conversation.Channel,
			&conversation.IsActive, &conversation.Context, &conversation.Metadata,
			&conversation.CreatedAt, &conversation.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		conversations = append(conversations, &conversation)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return conversations, nil
}

// FindActiveByUserID finds all active conversations by user ID
func (r *conversationRepository) FindActiveByUserID(ctx context.Context, userID uuid.UUID) ([]*entities.Conversation, error) {
	query := `
		SELECT conversation_id, user_id, channel, is_active, context, metadata, created_at, updated_at
		FROM "vasst_expense".conversations 
		WHERE user_id = $1 AND is_active = true
		ORDER BY updated_at DESC
	`

	rows, err := r.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var conversations []*entities.Conversation
	for rows.Next() {
		var conversation entities.Conversation
		err := rows.Scan(
			&conversation.ConversationID, &conversation.UserID, &conversation.Channel,
			&conversation.IsActive, &conversation.Context, &conversation.Metadata,
			&conversation.CreatedAt, &conversation.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		conversations = append(conversations, &conversation)
	}

	return conversations, nil
}

// FindByUserIDAndChannel finds a conversation by user ID and channel
func (r *conversationRepository) FindByUserIDAndChannel(ctx context.Context, userID uuid.UUID, channel string) (*entities.Conversation, error) {
	query := `
		SELECT conversation_id, user_id, channel, is_active, context, metadata, created_at, updated_at
		FROM "vasst_expense".conversations 
		WHERE user_id = $1 AND channel = $2 AND is_active = true
		ORDER BY updated_at DESC
		LIMIT 1
	`

	var conversation entities.Conversation
	err := r.DB.QueryRowContext(ctx, query, userID, channel).Scan(
		&conversation.ConversationID, &conversation.UserID, &conversation.Channel,
		&conversation.IsActive, &conversation.Context, &conversation.Metadata,
		&conversation.CreatedAt, &conversation.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &conversation, nil
}

// FindSimpleByUserID finds simplified conversations with message stats by user ID
func (r *conversationRepository) FindSimpleByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entities.ConversationSimple, error) {
	query := `
		SELECT 
			c.conversation_id,
			c.channel,
			c.is_active,
			COALESCE(m.last_content, '') as last_message,
			COALESCE(m.message_count, 0) as message_count,
			c.updated_at
		FROM "vasst_expense".conversations c
		LEFT JOIN (
			SELECT 
				conversation_id,
				MAX(content) as last_content,
				COUNT(*) as message_count
			FROM "vasst_expense".messages 
			GROUP BY conversation_id
		) m ON c.conversation_id = m.conversation_id
		WHERE c.user_id = $1
		ORDER BY c.updated_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.DB.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var conversations []*entities.ConversationSimple
	for rows.Next() {
		var conversation entities.ConversationSimple
		err := rows.Scan(
			&conversation.ConversationID, &conversation.Channel, &conversation.IsActive,
			&conversation.LastMessage, &conversation.MessageCount, &conversation.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		conversations = append(conversations, &conversation)
	}

	return conversations, nil
}

// CountByUserID counts conversations by user ID
func (r *conversationRepository) CountByUserID(ctx context.Context, userID uuid.UUID) (int64, error) {
	query := `
		SELECT COUNT(*)
		FROM "vasst_expense".conversations 
		WHERE user_id = $1
	`

	var count int64
	err := r.DB.QueryRowContext(ctx, query, userID).Scan(&count)
	return count, err
}
