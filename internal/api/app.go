package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	// "github.com/getsentry/sentry-go"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	swaggofiles "github.com/swaggo/files"
	ginswagger "github.com/swaggo/gin-swagger"
	gintrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/gin-gonic/gin"

	grace "github.com/julofinance/grace/v2"
	"github.com/vasst-id/vasst-expense-api/internal/repositories"
	health "github.com/vasst-id/vasst-expense-api/internal/utils/healthcheck"
	"github.com/vasst-id/vasst-expense-api/internal/utils/httpclient"
	logs "github.com/vasst-id/vasst-expense-api/internal/utils/logger"
	"github.com/vasst-id/vasst-expense-api/internal/utils/postgres"

	"github.com/getsentry/sentry-go"
	"github.com/vasst-id/vasst-expense-api/config"
	httpRouter "github.com/vasst-id/vasst-expense-api/internal/controller/http/v1"
	"github.com/vasst-id/vasst-expense-api/internal/middleware"
	"github.com/vasst-id/vasst-expense-api/internal/services"
	"github.com/vasst-id/vasst-expense-api/internal/utils"
)

const (
	ServiceName           = "vasst-expense-api-api"
	ServiceNamePostgres   = ServiceName + "-" + "postgres"
	ServiceNameRedis      = ServiceName + "-" + "redis"
	ServiceNameThirdParty = ServiceName + "-" + "third-party-app"

	SentryDSN = "https://b772fc0746f1e67d2ecac68f0c3d41bd@o4509568572653568.ingest.us.sentry.io/4509568576126976"
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
		// DataDogTracer: false,
	}

	pg, err := postgres.New(pgCfg)
	if err != nil {
		log.Fatalf("error init postgres %s", err.Error())
	}

	// redis
	// redisOpts := &redis.ClientOptions{
	// 	ServiceName:   ServiceNameRedis,
	// 	Address:       config.RedisHost,
	// 	Username:      config.RedisUsername,
	// 	Password:      config.RedisPassword,
	// 	DB:            config.RedisDB,
	// 	DataDogTracer: true,
	// }

	// reds, err := redis.New(redisOpts)
	// if err != nil {
	// 	log.Fatalf("error init redis %s", err.Error())
	// }

	// logger
	logger := logs.New(initLoggerOptions(config))
	if logger == nil {
		log.Fatalf("error init logger")
	}

	// set logger singleton
	utils.SetLogger(logger)

	// Initialize Google Pub/Sub client for webhook events
	// pubsubClient, err := pubsub.NewGooglePubSubClientWithCredentials(config.GoogleCloudProjectID, config.GoogleCloudCredentialsFile)
	// if err != nil {
	// 	log.Fatalf("error init pubsub client %s", err.Error())
	// }
	// defer pubsubClient.Close()

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

	// services
	authMiddleware := middleware.NewAuthMiddleware(config.JWTSecret)
	userService := services.NewUserService(repositories.NewUserRepository(pg), authMiddleware)
	workspaceService := services.NewWorkspaceService(repositories.NewWorkspaceRepository(pg))
	accountService := services.NewAccountService(repositories.NewAccountRepository(pg))
	bankService := services.NewBankService(repositories.NewBankRepository(pg))
	currencyService := services.NewCurrencyService(repositories.NewCurrencyRepository(pg))
	subscriptionPlanService := services.NewSubscriptionPlanService(repositories.NewSubscriptionPlanRepository(pg))
	budgetService := services.NewBudgetService(repositories.NewBudgetRepository(pg))
	categoryService := services.NewCategoryService(repositories.NewCategoryRepository(pg))
	transactionService := services.NewTransactionService(repositories.NewTransactionRepository(pg), repositories.NewWorkspaceRepository(pg), repositories.NewAccountRepository(pg))
	conversationService := services.NewConversationService(repositories.NewConversationRepository(pg), repositories.NewUserRepository(pg))
	messageService := services.NewMessageService(repositories.NewMessageRepository(pg), repositories.NewConversationRepository(pg), repositories.NewUserRepository(pg))
	taxonomyService := services.NewTaxonomyService(repositories.NewTaxonomyRepository(pg))
	userTagsService := services.NewUserTagsService(repositories.NewUserTagsRepository(pg))
	transactionTagsService := services.NewTransactionTagsService(repositories.NewTransactionTagsRepository(pg), repositories.NewUserTagsRepository(pg))
	verificationCodeService := services.NewVerificationCodeService(repositories.NewVerificationCodeRepository(pg), repositories.NewUserRepository(pg))
	// openAIService, err := services.NewOpenAIService(config, messageService)
	// if err != nil {
	// 	log.Fatalf("error init openai service %s", err.Error())
	// }
	// geminiService, err := services.NewGeminiService(config, messageService)
	// if err != nil {
	// 	log.Fatalf("error init gemini service %s", err.Error())
	// }
	// whatsAppService, err := services.NewWhatsAppService(config, userService, messageService)
	// if err != nil {
	// 	log.Fatalf("error init whatsapp service %s", err.Error())
	// }

	// Initialize unified webhook event handler
	// publisherAdapter := pubsub.NewPublisherAdapter(pubsubClient)
	// webhookEventHandler := handlers.NewWebhookEventHandler(publisherAdapter)

	httpClient := httpclient.New(httpClientConfig(config))

	// gin
	gin.SetMode(gin.ReleaseMode)
	handler := gin.New()

	// middlewares
	handler.Use(gintrace.Middleware(ServiceName))
	handler.Use(gin.Logger())
	handler.Use(gzip.Gzip(gzip.DefaultCompression))
	handler.Use(gin.Recovery())
	handler.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:3000", "http://localhost:3001"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// swagger
	handler.GET("/swagger/*any", ginswagger.WrapHandler(swaggofiles.Handler))

	// health check
	healthCheck := health.New(
		health.WithDB(pg.DB, health.Config{Name: ServiceNamePostgres}),
		// health.WithRedis(reds, health.Config{Name: ServiceNameRedis}),
		health.WithLogger(logger),
		health.WithComponent(health.Component{
			Name: ServiceNameThirdParty,
			CheckFunc: func(ctx context.Context) error {
				req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://httpstat.us/200", nil)
				if err != nil {
					return err
				}

				resp, err := httpClient.Do(req)
				if err != nil {
					return err
				}
				defer resp.Body.Close()

				utils.Log().Info().Msg(resp.Status)

				return err
			},
		}),
	)

	handler.GET("/health-check", gin.WrapF(healthCheck.HandlerFunc))

	httpRouter.NewRouter(handler, httpRouter.Services{
		Cfg:                     config,
		UserService:             userService,
		WorkspaceService:        workspaceService,
		AccountService:          accountService,
		BankService:             bankService,
		CurrencyService:         currencyService,
		SubscriptionPlanService: subscriptionPlanService,
		BudgetService:           budgetService,
		AuthMiddleware:          authMiddleware,
		CategoryService:         categoryService,
		TransactionService:      transactionService,
		ConversationService:     conversationService,
		MessageService:          messageService,
		TaxonomyService:         taxonomyService,
		UserTagsService:         userTagsService,
		TransactionTagsService:  transactionTagsService,
		VerificationCodeService: verificationCodeService,
	})

	fmt.Printf("Starting server on port %s\n", config.Port)

	grace.Serve(config.Port, handler)

	fmt.Println("Server started successfully")

}

func initLoggerOptions(cfg *config.Config) logs.Options {
	return logs.Options{
		ConsoleLog:    cfg.LogConsole,
		FileDirectory: cfg.LogDirectory,
		FileName:      cfg.LogFileName,
		MaxSize:       cfg.LogMaxSize,
		MaxAge:        cfg.LogMaxAge,
		MaxBackups:    cfg.LogMaxBackups,
	}
}

func httpClientConfig(config *config.Config) *httpclient.Config {
	httpClientCfg := &httpclient.Config{
		Timeout:     config.HttpClientTimeout,
		ServiceName: ServiceNameThirdParty,
		Transport: struct {
			DisableKeepAlives   bool
			MaxIdleConns        int
			MaxConnsPerHost     int
			MaxIdleConnsPerHost int
			IdleConnTimeout     time.Duration
		}{
			DisableKeepAlives:   config.HttpClientDisableKeepAlives,
			MaxIdleConns:        config.HttpClientMaxIdleConns,
			MaxConnsPerHost:     config.HttpClientMaxConnsPerHost,
			MaxIdleConnsPerHost: config.HttpClientMaxIdleConnsPerHost,
			IdleConnTimeout:     time.Duration(config.HttpClientIdleConnTimeout) * time.Second,
		},
		// DataDogTracer: false,
	}

	return httpClientCfg
}
