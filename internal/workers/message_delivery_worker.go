package workers

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/vasst-id/vasst-expense-api/config"
	"github.com/vasst-id/vasst-expense-api/internal/entities"
	"github.com/vasst-id/vasst-expense-api/internal/pubsub"
	"github.com/vasst-id/vasst-expense-api/internal/services"
	"github.com/vasst-id/vasst-expense-api/internal/utils"
	logs "github.com/vasst-id/vasst-expense-api/internal/utils/logger"
)

type MessageDeliveryWorker struct {
	pubsubClient        pubsub.Client
	messageService      services.MessageService
	conversationService services.ConversationService
	contactService      services.ContactService
	whatsappService     services.WhatsAppService
	config              *config.Config
	logger              *logs.Logger
	ctx                 context.Context
	cancel              context.CancelFunc
	wg                  sync.WaitGroup
}

func NewMessageDeliveryWorker(
	pubsubClient pubsub.Client,
	messageService services.MessageService,
	conversationService services.ConversationService,
	contactService services.ContactService,
	whatsappService services.WhatsAppService,
	config *config.Config,
	logger *logs.Logger,
) *MessageDeliveryWorker {
	ctx, cancel := context.WithCancel(context.Background())
	return &MessageDeliveryWorker{
		pubsubClient:        pubsubClient,
		messageService:      messageService,
		conversationService: conversationService,
		contactService:      contactService,
		whatsappService:     whatsappService,
		config:              config,
		logger:              logger,
		ctx:                 ctx,
		cancel:              cancel,
	}
}

func (w *MessageDeliveryWorker) Start() error {
	w.logger.Info().Msg("Starting Message Delivery Worker")

	// Create subscription if not exist
	if err := w.pubsubClient.CreateSubscriptionIfNotExists(w.ctx, entities.SubscriptionMessageDelivery, entities.TopicMessageDelivery); err != nil {
		return fmt.Errorf("failed to create message delivery subscription: %w", err)
	}

	// Start message delivery processor
	w.wg.Add(1)
	go func() {
		defer w.wg.Done()
		if err := w.pubsubClient.Subscribe(w.ctx, entities.SubscriptionMessageDelivery, func(data []byte) error {
			return w.HandleMessageDelivery(w.ctx, data)
		}); err != nil {
			w.logger.Error().Err(err).Msg("Message delivery subscription failed")
		}
	}()

	w.logger.Info().Msg("Message Delivery Worker started successfully")
	return nil
}

func (w *MessageDeliveryWorker) Stop() {
	w.logger.Info().Msg("Stopping Message Delivery Worker")
	w.cancel()
	w.wg.Wait()
	w.logger.Info().Msg("Message Delivery Worker stopped")
}

// HandleMessageDelivery processes message delivery events and sends messages via appropriate channels
func (w *MessageDeliveryWorker) HandleMessageDelivery(ctx context.Context, data []byte) error {
	var deliveryEvent entities.MessageDeliveryEvent
	if err := json.Unmarshal(data, &deliveryEvent); err != nil {
		w.logger.Error().Err(err).Msg("Failed to unmarshal message delivery event")
		return fmt.Errorf("failed to unmarshal message delivery event: %w", err)
	}

	w.logger.Info().
		Str("event_id", deliveryEvent.EventID.String()).
		Str("message_id", deliveryEvent.MessageID.String()).
		Str("conversation_id", deliveryEvent.ConversationID.String()).
		Str("medium", deliveryEvent.Medium).
		Msg("Processing message delivery")

	// Get message details
	message, err := w.messageService.GetMessageByID(ctx, deliveryEvent.MessageID, deliveryEvent.OrganizationID)
	if err != nil {
		w.logger.Error().Err(err).Str("message_id", deliveryEvent.MessageID.String()).Msg("Failed to get message for delivery")
		return fmt.Errorf("failed to get message for delivery: %w", err)
	}

	// Get contact for phone number info with tenant isolation
	contact, err := w.contactService.GetContactByIDAndOrganization(ctx, deliveryEvent.ContactID, deliveryEvent.OrganizationID)
	if err != nil {
		w.logger.Error().Err(err).Str("contact_id", deliveryEvent.ContactID.String()).Msg("Failed to get contact for delivery")
		return fmt.Errorf("failed to get contact for delivery: %w", err)
	}

	// Send via appropriate channel
	var deliveryErr error
	switch deliveryEvent.Medium {
	case "whatsapp":
		deliveryErr = w.sendWhatsAppMessage(ctx, contact, message)
	case "email":
		// TODO: Implement email delivery
		w.logger.Warn().Str("medium", deliveryEvent.Medium).Msg("Email delivery not implemented yet")
		deliveryErr = fmt.Errorf("email delivery not implemented")
	case "sms":
		// TODO: Implement SMS delivery
		w.logger.Warn().Str("medium", deliveryEvent.Medium).Msg("SMS delivery not implemented yet")
		deliveryErr = fmt.Errorf("SMS delivery not implemented")
	default:
		deliveryErr = fmt.Errorf("unsupported delivery medium: %s", deliveryEvent.Medium)
	}

	// Update message status based on delivery result
	if deliveryErr != nil {
		w.logger.Error().Err(deliveryErr).
			Str("message_id", message.MessageID.String()).
			Str("medium", deliveryEvent.Medium).
			Msg("Failed to deliver message")

		// Update message status to failed
		failureReason := deliveryErr.Error()
		updateInput := &entities.UpdateMessageStatusInput{
			Status:        int(entities.MessageStatusFailed),
			FailureReason: &failureReason,
		}
		if updateErr := w.messageService.UpdateMessageStatus(ctx, message.MessageID, message.OrganizationID, updateInput); updateErr != nil {
			w.logger.Error().Err(updateErr).Str("message_id", message.MessageID.String()).Msg("Failed to update message status to failed")
		}
		return deliveryErr
	}

	// Update message status to sent
	updateInput := &entities.UpdateMessageStatusInput{
		Status: int(entities.MessageStatusSent),
	}
	if updateErr := w.messageService.UpdateMessageStatus(ctx, message.MessageID, message.OrganizationID, updateInput); updateErr != nil {
		w.logger.Error().Err(updateErr).Str("message_id", message.MessageID.String()).Msg("Failed to update message status to sent")
		return fmt.Errorf("failed to update message status: %w", updateErr)
	}

	w.logger.Info().
		Str("message_id", message.MessageID.String()).
		Str("medium", deliveryEvent.Medium).
		Msg("Message delivered successfully")

	return nil
}

// sendWhatsAppMessage handles WhatsApp message delivery with chunking
func (w *MessageDeliveryWorker) sendWhatsAppMessage(ctx context.Context, contact *entities.Contact, message *entities.Message) error {
	// Check if features are enabled
	if !w.config.EnableMultiMessage {
		// Send single message without chunking
		return w.sendSingleMessage(ctx, contact.PhoneNumber, message.Content)
	}

	// Chunk message if multi-message is enabled
	chunks := utils.ChunkMessage(message.Content, w.config.MessageChunkLength)
	
	// Limit chunks to configured maximum
	if len(chunks) > w.config.MaxMessagesPerResponse {
		chunks = chunks[:w.config.MaxMessagesPerResponse]
		w.logger.Warn().
			Str("message_id", message.MessageID.String()).
			Int("original_chunks", len(chunks)).
			Int("limited_to", w.config.MaxMessagesPerResponse).
			Msg("Message chunks limited to configured maximum")
	}

	// If only one chunk, send as single message
	if len(chunks) <= 1 {
		content := message.Content
		if len(chunks) == 1 {
			content = chunks[0]
		}
		return w.sendSingleMessage(ctx, contact.PhoneNumber, content)
	}

	// Send multiple chunks with delays (typing indicator already sent before AI processing)
	for i, chunk := range chunks {
		// Add delay before sending (for all chunks to simulate natural typing)
		if w.config.MessageDelay > 0 {
			select {
			case <-time.After(w.config.MessageDelay):
			case <-ctx.Done():
				return ctx.Err()
			}
		}

		// Send the chunk
		if err := w.whatsappService.SendTextMessage(ctx, contact.PhoneNumber, chunk); err != nil {
			w.logger.Error().Err(err).
				Str("message_id", message.MessageID.String()).
				Str("phone_number", contact.PhoneNumber).
				Int("chunk_index", i).
				Msg("Failed to send message chunk")
			return fmt.Errorf("failed to send message chunk %d: %w", i+1, err)
		}

		w.logger.Info().
			Str("message_id", message.MessageID.String()).
			Str("phone_number", contact.PhoneNumber).
			Int("chunk_index", i+1).
			Int("total_chunks", len(chunks)).
			Msg("Message chunk sent successfully")
	}

	w.logger.Info().
		Str("message_id", message.MessageID.String()).
		Str("phone_number", contact.PhoneNumber).
		Int("total_chunks", len(chunks)).
		Msg("All message chunks sent successfully")

	return nil
}

// sendSingleMessage sends a single message without chunking
func (w *MessageDeliveryWorker) sendSingleMessage(ctx context.Context, phoneNumber, content string) error {
	// Send the message (typing indicator already sent before AI processing)
	if err := w.whatsappService.SendTextMessage(ctx, phoneNumber, content); err != nil {
		return fmt.Errorf("failed to send WhatsApp message: %w", err)
	}

	return nil
}