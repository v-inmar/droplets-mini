package dtomodels

import "time"

type CreateRequest struct {
	Value string `json:"value"`
}

type CreateResponse struct {
	PID       string    `json:"pid"`
	Value     string    `json:"value"`
	CreatedAt time.Time `json:"created_at"`
	Completed bool      `json:"completed"`
}

type GetTaskResponse struct {
	PID       string    `json:"pid"`
	Value     string    `json:"value"`
	CreatedAt time.Time `json:"created_at"`
	Completed bool      `json:"completed"`
}

type GetAllResponse struct {
	Tasks []GetTaskResponse `json:"tasks"`
}
