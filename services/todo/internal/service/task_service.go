package service

import (
	"context"
	"droplets_mini/services/todo/internal/models/db"
	"droplets_mini/services/todo/internal/models/dto"
	"droplets_mini/services/todo/internal/repo"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type TaskService struct {
	db  *sqlx.DB
	ctx context.Context
}

func NewTaskService(db *sqlx.DB, ctx context.Context) *TaskService {
	return &TaskService{
		db:  db,
		ctx: ctx,
	}
}

func (service *TaskService) CreateNewTaskItem(req *dto.CreateRequest) (*db.TaskModel, error) {
	tx, err := service.db.BeginTxx(service.ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	repoTx := repo.NewTaskRepoTx(tx)

	// pid collision is a real problem here
	// but for this, its ok, it wont be deploy to multi users
	task, err := repoTx.Create(service.ctx, req.Value, fmt.Sprintf("%d", time.Now().UnixNano()))
	if err != nil {
		return nil, err
	}

	// TODO: publish event here

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return task, nil
}

func (service *TaskService) ReadAllTasks() (*dto.GetAllResponse, error) {
	repoDB := repo.NewTaskRepoDB(service.db)
	tasks, err := repoDB.ReadAll(service.ctx)
	if err != nil {
		return nil, err
	}

	taskLists := []dto.GetTaskResponse{}
	if len(tasks) > 0 {
		for _, t := range tasks {
			taskLists = append(taskLists, dto.GetTaskResponse{
				PID:       t.PID,
				Value:     t.Value,
				Completed: t.Completed,
				CreatedAt: t.CreatedAt,
			})
		}
	}

	all := dto.GetAllResponse{
		Tasks: taskLists,
	}

	return &all, nil
}

func (service *TaskService) ReadTaskByPID(pid string) (*dto.GetTaskResponse, error) {
	repoDB := repo.NewTaskRepoDB(service.db)
	task, err := repoDB.ReadByPID(service.ctx, pid)
	if err != nil {
		return nil, err
	}

	return &dto.GetTaskResponse{
		PID:       task.PID,
		Value:     task.Value,
		Completed: task.Completed,
		CreatedAt: task.CreatedAt,
	}, nil
}

func (service *TaskService) UpdateTask(task *dto.UpdateRequest, pid string) (*dto.UpdateResponse, error) {
	tx, err := service.db.BeginTxx(service.ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	repoTx := repo.NewTaskRepoTx(tx)

	taskModel, err := repoTx.Update(service.ctx, task.Value, pid, task.Completed)
	if err != nil {
		return nil, err
	}

	// TODO: publish to event

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	resp := dto.UpdateResponse{
		PID:       taskModel.PID,
		Value:     taskModel.Value,
		Completed: taskModel.Completed,
		CreatedAt: taskModel.CreatedAt,
	}

	return &resp, nil

}

func (service *TaskService) DeleteTask(pid string) error {
	tx, err := service.db.BeginTxx(service.ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	repoTx := repo.NewTaskRepoTx(tx)

	if err := repoTx.Delete(service.ctx, pid); err != nil {
		return err
	}

	// TODO: publish event

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
