package subscriber

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/vasst-id/vasst-expense-api/config"
	"github.com/vasst-id/vasst-expense-api/internal/pubsub"
	"github.com/vasst-id/vasst-expense-api/internal/utils"
	logs "github.com/vasst-id/vasst-expense-api/internal/utils/logger"
	"github.com/vasst-id/vasst-expense-api/internal/utils/postgres"
)

const (
	ServiceName         = "vasst-expense-api-subscriber"
	ServiceNamePostgres = ServiceName + "-" + "postgres"
	SentryDSN           = "https://b772fc0746f1e67d2ecac68f0c3d41bd@o4509568572653568.ingest.us.sentry.io/4509568576126976"
)

func Run(config *config.Config) {
	// sentry
	sentryOpts := sentry.ClientOptions{
		Dsn:              SentryDSN,
		Environment:      config.Env,
		TracesSampleRate: 1.0,
	}

	err := sentry.Init(sentryOpts)
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}

	// local postgres
	pgCfg := &postgres.Config{
		ServiceName: ServiceNamePostgres,
		Dsn:         config.PostgreHost,
		MaxConn:     config.DBMaxOpenConnection,
		MaxIdle:     config.DBMaxIdleConnection,
	}

	pg, err := postgres.New(pgCfg)
	if err != nil {
		log.Fatalf("error init postgres %s", err.Error())
	}

	// logger
	logger := logs.New(initLoggerOptions(config))
	if logger == nil {
		log.Fatalf("error init logger")
	}

	// set logger singleton
	utils.SetLogger(logger)

	// Initialize Google Pub/Sub client
	pubsubClient, err := pubsub.NewGooglePubSubClientWithCredentials(config.GoogleCloudProjectID, config.GoogleCloudCredentialsFile)
	if err != nil {
		log.Fatalf("error init pubsub client %s", err.Error())
	}
	defer pubsubClient.Close()

	// Create topics if they don't exist
	// ctx := context.Background()
	// topics := []string{
	// 	entities.TopicWebhookReceived,
	// 	entities.TopicMessageCreated,
	// 	entities.TopicAIResponseReceived,
	// 	entities.TopicMessageDelivery,
	// }

	// for _, topic := range topics {
	// 	if err := pubsubClient.CreateTopicIfNotExists(ctx, topic); err != nil {
	// 		log.Fatalf("error creating topic %s: %s", topic, err.Error())
	// 	}
	// }

	// Initialize Google Cloud Storage service
	// storageService, err := services.NewGoogleStorageService(
	// 	config.GoogleCloudProjectID,
	// 	config.GoogleCloudBucketPrefix,
	// 	config.GoogleCloudCredentialsFile,
	// 	config.GoogleCloudRegion,
	// )
	// if err != nil {
	// 	log.Fatalf("error init google cloud storage service %s", err.Error())
	// }

	// Initialize services
	// messageService := services.NewMessageService(repositories.NewMessageRepository(pg), storageService)

	// geminiService, err := services.NewGeminiService(config, messageService)
	// if err != nil {
	// 	log.Fatalf("error init gemini service %s", err.Error())
	// }

	// openAIService, err := services.NewOpenAIService(config, messageService)
	// if err != nil {
	// 	log.Fatalf("error init openai service %s", err.Error())
	// }

	// userService := services.NewUserService(repositories.NewUserRepository(pg), nil) // TODO: Add auth middleware if needed
	// whatsAppMediaService := services.NewWhatsAppMediaService(config)
	// whatsAppService, err := services.NewWhatsAppService(config, userService, messageService, openAIService, geminiService, storageService)
	// if err != nil {
	// 	log.Fatalf("error init whatsapp service %s", err.Error())
	// }

	// // Initialize event handlers
	// messageEventHandler := handlers.NewMessageEventHandler(pubsubClient, messageService, userService, whatsAppMediaService, whatsAppService, storageService, organizationService, config, logger)
	// aiEventHandler := handlers.NewAIEventHandler(pubsubClient, geminiService, messageService, whatsAppService)

	// // Initialize workers
	// messageWorker := workers.NewMessageWorker(pubsubClient, messageEventHandler)
	// aiWorker := workers.NewAIWorker(pubsubClient, aiEventHandler)

	// Start workers
	logger.Info().Msg("Starting event-driven workers...")

	// if err := messageWorker.Start(); err != nil {
	// 	log.Fatalf("error starting message worker: %s", err.Error())
	// }

	// if err := aiWorker.Start(); err != nil {
	// 	log.Fatalf("error starting AI worker: %s", err.Error())
	// }

	// if err := messageDeliveryWorker.Start(); err != nil {
	// 	log.Fatalf("error starting message delivery worker: %s", err.Error())
	// }

	// logger.Info().Msg("All workers started successfully")

	// Wait for shutdown signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	logger.Info().Msg("Shutdown signal received, stopping workers...")

	// Stop workers gracefully
	var wg sync.WaitGroup
	wg.Add(3)

	// go func() {
	// 	defer wg.Done()
	// 	messageWorker.Stop()
	// }()

	// go func() {
	// 	defer wg.Done()
	// 	aiWorker.Stop()
	// }()

	// go func() {
	// 	defer wg.Done()
	// 	messageDeliveryWorker.Stop()
	// }()

	// Wait for workers to stop with timeout
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		logger.Info().Msg("All workers stopped successfully")
	case <-time.After(30 * time.Second):
		logger.Warn().Msg("Timeout waiting for workers to stop")
	}

	logger.Info().Msg("Subscriber service stopped")
}

func initLoggerOptions(cfg *config.Config) logs.Options {
	return logs.Options{
		FileDirectory:   cfg.LogDirectory,
		FileName:        cfg.LogFileName,
		MaxSize:         cfg.LogMaxSize,
		MaxAge:          cfg.LogMaxAge,
		MaxBackups:      cfg.LogMaxBackups,
		ConsoleLog:      cfg.LogConsole,
		OutputLevel:     0, // Default to info level
		DisableCompress: false,
	}
}
