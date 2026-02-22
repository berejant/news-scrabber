package transcribe

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"
)

// Service exposes an API to run ingest/transcribe jobs.
// Concurrency and queuing are managed by NATS JetStream (see Dispatcher).
type Service struct {
	deps JobParams
	log  *zap.Logger
}

func NewService(p JobParams) (*Service, error) {
	return &Service{
		log:  p.Log.With(zap.String("component", "transcribe")),
		deps: p,
	}, nil
}

// IngestURL runs a single ingest/transcribe job synchronously and returns when finished.
// The dispatcher (NATS consumer) should call this under its own concurrency control.
func (s *Service) IngestURL(ctx context.Context, url, jobID string) error {
	if url == "" {
		return errors.New("url is required")
	}
	if jobID == "" {
		jobID = fmt.Sprintf("job-%d", time.Now().UnixNano())
	}
	return NewIngestJob(s.deps, jobID, url).Start(ctx)
}
