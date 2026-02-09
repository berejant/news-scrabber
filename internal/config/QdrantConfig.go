package config

type QdrantConfig struct {
	URL       string `env:"URL" envDefault:"http://localhost:6333"`
	APIKey    string `env:"API_KEY"`
	Collection string `env:"COLLECTION" envDefault:"news"`
}
