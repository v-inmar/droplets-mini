package dto

import "time"

type UpdateRequest struct {
	Value     string `json:"value"`
	Completed bool   `json:"completed"`
}

type UpdateResponse struct {
	PID       string    `json:"pid"`
	Value     string    `json:"value"`
	CreatedAt time.Time `json:"created_at"`
	Completed bool      `json:"completed"`
}
