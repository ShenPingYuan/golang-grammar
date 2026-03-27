package notify

import (
	"context"
	"log/slog"
)

// EmailNotifier 邮件通知
// 生产环境请集成 net/smtp 或第三方邮件服务 (SendGrid, SES)
type EmailNotifier struct {
	Host     string
	Port     int
	Username string
	Password string
}

func NewEmail(host string, port int, username, password string) Notifier {
	return &EmailNotifier{Host: host, Port: port, Username: username, Password: password}
}

func (n *EmailNotifier) Send(_ context.Context, to, subject, body string) error {
	slog.Info("[Email] send", "to", to, "subject", subject, "body_len", len(body))
	// TODO: 使用 net/smtp.SendMail 发送
	return nil
}