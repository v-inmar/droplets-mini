package dbrepo

import (
	"context"
	dbmodels "droplets_mini_todoservice/internal/models/db_models"

	"github.com/jmoiron/sqlx"
)

type TaskRepo interface {
	CreateTx(ctx context.Context, tx *sqlx.Tx, value string, pid int64) (*dbmodels.TaskModel, error)
	ReadAll(ctx context.Context) ([]dbmodels.TaskModel, error)
	ReadByPID(ctx context.Context, pid int64) (*dbmodels.TaskModel, error)
	UpdateTx(ctx context.Context, tx *sqlx.Tx, value string, pid int64, completed bool) (*dbmodels.TaskModel, error)
	DeleteTx(ctx context.Context, tx *sqlx.Tx, pid int64) error
	GetDBInstance() *sqlx.DB
}

type PostgresTaskDBRepo struct {
	db *sqlx.DB
}

func NewPostgresTaskDBRepo(db *sqlx.DB) *PostgresTaskDBRepo {
	return &PostgresTaskDBRepo{
		db: db,
	}
}

func (r *PostgresTaskDBRepo) GetDBInstance() *sqlx.DB {
	return r.db
}

func (r *PostgresTaskDBRepo) CreateTx(ctx context.Context, tx *sqlx.Tx, value string, pid int64) (*dbmodels.TaskModel, error) {
	query := `
	INSERT INTO task_model (value, pid)
	VALUES ($1, $2)
	RETURNING id, pid, value, completed, created_at
	`

	var task dbmodels.TaskModel
	if err := tx.GetContext(ctx, &task, query, value, pid); err != nil {
		return nil, err
	}

	return &task, nil
}

func (r *PostgresTaskDBRepo) ReadAll(ctx context.Context) ([]dbmodels.TaskModel, error) {
	query := `
	SELECT * FROM task_model
	`

	var tasks []dbmodels.TaskModel

	if err := r.db.SelectContext(ctx, &tasks, query); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *PostgresTaskDBRepo) ReadByPID(ctx context.Context, pid int64) (*dbmodels.TaskModel, error) {
	query := `
	SELECT * FROM task_model
	WHERE pid=$1
	`

	var task dbmodels.TaskModel
	if err := r.db.GetContext(ctx, &task, query, pid); err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *PostgresTaskDBRepo) UpdateTx(ctx context.Context, tx *sqlx.Tx, value string, pid int64, completed bool) (*dbmodels.TaskModel, error) {
	query := `
	UPDATE task_model
	SET value=$1, completed=$2
	WHERE pid=$3
	RETURNING id, pid, value, completed, created_at
	`

	var task dbmodels.TaskModel

	if err := tx.GetContext(ctx, &task, query, value, completed, pid); err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *PostgresTaskDBRepo) DeleteTx(ctx context.Context, tx *sqlx.Tx, pid int64) error {
	query := `
	DELETE FROM task_model
	WHERE pid=$1
	`

	if _, err := tx.ExecContext(ctx, query, pid); err != nil {
		return err
	}
	return nil
}
