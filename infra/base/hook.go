package base

import (
	"os"
	"os/signal"
	"reflect"
	"resk/infra"
	"syscall"

	log "github.com/sirupsen/logrus"
)

var callbacks []func()

func Register(callback func()) {
	callbacks = append(callbacks, callback)
}

type HookStarter struct {
	infra.BaseStarter
}

func (s *HookStarter) Init(ctx infra.StarterContext) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		for {
			c := <-sigs
			log.Info("Received signal: ", c)
			for _, fn := range callbacks {
				fn()
			}
			os.Exit(0)
		}
	}()
}

func (s *HookStarter) Start(ctx infra.StarterContext) {
	starts := infra.GetStarters()
	for _, start := range starts {
		typ := reflect.TypeOf(start)
		log.Infof("[Register Notify Stop]: %s.Stop()", typ.String())
		Register(func() {
			start.Stop(ctx)
		})
	}
}
