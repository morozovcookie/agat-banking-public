package jaeger

import (
	"net/http"
)

var _ http.ResponseWriter = (*ResponseWriter)(nil)

// A ResponseWriter interface is used by an HTTP handler to construct an HTTP response.
type ResponseWriter struct {
	wrapped http.ResponseWriter

	statusCode int
}

// NewResponseWriter returns a new ResponseWriter instance.
func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{
		wrapped:    w,
		statusCode: http.StatusOK,
	}
}

// Header returns the header map that will be sent by WriteHeader. The Header map also is the mechanism with which
// Handlers can set HTTP trailers.
func (r *ResponseWriter) Header() http.Header {
	return r.wrapped.Header()
}

// Write writes the data to the connection as part of an HTTP reply.
func (r *ResponseWriter) Write(bytes []byte) (int, error) {
	return r.wrapped.Write(bytes)
}

// WriteHeader sends an HTTP response header with the provided status code.
func (r *ResponseWriter) WriteHeader(statusCode int) {
	r.statusCode = statusCode

	r.wrapped.WriteHeader(statusCode)
}

func (r *ResponseWriter) StatusCode() int {
	return r.statusCode
}
