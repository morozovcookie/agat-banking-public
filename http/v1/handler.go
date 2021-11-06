package v1

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

const (
	// BasePathPrefix is the path prefix that places at the beginning of any API prefixes.
	BasePathPrefix = "/api/v1"
)

var _ http.Handler = (*Handler)(nil)

// Handler represents a base type for any HTTP handler.
type Handler struct {
	router chi.Router
}

// NewHandler returns a new Handler instance.
func NewHandler() *Handler {
	return &Handler{
		router: chi.NewRouter(),
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	h.router.ServeHTTP(w, r)
}

func encodeResponse(_ context.Context, w http.ResponseWriter, status int, resp interface{}) {
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(resp)
}

func badRequestError(_ context.Context, w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
}
