package api

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"smartgrid/backend/internal/store"
)

// A database address that nothing listens on, so every query fails at once.
const unreachableDatabase = "postgres://nobody@127.0.0.1:1/none?sslmode=disable&connect_timeout=1"

func TestDatabaseFailureIsLoggedWithItsCause(t *testing.T) {
	db, err := store.Open(unreachableDatabase)
	if err != nil {
		t.Fatal(err)
	}
	var output bytes.Buffer
	handler := NewHandler(db, slog.New(slog.NewTextHandler(&output, nil)))

	response := httptest.NewRecorder()
	handler.ListAssets(response, httptest.NewRequest("GET", "/api/assets", nil))

	if response.Code != http.StatusServiceUnavailable {
		t.Errorf("status = %d, want %d", response.Code, http.StatusServiceUnavailable)
	}
	line := output.String()
	for _, want := range []string{"level=ERROR", `msg="the database is not available"`, "path=/api/assets", "error=", "connect"} {
		if !strings.Contains(line, want) {
			t.Errorf("logged %q, want it to contain %q", line, want)
		}
	}
}

func TestFailureAndRequestLinesShareTheRequestID(t *testing.T) {
	db, err := store.Open(unreachableDatabase)
	if err != nil {
		t.Fatal(err)
	}
	var output bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&output, nil))
	handler := LogRequests(logger, http.HandlerFunc(NewHandler(db, logger).ListAssets))

	response := httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest("GET", "/api/assets", nil))

	id := response.Header().Get("X-Request-ID")
	lines := strings.Split(strings.TrimSpace(output.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("logged %d lines, want the failure and the request: %q", len(lines), output.String())
	}
	for _, line := range lines {
		if id == "" || !strings.Contains(line, "request_id="+id) {
			t.Errorf("logged %q, want it to contain request ID %q", line, id)
		}
	}

	// The response tells the user the same ID, to quote when reporting it.
	var body struct {
		RequestID string `json:"request_id"`
	}
	json.Unmarshal(response.Body.Bytes(), &body)
	if body.RequestID != id {
		t.Errorf("response request_id = %q, want %q", body.RequestID, id)
	}
}

func TestMistakeInRequestIsNotLoggedAsFailure(t *testing.T) {
	db, err := store.Open(unreachableDatabase)
	if err != nil {
		t.Fatal(err)
	}
	var output bytes.Buffer
	handler := NewHandler(db, slog.New(slog.NewTextHandler(&output, nil)))

	response := httptest.NewRecorder()
	handler.ListAssets(response, httptest.NewRequest("GET", "/api/assets?limit=abc", nil))

	if response.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", response.Code, http.StatusBadRequest)
	}
	if output.Len() != 0 {
		t.Errorf("logged %q, want nothing", output.String())
	}
}

func TestAllowOriginNamesOnlyTheFrontend(t *testing.T) {
	handler := AllowOrigin("http://localhost:3000", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusNotFound, "not found")
	}))

	response := httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest("GET", "/api/assets/missing", nil))

	if got := response.Header().Get("Access-Control-Allow-Origin"); got != "http://localhost:3000" {
		t.Errorf("Access-Control-Allow-Origin = %q, want the frontend's address", got)
	}
}
