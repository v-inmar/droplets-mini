package services

import (
	"context"
	dbrepo "droplets_mini/services/todo/internal/db_repo"
	dtomodels "droplets_mini/services/todo/internal/models/dto_models"
	"fmt"
	"time"
)

type TaskServiceRoot interface {
	CreateNewTaskItem(ctx context.Context, reqDto *dtomodels.CreateRequest) (*dtomodels.CreateResponse, error)
	ReadAllTasks() (*dtomodels.GetAllResponse, error)
}

type TaskService struct {
	repo dbrepo.TaskRepo
}

func NewTaskService(repo dbrepo.TaskRepo) *TaskService {
	return &TaskService{
		repo: repo,
	}
}

func (s *TaskService) CreateNewTaskItem(ctx context.Context, reqDto *dtomodels.CreateRequest) (*dtomodels.CreateResponse, error) {
	db := s.repo.GetDBInstance()
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// pid collision is a real problem here
	// but for this, its ok, it wont be deploy to multi users
	task, err := s.repo.CreateTx(ctx, tx, reqDto.Value, fmt.Sprintf("%d", time.Now().UnixNano()))
	if err != nil {
		return nil, err
	}

	// TODO: publish event here

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &dtomodels.CreateResponse{
		PID:       task.PID,
		Value:     task.Value,
		CreatedAt: task.CreatedAt,
		Completed: task.Completed,
	}, nil
}

func (s *TaskService) ReadAllTasks() (*dtomodels.GetAllResponse, error) {
	return nil, nil
}
