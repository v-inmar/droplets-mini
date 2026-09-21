package services

import (
	"context"
	"database/sql"
	"droplets_mini_todoservice/internal/broker"
	dbrepo "droplets_mini_todoservice/internal/db_repo"
	dbmodels "droplets_mini_todoservice/internal/models/db_models"
	dtomodels "droplets_mini_todoservice/internal/models/dto_models"
	eventmodels "droplets_mini_todoservice/internal/models/event_models"
	"droplets_mini_todoservice/internal/utils"
	"errors"
	"fmt"
	"log"
	"time"
)

type TaskServiceRoot interface {
	CreateNewTaskItem(ctx context.Context, reqDto *dtomodels.CreateRequest) (*dtomodels.CreateResponse, error)
	ReadAllTasks(ctx context.Context) (*dtomodels.GetAllResponse, error)
	ReadTaskByPID(ctx context.Context, pid int64) (*dtomodels.GetTaskResponse, error)
	ReadTaskByPIDModel(ctx context.Context, pid int64) (*dbmodels.TaskModel, error)
	UpdateTask(ctx context.Context, taskUpdate *dtomodels.UpdateRequest, task *dtomodels.GetTaskResponse) (*dtomodels.UpdateResponse, error)
	DeleteTask(ctx context.Context, task *dbmodels.TaskModel) error
}

type TaskService struct {
	repo     dbrepo.TaskRepo
	producer broker.BrokerProducer
}

func NewTaskService(repo dbrepo.TaskRepo, producer broker.BrokerProducer) *TaskService {
	return &TaskService{
		repo:     repo,
		producer: producer,
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
	task, err := s.repo.CreateTx(ctx, tx, reqDto.Value, time.Now().UnixNano())
	if err != nil {
		return nil, err
	}

	eventPID := time.Now().UTC().UnixNano()
	eventModel := eventmodels.EventTask{
		Event:          utils.TopicTaskCreated,
		EventPID:       eventPID,
		EventCreatedAt: time.Now().UTC(),
		TaskID:         task.ID,
		TaskValue:      task.Value,
		TaskPID:        task.PID,
		TaskCompleted:  task.Completed,
		TaskCreatedAt:  task.CreatedAt,
	}

	key := fmt.Sprintf("tasks-%d", eventPID)
	log.Printf("[Todo] inside service, publishing task with key: %s", key)
	if err := s.producer.Produce(ctx, eventModel, []byte(key)); err != nil {
		return nil, err
	}

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

func (s *TaskService) ReadTaskByPID(ctx context.Context, pid int64) (*dtomodels.GetTaskResponse, error) {
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

func (s *TaskService) ReadTaskByPIDModel(ctx context.Context, pid int64) (*dbmodels.TaskModel, error) {
	taskModel, err := s.repo.ReadByPID(ctx, pid)
	if err != nil {
		return nil, err
	}

	return taskModel, nil
}

func (s *TaskService) UpdateTask(ctx context.Context, taskUpdate *dtomodels.UpdateRequest, task *dtomodels.GetTaskResponse) (*dtomodels.UpdateResponse, error) {
	db := s.repo.GetDBInstance()
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	taskModel, err := s.repo.UpdateTx(ctx, tx, taskUpdate.Value, task.PID, taskUpdate.Completed)
	if err != nil {
		return nil, err
	}

	eventPID := time.Now().UTC().UnixNano()
	eventModel := eventmodels.EventTask{
		Event:          utils.TopicTaskUpdated,
		EventPID:       eventPID,
		EventCreatedAt: time.Now().UTC(),
		TaskID:         taskModel.ID,
		TaskValue:      taskModel.Value,
		TaskPID:        taskModel.PID,
		TaskCompleted:  taskModel.Completed,
		TaskCreatedAt:  taskModel.CreatedAt,
	}

	key := fmt.Sprintf("tasks-%d", eventPID)
	if err := s.producer.Produce(ctx, eventModel, []byte(key)); err != nil {
		return nil, err
	}

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

func (s *TaskService) DeleteTask(ctx context.Context, task *dbmodels.TaskModel) error {
	db := s.repo.GetDBInstance()
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := s.repo.DeleteTx(ctx, tx, task.PID); err != nil {
		return err
	}

	eventPID := time.Now().UTC().UnixNano()
	eventModel := eventmodels.EventTask{
		Event:          utils.TopicTaskDeleted,
		EventPID:       eventPID,
		EventCreatedAt: time.Now().UTC(),
		TaskID:         task.ID,
		TaskValue:      task.Value,
		TaskPID:        task.PID,
		TaskCompleted:  task.Completed,
		TaskCreatedAt:  task.CreatedAt,
	}

	key := fmt.Sprintf("tasks-%d", eventPID)
	if err := s.producer.Produce(ctx, eventModel, []byte(key)); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
