package dbrepo

import (
	"context"
	dbmodels "droplets_mini_eventservice/internal/models/db_models"

	"github.com/jmoiron/sqlx"
)

type ProcessedRepo interface {
	CreateTx(ctx context.Context, tx *sqlx.Tx, event_pid int64, event_retry int, event_completed bool) (*dbmodels.EventProcessedModel, error)
	ReadByEventPID(ctx context.Context, event_pid int64) (*dbmodels.EventProcessedModel, error)
	ReadByEventPIDTx(ctx context.Context, tx *sqlx.Tx, event_pid int64) (*dbmodels.EventProcessedModel, error)
	UpdateTx(ctx context.Context, tx *sqlx.Tx, event_pid int64, event_retry int, event_processed bool) (*dbmodels.EventProcessedModel, error)
	GetDBInstance() *sqlx.DB
}

type PostgresEventProcessedDBRepo struct {
	db *sqlx.DB
}

func NewPostgresEventProcessedDBRepo(db *sqlx.DB) *PostgresEventProcessedDBRepo {
	return &PostgresEventProcessedDBRepo{
		db: db,
	}
}

func (repo *PostgresEventProcessedDBRepo) GetDBInstance() *sqlx.DB {
	return repo.db
}

func (repo *PostgresEventProcessedDBRepo) CreateTx(ctx context.Context, tx *sqlx.Tx, event_pid int64, event_retry int, event_processed bool) (*dbmodels.EventProcessedModel, error) {
	query := `
	INSERT INTO event_processed_model(event_pid, event_rety, event_processed)
	VALUES ($1, $2, $3)
	RETURNING id, event_pid, event_retry, event_processed, created_at
	`

	var model dbmodels.EventProcessedModel
	if err := tx.GetContext(ctx, &model, query, event_pid, event_retry, event_processed); err != nil {
		return nil, err
	}
	return &model, nil

}

func (repo *PostgresEventProcessedDBRepo) ReadByEventPID(ctx context.Context, event_pid int64) (*dbmodels.EventProcessedModel, error) {
	query := `
	SELECT * FROM event_processed_model
	WHERE event_pid=$1
	`

	var model dbmodels.EventProcessedModel
	if err := repo.db.GetContext(ctx, &model, query, event_pid); err != nil {
		return nil, err
	}
	return &model, nil
}

func (repo *PostgresEventProcessedDBRepo) ReadByEventPIDTx(ctx context.Context, tx *sqlx.Tx, event_pid int64) (*dbmodels.EventProcessedModel, error) {
	query := `
	SELECT * FROM event_processed_model
	WHERE event_pid=$1
	`

	var model dbmodels.EventProcessedModel
	if err := tx.GetContext(ctx, &model, query, event_pid); err != nil {
		return nil, err
	}
	return &model, nil
}

func (repo *PostgresEventProcessedDBRepo) UpdateTx(ctx context.Context, tx *sqlx.Tx, event_pid int64, event_retry int, event_processed bool) (*dbmodels.EventProcessedModel, error) {
	query := `
	UPDATE event_processed_model
	SET event_retry=$1, event_processed=$2
	WHERE event_pid=$3
	RETURNING id, event_pid, event_retry, event_processed, created_at
	`

	var model dbmodels.EventProcessedModel
	if err := tx.GetContext(ctx, &model, query, event_retry, event_processed, event_pid); err != nil {
		return nil, err
	}
	return &model, nil
}
