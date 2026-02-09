package config

import (
	"reflect"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

// ParseSASLMechanism parses SASL mechanism.
func ParseSASLMechanism(text string) (interface{}, error) {
	var mechanism SASLMechanism
	return mechanism, mechanism.UnmarshalText(text)
}

// LoadConfig loads configuration from environment variables.
func LoadConfig() (cfg *Config, err error) {
	_ = godotenv.Load()

	funcMap := map[reflect.Type]env.ParserFunc{
		reflect.TypeOf(SASLMechanism(0)): ParseSASLMechanism,
	}

	opts := env.Options{
		UseFieldNameByDefault: true,
		FuncMap:               funcMap,
	}

	cfg = &Config{}
	err = env.ParseWithOptions(cfg, opts)

	return cfg, err
}
