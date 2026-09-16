package dto

import "time"

type GetTaskResponse struct {
	PID       string    `json:"pid"`
	Value     string    `json:"value"`
	CreatedAt time.Time `json:"created_at"`
	Completed bool      `json:"completed"`
}

type GetAllResponse struct {
	Tasks []GetTaskResponse `json:"tasks"`
}
