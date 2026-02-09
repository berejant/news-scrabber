package scraper

import (
	"context"
	"net/http"
	"news-scrabber/internal/config"
	"news-scrabber/internal/natsx"
	"time"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Service crawls HTML websites and publishes raw items to the event bus (NATS subjects under news.raw).
// This is a minimal placeholder; replace HTTP fetching and parsing with a real implementation.
type Service struct {
	cfg  *config.Config
	log  *zap.Logger
	js   *natsx.JetStream
	http *http.Client
}

func NewService(lc fx.Lifecycle, cfg *config.Config, log *zap.Logger, js *natsx.JetStream) (*Service, error) {
	s := &Service{
		cfg:  cfg,
		log:  log.With(zap.String("component", "scraper")),
		js:   js,
		http: &http.Client{Timeout: time.Duration(cfg.Scraper.RequestTimeoutSec) * time.Second},
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go s.run()
			return nil
		},
		OnStop: func(ctx context.Context) error { return nil },
	})

	return s, nil
}

func (s *Service) run() {
	s.log.Info("scraper started", zap.Int("seeds", len(s.cfg.Scraper.Seeds)))
	// Placeholder: iterate seeds and log; real impl would fetch HTML and publish to NATS subjects
	for _, u := range s.cfg.Scraper.Seeds {
		s.log.Info("would scrape", zap.String("url", u))
	}
}
