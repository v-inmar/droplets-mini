package handler

import (
	"droplets_mini/services/todo/internal/service"

	"github.com/jmoiron/sqlx"
)

type HandlerRepo struct {
	db      *sqlx.DB
	service *service.TaskService
}

func NewHandlerRepo(db *sqlx.DB, srvc *service.TaskService) *HandlerRepo {
	return &HandlerRepo{
		db:      db,
		service: srvc,
	}
}
