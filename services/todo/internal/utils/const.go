package utils

type EventType string

const (
	TopicTaskCreated EventType = "task.created"
	TopicTaskUpdated EventType = "task.updated"
	TopicTaskDeleted EventType = "task.deleted"
)
