package jaeger

import (
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
)

// NewProvider configure and returns a new trace.TracerProvider instance.
func NewProvider(exporter tracesdk.SpanExporter) *tracesdk.TracerProvider {
	return tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exporter),
		tracesdk.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String("opentelemetry-jaeger"),
				semconv.ServiceVersionKey.String("1.0.0"),
				semconv.DeploymentEnvironmentKey.String("production"))))
}
