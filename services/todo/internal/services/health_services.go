package services

import (
	dtomodels "droplets_mini/services/todo/internal/models/dto_models"
	"net/http"
	"time"
)

type HealthServiceRoot interface {
	GetHealth() (*dtomodels.GetHealthResponse, error)
}

type HealthService struct{}

func NewHealthService() *HealthService {
	return &HealthService{}
}

func (service *HealthService) GetHealth() (*dtomodels.GetHealthResponse, error) {
	health := dtomodels.GetHealthResponse{
		Message: "todo service up and running",
		Status:  http.StatusOK,
		Time:    time.Now().UTC(),
	}

	return &health, nil
}
