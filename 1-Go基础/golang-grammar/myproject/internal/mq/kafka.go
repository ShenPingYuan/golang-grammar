package mq

import (
	"context"
	"log/slog"
)

// KafkaProducer Kafka 生产者
// 生产环境请使用 github.com/segmentio/kafka-go 或 github.com/IBM/sarama
type KafkaProducer struct {
	brokers []string
}

func NewKafkaProducer(brokers []string) Producer {
	return &KafkaProducer{brokers: brokers}
}

func (p *KafkaProducer) Publish(_ context.Context, topic string, message []byte) error {
	slog.Info("[Kafka] publish", "topic", topic, "size", len(message))
	// TODO: 使用 kafka-go writer.WriteMessages()
	return nil
}

func (p *KafkaProducer) Close() error {
	slog.Info("[Kafka] producer closed")
	return nil
}

// KafkaConsumer Kafka 消费者
type KafkaConsumer struct {
	brokers  []string
	handlers map[string]MessageHandler
}

func NewKafkaConsumer(brokers []string) Consumer {
	return &KafkaConsumer{brokers: brokers, handlers: make(map[string]MessageHandler)}
}

func (c *KafkaConsumer) Subscribe(topic string, handler MessageHandler) error {
	c.handlers[topic] = handler
	return nil
}

func (c *KafkaConsumer) Start(ctx context.Context) error {
	slog.Info("[Kafka] consumer started", "topics", len(c.handlers))
	<-ctx.Done()
	return nil
}

func (c *KafkaConsumer) Close() error {
	slog.Info("[Kafka] consumer closed")
	return nil
}