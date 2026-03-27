package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"myproject/internal/config"
	"myproject/internal/mq"
	"myproject/pkg/logger"
)

func main() {
	cfg := config.Load("configs/config.yaml")
	l := logger.New(cfg.Log.Level)

	consumer := mq.NewMemoryConsumer()

	_ = consumer.Subscribe("user.created", func(msg []byte) error {
		l.Info("worker received user.created", "payload", string(msg))
		return nil
	})
	_ = consumer.Subscribe("order.paid", func(msg []byte) error {
		l.Info("worker received order.paid", "payload", string(msg))
		return nil
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := consumer.Start(ctx); err != nil {
			slog.Error("consumer error", "err", err)
		}
	}()

	l.Info("Worker started, waiting for messages...")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	l.Info("Worker shutting down...")
	cancel()
	_ = consumer.Close()
	time.Sleep(500 * time.Millisecond)
}