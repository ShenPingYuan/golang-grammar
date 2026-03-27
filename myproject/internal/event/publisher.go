package event

import (
	"context"
	"encoding/json"
	"log/slog"

	"myproject/internal/mq"
)

// Publisher 组合本地事件总线和远程 MQ 的发布者
type Publisher struct {
	bus      Bus
	producer mq.Producer
}

func NewPublisher(bus Bus, producer mq.Producer) *Publisher {
	return &Publisher{bus: bus, producer: producer}
}

func (p *Publisher) Publish(ctx context.Context, evt Event) error {
	// 先发布到本地总线（同步处理）
	if err := p.bus.Publish(ctx, evt); err != nil {
		slog.Error("local bus publish failed", "event", evt.EventName(), "err", err)
	}
	// 再发布到 MQ（异步处理 / 跨服务）
	if p.producer != nil {
		data, err := json.Marshal(evt)
		if err != nil {
			return err
		}
		return p.producer.Publish(ctx, evt.EventName(), data)
	}
	return nil
}