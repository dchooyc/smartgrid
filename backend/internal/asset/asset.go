// Package asset holds the asset model, its validation rules and its
// hierarchy logic.
package asset

// Asset is one data row of the CSV, with surrounding spaces trimmed.
type Asset struct {
	Line             int // line number in the CSV file
	ID               string
	ParentID         string
	Type             string
	Name             string
	Status           string
	VoltageKV        string
	RatingKVA        string
	Manufacturer     string
	Model            string
	SerialNumber     string
	CommissionedDate string
	Record           []string // the original row, unmodified
}

// assetsByID indexes assets by ID. If an ID is duplicated, the first row wins.
func assetsByID(assets []Asset) map[string]Asset {
	byID := map[string]Asset{}
	for _, a := range assets {
		if _, taken := byID[a.ID]; a.ID != "" && !taken {
			byID[a.ID] = a
		}
	}
	return byID
}
