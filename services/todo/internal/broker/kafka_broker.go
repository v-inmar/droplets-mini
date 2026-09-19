package broker

import (
	"context"
	eventmodels "droplets_mini_todoservice/internal/models/event_models"
	"encoding/json"
	"errors"

	"github.com/segmentio/kafka-go"
)

type KafkaBroker struct {
	writer *kafka.Writer
	reader *kafka.Reader
}

func NewKafkaBroker(w *kafka.Writer, r *kafka.Reader) *KafkaBroker {
	return &KafkaBroker{
		writer: w,
		reader: r,
	}
}

func (b *KafkaBroker) Produce(ctx context.Context, eventModel eventmodels.TaskEvent, key []byte) error {
	data, err := json.Marshal(eventModel)
	if err != nil {
		return err
	}

	return b.writer.WriteMessages(ctx, kafka.Message{
		Key:   key,
		Value: data,
	})
}

func (b *KafkaBroker) Consume() error {
	return errors.New("not yet implemented")
}

func (b *KafkaBroker) CloseProcuder() error {
	return b.writer.Close()
}

func (b *KafkaBroker) CloseConsumer() error {
	return b.reader.Close()
}
