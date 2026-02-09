package natsx

import (
	"context"
	"news-scrabber/internal/config"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func BuildNATSOptions(lc fx.Lifecycle, cfg *config.Config, log *zap.Logger) []nats.Option {
	// Build connection options
	var opts []nats.Option
	opts = append(opts,
		nats.MaxReconnects(10),
		nats.Name("news-scrabber"),
	)
	if cfg.NATS.User != "" {
		opts = append(opts, nats.UserInfo(cfg.NATS.User, cfg.NATS.Password))
	}

	return opts
}

func NewJetStream(lc fx.Lifecycle, log *zap.Logger, cfg *config.Config, opts []nats.Option) (jetstream.JetStream, error) {
	nc, err := nats.Connect(cfg.NATS.URL, opts...)
	if err != nil {
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			// Best-effort drain then close
			if err := nc.Drain(); err != nil {
				log.Warn("nats drain error", zap.Error(err))
			}
			nc.Close()
			return nil
		},
	})

	return jetstream.New(nc)
}
