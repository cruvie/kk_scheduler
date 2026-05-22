package main

import (
	"fmt"
	"net/http"

	"gitee.com/cruvie/kk_go_kit/kk_server"
	"gitee.com/cruvie/kk_go_kit/kk_stage"
	"github.com/cruvie/kk-scheduler/go/internal/g_config"
)

func NewWebServer(stage *kk_stage.Stage) *kk_server.KKRunServer {
	fs := http.FileServer(http.Dir("public"))
	http.Handle("/", fs)

	run := func() {
		err := http.ListenAndServe(fmt.Sprintf(":%d", g_config.Config.WebPort), nil)
		if err != nil {
			panic(err)
		}
	}
	done := func(quitCh <-chan struct{}) {
		<-quitCh
	}

	return &kk_server.KKRunServer{
		Run:  run,
		Done: done,
	}
}
