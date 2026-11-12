package csvfile

import (
	"bytes"
	"encoding/csv"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"smartgrid/backend/internal/asset"
)

const errorsColumn = "Validation Errors (Remove Before Reimporting)"

// File is a generated CSV, held in memory and sent to the frontend.
type File struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

// BuildFiles splits the uploaded rows into two CSV files using the errors
// validation already found: rows without errors go to the valid file, and the
// rest go to the invalid file with their error labels in a new first column.
// Every row lands in exactly one file, in its original order.
func BuildFiles(originalName string, header []string, assets []asset.Asset, errs []asset.ValidationError, now time.Time) (valid, invalid File, err error) {
	labels := labelsByLine(errs)

	// Excel puts a byte-order mark in front of the first title.
	header = slices.Clone(header)
	header[0] = strings.TrimPrefix(header[0], "\uFEFF")

	validRows := [][]string{header}
	invalidRows := [][]string{slices.Concat([]string{errorsColumn}, header)}
	for _, a := range assets {
		if rowLabels, hasErrors := labels[a.Line]; hasErrors {
			errorsCell := strings.Join(rowLabels, "; ")
			invalidRows = append(invalidRows, slices.Concat([]string{errorsCell}, a.Record))
		} else {
			validRows = append(validRows, a.Record)
		}
	}

	if valid.Content, err = writeCSV(validRows); err != nil {
		return valid, invalid, err
	}
	if invalid.Content, err = writeCSV(invalidRows); err != nil {
		return valid, invalid, err
	}

	// e.g. grid_assets.csv -> grid_assets_valid_20260115_093000.csv
	name := filepath.Base(originalName)
	if strings.EqualFold(filepath.Ext(name), ".csv") {
		name = name[:len(name)-len(".csv")]
	}
	timestamp := now.Format("20060102_150405")
	valid.Name = name + "_valid_" + timestamp + ".csv"
	invalid.Name = name + "_invalid_" + timestamp + ".csv"
	return valid, invalid, nil
}

// labelsByLine lists the distinct error labels for each CSV line, in the
// order the errors were found.
func labelsByLine(errs []asset.ValidationError) map[int][]string {
	labels := map[int][]string{}
	for _, e := range errs {
		if !slices.Contains(labels[e.Line], e.Label) {
			labels[e.Line] = append(labels[e.Line], e.Label)
		}
	}
	return labels
}

func writeCSV(rows [][]string) (string, error) {
	var buffer bytes.Buffer
	if err := csv.NewWriter(&buffer).WriteAll(rows); err != nil {
		return "", err
	}
	return buffer.String(), nil
}
