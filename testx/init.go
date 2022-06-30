package testx

import (
	"resk/infra"
	"resk/infra/base"

	"github.com/tietang/props/ini"
)

func init() {
	file := "/Users/mds/Documents/study/golang/resk/brun/config.ini"
	conf := ini.NewIniFileConfigSource(file)

	infra.Register(&base.PropsStarter{})
	infra.Register(&base.DbxDatabaseStarter{})
	infra.Register(&base.ValidatorStarter{})
	//infra.Register(&base.IrisServerStarter{})

	app := infra.New(conf)
	app.Start()
}
