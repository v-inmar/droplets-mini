package handlers

import (
	dtomodels "droplets_mini_historyservice/internal/models/dto_models"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

type mockHealthService struct {
	serviceModel *dtomodels.GetHealthResponse
	serviceErr   error
}

func (s *mockHealthService) GetHealth() (*dtomodels.GetHealthResponse, error) {
	return s.serviceModel, s.serviceErr
}

func TestGetHealthHandler(t *testing.T) {
	serviceModel := dtomodels.GetHealthResponse{
		Message: "test message 1",
		Status:  200,
		Time:    time.Now().UTC(),
	}

	successResponse := dtomodels.GetHealthResponse{
		Message: serviceModel.Message,
		Status:  serviceModel.Status,
		Time:    serviceModel.Time,
	}
	errorResponse := dtomodels.ErrorResponse{
		Message: "internal server error",
		Status:  http.StatusInternalServerError,
	}

	serverError := errors.New("server error")
	testCases := []struct {
		name                    string
		serviceModel            *dtomodels.GetHealthResponse
		serviceErr              error
		expectedErrResponse     *dtomodels.ErrorResponse
		expectedSuccessResponse *dtomodels.GetHealthResponse
		expectedStatus          int
	}{
		{
			name:                    "success response",
			serviceModel:            &serviceModel,
			serviceErr:              nil,
			expectedErrResponse:     nil,
			expectedSuccessResponse: &successResponse,
			expectedStatus:          http.StatusOK,
		},
		{
			name:                    "error response",
			serviceModel:            nil,
			serviceErr:              serverError,
			expectedErrResponse:     &errorResponse,
			expectedSuccessResponse: nil,
			expectedStatus:          http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			service := &mockHealthService{
				tc.serviceModel, tc.serviceErr,
			}

			handler := NewHealthHandler(service)

			req := httptest.NewRequest(http.MethodGet, "/health", nil)
			rec := httptest.NewRecorder()

			handler.GetHealthHandler(rec, req)

			if tc.expectedStatus != rec.Code {
				t.Errorf("expected status: %d, got: %d", tc.expectedStatus, rec.Code)
			}

			if tc.expectedSuccessResponse != nil && tc.expectedErrResponse == nil {
				var got dtomodels.GetHealthResponse
				if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
					t.Fatalf("decoding rec body failed: %v", err)
				}

				if !reflect.DeepEqual(*tc.expectedSuccessResponse, got) {
					t.Errorf("expected service model did not match recorded response body. expected: %v, got: %v", tc.expectedSuccessResponse, got)
				}
			}

			if tc.expectedSuccessResponse == nil && tc.expectedErrResponse != nil {
				var got dtomodels.ErrorResponse
				if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
					t.Fatalf("decoding rec body failed: %v", err)
				}

				if !reflect.DeepEqual(*tc.expectedErrResponse, got) {
					t.Errorf("expected service model did not match recorded response body")
				}
			}

		})
	}
}
