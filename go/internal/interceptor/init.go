package interceptor

import (
	"context"
	"net"

	"gitee.com/cruvie/kk_go_kit/kk_ctx"
	"gitee.com/cruvie/kk_go_kit/kk_grpc"
	"github.com/cruvie/kk-scheduler/go/internal/common_go"
	"github.com/cruvie/kk-scheduler/go/kk_scheduler"
	middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

// contextKey 是未导出的类型，用于防止 key 冲突
type contextKey string

// UnaryInit 第一个中间件，用于
// 记录请求所涉及的全部拦截器类型
// 记录客户端ip
func UnaryInit(fileDescHub *kk_grpc.FileDescHub) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		newCtx, err := getNewCtx(ctx, info.FullMethod, fileDescHub)
		if err != nil {
			return nil, err
		}
		return handler(newCtx, req)
	}
}

func StreamInit(fileDescHub *kk_grpc.FileDescHub) grpc.StreamServerInterceptor {
	return func(srv any, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		newCtx, err := getNewCtx(stream.Context(), info.FullMethod, fileDescHub)
		if err != nil {
			return err
		}
		wrapped := middleware.WrapServerStream(stream)
		wrapped.WrappedContext = newCtx
		return handler(srv, wrapped)
	}
}

func getNewCtx(ctx context.Context, fullMethod string, fileDescHub *kk_grpc.FileDescHub) (context.Context, error) {
	{
		methodDesc := fileDescHub.GetMethodDescriptor(fullMethod)
		if methodDesc == nil {
			return nil, status.Error(codes.Unavailable, "kk_grpc method not found")
		}

		interceptorAuth := common_go.MethodDescGetInterceptorAuth(methodDesc)
		ctx = context.WithValue(ctx, interceptorAuthKey, interceptorAuth)
	}
	{
		ctx = context.WithValue(ctx, kk_grpc.FullMethodKey, fullMethod)
	}
	{
		// todo 如果使用了代理需要获取真实 ip
		if p, ok := peer.FromContext(ctx); ok {
			if tcpAddr, ok := p.Addr.(*net.TCPAddr); ok {
				ctx = context.WithValue(ctx, realIPKey, tcpAddr)
			}
		}
	}

	return ctx, nil
}

const (
	interceptorAuthKey contextKey = "interceptorAuthKey"
	realIPKey          contextKey = "realIPKey"
)

func getInterceptorAuth(ctx context.Context) (kk_scheduler.InterceptorAuth, error) {
	interceptorAuth, ok := kk_ctx.Value[kk_scheduler.InterceptorAuth](ctx, interceptorAuthKey)
	if !ok {
		return kk_scheduler.InterceptorAuth_UNSPECIFIED, status.Error(codes.NotFound, "kk_grpc interceptorAuth type not found")
	}
	return interceptorAuth, nil
}

func GetRealIP(ctx context.Context) (realIP *net.TCPAddr, err error) {
	realIP, ok := kk_ctx.Value[*net.TCPAddr](ctx, realIPKey)
	if !ok {
		return nil, status.Error(codes.NotFound, "real ip not found")
	}
	return realIP, nil
}
