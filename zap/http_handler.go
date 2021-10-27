package zap

import (
	"go.uber.org/zap"
	"net/http"
)

// HTTPHandler represents a wrapper on http.Handler which could write logs.
type HTTPHandler struct {
	wrapped http.Handler
}

// NewHTTPHandler returns a new HTTPHandler instance.
func NewHTTPHandler(handler http.Handler, logger *zap.Logger) *HTTPHandler {
	return &HTTPHandler{}
}

func (h *HTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.wrapped.ServeHTTP(w, r)
}
