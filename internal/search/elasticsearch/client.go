package elasticsearch

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"news-scrabber/internal/config"
)

// Client is a minimal HTTP-based client for Elasticsearch.
// For now we rely on simple HTTP calls to keep dependencies light.
// You can swap to the official go-elasticsearch client later if needed.
type Client struct {
	HTTP    *http.Client
	BaseURL string
}

func NewClient(cfg *config.Config) (*Client, error) {
	return &Client{
		HTTP: &http.Client{Timeout: 10 * time.Second},
		BaseURL: cfg.Elasticsearch.URL,
	}, nil
}

// Ping checks that the cluster is reachable.
func (c *Client) Ping(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL, nil)
	if err != nil {
		return err
	}
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}
	return fmt.Errorf("elasticsearch ping failed: status=%d", resp.StatusCode)
}
