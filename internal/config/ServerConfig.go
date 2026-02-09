package config

// ServerConfig is a config for server.
type ServerConfig struct {
	Host      string `env:"HOST"`
	Port      uint   `env:"PORT"  envDefault:"9052"`
	PublicURI string `env:"PUBLIC_URI"`

	OmniAPIToken    string `env:"OMNI_API_TOKEN"`
	CwAPIToken      string `env:"CW_API_TOKEN"`
	BiAPIToken      string `env:"BI_API_TOKEN"`
	WebhookToken    string `env:"WEBHOOK_TOKEN"`
	SwaggerPassword string `env:"SWAGGER_PASSWORD"`
}
