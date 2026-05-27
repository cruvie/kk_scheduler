package job

import (
	"gitee.com/cruvie/kk_kit/go/kk_stage"
	"github.com/cruvie/kk-scheduler/go/internal/scheduler"
	"github.com/cruvie/kk-scheduler/go/kk_scheduler"
)

func (x *ApiJobDelete) Service(stage *kk_stage.Stage) error {
	span := stage.StartTrace("Service")
	defer span.End()

	err := scheduler.GClient.JobDelete(x.In.GetId())
	return err
}

func (x *ApiJobDisable) Service(stage *kk_stage.Stage) error {
	span := stage.StartTrace("Service")
	defer span.End()

	err := scheduler.GClient.JobDisable(x.In.GetId())
	return err
}

func (x *ApiJobEnable) Service(stage *kk_stage.Stage) error {
	span := stage.StartTrace("Service")
	defer span.End()

	return scheduler.GClient.JobEnable(x.In.GetId())
}

func (x *ApiJobGet) Service(stage *kk_stage.Stage) (*kk_scheduler.PBJob, error) {
	span := stage.StartTrace("Service")
	defer span.End()

	job, err := scheduler.GClient.JobGet(x.In.GetId())
	return job, err
}

func (x *ApiJobList) Service(stage *kk_stage.Stage) ([]*kk_scheduler.PBJob, error) {
	span := stage.StartTrace("Service")
	defer span.End()

	jobList, err := scheduler.GClient.JobList(x.In.GetServiceName())
	return jobList, err
}

func (x *ApiJobSetSpec) Service(stage *kk_stage.Stage) error {
	span := stage.StartTrace("Service")
	defer span.End()

	return scheduler.GClient.JobSetSpec(x.In.GetId(), x.In.GetSpec())
}

func (x *ApiJobTrigger) Service(stage *kk_stage.Stage) error {
	span := stage.StartTrace("Service")
	defer span.End()

	return scheduler.GClient.JobTrigger(x.In.GetId())
}

func (x *ApiJobPut) Service(stage *kk_stage.Stage) error {
	span := stage.StartTrace("Service")
	defer span.End()

	err := scheduler.GClient.JobPut(x.In.GetJob())
	if err != nil {
		return err
	}

	return nil
}
