package resk

import (
	"resk/apis/gorpc"
	_ "resk/apis/gorpc"
	_ "resk/apis/web"
	_ "resk/core/accounts"
	_ "resk/core/envelops"
	"resk/infra"
	"resk/infra/base"
	"resk/jobs"
)

func init() {
	infra.Register(&base.LogoStarter{})
	infra.Register(&base.PropsStarter{})
	infra.Register(&base.DbxDatabaseStarter{})
	infra.Register(&base.ValidatorStarter{})
	infra.Register(&base.GoRPCStarter{})
	infra.Register(&gorpc.GoRpcApiStarter{})
	infra.Register(&jobs.RefundExpiredJobStarter{})
	infra.Register(&base.IrisServerStarter{})
	infra.Register(&infra.WebApiStarter{})
	infra.Register(&base.EurekaStarter{})
	infra.Register(&base.HookStarter{})
}
