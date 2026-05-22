package query

import (
	"context"
	"os"
	"testing"

	"gitee.com/cruvie/kk_go_kit/kk_env"
	"gitee.com/cruvie/kk_go_kit/kk_pg"
	"gitee.com/cruvie/kk_go_kit/kk_stage"
	"github.com/cruvie/kk-scheduler/go/internal/models"
)

func init() {
	kk_env.SetEnv(kk_env.Env(os.Getenv("KK_Schedule")))
}

var pg = &kk_pg.Config{DSN: kk_pg.PostgresDSN{
	Host:     "127.0.0.1",
	Port:     5432,
	User:     "postgres",
	Password: "testpg",
	DBName:   "kk_scheduler",
	Schema:   "",
	SSLMode:  "disable",
	TimeZone: "UTC",
	Addition: nil,
}}

func TestGen(t *testing.T) {
	kk_env.SetEnv(kk_env.Env(os.Getenv("KK_Schedule")))
	stage := kk_stage.NewStage(context.Background(), "test")
	kk_pg.GenQuery(stage, pg,
		models.TaskExecution{},
		models.Job{},
		models.Service{},
	)
}

func TestCreateTable(t *testing.T) {
	pg.Init(kk_stage.NewNoopStage())
	kk_pg.CreateTables(kk_pg.GormClient,
		models.TaskExecution{},
		models.Job{},
		models.Service{},
	)
}
