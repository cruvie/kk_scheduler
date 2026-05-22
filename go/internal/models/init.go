package models

import (
	"gitee.com/cruvie/kk_go_kit/kk_pg"
	"gitee.com/cruvie/kk_go_kit/kk_stage"
)

func InitDB(stage *kk_stage.Stage, pg *kk_pg.Config) {
	kk_pg.CreateDB(pg, pg.DSN.DBName)
	kk_pg.CreateTables(pg.NewDefaultDB(stage),
		TaskExecution{},
		Job{},
		Service{},
	)
}
