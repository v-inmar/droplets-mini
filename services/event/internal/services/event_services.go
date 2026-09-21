package services

import (
	"bytes"
	"context"
	dbrepo "droplets_mini_eventservice/internal/db_repo"
	dbmodels "droplets_mini_eventservice/internal/models/db_models"
	eventmodels "droplets_mini_eventservice/internal/models/event_models"
	"droplets_mini_eventservice/internal/utils"
	"encoding/json"
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

func (s *ProcessEventService) CreateEventProcess(ctx context.Context, e *eventmodels.EventTask) (*dbmodels.EventProcessedModel, error) {
	tx, err := s.repo.GetDBInstance().BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	model, err := s.repo.CreateTx(ctx, tx, e.EventPID, 1, false)

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return model, nil

}

func (s *ProcessEventService) UpdateEventProcess(ctx context.Context, pid int64, retryCount int, processed bool) (*dbmodels.EventProcessedModel, error) {
	tx, err := s.repo.GetDBInstance().BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	model, err := s.repo.UpdateTx(ctx, tx, pid, retryCount, processed)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return model, nil

}

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

*/

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
