package dbmodels

import "time"

type EventProcessedModel struct {
	ID             int64     `db:"id"`
	EventPID       int64     `db:"event_pid"`
	EventRetry     int       `db:"event_retry"`
	EventProcessed bool      `db:"event_processed"`
	CreatedAt      time.Time `db:"created_at"`
}
