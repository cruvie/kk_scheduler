package scheduler

import (
	"fmt"
	"log/slog"

	"gitee.com/cruvie/kk_go_kit/kk_server"
	"github.com/cruvie/kk-scheduler/go/internal/store_driver"
	"github.com/robfig/cron/v3"
)

type Config struct {
	Opts        []cron.Option            `json:"-" yaml:"-"`
	StoreDriver store_driver.StoreDriver `json:"-" yaml:"-"`
}

func (x *Config) check() {
	if x.StoreDriver == nil {
		panic("StoreDriver is nil")
	}
}

func NewScheduleServer() *kk_server.KKRunServer {
	run := func() {
		cfg := &Config{
			StoreDriver: store_driver.NewStoreDriver(),
		}
		logger := kKScheduleLog{}

		cfg.Opts = append(cfg.Opts,
			cron.WithChain(cron.Recover(logger)),
			cron.WithLogger(logger),
		)
		InitGClient(cfg)
		GClient.Start()
	}
	done := func(quitCh <-chan struct{}) {
		<-quitCh
		GClient.Close()
	}
	return &kk_server.KKRunServer{
		Run:  run,
		Done: done,
	}
}

type kKScheduleLog struct{}

func (K kKScheduleLog) Info(msg string, keysAndValues ...any) {
	slog.Info(fmt.Sprintf("KKSchedule Msg:[%s]", msg), keysAndValues...)
}

func (K kKScheduleLog) Error(err error, msg string, keysAndValues ...any) {
	slog.Error(fmt.Sprintf("KKSchedule Msg:[%s]: Error:[%v]", msg, err), keysAndValues...)
}
