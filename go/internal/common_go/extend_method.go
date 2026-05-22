package common_go

import (
	"github.com/cruvie/kk-scheduler/go/kk_scheduler"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func MethodDescGetInterceptorAuth(method protoreflect.MethodDescriptor) kk_scheduler.InterceptorAuth {
	options := method.Options()
	if options == nil {
		return kk_scheduler.InterceptorAuth_UNSPECIFIED
	}

	interceptorList := proto.GetExtension(options, kk_scheduler.E_InterceptorAuthList)
	if interceptorList == nil {
		return kk_scheduler.InterceptorAuth_UNSPECIFIED
	}

	if list, ok := interceptorList.([]kk_scheduler.InterceptorAuth); ok {
		switch len(list) {
		case 0:
			return kk_scheduler.InterceptorAuth_UNSPECIFIED
		case 1:
			return list[0]
		}
		if len(list) > 1 {
			// only support one InterceptorAuth
			panic("MethodDescGetInterceptorAuth: len(list) > 1")
		}
	}

	return kk_scheduler.InterceptorAuth_UNSPECIFIED
}

func MethodDescGetApiName(method protoreflect.MethodDescriptor) string {
	options := method.Options()
	if options == nil {
		return ""
	}

	serviceName := proto.GetExtension(options, kk_scheduler.E_ApiName)
	if serviceName == nil {
		return ""
	}

	if name, ok := serviceName.(string); ok {
		return name
	}

	return ""
}
