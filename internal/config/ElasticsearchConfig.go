package config

// ElasticsearchConfig holds connection settings for Elasticsearch cluster.
type ElasticsearchConfig struct {
	URL         string `env:"URL" envDefault:"http://localhost:9200"`
	APIKey      string `env:"API_KEY"`
	Username    string `env:"USERNAME"`
	Password    string `env:"PASSWORD"`
	IndexPrefix string `env:"INDEX_PREFIX" envDefault:"news"`
}
