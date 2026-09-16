package repo

import (
	"context"
	"droplets_mini/services/todo/internal/models/db"

	"github.com/jmoiron/sqlx"
)

type TaskRepoTx struct {
	tx *sqlx.Tx
}

func NewTaskRepoTx(tx *sqlx.Tx) *TaskRepoTx {
	return &TaskRepoTx{
		tx: tx,
	}
}

// Create new task
func (repo *TaskRepoTx) Create(ctx context.Context, value, pid string) (*db.TaskModel, error) {
	query := `
	INSERT INTO task_model (value, pid)
	VALUES ($1, $2)
	RETURNING id, pid, value, completed, created_at
	`

	var task db.TaskModel
	if err := repo.tx.GetContext(ctx, &task, query, value, pid); err != nil {
		return nil, err
	}

	return &task, nil
}

// Update task
func (repo *TaskRepoTx) Update(ctx context.Context, value, pid string, completed bool) (*db.TaskModel, error) {
	query := `
	UPDATE task_model
	SET value=$1, completed=$2
	WHERE pid=$3
	RETURNING id, pid, value, completed, created_at
	`

	var task db.TaskModel

	if err := repo.tx.GetContext(ctx, &task, query, value, completed, pid); err != nil {
		return nil, err
	}
	return &task, nil

}

func (repo *TaskRepoTx) Delete(ctx context.Context, pid string) error {
	query := `
	DELETE FROM task_model
	WHERE pid=$1
	`

	if _, err := repo.tx.ExecContext(ctx, query, pid); err != nil {
		return err
	}
	return nil
}

type TaskRepoDB struct {
	db *sqlx.DB
}

func NewTaskRepoDB(db *sqlx.DB) *TaskRepoDB {
	return &TaskRepoDB{
		db: db,
	}
}

// Read all tasks
func (repo *TaskRepoDB) ReadAll(ctx context.Context) ([]db.TaskModel, error) {
	query := `
	SELECT * FROM task_model
	`

	var tasks []db.TaskModel

	if err := repo.db.SelectContext(ctx, &tasks, query); err != nil {
		return nil, err
	}

	return tasks, nil
}

// Read task by pid
func (repo *TaskRepoDB) ReadByPID(ctx context.Context, pid string) (*db.TaskModel, error) {
	query := `
	SELECT * FROM task_model
	WHERE pid=$1
	`

	var task db.TaskModel
	if err := repo.db.GetContext(ctx, &task, query, pid); err != nil {
		return nil, err
	}
	return &task, nil
}
