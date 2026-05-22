package interceptor

import (
	"context"
	"log/slog"

	"gitee.com/cruvie/kk_go_kit/kk_stage"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"google.golang.org/grpc"
)

// interceptorLogger adapts slog logger to interceptor logger.
// This code is simple enough to be copied and not imported.
func interceptorLogger(configSlog *kk_stage.ConfigLog) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		configSlog.Logger.Log(ctx, slog.Level(lvl), msg, fields...)
	})
}

func UnaryLogging(configSlog *kk_stage.ConfigLog) grpc.UnaryServerInterceptor {
	return logging.UnaryServerInterceptor(interceptorLogger(configSlog))
}
