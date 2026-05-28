package api_impl

import (
	"context"

	"gitee.com/cruvie/kk_kit/go/kk_grpc"
	"github.com/cruvie/kk_scheduler/go/internal/api_handlers/job"
	"github.com/cruvie/kk_scheduler/go/internal/api_handlers/service"
	"github.com/cruvie/kk_scheduler/go/internal/api_handlers/task_execution"
	"github.com/cruvie/kk_scheduler/go/kk_scheduler"
)

func (x *server) JobList(ctx context.Context, input *kk_scheduler.JobList_Input) (*kk_scheduler.JobList_Output, error) {
	return kk_grpc.GrpcHandler(
		ctx,
		input,
		job.NewApiJobList,
	)
}

func (x *server) JobGet(ctx context.Context, input *kk_scheduler.JobGet_Input) (*kk_scheduler.JobGet_Output, error) {
	return kk_grpc.GrpcHandler(
		ctx,
		input,
		job.NewApiJobGet,
	)
}

func (x *server) JobSetSpec(ctx context.Context, input *kk_scheduler.JobSetSpec_Input) (*kk_scheduler.JobSetSpec_Output, error) {
	return kk_grpc.GrpcHandler(
		ctx,
		input,
		job.NewApiJobSetSpec,
	)
}

func (x *server) JobEnable(ctx context.Context, input *kk_scheduler.JobEnable_Input) (*kk_scheduler.JobEnable_Output, error) {
	return kk_grpc.GrpcHandler(
		ctx,
		input,
		job.NewApiJobEnable,
	)
}

func (x *server) JobDisable(ctx context.Context, input *kk_scheduler.JobDisable_Input) (*kk_scheduler.JobDisable_Output, error) {
	return kk_grpc.GrpcHandler(
		ctx,
		input,
		job.NewApiJobDisable,
	)
}

func (x *server) JobPut(ctx context.Context, input *kk_scheduler.JobPut_Input) (*kk_scheduler.JobPut_Output, error) {
	return kk_grpc.GrpcHandler(
		ctx,
		input,
		job.NewApiJobPut,
	)
}

func (x *server) JobDelete(ctx context.Context, input *kk_scheduler.JobDelete_Input) (*kk_scheduler.JobDelete_Output, error) {
	return kk_grpc.GrpcHandler(
		ctx,
		input,
		job.NewApiJobDelete,
	)
}

func (x *server) JobTrigger(ctx context.Context, input *kk_scheduler.JobTrigger_Input) (*kk_scheduler.JobTrigger_Output, error) {
	return kk_grpc.GrpcHandler(
		ctx,
		input,
		job.NewApiJobTrigger,
	)
}

func (x *server) ServiceList(ctx context.Context, input *kk_scheduler.ServiceList_Input) (*kk_scheduler.ServiceList_Output, error) {
	return kk_grpc.GrpcHandler(
		ctx,
		input,
		service.NewApiServiceList,
	)
}

func (x *server) ServicePut(ctx context.Context, input *kk_scheduler.ServicePut_Input) (*kk_scheduler.ServicePut_Output, error) {
	return kk_grpc.GrpcHandler(
		ctx,
		input,
		service.NewApiServicePut,
	)
}

func (x *server) ServiceGet(ctx context.Context, input *kk_scheduler.ServiceGet_Input) (*kk_scheduler.ServiceGet_Output, error) {
	return kk_grpc.GrpcHandler(
		ctx,
		input,
		service.NewApiServiceGet,
	)
}

func (x *server) ServiceDelete(ctx context.Context, input *kk_scheduler.ServiceDelete_Input) (*kk_scheduler.ServiceDelete_Output, error) {
	return kk_grpc.GrpcHandler(
		ctx,
		input,
		service.NewApiServiceDelete,
	)
}

func (x *server) TaskCreate(ctx context.Context, input *kk_scheduler.TaskCreate_Input) (*kk_scheduler.TaskCreate_Output, error) {
	return kk_grpc.GrpcHandler(
		ctx,
		input,
		task_execution.NewApiTaskCreate,
	)
}

func (x *server) TaskUpdateStatus(ctx context.Context, input *kk_scheduler.TaskUpdateStatus_Input) (*kk_scheduler.TaskUpdateStatus_Output, error) {
	return kk_grpc.GrpcHandler(
		ctx,
		input,
		task_execution.NewApiTaskUpdateStatus,
	)
}

func (x *server) TaskAppendLog(ctx context.Context, input *kk_scheduler.TaskAppendLog_Input) (*kk_scheduler.TaskAppendLog_Output, error) {
	return kk_grpc.GrpcHandler(
		ctx,
		input,
		task_execution.NewApiTaskAppendLog,
	)
}

func (x *server) TaskExecutionList(ctx context.Context, input *kk_scheduler.TaskExecutionList_Input) (*kk_scheduler.TaskExecutionList_Output, error) {
	return kk_grpc.GrpcHandler(
		ctx,
		input,
		task_execution.NewApiTaskExecutionList,
	)
}

func (x *server) TaskExecutionGet(ctx context.Context, input *kk_scheduler.TaskExecutionGet_Input) (*kk_scheduler.TaskExecutionGet_Output, error) {
	return kk_grpc.GrpcHandler(
		ctx,
		input,
		task_execution.NewApiTaskExecutionGet,
	)
}

func (x *server) TaskExecutionDelete(ctx context.Context, input *kk_scheduler.TaskExecutionDelete_Input) (*kk_scheduler.TaskExecutionDelete_Output, error) {
	return kk_grpc.GrpcHandler(
		ctx,
		input,
		task_execution.NewApiTaskExecutionDelete,
	)
}
