package main

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/sirupsen/logrus"
	"strconv"
)

func main() {
	app := iris.Default()

	app.Use(func(ctx iris.Context) {
		logrus.Info("中间件")
		ctx.Next()
	})
	app.Get("/hello", func(ctx iris.Context) {
		ctx.WriteString("Hello from Iris!")
	})

	app.Get("/users/{id:uint64}", func(ctx iris.Context) {
		id := ctx.Params().GetUint64Default("id", 0)
		ctx.WriteString(strconv.Itoa(int(id)))
	})

	app.Get("/users/{action:path}", func(ctx iris.Context) {
		action := ctx.Params().Get("action")
		ctx.WriteString(action)
	})

	err := app.Run(iris.Addr(":8080"))
	if err != nil {
		fmt.Println(err)
	}
}
