package models

import "time"

// todo_model
type Todo struct {
	ID        int64     `db:"id"`
	Value     string    `db:"value"`
	Completed bool      `db:"completed"`
	CreatedAt time.Time `db:"created_at"`
}
