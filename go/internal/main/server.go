package main

import (
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"

	"gitee.com/cruvie/kk_kit/go/kk_grpc"
	"gitee.com/cruvie/kk_kit/go/kk_server"
	"gitee.com/cruvie/kk_kit/go/kk_stage"
	"github.com/cruvie/kk_scheduler/go/internal/api_impl"
	"github.com/cruvie/kk_scheduler/go/internal/common_go"
	"github.com/cruvie/kk_scheduler/go/internal/g_config"
	"github.com/cruvie/kk_scheduler/go/internal/interceptor"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

func getGrpcServer(stage *kk_stage.Stage) *grpc.Server {
	grpcServer := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),

		grpc.ChainUnaryInterceptor(
			interceptor.UnaryInit(kk_grpc.GFileDescHub),
			interceptor.UnaryLogging(g_config.Config.ConfigSlog),
			recovery.UnaryServerInterceptor(recovery.WithRecoveryHandler(interceptor.PanicRecovery)),
		),
	)
	{
		kk_grpc.SetCheckFieldsFn(common_go.CheckFields)
		kk_grpc.RegisterReflectionServer(stage, g_config.Config.GrpcPort, grpcServer)
		kk_grpc.RegisterKKHealthCheckServer(grpcServer)
		api_impl.RegisterServer(grpcServer)
	}

	return grpcServer
}

func NewGrpcServer(stage *kk_stage.Stage) *kk_server.KKRunServer {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", g_config.Config.GrpcPort))
	if err != nil {
		panic(err)
	}
	grpcServer := getGrpcServer(stage)

	run := func() {
		if err := grpcServer.Serve(listener); err != nil {
			panic(err)
		}
	}
	done := func(quitCh <-chan struct{}) {
		<-quitCh
		grpcServer.GracefulStop()
	}

	return &kk_server.KKRunServer{
		Run:  run,
		Done: done,
	}
}

func NewHttpServer(stage *kk_stage.Stage) *kk_server.KKRunServer {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", g_config.Config.HttpPort))
	if err != nil {
		panic(err)
	}
	grpcServer := getGrpcServer(stage)
	grpcWebServer := grpcweb.WrapServer(grpcServer)
	httpServer := &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "*")
			w.Header().Set("Access-Control-Allow-Headers", "*")
			if r.Method == http.MethodOptions {
				return
			}
			if grpcWebServer.IsGrpcWebRequest(r) {
				grpcWebServer.ServeHTTP(w, r)
			} else {
				http.Error(w, "Not a valid gRPC-Web request", http.StatusBadRequest)
			}
		}),
	}

	run := func() {
		if err := httpServer.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}
	done := func(quitCh <-chan struct{}) {
		<-quitCh
		grpcServer.GracefulStop()
		err := httpServer.Shutdown(stage.Ctx)
		if err != nil {
			slog.Error("httpServer.Shutdown", "err", err)
		}
	}

	return &kk_server.KKRunServer{
		Run:  run,
		Done: done,
	}
}
