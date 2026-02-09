package whisper

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"news-scrabber/internal/config"
)

// Client is a minimal HTTP client for linuxserver/faster-whisper REST API.
// See: https://github.com/linuxserver/docker-faster-whisper
// This provides a basic health check and is extendable for transcription endpoints.
type Client struct {
	HTTP    *http.Client
	BaseURL string
	Model   string
}

func NewClient(cfg *config.Config) (*Client, error) {
	return &Client{
		HTTP:    &http.Client{Timeout: 30 * time.Second},
		BaseURL: cfg.Whisper.URL,
		Model:   cfg.Whisper.Model,
	}, nil
}

// Health checks service availability; linuxserver/faster-whisper exposes root page.
func (c *Client) Health(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+"/", nil)
	if err != nil {
		return err
	}
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 200 && resp.StatusCode < 500 { // returns 200 or 404 depending on route; just reachability check
		return nil
	}
	return fmt.Errorf("whisper health failed: status=%d", resp.StatusCode)
}
