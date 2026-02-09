package enrich

import (
	"context"
	"news-scrabber/internal/config"
	"news-scrabber/internal/natsx"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Service performs enrichment of parsed/transcribed items and republishes enriched events.
type Service struct {
	cfg *config.Config
	js  *natsx.JetStream
	log *zap.Logger
}

func NewService(lc fx.Lifecycle, cfg *config.Config, js *natsx.JetStream, log *zap.Logger) (*Service, error) {
	s := &Service{cfg: cfg, js: js, log: log.With(zap.String("component", "enrich"))}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			s.log.Info("enrichment service started (placeholder)")
			return nil
		},
		OnStop: func(ctx context.Context) error { return nil },
	})
	return s, nil
}
