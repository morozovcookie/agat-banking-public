package zap

import (
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

var _ http.Handler = (*HTTPHandler)(nil)

// A HTTPHandler responds to an HTTP request.
type HTTPHandler struct {
	wrapped http.Handler
	logger  *zap.Logger
}

// NewHTTPHandler returns a new HTTPHandler instance.
func NewHTTPHandler(handler http.Handler, logger *zap.Logger) *HTTPHandler {
	return &HTTPHandler{
		wrapped: handler,
		logger:  logger,
	}
}

func (h *HTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var handler http.Handler = h.wrapped
	{
		handler = logRequest(h.logger, handler)
		handler = recoverRequest(h.logger, handler)
	}

	handler.ServeHTTP(w, r)
}

func recoverRequest(logger *zap.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err == nil {
				return
			}

			if err, ok := err.(error); ok && errors.Is(err, http.ErrAbortHandler) {
				return
			}

			var brokenPipe bool
			if ne, ok := err.(*net.OpError); ok {
				if se, ok := ne.Err.(*os.SyscallError); ok {
					if strings.Contains(strings.ToLower(se.Error()), "broken pipe") ||
						strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
						brokenPipe = true
					}
				}
			}

			httpRequest, _ := httputil.DumpRequest(r, false)
			if brokenPipe {
				logger.Error(r.URL.Path,
					zap.Any("error", err),
					zap.String("request", string(httpRequest)))

				return
			}

			logger.Error("[Recovery from panic]",
				zap.Time("time", time.Now().UTC()),
				zap.Any("error", err),
				zap.String("request", string(httpRequest)),
				zap.String("stack", string(debug.Stack())))

			w.WriteHeader(http.StatusInternalServerError)
		}()

		next.ServeHTTP(w, r)
	})
}

func logRequest(logger *zap.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			start = time.Now().UTC()
			path  = r.URL.Path
			query = r.URL.RawQuery

			wrappedResponseWriter = NewResponseWriter(w)
		)

		next.ServeHTTP(wrappedResponseWriter, r)

		end := time.Now().UTC()

		logger.Info(path,
			zap.Int("status", wrappedResponseWriter.StatusCode()),
			zap.String("method", r.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", realIP(r)),
			zap.String("user-agent", r.UserAgent()),
			zap.Time("time", end),
			zap.Duration("latency", end.Sub(start)),
		)
	})
}

func realIP(r *http.Request) string {
	if xrip := r.Header.Get("X-Real-IP"); xrip != "" {
		return xrip
	}

	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		if i := strings.Index(xff, ", "); i != -1 {
			return xff[:i]
		}

		return xff
	}

	return ""
}
