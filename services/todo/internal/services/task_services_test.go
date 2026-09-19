package services

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	dbrepo "droplets_mini_todoservice/internal/db_repo"
	dtomodels "droplets_mini_todoservice/internal/models/dto_models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

func newMockedTaskService(t *testing.T) (*TaskService, sqlmock.Sqlmock) {
	t.Helper()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	sqlxDB := sqlx.NewDb(db, "postgres")
	return NewTaskService(dbrepo.NewPostgresTaskDBRepo(sqlxDB)), mock
}

func taskRows(pid, value string, completed bool, createdAt time.Time) *sqlmock.Rows {
	return sqlmock.NewRows([]string{"id", "pid", "value", "completed", "created_at"}).
		AddRow(1, pid, value, completed, createdAt)
}

func TestTaskService_CreateNewTaskItem(t *testing.T) {
	ctx := context.Background()
	created := time.Now().UTC()

	t.Run("success", func(t *testing.T) {
		s, mock := newMockedTaskService(t)

		mock.ExpectBegin()
		mock.ExpectQuery("INSERT INTO task_model").
			WithArgs("buy milk", sqlmock.AnyArg()).
			WillReturnRows(taskRows("p-1", "buy milk", false, created))
		mock.ExpectCommit()

		resp, err := s.CreateNewTaskItem(ctx, &dtomodels.CreateRequest{Value: "buy milk"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.PID != "p-1" || resp.Value != "buy milk" || resp.Completed {
			t.Errorf("unexpected response: %+v", resp)
		}
		if !resp.CreatedAt.Equal(created) {
			t.Errorf("created_at = %v, want %v", resp.CreatedAt, created)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	})

	t.Run("begin error", func(t *testing.T) {
		s, mock := newMockedTaskService(t)

		mock.ExpectBegin().WillReturnError(errors.New("begin failed"))

		resp, err := s.CreateNewTaskItem(ctx, &dtomodels.CreateRequest{Value: "buy milk"})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if resp != nil {
			t.Errorf("expected nil response, got %+v", resp)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	})

	t.Run("insert error rolls back", func(t *testing.T) {
		s, mock := newMockedTaskService(t)

		mock.ExpectBegin()
		mock.ExpectQuery("INSERT INTO task_model").
			WithArgs("buy milk", sqlmock.AnyArg()).
			WillReturnError(errors.New("insert failed"))
		mock.ExpectRollback()

		resp, err := s.CreateNewTaskItem(ctx, &dtomodels.CreateRequest{Value: "buy milk"})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if resp != nil {
			t.Errorf("expected nil response, got %+v", resp)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	})

	t.Run("commit error", func(t *testing.T) {
		s, mock := newMockedTaskService(t)

		mock.ExpectBegin()
		mock.ExpectQuery("INSERT INTO task_model").
			WithArgs("buy milk", sqlmock.AnyArg()).
			WillReturnRows(taskRows("p-1", "buy milk", false, created))
		mock.ExpectCommit().WillReturnError(errors.New("commit failed"))

		resp, err := s.CreateNewTaskItem(ctx, &dtomodels.CreateRequest{Value: "buy milk"})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if resp != nil {
			t.Errorf("expected nil response, got %+v", resp)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	})
}

func TestTaskService_ReadAllTasks(t *testing.T) {
	ctx := context.Background()
	t1 := time.Now().UTC()
	t2 := t1.Add(time.Hour)

	t.Run("success", func(t *testing.T) {
		s, mock := newMockedTaskService(t)

		rows := sqlmock.NewRows([]string{"id", "pid", "value", "completed", "created_at"}).
			AddRow(1, "p-1", "first", false, t1).
			AddRow(2, "p-2", "second", true, t2)
		mock.ExpectQuery("SELECT \\* FROM task_model").WillReturnRows(rows)

		resp, err := s.ReadAllTasks(ctx)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(resp.Tasks) != 2 {
			t.Fatalf("len(tasks) = %d, want 2", len(resp.Tasks))
		}
		if resp.Tasks[0].PID != "p-1" || resp.Tasks[0].Value != "first" {
			t.Errorf("unexpected first task: %+v", resp.Tasks[0])
		}
		if resp.Tasks[1].PID != "p-2" || resp.Tasks[1].Value != "second" || !resp.Tasks[1].Completed {
			t.Errorf("unexpected second task: %+v", resp.Tasks[1])
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	})

	t.Run("empty result", func(t *testing.T) {
		s, mock := newMockedTaskService(t)

		mock.ExpectQuery("SELECT \\* FROM task_model").
			WillReturnRows(sqlmock.NewRows([]string{"id", "pid", "value", "completed", "created_at"}))

		resp, err := s.ReadAllTasks(ctx)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(resp.Tasks) != 0 {
			t.Errorf("len(tasks) = %d, want 0", len(resp.Tasks))
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	})

	t.Run("no rows is not an error", func(t *testing.T) {
		s, mock := newMockedTaskService(t)

		mock.ExpectQuery("SELECT \\* FROM task_model").WillReturnError(sql.ErrNoRows)

		resp, err := s.ReadAllTasks(ctx)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(resp.Tasks) != 0 {
			t.Errorf("len(tasks) = %d, want 0", len(resp.Tasks))
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	})

	t.Run("other error propagates", func(t *testing.T) {
		s, mock := newMockedTaskService(t)

		mock.ExpectQuery("SELECT \\* FROM task_model").WillReturnError(errors.New("db down"))

		resp, err := s.ReadAllTasks(ctx)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if resp != nil {
			t.Errorf("expected nil response, got %+v", resp)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	})
}

func TestTaskService_ReadTaskByPID(t *testing.T) {
	ctx := context.Background()
	created := time.Now().UTC()

	t.Run("success", func(t *testing.T) {
		s, mock := newMockedTaskService(t)

		mock.ExpectQuery("SELECT \\* FROM task_model").
			WithArgs("p-1").
			WillReturnRows(taskRows("p-1", "task value", true, created))

		resp, err := s.ReadTaskByPID(ctx, "p-1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.PID != "p-1" || resp.Value != "task value" || !resp.Completed {
			t.Errorf("unexpected response: %+v", resp)
		}
		if !resp.CreatedAt.Equal(created) {
			t.Errorf("created_at = %v, want %v", resp.CreatedAt, created)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	})

	t.Run("error propagates", func(t *testing.T) {
		s, mock := newMockedTaskService(t)

		mock.ExpectQuery("SELECT \\* FROM task_model").
			WithArgs("p-1").
			WillReturnError(sql.ErrNoRows)

		resp, err := s.ReadTaskByPID(ctx, "p-1")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if resp != nil {
			t.Errorf("expected nil response, got %+v", resp)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	})
}

func TestTaskService_UpdateTask(t *testing.T) {
	ctx := context.Background()
	created := time.Now().UTC()

	t.Run("success", func(t *testing.T) {
		s, mock := newMockedTaskService(t)

		mock.ExpectBegin()
		mock.ExpectQuery("UPDATE task_model").
			WithArgs("updated value", true, "p-1").
			WillReturnRows(taskRows("p-1", "updated value", true, created))
		mock.ExpectCommit()

		resp, err := s.UpdateTask(ctx, &dtomodels.UpdateRequest{Value: "updated value", Completed: true}, "p-1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.PID != "p-1" || resp.Value != "updated value" || !resp.Completed {
			t.Errorf("unexpected response: %+v", resp)
		}
		if !resp.CreatedAt.Equal(created) {
			t.Errorf("created_at = %v, want %v", resp.CreatedAt, created)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	})

	t.Run("begin error", func(t *testing.T) {
		s, mock := newMockedTaskService(t)

		mock.ExpectBegin().WillReturnError(errors.New("begin failed"))

		resp, err := s.UpdateTask(ctx, &dtomodels.UpdateRequest{Value: "v", Completed: false}, "p-1")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if resp != nil {
			t.Errorf("expected nil response, got %+v", resp)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	})

	t.Run("update error rolls back", func(t *testing.T) {
		s, mock := newMockedTaskService(t)

		mock.ExpectBegin()
		mock.ExpectQuery("UPDATE task_model").
			WithArgs("v", false, "p-1").
			WillReturnError(errors.New("update failed"))
		mock.ExpectRollback()

		resp, err := s.UpdateTask(ctx, &dtomodels.UpdateRequest{Value: "v", Completed: false}, "p-1")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if resp != nil {
			t.Errorf("expected nil response, got %+v", resp)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	})

	t.Run("commit error", func(t *testing.T) {
		s, mock := newMockedTaskService(t)

		mock.ExpectBegin()
		mock.ExpectQuery("UPDATE task_model").
			WithArgs("v", false, "p-1").
			WillReturnRows(taskRows("p-1", "v", false, created))
		mock.ExpectCommit().WillReturnError(errors.New("commit failed"))

		resp, err := s.UpdateTask(ctx, &dtomodels.UpdateRequest{Value: "v", Completed: false}, "p-1")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if resp != nil {
			t.Errorf("expected nil response, got %+v", resp)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	})
}

func TestTaskService_DeleteTask(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		s, mock := newMockedTaskService(t)

		mock.ExpectBegin()
		mock.ExpectExec("DELETE FROM task_model").
			WithArgs("p-1").
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		if err := s.DeleteTask(ctx, "p-1"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	})

	t.Run("begin error", func(t *testing.T) {
		s, mock := newMockedTaskService(t)

		mock.ExpectBegin().WillReturnError(errors.New("begin failed"))

		if err := s.DeleteTask(ctx, "p-1"); err == nil {
			t.Fatal("expected error, got nil")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	})

	t.Run("delete error propagates", func(t *testing.T) {
		s, mock := newMockedTaskService(t)

		mock.ExpectBegin()
		mock.ExpectExec("DELETE FROM task_model").
			WithArgs("p-1").
			WillReturnError(errors.New("delete failed"))
		mock.ExpectRollback()

		if err := s.DeleteTask(ctx, "p-1"); err == nil {
			t.Fatal("expected error, got nil")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	})

	t.Run("commit error", func(t *testing.T) {
		s, mock := newMockedTaskService(t)

		mock.ExpectBegin()
		mock.ExpectExec("DELETE FROM task_model").
			WithArgs("p-1").
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit().WillReturnError(errors.New("commit failed"))

		if err := s.DeleteTask(ctx, "p-1"); err == nil {
			t.Fatal("expected error, got nil")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	})
}