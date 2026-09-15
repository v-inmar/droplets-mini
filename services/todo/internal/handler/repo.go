package handler

import "github.com/jmoiron/sqlx"

type HandlerRepo struct {
	db *sqlx.DB
}

func NewHandlerRepo(db *sqlx.DB) *HandlerRepo {
	return &HandlerRepo{
		db: db,
	}
}
