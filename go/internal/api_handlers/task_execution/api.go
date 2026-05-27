package task_execution

import (
	"gitee.com/cruvie/kk_kit/go/kk_grpc"
	"github.com/cruvie/kk-scheduler/go/kk_scheduler"
)

type ApiTaskCreate struct {
	*kk_grpc.DefaultApi[kk_scheduler.TaskCreate_Input]
}

func NewApiTaskCreate() *ApiTaskCreate {
	return &ApiTaskCreate{
		DefaultApi: kk_grpc.NewDefaultApi[kk_scheduler.TaskCreate_Input](),
	}
}

type ApiTaskUpdateStatus struct {
	*kk_grpc.DefaultApi[kk_scheduler.TaskUpdateStatus_Input]
}

func NewApiTaskUpdateStatus() *ApiTaskUpdateStatus {
	return &ApiTaskUpdateStatus{
		DefaultApi: kk_grpc.NewDefaultApi[kk_scheduler.TaskUpdateStatus_Input](),
	}
}

type ApiTaskAppendLog struct {
	*kk_grpc.DefaultApi[kk_scheduler.TaskAppendLog_Input]
}

func NewApiTaskAppendLog() *ApiTaskAppendLog {
	return &ApiTaskAppendLog{
		DefaultApi: kk_grpc.NewDefaultApi[kk_scheduler.TaskAppendLog_Input](),
	}
}

type ApiTaskExecutionList struct {
	*kk_grpc.DefaultApi[kk_scheduler.TaskExecutionList_Input]
}

func NewApiTaskExecutionList() *ApiTaskExecutionList {
	return &ApiTaskExecutionList{
		DefaultApi: kk_grpc.NewDefaultApi[kk_scheduler.TaskExecutionList_Input](),
	}
}

type ApiTaskExecutionGet struct {
	*kk_grpc.DefaultApi[kk_scheduler.TaskExecutionGet_Input]
}

func NewApiTaskExecutionGet() *ApiTaskExecutionGet {
	return &ApiTaskExecutionGet{
		DefaultApi: kk_grpc.NewDefaultApi[kk_scheduler.TaskExecutionGet_Input](),
	}
}

type ApiTaskExecutionDelete struct {
	*kk_grpc.DefaultApi[kk_scheduler.TaskExecutionDelete_Input]
}

func NewApiTaskExecutionDelete() *ApiTaskExecutionDelete {
	return &ApiTaskExecutionDelete{
		DefaultApi: kk_grpc.NewDefaultApi[kk_scheduler.TaskExecutionDelete_Input](),
	}
}
