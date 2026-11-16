package api

import (
	"fmt"
	"net/http"

	"smartgrid/backend/internal/asset"
)

type importResponse struct {
	Committed     bool                    `json:"committed"`
	RowsCommitted int                     `json:"rows_committed"`
	Errors        []asset.ValidationError `json:"errors"`
	// Error says why nothing was committed.
	Error string `json:"error,omitempty"`
}

// Import validates the uploaded rows again, inside the transaction that
// saves them, and commits them all or not at all.
func (h *Handler) Import(w http.ResponseWriter, r *http.Request) {
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
	tx, err := h.store.Begin(ctx)
	if err != nil {
		h.fail(w, r, http.StatusServiceUnavailable, "the database is not available", err)
		return
	}
	defer tx.Rollback() // does nothing once the transaction is committed

	existing, err := tx.Assets(ctx)
	if err != nil {
		h.fail(w, r, http.StatusInternalServerError, "the existing assets could not be read", err)
		return
	}

	if errs := asset.ValidateWithStored(assets, existing); len(errs) > 0 {
		writeJSON(w, http.StatusUnprocessableEntity, importResponse{
			Errors: errs,
			Error:  fmt.Sprintf("validation failed with %d errors, so nothing was committed", len(errs)),
		})
		return
	}

	err = tx.InsertAssets(ctx, asset.SortParentsFirst(assets))
	if err == nil {
		err = tx.Commit()
	}
	if err != nil {
		h.fail(w, r, http.StatusInternalServerError, "the rows could not be saved, so nothing was committed", err)
		return
	}

	writeJSON(w, http.StatusOK, importResponse{
		Committed:     true,
		RowsCommitted: len(assets),
		Errors:        []asset.ValidationError{},
	})
}
