package config

type S3Config struct {
	Endpoint  string `env:"ENDPOINT" envDefault:"http://localhost:8333"`
	Region    string `env:"REGION" envDefault:"us-east-1"`
	Bucket    string `env:"BUCKET" envDefault:"news"`
	AccessKey string `env:"ACCESS_KEY"`
	SecretKey string `env:"SECRET_KEY"`
	UseSSL    bool   `env:"USE_SSL" envDefault:"false"`
}
