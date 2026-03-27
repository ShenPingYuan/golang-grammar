package jobs

import (
	"context"

	"myproject/pkg/logger"
)

type ReportJob struct {
	logger *logger.Logger
}

func NewReportJob(l *logger.Logger) *ReportJob {
	return &ReportJob{logger: l}
}

func (j *ReportJob) Name() string { return "report" }

func (j *ReportJob) Run(_ context.Context) error {
	j.logger.Info("report: generating daily summary")
	// TODO: 统计数据并生成报表 / 发送邮件
	return nil
}