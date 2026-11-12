package asset

import (
	"slices"
	"testing"
)

func TestSortParentsFirst(t *testing.T) {
	ids := func(assets []Asset) []string {
		var ids []string
		for _, a := range assets {
			ids = append(ids, a.ID)
		}
		return ids
	}

	cases := []struct {
		name   string
		assets []Asset
		want   []string
	}{
		{
			name: "children listed before their parents",
			assets: []Asset{
				{ID: "PNL-1", ParentID: "SWB-1"},
				{ID: "SWB-1", ParentID: "SS-1"},
				{ID: "SS-1"},
			},
			want: []string{"SS-1", "SWB-1", "PNL-1"},
		},
		{
			name: "already ordered rows keep their order",
			assets: []Asset{
				{ID: "SS-1"},
				{ID: "TX-1", ParentID: "SS-1"},
				{ID: "SS-2"},
				{ID: "TX-2", ParentID: "SS-2"},
			},
			want: []string{"SS-1", "SS-2", "TX-1", "TX-2"},
		},
		{
			name: "a parent outside the batch counts as already stored",
			assets: []Asset{
				{ID: "PNL-1", ParentID: "SWB-1"},
				{ID: "SWB-1", ParentID: "SS-IN-DATABASE"},
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

			got := ids(SortParentsFirst(tc.assets))
			if !slices.Equal(got, tc.want) {
				t.Errorf("order = %v, want %v", got, tc.want)
			}
			if !slices.Equal(ids(tc.assets), original) {
				t.Errorf("the input was reordered: %v", ids(tc.assets))
			}
		})
	}
}
