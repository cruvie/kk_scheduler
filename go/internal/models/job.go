package models

import (
	"github.com/cruvie/kk_scheduler/go/kk_scheduler"
)

type Job struct {
	Id          string `gorm:"primaryKey;column:id;type:uuid"`
	ServiceName string `gorm:"column:service_name;type:text;not null;uniqueIndex:idx_service_func"`
	FuncName    string `gorm:"column:func_name;type:text;not null;uniqueIndex:idx_service_func"`
	EntryID     int32  `gorm:"column:entry_id;type:integer"`
	Enabled     bool   `gorm:"column:enabled;type:boolean;not null"`
	Spec        string `gorm:"column:spec;type:text;not null"`
	Description string `gorm:"column:description;type:text"`
}

func (*Job) TableName() string {
	return "job"
}

func (x *Job) ToPB() *kk_scheduler.PBJob {
	pb := &kk_scheduler.PBJob{}
	pb.SetId(x.Id)
	pb.SetEntryID(x.EntryID)
	pb.SetEnabled(x.Enabled)
	pb.SetSpec(x.Spec)
	pb.SetDescription(x.Description)
	pb.SetFuncName(x.FuncName)
	pb.SetServiceName(x.ServiceName)
	return pb
}

func (x *Job) FromPB(pb *kk_scheduler.PBJob) {
	x.Id = pb.GetId()
	x.EntryID = pb.GetEntryID()
	x.Enabled = pb.GetEnabled()
	x.Spec = pb.GetSpec()
	x.Description = pb.GetDescription()
	x.FuncName = pb.GetFuncName()
	x.ServiceName = pb.GetServiceName()
}
