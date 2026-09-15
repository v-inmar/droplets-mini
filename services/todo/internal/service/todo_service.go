package service

import (
	"context"
	"droplets_mini/services/todo/internal/models/db"
	"droplets_mini/services/todo/internal/repo"

	"github.com/jmoiron/sqlx"
)

type TodoService struct {
	db  *sqlx.DB
	ctx context.Context
}

func NewTodoService(db *sqlx.DB) *TodoService {
	return &TodoService{
		db: db,
	}
}

func (service *TodoService) CreateNewTodoItem(value string) (*db.TodoModel, error) {
	tx, err := service.db.BeginTxx(service.ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	repoTx := repo.NewTodoRepoTx(tx)
	todo, err := repoTx.Create(service.ctx, value)
	if err != nil {
		return nil, err
	}

	//publish event here

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return todo, nil
}
