package repo

import (
	"context"
	"droplets_mini/history/internal/models"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

var ErrEventIDExist = errors.New("event id exist")

type ProcessedRepoTx struct {
	tx  *sqlx.Tx
	ctx context.Context
}

func NewProcessedRepoTx(tx *sqlx.Tx, ctx context.Context) *ProcessedRepoTx {
	return &ProcessedRepoTx{
		tx:  tx,
		ctx: ctx,
	}
}

func (p *ProcessedRepoTx) Create(eventID int64) (*models.ProcessedEvent, error) {
	var proc models.ProcessedEvent

	query := `
	INSERT INTO processed_event_model (event_id)
	VALUES ($1)
	RETURNING id, event_id, created_at
	`

	if err := p.tx.GetContext(p.ctx, &proc, query, eventID); err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" && pqErr.Constraint == "processed_event_model_eventid_unique" {
			return nil, ErrEventIDExist
		}
		return nil, err
	}

	return &proc, nil
}
