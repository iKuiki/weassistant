package router

import (
	"github.com/kataras/iris"
	"weassistant/conf"
)

// RegisterAPI1Router 注册api路由
func RegisterAPI1Router(app *iris.Application, extraConf conf.ExtraConfig) {
	// v1版api
	APIv1 := app.Party("/api/v1")
	// Front
	registerFrontRoutes(APIv1, extraConf)
	// Manager
	managerAPI := app.Party("/api/v1/manager")
	registerManagerRoutes(managerAPI, extraConf)
}
