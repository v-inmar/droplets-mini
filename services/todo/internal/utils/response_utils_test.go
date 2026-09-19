package utils

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestResponseJSON(t *testing.T) {
	t.Run("writes marshaled body with status and headers", func(t *testing.T) {
		rec := httptest.NewRecorder()
		headers := http.Header{}
		headers.Set("X-Custom", "abc")

		data := map[string]string{"key": "value"}

		if err := ResponseJSON(rec, data, http.StatusCreated, headers); err != nil {
			t.Fatalf("ResponseJSON error: %v", err)
		}

		if rec.Code != http.StatusCreated {
			t.Errorf("status = %d, want %d", rec.Code, http.StatusCreated)
		}
		if got := rec.Header().Get("Content-Type"); got != "application/json" {
			t.Errorf("content-type = %q, want application/json", got)
		}
		if got := rec.Header().Get("X-Custom"); got != "abc" {
			t.Errorf("X-Custom = %q, want abc", got)
		}

		var got map[string]string
		if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
			t.Fatalf("unmarshal response: %v", err)
		}
		if got["key"] != "value" {
			t.Errorf("body = %v, want key=value", got)
		}
		if len(rec.Body.Bytes()) == 0 || rec.Body.Bytes()[len(rec.Body.Bytes())-1] != '\n' {
			t.Error("expected trailing newline in body")
		}
	})

	t.Run("returns error for non-marshalable data", func(t *testing.T) {
		rec := httptest.NewRecorder()
		if err := ResponseJSON(rec, make(chan int), http.StatusOK, nil); err == nil {
			t.Error("expected error, got nil")
		}
	})
}