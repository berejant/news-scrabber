package config

import (
	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

// LoadConfig loads configuration from environment variables.
func LoadConfig() (cfg *Config, err error) {
	_ = godotenv.Load()

	opts := env.Options{
		UseFieldNameByDefault: true,
	}

	cfg = &Config{}
	err = env.ParseWithOptions(cfg, opts)

	return cfg, err
}
