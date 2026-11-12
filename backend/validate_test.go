package main

import (
	"fmt"
	"slices"
	"testing"
)

// checkCase is one table entry: the rows given to a check and the errors it
// should return, in order.
type checkCase struct {
	name   string
	assets []asset
	want   []validationError
}

// runCheckCases runs every case as a subtest against one validation function.
func runCheckCases(t *testing.T, check func([]asset) []validationError, cases []checkCase) {
	t.Helper()
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := check(tc.assets)
			if !slices.Equal(got, tc.want) {
				t.Errorf("unexpected errors\n got: %s\nwant: %s", describe(got), describe(tc.want))
			}
		})
	}
}

func describe(errs []validationError) string {
	if len(errs) == 0 {
		return "(none)"
	}
	text := ""
	for _, e := range errs {
		text += fmt.Sprintf("\n  line %d [%s] %s", e.Line, e.AssetID, e.Message)
	}
	return text
}

func TestCheckRequiredFields(t *testing.T) {
	runCheckCases(t, checkRequiredFields, []checkCase{
		{
			name: "all required fields present",
			assets: []asset{
				{line: 2, id: "SS-1", assetType: "SUBSTATION", name: "North"},
			},
			want: nil,
		},
		{
			name: "optional fields may be empty",
			assets: []asset{
				{line: 2, id: "TX-1", assetType: "TRANSFORMER", name: "T1", parentID: "", status: "", voltageKV: "", ratingKVA: "", commissionedDate: ""},
			},
			want: nil,
		},
		{
			name: "missing asset_id",
			assets: []asset{
				{line: 5, id: "", assetType: "TRANSFORMER", name: "T1"},
			},
			want: []validationError{
				{Line: 5, AssetID: "", Message: "This asset is missing asset_id."},
			},
		},
		{
			name: "missing asset_type",
			assets: []asset{
				{line: 6, id: "TX-1", assetType: "", name: "T1"},
			},
			want: []validationError{
				{Line: 6, AssetID: "TX-1", Message: "Asset TX-1 is missing asset_type."},
			},
		},
		{
			name: "missing asset_name",
			assets: []asset{
				{line: 7, id: "TX-1", assetType: "TRANSFORMER", name: ""},
			},
			want: []validationError{
				{Line: 7, AssetID: "TX-1", Message: "Asset TX-1 is missing asset_name."},
			},
		},
		{
			name: "every required field missing on one row",
			assets: []asset{
				{line: 2, id: "SS-1", assetType: "SUBSTATION", name: "North"},
				{line: 3},
			},
			want: []validationError{
				{Line: 3, Message: "This asset is missing asset_id."},
				{Line: 3, Message: "This asset is missing asset_type."},
				{Line: 3, Message: "This asset is missing asset_name."},
			},
		},
		{
			name:   "no rows",
			assets: nil,
			want:   nil,
		},
	})
}

func TestCheckUniqueIDs(t *testing.T) {
	runCheckCases(t, checkUniqueIDs, []checkCase{
		{
			name: "all IDs unique",
			assets: []asset{
				{line: 2, id: "SS-1"},
				{line: 3, id: "TX-1"},
			},
			want: nil,
		},
		{
			name: "two rows share an ID and both are reported",
			assets: []asset{
				{line: 2, id: "TX-1"},
				{line: 3, id: "SS-1"},
				{line: 4, id: "TX-1"},
			},
			want: []validationError{
				{Line: 2, AssetID: "TX-1", Message: "Asset TX-1 has a duplicate asset_id, also used on line 4."},
				{Line: 4, AssetID: "TX-1", Message: "Asset TX-1 has a duplicate asset_id, also used on line 2."},
			},
		},
		{
			name: "three rows share an ID",
			assets: []asset{
				{line: 2, id: "TX-1"},
				{line: 5, id: "TX-1"},
				{line: 9, id: "TX-1"},
			},
			want: []validationError{
				{Line: 2, AssetID: "TX-1", Message: "Asset TX-1 has a duplicate asset_id, also used on line 5, 9."},
				{Line: 5, AssetID: "TX-1", Message: "Asset TX-1 has a duplicate asset_id, also used on line 2, 9."},
				{Line: 9, AssetID: "TX-1", Message: "Asset TX-1 has a duplicate asset_id, also used on line 2, 5."},
			},
		},
		{
			name: "two different IDs each duplicated",
			assets: []asset{
				{line: 2, id: "A"},
				{line: 3, id: "B"},
				{line: 4, id: "A"},
				{line: 5, id: "B"},
			},
			want: []validationError{
				{Line: 2, AssetID: "A", Message: "Asset A has a duplicate asset_id, also used on line 4."},
				{Line: 3, AssetID: "B", Message: "Asset B has a duplicate asset_id, also used on line 5."},
				{Line: 4, AssetID: "A", Message: "Asset A has a duplicate asset_id, also used on line 2."},
				{Line: 5, AssetID: "B", Message: "Asset B has a duplicate asset_id, also used on line 3."},
			},
		},
		{
			name: "IDs are case-sensitive",
			assets: []asset{
				{line: 2, id: "TX-1"},
				{line: 3, id: "tx-1"},
			},
			want: nil,
		},
		{
			name: "rows with no ID are not duplicates of each other",
			assets: []asset{
				{line: 2, id: ""},
				{line: 3, id: ""},
			},
			want: nil,
		},
	})
}

func TestCheckSupportedValues(t *testing.T) {
	runCheckCases(t, checkSupportedValues, []checkCase{
		{
			name: "every supported type and status",
			assets: []asset{
				{line: 2, id: "SS-1", assetType: "SUBSTATION", status: "IN_SERVICE"},
				{line: 3, id: "TX-1", assetType: "TRANSFORMER", status: "MAINTENANCE"},
				{line: 4, id: "LVB-1", assetType: "LV_BOARD", status: "OUT_OF_SERVICE"},
				{line: 5, id: "SWB-1", assetType: "SWITCHBOARD", status: "IN_SERVICE"},
				{line: 6, id: "PNL-1", assetType: "SWITCHBOARD_PANEL", status: "IN_SERVICE"},
			},
			want: nil,
		},
		{
			name: "unsupported asset type",
			assets: []asset{
				{line: 8, id: "GEN-1", assetType: "GENERATOR", status: "IN_SERVICE"},
			},
			want: []validationError{
				{Line: 8, AssetID: "GEN-1", Message: `Asset GEN-1 has unsupported asset_type "GENERATOR".`},
			},
		},
		{
			name: "asset type is case-sensitive",
			assets: []asset{
				{line: 8, id: "TX-1", assetType: "transformer", status: "IN_SERVICE"},
			},
			want: []validationError{
				{Line: 8, AssetID: "TX-1", Message: `Asset TX-1 has unsupported asset_type "transformer".`},
			},
		},
		{
			name: "missing asset type is left to the required-fields check",
			assets: []asset{
				{line: 8, id: "TX-1", assetType: "", status: "IN_SERVICE"},
			},
			want: nil,
		},
		{
			name: "unsupported status",
			assets: []asset{
				{line: 9, id: "PNL-1", assetType: "SWITCHBOARD_PANEL", status: "BROKEN"},
			},
			want: []validationError{
				{Line: 9, AssetID: "PNL-1", Message: `Asset PNL-1 has unsupported operational_status "BROKEN".`},
			},
		},
		{
			name: "empty status is unsupported",
			assets: []asset{
				{line: 9, id: "PNL-1", assetType: "SWITCHBOARD_PANEL", status: ""},
			},
			want: []validationError{
				{Line: 9, AssetID: "PNL-1", Message: `Asset PNL-1 has unsupported operational_status "".`},
			},
		},
		{
			name: "unsupported type and status on the same row",
			assets: []asset{
				{line: 4, id: "X-1", assetType: "CABLE", status: "RETIRED"},
			},
			want: []validationError{
				{Line: 4, AssetID: "X-1", Message: `Asset X-1 has unsupported asset_type "CABLE".`},
				{Line: 4, AssetID: "X-1", Message: `Asset X-1 has unsupported operational_status "RETIRED".`},
			},
		},
	})
}

func TestCheckDatesAndNumbers(t *testing.T) {
	runCheckCases(t, checkDatesAndNumbers, []checkCase{
		{
			name: "valid dates and numbers",
			assets: []asset{
				{line: 2, id: "A", commissionedDate: "15/2/09", voltageKV: "22", ratingKVA: "2000"},
				{line: 3, id: "B", commissionedDate: "1/12/2015", voltageKV: "0.4", ratingKVA: "0"},
				{line: 4, id: "C", commissionedDate: "05/03/11", voltageKV: "1e2", ratingKVA: "1500.5"},
				{line: 5, id: "D", commissionedDate: "29/2/20", voltageKV: "22", ratingKVA: "1000"},
			},
			want: nil,
		},
		{
			name: "empty date and numbers are allowed",
			assets: []asset{
				{line: 2, id: "A"},
			},
			want: nil,
		},
		{
			name: "date in year-month-day order",
			assets: []asset{
				{line: 6, id: "TX-BD", commissionedDate: "2026-13-41"},
			},
			want: []validationError{
				{Line: 6, AssetID: "TX-BD", Message: `Asset TX-BD has invalid commissioned_date "2026-13-41" (expected day/month/year, e.g. 15/2/09).`},
			},
		},
		{
			name: "date that does not exist",
			assets: []asset{
				{line: 6, id: "A", commissionedDate: "31/2/20"},
				{line: 7, id: "B", commissionedDate: "29/2/21"},
				{line: 8, id: "C", commissionedDate: "15/13/09"},
			},
			want: []validationError{
				{Line: 6, AssetID: "A", Message: `Asset A has invalid commissioned_date "31/2/20" (expected day/month/year, e.g. 15/2/09).`},
				{Line: 7, AssetID: "B", Message: `Asset B has invalid commissioned_date "29/2/21" (expected day/month/year, e.g. 15/2/09).`},
				{Line: 8, AssetID: "C", Message: `Asset C has invalid commissioned_date "15/13/09" (expected day/month/year, e.g. 15/2/09).`},
			},
		},
		{
			name: "date that is not a date at all",
			assets: []asset{
				{line: 6, id: "A", commissionedDate: "yesterday"},
			},
			want: []validationError{
				{Line: 6, AssetID: "A", Message: `Asset A has invalid commissioned_date "yesterday" (expected day/month/year, e.g. 15/2/09).`},
			},
		},
		{
			name: "non-numeric voltage",
			assets: []asset{
				{line: 3, id: "A", voltageKV: "abc"},
				{line: 4, id: "B", voltageKV: "22kV"},
				{line: 5, id: "C", voltageKV: "NaN"},
			},
			want: []validationError{
				{Line: 3, AssetID: "A", Message: `Asset A has non-numeric voltage_kv "abc".`},
				{Line: 4, AssetID: "B", Message: `Asset B has non-numeric voltage_kv "22kV".`},
				{Line: 5, AssetID: "C", Message: `Asset C has non-numeric voltage_kv "NaN".`},
			},
		},
		{
			name: "non-numeric rating",
			assets: []asset{
				{line: 3, id: "A", ratingKVA: "lots"},
				{line: 4, id: "B", ratingKVA: "Inf"},
			},
			want: []validationError{
				{Line: 3, AssetID: "A", Message: `Asset A has non-numeric rating_kva "lots".`},
				{Line: 4, AssetID: "B", Message: `Asset B has non-numeric rating_kva "Inf".`},
			},
		},
		{
			name: "negative rating",
			assets: []asset{
				{line: 9, id: "TX-NR", ratingKVA: "-500"},
				{line: 10, id: "TX-2", ratingKVA: "-0.5"},
			},
			want: []validationError{
				{Line: 9, AssetID: "TX-NR", Message: "Asset TX-NR has negative rating_kva -500."},
				{Line: 10, AssetID: "TX-2", Message: "Asset TX-2 has negative rating_kva -0.5."},
			},
		},
		{
			name: "bad date, voltage and rating on the same row",
			assets: []asset{
				{line: 4, id: "A", commissionedDate: "soon", voltageKV: "high", ratingKVA: "-1"},
			},
			want: []validationError{
				{Line: 4, AssetID: "A", Message: `Asset A has invalid commissioned_date "soon" (expected day/month/year, e.g. 15/2/09).`},
				{Line: 4, AssetID: "A", Message: `Asset A has non-numeric voltage_kv "high".`},
				{Line: 4, AssetID: "A", Message: "Asset A has negative rating_kva -1."},
			},
		},
	})
}

func TestCheckParentRequirements(t *testing.T) {
	runCheckCases(t, checkParentRequirements, []checkCase{
		{
			name: "substation without parent and others with one",
			assets: []asset{
				{line: 2, id: "SS-1", assetType: "SUBSTATION"},
				{line: 3, id: "TX-1", assetType: "TRANSFORMER", parentID: "SS-1"},
				{line: 4, id: "SWB-1", assetType: "SWITCHBOARD", parentID: "SS-1"},
				{line: 5, id: "PNL-1", assetType: "SWITCHBOARD_PANEL", parentID: "SWB-1"},
			},
			want: nil,
		},
		{
			name: "substation with a parent",
			assets: []asset{
				{line: 2, id: "SS-1", assetType: "SUBSTATION"},
				{line: 3, id: "SS-2", assetType: "SUBSTATION", parentID: "SS-1"},
			},
			want: []validationError{
				{Line: 3, AssetID: "SS-2", Message: "Asset SS-2 is a substation and cannot have a parent (SS-1)."},
			},
		},
		{
			name: "non-substations without a parent",
			assets: []asset{
				{line: 2, id: "TX-1", assetType: "TRANSFORMER"},
				{line: 3, id: "LVB-1", assetType: "LV_BOARD"},
				{line: 4, id: "SWB-1", assetType: "SWITCHBOARD"},
				{line: 5, id: "PNL-1", assetType: "SWITCHBOARD_PANEL"},
			},
			want: []validationError{
				{Line: 2, AssetID: "TX-1", Message: "Asset TX-1 has no parent; only substations can be without one."},
				{Line: 3, AssetID: "LVB-1", Message: "Asset LVB-1 has no parent; only substations can be without one."},
				{Line: 4, AssetID: "SWB-1", Message: "Asset SWB-1 has no parent; only substations can be without one."},
				{Line: 5, AssetID: "PNL-1", Message: "Asset PNL-1 has no parent; only substations can be without one."},
			},
		},
		{
			name: "unsupported type without a parent still needs one",
			assets: []asset{
				{line: 2, id: "GEN-1", assetType: "GENERATOR"},
			},
			want: []validationError{
				{Line: 2, AssetID: "GEN-1", Message: "Asset GEN-1 has no parent; only substations can be without one."},
			},
		},
		{
			name: "row with no type is left to the required-fields check",
			assets: []asset{
				{line: 2, id: "X-1", assetType: ""},
			},
			want: nil,
		},
		{
			name: "does not check whether the parent exists",
			assets: []asset{
				{line: 2, id: "TX-1", assetType: "TRANSFORMER", parentID: "SS-404"},
			},
			want: nil,
		},
	})
}

func TestCheckCycles(t *testing.T) {
	runCheckCases(t, checkCycles, []checkCase{
		{
			name: "valid tree",
			assets: []asset{
				{line: 2, id: "SS-1"},
				{line: 3, id: "SWB-1", parentID: "SS-1"},
				{line: 4, id: "PNL-1", parentID: "SWB-1"},
				{line: 5, id: "PNL-2", parentID: "SWB-1"},
			},
			want: nil,
		},
		{
			name: "self-parenting",
			assets: []asset{
				{line: 2, id: "SS-1"},
				{line: 3, id: "SWB-1", parentID: "SWB-1"},
			},
			want: []validationError{
				{Line: 3, AssetID: "SWB-1", Message: "Asset SWB-1 is its own parent."},
			},
		},
		{
			name: "two assets that are each other's parent",
			assets: []asset{
				{line: 2, id: "A", parentID: "B"},
				{line: 3, id: "B", parentID: "A"},
			},
			want: []validationError{
				{Line: 2, AssetID: "A", Message: "Asset A is in a circular parent relationship: A → B → A."},
				{Line: 3, AssetID: "B", Message: "Asset B is in a circular parent relationship: B → A → B."},
			},
		},
		{
			name: "loop of three",
			assets: []asset{
				{line: 2, id: "A", parentID: "B"},
				{line: 3, id: "B", parentID: "C"},
				{line: 4, id: "C", parentID: "A"},
			},
			want: []validationError{
				{Line: 2, AssetID: "A", Message: "Asset A is in a circular parent relationship: A → B → C → A."},
				{Line: 3, AssetID: "B", Message: "Asset B is in a circular parent relationship: B → C → A → B."},
				{Line: 4, AssetID: "C", Message: "Asset C is in a circular parent relationship: C → A → B → C."},
			},
		},
		{
			name: "asset below a loop is not itself in the loop",
			assets: []asset{
				{line: 2, id: "CHILD", parentID: "A"},
				{line: 3, id: "A", parentID: "B"},
				{line: 4, id: "B", parentID: "A"},
			},
			want: []validationError{
				{Line: 3, AssetID: "A", Message: "Asset A is in a circular parent relationship: A → B → A."},
				{Line: 4, AssetID: "B", Message: "Asset B is in a circular parent relationship: B → A → B."},
			},
		},
		{
			name: "two separate loops",
			assets: []asset{
				{line: 2, id: "A", parentID: "B"},
				{line: 3, id: "B", parentID: "A"},
				{line: 4, id: "SS-1"},
				{line: 5, id: "X", parentID: "X"},
			},
			want: []validationError{
				{Line: 2, AssetID: "A", Message: "Asset A is in a circular parent relationship: A → B → A."},
				{Line: 3, AssetID: "B", Message: "Asset B is in a circular parent relationship: B → A → B."},
				{Line: 5, AssetID: "X", Message: "Asset X is its own parent."},
			},
		},
		{
			name: "missing parent is not a cycle",
			assets: []asset{
				{line: 2, id: "A", parentID: "NOPE"},
			},
			want: nil,
		},
		{
			name: "rows with no ID are ignored",
			assets: []asset{
				{line: 2, id: "", parentID: ""},
				{line: 3, id: "", parentID: "A"},
				{line: 4, id: "A"},
			},
			want: nil,
		},
	})
}

func TestFindCycle(t *testing.T) {
	// parents maps each asset ID to its parent ID.
	cases := []struct {
		name    string
		parents map[string]string
		start   string
		want    []string
	}{
		{
			name:    "chain ending at a root",
			parents: map[string]string{"PNL": "SWB", "SWB": "SS", "SS": ""},
			start:   "PNL",
			want:    nil,
		},
		{
			name:    "chain ending at a missing parent",
			parents: map[string]string{"PNL": "SWB", "SWB": "GONE"},
			start:   "PNL",
			want:    nil,
		},
		{
			name:    "self loop",
			parents: map[string]string{"A": "A"},
			start:   "A",
			want:    []string{"A"},
		},
		{
			name:    "loop of three from a member",
			parents: map[string]string{"A": "B", "B": "C", "C": "A"},
			start:   "B",
			want:    []string{"B", "C", "A"},
		},
		{
			name:    "loop reached from below excludes the starting asset",
			parents: map[string]string{"CHILD": "A", "A": "B", "B": "A"},
			start:   "CHILD",
			want:    []string{"A", "B"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			byID := map[string]asset{}
			for id, parentID := range tc.parents {
				byID[id] = asset{id: id, parentID: parentID}
			}
			state := map[string]visitState{}

			got := findCycle(tc.start, byID, state, nil)
			if !slices.Equal(got, tc.want) {
				t.Errorf("findCycle(%q) = %v, want %v", tc.start, got, tc.want)
			}
			if state[tc.start] != visited {
				t.Errorf("start asset %q was not marked visited", tc.start)
			}
		})
	}

	t.Run("already explored assets are not searched again", func(t *testing.T) {
		byID := map[string]asset{
			"A":     {id: "A", parentID: "B"},
			"B":     {id: "B", parentID: "A"},
			"CHILD": {id: "CHILD", parentID: "A"},
		}
		state := map[string]visitState{}

		if got := findCycle("A", byID, state, nil); !slices.Equal(got, []string{"A", "B"}) {
			t.Fatalf("first search = %v, want [A B]", got)
		}
		if got := findCycle("CHILD", byID, state, nil); got != nil {
			t.Errorf("second search = %v, want nil because the loop was already found", got)
		}
	})
}

func TestCheckHierarchy(t *testing.T) {
	runCheckCases(t, checkHierarchy, []checkCase{
		{
			name: "every permitted relationship",
			assets: []asset{
				{line: 2, id: "SS-1", assetType: "SUBSTATION"},
				{line: 3, id: "TX-1", assetType: "TRANSFORMER", parentID: "SS-1"},
				{line: 4, id: "LVB-1", assetType: "LV_BOARD", parentID: "SS-1"},
				{line: 5, id: "SWB-1", assetType: "SWITCHBOARD", parentID: "SS-1"},
				{line: 6, id: "PNL-1", assetType: "SWITCHBOARD_PANEL", parentID: "SWB-1"},
			},
			want: nil,
		},
		{
			name: "child may appear before its parent in the file",
			assets: []asset{
				{line: 2, id: "PNL-1", assetType: "SWITCHBOARD_PANEL", parentID: "SWB-1"},
				{line: 3, id: "SWB-1", assetType: "SWITCHBOARD", parentID: "SS-1"},
				{line: 4, id: "SS-1", assetType: "SUBSTATION"},
			},
			want: nil,
		},
		{
			name: "parent does not exist",
			assets: []asset{
				{line: 2, id: "SS-1", assetType: "SUBSTATION"},
				{line: 3, id: "PNL-M", assetType: "SWITCHBOARD_PANEL", parentID: "SWB-404"},
			},
			want: []validationError{
				{Line: 3, AssetID: "PNL-M", Message: "Asset PNL-M has parent SWB-404, which does not exist in the file."},
			},
		},
		{
			name: "panel directly under a substation",
			assets: []asset{
				{line: 2, id: "SS-1", assetType: "SUBSTATION"},
				{line: 3, id: "PNL-BP", assetType: "SWITCHBOARD_PANEL", parentID: "SS-1"},
			},
			want: []validationError{
				{Line: 3, AssetID: "PNL-BP", Message: "Asset PNL-BP is a SWITCHBOARD_PANEL, so its parent must be a SWITCHBOARD, but SS-1 is a SUBSTATION."},
			},
		},
		{
			name: "transformer under a switchboard",
			assets: []asset{
				{line: 2, id: "SS-1", assetType: "SUBSTATION"},
				{line: 3, id: "SWB-1", assetType: "SWITCHBOARD", parentID: "SS-1"},
				{line: 4, id: "TX-BP", assetType: "TRANSFORMER", parentID: "SWB-1"},
			},
			want: []validationError{
				{Line: 4, AssetID: "TX-BP", Message: "Asset TX-BP is a TRANSFORMER, so its parent must be a SUBSTATION, but SWB-1 is a SWITCHBOARD."},
			},
		},
		{
			name: "top-level ancestor is not a substation",
			assets: []asset{
				{line: 2, id: "SWB-1", assetType: "SWITCHBOARD"},
				{line: 3, id: "PNL-1", assetType: "SWITCHBOARD_PANEL", parentID: "SWB-1"},
			},
			want: []validationError{
				{Line: 3, AssetID: "PNL-1", Message: "Asset PNL-1 does not resolve to a substation: its top-level ancestor SWB-1 is a SWITCHBOARD."},
			},
		},
		{
			name: "ancestor further up has a missing parent",
			assets: []asset{
				{line: 2, id: "SWB-1", assetType: "SWITCHBOARD", parentID: "SS-404"},
				{line: 3, id: "PNL-1", assetType: "SWITCHBOARD_PANEL", parentID: "SWB-1"},
			},
			want: []validationError{
				{Line: 2, AssetID: "SWB-1", Message: "Asset SWB-1 has parent SS-404, which does not exist in the file."},
				{Line: 3, AssetID: "PNL-1", Message: "Asset PNL-1 does not resolve to a substation: its ancestor SWB-1 has a parent that does not exist."},
			},
		},
		{
			name: "circular chain never reaches a substation",
			assets: []asset{
				{line: 2, id: "SWB-A", assetType: "SWITCHBOARD", parentID: "SWB-B"},
				{line: 3, id: "SWB-B", assetType: "SWITCHBOARD", parentID: "SWB-A"},
				{line: 4, id: "PNL-1", assetType: "SWITCHBOARD_PANEL", parentID: "SWB-A"},
			},
			want: []validationError{
				{Line: 2, AssetID: "SWB-A", Message: "Asset SWB-A is a SWITCHBOARD, so its parent must be a SUBSTATION, but SWB-B is a SWITCHBOARD."},
				{Line: 2, AssetID: "SWB-A", Message: "Asset SWB-A does not resolve to a substation: its parent chain is circular."},
				{Line: 3, AssetID: "SWB-B", Message: "Asset SWB-B is a SWITCHBOARD, so its parent must be a SUBSTATION, but SWB-A is a SWITCHBOARD."},
				{Line: 3, AssetID: "SWB-B", Message: "Asset SWB-B does not resolve to a substation: its parent chain is circular."},
				{Line: 4, AssetID: "PNL-1", Message: "Asset PNL-1 does not resolve to a substation: its parent chain is circular."},
			},
		},
		{
			name: "self-parenting asset",
			assets: []asset{
				{line: 2, id: "SWB-1", assetType: "SWITCHBOARD", parentID: "SWB-1"},
			},
			want: []validationError{
				{Line: 2, AssetID: "SWB-1", Message: "Asset SWB-1 is a SWITCHBOARD, so its parent must be a SUBSTATION, but SWB-1 is a SWITCHBOARD."},
				{Line: 2, AssetID: "SWB-1", Message: "Asset SWB-1 does not resolve to a substation: its parent chain is circular."},
			},
		},
		{
			name: "parent with no type",
			assets: []asset{
				{line: 2, id: "X-1", assetType: ""},
				{line: 3, id: "TX-1", assetType: "TRANSFORMER", parentID: "X-1"},
			},
			want: []validationError{
				{Line: 3, AssetID: "TX-1", Message: "Asset TX-1 is a TRANSFORMER, so its parent must be a SUBSTATION, but X-1 has no asset_type."},
				{Line: 3, AssetID: "TX-1", Message: "Asset TX-1 does not resolve to a substation: its top-level ancestor X-1 has no asset_type."},
			},
		},
		{
			name: "unsupported type has no parent-type rule but must still resolve",
			assets: []asset{
				{line: 2, id: "SS-1", assetType: "SUBSTATION"},
				{line: 3, id: "GEN-1", assetType: "GENERATOR", parentID: "SS-1"},
				{line: 4, id: "GEN-2", assetType: "GENERATOR", parentID: "SS-404"},
			},
			want: []validationError{
				{Line: 4, AssetID: "GEN-2", Message: "Asset GEN-2 has parent SS-404, which does not exist in the file."},
			},
		},
		{
			name: "assets without a parent are left to the parent-requirements check",
			assets: []asset{
				{line: 2, id: "TX-1", assetType: "TRANSFORMER"},
				{line: 3, id: "SS-1", assetType: "SUBSTATION"},
			},
			want: nil,
		},
	})
}

func TestCheckDuplicateSerials(t *testing.T) {
	runCheckCases(t, checkDuplicateSerials, []checkCase{
		{
			name: "all serial numbers unique",
			assets: []asset{
				{line: 2, id: "TX-1", serialNumber: "SER-1"},
				{line: 3, id: "TX-2", serialNumber: "SER-2"},
			},
			want: nil,
		},
		{
			name: "two assets share a serial number and both are reported",
			assets: []asset{
				{line: 2, id: "TX-D1", serialNumber: "SER-TX-DUP"},
				{line: 3, id: "TX-3", serialNumber: "SER-3"},
				{line: 4, id: "TX-D2", serialNumber: "SER-TX-DUP"},
			},
			want: []validationError{
				{Line: 2, AssetID: "TX-D1", Message: "Asset TX-D1 has duplicate serial_number SER-TX-DUP, also used on line 4."},
				{Line: 4, AssetID: "TX-D2", Message: "Asset TX-D2 has duplicate serial_number SER-TX-DUP, also used on line 2."},
			},
		},
		{
			name: "empty serial numbers are not duplicates of each other",
			assets: []asset{
				{line: 2, id: "TX-1", serialNumber: ""},
				{line: 3, id: "TX-2", serialNumber: ""},
			},
			want: nil,
		},
	})
}
