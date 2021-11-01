package prometheus

import (
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/histogram"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	selector "go.opentelemetry.io/otel/sdk/metric/selector/simple"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
)

// NewExporter configure and return a new prometheus.Exporter instance.
func NewExporter() (*prometheus.Exporter, error) {
	var (
		config = prometheus.Config{}
		ctrl   = controller.New(
			processor.NewFactory(
				selector.NewWithHistogramDistribution(
					histogram.WithExplicitBoundaries(config.DefaultHistogramBoundaries),
				),
				metric.CumulativeExportKindSelector(),
				processor.WithMemory(true),
			),
			controller.WithResource(
				resource.NewWithAttributes(
					semconv.SchemaURL,
					semconv.ServiceNameKey.String("opentelemetry-prom"),
					semconv.ServiceVersionKey.String("1.0.0"),
					semconv.DeploymentEnvironmentKey.String("production"),
				),
			),
		)
	)

	exporter, err := prometheus.New(config, ctrl)
	if err != nil {
		return nil, errors.Wrap(err, "init prometheus exporter")
	}

	return exporter, nil
}
