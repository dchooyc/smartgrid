// A minimal backend that validates uploaded asset CSV files.
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
	RowCount     int    `json:"row_count"`
	ColumnsValid bool   `json:"columns_valid"`
	Error        string `json:"error,omitempty"`
}

func main() {
	http.HandleFunc("POST /api/validate", handleValidate)

	log.Printf("listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func handleValidate(w http.ResponseWriter, r *http.Request) {
	// Let the frontend, which runs on a different port, read the response.
	w.Header().Set("Access-Control-Allow-Origin", "*")

	r.Body = http.MaxBytesReader(w, r.Body, maxUploadBytes)
	file, _, err := r.FormFile("file")
	if err != nil {
		var tooLarge *http.MaxBytesError
		if errors.As(err, &tooLarge) {
			writeError(w, http.StatusRequestEntityTooLarge, "the file is larger than 10 MB")
			return
		}
		writeError(w, http.StatusBadRequest, `no CSV file was uploaded in the "file" field`)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)

	header, err := reader.Read()
	if err == io.EOF {
		writeError(w, http.StatusBadRequest, "the file is empty")
		return
	}
	if err != nil {
		writeError(w, http.StatusBadRequest, "the file could not be read as CSV: "+err.Error())
		return
	}
	if err := checkColumns(header); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Count the data rows. The reader also rejects rows whose number of
	// fields differs from the header.
	rowCount := 0
	for {
		_, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			writeJSON(w, http.StatusBadRequest, validateResponse{
				ColumnsValid: true,
				Error:        "the file could not be read as CSV: " + err.Error(),
			})
			return
		}
		rowCount++
	}

	writeJSON(w, http.StatusOK, validateResponse{RowCount: rowCount, ColumnsValid: true})
}

// checkColumns reports an error unless the header holds exactly the required
// columns, in any order, with none missing, added or repeated.
func checkColumns(header []string) error {
	columns := make([]string, len(header))
	for i, title := range header {
		columns[i] = strings.TrimSpace(title)
	}
	// Excel puts a byte-order mark in front of the first title.
	columns[0] = strings.TrimPrefix(columns[0], "\uFEFF")

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

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, validateResponse{Error: message})
}

func writeJSON(w http.ResponseWriter, status int, body validateResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
}
