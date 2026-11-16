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
func (h *Handler) fail(w http.ResponseWriter, r *http.Request, status int, message string, err error) {
	h.logger.ErrorContext(r.Context(), message,
		"method", r.Method,
		"path", r.URL.Path,
		"error", err)
	writeError(w, status, message+": "+err.Error())
}

// writeError reports that a request failed before producing any result.
func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]any{
		"error":  message,
		"errors": []asset.ValidationError{},
	})
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	// Let the frontend, which runs on a different port, read the response.
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
}
