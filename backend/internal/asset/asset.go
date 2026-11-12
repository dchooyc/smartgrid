// Package asset holds the asset model, its validation rules and its
// hierarchy logic.
package asset

import (
	"maps"
	"slices"
	"time"
)

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

// Stored is an asset as it is saved in the database, with all its fields.
// Fields that are empty in the database are null in JSON.
type Stored struct {
	ID               string    `json:"asset_id"`
	ParentID         *string   `json:"parent_asset_id"`
	Type             string    `json:"asset_type"`
	Name             string    `json:"asset_name"`
	Status           string    `json:"operational_status"`
	VoltageKV        *float64  `json:"voltage_kv"`
	RatingKVA        *float64  `json:"rating_kva"`
	Manufacturer     *string   `json:"manufacturer"`
	Model            *string   `json:"model"`
	SerialNumber     *string   `json:"serial_number"`
	CommissionedDate *string   `json:"commissioned_date"` // year-month-day
	CreatedAt        time.Time `json:"created_at"`
}

// SupportedTypes lists the asset types, in alphabetical order.
func SupportedTypes() []string {
	return slices.Sorted(maps.Keys(permittedParent))
}
