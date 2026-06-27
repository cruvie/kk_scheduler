package api_impl

import (
	"github.com/cruvie/kk_scheduler/go/kk_scheduler"
	"google.golang.org/grpc"
)

type server struct {
	kk_scheduler.UnimplementedKKScheduleServer
}

func RegisterServer(grpcServer *grpc.Server) {
	kk_scheduler.RegisterKKScheduleServer(grpcServer, &server{})
}
