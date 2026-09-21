package handlers

import (
	"context"
	dtomodels "droplets_mini_historyservice/internal/models/dto_models"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

type mockHistoryService struct {
	serviceModel *dtomodels.GetAllHistoryResponse
	serviceErr   error
}

func (s *mockHistoryService) GetAllHistory(ctx context.Context) (*dtomodels.GetAllHistoryResponse, error) {
	return s.serviceModel, s.serviceErr
}

func TestGetAllHistoryHandler(t *testing.T) {
	serviceResponse := dtomodels.GetHistoryResponse{
		EventPID:        123456789,
		EventHappened:   "CREATED",
		EventHappenedAt: time.Now().UTC(),
		TaskValue:       "task1",
		TaskPID:         987654321,
	}

	serviceErr := errors.New("service error")
	serviceModel := dtomodels.GetAllHistoryResponse{
		History: []dtomodels.GetHistoryResponse{serviceResponse},
	}

	expectedSuccessResponse := serviceModel
	expectedErrorResponse := dtomodels.ErrorResponse{
		Message: "internal server error",
		Status:  http.StatusInternalServerError,
	}

	testCases := []struct {
		name                    string
		serviceModel            *dtomodels.GetAllHistoryResponse
		serviceError            error
		expectedSuccessResponse *dtomodels.GetAllHistoryResponse
		expectedErrorResponse   *dtomodels.ErrorResponse
		expectedErr             error
		expectedStatus          int
	}{
		{
			name:                    "success",
			serviceModel:            &serviceModel,
			serviceError:            nil,
			expectedSuccessResponse: &expectedSuccessResponse,
			expectedErrorResponse:   nil,
			expectedStatus:          http.StatusOK,
		},
		{
			name:                    "error",
			serviceModel:            nil,
			serviceError:            serviceErr,
			expectedSuccessResponse: nil,
			expectedErrorResponse:   &expectedErrorResponse,
			expectedStatus:          http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			service := &mockHistoryService{
				serviceModel: tc.serviceModel,
				serviceErr:   tc.serviceError,
			}

			handler := NewHistoryHandler(service)
			req := httptest.NewRequest(http.MethodGet, "/items", nil)
			rec := httptest.NewRecorder()

			handler.GetAllHistoryHandler(rec, req)

			if tc.expectedStatus != rec.Code {
				t.Fatalf("expected status: %d, got: %d", tc.expectedStatus, rec.Code)
			}

			if tc.expectedSuccessResponse != nil && tc.expectedErrorResponse == nil {
				var got dtomodels.GetAllHistoryResponse

				if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
					t.Fatalf("decoding rec body failed: %v", err)
				}

				if !reflect.DeepEqual(*tc.expectedSuccessResponse, got) {
					t.Errorf("expected service model did not match recorded response body. expected: %v, got: %v", tc.expectedSuccessResponse, got)
				}

			}

			if tc.expectedSuccessResponse == nil && tc.expectedErrorResponse != nil {
				var got dtomodels.ErrorResponse

				if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
					t.Fatalf("decoding rec body failed: %v", err)
				}

				if !reflect.DeepEqual(*tc.expectedErrorResponse, got) {
					t.Errorf("expected service model did not match recorded response body. expected: %v, got: %v", tc.expectedErrorResponse, got)
				}
			}
		})
	}
}
