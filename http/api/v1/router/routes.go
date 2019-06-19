package router

import (
	"github.com/kataras/iris"
	"weassistant/conf"
)

// RegisterAPI1Router 注册api路由
func RegisterAPI1Router(app *iris.Application, extraConf conf.ExtraConfig) {
	// Front
	registerFrontUserRoutes(app, extraConf)
	// Manager
	registerManagerAdministratorRoutes(app, extraConf)
}
