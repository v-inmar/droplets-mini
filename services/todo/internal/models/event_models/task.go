package eventmodels

import (
	"droplets_mini_todoservice/internal/utils"
	"time"
)

type TaskEvent interface {
	Repr() string
}

type TaskEventCreated struct {
	Event         utils.EventType `json:"event_type"`
	EventID       int64           `json:"event_id"`
	TaskID        int64           `json:"task_id"`
	TaskValue     string          `json:"task_value"`
	TaskPID       string          `json:"task_pid"`
	TaskCompleted bool            `json:"task_completed"`
	TaskCreatedAt time.Time       `json:"task_created_at"`
}

func (t TaskEventCreated) Repr() string {
	return "No implementation yet..coming soon"
}

type TaskEventUpdated struct {
	Event            utils.EventType `json:"event_type"`
	EventID          int64           `json:"event_id"`
	OldTaskValue     string          `json:"old_task_value"`
	OldTaskCompleted bool            `json:"old_task_completed"`
	TaskID           int64           `json:"task_id"`
	TaskValue        string          `json:"task_value"`
	TaskPID          string          `json:"task_pid"`
	TaskCompleted    bool            `json:"task_completed"`
	TaskCreatedAt    time.Time       `json:"task_created_at"`
}

func (t TaskEventUpdated) Repr() string {
	return "No implementation yet..coming soon"
}

type TaskEventDeleted struct {
	Event                utils.EventType `json:"event_type"`
	EventID              int64           `json:"event_id"`
	DeletedTaskID        int64           `json:"deleted_task_id"`
	DeletedTaskValue     string          `json:"deleted_task_value"`
	DeletedTaskPID       string          `json:"deleted_task_pid"`
	DeletedTaskCompleted bool            `json:"deleted_task_completed"`
	DeletedTaskCreatedAt time.Time       `json:"deleted_task_created_at"`
}

func (t TaskEventDeleted) Repr() string {
	return "No implementation yet..coming soon"
}

// type TaskCreatedEvent struct {

// }

// type TaskUpdatedEvent struct {
// 	EventID       int64     `json:"event_id"`
// 	TodoID        int64     `json:"todo_id"`
// 	TodoValue     string    `json:"todo_value"`
// 	TodoCompleted bool      `json:"todo_completed"`
// 	TodoCreatedAt time.Time `json:"todo_created_at"`
// }

// type TaskDeletedEvent struct {
// 	EventID       int64     `json:"event_id"`
// 	TodoID        int64     `json:"todo_id"`
// 	TodoValue     string    `json:"todo_value"`
// 	TodoCompleted bool      `json:"todo_completed"`
// 	TodoCreatedAt time.Time `json:"todo_created_at"`
// }
