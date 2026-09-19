package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	dtomodels "droplets_mini_todoservice/internal/models/dto_models"
)

type fakeHealthService struct {
	resp *dtomodels.GetHealthResponse
	err  error
}

func (f *fakeHealthService) GetHealth() (*dtomodels.GetHealthResponse, error) {
	return f.resp, f.err
}

func TestHealthHandler_GetHealthHandler(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		h := NewHealthHandler(&fakeHealthService{
			resp: &dtomodels.GetHealthResponse{
				Message: "todo service up and running",
				Status:  http.StatusOK,
				Time:    time.Now().UTC(),
			},
		})

		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		rec := httptest.NewRecorder()
		h.GetHealthHandler(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
		}

		var body dtomodels.GetHealthResponse
		if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
			t.Fatalf("unmarshal: %v", err)
		}
		if body.Message != "todo service up and running" {
			t.Errorf("message = %q", body.Message)
		}
		if body.Status != http.StatusOK {
			t.Errorf("status field = %d, want %d", body.Status, http.StatusOK)
		}
	})

	t.Run("service error", func(t *testing.T) {
		h := NewHealthHandler(&fakeHealthService{err: errors.New("boom")})

		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		rec := httptest.NewRecorder()
		h.GetHealthHandler(rec, req)

		if rec.Code != http.StatusInternalServerError {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusInternalServerError)
		}
	})
}