package cron

import (
	"context"
	"time"

	"myproject/pkg/logger"
)

type Job interface {
	Name() string
	Run(ctx context.Context) error
}

type entry struct {
	job      Job
	interval time.Duration
}

type Scheduler struct {
	entries []entry
	logger  *logger.Logger
}

func NewScheduler(l *logger.Logger) *Scheduler {
	return &Scheduler{logger: l}
}

func (s *Scheduler) Register(name, interval string, job Job) {
	d, err := time.ParseDuration(interval)
	if err != nil {
		s.logger.Error("invalid interval", "name", name, "interval", interval)
		return
	}
	s.entries = append(s.entries, entry{job: job, interval: d})
	s.logger.Info("job registered", "name", name, "interval", interval)
}

func (s *Scheduler) Start(ctx context.Context) {
	for _, e := range s.entries {
		go s.run(ctx, e)
	}
	<-ctx.Done()
}

func (s *Scheduler) run(ctx context.Context, e entry) {
	ticker := time.NewTicker(e.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.logger.Info("running job", "name", e.job.Name())
			if err := e.job.Run(ctx); err != nil {
				s.logger.Error("job failed", "name", e.job.Name(), "err", err)
			}
		}
	}
}