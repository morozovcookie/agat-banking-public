package jaeger

import (
	"net/http"

	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace"
)

var _ http.Handler = (*HTTPHandler)(nil)

// A HTTPHandler responds to an HTTP request.
type HTTPHandler struct {
	wrapped http.Handler

	tracer trace.Tracer
	attrs  []attribute.KeyValue
}

// NewHTTPHandler returns a new HTTPHandler instance.
func NewHTTPHandler(handler http.Handler, tracer trace.Tracer, attrs ...attribute.KeyValue) *HTTPHandler {
	return &HTTPHandler{
		wrapped: handler,

		tracer: tracer,
		attrs:  attrs,
	}
}

func (h *HTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		ctx   = r.Context()
		attrs = append(h.attrs, semconv.HTTPClientAttributesFromHTTPRequest(r)...)
	)

	ctx, span := h.tracer.Start(ctx, "Handler.ServeHTTP", trace.WithAttributes(attrs...))
	defer span.End()

	wrappedResponseWriter := NewResponseWriter(w)

	h.wrapped.ServeHTTP(wrappedResponseWriter, r.WithContext(ctx))

	statusCode := wrappedResponseWriter.StatusCode()

	span.SetAttributes(semconv.HTTPAttributesFromHTTPStatusCode(statusCode)...)
	span.SetStatus(semconv.SpanStatusFromHTTPStatusCode(statusCode))
}
