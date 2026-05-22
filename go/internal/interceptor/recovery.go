package interceptor

import (
	"fmt"
	"runtime/debug"

	"gitee.com/cruvie/kk_go_kit/kk_stage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func PanicRecovery(p any) (err error) {
	kk_stage.Print2Std(fmt.Sprintf("panic: %v\nstack: %v\n", p, string(debug.Stack())), kk_stage.WithRaw())
	return status.Errorf(codes.Internal, "%s", p)
}
