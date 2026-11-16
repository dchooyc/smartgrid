package api

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLogRequests(t *testing.T) {
	tests := []struct {
		name   string
		status int
		want   string // the level the request is logged at
	}{
		{"success is information", http.StatusOK, "level=INFO"},
		{"a mistake in the request is a warning", http.StatusBadRequest, "level=WARN"},
		{"a missing asset is a warning", http.StatusNotFound, "level=WARN"},
		{"a server failure is an error", http.StatusServiceUnavailable, "level=ERROR"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var output bytes.Buffer
			logger := slog.New(slog.NewTextHandler(&output, nil))
			handler := LogRequests(logger, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(test.status)
			}))

			request := httptest.NewRequest("GET", "/api/assets?limit=5", nil)
			handler.ServeHTTP(httptest.NewRecorder(), request)

			line := output.String()
			if strings.Count(line, "\n") != 1 {
				t.Fatalf("logged %q, want exactly one line", line)
			}
			for _, want := range []string{test.want, "method=GET", "path=/api/assets", `query="limit=5"`, "status=", "duration_ms="} {
				if !strings.Contains(line, want) {
					t.Errorf("logged %q, want it to contain %q", line, want)
				}
			}
		})
	}
}

func TestLogRequestsWithoutExplicitStatus(t *testing.T) {
	var output bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&output, nil))
	handler := LogRequests(logger, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok")) // no WriteHeader call means 200
	}))

	handler.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/api/assets", nil))

	if line := output.String(); !strings.Contains(line, "status=200") {
		t.Errorf("logged %q, want status=200", line)
	}
}
