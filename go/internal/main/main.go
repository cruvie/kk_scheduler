package main

import (
	"time"

	"gitee.com/cruvie/kk_kit/go/kk_server"
	"github.com/cruvie/kk-scheduler/go/internal/g_config"
	"github.com/cruvie/kk-scheduler/go/internal/scheduler"
)

func main() {
	stage := g_config.InitConfig()
	defer g_config.CloseConfig()

	kkServer := kk_server.NewKKServer(10*time.Second, stage)
	kkServer.Add("kk-scheduler", 0, scheduler.NewScheduleServer())
	kkServer.Add("kk-scheduler-grpc", 0, NewGrpcServer(stage))
	kkServer.Add("kk-scheduler-http", 0, NewHttpServer(stage))
	kkServer.Add("kk-scheduler-web", 0, NewWebServer(stage))
	kkServer.ServeAndWait()
}
