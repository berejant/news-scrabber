package bootstrap

import (
	"context"
	"news-scrabber/internal/config"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(lc fx.Lifecycle, cfg *config.Config) (*zap.Logger, error) {
	var loggerConfig zap.Config
	if cfg.IsDevelopment() || cfg.IsLocal() || cfg.IsTest() {
		loggerConfig = zap.NewDevelopmentConfig()
	} else {
		loggerConfig = zap.NewProductionConfig()

		if cfg.IsStaging() {
			loggerConfig.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
		}
	}
	loggerConfig.EncoderConfig.TimeKey = "timestamp"

	logger, err := loggerConfig.Build()
	lc.Append(fx.Hook{
		OnStop: func(_ context.Context) error {
			_ = logger.Sync()
			return nil
		},
	})

	return logger, err
}

func NewFxLogger(log *zap.Logger, cfg *config.Config) fxevent.Logger {
	zapLogger := &fxevent.ZapLogger{
		Logger: log.With(
			zap.String("component", "fx"),
		),
	}
	if cfg.IsTest() {
		zapLogger.UseLogLevel(zapcore.DebugLevel)
	}

	return zapLogger
}
