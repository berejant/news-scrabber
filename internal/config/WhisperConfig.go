package config

// WhisperConfig holds connection for Faster-Whisper HTTP server (linuxserver/faster-whisper)
type WhisperConfig struct {
	URL            string `env:"URL" envDefault:"http://localhost:10300"`
	Model          string `env:"MODEL" envDefault:"base"`
	TimeoutSeconds int    `env:"TIMEOUT_SECONDS" envDefault:"600"`
}
