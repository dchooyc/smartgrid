package api

import (
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"smartgrid/backend/internal/asset"
)

const (
	// maxPageSize is the most assets one request can ask for with "limit".
	maxPageSize = 100

	// searchLimit is the most assets one search returns.
	searchLimit = 20

	// cacheAssets lets the browser reuse an asset response for 10 seconds
	// without asking again: long enough to absorb repeated refreshes, short
	// enough that data changed elsewhere shows up almost at once. The
	// frontend changes its request URLs when it commits new rows, so its own
	// imports show up immediately.
	cacheAssets = "private, max-age=10"
)

type assetListResponse struct {
	Assets []asset.Stored `json:"assets"`
	// Counts covers the whole list, not just the assets returned.
	Counts asset.Counts `json:"counts"`
	// NextAfter is the "after" value that fetches the rest of the list, or
	// null when this response reaches its end.
	NextAfter *string `json:"next_after"`
}

type assetDetailsResponse struct {
	Asset    asset.Stored   `json:"asset"`
	Parent   *asset.Stored  `json:"parent"`   // null for substations
	Children []asset.Stored `json:"children"` // immediate children only
	// AncestorIDs runs from the asset's substation down to its parent.
	AncestorIDs []string `json:"ancestor_ids"`
	// DescendantCounts counts the assets beneath it at any depth, by type.
	DescendantCounts map[string]int `json:"descendant_counts"`
}

type searchResponse struct {
	Assets []asset.Stored `json:"assets"`
	// More says that more assets match than were returned.
	More bool `json:"more"`
}

// Search returns the first assets whose ID or name contains the text in the
// "q" query parameter, ignoring case.
func (h *Handler) Search(w http.ResponseWriter, r *http.Request) {
	term := strings.TrimSpace(r.URL.Query().Get("q"))
	if term == "" {
		writeError(w, http.StatusBadRequest, "q is required: the text to search for")
		return
	}

	// One asset more than the limit is fetched, to learn whether more match.
	assets, err := h.store.SearchAssets(r.Context(), term, searchLimit+1)
	if err != nil {
		writeError(w, http.StatusServiceUnavailable, "the database is not available: "+err.Error())
		return
	}

	response := searchResponse{Assets: assets}
	if len(assets) > searchLimit {
		response.Assets = assets[:searchLimit]
		response.More = true
	}
	w.Header().Set("Cache-Control", cacheAssets)
	writeJSON(w, http.StatusOK, response)
}

// ListAssets returns the stored assets, or only those of the type given in
// the optional "type" query parameter.
//
// The list can be read a page at a time: "limit" caps how many assets are
// returned and "after" continues from the asset ID given in the previous
// response's next_after. Without a limit, every asset is returned.
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

	query := r.URL.Query()
	limit := 0
	if value := query.Get("limit"); value != "" {
		var err error
		if limit, err = strconv.Atoi(value); err != nil || limit < 1 || limit > maxPageSize {
			writeError(w, http.StatusBadRequest, fmt.Sprintf("limit must be a whole number from 1 to %d", maxPageSize))
			return
		}
	}

	// One asset more than the limit is fetched, to learn whether more follow.
	fetch := limit
	if limit > 0 {
		fetch = limit + 1
	}
	assets, err := h.store.ListAssets(r.Context(), assetType, query.Get("after"), fetch)
	var counts asset.Counts
	if err == nil {
		counts, err = h.store.CountAssets(r.Context(), assetType)
	}
	if err != nil {
		writeError(w, http.StatusServiceUnavailable, "the database is not available: "+err.Error())
		return
	}

	response := assetListResponse{Assets: assets, Counts: counts}
	if limit > 0 && len(assets) > limit {
		response.Assets = assets[:limit]
		response.NextAfter = &assets[limit-1].ID
	}
	w.Header().Set("Cache-Control", cacheAssets)
	writeJSON(w, http.StatusOK, response)
}

// GetAsset returns one asset with its direct parent and immediate children,
// and where it sits in the hierarchy.
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
	var ancestorIDs []string
	if err == nil {
		ancestorIDs, err = h.store.ListAncestorIDs(ctx, id)
	}
	var descendantCounts map[string]int
	if err == nil {
		descendantCounts, err = h.store.CountDescendants(ctx, id)
	}
	if err != nil {
		writeError(w, http.StatusServiceUnavailable, "the database is not available: "+err.Error())
		return
	}

	w.Header().Set("Cache-Control", cacheAssets)
	writeJSON(w, http.StatusOK, assetDetailsResponse{
		Asset:            *found,
		Parent:           parent,
		Children:         children,
		AncestorIDs:      ancestorIDs,
		DescendantCounts: descendantCounts,
	})
}
