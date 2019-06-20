package router

import (
	"github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris/core/router"
	"weassistant/http/api/v1/manager"
	api1ManagerMiddleware "weassistant/http/api/v1/manager/middleware"
	"weassistant/services/orm"

	"github.com/kataras/iris/mvc"
)

type managerModuleConf interface {
	GetMgrJwtMiddleware() *jwt.Middleware
	GetAdministratorService() orm.AdministratorService
	GetAdministratorSessionService() orm.AdministratorSessionService
	GetUserService() orm.UserService
	GetJwtValidationKey() string
}

// 后台管理模块路由总成
func registerManagerRoutes(managerAPI router.Party, managerConf managerModuleConf) {
	// 先是不需要login的Auth相关接口
	registerManagerAuthRoutes(managerAPI, managerConf)
	// 创建登录验证中间件
	needLoginMiddleware := api1ManagerMiddleware.MustNewNeedLoginMiddleware(managerConf.GetMgrJwtMiddleware(), managerConf.GetAdministratorSessionService())
	// 之后的用户信息相关的接口需要登录后才能调用
	managerAPI.Use(needLoginMiddleware.Serve)
	registerManagerMyRoutes(managerAPI, managerConf, needLoginMiddleware)
	registerManagerUserRoutes(managerAPI, managerConf)
}

// 认证服务路由
func registerManagerAuthRoutes(managerAPI router.Party, managerConf managerModuleConf) {
	mvc.New(managerAPI.Party("/auth")).
		Register(managerConf.GetMgrJwtMiddleware()).
		Register(managerConf.GetAdministratorService()).
		Register(managerConf.GetAdministratorSessionService()).
		Register(managerConf.GetUserService()).
		Register(managerConf.GetJwtValidationKey()).
		Handle(new(manager.AuthAPI))
}

// registerManagerMyRoutes 注册后台我的信息服务路由
func registerManagerMyRoutes(managerAPI router.Party, managerConf managerModuleConf, needLoginMiddleware api1ManagerMiddleware.NeedLoginMiddleware) {
	mvc.New(managerAPI.Party("/my")).
		Register(managerConf.GetMgrJwtMiddleware()).
		Register(managerConf.GetAdministratorService()).
		Handle(new(manager.AdministratorAPI))
}

// 用户管理模块路由
func registerManagerUserRoutes(managerAPI router.Party, managerConf managerModuleConf) {
	mvc.New(managerAPI.Party("/user")).
		Register(managerConf.GetUserService()).
		Handle(new(manager.UserAPI))
}
