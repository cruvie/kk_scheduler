package job

import (
	"gitee.com/cruvie/kk_go_kit/kk_stage"
	"github.com/robfig/cron/v3"
)

func (x *ApiJobDelete) CheckInput(stage *kk_stage.Stage) error {
	return nil
}

func (x *ApiJobDisable) CheckInput(stage *kk_stage.Stage) error {
	return nil
}

func (x *ApiJobEnable) CheckInput(stage *kk_stage.Stage) error {
	return nil
}

func (x *ApiJobGet) CheckInput(stage *kk_stage.Stage) error {
	return nil
}

func (x *ApiJobList) CheckInput(stage *kk_stage.Stage) error {
	return nil
}

func (x *ApiJobSetSpec) CheckInput(stage *kk_stage.Stage) error {
	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	_, err := parser.Parse(x.In.GetSpec())
	if err != nil {
		return err
	}
	return nil
}

func (x *ApiJobTrigger) CheckInput(stage *kk_stage.Stage) error {
	return nil
}

func (x *ApiJobPut) CheckInput(stage *kk_stage.Stage) error {
	return nil
}
