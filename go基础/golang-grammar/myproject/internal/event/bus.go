package event

import (
	"context"
	"log/slog"
	"sync"
)

// Bus 事件总线接口
type Bus interface {
	Publish(ctx context.Context, evt Event) error
	Subscribe(eventName string, h Handler)
}

// --- 内存实现 ---

type memoryBus struct {
	mu       sync.RWMutex
	handlers map[string][]Handler
}

func NewBus() Bus {
	return &memoryBus{handlers: make(map[string][]Handler)}
}

func (b *memoryBus) Subscribe(eventName string, h Handler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[eventName] = append(b.handlers[eventName], h)
}

func (b *memoryBus) Publish(ctx context.Context, evt Event) error {
	b.mu.RLock()
	handlers := make([]Handler, len(b.handlers[evt.EventName()]))
	copy(handlers, b.handlers[evt.EventName()])
	b.mu.RUnlock()

	for _, h := range handlers {
		if err := h(ctx, evt); err != nil {
			slog.Error("event handler error",
				"event", evt.EventName(),
				"err", err,
			)
		}
	}
	return nil
}