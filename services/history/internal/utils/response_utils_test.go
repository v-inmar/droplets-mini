package utils

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestResponseJSON(t *testing.T) {
	tests := []struct {
		name             string
		data             any
		status           int
		headers          http.Header
		wantStatus       int
		wantBody         string
		wantContentType  string
		wantCustomHeader string
	}{
		{
			name:            "writes JSON response",
			data:            map[string]string{"message": "hello"},
			status:          http.StatusOK,
			wantStatus:      http.StatusOK,
			wantBody:        "{\n\t\"message\": \"hello\"\n}",
			wantContentType: "application/json",
		},
		{
			name:   "copies custom headers",
			data:   map[string]int{"count": 42},
			status: http.StatusCreated,
			headers: http.Header{
				"X-Custom-Header": []string{"custom-value"},
			},
			wantStatus:       http.StatusCreated,
			wantBody:         "{\n\t\"count\": 42\n}",
			wantContentType:  "application/json",
			wantCustomHeader: "custom-value",
		},
		{
			name:   "overrides Content-Type",
			data:   map[string]bool{"ok": true},
			status: http.StatusOK,
			headers: http.Header{
				"Content-Type":    []string{"text/plain"},
				"X-Custom-Header": []string{"custom-value"},
			},
			wantStatus:       http.StatusOK,
			wantBody:         "{\n\t\"ok\": true\n}",
			wantContentType:  "application/json",
			wantCustomHeader: "custom-value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			err := ResponseJSON(w, tt.data, tt.status, tt.headers)
			if err != nil {
				t.Fatalf("ResponseJSON() error = %v", err)
			}

			if w.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", w.Code, tt.wantStatus)
			}

			if got := w.Body.String(); got != tt.wantBody {
				t.Errorf("body = %q, want %q", got, tt.wantBody)
			}

			if got := w.Header().Get("Content-Type"); got != tt.wantContentType {
				t.Errorf("Content-Type = %q, want %q", got, tt.wantContentType)
			}

			if tt.wantCustomHeader != "" {
				if got := w.Header().Get("X-Custom-Header"); got != tt.wantCustomHeader {
					t.Errorf("X-Custom-Header = %q, want %q", got, tt.wantCustomHeader)
				}
			}
		})
	}
}

func TestResponseJSON_MarshalError(t *testing.T) {
	w := httptest.NewRecorder()

	err := ResponseJSON(w, func() {}, http.StatusInternalServerError, nil)
	if err == nil {
		t.Fatal("ResponseJSON() error = nil, want error")
	}

	if w.Body.Len() != 0 {
		t.Errorf("body = %q, want empty body", w.Body.String())
	}
}
