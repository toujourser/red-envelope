package base

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	recover2 "github.com/kataras/iris/v12/middleware/recover"
	"github.com/sirupsen/logrus"
	"resk/infra"
	"time"
)

var irisApplication *iris.Application

func Iris() *iris.Application {
	return irisApplication
}

type IrisServerStarter struct {
	infra.BaseStarter
}

func (i *IrisServerStarter) Init(ctx infra.StarterContext) {
	// 创建Iris 实例
	irisApplication = initIris()
	// 日志组件配置和扩展
	logger := irisApplication.Logger()
	logger.Install(logrus.StandardLogger())
}

func (i *IrisServerStarter) Start(ctx infra.StarterContext) {
	Iris().Logger().SetLevel(ctx.Props().GetDefault("log.level", "info"))

	// 路由信息打印到控制台
	routes := Iris().GetRoutes()
	for _, r := range routes {
		logrus.Info(r.Path)
	}
	// 启动iris
	port := ctx.Props().GetDefault("app.server.port", "8080")
	Iris().Run(iris.Addr(":" + port))
}

func (i *IrisServerStarter) StartBlocking() bool {
	return true
}

func initIris() *iris.Application {
	app := iris.New()
	app.Use(recover2.New())
	cfg := logger.Config{
		Status:             true,
		IP:                 true,
		Method:             true,
		Path:               false,
		PathAfterHandler:   false,
		Query:              true,
		TraceRoute:         false,
		MessageContextKeys: nil,
		MessageHeaderKeys:  nil,
		LogFunc: func(endTime time.Time, latency time.Duration, status, ip, method, path string, message interface{}, headerMessage interface{}) {
			app.Logger().Infof("| %s | %s | %3d | %15s | %s | %s | %s | %s |",
				endTime.Format("2006-01-02 15:04:05.000000"), latency.String(),
				status, ip, method, path, message, headerMessage,
			)
		},
		LogFuncCtx: nil,
		Skippers:   nil,
	}
	app.Use(logger.New(cfg))
	return app
}
