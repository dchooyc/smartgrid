package api

import (
	"fmt"
	"net/http"
	"slices"
	"strings"

	"smartgrid/backend/internal/asset"
)

type assetListResponse struct {
	Assets []asset.Stored `json:"assets"`
}

type assetDetailsResponse struct {
	Asset    asset.Stored   `json:"asset"`
	Parent   *asset.Stored  `json:"parent"`   // null for substations
	Children []asset.Stored `json:"children"` // immediate children only
}

// ListAssets returns every stored asset, or only those of the type given in
// the optional "type" query parameter.
func (h *Handler) ListAssets(w http.ResponseWriter, r *http.Request) {
	h.listAssets(w, r, r.URL.Query().Get("type"))
}

// ListAssetsOfType returns a handler listing the stored assets of one type,
// for routes such as /api/substations.
func (h *Handler) ListAssetsOfType(assetType string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.listAssets(w, r, assetType)
	}
}

// listAssets serves both styles of type filtering. An empty type lists all
// assets.
func (h *Handler) listAssets(w http.ResponseWriter, r *http.Request, assetType string) {
	supported := asset.SupportedTypes()
	if assetType != "" && !slices.Contains(supported, assetType) {
		writeError(w, http.StatusBadRequest, fmt.Sprintf(
			"unsupported asset type %q (supported types: %s)", assetType, strings.Join(supported, ", ")))
		return
	}

	assets, err := h.store.ListAssets(r.Context(), assetType)
	if err != nil {
		writeError(w, http.StatusServiceUnavailable, "the database is not available: "+err.Error())
		return
	}
	writeJSON(w, http.StatusOK, assetListResponse{Assets: assets})
}

// GetAsset returns one asset with its direct parent and immediate children.
func (h *Handler) GetAsset(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.PathValue("assetId")

	found, err := h.store.GetAsset(ctx, id)
	if err == nil && found == nil {
		writeError(w, http.StatusNotFound, fmt.Sprintf("asset %q was not found", id))
		return
	}

	var parent *asset.Stored
	if err == nil && found.ParentID != nil {
		parent, err = h.store.GetAsset(ctx, *found.ParentID)
	}
	var children []asset.Stored
	if err == nil {
		children, err = h.store.ListChildren(ctx, id)
	}
	if err != nil {
		writeError(w, http.StatusServiceUnavailable, "the database is not available: "+err.Error())
		return
	}

	writeJSON(w, http.StatusOK, assetDetailsResponse{Asset: *found, Parent: parent, Children: children})
}
