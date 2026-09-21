package broker

import (
	"context"
	eventmodels "droplets_mini_eventservice/internal/models/event_models"

	"github.com/segmentio/kafka-go"
)

type BrokerProducer interface {
	Produce(ctx context.Context, model eventmodels.EventTask, key []byte) error
	Close() error
}

type BrokerConsumer interface {
	Consume(ctx context.Context) (*kafka.Message, error)
	Commit(ctx context.Context, msg *kafka.Message) error
	Close() error
}
