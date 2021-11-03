package prometheus

import (
	"bytes"
	"net/http"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/unit"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
)

var _ http.Handler = (*HTTPHandler)(nil)

// A HTTPHandler responds to an HTTP request.
type HTTPHandler struct {
	wrapped http.Handler

	meter metric.Meter
	attrs []attribute.KeyValue

	requestCount    metric.Int64Counter
	requestDuration metric.Int64Histogram
}

// NewHTTPHandler returns a new HTTPHandler instance.
func NewHTTPHandler(
	handler http.Handler,
	meter metric.Meter,
	attrs ...attribute.KeyValue,
) (
	wrapper *HTTPHandler,
	err error,
) {
	wrapper = &HTTPHandler{
		wrapped: handler,

		meter: meter,
		attrs: attrs,

		requestCount:    metric.Int64Counter{},
		requestDuration: metric.Int64Histogram{},
	}

	wrapper.requestCount, err = meter.NewInt64Counter("http.server.active_requests",
		metric.WithDescription("measures the number of concurrent HTTP requests that are currently in-flight"),
		metric.WithUnit(unit.Dimensionless))
	if err != nil {
		return nil, errors.Wrap(err, "init handler")
	}

	wrapper.requestDuration, err = meter.NewInt64Histogram("http.server.duration",
		metric.WithDescription("measures the duration of the inbound HTTP request"),
		metric.WithUnit(unit.Milliseconds))
	if err != nil {
		return nil, errors.Wrap(err, "init handler")
	}

	return wrapper, nil
}

func (h *HTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		now   = time.Now()
		attrs = append(h.attrs, semconv.HTTPMethodKey.String(r.Method))

		wrappedResponseWriter = NewResponseWriter(w)
	)

	if ua := r.UserAgent(); ua != "" {
		attrs = append(attrs, semconv.HTTPUserAgentKey.String(ua))
	}

	if r.Host != "" {
		attrs = append(attrs, semconv.HTTPHostKey.String(r.Host))
	}

	schema := semconv.HTTPSchemeHTTP
	if r.TLS != nil {
		schema = semconv.HTTPSchemeHTTPS
	}

	attrs = append(attrs, schema)

	userInfo := r.URL.User
	r.URL.User = nil

	attrs = append(attrs, semconv.HTTPURLKey.String(r.URL.String()), semconv.HTTPTargetKey.String(r.URL.Path))

	r.URL.User = userInfo

	flavor := new(bytes.Buffer)
	_, _ = flavor.WriteString(strconv.Itoa(r.ProtoMajor))

	if r.ProtoMajor == 1 {
		_, _ = flavor.WriteRune('.')
		_, _ = flavor.WriteString(strconv.Itoa(r.ProtoMinor))
	}

	if val := flavor.String(); val != "" {
		attrs = append(attrs, semconv.HTTPFlavorKey.String(val))
	}

	h.wrapped.ServeHTTP(wrappedResponseWriter, r)

	attrs = append(attrs, semconv.HTTPStatusCodeKey.Int(wrappedResponseWriter.StatusCode()))

	ctx := r.Context()

	h.requestCount.Add(ctx, 1, attrs...)
	h.requestDuration.Record(ctx, time.Since(now).Milliseconds(), attrs...)
}
