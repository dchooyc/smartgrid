package main

import (
	"fmt"
	"net/http"
	"slices"
)

type importResponse struct {
	Committed     bool              `json:"committed"`
	RowsCommitted int               `json:"rows_committed"`
	Errors        []validationError `json:"errors"`
	// Error says why nothing was committed.
	Error string `json:"error,omitempty"`
}

// handleImport validates the uploaded rows again, this time against the
// database too, and commits them all in one transaction or not at all.
func handleImport(w http.ResponseWriter, r *http.Request) {
	upload, ok := readUpload(w, r)
	if !ok {
		return
	}
	assets := upload.assets
	if len(assets) == 0 {
		writeError(w, http.StatusBadRequest, "the file has no rows to commit")
		return
	}

	ctx := r.Context()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		writeError(w, http.StatusServiceUnavailable, "the database is not available: "+err.Error())
		return
	}
	defer tx.Rollback() // does nothing once the transaction is committed

	existing, err := loadExistingAssets(ctx, tx)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "the existing assets could not be read: "+err.Error())
		return
	}

	if errs := validateWithStored(assets, existing); len(errs) > 0 {
		writeJSON(w, http.StatusUnprocessableEntity, importResponse{
			Errors: errs,
			Error:  fmt.Sprintf("validation failed with %d errors, so nothing was committed", len(errs)),
		})
		return
	}

	err = insertAssets(ctx, tx, sortParentsFirst(assets))
	if err == nil {
		err = tx.Commit()
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "the rows could not be saved, so nothing was committed: "+err.Error())
		return
	}

	writeJSON(w, http.StatusOK, importResponse{
		Committed:     true,
		RowsCommitted: len(assets),
		Errors:        []validationError{},
	})
}

// validateWithStored runs every validation on the incoming rows, treating the
// assets already in the database as part of the hierarchy, and rejects rows
// whose ID is already stored. Both validating and importing use it.
func validateWithStored(incoming, existing []asset) []validationError {
	incomingIDs := assetsByID(incoming)

	// Stored assets take part in validation so that parents can be found in
	// the database. Those sharing an ID with an incoming row are left out and
	// reported by checkAlreadyStored instead.
	all := slices.Clone(incoming)
	for _, a := range existing {
		if _, clash := incomingIDs[a.id]; !clash {
			all = append(all, a)
		}
	}

	errs := []validationError{}
	for _, e := range validateAssets(all) {
		if e.Line != 0 { // line 0 is a stored asset, not a row of the upload
			errs = append(errs, e)
		}
	}
	return append(errs, withType("Already in database", checkAlreadyStored(incoming, existing))...)
}

// checkAlreadyStored reports every incoming row whose asset_id is already in
// the database.
func checkAlreadyStored(incoming, existing []asset) []validationError {
	stored := assetsByID(existing)

	var errs []validationError
	for _, a := range incoming {
		if _, exists := stored[a.id]; exists {
			errs = append(errs, newError(a, "Already Exists", "already exists in the database."))
		}
	}
	return errs
}

// sortParentsFirst orders assets so that every parent comes before its
// children, otherwise keeping the original order. The assets must already be
// validated, so they contain no cycles.
func sortParentsFirst(assets []asset) []asset {
	byID := assetsByID(assets)

	// depth counts how many ancestors of an asset are in this batch.
	depth := func(a asset) int {
		levels := 0
		for parent, ok := byID[a.parentID]; ok && levels < len(assets); parent, ok = byID[parent.parentID] {
			levels++
		}
		return levels
	}

	sorted := slices.Clone(assets)
	slices.SortStableFunc(sorted, func(a, b asset) int {
		return depth(a) - depth(b)
	})
	return sorted
}
