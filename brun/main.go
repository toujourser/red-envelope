package main

import (
	"net/http"
	_ "resk"
	"resk/infra"

	_ "net/http/pprof"

	log "github.com/sirupsen/logrus"
	"github.com/tietang/props/ini"
	"github.com/tietang/props/kvs"
)

func main() {
	//通过HTTP服务来开启运行时性能剖析
	go func() {
		log.Info(http.ListenAndServe(":6060", nil))
	}()
	file := kvs.GetCurrentFilePath("config.ini", 1)
	conf := ini.NewIniFileConfigSource(file)
	app := infra.New(conf)
	app.Start()

}
