package dtomodels

import "time"

type GetHealthResponse struct {
	Message string    `json:"message"`
	Status  int       `json:"status"`
	Time    time.Time `json:"time"`
}
