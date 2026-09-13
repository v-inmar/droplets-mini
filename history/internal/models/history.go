package models

import "time"

type History struct {
	ID        int64     `db:"id"`
	Event     string    `db:"event"`
	TodoValue string    `db:"todo_value"`
	TodoID    int64     `db:"todo_id"`
	CreatedAt time.Time `db:"created_at"`
}
