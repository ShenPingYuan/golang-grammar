package eventhandler

import (
	"myproject/internal/event"
	"myproject/internal/notify"
	"myproject/pkg/logger"
)

// Register 注册所有事件处理器到总线
func Register(bus event.Bus, notifier notify.Notifier, l *logger.Logger) {
	bus.Subscribe("user.created", onUserCreated(notifier))
	bus.Subscribe("order.paid", onOrderPaid(notifier))

	l.Info("event handlers registered")
}