package broker

import (
	"context"
	eventmodels "droplets_mini_todoservice/internal/models/event_models"
)

type Broker interface {
	Produce(ctx context.Context, eventModel eventmodels.EventTask, key []byte) error
	Consume() error
	CloseProcuder() error
	CloseConsumer() error
}

type BrokerProducer interface {
	Produce(ctx context.Context, model eventmodels.EventTask, key []byte) error
	Close() error
}
