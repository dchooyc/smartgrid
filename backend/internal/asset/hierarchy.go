package asset

import "slices"

// SortParentsFirst orders assets so that every parent comes before its
// children, otherwise keeping the original order. The assets must already be
// validated, so they contain no cycles.
func SortParentsFirst(assets []Asset) []Asset {
	byID := assetsByID(assets)

	// depth counts how many ancestors of an asset are in this batch.
	depth := func(a Asset) int {
		levels := 0
		for parent, ok := byID[a.ParentID]; ok && levels < len(assets); parent, ok = byID[parent.ParentID] {
			levels++
		}
		return levels
	}

	sorted := slices.Clone(assets)
	slices.SortStableFunc(sorted, func(a, b Asset) int {
		return depth(a) - depth(b)
	})
	return sorted
}
