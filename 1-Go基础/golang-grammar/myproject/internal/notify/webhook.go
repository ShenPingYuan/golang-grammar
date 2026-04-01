package notify

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
)

// WebhookNotifier 通过 Webhook 回调通知
type WebhookNotifier struct {
	URL string
}

func NewWebhook(url string) Notifier {
	return &WebhookNotifier{URL: url}
}

func (n *WebhookNotifier) Send(ctx context.Context, to, subject, body string) error {
	payload, _ := json.Marshal(map[string]string{
		"to": to, "subject": subject, "body": body,
	})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, n.URL, bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	slog.Info("[Webhook] sent", "url", n.URL, "status", resp.StatusCode)
	return nil
}

// ConsoleNotifier 控制台输出（开发用）
type ConsoleNotifier struct{}

func NewConsole() Notifier {
	return &ConsoleNotifier{}
}

func (n *ConsoleNotifier) Send(_ context.Context, to, subject, body string) error {
	slog.Info("[Notify] console", "to", to, "subject", subject, "body", body)
	return nil
}