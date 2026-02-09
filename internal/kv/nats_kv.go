package kv

import (
	"news-scrabber/internal/config"

	"github.com/gofiber/storage/nats"
	natsgo "github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

// NewKVStore is a simplified factory that wires NATS/JetStream KV using the existing
// NATS connection from the natsx package and lifecycle hooks. It prepares KV resources
// on start and cleans up on stop. Placeholder implementation returns nil KVStore for now.
func NewKVStore(cfg *config.Config, opts []natsgo.Option) KVStore {
	return nats.New(nats.Config{
		URLs:        cfg.NATS.URL,
		NatsOptions: opts,
		KeyValueConfig: jetstream.KeyValueConfig{
			Bucket:  cfg.JetStream.KVBucket,
			Storage: jetstream.FileStorage,
		},
	})
}
