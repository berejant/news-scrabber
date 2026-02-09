package bootstrap

import (
	"context"

	"github.com/DataDog/dd-trace-go/v2/ddtrace/tracer"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type ddLogger struct {
	logger *zap.Logger
}

func (d *ddLogger) Log(msg string) {
	d.logger.Info(msg)
}

func StartTracer(lc fx.Lifecycle, logger *zap.Logger) error {
	lc.Append(fx.Hook{
		OnStop: func(_ context.Context) error {
			tracer.Stop()
			return nil
		},
	})

	return tracer.Start(
		tracer.WithLogger(
			&ddLogger{
				logger: logger.With(
					zap.String("component", "tracer"),
				),
			},
		),
	)
}
