package config

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfigFromEnvVars(t *testing.T) {
	expectedConfig := Config{
		App:         "social-posting-service",
		Environment: "test",

		Server: ServerConfig{
			Host:            "localhost",
			Port:            8080,
			PublicURI:       "http://localhost:8080",
			OmniAPIToken:    "omni-token",
			CwAPIToken:      "cw-token",
			SwaggerPassword: "swagger-password",
		},

		MaxPublishRetryAttempts: 5,
	}

	getEnvValues := func() map[string]string {
		return map[string]string{
			"APP":         expectedConfig.App,
			"ENVIRONMENT": expectedConfig.Environment,

			"MAX_PUBLISH_RETRY_ATTEMPTS": strconv.FormatUint(uint64(expectedConfig.MaxPublishRetryAttempts), 10),

			"SERVER_HOST":       expectedConfig.Server.Host,
			"SERVER_PORT":       strconv.FormatUint(uint64(expectedConfig.Server.Port), 10),
			"SERVER_PUBLIC_URI": expectedConfig.Server.PublicURI,

			"SERVER_OMNI_API_TOKEN":   expectedConfig.Server.OmniAPIToken,
			"SERVER_CW_API_TOKEN":     expectedConfig.Server.CwAPIToken,
			"SERVER_SWAGGER_PASSWORD": expectedConfig.Server.SwaggerPassword,
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

	})
}

func assertConfig(t *testing.T, expected Config, actual *Config) {
	assert.Equal(t, expected, *actual)
}
