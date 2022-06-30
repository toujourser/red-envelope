package base

import (
	"github.com/common-nighthawk/go-figure"
	"resk/infra"
)

type LogoStarter struct {
	infra.BaseStarter
}

func (p *LogoStarter) Init(ctx infra.StarterContext) {
	figure.NewFigure("RESK", "", true).Print()
}
