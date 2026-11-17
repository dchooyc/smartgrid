// Package api holds the HTTP handlers and their request and response types.
package api

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"smartgrid/backend/internal/asset"
	"smartgrid/backend/internal/store"
)

// Handler serves the API. Its methods are the handlers for each route.
type Handler struct {
	store  *store.Store
	logger *slog.Logger
}

func NewHandler(store *store.Store, logger *slog.Logger) *Handler {
	return &Handler{store: store, logger: logger}
}

// fail reports that the server could not do what a request asked, and logs
// why. The message says what could not be done and err is the cause. Use it
// for failures on the server's side, not for mistakes in the request.
//
// The response carries the request's ID, so that whoever sees the failure
// can quote it and whoever investigates can find it in the logs.
func (h *Handler) fail(w http.ResponseWriter, r *http.Request, status int, message string, err error) {
	h.logFor(r).Error(message,
		"method", r.Method,
		"path", r.URL.Path,
		"error", err)
	writeFailure(w, status, message+": "+err.Error(), requestID(r.Context()))
}

// logFor returns the logger for a request. Every line it writes carries the
// request's ID.
func (h *Handler) logFor(r *http.Request) *slog.Logger {
	return h.logger.With("request_id", requestID(r.Context()))
}

// writeFailure reports a failure on the server's side, with the ID of the
// request it happened in.
func writeFailure(w http.ResponseWriter, status int, message, requestID string) {
	writeJSON(w, status, map[string]any{
		"error":      message,
		"errors":     []asset.ValidationError{},
		"request_id": requestID,
	})
}

// writeError reports that a request failed before producing any result.
func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]any{
		"error":  message,
		"errors": []asset.ValidationError{},
	})
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
}

// AllowOrigin wraps a handler so that pages served from one other address,
// the frontend's, may read its responses. Browsers refuse to let a page read
// a response from a different address unless the response names that page's
// address, so pages on any other website are refused.
//
// Before a request that changes data in a way a plain form cannot, such as
// DELETE, a browser first asks with OPTIONS whether it is allowed. That
// question is answered here and never reaches the handlers.
func AllowOrigin(origin string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE")
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
