package s3client

import (
	"context"
	"net/http"
	"news-scrabber/internal/config"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Client is a minimal placeholder S3 client suitable for Seaweed S3 gateway.
type Client struct {
	HTTP     *http.Client
	Endpoint string
	Region   string
	Bucket   string
	AccessKey string
	SecretKey string
	UseSSL    bool
}

func New(lc fx.Lifecycle, cfg *config.Config, log *zap.Logger) (*Client, error) {
	c := &Client{
		HTTP:      &http.Client{},
		Endpoint:  cfg.S3.Endpoint,
		Region:    cfg.S3.Region,
		Bucket:    cfg.S3.Bucket,
		AccessKey: cfg.S3.AccessKey,
		SecretKey: cfg.S3.SecretKey,
		UseSSL:    cfg.S3.UseSSL,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Info("S3 client ready (placeholder)",
				zap.String("endpoint", c.Endpoint),
				zap.String("bucket", c.Bucket),
			)
			return nil
		},
		OnStop: func(ctx context.Context) error { return nil },
	})

	return c, nil
}
