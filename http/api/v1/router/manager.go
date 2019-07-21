package router

import (
	"github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris/context"
	"github.com/kataras/iris/core/router"
	"weassistant/http/api/v1/manager"
	api1ManagerMiddleware "weassistant/http/api/v1/manager/middleware"
	"weassistant/services/locker"
	"weassistant/services/orm"

	"github.com/kataras/iris/mvc"
)

type managerModuleConf interface {
	GetMgrJwtMiddleware() *jwt.Middleware
	GetAdministratorService() orm.AdministratorService
	GetAdministratorSessionService() orm.AdministratorSessionService
	GetUserService() orm.UserService
	GetJwtValidationKey() string
	GetRegisterLockerService() locker.CommonLockerService
}

// 后台管理模块路由总成
func registerManagerRoutes(managerAPI router.Party, managerConf managerModuleConf, handles ...context.Handler) {
	// 先是不需要login的Auth相关接口
	registerManagerAuthRoutes(managerAPI, managerConf, handles...)
	// 创建登录验证中间件
	needLoginMiddleware := api1ManagerMiddleware.MustNewNeedLoginMiddleware(managerConf.GetMgrJwtMiddleware(), managerConf.GetAdministratorSessionService())
	// 之后的用户信息相关的接口需要登录后才能调用
	managerAPI.Use(needLoginMiddleware.Serve)
	registerManagerMyRoutes(managerAPI, managerConf, needLoginMiddleware, handles...)
	registerManagerUserRoutes(managerAPI, managerConf, handles...)
}

// 认证服务路由
func registerManagerAuthRoutes(managerAPI router.Party, managerConf managerModuleConf, handles ...context.Handler) {
	mvc.New(managerAPI.Party("/auth", handles...)).
		Register(managerConf.GetMgrJwtMiddleware()).
		Register(managerConf.GetAdministratorService()).
		Register(managerConf.GetAdministratorSessionService()).
		Register(managerConf.GetUserService()).
		Register(managerConf.GetJwtValidationKey()).
		Handle(new(manager.AuthAPI))
}

// registerManagerMyRoutes 注册后台我的信息服务路由
func registerManagerMyRoutes(managerAPI router.Party, managerConf managerModuleConf, needLoginMiddleware api1ManagerMiddleware.NeedLoginMiddleware, handles ...context.Handler) {
	mvc.New(managerAPI.Party("/my", handles...)).
		Register(managerConf.GetAdministratorService()).
		Register(managerConf.GetAdministratorSessionService()).
		Handle(new(manager.MyAPI))
}

// 用户管理模块路由
func registerManagerUserRoutes(managerAPI router.Party, managerConf managerModuleConf, handles ...context.Handler) { // 允许cors跨域访问
	mvc.New(managerAPI.Party("/user", handles...)).
		Register(managerConf.GetUserService()).
		Register(managerConf.GetRegisterLockerService()).
		Handle(new(manager.UserAPI))
}
