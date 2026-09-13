package repo

import (
	"context"
	"droplets_mini/app/internal/models"

	"github.com/jmoiron/sqlx"
)

type TodoRepoTx struct {
	Tx  *sqlx.Tx
	Ctx context.Context
}

func NewTodoRepoTx(ctx context.Context, tx *sqlx.Tx) *TodoRepoTx {
	return &TodoRepoTx{
		Tx:  tx,
		Ctx: ctx,
	}
}

func (t *TodoRepoTx) Create(value string) (*models.Todo, error) {
	var todo models.Todo

	query := `
	INSERT INTO todo_model (value)
	VALUES ($1)
	RETURNING id, value, completed, created_at
	`

	if err := t.Tx.GetContext(t.Ctx, &todo, query, value); err != nil {
		return nil, err
	}

	return &todo, nil
}
