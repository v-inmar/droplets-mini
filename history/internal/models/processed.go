package models

import "time"

type ProcessedEvent struct {
	ID        int64     `db:"id"`
	EventID   int64     `db:"event_id"`
	CreatedAt time.Time `db:"created_at"`
}
