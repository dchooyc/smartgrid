package api

import (
	"errors"
	"net/http"

	"smartgrid/backend/internal/asset"
	"smartgrid/backend/internal/csvfile"
)

// Uploads are capped so the whole file always fits in memory and is never
// spilled to a temporary file on disk.
const maxUploadBytes = 10 << 20 // 10 MB

// upload is a CSV file received from the frontend.
type upload struct {
	filename string
	header   []string // as it appears in the file
	assets   []asset.Asset
}

// readUpload reads and parses the CSV in the request's "file" field. If the
// file cannot be used it writes the error response and returns false.
func readUpload(w http.ResponseWriter, r *http.Request) (upload, bool) {
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

	header, assets, err := csvfile.Read(file)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return upload{}, false
	}
	return upload{filename: fileHeader.Filename, header: header, assets: assets}, true
}
