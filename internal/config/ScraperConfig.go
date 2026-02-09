package config

type ScraperConfig struct {
	UserAgent   string   `env:"USER_AGENT" envDefault:"news-scrapper-bot/1.0"`
	Seeds       []string `env:"SEEDS" envSeparator:","`
	Concurrency int      `env:"CONCURRENCY" envDefault:"4"`
	RequestTimeoutSec int `env:"REQUEST_TIMEOUT_SEC" envDefault:"10"`
}
