package config

type OpenAIConfig struct {
	APIKey       string `env:"API_KEY"`
	BaseURL      string `env:"BASE_URL" envDefault:"https://api.openai.com/v1"`
	Model        string `env:"MODEL" envDefault:"whisper-1"`
	TimeoutSec   int    `env:"TIMEOUT_SEC" envDefault:"120"`
}
