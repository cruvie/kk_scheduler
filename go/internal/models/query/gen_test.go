package query

import (
	"os"
	"testing"

	"gitee.com/cruvie/kk_kit/go/kk_env"
	"gitee.com/cruvie/kk_kit/go/kk_pg"
	"gitee.com/cruvie/kk_kit/go/kk_stage"
	"github.com/cruvie/kk_scheduler/go/internal/models"
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
	kk_pg.GenQuery(
		models.TaskExecution{},
		models.Job{},
		models.Service{},
	)
}

func TestCreateTable(t *testing.T) {
	pg.Init(kk_stage.NewNoopStage())
	kk_pg.CreateTables(
		kk_pg.GormClient,
		models.TaskExecution{},
		models.Job{},
		models.Service{},
	)
}
