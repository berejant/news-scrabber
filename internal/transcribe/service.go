package transcribe

import (
	"context"
	"news-scrabber/internal/config"
	"news-scrabber/internal/natsx"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Service orchestrates ffmpeg-based audio extraction and OpenAI transcription.
// This is a placeholder; implement actual ffmpeg invocation and OpenAI API calls as needed.
type Service struct {
	cfg *config.Config
	js  *natsx.JetStream
	log *zap.Logger
}

func NewService(lc fx.Lifecycle, cfg *config.Config, js *natsx.JetStream, log *zap.Logger) (*Service, error) {
	s := &Service{cfg: cfg, js: js, log: log.With(zap.String("component", "transcribe"))}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			s.log.Info("transcribe service started (placeholder)",
				zap.String("ffmpeg", cfg.Transcribe.FFmpegPath),
				zap.String("model", cfg.OpenAI.Model),
			)
			return nil
		},
		OnStop: func(ctx context.Context) error { return nil },
	})
	return s, nil
}
