package asset

import (
	"fmt"
	"slices"
	"testing"
)

// checkCase is one table entry: the rows given to a check and the errors it
// should return, in order.
type checkCase struct {
	name   string
	assets []Asset
	want   []ValidationError
}

// runCheckCases runs every case as a subtest against one validation function.
func runCheckCases(t *testing.T, check func([]Asset) []ValidationError, cases []checkCase) {
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

func describe(errs []ValidationError) string {
	if len(errs) == 0 {
		return "(none)"
	}
	text := ""
	for _, e := range errs {
		text += fmt.Sprintf("\n  line %d [%s] (%s) %s", e.Line, e.AssetID, e.Label, e.Message)
	}
	return text
}

func TestCheckRequiredFields(t *testing.T) {
	runCheckCases(t, checkRequiredFields, []checkCase{
		{
			name: "all required fields present",
			assets: []Asset{
				{Line: 2, ID: "SS-1", Type: "SUBSTATION", Name: "North"},
			},
			want: nil,
		},
		{
			name: "optional fields may be empty",
			assets: []Asset{
				{Line: 2, ID: "TX-1", Type: "TRANSFORMER", Name: "T1", ParentID: "", Status: "", VoltageKV: "", RatingKVA: "", CommissionedDate: ""},
			},
			want: nil,
		},
		{
			name: "missing asset_id",
			assets: []Asset{
				{Line: 5, ID: "", Type: "TRANSFORMER", Name: "T1"},
			},
			want: []ValidationError{
				{Line: 5, AssetID: "", Label: "Missing Asset ID", Message: "This asset is missing asset_id."},
			},
		},
		{
			name: "missing asset_type",
			assets: []Asset{
				{Line: 6, ID: "TX-1", Type: "", Name: "T1"},
			},
			want: []ValidationError{
				{Line: 6, AssetID: "TX-1", Label: "Missing Asset Type", Message: "Asset TX-1 is missing asset_type."},
			},
		},
		{
			name: "missing asset_name",
			assets: []Asset{
				{Line: 7, ID: "TX-1", Type: "TRANSFORMER", Name: ""},
			},
			want: []ValidationError{
				{Line: 7, AssetID: "TX-1", Label: "Missing Asset Name", Message: "Asset TX-1 is missing asset_name."},
			},
		},
		{
			name: "every required field missing on one row",
			assets: []Asset{
				{Line: 2, ID: "SS-1", Type: "SUBSTATION", Name: "North"},
				{Line: 3},
			},
			want: []ValidationError{
				{Line: 3, Label: "Missing Asset ID", Message: "This asset is missing asset_id."},
				{Line: 3, Label: "Missing Asset Type", Message: "This asset is missing asset_type."},
				{Line: 3, Label: "Missing Asset Name", Message: "This asset is missing asset_name."},
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
			assets: []Asset{
				{Line: 2, ID: "SS-1"},
				{Line: 3, ID: "TX-1"},
			},
			want: nil,
		},
		{
			name: "two rows share an ID and both are reported",
			assets: []Asset{
				{Line: 2, ID: "TX-1"},
				{Line: 3, ID: "SS-1"},
				{Line: 4, ID: "TX-1"},
			},
			want: []ValidationError{
				{Line: 2, AssetID: "TX-1", Label: "Duplicate ID", Message: "Asset TX-1 has a duplicate asset_id, also used on line 4."},
				{Line: 4, AssetID: "TX-1", Label: "Duplicate ID", Message: "Asset TX-1 has a duplicate asset_id, also used on line 2."},
			},
		},
		{
			name: "three rows share an ID",
			assets: []Asset{
				{Line: 2, ID: "TX-1"},
				{Line: 5, ID: "TX-1"},
				{Line: 9, ID: "TX-1"},
			},
			want: []ValidationError{
				{Line: 2, AssetID: "TX-1", Label: "Duplicate ID", Message: "Asset TX-1 has a duplicate asset_id, also used on line 5, 9."},
				{Line: 5, AssetID: "TX-1", Label: "Duplicate ID", Message: "Asset TX-1 has a duplicate asset_id, also used on line 2, 9."},
				{Line: 9, AssetID: "TX-1", Label: "Duplicate ID", Message: "Asset TX-1 has a duplicate asset_id, also used on line 2, 5."},
			},
		},
		{
			name: "two different IDs each duplicated",
			assets: []Asset{
				{Line: 2, ID: "A"},
				{Line: 3, ID: "B"},
				{Line: 4, ID: "A"},
				{Line: 5, ID: "B"},
			},
			want: []ValidationError{
				{Line: 2, AssetID: "A", Label: "Duplicate ID", Message: "Asset A has a duplicate asset_id, also used on line 4."},
				{Line: 3, AssetID: "B", Label: "Duplicate ID", Message: "Asset B has a duplicate asset_id, also used on line 5."},
				{Line: 4, AssetID: "A", Label: "Duplicate ID", Message: "Asset A has a duplicate asset_id, also used on line 2."},
				{Line: 5, AssetID: "B", Label: "Duplicate ID", Message: "Asset B has a duplicate asset_id, also used on line 3."},
			},
		},
		{
			name: "IDs are case-sensitive",
			assets: []Asset{
				{Line: 2, ID: "TX-1"},
				{Line: 3, ID: "tx-1"},
			},
			want: nil,
		},
		{
			name: "rows with no ID are not duplicates of each other",
			assets: []Asset{
				{Line: 2, ID: ""},
				{Line: 3, ID: ""},
			},
			want: nil,
		},
	})
}

func TestCheckSupportedValues(t *testing.T) {
	runCheckCases(t, checkSupportedValues, []checkCase{
		{
			name: "every supported type and status",
			assets: []Asset{
				{Line: 2, ID: "SS-1", Type: "SUBSTATION", Status: "IN_SERVICE"},
				{Line: 3, ID: "TX-1", Type: "TRANSFORMER", Status: "MAINTENANCE"},
				{Line: 4, ID: "LVB-1", Type: "LV_BOARD", Status: "OUT_OF_SERVICE"},
				{Line: 5, ID: "SWB-1", Type: "SWITCHBOARD", Status: "IN_SERVICE"},
				{Line: 6, ID: "PNL-1", Type: "SWITCHBOARD_PANEL", Status: "IN_SERVICE"},
			},
			want: nil,
		},
		{
			name: "unsupported asset type",
			assets: []Asset{
				{Line: 8, ID: "GEN-1", Type: "GENERATOR", Status: "IN_SERVICE"},
			},
			want: []ValidationError{
				{Line: 8, AssetID: "GEN-1", Label: "Invalid Asset Type", Message: `Asset GEN-1 has unsupported asset_type "GENERATOR".`},
			},
		},
		{
			name: "asset type is case-sensitive",
			assets: []Asset{
				{Line: 8, ID: "TX-1", Type: "transformer", Status: "IN_SERVICE"},
			},
			want: []ValidationError{
				{Line: 8, AssetID: "TX-1", Label: "Invalid Asset Type", Message: `Asset TX-1 has unsupported asset_type "transformer".`},
			},
		},
		{
			name: "missing asset type is left to the required-fields check",
			assets: []Asset{
				{Line: 8, ID: "TX-1", Type: "", Status: "IN_SERVICE"},
			},
			want: nil,
		},
		{
			name: "unsupported status",
			assets: []Asset{
				{Line: 9, ID: "PNL-1", Type: "SWITCHBOARD_PANEL", Status: "BROKEN"},
			},
			want: []ValidationError{
				{Line: 9, AssetID: "PNL-1", Label: "Invalid Status", Message: `Asset PNL-1 has unsupported operational_status "BROKEN".`},
			},
		},
		{
			name: "empty status is unsupported",
			assets: []Asset{
				{Line: 9, ID: "PNL-1", Type: "SWITCHBOARD_PANEL", Status: ""},
			},
			want: []ValidationError{
				{Line: 9, AssetID: "PNL-1", Label: "Invalid Status", Message: `Asset PNL-1 has unsupported operational_status "".`},
			},
		},
		{
			name: "unsupported type and status on the same row",
			assets: []Asset{
				{Line: 4, ID: "X-1", Type: "CABLE", Status: "RETIRED"},
			},
			want: []ValidationError{
				{Line: 4, AssetID: "X-1", Label: "Invalid Asset Type", Message: `Asset X-1 has unsupported asset_type "CABLE".`},
				{Line: 4, AssetID: "X-1", Label: "Invalid Status", Message: `Asset X-1 has unsupported operational_status "RETIRED".`},
			},
		},
	})
}

func TestCheckDatesAndNumbers(t *testing.T) {
	runCheckCases(t, checkDatesAndNumbers, []checkCase{
		{
			name: "valid dates and numbers",
			assets: []Asset{
				{Line: 2, ID: "A", CommissionedDate: "15/2/09", VoltageKV: "22", RatingKVA: "2000"},
				{Line: 3, ID: "B", CommissionedDate: "1/12/2015", VoltageKV: "0.4", RatingKVA: "0"},
				{Line: 4, ID: "C", CommissionedDate: "05/03/11", VoltageKV: "1e2", RatingKVA: "1500.5"},
				{Line: 5, ID: "D", CommissionedDate: "29/2/20", VoltageKV: "22", RatingKVA: "1000"},
			},
			want: nil,
		},
		{
			name: "empty date and numbers are allowed",
			assets: []Asset{
				{Line: 2, ID: "A"},
			},
			want: nil,
		},
		{
			name: "date in year-month-day order",
			assets: []Asset{
				{Line: 6, ID: "TX-BD", CommissionedDate: "2026-13-41"},
			},
			want: []ValidationError{
				{Line: 6, AssetID: "TX-BD", Label: "Invalid Date", Message: `Asset TX-BD has invalid commissioned_date "2026-13-41" (expected day/month/year, e.g. 15/2/09).`},
			},
		},
		{
			name: "date that does not exist",
			assets: []Asset{
				{Line: 6, ID: "A", CommissionedDate: "31/2/20"},
				{Line: 7, ID: "B", CommissionedDate: "29/2/21"},
				{Line: 8, ID: "C", CommissionedDate: "15/13/09"},
			},
			want: []ValidationError{
				{Line: 6, AssetID: "A", Label: "Invalid Date", Message: `Asset A has invalid commissioned_date "31/2/20" (expected day/month/year, e.g. 15/2/09).`},
				{Line: 7, AssetID: "B", Label: "Invalid Date", Message: `Asset B has invalid commissioned_date "29/2/21" (expected day/month/year, e.g. 15/2/09).`},
				{Line: 8, AssetID: "C", Label: "Invalid Date", Message: `Asset C has invalid commissioned_date "15/13/09" (expected day/month/year, e.g. 15/2/09).`},
			},
		},
		{
			name: "date that is not a date at all",
			assets: []Asset{
				{Line: 6, ID: "A", CommissionedDate: "yesterday"},
			},
			want: []ValidationError{
				{Line: 6, AssetID: "A", Label: "Invalid Date", Message: `Asset A has invalid commissioned_date "yesterday" (expected day/month/year, e.g. 15/2/09).`},
			},
		},
		{
			name: "non-numeric voltage",
			assets: []Asset{
				{Line: 3, ID: "A", VoltageKV: "abc"},
				{Line: 4, ID: "B", VoltageKV: "22kV"},
				{Line: 5, ID: "C", VoltageKV: "NaN"},
			},
			want: []ValidationError{
				{Line: 3, AssetID: "A", Label: "Invalid Voltage", Message: `Asset A has non-numeric voltage_kv "abc".`},
				{Line: 4, AssetID: "B", Label: "Invalid Voltage", Message: `Asset B has non-numeric voltage_kv "22kV".`},
				{Line: 5, AssetID: "C", Label: "Invalid Voltage", Message: `Asset C has non-numeric voltage_kv "NaN".`},
			},
		},
		{
			name: "non-numeric rating",
			assets: []Asset{
				{Line: 3, ID: "A", RatingKVA: "lots"},
				{Line: 4, ID: "B", RatingKVA: "Inf"},
			},
			want: []ValidationError{
				{Line: 3, AssetID: "A", Label: "Invalid Rating", Message: `Asset A has non-numeric rating_kva "lots".`},
				{Line: 4, AssetID: "B", Label: "Invalid Rating", Message: `Asset B has non-numeric rating_kva "Inf".`},
			},
		},
		{
			name: "negative rating",
			assets: []Asset{
				{Line: 9, ID: "TX-NR", RatingKVA: "-500"},
				{Line: 10, ID: "TX-2", RatingKVA: "-0.5"},
			},
			want: []ValidationError{
				{Line: 9, AssetID: "TX-NR", Label: "Negative Rating", Message: "Asset TX-NR has negative rating_kva -500."},
				{Line: 10, AssetID: "TX-2", Label: "Negative Rating", Message: "Asset TX-2 has negative rating_kva -0.5."},
			},
		},
		{
			name: "bad date, voltage and rating on the same row",
			assets: []Asset{
				{Line: 4, ID: "A", CommissionedDate: "soon", VoltageKV: "high", RatingKVA: "-1"},
			},
			want: []ValidationError{
				{Line: 4, AssetID: "A", Label: "Invalid Date", Message: `Asset A has invalid commissioned_date "soon" (expected day/month/year, e.g. 15/2/09).`},
				{Line: 4, AssetID: "A", Label: "Invalid Voltage", Message: `Asset A has non-numeric voltage_kv "high".`},
				{Line: 4, AssetID: "A", Label: "Negative Rating", Message: "Asset A has negative rating_kva -1."},
			},
		},
	})
}

func TestCheckParentRequirements(t *testing.T) {
	runCheckCases(t, checkParentRequirements, []checkCase{
		{
			name: "substation without parent and others with one",
			assets: []Asset{
				{Line: 2, ID: "SS-1", Type: "SUBSTATION"},
				{Line: 3, ID: "TX-1", Type: "TRANSFORMER", ParentID: "SS-1"},
				{Line: 4, ID: "SWB-1", Type: "SWITCHBOARD", ParentID: "SS-1"},
				{Line: 5, ID: "PNL-1", Type: "SWITCHBOARD_PANEL", ParentID: "SWB-1"},
			},
			want: nil,
		},
		{
			name: "substation with a parent",
			assets: []Asset{
				{Line: 2, ID: "SS-1", Type: "SUBSTATION"},
				{Line: 3, ID: "SS-2", Type: "SUBSTATION", ParentID: "SS-1"},
			},
			want: []ValidationError{
				{Line: 3, AssetID: "SS-2", Label: "Substation Has Parent", Message: "Asset SS-2 is a substation and cannot have a parent (SS-1)."},
			},
		},
		{
			name: "non-substations without a parent",
			assets: []Asset{
				{Line: 2, ID: "TX-1", Type: "TRANSFORMER"},
				{Line: 3, ID: "LVB-1", Type: "LV_BOARD"},
				{Line: 4, ID: "SWB-1", Type: "SWITCHBOARD"},
				{Line: 5, ID: "PNL-1", Type: "SWITCHBOARD_PANEL"},
			},
			want: []ValidationError{
				{Line: 2, AssetID: "TX-1", Label: "Missing Parent", Message: "Asset TX-1 has no parent; only substations can be without one."},
				{Line: 3, AssetID: "LVB-1", Label: "Missing Parent", Message: "Asset LVB-1 has no parent; only substations can be without one."},
				{Line: 4, AssetID: "SWB-1", Label: "Missing Parent", Message: "Asset SWB-1 has no parent; only substations can be without one."},
				{Line: 5, AssetID: "PNL-1", Label: "Missing Parent", Message: "Asset PNL-1 has no parent; only substations can be without one."},
			},
		},
		{
			name: "unsupported type without a parent still needs one",
			assets: []Asset{
				{Line: 2, ID: "GEN-1", Type: "GENERATOR"},
			},
			want: []ValidationError{
				{Line: 2, AssetID: "GEN-1", Label: "Missing Parent", Message: "Asset GEN-1 has no parent; only substations can be without one."},
			},
		},
		{
			name: "row with no type is left to the required-fields check",
			assets: []Asset{
				{Line: 2, ID: "X-1", Type: ""},
			},
			want: nil,
		},
		{
			name: "does not check whether the parent exists",
			assets: []Asset{
				{Line: 2, ID: "TX-1", Type: "TRANSFORMER", ParentID: "SS-404"},
			},
			want: nil,
		},
	})
}

func TestCheckCycles(t *testing.T) {
	runCheckCases(t, checkCycles, []checkCase{
		{
			name: "valid tree",
			assets: []Asset{
				{Line: 2, ID: "SS-1"},
				{Line: 3, ID: "SWB-1", ParentID: "SS-1"},
				{Line: 4, ID: "PNL-1", ParentID: "SWB-1"},
				{Line: 5, ID: "PNL-2", ParentID: "SWB-1"},
			},
			want: nil,
		},
		{
			name: "self-parenting",
			assets: []Asset{
				{Line: 2, ID: "SS-1"},
				{Line: 3, ID: "SWB-1", ParentID: "SWB-1"},
			},
			want: []ValidationError{
				{Line: 3, AssetID: "SWB-1", Label: "Self Parent", Message: "Asset SWB-1 is its own parent."},
			},
		},
		{
			name: "two assets that are each other's parent",
			assets: []Asset{
				{Line: 2, ID: "A", ParentID: "B"},
				{Line: 3, ID: "B", ParentID: "A"},
			},
			want: []ValidationError{
				{Line: 2, AssetID: "A", Label: "Cycle Detected", Message: "Asset A is in a circular parent relationship: A → B → A."},
				{Line: 3, AssetID: "B", Label: "Cycle Detected", Message: "Asset B is in a circular parent relationship: B → A → B."},
			},
		},
		{
			name: "loop of three",
			assets: []Asset{
				{Line: 2, ID: "A", ParentID: "B"},
				{Line: 3, ID: "B", ParentID: "C"},
				{Line: 4, ID: "C", ParentID: "A"},
			},
			want: []ValidationError{
				{Line: 2, AssetID: "A", Label: "Cycle Detected", Message: "Asset A is in a circular parent relationship: A → B → C → A."},
				{Line: 3, AssetID: "B", Label: "Cycle Detected", Message: "Asset B is in a circular parent relationship: B → C → A → B."},
				{Line: 4, AssetID: "C", Label: "Cycle Detected", Message: "Asset C is in a circular parent relationship: C → A → B → C."},
			},
		},
		{
			name: "asset below a loop is not itself in the loop",
			assets: []Asset{
				{Line: 2, ID: "CHILD", ParentID: "A"},
				{Line: 3, ID: "A", ParentID: "B"},
				{Line: 4, ID: "B", ParentID: "A"},
			},
			want: []ValidationError{
				{Line: 3, AssetID: "A", Label: "Cycle Detected", Message: "Asset A is in a circular parent relationship: A → B → A."},
				{Line: 4, AssetID: "B", Label: "Cycle Detected", Message: "Asset B is in a circular parent relationship: B → A → B."},
			},
		},
		{
			name: "two separate loops",
			assets: []Asset{
				{Line: 2, ID: "A", ParentID: "B"},
				{Line: 3, ID: "B", ParentID: "A"},
				{Line: 4, ID: "SS-1"},
				{Line: 5, ID: "X", ParentID: "X"},
			},
			want: []ValidationError{
				{Line: 2, AssetID: "A", Label: "Cycle Detected", Message: "Asset A is in a circular parent relationship: A → B → A."},
				{Line: 3, AssetID: "B", Label: "Cycle Detected", Message: "Asset B is in a circular parent relationship: B → A → B."},
				{Line: 5, AssetID: "X", Label: "Self Parent", Message: "Asset X is its own parent."},
			},
		},
		{
			name: "missing parent is not a cycle",
			assets: []Asset{
				{Line: 2, ID: "A", ParentID: "NOPE"},
			},
			want: nil,
		},
		{
			name: "rows with no ID are ignored",
			assets: []Asset{
				{Line: 2, ID: "", ParentID: ""},
				{Line: 3, ID: "", ParentID: "A"},
				{Line: 4, ID: "A"},
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
			byID := map[string]Asset{}
			for id, parentID := range tc.parents {
				byID[id] = Asset{ID: id, ParentID: parentID}
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
		byID := map[string]Asset{
			"A":     {ID: "A", ParentID: "B"},
			"B":     {ID: "B", ParentID: "A"},
			"CHILD": {ID: "CHILD", ParentID: "A"},
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
			assets: []Asset{
				{Line: 2, ID: "SS-1", Type: "SUBSTATION"},
				{Line: 3, ID: "TX-1", Type: "TRANSFORMER", ParentID: "SS-1"},
				{Line: 4, ID: "LVB-1", Type: "LV_BOARD", ParentID: "SS-1"},
				{Line: 5, ID: "SWB-1", Type: "SWITCHBOARD", ParentID: "SS-1"},
				{Line: 6, ID: "PNL-1", Type: "SWITCHBOARD_PANEL", ParentID: "SWB-1"},
			},
			want: nil,
		},
		{
			name: "child may appear before its parent in the file",
			assets: []Asset{
				{Line: 2, ID: "PNL-1", Type: "SWITCHBOARD_PANEL", ParentID: "SWB-1"},
				{Line: 3, ID: "SWB-1", Type: "SWITCHBOARD", ParentID: "SS-1"},
				{Line: 4, ID: "SS-1", Type: "SUBSTATION"},
			},
			want: nil,
		},
		{
			name: "parent does not exist",
			assets: []Asset{
				{Line: 2, ID: "SS-1", Type: "SUBSTATION"},
				{Line: 3, ID: "PNL-M", Type: "SWITCHBOARD_PANEL", ParentID: "SWB-404"},
			},
			want: []ValidationError{
				{Line: 3, AssetID: "PNL-M", Label: "Parent Not Found", Message: "Asset PNL-M has parent SWB-404, which does not exist in the file."},
			},
		},
		{
			name: "panel directly under a substation",
			assets: []Asset{
				{Line: 2, ID: "SS-1", Type: "SUBSTATION"},
				{Line: 3, ID: "PNL-BP", Type: "SWITCHBOARD_PANEL", ParentID: "SS-1"},
			},
			want: []ValidationError{
				{Line: 3, AssetID: "PNL-BP", Label: "Invalid Parent Type", Message: "Asset PNL-BP is a SWITCHBOARD_PANEL, so its parent must be a SWITCHBOARD, but SS-1 is a SUBSTATION."},
			},
		},
		{
			name: "transformer under a switchboard",
			assets: []Asset{
				{Line: 2, ID: "SS-1", Type: "SUBSTATION"},
				{Line: 3, ID: "SWB-1", Type: "SWITCHBOARD", ParentID: "SS-1"},
				{Line: 4, ID: "TX-BP", Type: "TRANSFORMER", ParentID: "SWB-1"},
			},
			want: []ValidationError{
				{Line: 4, AssetID: "TX-BP", Label: "Invalid Parent Type", Message: "Asset TX-BP is a TRANSFORMER, so its parent must be a SUBSTATION, but SWB-1 is a SWITCHBOARD."},
			},
		},
		{
			name: "top-level ancestor is not a substation",
			assets: []Asset{
				{Line: 2, ID: "SWB-1", Type: "SWITCHBOARD"},
				{Line: 3, ID: "PNL-1", Type: "SWITCHBOARD_PANEL", ParentID: "SWB-1"},
			},
			want: []ValidationError{
				{Line: 3, AssetID: "PNL-1", Label: "No Substation Root", Message: "Asset PNL-1 does not resolve to a substation: its top-level ancestor SWB-1 is a SWITCHBOARD."},
			},
		},
		{
			name: "ancestor further up has a missing parent",
			assets: []Asset{
				{Line: 2, ID: "SWB-1", Type: "SWITCHBOARD", ParentID: "SS-404"},
				{Line: 3, ID: "PNL-1", Type: "SWITCHBOARD_PANEL", ParentID: "SWB-1"},
			},
			want: []ValidationError{
				{Line: 2, AssetID: "SWB-1", Label: "Parent Not Found", Message: "Asset SWB-1 has parent SS-404, which does not exist in the file."},
				{Line: 3, AssetID: "PNL-1", Label: "No Substation Root", Message: "Asset PNL-1 does not resolve to a substation: its ancestor SWB-1 has a parent that does not exist."},
			},
		},
		{
			name: "circular chain never reaches a substation",
			assets: []Asset{
				{Line: 2, ID: "SWB-A", Type: "SWITCHBOARD", ParentID: "SWB-B"},
				{Line: 3, ID: "SWB-B", Type: "SWITCHBOARD", ParentID: "SWB-A"},
				{Line: 4, ID: "PNL-1", Type: "SWITCHBOARD_PANEL", ParentID: "SWB-A"},
			},
			want: []ValidationError{
				{Line: 2, AssetID: "SWB-A", Label: "Invalid Parent Type", Message: "Asset SWB-A is a SWITCHBOARD, so its parent must be a SUBSTATION, but SWB-B is a SWITCHBOARD."},
				{Line: 2, AssetID: "SWB-A", Label: "No Substation Root", Message: "Asset SWB-A does not resolve to a substation: its parent chain is circular."},
				{Line: 3, AssetID: "SWB-B", Label: "Invalid Parent Type", Message: "Asset SWB-B is a SWITCHBOARD, so its parent must be a SUBSTATION, but SWB-A is a SWITCHBOARD."},
				{Line: 3, AssetID: "SWB-B", Label: "No Substation Root", Message: "Asset SWB-B does not resolve to a substation: its parent chain is circular."},
				{Line: 4, AssetID: "PNL-1", Label: "No Substation Root", Message: "Asset PNL-1 does not resolve to a substation: its parent chain is circular."},
			},
		},
		{
			name: "self-parenting asset",
			assets: []Asset{
				{Line: 2, ID: "SWB-1", Type: "SWITCHBOARD", ParentID: "SWB-1"},
			},
			want: []ValidationError{
				{Line: 2, AssetID: "SWB-1", Label: "Invalid Parent Type", Message: "Asset SWB-1 is a SWITCHBOARD, so its parent must be a SUBSTATION, but SWB-1 is a SWITCHBOARD."},
				{Line: 2, AssetID: "SWB-1", Label: "No Substation Root", Message: "Asset SWB-1 does not resolve to a substation: its parent chain is circular."},
			},
		},
		{
			name: "parent with no type",
			assets: []Asset{
				{Line: 2, ID: "X-1", Type: ""},
				{Line: 3, ID: "TX-1", Type: "TRANSFORMER", ParentID: "X-1"},
			},
			want: []ValidationError{
				{Line: 3, AssetID: "TX-1", Label: "Invalid Parent Type", Message: "Asset TX-1 is a TRANSFORMER, so its parent must be a SUBSTATION, but X-1 has no asset_type."},
				{Line: 3, AssetID: "TX-1", Label: "No Substation Root", Message: "Asset TX-1 does not resolve to a substation: its top-level ancestor X-1 has no asset_type."},
			},
		},
		{
			name: "unsupported type has no parent-type rule but must still resolve",
			assets: []Asset{
				{Line: 2, ID: "SS-1", Type: "SUBSTATION"},
				{Line: 3, ID: "GEN-1", Type: "GENERATOR", ParentID: "SS-1"},
				{Line: 4, ID: "GEN-2", Type: "GENERATOR", ParentID: "SS-404"},
			},
			want: []ValidationError{
				{Line: 4, AssetID: "GEN-2", Label: "Parent Not Found", Message: "Asset GEN-2 has parent SS-404, which does not exist in the file."},
			},
		},
		{
			name: "assets without a parent are left to the parent-requirements check",
			assets: []Asset{
				{Line: 2, ID: "TX-1", Type: "TRANSFORMER"},
				{Line: 3, ID: "SS-1", Type: "SUBSTATION"},
			},
			want: nil,
		},
	})
}

func TestCheckDuplicateSerials(t *testing.T) {
	runCheckCases(t, checkDuplicateSerials, []checkCase{
		{
			name: "all serial numbers unique",
			assets: []Asset{
				{Line: 2, ID: "TX-1", SerialNumber: "SER-1"},
				{Line: 3, ID: "TX-2", SerialNumber: "SER-2"},
			},
			want: nil,
		},
		{
			name: "two assets share a serial number and both are reported",
			assets: []Asset{
				{Line: 2, ID: "TX-D1", SerialNumber: "SER-TX-DUP"},
				{Line: 3, ID: "TX-3", SerialNumber: "SER-3"},
				{Line: 4, ID: "TX-D2", SerialNumber: "SER-TX-DUP"},
			},
			want: []ValidationError{
				{Line: 2, AssetID: "TX-D1", Label: "Duplicate Serial", Message: "Asset TX-D1 has duplicate serial_number SER-TX-DUP, also used on line 4."},
				{Line: 4, AssetID: "TX-D2", Label: "Duplicate Serial", Message: "Asset TX-D2 has duplicate serial_number SER-TX-DUP, also used on line 2."},
			},
		},
		{
			name: "empty serial numbers are not duplicates of each other",
			assets: []Asset{
				{Line: 2, ID: "TX-1", SerialNumber: ""},
				{Line: 3, ID: "TX-2", SerialNumber: ""},
			},
			want: nil,
		},
	})
}
