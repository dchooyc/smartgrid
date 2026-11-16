package api

import (
	"bytes"
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
