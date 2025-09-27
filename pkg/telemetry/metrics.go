package telemetry

import (
	"shantaram/pkg/config"

	otelmetric "go.opentelemetry.io/otel/metric"
)

type Metrics struct{}

func NewMetrics(_ *config.Config, _ otelmetric.Meter) (*Metrics, error) {
	return &Metrics{}, nil
}
