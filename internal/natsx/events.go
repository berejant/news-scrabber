package natsx

import (
	"context"

	"news-scrabber/internal/config"

	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// EnsureEventsStream creates or updates a JetStream stream to capture application events.
// It uses cfg.JetStream.EventsStream for the stream name and cfg.JetStream.EventsSubjects for subjects pattern.
func EnsureEventsStream(lc fx.Lifecycle, js jetstream.JetStream, cfg *config.Config, log *zap.Logger) error {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			name := cfg.JetStream.EventsStream
			subj := cfg.JetStream.EventsSubjects
			if name == "" { name = "NEWS" }
			if subj == "" { subj = "news.*" }
			_, err := js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
				Name:     name,
				Subjects: []string{subj},
			})
			if err != nil { log.Warn("ensure events stream failed", zap.Error(err)) }
			return nil
		},
	})
	return nil
}
