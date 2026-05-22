package models

import (
	"github.com/cruvie/kk-scheduler/go/kk_scheduler"
)

// Service 注册的服务
type Service struct {
	ServiceName string `gorm:"primaryKey;column:service_name;type:text;not null"`
	Target      string `gorm:"column:target;type:text;not null"`
	AuthToken   string `gorm:"column:auth_token;type:text"`
}

func (*Service) TableName() string {
	return "service"
}

func (x *Service) ToPB() *kk_scheduler.PBRegisterService {
	pb := &kk_scheduler.PBRegisterService{}
	pb.SetServiceName(x.ServiceName)
	pb.SetTarget(x.Target)
	pb.SetAuthToken(x.AuthToken)
	return pb
}

func (x *Service) FromPB(pb *kk_scheduler.PBRegisterService) {
	x.ServiceName = pb.GetServiceName()
	x.Target = pb.GetTarget()
	x.AuthToken = pb.GetAuthToken()
}
