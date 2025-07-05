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

type AIWorker struct {
	pubsubClient pubsub.Client
	aiHandler    *handlers.AIEventHandler
	logger       *logs.Logger
	ctx          context.Context
	cancel       context.CancelFunc
	wg           sync.WaitGroup
}

func NewAIWorker(
	pubsubClient pubsub.Client,
	aiHandler *handlers.AIEventHandler,
	logger *logs.Logger,
) *AIWorker {
	ctx, cancel := context.WithCancel(context.Background())
	return &AIWorker{
		pubsubClient: pubsubClient,
		aiHandler:    aiHandler,
		logger:       logger,
		ctx:          ctx,
		cancel:       cancel,
	}
}

func (w *AIWorker) Start() error {
	w.logger.Info().Msg("Starting AI Worker")

	// Create subscriptions if not exist
	if err := w.pubsubClient.CreateSubscriptionIfNotExists(w.ctx, entities.SubscriptionAIProcessor, entities.TopicMessageCreated); err != nil {
		return fmt.Errorf("failed to create AI processor subscription: %w", err)
	}

	if err := w.pubsubClient.CreateSubscriptionIfNotExists(w.ctx, entities.SubscriptionAIResponseProcessor, entities.TopicAIResponseReceived); err != nil {
		return fmt.Errorf("failed to create AI response processor subscription: %w", err)
	}

	// Start message processor for AI responses
	w.wg.Add(1)
	go func() {
		defer w.wg.Done()
		w.processMessageEvents()
	}()

	// Start AI response processor
	w.wg.Add(1)
	go func() {
		defer w.wg.Done()
		w.processAIResponseEvents()
	}()

	w.logger.Info().Msg("AI Worker started successfully")
	return nil
}

func (w *AIWorker) Stop() {
	w.logger.Info().Msg("Stopping AI Worker")
	w.cancel()
	w.wg.Wait()
	w.logger.Info().Msg("AI Worker stopped")
}

func (w *AIWorker) processMessageEvents() {
	w.logger.Info().Msg("Starting message event processor for AI")

	err := w.pubsubClient.Subscribe(w.ctx, entities.SubscriptionAIProcessor, func(data []byte) error {

		fmt.Println("\n\n\n data", string(data))

		if err := w.aiHandler.HandleMessageCreated(w.ctx, data); err != nil {
			w.logger.Error().Err(err).Msg("Failed to process message event for AI")
			return err
		}
		return nil
	})

	if err != nil {
		w.logger.Error().Err(err).Msg("Message event processor for AI stopped with error")
	} else {
		w.logger.Info().Msg("Message event processor for AI stopped")
	}
}

func (w *AIWorker) processAIResponseEvents() {
	w.logger.Info().Msg("Starting AI response event processor")

	err := w.pubsubClient.Subscribe(w.ctx, entities.SubscriptionAIResponseProcessor, func(data []byte) error {
		if err := w.aiHandler.HandleAIResponseReceived(w.ctx, data); err != nil {
			w.logger.Error().Err(err).Msg("Failed to process AI response event")
			return err
		}
		return nil
	})

	if err != nil {
		w.logger.Error().Err(err).Msg("AI response event processor stopped with error")
	} else {
		w.logger.Info().Msg("AI response event processor stopped")
	}
}
