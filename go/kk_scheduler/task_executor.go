package kk_scheduler

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"gitee.com/cruvie/kk_go_kit/kk_id"
	"gitee.com/cruvie/kk_go_kit/kk_stage"
	"github.com/samber/lo"
)

type StepCtl struct {
	stage     *kk_stage.Stage
	reportLog func(errOr error, msgOr string)
	hasError  bool
}

// Log adds a log message to the execution record
func (c *StepCtl) Log(errOr error, msgOr string) {
	if errOr != nil {
		c.reportLog(errOr, msgOr)
		slog.Error("TaskMsg", kk_stage.NewLog(c.stage).Error(errOr).Args()...)
	} else {
		c.reportLog(nil, msgOr)
		slog.Info("TaskMsg", kk_stage.NewLog(c.stage).String("msg", msgOr).Args()...)
	}
}

// Step represents a single execution step
type Step struct {
	name     string
	handler  func(ctl *StepCtl) error
	fallback func(ctl *StepCtl) error
}

// TaskExecutor manages task lifecycle and step execution
type TaskExecutor struct {
	id     string
	jobId  string
	steps  []*Step
	client KKScheduleClient
	stage  *kk_stage.Stage
}

// NewTaskExecutor creates a new task executor
func NewTaskExecutor(opts ...TaskExecutorOption) *TaskExecutor {
	t := &TaskExecutor{
		id: kk_id.GenUUID7(),
	}

	for _, opt := range opts {
		opt(t)
	}

	return t
}

// AddStep adds a step to the task
func (t *TaskExecutor) AddStep(
	name string,
	handler func(ctl *StepCtl) error,
	fallback func(ctl *StepCtl) error,
) *TaskExecutor {
	t.steps = append(t.steps, &Step{
		name:     name,
		handler:  handler,
		fallback: fallback,
	})
	return t
}

// Run executes all steps sequentially
func (t *TaskExecutor) Run(ctx context.Context) error {
	{
		if t.client == nil {
			return fmt.Errorf("scheduler client is not set")
		}
		if t.jobId == "" {
			return fmt.Errorf("job id is not set")
		}
	}
	{
		// Create execution record
		input := &TaskCreate_Input{}
		input.SetId(t.id)
		input.SetJobId(t.jobId)
		_, err := t.client.TaskCreate(ctx, input)
		if err != nil {
			return fmt.Errorf("failed to create execution record: %w", err)
		}
	}

	// Create shared StepCtl for all steps
	ctl := &StepCtl{
		stage: t.stage,
		reportLog: func(errOr error, msgOr string) {
			in := &TaskAppendLog_Input{}
			in.SetId(t.id)
			in.SetLog(lo.TernaryF(errOr != nil,
				func() string {
					return errOr.Error()
				}, func() string {
					return msgOr
				}))
			_, err := t.client.TaskAppendLog(ctx, in)
			if err != nil {
				slog.Error("failed to report log", "err", err)
			}
		},
	}
	ctl.Log(nil, "[🚀Task Start]")
	// Execute steps
	for _, step := range t.steps {
		ctl.Log(nil, fmt.Sprintf("[🧩Step Begin %s]", step.name))
		err := step.handler(ctl)
		if err != nil {
			if errors.Is(err, ErrStopTask) {
				ctl.Log(nil, "⚠️stop step manually")
			} else {
				ctl.Log(fmt.Errorf("❌%w", err), "")
				ctl.hasError = true
				if step.fallback != nil {
					ctl.Log(nil, fmt.Sprintf("[🧯Step Fallback %s]", step.name))
					err := step.fallback(ctl)
					if err != nil {
						ctl.Log(fmt.Errorf("❌%w", err), "")
					}
				}
			}
			break
		}
	}

	in := &TaskUpdateStatus_Input{}
	in.SetId(t.id)
	if ctl.hasError {
		ctl.Log(nil, "[😭Task Failed]")
		in.SetStatus(TaskExecutionStatus_TASK_EXECUTION_STATUS_FAILED)
	} else {
		ctl.Log(nil, "[✅Task Finished]")
		in.SetStatus(TaskExecutionStatus_TASK_EXECUTION_STATUS_COMPLETED)
	}
	_, err := t.client.TaskUpdateStatus(ctx, in)
	if err != nil {
		ctl.Log(fmt.Errorf("failed to update task status %w", err), "")
	}

	return nil
}
