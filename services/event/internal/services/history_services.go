package services

import (
	"context"
	dbrepo "droplets_mini_eventservice/internal/db_repo"
	eventmodels "droplets_mini_eventservice/internal/models/event_models"
	"droplets_mini_eventservice/internal/utils"
)

type HistoryService struct {
	repo dbrepo.HistoryRepo
}

func NewHistoryService(repo dbrepo.HistoryRepo) *HistoryService {
	return &HistoryService{
		repo: repo,
	}
}

/*

// type HistoryItem struct {
// 	ID              int64     `db:"id"`
// 	EventPID        int64     `db:"event_pid"`
// 	EventHappened   string    `db:"event_happened"`
// 	EventHappenedAt time.Time `db:"event_happened_at"`
// 	TaskID          int64     `db:"task_id"`
// 	TaskValue       string    `db:"task_value"`
// 	TaskPID         int64     `db:"task_pid"`
// 	TaskCompleted   bool      `db:"task_completed"`
// 	CreatedAt       time.Time `db:"created_at"`
// }

*/

func (s *HistoryService) CreateHistoryItem(ctx context.Context, event *eventmodels.EventTask) error {
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

	tx, err := s.repo.GetDBInstance().BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = s.repo.CreateTx(
		ctx,
		tx,
		event.EventPID,
		eventHappened,
		event.EventCreatedAt,
		event.TaskID,
		event.TaskValue,
		event.TaskPID,
		event.TaskCompleted,
	)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil

}
