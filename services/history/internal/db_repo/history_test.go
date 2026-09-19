package dbrepo

import (
	"context"
	dbmodels "droplets_mini_historyservice/internal/models/db_models"
	"errors"
	"reflect"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

func TestPostgresHistoryDBRepo(t *testing.T) {
	testHistory := dbmodels.HistoryModel{
		ID:              1,
		EventPID:        "epid123",
		EventHappened:   "CREATED",
		EventHappenedAt: time.Now().UTC(),
		TaskValue:       "item1",
		TaskPID:         "tpid123",
		CreatedAt:       time.Now().UTC(),
	}

	testCases := []struct {
		name          string
		data          *dbmodels.HistoryModel
		setupMock     func(mock sqlmock.Sqlmock)
		expectedCount int
		expectedErr   error
	}{
		{
			name: "returns non-empty histories",
			data: &testHistory,
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id",
					"event_pid",
					"event_happened",
					"event_happened_at",
					"task_value",
					"task_pid",
					"created_at",
				}).AddRow(
					testHistory.ID,
					testHistory.EventPID,
					testHistory.EventHappened,
					testHistory.EventHappenedAt,
					testHistory.TaskValue,
					testHistory.TaskPID,
					testHistory.CreatedAt,
				)

				mock.ExpectQuery(
					regexp.QuoteMeta(`SELECT * FROM history_model`),
				).WillReturnRows(rows)
			},
			expectedCount: 1,
			expectedErr:   nil,
		},

		{
			name: "returns empty histories",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id",
					"event_pid",
					"event_happened",
					"event_happened_at",
					"task_value",
					"task_pid",
					"created_at",
				})

				mock.ExpectQuery(
					regexp.QuoteMeta(`SELECT * FROM history_model`),
				).WillReturnRows(rows)
			},
			expectedCount: 0,
			expectedErr:   nil,
			data:          nil,
		},

		{
			name: "returns database error",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(
					regexp.QuoteMeta(`SELECT * FROM history_model`),
				).WillReturnError(errors.New("database error"))
			},
			expectedCount: 0,
			expectedErr:   errors.New("database error"),
			data:          nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to create sqlmock: %v", err)
			}
			defer db.Close()

			tc.setupMock(mock)

			sqlxDB := sqlx.NewDb(db, "postgres-mock")
			repo := NewPostgresHistoryDBRepo(sqlxDB)

			histories, err := repo.ReadAll(context.Background())

			if tc.expectedErr != nil {
				if err == nil {
					t.Fatalf("expected error %v, got nil", tc.expectedErr)
				}

				if err.Error() != tc.expectedErr.Error() {
					t.Fatalf("expected error %v, got %v", tc.expectedErr, err)
				}
			} else if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if len(histories) != tc.expectedCount {
				t.Fatalf(
					"expected %d history, got %d",
					tc.expectedCount,
					len(histories),
				)
			}

			if tc.data != nil {
				if len(histories) == 0 {
					t.Fatalf("expected history data, got empty result")
				}

				if !reflect.DeepEqual(histories[0], *tc.data) {
					t.Errorf(
						"expected row values: %v, got: %v",
						*tc.data,
						histories[0],
					)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("mock expectations were not met: %v", err)
			}
		})
	}
}
