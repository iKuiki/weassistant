package router

import (
	"github.com/iris-contrib/middleware/jwt"
	"weassistant/http/api/v1/manager"
	api1ManagerMiddleware "weassistant/http/api/v1/manager/middleware"
	"weassistant/services/orm"

	"github.com/kataras/iris/mvc"

	"github.com/kataras/iris"
)

type managerModuleConf interface {
	GetMgrJwtMiddleware() *jwt.Middleware
	GetAdministratorService() orm.AdministratorService
	GetAdministratorSessionService() orm.AdministratorSessionService
	GetJwtValidationKey() string
}

// registerManagerAdministratorRoutes 注册后台管理员服务路由
func registerManagerAdministratorRoutes(app *iris.Application, managerConf managerModuleConf) {
	APIv1 := app.Party("/api/v1")
	managerAPI := APIv1.Party("/manager")
	// 创建登录验证中间件
	needLoginMiddleware := api1ManagerMiddleware.MustNewNeedLoginMiddleware(managerConf.GetMgrJwtMiddleware(), managerConf.GetAdministratorSessionService())
	// 先是不需要login的Auth相关接口
	mvc.New(managerAPI.Party("/administrator")).
		Register(managerConf.GetMgrJwtMiddleware()).
		Register(managerConf.GetAdministratorService()).
		Register(managerConf.GetAdministratorSessionService()).
		Register(managerConf.GetJwtValidationKey()).
		Handle(new(manager.AuthAPI))
	// 之后的用户信息相关的接口需要登录后才能调用
	managerAPI.Use(needLoginMiddleware.Serve)
	mvc.New(managerAPI.Party("/administrator")).
		Register(managerConf.GetMgrJwtMiddleware()).
		Register(managerConf.GetAdministratorService()).
		Register(needLoginMiddleware).
		Handle(new(manager.AdministratorAPI))
}
