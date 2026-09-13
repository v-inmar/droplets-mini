package service

import (
	"context"
	"droplets_mini/app/internal/models"
	"droplets_mini/pkg/events"
	"droplets_mini/pkg/subjects"
	"encoding/json"
	"time"

	"github.com/nats-io/nats.go/jetstream"
)

type JetstreamService struct {
	js jetstream.JetStream
}

func NewJetstreamService(js jetstream.JetStream) *JetstreamService {
	return &JetstreamService{
		js: js,
	}
}

func (service *JetstreamService) PublishTodoCreated(ctx context.Context, todo *models.Todo) error {
	event := events.TodoCreatedEvent{
		EventID:       time.Now().UnixNano(),
		TodoID:        todo.ID,
		TodoValue:     todo.Value,
		TodoCompleted: todo.Completed,
		TodoCreatedAt: todo.CreatedAt,
	}

	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	_, err = service.js.Publish(ctx, subjects.SubjectTodoCreated, data)
	if err != nil {
		return err
	}

	return nil
}

func (service *JetstreamService) CreateOrUpdateStream(ctx context.Context, streamName string) error {
	cfg := jetstream.StreamConfig{
		Name:     streamName,
		Subjects: []string{"todo.>"}, // > wildcard
		Storage:  jetstream.FileStorage,
	}

	_, err := service.js.CreateOrUpdateStream(ctx, cfg)
	if err != nil {
		return err
	}

	// configStreamName := stream.CachedInfo().Config.Name
	// configStreamSubjects := stream.CachedInfo().Config.Subjects
	// configStreamStorage := stream.CachedInfo().Config.Storage
	// log.Printf("[BROKER] Jetstream stream %s ready (subjects=%v, storage=%v)\n", configStreamName, configStreamSubjects, configStreamStorage)

	return nil
}
