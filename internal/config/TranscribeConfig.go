package config

type TranscribeConfig struct {
	FFmpegPath   string `env:"FFMPEG_PATH" envDefault:"ffmpeg"`
	TempDir      string `env:"TEMP_DIR" envDefault:"/tmp/news-scrabber"`
	MaxConcurrent int    `env:"MAX_CONCURRENT" envDefault:"2"`
	QueueSize     int    `env:"QUEUE_SIZE" envDefault:"100"`
}
