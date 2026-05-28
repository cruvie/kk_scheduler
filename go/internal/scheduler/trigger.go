package scheduler

import (
	"context"
	"log/slog"

	"gitee.com/cruvie/kk_kit/go/kk_grpc"
	"gitee.com/cruvie/kk_kit/go/kk_stage"
	"github.com/cruvie/kk_scheduler/go/kk_scheduler"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func triggerClient(service *kk_scheduler.PBRegisterService) (conn *grpc.ClientConn, client kk_scheduler.KKScheduleTriggerClient, err error) {
	var opts []grpc.DialOption
	if service.GetAuthToken() != "" {
		opts = append(opts, grpc.WithAuthority(service.GetAuthToken()))
	}
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err = grpc.NewClient(service.GetTarget(), opts...)
	if err != nil {
		return nil, nil, err
	}
	return conn, kk_scheduler.NewKKScheduleTriggerClient(conn), nil
}

func triggerFunc(service *kk_scheduler.PBRegisterService, pbJob *kk_scheduler.PBJob) func() {
	return func() {
		conn, client, err := triggerClient(service)
		if err != nil {
			slog.Error(err.Error())
			return
		}
		defer func() {
			err := conn.Close()
			if err != nil {
				slog.Error(err.Error())
			}
		}()
		stage := kk_stage.NewStage(context.Background(), "kk_scheduler")
		ctx, cancelFunc := kk_grpc.NewCallGrpcCtx(stage)
		defer cancelFunc()

		input := &kk_scheduler.Trigger_Input{}
		input.SetFuncName(pbJob.GetFuncName())
		input.SetJobId(pbJob.GetId())
		_, err = client.Trigger(ctx, input)
		if err != nil {
			slog.Error(err.Error())
		}
	}
}
