package dtomodels

import "time"

type GetHistoryResponse struct {
	EventPID        int64     `json:"event_pid"`
	EventHappened   string    `json:"event_happened"`
	EventHappenedAt time.Time `json:"event_happened_at"`
	TaskValue       string    `json:"task_value"`
	TaskPID         int64     `json:"task_pid"`
	TaskCompleted   bool      `json:"task_completed"`
}

type GetAllHistoryResponse struct {
	History []GetHistoryResponse `json:"history"`
}

/*
type HistoryModel struct {
	ID              int64     `db:"id"`
	EventPID        int64     `db:"event_pid"`
	EventHappened   string    `db:"event_happened"`
	EventHappenedAt time.Time `db:"event_happened_at"`
	TaskID          int64     `db:"task_id"`
	TaskValue       string    `db:"task_value"`
	TaskPID         int64     `db:"task_pid"`
	TaskCompleted   bool      `db:"task_completed"`
	CreatedAt       time.Time `db:"created_at"`
}

*/
