package service

import (
	"context"
	"droplets_mini/app/internal/models"
	"droplets_mini/app/internal/repo"

	"github.com/jmoiron/sqlx"
)

type DBService struct {
	DB *sqlx.DB
}

func NewDBService(db *sqlx.DB) *DBService {
	return &DBService{
		DB: db,
	}
}

func (s *DBService) CreateNewTodo(ctx context.Context, value string) (*models.Todo, error) {
	tx, err := s.DB.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	repoTx := repo.NewTodoRepoTx(ctx, tx)

	todo, err := repoTx.Create(value)
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return todo, nil
}

func (s *DBService) CreateNewTodoAndPublisEvent(ctx context.Context, value string) (*models.Todo, error) {
	tx, err := s.DB.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	repoTx := repo.NewTodoRepoTx(ctx, tx)
	todo, err := repoTx.Create(value)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return todo, nil
}
