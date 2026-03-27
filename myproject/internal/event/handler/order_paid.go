package eventhandler

import (
	"context"
	"fmt"
	"log/slog"

	"myproject/internal/event"
	"myproject/internal/notify"
)

func onOrderPaid(notifier notify.Notifier) event.Handler {
	return func(ctx context.Context, evt event.Event) error {
		e, ok := evt.(event.OrderPaidEvent)
		if !ok {
			return nil
		}
		slog.Info("handling order.paid", "order_id", e.OrderID, "amount", e.Amount)
		body := fmt.Sprintf("Your order %s has been paid. Amount: %.2f", e.OrderID, e.Amount)
		return notifier.Send(ctx, e.UserID, "Order Paid", body)
	}
}