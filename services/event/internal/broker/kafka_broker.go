package broker

import (
	"context"
	eventmodels "droplets_mini_eventservice/internal/models/event_models"
	"encoding/json"

	"github.com/segmentio/kafka-go"
)

type KafkaBrokerProducer struct {
	writer *kafka.Writer
}

func NewKafkaBrokerProducer(w *kafka.Writer) *KafkaBrokerProducer {
	return &KafkaBrokerProducer{
		writer: w,
	}
}

func (k *KafkaBrokerProducer) Produce(ctx context.Context, model eventmodels.EventTask, key []byte) error {
	data, err := json.Marshal(model)
	if err != nil {
		return err
	}

	return k.writer.WriteMessages(ctx, kafka.Message{
		Key:   key,
		Value: data,
	})
}

func (k *KafkaBrokerProducer) Close() error {
	return k.writer.Close()
}

type KafkaBrokerConsumer struct {
	reader *kafka.Reader
}

func NewKafkaBrokerConsumer(r *kafka.Reader) *KafkaBrokerConsumer {
	return &KafkaBrokerConsumer{
		reader: r,
	}
}

func (k *KafkaBrokerConsumer) Consume(ctx context.Context) (*kafka.Message, error) {
	msg, err := k.reader.FetchMessage(ctx)
	if err != nil {
		return nil, err
	}

	return &msg, nil
}

func (k *KafkaBrokerConsumer) Commit(ctx context.Context, msg *kafka.Message) error {
	return k.reader.CommitMessages(ctx, *msg)
}

func (k *KafkaBrokerConsumer) Close() error {
	return k.reader.Close()
}
