// A minimal backend that validates uploaded asset CSV files and commits the
// valid rows to PostgreSQL.
package main

import (
	"log"
	"net/http"
	"os"

	"smartgrid/backend/internal/api"
	"smartgrid/backend/internal/store"
)

const (
	addr = ":8080"

	// defaultDatabaseURL is used when DATABASE_URL is not set.
	defaultDatabaseURL = "postgres://localhost:5432/smartgrid?sslmode=disable"
)

func main() {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = defaultDatabaseURL
	}
	db, err := store.Open(databaseURL)
	if err != nil {
		log.Fatal(err)
	}

	handler := api.NewHandler(db)
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/validate", handler.Validate)
	mux.HandleFunc("POST /api/imports", handler.Import)

	log.Printf("listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
