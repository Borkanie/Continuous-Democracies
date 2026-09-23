// Package httplog provides a chi middleware that logs the arrival and
// completion of every HTTP request, giving controller-layer trace coverage
// for all generated routes without any per-handler code.
package httplog

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
)

// TraceIDMiddleware assigns a UUID trace ID to every request, stored under
// chi's own middleware.RequestIDKey so middleware.GetReqID(ctx) keeps
// working everywhere unchanged. This replaces middleware.RequestID entirely
// rather than wrapping it: chi's built-in generator produces a
// hostname/counter-style ID with no way to override its format in this
// version, and a UUID is what a trace ID should look like. As with chi's
// own RequestID, an incoming X-Request-Id header is honored if present, so
// an upstream gateway's trace ID propagates instead of being replaced.
func TraceIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		traceID := request.Header.Get(middleware.RequestIDHeader)
		if traceID == "" {
			traceID = uuid.NewString()
		}
		ctx := context.WithValue(request.Context(), middleware.RequestIDKey, traceID)
		next.ServeHTTP(responseWriter, request.WithContext(ctx))
	})
}

// Middleware logs one debug-level line when a request arrives and another
// when it completes (with status code and duration).
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		start := time.Now()
		slog.DebugContext(request.Context(), "http request received",
			"method", request.Method,
			"path", request.URL.Path,
		)

		wrapped := middleware.NewWrapResponseWriter(responseWriter, request.ProtoMajor)
		next.ServeHTTP(wrapped, request)

		slog.DebugContext(request.Context(), "http request completed",
			"method", request.Method,
			"path", request.URL.Path,
			"status", wrapped.Status(),
			"durationMs", time.Since(start).Milliseconds(),
		)
	})
}
