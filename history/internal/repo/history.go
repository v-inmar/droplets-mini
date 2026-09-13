package repo

import (
	"context"
	"droplets_mini/history/internal/models"

	"github.com/jmoiron/sqlx"
)

type HistoryRepoTx struct {
	tx  *sqlx.Tx
	ctx context.Context
}

func NewHistoryRepoTx(tx *sqlx.Tx, ctx context.Context) *HistoryRepoTx {
	return &HistoryRepoTx{
		tx:  tx,
		ctx: ctx,
	}
}

func (repo *HistoryRepoTx) Create(event, todoValue string, todoID int64) (*models.History, error) {
	var h models.History

	query := `
	INSERT INTO history_model (event, todo_value, todo_id)
	VALUES ($1, $2, $3)
	RETURNING id, event, todo_value, todo_id, created_at
	`

	if err := repo.tx.GetContext(repo.ctx, &h, query, event, todoValue, todoID); err != nil {
		return nil, err
	}

	return &h, nil
}

type HistoryRepo struct {
	db  *sqlx.DB
	ctx context.Context
}

func NewHistoryRepo(db *sqlx.DB, ctx context.Context) *HistoryRepo {
	return &HistoryRepo{
		db:  db,
		ctx: ctx,
	}
}

func (repo *HistoryRepo) ReadAll() ([]models.History, error) {
	query := `
	SELECT id, event, todo_value, todo_id, created_at
	FROM history_model
	`
	var hs []models.History
	if err := repo.db.SelectContext(repo.ctx, &hs, query); err != nil {
		return nil, err
	}

	return hs, nil
}
