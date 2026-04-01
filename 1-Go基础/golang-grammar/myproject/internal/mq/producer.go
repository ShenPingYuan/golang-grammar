package mq

import "context"

// Producer 消息生产者接口
type Producer interface {
	Publish(ctx context.Context, topic string, message []byte) error
	Close() error
}