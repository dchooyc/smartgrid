// A minimal backend that validates uploaded asset CSV files, commits the
// valid rows to PostgreSQL and serves the stored assets.
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

	// defaultDatabaseURL is used when DATABASE_URL is not set. It is the
	// development database that Docker Compose runs (see compose.yaml).
	defaultDatabaseURL = "postgres://smartgrid:smartgrid@localhost:5432/smartgrid?sslmode=disable"
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

	mux.HandleFunc("GET /api/assets", handler.ListAssets)
	mux.HandleFunc("GET /api/assets/{assetId}", handler.GetAsset)
	mux.HandleFunc("GET /api/substations", handler.ListAssetsOfType("SUBSTATION"))
	mux.HandleFunc("GET /api/transformers", handler.ListAssetsOfType("TRANSFORMER"))
	mux.HandleFunc("GET /api/lv-boards", handler.ListAssetsOfType("LV_BOARD"))
	mux.HandleFunc("GET /api/switchboards", handler.ListAssetsOfType("SWITCHBOARD"))
	mux.HandleFunc("GET /api/switchboard-panels", handler.ListAssetsOfType("SWITCHBOARD_PANEL"))

	log.Printf("listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
