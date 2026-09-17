package dbmodels

import "time"

type HistoryModel struct {
	ID              int64     `db:"id"`
	EventPID        string    `db:"event_pid"`
	EventHappened   string    `db:"event_happened"`
	EventHappenedAt time.Time `db:"event_happened_at"`
	TaskValue       string    `db:"task_value"`
	TaskPID         string    `db:"task_pid"`
	CreatedAt       time.Time `db:"created_at"`
}
