package v1

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Handler represents a base type for any HTTP handler.
type Handler struct {
	router chi.Router
}

// NewHandler returns a new Handler instance.
func NewHandler() *Handler {
	h := &Handler{
		router: chi.NewRouter(),
	}

	for _, fn := range []func(http.Handler) http.Handler {
		middleware.RealIP,
		middleware.Recoverer,
		middleware.SetHeader("Content-Type", "application/json"),
		middleware.RequestID,
	} {
		h.router.Use(fn)
	}

	h.router.Mount("/debug", middleware.Profiler())

	return h
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}
