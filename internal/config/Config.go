package config

type Config struct {
	App         string `env:"APP" envDefault:"news-scrabber"`
	Environment string `env:"ENVIRONMENT" envDefault:"development"`

	// Legacy/other modules (kept for compatibility if used elsewhere)
	Server ServerConfig `envPrefix:"SERVER_"`

	// New modules for news-scrabber app
	NATS       NATSConfig       `envPrefix:"NATS_"`
	JetStream  JetStreamConfig  `envPrefix:"JS_"`
	S3         S3Config         `envPrefix:"S3_"`
	Qdrant     QdrantConfig     `envPrefix:"QDRANT_"`
	OpenAI     OpenAIConfig     `envPrefix:"OPENAI_"`
	Transcribe TranscribeConfig `envPrefix:"TRANSCRIBE_"`
	Scraper    ScraperConfig    `envPrefix:"SCRAPER_"`

	MaxPublishRetryAttempts uint8 `env:"MAX_PUBLISH_RETRY_ATTEMPTS" envDefault:"10"`

	DatadogHost string `env:"DD_DOGSTATSD_HOST"`
}

func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}

func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

func (c *Config) IsTest() bool {
	return c.Environment == "test"
}

func (c *Config) IsLocal() bool {
	return c.Environment == "local"
}

func (c *Config) IsStaging() bool {
	return c.Environment == "staging"
}

func (c *Config) AppName() string {
	if c.Environment == "production" {
		return c.App
	}
	return c.App + "-" + c.Environment
}
