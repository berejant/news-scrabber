package s3client

import (
	"context"
	"net/http"
	"news-scrabber/internal/config"
	"os"
	"path/filepath"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Client is a minimal placeholder S3 client suitable for Seaweed S3 gateway.
type Client struct {
	HTTP      *http.Client
	Endpoint  string
	Region    string
	Bucket    string
	AccessKey string
	SecretKey string
	UseSSL    bool
	log       *zap.Logger
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
		log:       log.With(zap.String("component", "s3")),
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

// Upload uploads a local file under the provided key. Placeholder returns the key without real upload.
// Replace with a proper S3 SDK or SeaweedFS S3 signed request.
func (c *Client) Upload(ctx context.Context, key, localPath string) (string, error) {
	if _, err := os.Stat(localPath); err != nil {
		return "", err
	}
	if key == "" {
		key = filepath.Base(localPath)
	}
	c.log.Info("S3 upload (placeholder)", zap.String("key", key), zap.String("path", localPath))
	// TODO: implement actual upload using AWS SDK or MinIO SDK.
	return key, nil
}
