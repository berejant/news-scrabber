package bootstrap

import (
	"context"

	"github.com/DataDog/datadog-go/v5/statsd"
	"go.uber.org/fx"
)

func NewStatsd(lc fx.Lifecycle) (*statsd.Client, error) {
	stat, err := statsd.New("")

	lc.Append(fx.Hook{
		OnStop: func(_ context.Context) error {
			return stat.Close()
		},
	})

	return stat, err
}
