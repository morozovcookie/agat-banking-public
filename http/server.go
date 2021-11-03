package http

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/pkg/errors"
)

// A Server defines parameters for running an HTTP server.
type Server struct {
	addr string

	router chi.Router
	server *http.Server
}

// NewServer returns a new Server instance.
func NewServer(addr string, opts ...ServerOption) *Server {
	srv := &Server{
		addr:   addr,
		router: chi.NewRouter(),

		server: nil,
	}

	srv.router.Use(middleware.RealIP, middleware.RequestID)
	srv.router.Mount("/debug", middleware.Profiler())

	srv.server = &http.Server{
		Addr:              addr,
		Handler:           srv.router,
		TLSConfig:         nil,
		ReadTimeout:       DefaultReadTimeout,
		ReadHeaderTimeout: DefaultReadHeaderTimeout,
		WriteTimeout:      DefaultWriteTimeout,
		IdleTimeout:       DefaultIdleTimeout,
		MaxHeaderBytes:    DefaultMaxHeaderBytes,
		TLSNextProto:      nil,
		ConnState:         nil,
		ErrorLog:          nil,
		BaseContext:       nil,
		ConnContext:       nil,
	}

	for _, opt := range opts {
		opt.apply(srv)
	}

	return srv
}

// Start listens on the TCP network address srv.Addr and then calls Serve to handle requests on incoming connections.
// Accepted connections are configured to enable TCP keep-alives.
func (srv *Server) Start() error {
	if err := srv.server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		return errors.Wrap(err, "start http server")
	}

	return nil
}

// Shutdown gracefully shuts down the server without interrupting any active connections.
func (srv *Server) Shutdown(ctx context.Context) error {
	if err := srv.server.Shutdown(ctx); err != nil {
		return errors.Wrap(err, "shutdown http server")
	}

	return nil
}
