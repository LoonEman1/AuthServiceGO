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

func (p *KafkaProducer) EnsureTopicExists(ctx context.Context, brokers []string, topic string) error {
	address := brokers[0]
	conn, err := kafka.DialContext(ctx, "tcp", address)
	if err != nil {
		return err
	}
	defer conn.Close()

	partitions, err := conn.ReadPartitions()
	if err == nil {
		for _, p := range partitions {
			if p.Topic == topic {
				return nil
			}
		}
	}

	topicConfig := kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     1,
		ReplicationFactor: 1,
	}

	return conn.CreateTopics(topicConfig)
}
