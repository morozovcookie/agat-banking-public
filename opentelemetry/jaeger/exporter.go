package jaeger

import (
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/trace"
)

// NewExporter returns a new NewExporter instance.
func NewExporter() (trace.SpanExporter, error) {
	exporter, err := jaeger.New(
		jaeger.WithAgentEndpoint(
			jaeger.WithAgentHost("127.0.0.1"),
			jaeger.WithAgentPort("6831")))
	if err != nil {
		return nil, errors.Wrap(err, "init exporter")
	}

	return exporter, nil
}
