// A minimal backend that validates uploaded asset CSV files and commits the
// valid rows to PostgreSQL.
package main

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"slices"
	"strings"
	"time"
)

const (
	addr = ":8080"

	// Uploads are capped so the whole file always fits in memory and is
	// never spilled to a temporary file on disk.
	maxUploadBytes = 10 << 20 // 10 MB
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

type validateResponse struct {
	Valid     bool              `json:"valid"`
	TotalRows int               `json:"total_rows"`
	Errors    []validationError `json:"errors"`
	Warnings  []validationError `json:"warnings"`
	// ValidFile holds the rows without errors and InvalidFile the rest.
	ValidFile   *csvFile `json:"valid_file,omitempty"`
	InvalidFile *csvFile `json:"invalid_file,omitempty"`
}

func main() {
	var err error
	if db, err = openDatabase(); err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("POST /api/validate", handleValidate)
	http.HandleFunc("POST /api/imports", handleImport)

	log.Printf("listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func handleValidate(w http.ResponseWriter, r *http.Request) {
	upload, ok := readUpload(w, r)
	if !ok {
		return
	}
	header, assets := upload.header, upload.assets

	// Check against the database as well, so that rows which could not be
	// committed are reported now rather than when committing.
	existing, err := loadExistingAssets(r.Context(), db)
	if err != nil {
		writeError(w, http.StatusServiceUnavailable, "the database is not available: "+err.Error())
		return
	}

	errs := validateWithStored(assets, existing)
	validFile, invalidFile, err := buildCSVFiles(upload.filename, header, assets, errs, time.Now())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "the result files could not be generated: "+err.Error())
		return
	}

	writeJSON(w, http.StatusOK, validateResponse{
		Valid:       len(errs) == 0,
		TotalRows:   len(assets),
		Errors:      errs,
		Warnings:    findWarnings(assets),
		ValidFile:   &validFile,
		InvalidFile: &invalidFile,
	})
}

// upload is a CSV file received from the frontend.
type upload struct {
	filename string
	header   []string // as it appears in the file
	assets   []asset
}

// readUpload reads and parses the CSV in the request's "file" field. If the
// file cannot be used it writes the error response and returns false.
func readUpload(w http.ResponseWriter, r *http.Request) (upload, bool) {
	// Let the frontend, which runs on a different port, read the response.
	w.Header().Set("Access-Control-Allow-Origin", "*")

	r.Body = http.MaxBytesReader(w, r.Body, maxUploadBytes)
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		var tooLarge *http.MaxBytesError
		if errors.As(err, &tooLarge) {
			writeError(w, http.StatusRequestEntityTooLarge, "the file is larger than 10 MB")
			return upload{}, false
		}
		writeError(w, http.StatusBadRequest, `no CSV file was uploaded in the "file" field`)
		return upload{}, false
	}
	defer file.Close()

	reader := csv.NewReader(file)

	header, err := reader.Read()
	if err == io.EOF {
		writeError(w, http.StatusBadRequest, "the file is empty")
		return upload{}, false
	}
	if err != nil {
		writeError(w, http.StatusBadRequest, "the file could not be read as CSV: "+err.Error())
		return upload{}, false
	}
	columns := cleanHeader(header)
	if err := checkColumns(columns); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return upload{}, false
	}

	assets, err := readAssets(reader, columns)
	if err != nil {
		writeError(w, http.StatusBadRequest, "the file could not be read as CSV: "+err.Error())
		return upload{}, false
	}

	return upload{filename: fileHeader.Filename, header: header, assets: assets}, true
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
func readAssets(reader *csv.Reader, columns []string) ([]asset, error) {
	var assets []asset
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
		assets = append(assets, asset{
			line:             line,
			id:               field("asset_id"),
			parentID:         field("parent_asset_id"),
			assetType:        field("asset_type"),
			name:             field("asset_name"),
			status:           field("operational_status"),
			voltageKV:        field("voltage_kv"),
			ratingKVA:        field("rating_kva"),
			manufacturer:     field("manufacturer"),
			model:            field("model"),
			serialNumber:     field("serial_number"),
			commissionedDate: field("commissioned_date"),
			record:           record,
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

// writeError reports that a request failed before producing any result.
func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]any{
		"error":  message,
		"errors": []validationError{},
	})
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
}
