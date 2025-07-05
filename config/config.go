package config

import (
	"time"

	"github.com/spf13/viper"
)

type (
	Config struct {
		App        string `mapstructure:"APP"`
		Env        string `mapstructure:"ENV"`
		SvcVersion string `mapstructure:"SVC_VERSION"`
		Port       string `mapstructure:"PORT"`
		JWTSecret  string `mapstructure:"JWT_SECRET"`

		// Postgres
		PostgreHost         string `mapstructure:"POSTGRES_URL"`
		DBMaxOpenConnection int    `mapstructure:"DB_MAX_OPEN_CONN"`
		DBMaxIdleConnection int    `mapstructure:"DB_MAX_IDLE_CONN"`

		// Redis
		RedisDB       int      `mapstructure:"REDIS_DB"`
		RedisHost     []string `mapstructure:"REDIS_URL"`
		RedisUsername string   `mapstructure:"REDIS_USERNAME"`
		RedisPassword string   `mapstructure:"REDIS_PASSWORD"`

		// HTTP client
		HttpClientTimeout             int  `mapstructure:"HTTP_CLIENT_TIMEOUT"`
		HttpClientDisableKeepAlives   bool `mapstructure:"HTTP_CLIENT_DISABLE_KEEP_ALIVE"`
		HttpClientMaxIdleConns        int  `mapstructure:"HTTP_CLIENT_MAX_IDLE_CONNS"`
		HttpClientMaxConnsPerHost     int  `mapstructure:"HTTP_CLIENT_MAX_CONNS_PER_HOST"`
		HttpClientMaxIdleConnsPerHost int  `mapstructure:"HTTP_CLIENT_MAX_IDLE_CONNS_PER_HOST"`
		HttpClientIdleConnTimeout     int  `mapstructure:"HTTP_CLIENT_IDLE_CONN_TIMEOUT"`

		// Logging
		LogDirectory  string `mapstructure:"LOG_DIR"`
		LogFileName   string `mapstructure:"LOG_FILENAME"`
		LogConsole    bool   `mapstructure:"LOG_CONSOLE"`
		LogMaxSize    int    `mapstructure:"LOG_MAX_SIZE"`
		LogMaxAge     int    `mapstructure:"LOG_MAX_AGE"`
		LogMaxBackups int    `mapstructure:"LOG_MAX_BACKUP"`

		// Datadog
		DatadogAgentHost string `mapstructure:"DATADOG_AGENT_HOST"`

		// WhatsApp
		WhatsAppPhoneNumberID string `mapstructure:"WHATSAPP_PHONE_NUMBER_ID"`
		WhatsAppAccessToken   string `mapstructure:"WHATSAPP_ACCESS_TOKEN"`

		// OpenAI
		OpenAIApiKey string `mapstructure:"OPENAI_API_KEY"`

		// Gemini
		GeminiApiKey string `mapstructure:"GEMINI_API_KEY"`

		// Google Cloud Storage
		GoogleCloudProjectID       string `mapstructure:"GOOGLE_CLOUD_PROJECT_ID"`
		GoogleCloudBucketPrefix    string `mapstructure:"GOOGLE_CLOUD_BUCKET_PREFIX"`
		GoogleCloudCredentialsFile string `mapstructure:"GOOGLE_CLOUD_CREDENTIALS_FILE"`
		GoogleCloudRegion          string `mapstructure:"GOOGLE_CLOUD_REGION"`

		// Message Processing & AI Features
		EnableTypingIndicators bool          `mapstructure:"ENABLE_TYPING_INDICATORS"`
		EnableMultiMessage     bool          `mapstructure:"ENABLE_MULTI_MESSAGE"`
		MessageDelay           time.Duration `mapstructure:"MESSAGE_DELAY"`
		MaxMessagesPerResponse int           `mapstructure:"MAX_MESSAGES_PER_RESPONSE"`
		MaxMessageWordCount    int           `mapstructure:"MAX_MESSAGE_WORD_COUNT"`
		MessageChunkLength     int           `mapstructure:"MESSAGE_CHUNK_LENGTH"`
	}
)

func NewConfig() (*Config, error) {
	viper.SetConfigFile(".env")

	// Set defaults for new AI and messaging features
	viper.SetDefault("ENABLE_TYPING_INDICATORS", true)
	viper.SetDefault("ENABLE_MULTI_MESSAGE", true)
	viper.SetDefault("MESSAGE_DELAY", "1s")
	viper.SetDefault("MAX_MESSAGES_PER_RESPONSE", 5)
	viper.SetDefault("MAX_MESSAGE_WORD_COUNT", 500)
	viper.SetDefault("MESSAGE_CHUNK_LENGTH", 1000)

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = viper.Unmarshal(&cfg)
	redisHost := viper.GetStringSlice("REDIS_URL")
	cfg.RedisHost = redisHost

	// Parse MESSAGE_DELAY duration
	if delayStr := viper.GetString("MESSAGE_DELAY"); delayStr != "" {
		if duration, parseErr := time.ParseDuration(delayStr); parseErr == nil {
			cfg.MessageDelay = duration
		} else {
			cfg.MessageDelay = 3 * time.Second // Default fallback
		}
	}

	return &cfg, err
}
