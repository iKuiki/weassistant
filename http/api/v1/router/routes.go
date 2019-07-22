package router

import (
	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris"
	"weassistant/conf"
)

// RegisterAPI1Router 注册api路由
func RegisterAPI1Router(app *iris.Application, extraConf conf.ExtraConfig) {
	// 处理cors请求
	// 此处创建的cors中间件负责handle所有options请求（为了避免需要登陆的路由无法正确处理options请求
	// 除此处外，每个Party也要单独添加cors中间件，因为复杂method（如delete）需要返回时也验证cors头
	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
	})
	// v1版api
	APIv1 := app.Party("/api/v1", crs)
	// 处理所有options请求
	APIv1.Options("/*", func(ctx iris.Context) {})
	// Front
	registerFrontRoutes(APIv1, extraConf)
	// Manager
	managerAPI := app.Party("/api/v1/manager")
	registerManagerRoutes(managerAPI, extraConf, crs)
}
