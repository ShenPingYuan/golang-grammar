package notify

import (
	"context"
	"log/slog"
)

// SMSNotifier 短信通知
type SMSNotifier struct {
	APIKey string
}

func NewSMS(apiKey string) Notifier {
	return &SMSNotifier{APIKey: apiKey}
}

func (n *SMSNotifier) Send(_ context.Context, to, subject, body string) error {
	slog.Info("[SMS] send", "to", to, "content", body)
	// TODO: 调用短信服务 API (阿里云、腾讯云等)
	return nil
}