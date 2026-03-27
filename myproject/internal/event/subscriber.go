package event

import (
	"myproject/internal/mq"
)

// RegisterConsumers 将 MQ 消息桥接到本地事件总线
func RegisterConsumers(bus Bus, consumer mq.Consumer) error {
	topics := []string{"user.created", "order.paid"}
	for _, topic := range topics {
		t := topic
		err := consumer.Subscribe(t, func(msg []byte) error {
			// MQ 消息桥接为本地日志（完整实现需反序列化为 Event）
			return nil
		})
		if err != nil {
			return err
		}
	}
	return nil
}