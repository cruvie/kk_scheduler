package service

import (
	"gitee.com/cruvie/kk_kit/go/kk_stage"
	"github.com/cruvie/kk_scheduler/go/internal/scheduler"
	"github.com/cruvie/kk_scheduler/go/kk_scheduler"
)

func (x *ApiServiceDelete) Service(stage *kk_stage.Stage) error {
	span := stage.StartTrace("Service")
	defer span.End()

	return scheduler.GClient.ServiceDelete(x.In.GetServiceName())
}

func (x *ApiServiceGet) Service(stage *kk_stage.Stage) (*kk_scheduler.PBRegisterService, error) {
	span := stage.StartTrace("Service")
	defer span.End()

	service, err := scheduler.GClient.ServiceGet(x.In.GetServiceName())
	if err != nil {
		return nil, err
	}

	return service, nil
}

func (x *ApiServiceList) Service(stage *kk_stage.Stage) ([]*kk_scheduler.PBRegisterService, error) {
	span := stage.StartTrace("Service")
	defer span.End()

	service, err := scheduler.GClient.ServiceList()
	return service, err
}

func (x *ApiServicePut) Service(stage *kk_stage.Stage) error {
	span := stage.StartTrace("Service")
	defer span.End()

	return scheduler.GClient.ServicePut(x.In.GetService())
}
