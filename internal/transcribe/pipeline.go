package transcribe

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"news-scrabber/internal/config"
	"news-scrabber/internal/search/elasticsearch"
	"news-scrabber/internal/storage/s3client"
	"news-scrabber/internal/transcribe/whisper"
	"news-scrabber/internal/vector/qdrant"
)

// RawContentReadyEvent is emitted to NATS on every processed 30s chunk.
// It contains the latest chunk text and a rolling window of the previous 6 chunks (total ~3.5 minutes).
// Subject: "RawContentReady"
// NOTE: Keep this backward compatible; evolve by adding new json fields.
type RawContentReadyEvent struct {
	Event        string    `json:"event"`
	SourceURL    string    `json:"source_url"`
	JobID        string    `json:"job_id"`
	ChunkIndex   int       `json:"chunk_index"`
	ChunkSeconds int       `json:"chunk_seconds"`
	ChunkText    string    `json:"chunk_text"`
	WindowText   string    `json:"window_text"`
	S3Key        string    `json:"s3_key"`
	CreatedAt    time.Time `json:"created_at"`
}

// JobParams bundles dependencies for ingest jobs.
type JobParams struct {
	fx.In

	Cfg *config.Config
	JS  jetstream.JetStream
	Log *zap.Logger
	S3  *s3client.Client
	WH  *whisper.Client
	ES  *elasticsearch.Client
	Vec *qdrant.Client
}

// IngestJob coordinates ffmpeg segmentation and per-chunk processing.
type IngestJob struct {
	jobID     string
	sourceURL string
	tempDir   string

	cfg *config.Config
	js  jetstream.JetStream
	log *zap.Logger
	s3  *s3client.Client
	wh  *whisper.Client
	es  *elasticsearch.Client
	vec *qdrant.Client

	// internal
	mu           sync.Mutex
	chunkWindow  []string // last 7 chunk texts
	processedSet map[string]struct{}
}

func NewIngestJob(params JobParams, jobID, sourceURL string) *IngestJob {
	tempDir := filepath.Join(params.Cfg.Transcribe.TempDir, jobID)
	return &IngestJob{
		jobID:        jobID,
		sourceURL:    sourceURL,
		tempDir:      tempDir,
		cfg:          params.Cfg,
		js:           params.JS,
		log:          params.Log,
		s3:           params.S3,
		wh:           params.WH,
		es:           params.ES,
		vec:          params.Vec,
		processedSet: make(map[string]struct{}),
	}
}

// Start launches ffmpeg segmenter and the watcher until context is done.
func (j *IngestJob) Start(ctx context.Context) error {
	if err := os.MkdirAll(j.tempDir, 0o755); err != nil {
		return err
	}

	j.log.Info("starting ffmpeg segmenter", zap.String("url", j.sourceURL), zap.String("dir", j.tempDir))
	segPattern := filepath.Join(j.tempDir, "segment_%05d.wav")
	cmd := exec.CommandContext(ctx, j.cfg.Transcribe.FFmpegPath,
		"-hide_banner", "-loglevel", "error",
		"-i", j.sourceURL,
		"-vn",
		"-ac", "1",
		"-ar", "16000",
		"-c:a", "pcm_s16le",
		"-f", "segment",
		"-segment_time", "30",
		"-reset_timestamps", "1",
		segPattern,
	)

	// Start watcher first to not miss early files
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		j.watchAndProcess(ctx)
	}()

	// Run ffmpeg (will exit when ctx is canceled or source ends)
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("ffmpeg start: %w", err)
	}
	go func() {
		if err := cmd.Wait(); err != nil {
			if ctx.Err() == nil { // log only if not due to context cancel
				j.log.Warn("ffmpeg exited with error", zap.Error(err))
			}
		} else {
			j.log.Info("ffmpeg finished")
		}
	}()

	// Block until context canceled
	<-ctx.Done()
	j.log.Info("ingest job context done, waiting watcher")
	wg.Wait()
	return nil
}

func (j *IngestJob) watchAndProcess(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			j.scanOnce(ctx)
		}
	}
}

func (j *IngestJob) scanOnce(ctx context.Context) {
	var files []string
	filepath.WalkDir(j.tempDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			return nil
		}
		if strings.HasPrefix(filepath.Base(path), "segment_") && strings.HasSuffix(path, ".wav") {
			if _, ok := j.processedSet[path]; !ok {
				files = append(files, path)
			}
		}
		return nil
	})
	if len(files) == 0 {
		return
	}
	sort.Strings(files)
	for _, f := range files {
		if err := j.processOne(ctx, f); err != nil {
			j.log.Warn("process chunk failed", zap.String("file", f), zap.Error(err))
			continue
		}
		j.processedSet[f] = struct{}{}
	}
}

func (j *IngestJob) processOne(ctx context.Context, path string) error {
	idx, err := parseIndex(path)
	if err != nil {
		return err
	}
	key := filepath.Join("raw", j.jobID, filepath.Base(path))

	// 1) Upload to S3 (placeholder client provides best-effort)
	s3Key, err := j.s3.Upload(ctx, key, path)
	if err != nil {
		return fmt.Errorf("s3 upload: %w", err)
	}

	// 2) Transcribe with Whisper
	text, err := j.wh.TranscribeFile(ctx, path)
	if err != nil {
		return fmt.Errorf("whisper: %w", err)
	}

	// 3) Upload transcribed text to S3 as well
	textKey := filepath.Join("raw", j.jobID, fmt.Sprintf("segment_%05d.txt", idx))
	// write to temp file then upload via placeholder client
	tf := filepath.Join(j.tempDir, fmt.Sprintf("segment_%05d.txt", idx))
	if err := os.WriteFile(tf, []byte(text), 0o644); err != nil {
		j.log.Warn("write txt temp failed", zap.Error(err))
	} else {
		if _, err := j.s3.Upload(ctx, textKey, tf); err != nil {
			j.log.Warn("s3 upload txt failed", zap.Error(err))
		}
	}

	// 4) Save text to Elasticsearch
	doc := map[string]any{
		"job_id":        j.jobID,
		"source_url":    j.sourceURL,
		"chunk_index":   idx,
		"chunk_seconds": 30,
		"text":          text,
		"s3_key":        s3Key,
		"text_s3_key":   textKey,
		"created_at":    time.Now().UTC().Format(time.RFC3339Nano),
	}
	if err := j.es.IndexText(ctx, "raw-content", fmt.Sprintf("%s-%05d", j.jobID, idx), doc); err != nil {
		j.log.Warn("elasticsearch index failed", zap.Error(err))
	}

	// 4) Upsert into Qdrant (placeholder: may be no-op)
	if err := j.vec.UpsertText(ctx, "raw-content", fmt.Sprintf("%s-%05d", j.jobID, idx), text, map[string]any{
		"job_id":      j.jobID,
		"chunk_index": idx,
		"source_url":  j.sourceURL,
	}); err != nil {
		j.log.Warn("qdrant upsert failed", zap.Error(err))
	}

	// 5) Emit NATS event with rolling window (last 7 chunks)
	window := j.appendAndWindow(text, 7)
	ev := RawContentReadyEvent{
		Event:        "RawContentReady",
		SourceURL:    j.sourceURL,
		JobID:        j.jobID,
		ChunkIndex:   idx,
		ChunkSeconds: 30,
		ChunkText:    text,
		WindowText:   window,
		S3Key:        s3Key,
		CreatedAt:    time.Now().UTC(),
	}
	b, _ := json.Marshal(ev)
	if _, err := j.js.Publish(ctx, "RawContentReady", b); err != nil {
		j.log.Warn("nats publish failed", zap.Error(err))
	}
	return nil
}

func (j *IngestJob) appendAndWindow(text string, n int) string {
	j.mu.Lock()
	defer j.mu.Unlock()
	j.chunkWindow = append(j.chunkWindow, text)
	if len(j.chunkWindow) > n {
		j.chunkWindow = j.chunkWindow[len(j.chunkWindow)-n:]
	}
	return strings.TrimSpace(strings.Join(j.chunkWindow, " "))
}

func parseIndex(path string) (int, error) {
	base := filepath.Base(path) // segment_00001.wav
	if !strings.HasPrefix(base, "segment_") || !strings.HasSuffix(base, ".wav") {
		return 0, errors.New("unexpected segment filename")
	}
	mid := strings.TrimSuffix(strings.TrimPrefix(base, "segment_"), ".wav")
	var i int
	_, err := fmt.Sscanf(mid, "%05d", &i)
	return i, err
}
