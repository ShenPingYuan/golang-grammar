package eventhandler

import (
	"context"
	"log/slog"

	"myproject/internal/event"
	"myproject/internal/notify"
)

func onUserCreated(notifier notify.Notifier) event.Handler {
	return func(ctx context.Context, evt event.Event) error {
		e, ok := evt.(event.UserCreatedEvent)
		if !ok {
			return nil
		}
		slog.Info("handling user.created", "user_id", e.UserID, "email", e.Email)
		return notifier.Send(ctx, e.Email, "Welcome!", "Hello "+e.Username+", welcome to our platform!")
	}
}