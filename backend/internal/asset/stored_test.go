package asset

import (
	"slices"
	"testing"
)

func TestValidateWithStored(t *testing.T) {
	// Assets already in the database have no line number.
	stored := []Asset{
		{ID: "SS-1", Type: "SUBSTATION", Name: "North", Status: "IN_SERVICE"},
		{ID: "SWB-1", ParentID: "SS-1", Type: "SWITCHBOARD", Name: "Board", Status: "IN_SERVICE"},
	}

	cases := []struct {
		name     string
		incoming []Asset
		existing []Asset
		want     []ValidationError
	}{
		{
			name: "complete hierarchy into an empty database",
			incoming: []Asset{
				{Line: 2, ID: "SS-2", Type: "SUBSTATION", Name: "South", Status: "IN_SERVICE"},
				{Line: 3, ID: "TX-2", ParentID: "SS-2", Type: "TRANSFORMER", Name: "T", Status: "IN_SERVICE"},
			},
			want: []ValidationError{},
		},
		{
			name: "parent found in the database",
			incoming: []Asset{
				{Line: 2, ID: "PNL-1", ParentID: "SWB-1", Type: "SWITCHBOARD_PANEL", Name: "Panel", Status: "IN_SERVICE"},
			},
			existing: stored,
			want:     []ValidationError{},
		},
		{
			name: "parent in neither the upload nor the database",
			incoming: []Asset{
				{Line: 2, ID: "PNL-1", ParentID: "SWB-404", Type: "SWITCHBOARD_PANEL", Name: "Panel", Status: "IN_SERVICE"},
			},
			existing: stored,
			want: []ValidationError{
				{Type: "Invalid hierarchy", Label: "Parent Not Found", Line: 2, AssetID: "PNL-1", Message: "Asset PNL-1 has parent SWB-404, which does not exist in the file."},
			},
		},
		{
			name: "parent in the database has the wrong type",
			incoming: []Asset{
				{Line: 2, ID: "PNL-1", ParentID: "SS-1", Type: "SWITCHBOARD_PANEL", Name: "Panel", Status: "IN_SERVICE"},
			},
			existing: stored,
			want: []ValidationError{
				{Type: "Invalid hierarchy", Label: "Invalid Parent Type", Line: 2, AssetID: "PNL-1", Message: "Asset PNL-1 is a SWITCHBOARD_PANEL, so its parent must be a SWITCHBOARD, but SS-1 is a SUBSTATION."},
			},
		},
		{
			name: "asset ID already in the database",
			incoming: []Asset{
				{Line: 2, ID: "SS-1", Type: "SUBSTATION", Name: "North again", Status: "IN_SERVICE"},
				{Line: 3, ID: "TX-1", ParentID: "SS-1", Type: "TRANSFORMER", Name: "T", Status: "IN_SERVICE"},
			},
			existing: stored,
			want: []ValidationError{
				{Type: "Already in database", Label: "Already Exists", Line: 2, AssetID: "SS-1", Message: "Asset SS-1 already exists in the database."},
			},
		},
		{
			name: "the existing validation rules still apply",
			incoming: []Asset{
				{Line: 2, ID: "TX-1", ParentID: "SS-1", Type: "TRANSFORMER", Name: "", Status: "BROKEN"},
			},
			existing: stored,
			want: []ValidationError{
				{Type: "Missing required fields", Label: "Missing Asset Name", Line: 2, AssetID: "TX-1", Message: "Asset TX-1 is missing asset_name."},
				{Type: "Unsupported values", Label: "Invalid Status", Line: 2, AssetID: "TX-1", Message: `Asset TX-1 has unsupported operational_status "BROKEN".`},
			},
		},
		{
			name: "problems with stored assets are not reported",
			incoming: []Asset{
				{Line: 2, ID: "SS-2", Type: "SUBSTATION", Name: "South", Status: "IN_SERVICE"},
			},
			existing: []Asset{{ID: "OLD-1", Type: "LEGACY", Name: "", Status: ""}},
			want:     []ValidationError{},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := ValidateWithStored(tc.incoming, tc.existing)
			if !slices.Equal(got, tc.want) {
				t.Errorf("unexpected errors\n got: %s\nwant: %s", describe(got), describe(tc.want))
			}
		})
	}
}
