package qdrant

import (
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
