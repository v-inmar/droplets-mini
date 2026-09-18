package services

import (
	"context"
	"database/sql"
	dbrepo "droplets_mini/services/todo/internal/db_repo"
	dtomodels "droplets_mini/services/todo/internal/models/dto_models"
	"errors"
	"fmt"
	"time"
)

type TaskServiceRoot interface {
	CreateNewTaskItem(ctx context.Context, reqDto *dtomodels.CreateRequest) (*dtomodels.CreateResponse, error)
	ReadAllTasks(ctx context.Context) (*dtomodels.GetAllResponse, error)
	ReadTaskByPID(ctx context.Context, pid string) (*dtomodels.GetTaskResponse, error)
	UpdateTask(ctx context.Context, task *dtomodels.UpdateRequest, pid string) (*dtomodels.UpdateResponse, error)
	DeleteTask(ctx context.Context, pid string) error
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

func (s *TaskService) ReadAllTasks(ctx context.Context) (*dtomodels.GetAllResponse, error) {
	taskModels, err := s.repo.ReadAll(ctx)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	taskLists := []dtomodels.GetTaskResponse{}

	for _, t := range taskModels {
		taskLists = append(taskLists, dtomodels.GetTaskResponse{
			PID:       t.PID,
			Value:     t.Value,
			Completed: t.Completed,
			CreatedAt: t.CreatedAt,
		})
	}

	all := dtomodels.GetAllResponse{
		Tasks: taskLists,
	}

	return &all, nil
}

func (s *TaskService) ReadTaskByPID(ctx context.Context, pid string) (*dtomodels.GetTaskResponse, error) {
	taskModel, err := s.repo.ReadByPID(ctx, pid)
	if err != nil {
		return nil, err
	}

	return &dtomodels.GetTaskResponse{
		PID:       taskModel.PID,
		Value:     taskModel.Value,
		Completed: taskModel.Completed,
		CreatedAt: taskModel.CreatedAt,
	}, nil
}

func (s *TaskService) UpdateTask(ctx context.Context, task *dtomodels.UpdateRequest, pid string) (*dtomodels.UpdateResponse, error) {
	db := s.repo.GetDBInstance()
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	taskModel, err := s.repo.UpdateTx(ctx, tx, task.Value, pid, task.Completed)
	if err != nil {
		return nil, err
	}

	// TODO: publish to event

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &dtomodels.UpdateResponse{
		PID:       taskModel.PID,
		Value:     taskModel.Value,
		CreatedAt: taskModel.CreatedAt,
		Completed: taskModel.Completed,
	}, nil

}

func (s *TaskService) DeleteTask(ctx context.Context, pid string) error {
	db := s.repo.GetDBInstance()
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := s.repo.DeleteTx(ctx, tx, pid); err != nil {
		return nil
	}

	// TODO: publish event

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
