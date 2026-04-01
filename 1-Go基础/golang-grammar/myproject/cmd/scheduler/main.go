package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"myproject/internal/config"
	"myproject/internal/cron"
	"myproject/internal/cron/jobs"
	"myproject/pkg/logger"
)

func main() {
	cfg := config.Load("configs/config.yaml")
	l := logger.New(cfg.Log.Level)

	s := cron.NewScheduler(l)
	s.Register("cleanup", "1h", jobs.NewCleanupJob(l))
	s.Register("report", "24h", jobs.NewReportJob(l))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go s.Start(ctx)

	l.Info("Scheduler started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	l.Info("Scheduler shutting down...")
	cancel()
}