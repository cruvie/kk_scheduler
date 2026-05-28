package task_execution

import (
	"gitee.com/cruvie/kk_kit/go/kk_stage"
	"github.com/cruvie/kk_scheduler/go/kk_scheduler"
)

func (x *ApiTaskCreate) Handler(stage *kk_stage.Stage) (*kk_scheduler.TaskCreate_Output, error) {
	err := x.Service(stage)
	if err != nil {
		return nil, err
	}
	return &kk_scheduler.TaskCreate_Output{}, nil
}

func (x *ApiTaskUpdateStatus) Handler(stage *kk_stage.Stage) (*kk_scheduler.TaskUpdateStatus_Output, error) {
	err := x.Service(stage)
	if err != nil {
		return nil, err
	}
	return &kk_scheduler.TaskUpdateStatus_Output{}, nil
}

func (x *ApiTaskAppendLog) Handler(stage *kk_stage.Stage) (*kk_scheduler.TaskAppendLog_Output, error) {
	err := x.Service(stage)
	if err != nil {
		return nil, err
	}
	return &kk_scheduler.TaskAppendLog_Output{}, nil
}

func (x *ApiTaskExecutionList) Handler(stage *kk_stage.Stage) (*kk_scheduler.TaskExecutionList_Output, error) {
	list, err := x.Service(stage)
	if err != nil {
		return nil, err
	}
	out := &kk_scheduler.TaskExecutionList_Output{}
	out.SetTaskExecutionList(list)
	return out, nil
}

func (x *ApiTaskExecutionGet) Handler(stage *kk_stage.Stage) (*kk_scheduler.TaskExecutionGet_Output, error) {
	execution, err := x.Service(stage)
	if err != nil {
		return nil, err
	}
	out := &kk_scheduler.TaskExecutionGet_Output{}
	out.SetTaskExecution(execution)
	return out, nil
}

func (x *ApiTaskExecutionDelete) Handler(stage *kk_stage.Stage) (*kk_scheduler.TaskExecutionDelete_Output, error) {
	err := x.Service(stage)
	if err != nil {
		return nil, err
	}
	return &kk_scheduler.TaskExecutionDelete_Output{}, nil
}
