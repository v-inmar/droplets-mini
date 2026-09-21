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
