package job

import (
	"gitee.com/cruvie/kk_kit/go/kk_stage"
	"github.com/cruvie/kk-scheduler/go/kk_scheduler"
)

func (x *ApiJobDelete) Handler(stage *kk_stage.Stage) (*kk_scheduler.JobDelete_Output, error) {
	err := x.Service(stage)
	if err != nil {
		return nil, err
	}
	return &kk_scheduler.JobDelete_Output{}, nil
}

func (x *ApiJobDisable) Handler(stage *kk_stage.Stage) (*kk_scheduler.JobDisable_Output, error) {
	err := x.Service(stage)
	if err != nil {
		return nil, err
	}
	return &kk_scheduler.JobDisable_Output{}, nil
}

func (x *ApiJobEnable) Handler(stage *kk_stage.Stage) (*kk_scheduler.JobEnable_Output, error) {
	err := x.Service(stage)
	if err != nil {
		return nil, err
	}
	return &kk_scheduler.JobEnable_Output{}, nil
}

func (x *ApiJobGet) Handler(stage *kk_stage.Stage) (*kk_scheduler.JobGet_Output, error) {
	job, err := x.Service(stage)
	if err != nil {
		return nil, err
	}
	out := &kk_scheduler.JobGet_Output{}
	out.SetJob(job)
	return out, nil
}

func (x *ApiJobList) Handler(stage *kk_stage.Stage) (*kk_scheduler.JobList_Output, error) {
	jobList, err := x.Service(stage)
	if err != nil {
		return nil, err
	}
	out := &kk_scheduler.JobList_Output{}
	out.SetJobList(jobList)
	return out, nil
}

func (x *ApiJobSetSpec) Handler(stage *kk_stage.Stage) (*kk_scheduler.JobSetSpec_Output, error) {
	err := x.Service(stage)
	if err != nil {
		return nil, err
	}
	return &kk_scheduler.JobSetSpec_Output{}, nil
}

func (x *ApiJobTrigger) Handler(stage *kk_stage.Stage) (*kk_scheduler.JobTrigger_Output, error) {
	err := x.Service(stage)
	if err != nil {
		return nil, err
	}
	return &kk_scheduler.JobTrigger_Output{}, nil
}

func (x *ApiJobPut) Handler(stage *kk_stage.Stage) (*kk_scheduler.JobPut_Output, error) {
	err := x.Service(stage)
	if err != nil {
		return nil, err
	}
	return &kk_scheduler.JobPut_Output{}, nil
}
