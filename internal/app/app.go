package app

import (
	"news-scrabber/internal/bootstrap"
	"news-scrabber/internal/config"
	"news-scrabber/internal/enrich"
	"news-scrabber/internal/kv"
	"news-scrabber/internal/natsx"
	"news-scrabber/internal/scraper"
	"news-scrabber/internal/search/elasticsearch"
	"news-scrabber/internal/server"
	"news-scrabber/internal/storage/s3client"
	"news-scrabber/internal/transcribe"
	"news-scrabber/internal/transcribe/whisper"
	"news-scrabber/internal/vector/qdrant"

	"go.uber.org/fx"
)

// Serve builds the News Scrabber application with all required services wired.
func Serve() *fx.App {
	return fx.New(
		fx.WithLogger(bootstrap.NewFxLogger),
		fx.Module("bootstrap",
			fx.Provide(config.LoadConfig),
			fx.Provide(bootstrap.NewLogger),
			fx.Invoke(bootstrap.StartTracer),
		),

		fx.Module("infra",
			fx.Provide(natsx.BuildNATSOptions), // shared NATS connection (event bus + KV)
			fx.Provide(natsx.NewJetStream),     // event bus JetStream for services
			fx.Provide(kv.NewKVStore),          // simplified KV factory (accepts existing NATS conn)
			fx.Provide(s3client.New),
			fx.Provide(qdrant.NewClient),
			fx.Provide(elasticsearch.NewClient), // Elasticsearch HTTP client
			fx.Provide(whisper.NewClient),       // Faster-Whisper HTTP client
		),

		fx.Module("http",
			fx.Provide(server.NewFiberApp),
			fx.Invoke(server.Start),
		),

		fx.Module("domain",
			fx.Provide(scraper.NewService),
			fx.Provide(transcribe.NewService),
			fx.Provide(enrich.NewService),
		),

		// start background workers if any
		fx.Invoke(func(
			_ *scraper.Service,
			_ *transcribe.Service,
			_ *enrich.Service,
		) {
			// constructors register lifecycle hooks; nothing else to invoke here
		}),
	)
}
