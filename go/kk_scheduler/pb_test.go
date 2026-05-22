package kk_scheduler

import (
	"testing"

	"gitee.com/cruvie/kk_go_kit/kk_protobuf"
	"github.com/robfig/cron/v3"
)

func TestGenerateProtoFile(t *testing.T) {
	packageName := "kk_scheduler"
	kk_protobuf.GenerateProtoFile(PBRegisterService{}, packageName, false, false)
	kk_protobuf.GenerateProtoFile(PBRegisterJob{}, packageName, false, false)
}

func TestModelToFromPB(t *testing.T) {
	kk_protobuf.GenToFromPB(cron.Entry{})
}
