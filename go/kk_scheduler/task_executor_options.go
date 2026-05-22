package kk_scheduler

import "gitee.com/cruvie/kk_go_kit/kk_stage"

// TaskExecutorOption configures TaskExecutor
type TaskExecutorOption func(*TaskExecutor)

// WithSchedulerClient sets the gRPC client for reporting
func WithSchedulerClient(client KKScheduleClient) TaskExecutorOption {
	return func(t *TaskExecutor) {
		t.client = client
	}
}

func WithJobId(id string) TaskExecutorOption {
	return func(t *TaskExecutor) {
		t.jobId = id
	}
}

func WithStage(stage *kk_stage.Stage) TaskExecutorOption {
	return func(t *TaskExecutor) {
		t.stage = stage
	}
}
