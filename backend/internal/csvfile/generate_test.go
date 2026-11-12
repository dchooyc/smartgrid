package csvfile

import (
	"testing"
	"time"

	"smartgrid/backend/internal/asset"
)

func TestBuildFiles(t *testing.T) {
	header := []string{"asset_id", "parent_asset_id", "asset_name"}
	now := time.Date(2026, time.January, 15, 9, 30, 5, 0, time.UTC)

	t.Run("splits rows between the two files in their original order", func(t *testing.T) {
		assets := []asset.Asset{
			{Line: 2, Record: []string{"SS-1", "", "North, Main"}},
			{Line: 3, Record: []string{"TX-1", "SS-404", " T1 "}},
			{Line: 4, Record: []string{"TX-2", "SS-1", "T2"}},
			{Line: 5, Record: []string{"", "TX-9", ""}},
		}
		errs := []asset.ValidationError{
			{Line: 5, Label: "Missing Asset ID"},
			{Line: 3, Label: "Parent Not Found"},
			{Line: 5, Label: "Missing Asset Name"},
			{Line: 5, Label: "Missing Asset ID"}, // repeated labels are listed once
		}

		valid, invalid, err := BuildFiles("grid_assets.csv", header, assets, errs, now)
		if err != nil {
			t.Fatal(err)
		}

		wantValid := "asset_id,parent_asset_id,asset_name\n" +
			"SS-1,,\"North, Main\"\n" +
			"TX-2,SS-1,T2\n"
		wantInvalid := "Validation Errors (Remove Before Reimporting),asset_id,parent_asset_id,asset_name\n" +
			"Parent Not Found,TX-1,SS-404,\" T1 \"\n" + // spaces are kept, so the value is quoted
			"Missing Asset ID; Missing Asset Name,,TX-9,\n"

		if valid.Content != wantValid {
			t.Errorf("valid file\n got: %q\nwant: %q", valid.Content, wantValid)
		}
		if invalid.Content != wantInvalid {
			t.Errorf("invalid file\n got: %q\nwant: %q", invalid.Content, wantInvalid)
		}
		if want := "grid_assets_valid_20260115_093005.csv"; valid.Name != want {
			t.Errorf("valid name = %q, want %q", valid.Name, want)
		}
		if want := "grid_assets_invalid_20260115_093005.csv"; invalid.Name != want {
			t.Errorf("invalid name = %q, want %q", invalid.Name, want)
		}
	})

	t.Run("no errors gives a header-only invalid file", func(t *testing.T) {
		assets := []asset.Asset{{Line: 2, Record: []string{"SS-1", "", "North"}}}

		valid, invalid, err := BuildFiles("assets.CSV", header, assets, nil, now)
		if err != nil {
			t.Fatal(err)
		}

		if want := "asset_id,parent_asset_id,asset_name\nSS-1,,North\n"; valid.Content != want {
			t.Errorf("valid file = %q, want %q", valid.Content, want)
		}
		if want := "Validation Errors (Remove Before Reimporting),asset_id,parent_asset_id,asset_name\n"; invalid.Content != want {
			t.Errorf("invalid file = %q, want %q", invalid.Content, want)
		}
		if want := "assets_valid_20260115_093005.csv"; valid.Name != want {
			t.Errorf("valid name = %q, want %q", valid.Name, want)
		}
	})

	t.Run("no valid rows gives a header-only valid file", func(t *testing.T) {
		assets := []asset.Asset{{Line: 2, Record: []string{"", "", ""}}}
		errs := []asset.ValidationError{{Line: 2, Label: "Missing Asset ID"}}

		valid, _, err := BuildFiles("assets.csv", header, assets, errs, now)
		if err != nil {
			t.Fatal(err)
		}
		if want := "asset_id,parent_asset_id,asset_name\n"; valid.Content != want {
			t.Errorf("valid file = %q, want %q", valid.Content, want)
		}
	})
}
