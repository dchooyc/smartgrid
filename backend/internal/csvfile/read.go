// Package csvfile reads uploaded asset CSV files and generates the valid and
// invalid CSV files returned after validation.
package csvfile

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"slices"
	"strings"

	"smartgrid/backend/internal/asset"
)

// requiredColumns are the column titles every asset file must have.
var requiredColumns = []string{
	"asset_id",
	"parent_asset_id",
	"asset_type",
	"asset_name",
	"operational_status",
	"voltage_kv",
	"rating_kva",
	"manufacturer",
	"model",
	"serial_number",
	"commissioned_date",
}

// Read parses an asset CSV. It returns the header as it appears in the file
// and one asset per data row, or an error saying why the file cannot be used.
func Read(file io.Reader) (header []string, assets []asset.Asset, err error) {
	reader := csv.NewReader(file)

	header, err = reader.Read()
	if err == io.EOF {
		return nil, nil, errors.New("the file is empty")
	}
	if err != nil {
		return nil, nil, fmt.Errorf("the file could not be read as CSV: %w", err)
	}
	columns := cleanHeader(header)
	if err := checkColumns(columns); err != nil {
		return nil, nil, err
	}

	assets, err = readAssets(reader, columns)
	if err != nil {
		return nil, nil, fmt.Errorf("the file could not be read as CSV: %w", err)
	}
	return header, assets, nil
}

// cleanHeader trims the column titles and drops the byte-order mark that
// Excel puts in front of the first one.
func cleanHeader(header []string) []string {
	columns := make([]string, len(header))
	for i, title := range header {
		columns[i] = strings.TrimSpace(title)
	}
	columns[0] = strings.TrimPrefix(columns[0], "\uFEFF")
	return columns
}

// readAssets reads the data rows. The reader rejects rows whose number of
// fields differs from the header.
func readAssets(reader *csv.Reader, columns []string) ([]asset.Asset, error) {
	var assets []asset.Asset
	for {
		record, err := reader.Read()
		if err == io.EOF {
			return assets, nil
		}
		if err != nil {
			return nil, err
		}

		field := func(column string) string {
			return strings.TrimSpace(record[slices.Index(columns, column)])
		}
		line, _ := reader.FieldPos(0)
		assets = append(assets, asset.Asset{
			Line:             line,
			ID:               field("asset_id"),
			ParentID:         field("parent_asset_id"),
			Type:             field("asset_type"),
			Name:             field("asset_name"),
			Status:           field("operational_status"),
			VoltageKV:        field("voltage_kv"),
			RatingKVA:        field("rating_kva"),
			Manufacturer:     field("manufacturer"),
			Model:            field("model"),
			SerialNumber:     field("serial_number"),
			CommissionedDate: field("commissioned_date"),
			Record:           record,
		})
	}
}

// checkColumns reports an error unless the header holds exactly the required
// columns, in any order, with none missing, added or repeated.
func checkColumns(columns []string) error {
	var missing, unexpected, duplicate []string
	for _, column := range requiredColumns {
		if !slices.Contains(columns, column) {
			missing = append(missing, column)
		}
	}
	for i, column := range columns {
		switch {
		case !slices.Contains(requiredColumns, column):
			unexpected = append(unexpected, fmt.Sprintf("%q", column))
		case slices.Index(columns, column) != i && !slices.Contains(duplicate, column):
			duplicate = append(duplicate, column)
		}
	}

	var problems []string
	if len(missing) > 0 {
		problems = append(problems, "missing columns: "+strings.Join(missing, ", "))
	}
	if len(unexpected) > 0 {
		problems = append(problems, "unexpected columns: "+strings.Join(unexpected, ", "))
	}
	if len(duplicate) > 0 {
		problems = append(problems, "duplicate columns: "+strings.Join(duplicate, ", "))
	}
	if len(problems) > 0 {
		return errors.New("the file does not have the required columns (" + strings.Join(problems, "; ") + ")")
	}
	return nil
}
