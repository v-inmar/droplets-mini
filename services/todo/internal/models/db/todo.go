package db

import "time"

type TodoModel struct {
	ID        int64     `db:"id"`
	Value     string    `db:"value"`
	Completed bool      `db:"completed"`
	CreatedAt time.Time `db:"created_at"`
}
