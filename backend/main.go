// A minimal backend that validates uploaded asset CSV files, commits the
// valid rows to PostgreSQL and serves the stored assets.
package main

import (
	"log/slog"
	"net/http"
	"net/url"
	"os"

	"smartgrid/backend/internal/api"
	"smartgrid/backend/internal/store"
)

const (
	addr = ":8080"

	// defaultDatabaseURL is used when DATABASE_URL is not set. It is the
	// development database that Docker Compose runs (see compose.yaml).
	defaultDatabaseURL = "postgres://smartgrid:smartgrid@localhost:5432/smartgrid?sslmode=disable"
)

func main() {
	// Logs go to standard output, one line per event, for whatever runs the
	// server to collect: the terminal in development, Docker in a container.
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = defaultDatabaseURL
	}
	db, err := store.Open(databaseURL)
	if err != nil {
		logger.Error("could not open the database", "error", err)
		os.Exit(1)
	}

	handler := api.NewHandler(db, logger)
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/validate", handler.Validate)
	mux.HandleFunc("POST /api/imports", handler.Import)

	mux.HandleFunc("GET /api/search", handler.Search)
	mux.HandleFunc("GET /api/assets", handler.ListAssets)
	mux.HandleFunc("GET /api/assets/{assetId}", handler.GetAsset)

	logger.Info("starting", "database", withoutPassword(databaseURL))
	logger.Info("listening", "addr", addr)
	err = http.ListenAndServe(addr, api.LogRequests(logger, mux))
	logger.Error("the server stopped", "error", err)
	os.Exit(1)
}

// withoutPassword returns a database URL with its password hidden, so that
// the logs can say which database the server uses without leaking it.
func withoutPassword(databaseURL string) string {
	parsed, err := url.Parse(databaseURL)
	if err != nil {
		// The error would quote the URL, password included.
		return "a URL that could not be read"
	}
	return parsed.Redacted()
}
