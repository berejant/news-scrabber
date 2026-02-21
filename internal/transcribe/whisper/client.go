package whisper

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
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
	to := 60 * time.Second
	if cfg.Whisper.TimeoutSeconds > 0 {
		to = time.Duration(cfg.Whisper.TimeoutSeconds) * time.Second
	}
	return &Client{
		HTTP:    &http.Client{Timeout: to},
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

// TranscribeFile sends a local audio file to the whisper server and returns plain text.
// This attempts linuxserver/faster-whisper compatible REST: POST /inference with multipart field "audio_file" and optional "model".
func (c *Client) TranscribeFile(ctx context.Context, path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	pr, pw := io.Pipe()
	mw := multipart.NewWriter(pw)

	// stream multipart body in background
	done := make(chan error, 1)
	go func() {
		defer pw.Close()
		defer mw.Close()
		// file part
		fw, err := mw.CreateFormFile("audio_file", filepath.Base(path))
		if err != nil {
			done <- err
			return
		}
		if _, err := io.Copy(fw, f); err != nil {
			done <- err
			return
		}
		// model (optional)
		_ = mw.WriteField("model", c.Model)
		done <- nil
	}()

	url := strings.TrimRight(c.BaseURL, "/") + "/inference"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, pr)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", mw.FormDataContentType())
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if derr := <-done; derr != nil {
		return "", derr
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("whisper http %d: %s", resp.StatusCode, string(b))
	}
	// best-effort parse JSON {"text":"..."} or {"segments":[{"text":"..."},...]}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var jmap map[string]any
	if err := json.Unmarshal(b, &jmap); err == nil {
		if t, ok := jmap["text"].(string); ok && t != "" {
			return t, nil
		}
		if segs, ok := jmap["segments"].([]any); ok {
			var sb strings.Builder
			for _, s := range segs {
				if m, ok := s.(map[string]any); ok {
					if t, ok := m["text"].(string); ok {
						sb.WriteString(t)
						sb.WriteByte(' ')
					}
				}
			}
			return strings.TrimSpace(sb.String()), nil
		}
	}
	// fallback: return raw body as string
	return strings.TrimSpace(string(b)), nil
}
