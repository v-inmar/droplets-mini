package repo

import (
	"context"
	"droplets_mini/services/todo/internal/models/db"

	"github.com/jmoiron/sqlx"
)

type TodoRepoTx struct {
	tx *sqlx.Tx
}

func NewTodoRepoTx(tx *sqlx.Tx) *TodoRepoTx {
	return &TodoRepoTx{
		tx: tx,
	}
}

func (repo *TodoRepoTx) Create(ctx context.Context, value string) (*db.TodoModel, error) {
	query := `
	INSERT INTO todo_model (value)
	VALUES ($1)
	RETURNING id, value, completed, create_at
	`

	var todo db.TodoModel
	if err := repo.tx.GetContext(ctx, &todo, query, value); err != nil {
		return nil, err
	}

	return &todo, nil
}
