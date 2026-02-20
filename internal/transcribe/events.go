package transcribe

import (
	"context"
	"encoding/json"
	"time"

	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/zap"
)

// SubjectVideoTranscribeRequested is the NATS subject for requesting a new video transcription job.
const SubjectVideoTranscribeRequested = "news.transcribe.request"

// VideoTranscribeRequested is an event requesting to start transcription for a given video/stream URL.
// Evolve by adding fields; keep existing fields backward compatible.
type VideoTranscribeRequested struct {
	Event       string    `json:"event"`
	URL         string    `json:"url"`
	JobID       string    `json:"job_id,omitempty"`
	RequestedAt time.Time `json:"requested_at"`
}

// TranscribeEventPublisher defines the interface to publish transcription-related events.
type TranscribeEventPublisher interface {
	// PublishVideoTranscribeRequested publishes a request event to start transcription.
	// If jobID is empty, it will be auto-generated and returned.
	PublishVideoTranscribeRequested(ctx context.Context, url, jobID string) (string, error)
}

// jsPublisher implements TranscribeEventPublisher using NATS JetStream.
type jsPublisher struct {
	js  jetstream.JetStream
	log *zap.Logger
}

// NewPublisher returns a JetStream-backed publisher for transcription events.
func NewPublisher(js jetstream.JetStream, log *zap.Logger) TranscribeEventPublisher {
	return &jsPublisher{js: js, log: log.With(zap.String("component", "transcribe.publisher"))}
}

func (p *jsPublisher) PublishVideoTranscribeRequested(ctx context.Context, url, jobID string) (string, error) {
	if jobID == "" {
		jobID = generateJobID()
	}
	ev := VideoTranscribeRequested{
		Event:       SubjectVideoTranscribeRequested,
		URL:         url,
		JobID:       jobID,
		RequestedAt: time.Now().UTC(),
	}
	b, err := json.Marshal(ev)
	if err != nil {
		return "", err
	}
	if _, err := p.js.Publish(ctx, SubjectVideoTranscribeRequested, b); err != nil {
		return "", err
	}
	p.log.Info("published VideoTranscribeRequested", zap.String("job", jobID), zap.String("url", url))
	return jobID, nil
}

// generateJobID returns a reasonably unique identifier for jobs.
func generateJobID() string {
	return time.Now().UTC().Format("20060102T150405.000000000Z07:00")
}
