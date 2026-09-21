package eventmodels

import (
	"droplets_mini_eventservice/internal/utils"
	"time"
)

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
