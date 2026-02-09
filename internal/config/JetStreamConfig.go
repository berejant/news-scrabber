package config

type JetStreamConfig struct {
	KVBucket       string `env:"KV_BUCKET" envDefault:"news_kv"`
	KVHistory      uint8  `env:"KV_HISTORY" envDefault:"1"`
	KVDescription  string `env:"KV_DESCRIPTION" envDefault:"Key-Value store for news-scrabber"`
	EventsStream   string `env:"EVENTS_STREAM" envDefault:"NEWS"`
	EventsSubjects string `env:"EVENTS_SUBJECTS" envDefault:"news.*"`
}
