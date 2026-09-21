package services

import (
	"bytes"
	"context"
	"database/sql"
	dbrepo "droplets_mini_eventservice/internal/db_repo"
	dbmodels "droplets_mini_eventservice/internal/models/db_models"
	eventmodels "droplets_mini_eventservice/internal/models/event_models"
	"droplets_mini_eventservice/internal/utils"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type ProcessEventService struct {
	repo dbrepo.ProcessedRepo
}

func NewProcessEventService(repo dbrepo.ProcessedRepo) *ProcessEventService {
	return &ProcessEventService{
		repo: repo,
	}
}

func (s *ProcessEventService) GetOrCreateEventProcess(ctx context.Context, pid int64) (*dbmodels.EventProcessedModel, error) {
	tx, err := s.repo.GetDBInstance().BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	model, err := s.repo.ReadByEventPIDTx(ctx, tx, pid)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	if model != nil {
		return model, nil
	}

	model, err = s.repo.CreateTx(ctx, tx, pid, 0, false)
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return model, nil
}

func (s *ProcessEventService) UpdateEventProcess(ctx context.Context, epid int64, retry int, processed bool) (*dbmodels.EventProcessedModel, error) {
	tx, err := s.repo.GetDBInstance().BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	model, err := s.repo.UpdateTx(ctx, tx, epid, retry, processed)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return model, nil

}

// func (s *ProcessEventService) CreateEventProcess(ctx context.Context, e *eventmodels.EventTask) (*dbmodels.EventProcessedModel, error) {
// 	tx, err := s.repo.GetDBInstance().BeginTxx(ctx, nil)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer tx.Rollback()

// 	model, err := s.repo.CreateTx(ctx, tx, e.EventPID, 1, false)

// 	if err := tx.Commit(); err != nil {
// 		return nil, err
// 	}
// 	return model, nil

// }

// func (s *ProcessEventService) UpdateEventProcess(ctx context.Context, pid int64, retryCount int, processed bool) (*dbmodels.EventProcessedModel, error) {
// 	tx, err := s.repo.GetDBInstance().BeginTxx(ctx, nil)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer tx.Rollback()

// 	model, err := s.repo.UpdateTx(ctx, tx, pid, retryCount, processed)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if err := tx.Commit(); err != nil {
// 		return nil, err
// 	}
// 	return model, nil

// }

/*



type EventTask struct {
	Event          utils.EventType `json:"event_type"`
	EventPID       int64           `json:"event_pid"`
	EventCreatedAt time.Time       `json:"event_created_at"`
	TaskID         int64           `json:"task_id"`
	TaskValue      string          `json:"task_value"`
	TaskPID        int64           `json:"task_pid"`
	TaskCompleted  bool            `json:"task_completed"`
	TaskCreatedAt  time.Time       `json:"task_created_at"`
}


type HistoryItem struct {
	ID              int64     `db:"id"`
	EventPID        int64     `db:"event_pid"`
	EventHappened   string    `db:"event_happened"`
	EventHappenedAt time.Time `db:"event_happened_at"`
	TaskID          int64     `db:"task_id"`
	TaskValue       string    `db:"task_value"`
	TaskPID         int64     `db:"task_pid"`
	CreatedAt       time.Time `db:"created_at"`
}

eventPID int64, eventHappened string, eventHappenedAt time.Time, taskID int64, taskPID int64, taskValue string

*/

// func (s *ProcessEventService) CreateHistoryItem(ctx context.Context, event *eventmodels.EventTask)error{
// 	var eventHappened string
// 	switch event.Event {
// 	case utils.TopicTaskCreated:
// 		eventHappened = "CREATED"
// 	case utils.TopicTaskUpdated:
// 		eventHappened = "UPDATED"
// 	case utils.TopicTaskDeleted:
// 		eventHappened = "DELETED"
// 	default:
// 		eventHappened = "UNKNOWN"
// 	}

// 	s.repo
// }

func (s *ProcessEventService) PostHistoryAPI(ctx context.Context, url string, event eventmodels.EventTask) error {
	var eventHappened string
	switch event.Event {
	case utils.TopicTaskCreated:
		eventHappened = "CREATED"
	case utils.TopicTaskUpdated:
		eventHappened = "UPDATED"
	case utils.TopicTaskDeleted:
		eventHappened = "DELETED"
	default:
		eventHappened = "UNKNOWN"
	}

	body := struct {
		EventHappened  string    `json:"event_happened"`
		EventPID       int64     `json:"event_pid"`
		EventCreatedAt time.Time `json:"event_created_at"`
		TaskID         int64     `json:"task_id"`
		TaskValue      string    `json:"task_value"`
		TaskPID        int64     `json:"task_pid"`
		TaskCompleted  bool      `json:"task_completed"`
		TaskCreatedAt  time.Time `json:"task_created_at"`
	}{
		EventHappened:  eventHappened,
		EventPID:       event.EventPID,
		EventCreatedAt: event.EventCreatedAt,
		TaskID:         event.TaskID,
		TaskValue:      event.TaskValue,
		TaskPID:        event.TaskPID,
		TaskCompleted:  event.TaskCompleted,
		TaskCreatedAt:  event.TaskCreatedAt,
	}

	data, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("status code returned: %d", resp.StatusCode)
	}

	return nil
}
