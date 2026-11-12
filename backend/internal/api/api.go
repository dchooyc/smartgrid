// Package api holds the HTTP handlers and their request and response types.
package api

import (
	"encoding/json"
	"net/http"

	"smartgrid/backend/internal/asset"
	"smartgrid/backend/internal/store"
)

// Handler serves the API. Its methods are the handlers for each route.
type Handler struct {
	store *store.Store
}

func NewHandler(store *store.Store) *Handler {
	return &Handler{store: store}
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
