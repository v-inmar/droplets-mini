package services

import (
	"context"
	dbmodels "droplets_mini_historyservice/internal/models/db_models"
	dtomodels "droplets_mini_historyservice/internal/models/dto_models"
	"errors"
	"log"
	"testing"
	"time"
)

type mockHistoryRepo struct {
	histories []dbmodels.HistoryModel
	err       error
}

func (m *mockHistoryRepo) ReadAll(ctx context.Context) ([]dbmodels.HistoryModel, error) {
	return m.histories, m.err
}

func TestHistoryService_GetAllHistory(t *testing.T) {
	testHistory := dbmodels.HistoryModel{
		ID:              1,
		EventPID:        123456789,
		EventHappened:   "CREATED",
		EventHappenedAt: time.Now().UTC(),
		TaskValue:       "task1",
		TaskPID:         987654321,
		CreatedAt:       time.Now().UTC(),
	}

	repoErr := errors.New("repo error")

	testCases := []struct {
		name          string
		repoHistories []dbmodels.HistoryModel
		repoError     error
		expectedResp  *dtomodels.GetAllHistoryResponse
		expectedErr   error
	}{
		{
			name: "returns success with 1 history",
			repoHistories: []dbmodels.HistoryModel{
				testHistory,
			},
			repoError: nil,
			expectedResp: &dtomodels.GetAllHistoryResponse{
				History: []dtomodels.GetHistoryResponse{
					{
						EventPID:        testHistory.EventPID,
						EventHappened:   testHistory.EventHappened,
						EventHappenedAt: testHistory.EventHappenedAt,
						TaskValue:       testHistory.TaskValue,
						TaskPID:         testHistory.TaskPID,
					},
				},
			},
			expectedErr: nil,
		},

		{
			name:          "returns success with empty history list",
			repoHistories: []dbmodels.HistoryModel{},
			repoError:     nil,
			expectedResp: &dtomodels.GetAllHistoryResponse{
				History: []dtomodels.GetHistoryResponse{},
			},
			expectedErr: nil,
		},

		{
			name:          "returns with error",
			repoHistories: nil,
			repoError:     repoErr,
			expectedResp:  nil,
			expectedErr:   repoErr,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := &mockHistoryRepo{
				histories: tc.repoHistories,
				err:       tc.repoError,
			}
			service := NewHistoryService(repo)
			got, err := service.GetAllHistory(context.Background())

			if tc.expectedErr == nil {
				if err != nil {
					t.Fatalf("expected not error but got: %v", err)
				}

				if tc.expectedResp != nil {
					if len(tc.expectedResp.History) != len(got.History) {
						t.Fatalf("expected response history lenght: %d, got: %d", len(tc.expectedResp.History), len(got.History))
					}
				} else {
					if len(got.History) > 0 {
						t.Errorf("expected to not have history but got history length: %d", len(got.History))
					}
				}
				return
			} else {
				if err == nil {
					log.Fatalf("expected error but got: %v", err)
				}

				if tc.expectedErr != err {
					log.Fatalf("errors did not match. expected: %v, got: %v", tc.expectedErr, err)
				}
			}

		})
	}
}
