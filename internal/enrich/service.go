package enrich

import (
	"context"
	"news-scrabber/internal/config"

	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Service performs enrichment of parsed/transcribed items and republishes enriched events.
type Service struct {
	log *zap.Logger
	cfg *config.Config
	js  jetstream.JetStream
}

func NewService(lc fx.Lifecycle, log *zap.Logger, cfg *config.Config, js jetstream.JetStream) (*Service, error) {
	s := &Service{
		log: log.With(zap.String("component", "enrich")),
		cfg: cfg,
		js:  js,
	}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			s.log.Info("enrichment service started (placeholder)")
			return nil
		},
		OnStop: func(ctx context.Context) error { return nil },
	})
	return s, nil
}
