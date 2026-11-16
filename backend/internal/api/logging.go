package api

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"
)

// requestIDHeader is the response header that carries a request's ID.
const requestIDHeader = "X-Request-ID"

// requestIDKey is where a request's ID is kept in its context.
type requestIDKey struct{}

// newRequestID returns a short random ID, e.g. "a3f9c2e1". It is long enough
// to tell apart the requests around a failure, which is all it is for.
func newRequestID() string {
	id := make([]byte, 4)
	rand.Read(id) // never fails
	return hex.EncodeToString(id)
}

// requestID returns the ID that LogRequests gave the request, or "" if the
// request did not come through it.
func requestID(ctx context.Context) string {
	id, _ := ctx.Value(requestIDKey{}).(string)
	return id
}

// LogRequests wraps a handler so that every request it serves is logged as
// one line, once the response has been written. A response that reports a
// server failure is logged as an error and one that reports a mistake in the
// request as a warning, so that problems stand out from normal traffic.
//
// It also gives each request an ID, which is on every line logged for that
// request and in the response's X-Request-ID header. Searching the logs for
// the ID shows everything that happened in the request.
//
// A handler that panics does not take the server down or leave the caller
// without an answer: the panic is logged and the response is a 500.
func LogRequests(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		id := newRequestID()
		w.Header().Set(requestIDHeader, id)
		r = r.WithContext(context.WithValue(r.Context(), requestIDKey{}, id))
		recorder := &statusRecorder{ResponseWriter: w, status: http.StatusOK}

		serve(logger, next, recorder, r)

		level := slog.LevelInfo
		switch {
		case recorder.status >= 500:
			level = slog.LevelError
		case recorder.status >= 400:
			level = slog.LevelWarn
		}
		logger.Log(r.Context(), level, "request",
			"request_id", id,
			"method", r.Method,
			"path", r.URL.Path,
			"query", r.URL.RawQuery,
			"status", recorder.status,
			"duration_ms", time.Since(start).Milliseconds())
	})
}

// serve runs the handler. If the handler panics, which means the code has a
// bug, serve logs the panic with where it happened and responds with a 500.
func serve(logger *slog.Logger, next http.Handler, w http.ResponseWriter, r *http.Request) {
	defer func() {
		cause := recover()
		if cause == nil {
			return
		}
		id := requestID(r.Context())
		logger.Error("panic",
			"request_id", id,
			"method", r.Method,
			"path", r.URL.Path,
			"error", cause,
			"stack", string(debug.Stack()))
		writeFailure(w, http.StatusInternalServerError, "the server failed unexpectedly", id)
	}()
	next.ServeHTTP(w, r)
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
