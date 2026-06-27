package interceptor

import (
	"context"
	"net"
	"strings"

	"gitee.com/cruvie/kk_kit/go/kk_ctx"
	"gitee.com/cruvie/kk_kit/go/kk_grpc"
	"github.com/cruvie/kk_scheduler/go/internal/common_go"
	"github.com/cruvie/kk_scheduler/go/kk_scheduler"
	middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

// contextKey 是未导出的类型，用于防止 key 冲突
type contextKey string

// UnaryInit 第一个中间件，用于
// 记录请求所涉及的全部拦截器类型
// 记录客户端ip
func UnaryInit() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		newCtx, err := getNewCtx(ctx, info.FullMethod)
		if err != nil {
			return nil, err
		}
		return handler(newCtx, req)
	}
}

func StreamInit() grpc.StreamServerInterceptor {
	return func(srv any, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		newCtx, err := getNewCtx(stream.Context(), info.FullMethod)
		if err != nil {
			return err
		}
		wrapped := middleware.WrapServerStream(stream)
		wrapped.WrappedContext = newCtx
		return handler(srv, wrapped)
	}
}

func getNewCtx(ctx context.Context, fullMethod string) (context.Context, error) {
	{
		methodDesc := findMethodDescriptor(fullMethod)
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

// findMethodDescriptor 通过 protoregistry.GlobalFiles 查找方法描述符，O(1) 哈希查找。
func findMethodDescriptor(fullMethod string) protoreflect.MethodDescriptor {
	// "/package.Service/Method" → "package.Service.Method"
	name := strings.TrimPrefix(fullMethod, "/")
	name = strings.ReplaceAll(name, "/", ".")
	desc, err := protoregistry.GlobalFiles.FindDescriptorByName(protoreflect.FullName(name))
	if err != nil {
		return nil
	}
	methodDesc, _ := desc.(protoreflect.MethodDescriptor)
	return methodDesc
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
