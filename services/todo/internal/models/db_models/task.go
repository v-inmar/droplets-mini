package dbmodels

import "time"

type TaskModel struct {
	ID        int64     `db:"id"`
	PID       int64     `db:"pid"`
	Value     string    `db:"value"`
	Completed bool      `db:"completed"`
	CreatedAt time.Time `db:"created_at"`
}
