package g_config

import (
	"context"
	"log/slog"
	"os"

	"gitee.com/cruvie/kk_kit/go/kk_env"
	"gitee.com/cruvie/kk_kit/go/kk_pg"
	"gitee.com/cruvie/kk_kit/go/kk_stage"
	"github.com/BurntSushi/toml"
	"github.com/cruvie/kk-scheduler/go/internal/models"
	"github.com/cruvie/kk-scheduler/go/internal/models/query"
)

func init() {
	kk_env.SetEnv(kk_env.Env(os.Getenv("KK_Schedule")))
}

var Config config

type config struct {
	HttpPort int
	GrpcPort int
	WebPort  int
	Store    struct {
		Choose string
		PG     *kk_pg.Config
	}
	ConfigSlog *kk_stage.ConfigLog `toml:"-"`
}

func InitConfig() *kk_stage.Stage {
	data, err := os.ReadFile("config.toml")
	if err != nil {
		slog.Error("unable to read config.toml", "err", err)
		panic(err)
	}

	_, err = toml.Decode(string(data), &Config)
	if err != nil {
		slog.Error("unable to decode config.toml", "err", err)
		panic(err)
	}

	stage := kk_stage.NewStage(context.Background(), "kk-scheduler")
	{
		Config.ConfigSlog = &kk_stage.ConfigLog{
			Lumberjack: kk_stage.DefaultLogConfig("kk-scheduler"),
			Format:     kk_stage.FormatJSON,
		}
		Config.ConfigSlog.Init()
	}
	{
		switch Config.Store.Choose {
		case "PG":
			models.InitDB(stage, Config.Store.PG)
			Config.Store.PG.Init(stage)
			query.SetDefault(kk_pg.GormClient)
		default:
			panic("store choose error")
		}
	}
	return stage
}

func CloseConfig() {
	{
		switch Config.Store.Choose {
		case "PG":
			Config.Store.PG.Close()
		default:
			panic("store choose error")
		}
	}
	Config.ConfigSlog.Close()
}
