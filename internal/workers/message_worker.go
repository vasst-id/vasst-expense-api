package workers

import (
	"context"
	"fmt"
	"sync"

	"github.com/vasst-id/vasst-expense-api/internal/entities"
	"github.com/vasst-id/vasst-expense-api/internal/events/handlers"
	"github.com/vasst-id/vasst-expense-api/internal/pubsub"
	logs "github.com/vasst-id/vasst-expense-api/internal/utils/logger"
)

type MessageWorker struct {
	pubsubClient   pubsub.Client
	messageHandler *handlers.MessageEventHandler
	logger         *logs.Logger
	ctx            context.Context
	cancel         context.CancelFunc
	wg             sync.WaitGroup
}

func NewMessageWorker(
	pubsubClient pubsub.Client,
	messageHandler *handlers.MessageEventHandler,
	logger *logs.Logger,
) *MessageWorker {
	ctx, cancel := context.WithCancel(context.Background())
	return &MessageWorker{
		pubsubClient:   pubsubClient,
		messageHandler: messageHandler,
		logger:         logger,
		ctx:            ctx,
		cancel:         cancel,
	}
}

func (w *MessageWorker) Start() error {
	w.logger.Info().Msg("Starting Message Worker")

	// Create subscription if not exists
	if err := w.pubsubClient.CreateSubscriptionIfNotExists(w.ctx, entities.SubscriptionWebhookProcessor, entities.TopicWebhookReceived); err != nil {
		return fmt.Errorf("failed to create webhook processor subscription: %w", err)
	}

	// Start webhook processor
	w.wg.Add(1)
	go func() {
		defer w.wg.Done()
		w.processWebhookEvents()
	}()

	w.logger.Info().Msg("Message Worker started successfully")
	return nil
}

func (w *MessageWorker) Stop() {
	w.logger.Info().Msg("Stopping Message Worker")
	w.cancel()
	w.wg.Wait()
	w.logger.Info().Msg("Message Worker stopped")
}

func (w *MessageWorker) processWebhookEvents() {
	w.logger.Info().Msg("Starting webhook event processor")

	err := w.pubsubClient.Subscribe(w.ctx, entities.SubscriptionWebhookProcessor, func(data []byte) error {
		if err := w.messageHandler.HandleWebhookReceived(w.ctx, data); err != nil {
			w.logger.Error().Err(err).Msg("Failed to process webhook event")
			return err
		}
		return nil
	})

	if err != nil {
		w.logger.Error().Err(err).Msg("Webhook event processor stopped with error")
	} else {
		w.logger.Info().Msg("Webhook event processor stopped")
	}
}
