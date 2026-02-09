package config

type NATSConfig struct {
	URL            string `env:"URL" envDefault:"nats://localhost:4222"`
	User           string `env:"USER"`
	Password       string `env:"PASSWORD"`
	EnableJetStream bool   `env:"ENABLE_JETSTREAM" envDefault:"true"`
}
