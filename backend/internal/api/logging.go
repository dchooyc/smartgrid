package api

import (
	"log/slog"
	"net/http"
	"time"
)

// LogRequests wraps a handler so that every request it serves is logged as
// one line, once the response has been written. A response that reports a
// server failure is logged as an error and one that reports a mistake in the
// request as a warning, so that problems stand out from normal traffic.
func LogRequests(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		recorder := &statusRecorder{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(recorder, r)

		level := slog.LevelInfo
		switch {
		case recorder.status >= 500:
			level = slog.LevelError
		case recorder.status >= 400:
			level = slog.LevelWarn
		}
		logger.Log(r.Context(), level, "request",
			"method", r.Method,
			"path", r.URL.Path,
			"query", r.URL.RawQuery,
			"status", recorder.status,
			"duration_ms", time.Since(start).Milliseconds())
	})
}

// statusRecorder remembers the status code a handler responds with, which
// http.ResponseWriter does not otherwise reveal.
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}
