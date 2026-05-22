package kk_scheduler

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func CheckAuthority(rpcCtx context.Context, token string) error {
	md, ok := metadata.FromIncomingContext(rpcCtx)
	if !ok {
		return status.Error(codes.Unauthenticated, "missing metadata")
	}

	var authority string
	if auth := md[":authority"]; len(auth) > 0 {
		authority = auth[0]
	} else {
		return status.Error(codes.Unauthenticated, "missing authority")
	}

	if authority != token {
		return status.Error(codes.Unauthenticated, "invalid authority")
	}
	return nil
}
