package qdrant

import (
	"context"
	"net/http"
	"news-scrabber/internal/config"
)

// Client is a minimal placeholder HTTP client for Qdrant.
type Client struct {
	HTTP       *http.Client
	BaseURL    string
	APIKey     string
	Collection string
}

func NewClient(cfg *config.Config) (*Client, error) {
	return &Client{
		HTTP:       &http.Client{},
		BaseURL:    cfg.Qdrant.URL,
		APIKey:     cfg.Qdrant.APIKey,
		Collection: cfg.Qdrant.Collection,
	}, nil
}

// UpsertText is a placeholder that should upsert the text embedding into Qdrant.
// Implement embedding generation and points upsert with your preferred embedding model.
func (c *Client) UpsertText(ctx context.Context, collection, id, text string, meta map[string]any) error {
	// TODO: Implement actual embedding + upsert.
	return nil
}
