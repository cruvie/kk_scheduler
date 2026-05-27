package internal

import (
	"testing"

	"gitee.com/cruvie/kk_kit/go/kk_grpc/grpc_api_gen"
	"gitee.com/cruvie/kk_kit/go/kk_system"
	"github.com/cruvie/kk-scheduler/go/kk_scheduler"
)

func TestName(t *testing.T) {
	kk_system.TerminatePort(3000)
}

func TestGeneratePermissionApi(t *testing.T) {
	//apiGroupModel := grpc_api_gen.ApiGroupModel{
	//	AdditionImports: `
	//`,
	//	TargetPath: "./api_handlers",
	//}

	//grpc_api_gen.GenerateHandler(apiGroupModel, grpc_api_gen.ApiModel{
	//	ApiPtr: &kk_scheduler.JobList{},
	//})
	//
	//grpc_api_gen.GenerateHandler(apiGroupModel, grpc_api_gen.ApiModel{
	//	ApiPtr: &kk_scheduler.JobGet{},
	//})
	//
	//grpc_api_gen.GenerateHandler(apiGroupModel, grpc_api_gen.ApiModel{
	//	ApiPtr: &kk_scheduler.JobSetSpec{},
	//})
	//grpc_api_gen.GenerateHandler(apiGroupModel, grpc_api_gen.ApiModel{
	//	ApiPtr: &kk_scheduler.JobEnable{},
	//})
	//
	//grpc_api_gen.GenerateHandler(apiGroupModel, grpc_api_gen.ApiModel{
	//	ApiPtr: &kk_scheduler.JobDisable{},
	//})
	//
	//grpc_api_gen.GenerateHandler(apiGroupModel, grpc_api_gen.ApiModel{
	//	ApiPtr: &kk_scheduler.JobPut{},
	//})
	//grpc_api_gen.GenerateHandler(apiGroupModel, grpc_api_gen.ApiModel{
	//	ApiPtr: &kk_scheduler.JobDelete{},
	//})
	//grpc_api_gen.GenerateHandler(apiGroupModel, grpc_api_gen.ApiModel{
	//	ApiPtr: &kk_scheduler.JobTrigger{},
	//})
	//grpc_api_gen.GenerateHandler(apiGroupModel, grpc_api_gen.ApiModel{
	//	ApiPtr: &kk_scheduler.ServiceList{},
	//})
	//grpc_api_gen.GenerateHandler(apiGroupModel, grpc_api_gen.ApiModel{
	//	ApiPtr: &kk_scheduler.ServicePut{},
	//})
	//grpc_api_gen.GenerateHandler(apiGroupModel, grpc_api_gen.ApiModel{
	//	ApiPtr: &kk_scheduler.ServiceGet{},
	//})
	//grpc_api_gen.GenerateHandler(apiGroupModel, grpc_api_gen.ApiModel{
	//	ApiPtr: &kk_scheduler.ServiceDelete{},
	//})
}

func TestGenImpl(t *testing.T) {
	grpc_api_gen.GenerateImpl(
		grpc_api_gen.GenerateImplInput{
			ServerName:      "KKSchedule",
			Methods:         kk_scheduler.KKSchedule_ServiceDesc.Methods,
			ApiDefPkgPath:   "github.com/cruvie/kk-scheduler/go/internal/api_def",
			HandlersPkgPath: "github.com/cruvie/kk-scheduler/go/internal/api_handlers",
		},
	)
}
