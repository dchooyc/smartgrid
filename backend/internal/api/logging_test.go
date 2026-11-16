package api

import (
	"bytes"
	"encoding/json"
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
			response := httptest.NewRecorder()
			handler.ServeHTTP(response, request)

			line := output.String()
			id := response.Header().Get("X-Request-ID")
			if len(id) != 8 {
				t.Errorf("X-Request-ID = %q, want 8 characters", id)
			}
			if !strings.Contains(line, "request_id="+id) {
				t.Errorf("logged %q, want it to contain the response's request ID %q", line, id)
			}
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

func TestPanicIsLoggedAndAnsweredWithServerError(t *testing.T) {
	var output bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&output, nil))
	handler := LogRequests(logger, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("something impossible happened")
	}))

	response := httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest("GET", "/api/assets", nil))

	if response.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", response.Code, http.StatusInternalServerError)
	}
	id := response.Header().Get("X-Request-ID")
	var body struct {
		RequestID string `json:"request_id"`
	}
	json.Unmarshal(response.Body.Bytes(), &body)
	if body.RequestID != id {
		t.Errorf("response request_id = %q, want %q", body.RequestID, id)
	}

	lines := strings.Split(strings.TrimSpace(output.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("logged %d lines, want the panic and the request: %q", len(lines), output.String())
	}
	for _, want := range []string{"level=ERROR", "msg=panic", "request_id=" + id, "something impossible happened", "stack="} {
		if !strings.Contains(lines[0], want) {
			t.Errorf("logged %q, want it to contain %q", lines[0], want)
		}
	}
	for _, want := range []string{"level=ERROR", "msg=request", "request_id=" + id, "status=500"} {
		if !strings.Contains(lines[1], want) {
			t.Errorf("logged %q, want it to contain %q", lines[1], want)
		}
	}
}
