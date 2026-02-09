package config

type TranscribeConfig struct {
	FFmpegPath string `env:"FFMPEG_PATH" envDefault:"ffmpeg"`
	TempDir    string `env:"TEMP_DIR" envDefault:"/tmp/news-scrabber"`
}
