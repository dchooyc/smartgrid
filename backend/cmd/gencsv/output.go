package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"slices"
	"strings"
)

// columnChanges are the ways of breaking the header. Any of them makes the
// app reject the whole file.
type columnChanges struct {
	missing   int // columns removed
	duplicate int // columns repeated
	extra     int // columns the app does not know added
}

func (c columnChanges) any() bool {
	return c.missing+c.duplicate+c.extra > 0
}

// table returns the header and rows of the file with the column changes
// made, and the message the app should reject it with ("" if none).
func (f *file) table(changes columnChanges) (records [][]string, rejection string) {
	// Columns are removed from the end and repeated from the start, so the
	// two changes never touch the same column.
	kept := columns[:len(columns)-changes.missing]
	repeated := kept[:min(changes.duplicate, len(kept))]

	header := slices.Concat(kept, repeated)
	var added []string
	for i := range changes.extra {
		added = append(added, fmt.Sprintf("extra_column_%d", i+1))
	}
	header = append(header, added...)
	records = append(records, header)

	for _, r := range f.rows {
		record := r.record()
		record = slices.Concat(record[:len(kept)], record[:len(repeated)])
		for range added {
			record = append(record, "extra")
		}
		records = append(records, record)
	}

	if !changes.any() {
		return records, ""
	}
	var problems []string
	if removed := columns[len(kept):]; len(removed) > 0 {
		problems = append(problems, "missing columns: "+strings.Join(removed, ", "))
	}
	if len(added) > 0 {
		quoted := make([]string, len(added))
		for i, column := range added {
			quoted[i] = fmt.Sprintf("%q", column)
		}
		problems = append(problems, "unexpected columns: "+strings.Join(quoted, ", "))
	}
	if len(repeated) > 0 {
		problems = append(problems, "duplicate columns: "+strings.Join(repeated, ", "))
	}
	return records, "the file does not have the required columns (" + strings.Join(problems, "; ") + ")"
}

func writeCSV(records [][]string) (string, error) {
	var buffer bytes.Buffer
	if err := csv.NewWriter(&buffer).WriteAll(records); err != nil {
		return "", err
	}
	return buffer.String(), nil
}

// answerSheet says what the app should report when the file is uploaded.
func (f *file) answerSheet(csvName string, substations int, rejection string) string {
	var sheet strings.Builder
	fmt.Fprintf(&sheet, "Answer sheet for %s\n\n", csvName)
	fmt.Fprintf(&sheet, "Substations: %d\nAssets (rows): %d\n\n", substations, len(f.rows))

	if rejection != "" {
		sheet.WriteString("The app should reject the whole file with:\n\n  " + rejection + "\n")
		if len(f.findings) > 0 {
			sheet.WriteString("\nThe rows also have problems, but the app never gets as far as checking them.\n")
		}
		return sheet.String()
	}

	findings := f.sortedFindings()
	invalid := map[int]bool{}
	warnings := 0
	for _, found := range findings {
		if found.warning {
			warnings++
		} else {
			invalid[found.line] = true
		}
	}
	fmt.Fprintf(&sheet, "The app should report:\n")
	fmt.Fprintf(&sheet, "  Lines that can be committed: %d\n", len(f.rows)-len(invalid))
	fmt.Fprintf(&sheet, "  Lines with errors:           %d\n", len(invalid))
	fmt.Fprintf(&sheet, "  Warnings:                    %d\n", warnings)
	if len(findings) == 0 {
		sheet.WriteString("\nEvery row is valid.\n")
		return sheet.String()
	}

	sheet.WriteString("\nThe lines below were changed. \"Expects\" lists what the app should show for each.\n")
	for _, found := range findings {
		id := found.assetID
		if id == "" {
			id = "(no asset ID)"
		}
		expects := "errors"
		if found.warning {
			expects = "warning"
		}
		fmt.Fprintf(&sheet, "\nLine %d  %s\n", found.line, id)
		fmt.Fprintf(&sheet, "  Changed: %s\n", found.problem)
		fmt.Fprintf(&sheet, "  Expects %s: %s\n", expects, strings.Join(found.labels, "; "))
	}
	return sheet.String()
}
