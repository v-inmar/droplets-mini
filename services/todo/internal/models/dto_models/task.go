package dtomodels

import "time"

type CreateRequest struct {
	Value string `json:"value"`
}

type CreateResponse struct {
	PID       int64     `json:"pid"`
	Value     string    `json:"value"`
	CreatedAt time.Time `json:"created_at"`
	Completed bool      `json:"completed"`
}

type GetTaskResponse struct {
	PID       int64     `json:"pid"`
	Value     string    `json:"value"`
	CreatedAt time.Time `json:"created_at"`
	Completed bool      `json:"completed"`
}

type GetAllResponse struct {
	Tasks []GetTaskResponse `json:"tasks"`
}

type UpdateRequest struct {
	Value     string `json:"value"`
	Completed bool   `json:"completed"`
}

type UpdateResponse struct {
	PID       int64     `json:"pid"`
	Value     string    `json:"value"`
	CreatedAt time.Time `json:"created_at"`
	Completed bool      `json:"completed"`
}
