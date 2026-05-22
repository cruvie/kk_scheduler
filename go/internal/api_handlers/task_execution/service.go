package task_execution

import (
	"gitee.com/cruvie/kk_go_kit/kk_stage"
	"github.com/cruvie/kk-scheduler/go/internal/scheduler"
	"github.com/cruvie/kk-scheduler/go/kk_scheduler"
)

func (x *ApiTaskCreate) Service(stage *kk_stage.Stage) error {
	span := stage.StartTrace("Service")
	defer span.End()

	return scheduler.GClient.TaskCreate(x.In)
}

func (x *ApiTaskUpdateStatus) Service(stage *kk_stage.Stage) error {
	span := stage.StartTrace("Service")
	defer span.End()

	return scheduler.GClient.TaskUpdateStatus(x.In.GetId(), x.In.GetStatus())
}

func (x *ApiTaskAppendLog) Service(stage *kk_stage.Stage) error {
	span := stage.StartTrace("Service")
	defer span.End()

	return scheduler.GClient.TaskAppendLog(x.In.GetId(), x.In.GetLog())
}

func (x *ApiTaskExecutionList) Service(stage *kk_stage.Stage) ([]*kk_scheduler.PBTaskExecution, error) {
	span := stage.StartTrace("Service")
	defer span.End()

	return scheduler.GClient.TaskExecutionList(x.In.GetJobId())
}

func (x *ApiTaskExecutionGet) Service(stage *kk_stage.Stage) (*kk_scheduler.PBTaskExecution, error) {
	span := stage.StartTrace("Service")
	defer span.End()

	return scheduler.GClient.TaskExecutionGet(x.In.GetId())
}

func (x *ApiTaskExecutionDelete) Service(stage *kk_stage.Stage) error {
	span := stage.StartTrace("Service")
	defer span.End()

	return scheduler.GClient.TaskExecutionDelete(x.In.GetId())
}
