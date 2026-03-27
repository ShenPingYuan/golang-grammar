package mq

import "context"

// MessageHandler 消息处理函数
type MessageHandler func(message []byte) error

// Consumer 消息消费者接口
type Consumer interface {
	Subscribe(topic string, handler MessageHandler) error
	Start(ctx context.Context) error
	Close() error
}