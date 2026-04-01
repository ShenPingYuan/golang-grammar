package jobs

import (
	"context"

	"myproject/pkg/logger"
)

type CleanupJob struct {
	logger *logger.Logger
}

func NewCleanupJob(l *logger.Logger) *CleanupJob {
	return &CleanupJob{logger: l}
}

func (j *CleanupJob) Name() string { return "cleanup" }

func (j *CleanupJob) Run(_ context.Context) error {
	j.logger.Info("cleanup: removing expired sessions and temp data")
	// TODO: 清理过期 session、临时文件、软删除数据等
	return nil
}