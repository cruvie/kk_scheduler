package api_impl

import (
	"gitee.com/cruvie/kk_kit/go/kk_grpc"
	"github.com/cruvie/kk_scheduler/go/kk_scheduler"
	"google.golang.org/grpc"
)

type server struct {
	kk_scheduler.UnimplementedKKScheduleServer
}

func RegisterServer(grpcServer *grpc.Server) {
	kk_scheduler.RegisterKKScheduleServer(grpcServer, &server{})
}

func init() {
	kk_grpc.GFileDescHub.RegisterFileDesc(kk_scheduler.File_kk_scheduler_rpc_service_proto)
}
