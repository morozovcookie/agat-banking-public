package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"time"

	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace"
)

func main() {
	exporter, err := jaeger.New(
		jaeger.WithAgentEndpoint(
			jaeger.WithAgentHost("127.0.0.1"),
			jaeger.WithAgentPort("6831")))
	if err != nil {
		log.Fatalln(err)
	}

	provider := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exporter),
		tracesdk.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String("opentelemetry-jaeger"),
				semconv.ServiceVersionKey.String("1.0.0"),
				semconv.DeploymentEnvironmentKey.String("production"))))

	tracer := provider.Tracer("main")

	http.HandleFunc("/process", func(w http.ResponseWriter, r *http.Request) {
		ctx, span := tracer.Start(r.Context(), "handleProcess",
			trace.WithAttributes(semconv.HTTPServerAttributesFromHTTPRequest("", "/process", r)...))
		defer span.End()

		_ = r.WithContext(ctx)

		time.Sleep(time.Second)

		w.WriteHeader(http.StatusOK)
		_, _ = io.Copy(w, bytes.NewBufferString("processed"))

		span.SetStatus(semconv.SpanStatusFromHTTPStatusCode(http.StatusOK))
	})

	log.Fatalln(http.ListenAndServe(":2222", nil))
}
