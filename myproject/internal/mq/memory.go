// 文件: internal/mq/memory.go
package mq

import (
	"context"
	"log/slog"
	"sync"
)

// --- 内存 Producer ---

type MemoryProducer struct {
	mu   sync.RWMutex
	subs map[string][]MessageHandler
}

func NewMemoryProducer() *MemoryProducer {
	return &MemoryProducer{subs: make(map[string][]MessageHandler)}
}

func (p *MemoryProducer) Publish(_ context.Context, topic string, message []byte) error {
	p.mu.RLock()
	defer p.mu.RUnlock()
	for _, h := range p.subs[topic] {
		if err := h(message); err != nil {
			slog.Error("[MemoryMQ] handler error", "topic", topic, "err", err)
		}
	}
	return nil
}

func (p *MemoryProducer) Close() error { return nil }

func (p *MemoryProducer) AddSubscriber(topic string, h MessageHandler) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.subs[topic] = append(p.subs[topic], h)
}

// --- 内存 Consumer ---

type MemoryConsumer struct {
	handlers map[string]MessageHandler
}

func NewMemoryConsumer() Consumer {
	return &MemoryConsumer{handlers: make(map[string]MessageHandler)}
}

func (c *MemoryConsumer) Subscribe(topic string, handler MessageHandler) error {
	c.handlers[topic] = handler
	return nil
}

func (c *MemoryConsumer) Start(ctx context.Context) error {
	slog.Info("[MemoryMQ] consumer started, waiting...")
	<-ctx.Done()
	return nil
}

func (c *MemoryConsumer) Close() error { return nil }