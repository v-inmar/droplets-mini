package events

import "time"

type TodoCreatedEvent struct {
	EventID       int64     `json:"event_id"`
	TodoID        int64     `json:"todo_id"`
	TodoValue     string    `json:"todo_value"`
	TodoCompleted bool      `json:"todo_completed"`
	TodoCreatedAt time.Time `json:"todo_created_at"`
}
