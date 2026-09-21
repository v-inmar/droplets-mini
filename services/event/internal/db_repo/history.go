package dbrepo

import (
	"context"
	dbmodels "droplets_mini_eventservice/internal/models/db_models"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

var ErrEvenPIDExist = errors.New("event id exist")

type HistoryRepo interface {
	CreateTx(ctx context.Context, tx *sqlx.Tx, eventPID int64, eventHappened string, eventHappenedAt time.Time, taskID int64, taskValue string, taskPID int64, taskCompleted bool) (*dbmodels.HistoryItem, error)
	GetDBInstance() *sqlx.DB
}

type PostgresHistoryDBRepo struct {
	db *sqlx.DB
}

func NewPostgresHistoryDBRepo(db *sqlx.DB) *PostgresHistoryDBRepo {
	return &PostgresHistoryDBRepo{
		db: db,
	}
}

func (repo *PostgresHistoryDBRepo) GetDBInstance() *sqlx.DB {
	return repo.db
}

func (repo *PostgresHistoryDBRepo) CreateTx(ctx context.Context, tx *sqlx.Tx, eventPID int64, eventHappened string, eventHappenedAt time.Time, taskID int64, taskValue string, taskPID int64, taskCompleted bool) (*dbmodels.HistoryItem, error) {
	query := `
	INSERT INTO history_model(event_pid, event_happened, event_happened_at, task_id,task_value, task_pid, task_completed)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id, event_pid, event_happened, event_happened_at, task_id, task_value, task_pid, task_completed, created_at
	`

	var model dbmodels.HistoryItem
	if err := tx.GetContext(ctx, &model, query, eventPID, eventHappened, eventHappenedAt, taskID, taskValue, taskPID, taskCompleted); err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" && pqErr.Constraint == "history_model_eventpid_unique" {
			return nil, ErrEvenPIDExist
		}
		return nil, err
	}

	return &model, nil
}
