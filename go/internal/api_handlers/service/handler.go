package service

import (
	"gitee.com/cruvie/kk_kit/go/kk_stage"
	"github.com/cruvie/kk_scheduler/go/kk_scheduler"
)

func (x *ApiServiceDelete) Handler(stage *kk_stage.Stage) (*kk_scheduler.ServiceDelete_Output, error) {
	err := x.Service(stage)
	if err != nil {
		return nil, err
	}
	return &kk_scheduler.ServiceDelete_Output{}, nil
}

func (x *ApiServiceGet) Handler(stage *kk_stage.Stage) (*kk_scheduler.ServiceGet_Output, error) {
	service, err := x.Service(stage)
	if err != nil {
		return nil, err
	}

	out := &kk_scheduler.ServiceGet_Output{}
	out.SetService(service)
	return out, nil
}

func (x *ApiServiceList) Handler(stage *kk_stage.Stage) (*kk_scheduler.ServiceList_Output, error) {
	service, err := x.Service(stage)
	if err != nil {
		return nil, err
	}
	out := &kk_scheduler.ServiceList_Output{}
	out.SetServiceList(service)
	return out, nil
}

func (x *ApiServicePut) Handler(stage *kk_stage.Stage) (*kk_scheduler.ServicePut_Output, error) {
	err := x.Service(stage)
	if err != nil {
		return nil, err
	}
	return &kk_scheduler.ServicePut_Output{}, nil
}
