package schedule_test

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"testing"

	"github.com/cruvie/kk_scheduler/go/kk_scheduler"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	kk_scheduler.UnimplementedKKScheduleTriggerServer
}

func (server) Trigger(ctx context.Context, input *kk_scheduler.Trigger_Input) (*kk_scheduler.Trigger_Output, error) {
	slog.Info("Trigger received", "FuncName", input.GetFuncName())
	switch input.GetFuncName() {
	case "Func1":
		go Func1()
	default:
		return nil, kk_scheduler.ErrJobNotFount
	}
	return &kk_scheduler.Trigger_Output{}, nil
}

func Func1() {
	slog.Info("Func1 start")
	defer slog.Info("Func1 end")
}

func authorityAuthInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	err := kk_scheduler.CheckAuthority(ctx, testAuthToken)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	return handler(ctx, req)
}

func TestClientServer(t *testing.T) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", 8000))
	if err != nil {
		panic(err)
	}
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			authorityAuthInterceptor,
		),
	)
	defer grpcServer.GracefulStop()
	kk_scheduler.RegisterKKScheduleTriggerServer(grpcServer, &server{})
	if err := grpcServer.Serve(listener); err != nil {
		panic(err)
	}
}
