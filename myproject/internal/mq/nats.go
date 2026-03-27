package mq

import (
	"context"
	"log/slog"
)

// NATSProducer NATS 生产者
// 生产环境请使用 github.com/nats-io/nats.go
type NATSProducer struct {
	url string
}

func NewNATSProducer(url string) Producer {
	return &NATSProducer{url: url}
}

func (p *NATSProducer) Publish(_ context.Context, topic string, message []byte) error {
	slog.Info("[NATS] publish", "subject", topic, "size", len(message))
	return nil
}

func (p *NATSProducer) Close() error { return nil }

type NATSConsumer struct {
	url      string
	handlers map[string]MessageHandler
}

func NewNATSConsumer(url string) Consumer {
	return &NATSConsumer{url: url, handlers: make(map[string]MessageHandler)}
}

func (c *NATSConsumer) Subscribe(topic string, handler MessageHandler) error {
	c.handlers[topic] = handler
	return nil
}

func (c *NATSConsumer) Start(ctx context.Context) error {
	slog.Info("[NATS] consumer started")
	<-ctx.Done()
	return nil
}

func (c *NATSConsumer) Close() error { return nil }