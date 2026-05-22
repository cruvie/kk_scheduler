package interceptor

import (
	"context"
	"strconv"

	"gitee.com/cruvie/kk_go_kit/kk_jwt"
	"github.com/cruvie/kk-scheduler/go/kk_scheduler"
	middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const internalOnlyTokenKey = "InternalOnlyTokenKey"

// SetInternalOnlyToken sets the internal only token in the outgoing context metadata.
func SetInternalOnlyToken(ctx context.Context, token string) context.Context {
	return metadata.AppendToOutgoingContext(ctx, internalOnlyTokenKey, token)
}

// InternalOnlyCheckFunc defines the function signature for internal token verification.
type InternalOnlyCheckFunc func(ctx context.Context, token string) (
	newCtx context.Context,
	err error)

// SetAccessToken sets the access token in the outgoing context metadata.
func SetAccessToken(ctx context.Context, token string) context.Context {
	return metadata.AppendToOutgoingContext(ctx, string(kk_jwt.AccessToken), token)
}

// GetAccessToken retrieves the access token from the incoming context metadata.
func GetAccessToken(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}

	tokens := md.Get(string(kk_jwt.AccessToken))
	if len(tokens) == 0 {
		return ""
	}

	return tokens[0]
}

// JWTCheckFunc defines the function signature for JWT token verification.
type JWTCheckFunc func(ctx context.Context, token string) (
	newCtx context.Context,
	needRefresh bool,
	err error)

// AuthConfig holds the configuration for authentication interceptors.
type AuthConfig struct {
	JWTChecker          JWTCheckFunc
	InternalOnlyChecker InternalOnlyCheckFunc
}

func verifyInternalOnlyToken(ctx context.Context, checker InternalOnlyCheckFunc) (newCtx context.Context, err error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "kk_grpc internalOnly missing metadata")
	}

	tokens := md.Get(internalOnlyTokenKey)
	if len(tokens) == 0 {
		return nil, status.Error(codes.Unauthenticated, "kk_grpc internalOnly missing token")
	}

	newCtx, err = checker(ctx, tokens[0])
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	return newCtx, nil
}

func verifyJWT(ctx context.Context, checker JWTCheckFunc) (context.Context, error) {
	newCtx, needRefresh, err := checker(ctx, GetAccessToken(ctx))
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	err = grpc.SetTrailer(ctx, metadata.New(map[string]string{
		kk_jwt.NeedRefresh: strconv.FormatBool(needRefresh),
	}))
	if err != nil {
		return nil, status.Error(codes.Internal, "kk_grpc SetTrailer")
	}

	return newCtx, nil
}

// UnaryAuth returns a unary server interceptor that handles authentication based on the AuthConfig.
func UnaryAuth(cfg *AuthConfig) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		interceptorAuth, err := getInterceptorAuth(ctx)
		if err != nil {
			return nil, err
		}
		switch interceptorAuth {
		case kk_scheduler.InterceptorAuth_InternalOnly:
			{
				if cfg.InternalOnlyChecker == nil {
					return nil, status.Error(codes.Internal, "kk_grpc InternalOnlyChecker not configured")
				}
				newCtx, err := verifyInternalOnlyToken(ctx, cfg.InternalOnlyChecker)
				if err != nil {
					return nil, err
				}
				return handler(newCtx, req)
			}
		case kk_scheduler.InterceptorAuth_JWT:
			{
				if cfg.JWTChecker == nil {
					return nil, status.Error(codes.Internal, "kk_grpc JWTChecker not configured")
				}
				newCtx, err := verifyJWT(ctx, cfg.JWTChecker)
				if err != nil {
					return nil, err
				}
				return handler(newCtx, req)
			}
		}

		return handler(ctx, req)
	}
}

// StreamAuth returns a stream server interceptor that handles authentication based on the AuthConfig.
func StreamAuth(cfg *AuthConfig) grpc.StreamServerInterceptor {
	return func(srv any, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		interceptorAuth, err := getInterceptorAuth(stream.Context())
		if err != nil {
			return err
		}

		switch interceptorAuth {
		case kk_scheduler.InterceptorAuth_InternalOnly:
			{
				if cfg.InternalOnlyChecker == nil {
					return status.Error(codes.Internal, "kk_grpc InternalOnlyChecker not configured")
				}
				newCtx, err := verifyInternalOnlyToken(stream.Context(), cfg.InternalOnlyChecker)
				if err != nil {
					return err
				}
				wrapped := middleware.WrapServerStream(stream)
				wrapped.WrappedContext = newCtx
				return handler(srv, wrapped)
			}
		case kk_scheduler.InterceptorAuth_JWT:
			{
				if cfg.JWTChecker == nil {
					return status.Error(codes.Internal, "kk_grpc JWTChecker not configured")
				}
				newCtx, err := verifyJWT(stream.Context(), cfg.JWTChecker)
				if err != nil {
					return err
				}
				wrapped := middleware.WrapServerStream(stream)
				wrapped.WrappedContext = newCtx
				return handler(srv, wrapped)
			}
		}

		return handler(srv, stream)
	}
}
