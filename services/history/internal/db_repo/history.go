package dbrepo

import (
	"context"
	dbmodels "droplets_mini_historyservice/internal/models/db_models"

	"github.com/jmoiron/sqlx"
)

type HistoryRepo interface {
	ReadAll(ctx context.Context) ([]dbmodels.HistoryModel, error)
}

type PostgresHistoryDBRepo struct {
	db *sqlx.DB
}

func NewPostgresHistoryDBRepo(db *sqlx.DB) *PostgresHistoryDBRepo {
	return &PostgresHistoryDBRepo{
		db: db,
	}
}

func (repo *PostgresHistoryDBRepo) ReadAll(ctx context.Context) ([]dbmodels.HistoryModel, error) {
	query := `
	SELECT * FROM history_model
	`

	var histories []dbmodels.HistoryModel

	if err := repo.db.SelectContext(ctx, &histories, query); err != nil {
		return nil, err
	}
	return histories, nil
}
