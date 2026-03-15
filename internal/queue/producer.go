package queue

import (
	"AuthService/internal/models"
	"context"
	"encoding/json"
	"errors"

	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	writer *kafka.Writer
}

func NewKafkaProducer(brokers []string, topic string) *KafkaProducer {
	return &KafkaProducer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		},
	}
}

func (p *KafkaProducer) SendEmailTask(ctx context.Context, task models.EmailTask) error {
	payload, err := json.Marshal(task)
	if err != nil {
		return errors.New("Ошибка преобразования в JSON")
	}

	err = p.writer.WriteMessages(ctx, kafka.Message{
		Value: payload,
	})
	if err != nil {
		return err
	}

	return nil
}

func (p *KafkaProducer) Close() error {
	return p.writer.Close()
}
