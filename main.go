package main

import (
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/i18n"
	"github.com/kataras/iris/middleware/logger"
	"weassistant/conf"
	apiCommon "weassistant/http/api/common"
	api1Router "weassistant/http/api/v1/router"
)

func main() {
	config := conf.MustNewConfig()
	err := config.Load("config.json")
	if err != nil {
		panic(err)
	}
	extraConf := conf.MustExtraNewConfig(config)
	// 创建iris app
	app := iris.New()
	if config.GetDebug() {
		app.Logger().SetLevel("debug")
		app.Logger().Info("enable debug logger level")
	}
	// log必须在错误容器上方，否则会失效
	app.Use(logger.New(logger.DefaultConfig()))
	// 错误捕获
	apiCommon.SetSentryClient(config.GetSentryClient())
	app.Use(apiCommon.ErrHandler)
	// 准备i18n
	app.Use(i18n.New(i18n.Config{
		Default:      "en",
		URLParameter: "lang",
		Languages: map[string]string{
			"en":    "http/locales/en-US.ini",
			"en-US": "http/locales/en-US.ini",
			"zh":    "http/locales/zh-CN.ini",
			"zh-CN": "http/locales/zh-CN.ini",
		},
	}))
	api1Router.RegisterAPI1Router(app, extraConf)
	app.Run(
		iris.Addr(
			fmt.Sprintf(":%d", config.GetAppPort()),
		),
	)
}
