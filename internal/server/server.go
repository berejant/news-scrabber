package server

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"news-scrabber/internal/config"
	"news-scrabber/internal/server/actions/transribe"
)

// NewFiberApp constructs a Fiber application and registers routes.
func NewFiberApp(cfg *config.Config, log *zap.Logger, act *transribe.RequestTranscribeAction) *fiber.App {
	app := fiber.New()

	// Register all routes
	RegisterRoutes(app, act)

	log.Info("fiber app initialized")
	return app
}

// Start attaches lifecycle hooks to start and stop the Fiber HTTP server.
func Start(lc fx.Lifecycle, app *fiber.App, cfg *config.Config, log *zap.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			host := cfg.Server.Host
			if host == "" {
				host = "0.0.0.0"
			}
			addr := fmt.Sprintf("%s:%d", host, cfg.Server.Port)
			log.Info("starting HTTP server", zap.String("addr", addr))
			go func() {
				if err := app.Listen(addr); err != nil {
					log.Error("fiber server stopped with error", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info("shutting down HTTP server")
			return app.Shutdown()
		},
	})
}
