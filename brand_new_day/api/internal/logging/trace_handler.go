// Package logging provides a slog.Handler wrapper that auto-attaches a
// trace ID and calling-function name to every log line, without requiring
// any change to existing slog call sites.
package logging

import (
	"context"
	"log/slog"
	"runtime"
	"strings"

	"github.com/go-chi/chi/v5/middleware"
)

// modulePrefix is trimmed from fully-qualified function names so log lines
// read e.g. "service.(*VotingService).GetVotingRoundByID" instead of the
// full github.com/borkanie/brand-new-day-api/internal/... symbol.
const modulePrefix = "github.com/borkanie/brand-new-day-api/internal/"

// TraceHandler wraps another slog.Handler, auto-attaching:
//   - "traceId": read from chi's request-ID context value, already
//     populated by middleware.RequestID for every request.
//   - "function": derived from the log call site's program counter, which
//     every slog.XContext call already captures on the record (that's what
//     HandlerOptions.AddSource normally reads) - so this comes for free
//     with no per-call boilerplate.
type TraceHandler struct {
	next slog.Handler
}

// NewTraceHandler wraps next.
func NewTraceHandler(next slog.Handler) *TraceHandler {
	return &TraceHandler{next: next}
}

func (handler *TraceHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return handler.next.Enabled(ctx, level)
}

func (handler *TraceHandler) Handle(ctx context.Context, record slog.Record) error {
	if traceID := middleware.GetReqID(ctx); traceID != "" {
		record.AddAttrs(slog.String("traceId", traceID))
	}
	if record.PC != 0 {
		frames := runtime.CallersFrames([]uintptr{record.PC})
		frame, _ := frames.Next()
		record.AddAttrs(slog.String("function", shortFunctionName(frame.Function)))
	}
	return handler.next.Handle(ctx, record)
}

func (handler *TraceHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &TraceHandler{next: handler.next.WithAttrs(attrs)}
}

func (handler *TraceHandler) WithGroup(name string) slog.Handler {
	return &TraceHandler{next: handler.next.WithGroup(name)}
}

// shortFunctionName trims the module path prefix from a fully-qualified
// function name for readability.
func shortFunctionName(fullName string) string {
	return strings.TrimPrefix(fullName, modulePrefix)
}
