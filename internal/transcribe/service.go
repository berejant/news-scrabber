package transcribe

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

// jobRequest represents a pending transcribe job submission queued to workers.
type jobRequest struct {
	ctx   context.Context
	url   string
	jobID string
}

// Service orchestrates ingest jobs; it exposes an API to submit URLs for processing.
// It also owns a worker-pool dispatcher to cap concurrent transcriptions.
type Service struct {
	deps JobParams
	log  *zap.Logger

	queue       chan jobRequest
	wg          sync.WaitGroup
	cancel      context.CancelFunc
	serviceCtx  context.Context
	workerCount int
}

func NewService(lc fx.Lifecycle, p JobParams) (*Service, error) {
	s := &Service{deps: p, log: p.Log.With(zap.String("component", "transcribe"))}

	// derive worker pool settings from config
	wc := p.Cfg.Transcribe.MaxConcurrent
	if wc <= 0 {
		wc = 2
	}
	qs := p.Cfg.Transcribe.QueueSize
	if qs <= 0 {
		qs = 100
	}
	s.workerCount = wc
	s.queue = make(chan jobRequest, qs)
	s.serviceCtx, s.cancel = context.WithCancel(context.Background())

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			s.log.Info("transcribe dispatcher starting", zap.Int("workers", wc), zap.Int("queue", qs))
			for i := 0; i < wc; i++ {
				wID := i
				s.wg.Add(1)
				go s.workerLoop(wID)
			}
			return nil
		},
		OnStop: func(ctx context.Context) error {
			s.log.Info("transcribe dispatcher stopping")
			s.cancel()
			// do not close queue to avoid panics if someone enqueues; just drain via ctx cancel
			done := make(chan struct{})
			go func() { s.wg.Wait(); close(done) }()
			select {
			case <-done:
				return nil
			case <-ctx.Done():
				return ctx.Err()
			}
		},
	})

	return s, nil
}

func (s *Service) workerLoop(workerID int) {
	defer s.wg.Done()
	log := s.log.With(zap.Int("worker", workerID))
	for {
		select {
		case <-s.serviceCtx.Done():
			return
		case req := <-s.queue:
			if req.ctx == nil {
				req.ctx = s.serviceCtx
			}
			ctx, cancel := context.WithCancel(req.ctx)
			// ensure jobID
			jid := req.jobID
			if jid == "" {
				jid = fmt.Sprintf("job-%d", time.Now().UnixNano())
			}
			log.Info("starting job", zap.String("job", jid), zap.String("url", req.url))
			job := NewIngestJob(s.deps, jid, req.url)
			if err := job.Start(ctx); err != nil && !errors.Is(err, context.Canceled) {
				log.Warn("job finished with error", zap.String("job", jid), zap.Error(err))
			} else {
				log.Info("job finished", zap.String("job", jid))
			}
			cancel()
		}
	}
}

var ErrQueueFull = errors.New("transcribe queue is full")

// IngestURL submits a URL for processing. If jobID is empty it's auto-generated.
// It enqueues the job to the dispatcher; returns ErrQueueFull if the queue is saturated.
func (s *Service) IngestURL(ctx context.Context, url, jobID string) error {
	if url == "" {
		return errors.New("url is required")
	}
	if jobID == "" {
		jobID = fmt.Sprintf("job-%d", time.Now().Unix())
	}
	req := jobRequest{ctx: ctx, url: url, jobID: jobID}
	select {
	case s.queue <- req:
		s.log.Info("enqueued ingest job", zap.String("job", jobID), zap.String("url", url), zap.Int("queue_len", len(s.queue)))
		return nil
	case <-ctx.Done():
		return ctx.Err()
	default:
		// queue full
		s.log.Warn("ingest queue full", zap.String("job", jobID), zap.String("url", url))
		return ErrQueueFull
	}
}
