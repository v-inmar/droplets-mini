package dto

import (
	"fmt"
	"time"
)

type HealthCheckResponse struct {
	Message string    `json:"message"`
	Status  int       `json:"status"`
	Time    time.Time `json:"time"`
}

func (dto *HealthCheckResponse) Repr() string {
	return fmt.Sprintf("Message: %s - Status: %d - Time: %v", dto.Message, dto.Status, dto.Time)
}
