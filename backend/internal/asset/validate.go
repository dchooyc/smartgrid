package asset

import (
	"fmt"
	"math"
	"slices"
	"strconv"
	"strings"
	"time"
)

const substation = "SUBSTATION"

// permittedParent maps each asset type to the type its parent must have.
// Substations have no parent.
var permittedParent = map[string]string{
	substation:          "",
	"TRANSFORMER":       substation,
	"LV_BOARD":          substation,
	"SWITCHBOARD":       substation,
	"SWITCHBOARD_PANEL": "SWITCHBOARD",
}

var supportedStatuses = []string{"IN_SERVICE", "MAINTENANCE", "OUT_OF_SERVICE"}

// Dates are day/month/year, with a two- or four-digit year.
var dateLayouts = []string{"2/1/06", "2/1/2006"}

// ValidationError describes a problem with one row. Warnings use it too.
type ValidationError struct {
	Type    string `json:"type"`  // the kind of problem, set by Validate
	Label   string `json:"label"` // short name, e.g. "Missing Parent"
	Line    int    `json:"line"`
	AssetID string `json:"asset_id,omitempty"`
	Message string `json:"message"`
}

// checks are the validations, in the order they run.
var checks = []struct {
	errorType string
	run       func([]Asset) []ValidationError
}{
	{"Missing required fields", checkRequiredFields},
	{"Duplicate asset IDs", checkUniqueIDs},
	{"Unsupported values", checkSupportedValues},
	{"Invalid dates and numbers", checkDatesAndNumbers},
	{"Parent requirements", checkParentRequirements},
	{"Circular relationships", checkCycles},
	{"Invalid hierarchy", checkHierarchy},
}

// Validate runs every check and returns all the errors found, each
// labelled with the type of the check that found it.
func Validate(assets []Asset) []ValidationError {
	errs := []ValidationError{}
	for _, check := range checks {
		errs = append(errs, withType(check.errorType, check.run(assets))...)
	}
	return errs
}

// FindWarnings returns problems worth flagging that do not fail validation.
func FindWarnings(assets []Asset) []ValidationError {
	warnings := []ValidationError{}
	return append(warnings, withType("Duplicate serial numbers", checkDuplicateSerials(assets))...)
}

func withType(errorType string, errs []ValidationError) []ValidationError {
	for i := range errs {
		errs[i].Type = errorType
	}
	return errs
}

// 1. asset_id, asset_type and asset_name cannot be empty.
func checkRequiredFields(assets []Asset) []ValidationError {
	var errs []ValidationError
	for _, a := range assets {
		if a.ID == "" {
			errs = append(errs, newError(a, "Missing Asset ID", "is missing asset_id."))
		}
		if a.Type == "" {
			errs = append(errs, newError(a, "Missing Asset Type", "is missing asset_type."))
		}
		if a.Name == "" {
			errs = append(errs, newError(a, "Missing Asset Name", "is missing asset_name."))
		}
	}
	return errs
}

// 2. Every row that shares its asset_id with another row is reported.
func checkUniqueIDs(assets []Asset) []ValidationError {
	linesByID := map[string][]int{}
	for _, a := range assets {
		if a.ID != "" {
			linesByID[a.ID] = append(linesByID[a.ID], a.Line)
		}
	}

	var errs []ValidationError
	for _, a := range assets {
		lines := linesByID[a.ID]
		if len(lines) < 2 {
			continue
		}
		var others []string
		for _, line := range lines {
			if line != a.Line {
				others = append(others, strconv.Itoa(line))
			}
		}
		errs = append(errs, newError(a, "Duplicate ID", "has a duplicate asset_id, also used on line %s.", strings.Join(others, ", ")))
	}
	return errs
}

// 3. Asset types and operational statuses must be supported values.
func checkSupportedValues(assets []Asset) []ValidationError {
	var errs []ValidationError
	for _, a := range assets {
		// A missing type is already reported by checkRequiredFields.
		if _, supported := permittedParent[a.Type]; !supported && a.Type != "" {
			errs = append(errs, newError(a, "Invalid Asset Type", "has unsupported asset_type %q.", a.Type))
		}
		if !slices.Contains(supportedStatuses, a.Status) {
			errs = append(errs, newError(a, "Invalid Status", "has unsupported operational_status %q.", a.Status))
		}
	}
	return errs
}

// 4. Dates and numbers must be readable when present, and rating_kva cannot
// be negative.
func checkDatesAndNumbers(assets []Asset) []ValidationError {
	var errs []ValidationError
	for _, a := range assets {
		if a.CommissionedDate != "" && !isDate(a.CommissionedDate) {
			errs = append(errs, newError(a, "Invalid Date", "has invalid commissioned_date %q (expected day/month/year, e.g. 15/2/09).", a.CommissionedDate))
		}
		if a.VoltageKV != "" {
			if _, ok := ParseNumber(a.VoltageKV); !ok {
				errs = append(errs, newError(a, "Invalid Voltage", "has non-numeric voltage_kv %q.", a.VoltageKV))
			}
		}
		if a.RatingKVA != "" {
			rating, ok := ParseNumber(a.RatingKVA)
			if !ok {
				errs = append(errs, newError(a, "Invalid Rating", "has non-numeric rating_kva %q.", a.RatingKVA))
			} else if rating < 0 {
				errs = append(errs, newError(a, "Negative Rating", "has negative rating_kva %s.", a.RatingKVA))
			}
		}
	}
	return errs
}

// 5. Substations cannot have a parent; every other asset must have one.
func checkParentRequirements(assets []Asset) []ValidationError {
	var errs []ValidationError
	for _, a := range assets {
		switch {
		case a.Type == substation && a.ParentID != "":
			errs = append(errs, newError(a, "Substation Has Parent", "is a substation and cannot have a parent (%s).", a.ParentID))
		case a.Type != substation && a.Type != "" && a.ParentID == "":
			errs = append(errs, newError(a, "Missing Parent", "has no parent; only substations can be without one."))
		}
	}
	return errs
}

// 6. An asset cannot be its own parent or its own ancestor.
func checkCycles(assets []Asset) []ValidationError {
	byID := assetsByID(assets)

	// cycles maps each asset ID that is part of a loop to the IDs in that loop.
	cycles := map[string][]string{}
	state := map[string]visitState{}
	for _, a := range assets {
		if a.ID == "" || state[a.ID] != unvisited {
			continue
		}
		cycle := findCycle(a.ID, byID, state, nil)
		for _, id := range cycle {
			cycles[id] = cycle
		}
	}

	var errs []ValidationError
	for _, a := range assets {
		cycle, inCycle := cycles[a.ID]
		if !inCycle || byID[a.ID].Line != a.Line {
			continue
		}
		if len(cycle) == 1 {
			errs = append(errs, newError(a, "Self Parent", "is its own parent."))
			continue
		}
		// Show the loop starting and ending at this asset.
		start := slices.Index(cycle, a.ID)
		loop := slices.Concat(cycle[start:], cycle[:start], []string{a.ID})
		errs = append(errs, newError(a, "Cycle Detected", "is in a circular parent relationship: %s.", strings.Join(loop, " → ")))
	}
	return errs
}

type visitState int

const (
	unvisited visitState = iota
	visiting             // on the path currently being walked
	visited
)

// findCycle is a depth-first search up the parent chain from id. It returns
// the IDs forming a loop if the chain runs back into the current path, or nil.
func findCycle(id string, byID map[string]Asset, state map[string]visitState, path []string) []string {
	state[id] = visiting
	defer func() { state[id] = visited }()
	path = append(path, id)

	parentID := byID[id].ParentID
	if _, exists := byID[parentID]; !exists {
		return nil
	}

	switch state[parentID] {
	case visiting:
		return slices.Clone(path[slices.Index(path, parentID):])
	case unvisited:
		return findCycle(parentID, byID, state, path)
	default:
		return nil // already fully explored
	}
}

// 7. Every parent must exist, be of the permitted type, and lead up to a
// substation.
func checkHierarchy(assets []Asset) []ValidationError {
	byID := assetsByID(assets)

	var errs []ValidationError
	for _, a := range assets {
		if a.ParentID == "" {
			continue // checkParentRequirements covers assets without a parent
		}

		parent, exists := byID[a.ParentID]
		if !exists {
			errs = append(errs, newError(a, "Parent Not Found", "has parent %s, which does not exist in the file.", a.ParentID))
			continue
		}

		// Substations and unsupported types are reported by earlier checks.
		wanted := permittedParent[a.Type]
		if wanted != "" && parent.Type != wanted {
			errs = append(errs, newError(a, "Invalid Parent Type", "is a %s, so its parent must be a %s, but %s %s.",
				a.Type, wanted, parent.ID, describeType(parent)))
		}

		if problem := rootProblem(a, byID); problem != "" {
			errs = append(errs, newError(a, "No Substation Root", "does not resolve to a substation: %s.", problem))
		}
	}
	return errs
}

// rootProblem follows a's parents to the top of the hierarchy. It returns ""
// if that ends at a substation, or the reason it does not.
func rootProblem(a Asset, byID map[string]Asset) string {
	seen := map[string]bool{a.ID: true}
	current := a
	for current.ParentID != "" {
		parent, exists := byID[current.ParentID]
		if !exists {
			return fmt.Sprintf("its ancestor %s has a parent that does not exist", current.ID)
		}
		if seen[parent.ID] {
			return "its parent chain is circular"
		}
		seen[parent.ID] = true
		current = parent
	}

	if current.Type != substation {
		return fmt.Sprintf("its top-level ancestor %s %s", current.ID, describeType(current))
	}
	return ""
}

// Warning: every row that shares its serial_number with another row.
func checkDuplicateSerials(assets []Asset) []ValidationError {
	linesBySerial := map[string][]int{}
	for _, a := range assets {
		if a.SerialNumber != "" {
			linesBySerial[a.SerialNumber] = append(linesBySerial[a.SerialNumber], a.Line)
		}
	}

	var warnings []ValidationError
	for _, a := range assets {
		lines := linesBySerial[a.SerialNumber]
		if len(lines) < 2 {
			continue
		}
		var others []string
		for _, line := range lines {
			if line != a.Line {
				others = append(others, strconv.Itoa(line))
			}
		}
		warnings = append(warnings, newError(a, "Duplicate Serial", "has duplicate serial_number %s, also used on line %s.", a.SerialNumber, strings.Join(others, ", ")))
	}
	return warnings
}

// newError builds an error whose message starts by naming the asset. The
// label is a short name for the problem, e.g. "Missing Parent".
func newError(a Asset, label, format string, args ...any) ValidationError {
	name := "This asset"
	if a.ID != "" {
		name = "Asset " + a.ID
	}
	return ValidationError{
		Label:   label,
		Line:    a.Line,
		AssetID: a.ID,
		Message: name + " " + fmt.Sprintf(format, args...),
	}
}

// describeType completes a sentence about an asset, e.g. "X is a SWITCHBOARD".
func describeType(a Asset) string {
	if a.Type == "" {
		return "has no asset_type"
	}
	return "is a " + a.Type
}

func isDate(value string) bool {
	_, ok := ParseDate(value)
	return ok
}

func ParseDate(value string) (time.Time, bool) {
	for _, layout := range dateLayouts {
		if date, err := time.Parse(layout, value); err == nil {
			return date, true
		}
	}
	return time.Time{}, false
}

func ParseNumber(value string) (float64, bool) {
	number, err := strconv.ParseFloat(value, 64)
	if err != nil || math.IsNaN(number) || math.IsInf(number, 0) {
		return 0, false
	}
	return number, true
}
