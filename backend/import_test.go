package main

import (
	"slices"
	"testing"
)

func TestValidateForImport(t *testing.T) {
	// Assets already in the database have no line number.
	stored := []asset{
		{id: "SS-1", assetType: "SUBSTATION", name: "North", status: "IN_SERVICE"},
		{id: "SWB-1", parentID: "SS-1", assetType: "SWITCHBOARD", name: "Board", status: "IN_SERVICE"},
	}

	cases := []struct {
		name     string
		incoming []asset
		existing []asset
		want     []validationError
	}{
		{
			name: "complete hierarchy into an empty database",
			incoming: []asset{
				{line: 2, id: "SS-2", assetType: "SUBSTATION", name: "South", status: "IN_SERVICE"},
				{line: 3, id: "TX-2", parentID: "SS-2", assetType: "TRANSFORMER", name: "T", status: "IN_SERVICE"},
			},
			want: []validationError{},
		},
		{
			name: "parent found in the database",
			incoming: []asset{
				{line: 2, id: "PNL-1", parentID: "SWB-1", assetType: "SWITCHBOARD_PANEL", name: "Panel", status: "IN_SERVICE"},
			},
			existing: stored,
			want:     []validationError{},
		},
		{
			name: "parent in neither the upload nor the database",
			incoming: []asset{
				{line: 2, id: "PNL-1", parentID: "SWB-404", assetType: "SWITCHBOARD_PANEL", name: "Panel", status: "IN_SERVICE"},
			},
			existing: stored,
			want: []validationError{
				{Type: "Invalid hierarchy", Label: "Parent Not Found", Line: 2, AssetID: "PNL-1", Message: "Asset PNL-1 has parent SWB-404, which does not exist in the file."},
			},
		},
		{
			name: "parent in the database has the wrong type",
			incoming: []asset{
				{line: 2, id: "PNL-1", parentID: "SS-1", assetType: "SWITCHBOARD_PANEL", name: "Panel", status: "IN_SERVICE"},
			},
			existing: stored,
			want: []validationError{
				{Type: "Invalid hierarchy", Label: "Invalid Parent Type", Line: 2, AssetID: "PNL-1", Message: "Asset PNL-1 is a SWITCHBOARD_PANEL, so its parent must be a SWITCHBOARD, but SS-1 is a SUBSTATION."},
			},
		},
		{
			name: "asset ID already in the database",
			incoming: []asset{
				{line: 2, id: "SS-1", assetType: "SUBSTATION", name: "North again", status: "IN_SERVICE"},
				{line: 3, id: "TX-1", parentID: "SS-1", assetType: "TRANSFORMER", name: "T", status: "IN_SERVICE"},
			},
			existing: stored,
			want: []validationError{
				{Type: "Already in database", Label: "Already Exists", Line: 2, AssetID: "SS-1", Message: "Asset SS-1 already exists in the database."},
			},
		},
		{
			name: "the existing validation rules still apply",
			incoming: []asset{
				{line: 2, id: "TX-1", parentID: "SS-1", assetType: "TRANSFORMER", name: "", status: "BROKEN"},
			},
			existing: stored,
			want: []validationError{
				{Type: "Missing required fields", Label: "Missing Asset Name", Line: 2, AssetID: "TX-1", Message: "Asset TX-1 is missing asset_name."},
				{Type: "Unsupported values", Label: "Invalid Status", Line: 2, AssetID: "TX-1", Message: `Asset TX-1 has unsupported operational_status "BROKEN".`},
			},
		},
		{
			name: "problems with stored assets are not reported",
			incoming: []asset{
				{line: 2, id: "SS-2", assetType: "SUBSTATION", name: "South", status: "IN_SERVICE"},
			},
			existing: []asset{{id: "OLD-1", assetType: "LEGACY", name: "", status: ""}},
			want:     []validationError{},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := validateForImport(tc.incoming, tc.existing)
			if !slices.Equal(got, tc.want) {
				t.Errorf("unexpected errors\n got: %s\nwant: %s", describe(got), describe(tc.want))
			}
		})
	}
}

func TestSortParentsFirst(t *testing.T) {
	ids := func(assets []asset) []string {
		var ids []string
		for _, a := range assets {
			ids = append(ids, a.id)
		}
		return ids
	}

	cases := []struct {
		name   string
		assets []asset
		want   []string
	}{
		{
			name: "children listed before their parents",
			assets: []asset{
				{id: "PNL-1", parentID: "SWB-1"},
				{id: "SWB-1", parentID: "SS-1"},
				{id: "SS-1"},
			},
			want: []string{"SS-1", "SWB-1", "PNL-1"},
		},
		{
			name: "already ordered rows keep their order",
			assets: []asset{
				{id: "SS-1"},
				{id: "TX-1", parentID: "SS-1"},
				{id: "SS-2"},
				{id: "TX-2", parentID: "SS-2"},
			},
			want: []string{"SS-1", "SS-2", "TX-1", "TX-2"},
		},
		{
			name: "a parent outside the batch counts as already stored",
			assets: []asset{
				{id: "PNL-1", parentID: "SWB-1"},
				{id: "SWB-1", parentID: "SS-IN-DATABASE"},
			},
			want: []string{"SWB-1", "PNL-1"},
		},
		{
			name:   "no rows",
			assets: nil,
			want:   nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			original := slices.Clone(ids(tc.assets))

			got := ids(sortParentsFirst(tc.assets))
			if !slices.Equal(got, tc.want) {
				t.Errorf("order = %v, want %v", got, tc.want)
			}
			if !slices.Equal(ids(tc.assets), original) {
				t.Errorf("the input was reordered: %v", ids(tc.assets))
			}
		})
	}
}

func TestFormatForDatabase(t *testing.T) {
	if got := formatDate("15/2/09"); got != "2009-02-15" {
		t.Errorf(`formatDate("15/2/09") = %q, want "2009-02-15"`, got)
	}
	if got := formatDate(""); got != "" {
		t.Errorf(`formatDate("") = %q, want ""`, got)
	}
	if got := formatNumber("1e3"); got != "1000" {
		t.Errorf(`formatNumber("1e3") = %q, want "1000"`, got)
	}
	if got := formatNumber("0.4"); got != "0.4" {
		t.Errorf(`formatNumber("0.4") = %q, want "0.4"`, got)
	}
	if got := formatNumber(""); got != "" {
		t.Errorf(`formatNumber("") = %q, want ""`, got)
	}
}
