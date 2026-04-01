package notify

import "context"

// Notifier 通知推送抽象接口
type Notifier interface {
	Send(ctx context.Context, to, subject, body string) error
}