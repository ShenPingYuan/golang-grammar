package mq

import (
	"context"
	"log/slog"
)

// RabbitMQProducer RabbitMQ 生产者
// 生产环境请使用 github.com/rabbitmq/amqp091-go
type RabbitMQProducer struct {
	url string
}

func NewRabbitMQProducer(url string) Producer {
	return &RabbitMQProducer{url: url}
}

func (p *RabbitMQProducer) Publish(_ context.Context, topic string, message []byte) error {
	slog.Info("[RabbitMQ] publish", "queue", topic, "size", len(message))
	return nil
}

func (p *RabbitMQProducer) Close() error { return nil }

type RabbitMQConsumer struct {
	url      string
	handlers map[string]MessageHandler
}

func NewRabbitMQConsumer(url string) Consumer {
	return &RabbitMQConsumer{url: url, handlers: make(map[string]MessageHandler)}
}

func (c *RabbitMQConsumer) Subscribe(topic string, handler MessageHandler) error {
	c.handlers[topic] = handler
	return nil
}

func (c *RabbitMQConsumer) Start(ctx context.Context) error {
	slog.Info("[RabbitMQ] consumer started")
	<-ctx.Done()
	return nil
}

func (c *RabbitMQConsumer) Close() error { return nil }