package api

import (
	"net/http"
	"time"

	"smartgrid/backend/internal/asset"
	"smartgrid/backend/internal/csvfile"
)

type validateResponse struct {
	Valid     bool                    `json:"valid"`
	TotalRows int                     `json:"total_rows"`
	Errors    []asset.ValidationError `json:"errors"`
	Warnings  []asset.ValidationError `json:"warnings"`
	// ValidFile holds the rows without errors and InvalidFile the rest.
	ValidFile   *csvfile.File `json:"valid_file,omitempty"`
	InvalidFile *csvfile.File `json:"invalid_file,omitempty"`
}

// Validate checks the uploaded rows and returns every error found, along
// with the valid and invalid rows as two CSV files.
func (h *Handler) Validate(w http.ResponseWriter, r *http.Request) {
	upload, ok := readUpload(w, r)
	if !ok {
		return
	}
	header, assets := upload.header, upload.assets

	// Check against the database as well, so that rows which could not be
	// committed are reported now rather than when committing.
	existing, err := h.store.Assets(r.Context())
	if err != nil {
		h.fail(w, r, http.StatusServiceUnavailable, "the database is not available", err)
		return
	}

	errs := asset.ValidateWithStored(assets, existing)
	validFile, invalidFile, err := csvfile.BuildFiles(upload.filename, header, assets, errs, time.Now())
	if err != nil {
		h.fail(w, r, http.StatusInternalServerError, "the result files could not be generated", err)
		return
	}

	warnings := asset.FindWarnings(assets)
	h.logFor(r).Info("validation",
		"file", upload.filename,
		"rows", len(assets),
		"errors", len(errs),
		"warnings", len(warnings))

	writeJSON(w, http.StatusOK, validateResponse{
		Valid:       len(errs) == 0,
		TotalRows:   len(assets),
		Errors:      errs,
		Warnings:    warnings,
		ValidFile:   &validFile,
		InvalidFile: &invalidFile,
	})
}
