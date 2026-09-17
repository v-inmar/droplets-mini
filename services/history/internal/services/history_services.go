package services

import (
	"context"
	"database/sql"
	dbrepo "droplets_mini/services/history/internal/db_repo"
	dtomodels "droplets_mini/services/history/internal/models/dto_models"
	"errors"
)

type HistoryServiceRoot interface {
	GetAllHistory(ctx context.Context) (*dtomodels.GetAllHistoryResponse, error)
}

type HistoryService struct {
	repo dbrepo.HistoryRepo
}

func NewHistoryService(repo dbrepo.HistoryRepo) *HistoryService {
	return &HistoryService{
		repo: repo,
	}
}

func (s *HistoryService) GetAllHistory(ctx context.Context) (*dtomodels.GetAllHistoryResponse, error) {
	itemModels, err := s.repo.ReadAll(ctx)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}

	itemResps := []dtomodels.GetHistoryResponse{}
	for _, h := range itemModels {
		itemResps = append(itemResps, dtomodels.GetHistoryResponse{
			EventPID:        h.EventPID,
			EventHappened:   h.EventHappened,
			EventHappenedAt: h.EventHappenedAt,
			TaskValue:       h.TaskValue,
			TaskPID:         h.TaskPID,
		})
	}

	return &dtomodels.GetAllHistoryResponse{History: itemResps}, nil
}
