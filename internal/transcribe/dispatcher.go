package transcribe

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"news-scrabber/internal/config"

	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Dispatcher subscribes to NATS and runs VideoTranscribeRequested with bounded concurrency.
type Dispatcher struct {
	js            jetstream.JetStream
	log           *zap.Logger
	svc           *Service
	consumer      jetstream.Consumer
	stream        string
	subjects      string
	ctx           context.Context
	cancel        context.CancelFunc
	sem           chan struct{}
	maxConcurrent int
}

func (d *Dispatcher) processMessages() {
	msgs, err := d.consumer.Messages()
	if err != nil {
		d.log.Error("failed to get consumer messages", zap.Error(err))
		return
	}
	for {
		msg, err := msgs.Next()
		if err != nil {
			if d.ctx.Err() != nil {
				return // context cancelled
			}
			d.log.Error("error receiving message", zap.Error(err))
			continue
		}
		var ev VideoTranscribeRequested
		if err := json.Unmarshal(msg.Data(), &ev); err != nil {
			d.log.Warn("bad event payload", zap.Error(err))
			_ = msg.Nak()
			continue
		}
		if ev.URL == "" {
			d.log.Warn("missing url in event")
			_ = msg.Nak()
			continue
		}

		// Acquire a worker slot for bounded concurrency.
		select {
		case d.sem <- struct{}{}:
			// acquired
		case <-d.ctx.Done():
			_ = msg.Nak()
			return
		}

		go d.handleMessage(ev, msg)
	}
}

func (d *Dispatcher) handleMessage(ev VideoTranscribeRequested, msg jetstream.Msg) {
	defer func() { <-d.sem }()

	errCh := make(chan error, 1)
	go func() { errCh <- d.svc.IngestURL(d.ctx, ev.URL, ev.JobID) }()

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case err := <-errCh:
			if err != nil && !errors.Is(err, context.Canceled) {
				d.log.Warn("job finished with error", zap.Error(err), zap.String("url", ev.URL), zap.String("job", ev.JobID))
				_ = msg.Nak()
			} else {
				d.log.Info("job finished", zap.String("url", ev.URL), zap.String("job", ev.JobID))
				_ = msg.Ack()
			}
			return
		case <-ticker.C:
			if err := msg.InProgress(); err != nil {
				d.log.Debug("in-progress heartbeat failed", zap.Error(err))
			}
		case <-d.ctx.Done():
			return
		}
	}
}

func NewDispatcher(lc fx.Lifecycle, js jetstream.JetStream, log *zap.Logger, cfg *config.Config, svc *Service) (*Dispatcher, error) {
	stream := cfg.JetStream.EventsStream
	if stream == "" {
		stream = "NEWS"
	}
	subjects := cfg.JetStream.EventsSubjects
	if subjects == "" {
		subjects = "news.*"
	}

	maxConc := cfg.Transcribe.MaxConcurrent
	if maxConc <= 0 {
		maxConc = 2
	}

	d := &Dispatcher{
		log:           log.With(zap.String("component", "transcribe.dispatcher")),
		js:            js,
		svc:           svc,
		stream:        stream,
		subjects:      subjects,
		sem:           make(chan struct{}, maxConc),
		maxConcurrent: maxConc,
	}
	d.ctx, d.cancel = context.WithCancel(context.Background())

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			// Ensure stream exists (idempotent) to avoid "no response from stream" errors.
			if _, err := d.js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
				Name:     d.stream,
				Subjects: []string{d.subjects, SubjectVideoTranscribeRequested},
			}); err != nil {
				d.log.Warn("ensure events stream failed", zap.Error(err), zap.String("stream", d.stream), zap.String("subjects", d.subjects))
				return err
			}

			consumer, err := d.js.CreateConsumer(ctx, d.stream, jetstream.ConsumerConfig{
				Durable:       "transcribe-dispatcher",
				AckPolicy:     jetstream.AckExplicitPolicy,
				AckWait:       30 * time.Second,
				MaxAckPending: d.maxConcurrent,
				FilterSubject: SubjectVideoTranscribeRequested,
			})
			if err != nil {
				return err
			}
			d.consumer = consumer
			go d.processMessages()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			if d.cancel != nil {
				d.cancel()
			}
			if d.consumer != nil {
				_ = d.js.DeleteConsumer(ctx, d.stream, "transcribe-dispatcher")
			}
			return nil
		},
	})

	return d, nil
}
