package dbrepo

import (
	"context"
	"database/sql"
	dbmodels "droplets_mini_todoservice/internal/models/db_models"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

func TestPostgresTaskDBRepo_CreateTx(t *testing.T) {
	expectedError := errors.New("db error")
	expectedOutput := dbmodels.TaskModel{
		ID:        1,
		PID:       12345,
		Value:     "task1",
		Completed: false,
		CreatedAt: time.Now().UTC(),
	}
	testCases := []struct {
		name           string
		setupMock      func(mock sqlmock.Sqlmock)
		expectedOutput *dbmodels.TaskModel
		expectedError  error
	}{
		{
			name: "successful creation",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()

				mock.ExpectQuery(regexp.QuoteMeta(`
				INSERT INTO task_model (value, pid)
				VALUES ($1, $2)
				RETURNING id, pid, value, completed, created_at
				`)).WithArgs(expectedOutput.Value, expectedOutput.PID).WillReturnRows(
					sqlmock.NewRows([]string{"id", "pid", "value", "completed", "created_at"}).AddRow(1, expectedOutput.PID, expectedOutput.Value, false, expectedOutput.CreatedAt),
				)
			},
			expectedError:  nil,
			expectedOutput: &expectedOutput,
		},
		{
			name: "db error",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()

				mock.ExpectQuery(regexp.QuoteMeta(`
				INSERT INTO task_model (value, pid)
				VALUES ($1, $2)
				RETURNING id, pid, value, completed, created_at
				`)).WithArgs(expectedOutput.Value, expectedOutput.PID).WillReturnError(expectedError)
			},
			expectedError:  expectedError,
			expectedOutput: &expectedOutput,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("unable to create sqlmock: %v", err)
			}
			defer db.Close()

			sqlxDB := sqlx.NewDb(db, "postgres-mock")
			tc.setupMock(mock)

			tx, err := sqlxDB.BeginTxx(context.Background(), nil)
			if err != nil {
				t.Fatalf("unable to begin transaction: %v", err)
			}
			defer tx.Rollback()

			repo := NewPostgresTaskDBRepo(sqlxDB)
			task, err := repo.CreateTx(context.Background(), tx, tc.expectedOutput.Value, tc.expectedOutput.PID)

			if tc.expectedError == nil {
				if err != nil {
					t.Fatalf("expected no error but got one: %v", err)
				}

				if !reflect.DeepEqual(*tc.expectedOutput, *task) {
					t.Errorf("expected output and task output did not match. expected: %v, got: %v", tc.expectedOutput, task)
				}
			} else {
				if tc.expectedError != err {
					t.Fatalf("errors did not match. expected: %v, got: %v", tc.expectedError, err)
				}

				if task != nil {
					t.Fatalf("resulting task is not nil")
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("mock expectations were not met: %v", err)
			}

		})
	}
}

func TestPostgresTaskDBRepo_ReadByPID(t *testing.T) {
	var testPID int64 = 12345
	testTask := dbmodels.TaskModel{
		ID:        1,
		PID:       testPID,
		Value:     "task1",
		Completed: false,
		CreatedAt: time.Now().UTC(),
	}

	testCases := []struct {
		name          string
		setupMock     func(mock sqlmock.Sqlmock)
		inputPID      int64
		expectedData  *dbmodels.TaskModel
		expectedError error
	}{
		{
			name:     "task with pid exist and data found",
			inputPID: testPID,
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "pid", "value", "completed", "created_at",
				}).AddRow(testTask.ID, testTask.PID, testTask.Value, testTask.Completed, testTask.CreatedAt)

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM task_model WHERE pid=$1`)).WithArgs(testPID).WillReturnRows(rows)
			},
			expectedData:  &testTask,
			expectedError: nil,
		},
		{
			name:     "not found",
			inputPID: 123123,
			setupMock: func(mock sqlmock.Sqlmock) {

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM task_model WHERE pid=$1`)).WithArgs(123123).WillReturnError(sql.ErrNoRows)
			},
			expectedData:  nil,
			expectedError: sql.ErrNoRows,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("unable to create sqlmock: %v", err)
			}
			defer db.Close()
			tc.setupMock(mock)

			sqlxDB := sqlx.NewDb(db, "postgres-mock")
			repo := NewPostgresTaskDBRepo(sqlxDB)

			task, err := repo.ReadByPID(context.Background(), tc.inputPID)
			if tc.expectedError == nil {

				if err != nil {
					t.Fatalf("expected no error but got one: %v", err)
				}

				if !reflect.DeepEqual(*task, *tc.expectedData) {
					t.Errorf("expected data: %v did not match returned data: %v", tc.expectedData, task)
				}

			} else {
				fmt.Println(task)
				if tc.expectedError != err {
					t.Errorf("expected error: %v, got: %v", tc.expectedError, err)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("mock expectations were not met: %v", err)
			}

		})
	}

}

func TestPostgresTaskDBRepo_ReadAll(t *testing.T) {
	errDB := errors.New("database error")
	testTask := dbmodels.TaskModel{
		ID:        1,
		PID:       123456789,
		Value:     "task1",
		Completed: false,
		CreatedAt: time.Now().UTC(),
	}

	testCases := []struct {
		name          string
		data          any
		setupMock     func(mock sqlmock.Sqlmock)
		expectedCount int
		expectedErr   error
	}{
		{
			name: "return non-empty list",
			data: testTask,
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "pid", "value", "completed", "created_at",
				}).AddRow(testTask.ID, testTask.PID, testTask.Value, testTask.Completed, testTask.CreatedAt)

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM task_model`)).WillReturnRows(rows)
			},
			expectedCount: 1,
			expectedErr:   nil,
		},
		{
			name: "return empty list",
			data: nil,
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "pid", "value", "completed", "created_at",
				})

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM task_model`)).WillReturnRows(rows)
			},
			expectedCount: 0,
			expectedErr:   nil,
		},
		{
			name: "return db error",
			data: nil,
			setupMock: func(mock sqlmock.Sqlmock) {

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM task_model`)).WillReturnError(errDB)
			},
			expectedCount: 0,
			expectedErr:   errDB,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("unable to create sqlmock: %v", err)
			}
			defer db.Close()
			tc.setupMock(mock)

			sqlxDB := sqlx.NewDb(db, "postgres-mock")
			repo := NewPostgresTaskDBRepo(sqlxDB)

			tasks, err := repo.ReadAll(context.Background())

			if tc.expectedErr == nil {
				if err != nil {
					t.Fatalf("no error expected but got: %v", err)
				}

				if len(tasks) != tc.expectedCount {
					t.Errorf("tasks length expected: %d, got: %d", tc.expectedCount, len(tasks))
				}

				if tc.expectedCount > 0 {
					if tasks == nil {
						t.Errorf("tasks is expected to be not nil but got nil tasks")
					}

					if !reflect.DeepEqual(tc.data, tasks[0]) {
						t.Error("test data and returned tasks[0] did not match")
					}
				}

			} else {
				if tc.expectedErr != err {
					t.Errorf("expected error: %v, got: %v", tc.expectedErr, err)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("mock expectations were not met: %v", err)
			}
		})
	}
}

func TestPostgresTaskDBRepo_UpdateTx(t *testing.T) {

	expectedError := errors.New("db error")

	expectedOutput := dbmodels.TaskModel{
		ID:        1,
		PID:       123456789,
		Value:     "task1_updated",
		Completed: true,
		CreatedAt: time.Now().UTC(),
	}

	testCases := []struct {
		name           string
		setupMock      func(mock sqlmock.Sqlmock)
		expectedOutput *dbmodels.TaskModel
		expectedError  error
	}{
		{
			name: "succesful update",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()

				mock.ExpectQuery(regexp.QuoteMeta(`
				UPDATE task_model
				SET value=$1, completed=$2
				WHERE pid=$3
				RETURNING id, pid, value, completed, created_at
				`)).WithArgs(expectedOutput.Value, expectedOutput.Completed, expectedOutput.PID).WillReturnRows(
					sqlmock.NewRows([]string{"id", "pid", "value", "completed", "created_at"}).AddRow(expectedOutput.ID, expectedOutput.PID, expectedOutput.Value, expectedOutput.Completed, expectedOutput.CreatedAt),
				)
			},
			expectedOutput: &expectedOutput,
			expectedError:  nil,
		},
		{
			name: "db error",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()

				mock.ExpectQuery(regexp.QuoteMeta(`
				UPDATE task_model
				SET value=$1, completed=$2
				WHERE pid=$3
				RETURNING id, pid, value, completed, created_at
				`)).WithArgs(expectedOutput.Value, expectedOutput.Completed, expectedOutput.PID).WillReturnError(expectedError)
			},
			expectedOutput: &expectedOutput,
			expectedError:  expectedError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("unable to create sqlmock: %v", err)
			}
			defer db.Close()

			sqlxDB := sqlx.NewDb(db, "postgres-mock")
			tc.setupMock(mock)

			tx, err := sqlxDB.BeginTxx(context.Background(), nil)
			if err != nil {
				t.Fatalf("unable to begin transaction: %v", err)
			}
			defer tx.Rollback()

			repo := NewPostgresTaskDBRepo(sqlxDB)

			updatedTask, err := repo.UpdateTx(context.Background(), tx, tc.expectedOutput.Value, tc.expectedOutput.PID, tc.expectedOutput.Completed)

			if tc.expectedError == nil {
				if err != nil {
					t.Fatalf("no error is expected but got one: %v", err)
				}

				if !reflect.DeepEqual(*tc.expectedOutput, *updatedTask) {
					t.Errorf("expected output and updated task output did not match. expected: %v, got: %v", tc.expectedOutput, updatedTask)
				}
			} else {
				if tc.expectedError != err {
					t.Fatalf("expected error and returned error did not match. expected: %v, got:%v", tc.expectedError, err)
				}
				if updatedTask != nil {
					t.Fatalf("resulting updated task is not nil")
				}

			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("mock expectations were not met: %v", err)
			}
		})
	}
}

/*
	func (r *PostgresTaskDBRepo) DeleteTx(ctx context.Context, tx *sqlx.Tx, pid string) error {
		query := `
		DELETE FROM task_model
		WHERE pid=$1
		`

		if _, err := tx.ExecContext(ctx, query, pid); err != nil {
			return err
		}
		return nil
	}
*/
func TestPostgresTaskDBRepo_DeleteTx(t *testing.T) {
	expectedError := errors.New("db error")
	var testPID int64 = 123456789
	testCases := []struct {
		name          string
		pid           int64
		setupMock     func(mock sqlmock.Sqlmock)
		expectedError error
	}{
		{
			name:          "successfully deleted",
			pid:           testPID,
			expectedError: nil,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()

				mock.ExpectExec(
					regexp.QuoteMeta(`
					DELETE FROM task_model
					WHERE pid=$1
					`),
				).WithArgs(testPID).WillReturnResult(sqlmock.NewResult(0, 1))
			},
		},
		{
			name:          "successfully deleted",
			pid:           testPID,
			expectedError: expectedError,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()

				mock.ExpectExec(
					regexp.QuoteMeta(`
					DELETE FROM task_model
					WHERE pid=$1
					`),
				).WithArgs(testPID).WillReturnError(expectedError)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("unable to create sqlmock: %v", err)
			}
			defer db.Close()

			sqlxDB := sqlx.NewDb(db, "postgres-mock")

			tc.setupMock(mock)

			tx, err := sqlxDB.BeginTxx(context.Background(), nil)
			if err != nil {
				t.Fatalf("unable to begin transaction: %v", err)
			}
			defer tx.Rollback()

			repo := NewPostgresTaskDBRepo(sqlxDB)

			err = repo.DeleteTx(context.Background(), tx, tc.pid)

			if tc.expectedError == nil {
				if err != nil {
					t.Fatalf("no error is expected but got one: %v", err)
				}
			} else {
				if tc.expectedError != err {
					t.Fatalf("expected error and returned error did not match. expected: %v, got:%v", tc.expectedError, err)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("mock expectations were not met: %v", err)
			}
		})
	}
}
