package main

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

// asset is one data row of the CSV, with surrounding spaces trimmed.
type asset struct {
	line             int // line number in the CSV file
	id               string
	parentID         string
	assetType        string
	name             string
	status           string
	voltageKV        string
	ratingKVA        string
	serialNumber     string
	commissionedDate string
}

// validationError describes a problem with one row. Warnings use it too.
type validationError struct {
	Line    int    `json:"line"`
	AssetID string `json:"asset_id,omitempty"`
	Message string `json:"message"`
}

// checks are the validations, in the order they run.
var checks = []func([]asset) []validationError{
	checkRequiredFields,
	checkUniqueIDs,
	checkSupportedValues,
	checkDatesAndNumbers,
	checkParentRequirements,
	checkCycles,
	checkHierarchy,
}

// validateAssets runs every check and returns all the errors found.
func validateAssets(assets []asset) []validationError {
	errs := []validationError{}
	for _, check := range checks {
		errs = append(errs, check(assets)...)
	}
	return errs
}

// findWarnings returns problems worth flagging that do not fail validation.
func findWarnings(assets []asset) []validationError {
	warnings := []validationError{}
	return append(warnings, checkDuplicateSerials(assets)...)
}

// 1. asset_id, asset_type and asset_name cannot be empty.
func checkRequiredFields(assets []asset) []validationError {
	var errs []validationError
	for _, a := range assets {
		if a.id == "" {
			errs = append(errs, newError(a, "is missing asset_id."))
		}
		if a.assetType == "" {
			errs = append(errs, newError(a, "is missing asset_type."))
		}
		if a.name == "" {
			errs = append(errs, newError(a, "is missing asset_name."))
		}
	}
	return errs
}

// 2. Every row that shares its asset_id with another row is reported.
func checkUniqueIDs(assets []asset) []validationError {
	linesByID := map[string][]int{}
	for _, a := range assets {
		if a.id != "" {
			linesByID[a.id] = append(linesByID[a.id], a.line)
		}
	}

	var errs []validationError
	for _, a := range assets {
		lines := linesByID[a.id]
		if len(lines) < 2 {
			continue
		}
		var others []string
		for _, line := range lines {
			if line != a.line {
				others = append(others, strconv.Itoa(line))
			}
		}
		errs = append(errs, newError(a, "has a duplicate asset_id, also used on line %s.", strings.Join(others, ", ")))
	}
	return errs
}

// 3. Asset types and operational statuses must be supported values.
func checkSupportedValues(assets []asset) []validationError {
	var errs []validationError
	for _, a := range assets {
		// A missing type is already reported by checkRequiredFields.
		if _, supported := permittedParent[a.assetType]; !supported && a.assetType != "" {
			errs = append(errs, newError(a, "has unsupported asset_type %q.", a.assetType))
		}
		if !slices.Contains(supportedStatuses, a.status) {
			errs = append(errs, newError(a, "has unsupported operational_status %q.", a.status))
		}
	}
	return errs
}

// 4. Dates and numbers must be readable when present, and rating_kva cannot
// be negative.
func checkDatesAndNumbers(assets []asset) []validationError {
	var errs []validationError
	for _, a := range assets {
		if a.commissionedDate != "" && !isDate(a.commissionedDate) {
			errs = append(errs, newError(a, "has invalid commissioned_date %q (expected day/month/year, e.g. 15/2/09).", a.commissionedDate))
		}
		if a.voltageKV != "" {
			if _, ok := parseNumber(a.voltageKV); !ok {
				errs = append(errs, newError(a, "has non-numeric voltage_kv %q.", a.voltageKV))
			}
		}
		if a.ratingKVA != "" {
			rating, ok := parseNumber(a.ratingKVA)
			if !ok {
				errs = append(errs, newError(a, "has non-numeric rating_kva %q.", a.ratingKVA))
			} else if rating < 0 {
				errs = append(errs, newError(a, "has negative rating_kva %s.", a.ratingKVA))
			}
		}
	}
	return errs
}

// 5. Substations cannot have a parent; every other asset must have one.
func checkParentRequirements(assets []asset) []validationError {
	var errs []validationError
	for _, a := range assets {
		switch {
		case a.assetType == substation && a.parentID != "":
			errs = append(errs, newError(a, "is a substation and cannot have a parent (%s).", a.parentID))
		case a.assetType != substation && a.assetType != "" && a.parentID == "":
			errs = append(errs, newError(a, "has no parent; only substations can be without one."))
		}
	}
	return errs
}

// 6. An asset cannot be its own parent or its own ancestor.
func checkCycles(assets []asset) []validationError {
	byID := assetsByID(assets)

	// cycles maps each asset ID that is part of a loop to the IDs in that loop.
	cycles := map[string][]string{}
	state := map[string]visitState{}
	for _, a := range assets {
		if a.id == "" || state[a.id] != unvisited {
			continue
		}
		cycle := findCycle(a.id, byID, state, nil)
		for _, id := range cycle {
			cycles[id] = cycle
		}
	}

	var errs []validationError
	for _, a := range assets {
		cycle, inCycle := cycles[a.id]
		if !inCycle || byID[a.id].line != a.line {
			continue
		}
		if len(cycle) == 1 {
			errs = append(errs, newError(a, "is its own parent."))
			continue
		}
		// Show the loop starting and ending at this asset.
		start := slices.Index(cycle, a.id)
		loop := slices.Concat(cycle[start:], cycle[:start], []string{a.id})
		errs = append(errs, newError(a, "is in a circular parent relationship: %s.", strings.Join(loop, " → ")))
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
func findCycle(id string, byID map[string]asset, state map[string]visitState, path []string) []string {
	state[id] = visiting
	defer func() { state[id] = visited }()
	path = append(path, id)

	parentID := byID[id].parentID
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
func checkHierarchy(assets []asset) []validationError {
	byID := assetsByID(assets)

	var errs []validationError
	for _, a := range assets {
		if a.parentID == "" {
			continue // checkParentRequirements covers assets without a parent
		}

		parent, exists := byID[a.parentID]
		if !exists {
			errs = append(errs, newError(a, "has parent %s, which does not exist in the file.", a.parentID))
			continue
		}

		// Substations and unsupported types are reported by earlier checks.
		wanted := permittedParent[a.assetType]
		if wanted != "" && parent.assetType != wanted {
			errs = append(errs, newError(a, "is a %s, so its parent must be a %s, but %s %s.",
				a.assetType, wanted, parent.id, describeType(parent)))
		}

		if problem := rootProblem(a, byID); problem != "" {
			errs = append(errs, newError(a, "does not resolve to a substation: %s.", problem))
		}
	}
	return errs
}

// rootProblem follows a's parents to the top of the hierarchy. It returns ""
// if that ends at a substation, or the reason it does not.
func rootProblem(a asset, byID map[string]asset) string {
	seen := map[string]bool{a.id: true}
	current := a
	for current.parentID != "" {
		parent, exists := byID[current.parentID]
		if !exists {
			return fmt.Sprintf("its ancestor %s has a parent that does not exist", current.id)
		}
		if seen[parent.id] {
			return "its parent chain is circular"
		}
		seen[parent.id] = true
		current = parent
	}

	if current.assetType != substation {
		return fmt.Sprintf("its top-level ancestor %s %s", current.id, describeType(current))
	}
	return ""
}

// Warning: every row that shares its serial_number with another row.
func checkDuplicateSerials(assets []asset) []validationError {
	linesBySerial := map[string][]int{}
	for _, a := range assets {
		if a.serialNumber != "" {
			linesBySerial[a.serialNumber] = append(linesBySerial[a.serialNumber], a.line)
		}
	}

	var warnings []validationError
	for _, a := range assets {
		lines := linesBySerial[a.serialNumber]
		if len(lines) < 2 {
			continue
		}
		var others []string
		for _, line := range lines {
			if line != a.line {
				others = append(others, strconv.Itoa(line))
			}
		}
		warnings = append(warnings, newError(a, "has duplicate serial_number %s, also used on line %s.", a.serialNumber, strings.Join(others, ", ")))
	}
	return warnings
}

// assetsByID indexes assets by ID. If an ID is duplicated, the first row wins.
func assetsByID(assets []asset) map[string]asset {
	byID := map[string]asset{}
	for _, a := range assets {
		if _, taken := byID[a.id]; a.id != "" && !taken {
			byID[a.id] = a
		}
	}
	return byID
}

// newError builds an error whose message starts by naming the asset.
func newError(a asset, format string, args ...any) validationError {
	name := "This asset"
	if a.id != "" {
		name = "Asset " + a.id
	}
	return validationError{
		Line:    a.line,
		AssetID: a.id,
		Message: name + " " + fmt.Sprintf(format, args...),
	}
}

// describeType completes a sentence about an asset, e.g. "X is a SWITCHBOARD".
func describeType(a asset) string {
	if a.assetType == "" {
		return "has no asset_type"
	}
	return "is a " + a.assetType
}

func isDate(value string) bool {
	for _, layout := range dateLayouts {
		if _, err := time.Parse(layout, value); err == nil {
			return true
		}
	}
	return false
}

func parseNumber(value string) (float64, bool) {
	number, err := strconv.ParseFloat(value, 64)
	if err != nil || math.IsNaN(number) || math.IsInf(number, 0) {
		return 0, false
	}
	return number, true
}
