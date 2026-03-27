package event

import "context"

// Event 事件接口
type Event interface {
	EventName() string
}

// Handler 事件处理函数
type Handler func(ctx context.Context, evt Event) error

// --- 具体事件类型 ---

type UserCreatedEvent struct {
	UserID   string
	Username string
	Email    string
}

func (e UserCreatedEvent) EventName() string { return "user.created" }

type OrderPaidEvent struct {
	OrderID string
	UserID  string
	Amount  float64
}

func (e OrderPaidEvent) EventName() string { return "order.paid" }