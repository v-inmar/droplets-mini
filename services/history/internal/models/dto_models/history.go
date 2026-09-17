package dtomodels

import "time"

type GetHistoryResponse struct {
	EventPID        string    `json:"event_pid"`
	EventHappened   string    `json:"event_happened"`
	EventHappenedAt time.Time `json:"event_happened_at"`
	TaskValue       string    `json:"task_value"`
	TaskPID         string    `json:"task_pid"`
}

type GetAllHistoryResponse struct {
	History []GetHistoryResponse `json:"history"`
}
