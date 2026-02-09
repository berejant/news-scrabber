package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfigFromEnvVars(t *testing.T) {
	expectedConfig := Config{
		App:         "social-posting-service",
		Environment: "test",

		OauthConfig: OauthConfig{
			StateSecret: "secret",
			StateExpire: time.Hour * 3,
		},

		Storage: StorageConfig{
			Driver: "badger",
			URI:    "/tmp/badger",
		},

		DB: DBConfig{
			Driver: "mysql",
			DSN:    "user:password@tcp(localhost:3306)/dbname",
		},

		Server: ServerConfig{
			Host:            "localhost",
			Port:            8080,
			PublicURI:       "http://localhost:8080",
			OmniAPIToken:    "omni-token",
			CwAPIToken:      "cw-token",
			SwaggerPassword: "swagger-password",
		},

		Kafka: KafkaConfig{
			Brokers: []string{"localhost:9092"},
			SASLConfig: SASLConfig{
				Mechanism: ScramSHA256,
				Username:  "test-user",
				Password:  "test-password",
			},
			TLSConfig:   TLSConfig{},
			MaxAttempts: 7,
			DialTimeout: time.Second * 3,

			TopicTransactions:          "order-created",
			TopicContentUpdated:        "fnd.contentUpdated",
			ConsumersCountTransactions: 2,
		},

		AMQP: AMQPConfig{
			URL:           "amqp://guest:guest@localhost:5672/",
			PrefetchCount: 1,
			PrefetchSize:  1,
		},

		X: XConfig{
			BaseURI:        "http://localhost:8080",
			ClientID:       "x-client-id",
			ClientSecret:   "x-client-secret",
			ConsumerKey:    "x-consumer-key",
			ConsumerSecret: "x-consumer-secret",
		},

		Bluesky: BlueskyConfig{
			BaseURI:      "https://bsky.example.com",
			ClientID:     "bluesky-client-id",
			ClientSecret: "bluesky-client-secret",
		},

		RectalogistConfig: RectalogistConfig{
			URI:    "http://rectalogist.example.com",
			CdnURI: "https://test-cdn.com",
		},

		ContentDataMartConfig: ContentDataMartConfig{
			URI: "https://content-data-mart.example.com",
			Key: "content-data-mart.key",
		},

		MaxPublishRetryAttempts: 5,
	}

	getEnvValues := func() map[string]string {
		return map[string]string{
			"APP":         expectedConfig.App,
			"ENVIRONMENT": expectedConfig.Environment,

			"OAUTH_STATE_SECRET": expectedConfig.OauthConfig.StateSecret,
			"OAUTH_STATE_EXPIRE": expectedConfig.OauthConfig.StateExpire.String(),

			"MAX_PUBLISH_RETRY_ATTEMPTS": strconv.FormatUint(uint64(expectedConfig.MaxPublishRetryAttempts), 10),

			"STORAGE_DRIVER": expectedConfig.Storage.Driver,
			"STORAGE_URI":    expectedConfig.Storage.URI,

			"DB_DRIVER": expectedConfig.DB.Driver,
			"DB_DSN":    expectedConfig.DB.DSN,

			"SERVER_HOST":       expectedConfig.Server.Host,
			"SERVER_PORT":       strconv.FormatUint(uint64(expectedConfig.Server.Port), 10),
			"SERVER_PUBLIC_URI": expectedConfig.Server.PublicURI,

			"SERVER_OMNI_API_TOKEN":   expectedConfig.Server.OmniAPIToken,
			"SERVER_CW_API_TOKEN":     expectedConfig.Server.CwAPIToken,
			"SERVER_SWAGGER_PASSWORD": expectedConfig.Server.SwaggerPassword,

			"AMQP_URL":            expectedConfig.AMQP.URL,
			"AMQP_PREFETCH_COUNT": strconv.FormatUint(uint64(expectedConfig.AMQP.PrefetchCount), 10),
			"AMQP_PREFETCH_SIZE":  strconv.FormatUint(uint64(expectedConfig.AMQP.PrefetchSize), 10),

			"KAFKA_BROKERS":                         strings.Join(expectedConfig.Kafka.Brokers, ","),
			"KAFKA_TOPIC_TRANSACTIONS":              expectedConfig.Kafka.TopicTransactions,
			"KAFKA_TOPIC_CONTENT_UPDATED":           expectedConfig.Kafka.TopicContentUpdated,
			"KAFKA_CONSUMERS_COUNT_TRANSACTIONS":    strconv.FormatUint(uint64(expectedConfig.Kafka.ConsumersCountTransactions), 10),
			"KAFKA_CONSUMERS_COUNT_CONTENT_UPDATED": strconv.FormatUint(uint64(expectedConfig.Kafka.ConsumersCountContentUpdated), 10),
			"KAFKA_SASL_MECHANISM":                  "SCRAM-SHA-256",
			"KAFKA_SASL_USERNAME":                   expectedConfig.Kafka.SASLConfig.Username,
			"KAFKA_SASL_PASSWORD":                   expectedConfig.Kafka.SASLConfig.Password,

			"KAFKA_MAX_ATTEMPTS": strconv.FormatUint(uint64(expectedConfig.Kafka.MaxAttempts), 10),
			"KAFKA_DIAL_TIMEOUT": expectedConfig.Kafka.DialTimeout.String(),

			"X_BASE_URI":        expectedConfig.X.BaseURI,
			"X_CLIENT_ID":       expectedConfig.X.ClientID,
			"X_CLIENT_SECRET":   expectedConfig.X.ClientSecret,
			"X_CONSUMER_KEY":    expectedConfig.X.ConsumerKey,
			"X_CONSUMER_SECRET": expectedConfig.X.ConsumerSecret,

			"BLUESKY_BASE_URI":      expectedConfig.Bluesky.BaseURI,
			"BLUESKY_CLIENT_ID":     expectedConfig.Bluesky.ClientID,
			"BLUESKY_CLIENT_SECRET": expectedConfig.Bluesky.ClientSecret,

			"RECTALOGIST_URI":     expectedConfig.RectalogistConfig.URI,
			"RECTALOGIST_CDN_URI": expectedConfig.RectalogistConfig.CdnURI,

			"CONTENT_DATA_MART_URI": expectedConfig.ContentDataMartConfig.URI,
			"CONTENT_DATA_MART_KEY": expectedConfig.ContentDataMartConfig.Key,
		}
	}

	unsetAll := func() {
		for key := range getEnvValues() {
			_ = os.Unsetenv(key)
		}
	}

	t.Run("FromEnvVars", func(t *testing.T) {
		unsetAll()
		for key, value := range getEnvValues() {
			t.Setenv(key, value)
		}

		cfg, err := LoadConfig()

		assert.NoErrorf(t, err, "got unexpected error %s", err)
		assertConfig(t, expectedConfig, cfg)
	})

	createTempWorkingDir := func() func() {
		dirName := t.TempDir()
		previousDir, _ := os.Getwd()
		t.Chdir(dirName)
		return func() {
			err := os.RemoveAll(dirName)
			if err == nil {
				t.Chdir(previousDir)
			}
			if err != nil {
				t.Fatal(err)
			}
		}
	}

	t.Run("FromFile", func(t *testing.T) {
		unsetAll()
		restoreWorkingDir := createTempWorkingDir()
		defer restoreWorkingDir()

		envFileContent := make([]byte, 0, 400)
		for key, value := range getEnvValues() {
			envFileContent = append(envFileContent, fmt.Sprintf("%s=%s\n", key, value)...)
		}

		err := os.WriteFile(".env", envFileContent, 0600)
		assert.NoErrorf(t, err, "got unexpected while write file .env error %s", err)

		config, err := LoadConfig()

		assert.NoError(t, err)
		assertConfig(t, expectedConfig, config)
	})

	t.Run("EmptyConfig", func(t *testing.T) {
		unsetAll()

		_, err := LoadConfig()
		assert.Error(t, err, "loadConfig() should exit with error, actual error is nil")
	})

	t.Run("EmptyPortAndEmptyTimeout", func(t *testing.T) {
		unsetAll()

		for key, value := range getEnvValues() {
			t.Setenv(key, value)
		}
		_ = os.Unsetenv("SERVER_PORT")
		_ = os.Unsetenv("KAFKA_DIAL_TIMEOUT")
		_ = os.Unsetenv("KAFKA_MAX_ATTEMPTS")

		config, err := LoadConfig()

		assert.NoError(t, err)
		assert.Equal(t, uint(9052), config.Server.Port)

		assert.Equal(t, time.Second*5, config.Kafka.DialTimeout)
		assert.Equal(t, 10, config.Kafka.MaxAttempts)
	})
}

func assertConfig(t *testing.T, expected Config, actual *Config) {
	assert.Equal(t, expected, *actual)
}
