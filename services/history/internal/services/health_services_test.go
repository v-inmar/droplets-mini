package services

import (
	dtomodels "droplets_mini_historyservice/internal/models/dto_models"
	"testing"
	"time"
)

func TestHealthService_GetHealth(t *testing.T) {
	respModel := dtomodels.GetHealthResponse{
		Message: "history service up and running",
		Status:  200,
		Time:    time.Now().UTC(), // service responds with its own time, so just check that it is populated or after
	}

	testCases := []struct {
		name          string
		expectedData  *dtomodels.GetHealthResponse
		expectedError error
	}{
		{
			name:          "returns non-empty model",
			expectedData:  &respModel,
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			service := NewHealthService()
			got, err := service.GetHealth()

			if tc.expectedError != nil {
				if err == nil {
					t.Fatalf("expected return error: %v got: %v", tc.expectedError, err)
				}
			}

			if got.Message != tc.expectedData.Message {
				t.Errorf("expected message: %s, got: %s", tc.expectedData.Message, got.Message)
			}

			if got.Status != tc.expectedData.Status {
				t.Errorf("expected status: %d, got: %d", tc.expectedData.Status, got.Status)
			}

			if !got.Time.After(tc.expectedData.Time) {
				t.Errorf("expected time is not before service was executed: %v, got: %v", tc.expectedData.Time, got.Time)
			}
		})
	}
}
