package main

import (
	"fmt"
	"slices"
)

// finding is something the app should report about a row of the file.
type finding struct {
	line    int
	assetID string
	problem string   // what was done to the row
	labels  []string // the errors the app should show for it
	warning bool     // a warning, which does not make the row invalid
}

// kind is one way of breaking rows, matching one of the app's validations.
// They are listed in the order the app checks them.
type kind struct {
	flag  string // command-line argument, and the folder its files go in
	code  string // short name used in file names
	title string
	apply func(f *file, count int) error
}

var kinds = []kind{
	{"missing-id", "NOID", "Missing asset ID", func(f *file, count int) error {
		return f.change(count, isEquipment, "Missing asset ID", []string{"Missing Asset ID"},
			func(r *row) { r.id = "" })
	}},
	{"missing-type", "NOTYPE", "Missing asset type", func(f *file, count int) error {
		return f.change(count, isEquipment, "Missing asset type", []string{"Missing Asset Type"},
			func(r *row) { r.assetType = "" })
	}},
	{"missing-name", "NONAME", "Missing asset name", func(f *file, count int) error {
		return f.change(count, isEquipment, "Missing asset name", []string{"Missing Asset Name"},
			func(r *row) { r.name = "" })
	}},
	{"duplicate-id", "DUPID", "Duplicate asset ID", func(f *file, count int) error {
		return f.changePairs(count, "Duplicate asset ID", func(first, second *row) {
			second.id = first.id
			f.report(first, "Asset ID repeated on line "+fmt.Sprint(second.line), "Duplicate ID")
			f.report(second, "Asset ID copied from line "+fmt.Sprint(first.line), "Duplicate ID")
		})
	}},
	{"invalid-type", "BADTYPE", "Unsupported asset type", func(f *file, count int) error {
		return f.change(count, isEquipment, "Unsupported asset type", []string{"Invalid Asset Type"},
			func(r *row) { r.assetType = "GENERATOR" })
	}},
	{"invalid-status", "BADSTATUS", "Unsupported status", func(f *file, count int) error {
		return f.change(count, isEquipment, "Unsupported status", []string{"Invalid Status"},
			func(r *row) { r.status = "BROKEN" })
	}},
	{"invalid-date", "BADDATE", "Invalid date", func(f *file, count int) error {
		return f.change(count, isEquipment, "Invalid date", []string{"Invalid Date"},
			func(r *row) { r.date = "31/2/20" })
	}},
	{"invalid-voltage", "BADVOLT", "Invalid voltage", func(f *file, count int) error {
		return f.change(count, isEquipment, "Invalid voltage", []string{"Invalid Voltage"},
			func(r *row) { r.voltage = "abc" })
	}},
	{"invalid-rating", "BADRATING", "Invalid rating", func(f *file, count int) error {
		return f.change(count, isType(transformer), "Invalid rating", []string{"Invalid Rating"},
			func(r *row) { r.rating = "abc" })
	}},
	{"negative-rating", "NEGRATING", "Negative rating", func(f *file, count int) error {
		return f.change(count, isType(transformer), "Negative rating", []string{"Negative Rating"},
			func(r *row) { r.rating = "-500" })
	}},
	{"substation-parent", "SSPARENT", "Substation with a parent", substationParents},
	{"no-parent", "NOPARENT", "No parent", func(f *file, count int) error {
		return f.change(count, isEquipment, "No parent", []string{"Missing Parent"},
			func(r *row) { r.parentID = "" })
	}},
	{"parent-not-found", "ORPHAN", "Parent not found", func(f *file, count int) error {
		return f.change(count, isEquipment, "Parent not found", []string{"Parent Not Found"},
			func(r *row) { r.parentID = "SWB-404-1-GEN-ID-" + f.hash })
	}},
	{"self-parent", "SELF", "Its own parent", func(f *file, count int) error {
		return f.change(count, isEquipment, "Its own parent",
			[]string{"Self Parent", "Invalid Parent Type", "No Substation Root"},
			func(r *row) { r.parentID = r.id })
	}},
	{"cycle", "CYCLE", "Circular parents", func(f *file, count int) error {
		return f.changePairs(count, "Circular parents", func(first, second *row) {
			first.parentID, second.parentID = second.id, first.id
			labels := []string{"Cycle Detected", "Invalid Parent Type", "No Substation Root"}
			f.report(first, "Circular parents with line "+fmt.Sprint(second.line), labels...)
			f.report(second, "Circular parents with line "+fmt.Sprint(first.line), labels...)
		})
	}},
	{"wrong-parent-type", "BADPARENT", "Wrong type of parent", wrongParentTypes},
	{"no-substation-root", "NOROOT", "No substation root", func(f *file, count int) error {
		return f.changePairs(count, "No substation root", func(top, below *row) {
			top.parentID = ""
			below.parentID = top.id
			f.report(top, "Parent removed, to head a chain without a substation", "Missing Parent")
			f.report(below, "Moved under line "+fmt.Sprint(top.line)+", which leads to no substation",
				"Invalid Parent Type", "No Substation Root")
		})
	}},
	{"duplicate-serial", "DUPSERIAL", "Duplicate serial number (warning)", func(f *file, count int) error {
		return f.changePairs(count, "Duplicate serial number", func(first, second *row) {
			second.serial = first.serial
			f.warn(first, "Serial number repeated on line "+fmt.Sprint(second.line))
			f.warn(second, "Serial number copied from line "+fmt.Sprint(first.line))
		})
	}},
}

// isEquipment selects transformers, LV boards and panels: the assets that
// are never anyone's parent in a valid file, so breaking one cannot cause
// errors in other rows.
func isEquipment(r *row) bool {
	return r.assetType == transformer || r.assetType == lvBoard || r.assetType == switchboardPanel
}

func isType(assetType string) func(*row) bool {
	return func(r *row) bool { return r.assetType == assetType }
}

// pick chooses count untouched rows without children that match, spread
// through the file, and marks them as used.
func (f *file) pick(count int, match func(*row) bool, title string) ([]*row, error) {
	var candidates []*row
	for _, r := range f.rows {
		if !r.used && r.children == 0 && match(r) {
			candidates = append(candidates, r)
		}
	}
	if len(candidates) < count {
		return nil, fmt.Errorf("%q needs %d suitable rows but only %d are left: generate more assets or fewer errors",
			title, count, len(candidates))
	}

	picked := make([]*row, count)
	for i := range picked {
		picked[i] = candidates[i*len(candidates)/count]
		picked[i].used = true
	}
	return picked, nil
}

// change breaks count rows in the same way, each with the same result.
func (f *file) change(count int, match func(*row) bool, problem string, labels []string, apply func(*row)) error {
	rows, err := f.pick(count, match, problem)
	if err != nil {
		return err
	}
	for _, r := range rows {
		apply(r)
		f.report(r, problem, labels...)
	}
	return nil
}

// changePairs breaks count pairs of rows; apply reports what it did.
func (f *file) changePairs(count int, title string, apply func(first, second *row)) error {
	rows, err := f.pick(2*count, isEquipment, title)
	if err != nil {
		return err
	}
	for i := 0; i < len(rows); i += 2 {
		apply(rows[i], rows[i+1])
	}
	return nil
}

// substationParents puts substations under the first substation, which is
// left alone so that no cycle forms.
func substationParents(f *file, count int) error {
	var substations []*row
	for _, r := range f.rows {
		if r.assetType == substation {
			substations = append(substations, r)
		}
	}

	var free []*row
	for _, r := range substations[1:] {
		if !r.used {
			free = append(free, r)
		}
	}
	if len(free) < count {
		return fmt.Errorf("%q needs %d substations besides the first but only %d are left: generate more substations",
			"Substation with a parent", count, len(free))
	}
	for _, r := range free[len(free)-count:] {
		r.used = true
		r.parentID = substations[0].id
		f.report(r, "Substation with a parent", "Substation Has Parent")
	}
	return nil
}

// wrongParentTypes moves panels from their switchboard to its substation.
func wrongParentTypes(f *file, count int) error {
	parentOf := map[string]*row{}
	for _, r := range f.rows {
		parentOf[r.id] = r
	}
	return f.change(count, isType(switchboardPanel), "Wrong type of parent", []string{"Invalid Parent Type"},
		func(r *row) { r.parentID = parentOf[r.parentID].parentID })
}

func (f *file) report(r *row, problem string, labels ...string) {
	f.findings = append(f.findings, finding{line: r.line, assetID: r.id, problem: problem, labels: labels})
}

func (f *file) warn(r *row, problem string) {
	f.findings = append(f.findings, finding{
		line: r.line, assetID: r.id, problem: problem,
		labels: []string{"Duplicate Serial"}, warning: true,
	})
}

// sortedFindings lists the findings in the order of the file.
func (f *file) sortedFindings() []finding {
	sorted := slices.Clone(f.findings)
	slices.SortStableFunc(sorted, func(a, b finding) int { return a.line - b.line })
	return sorted
}
