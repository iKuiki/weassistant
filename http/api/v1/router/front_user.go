package router

import (
	"github.com/iris-contrib/middleware/jwt"
	"weassistant/http/api/v1/front"
	api1FrontMiddleware "weassistant/http/api/v1/front/middleware"
	"weassistant/services/locker"
	"weassistant/services/orm"

	"github.com/kataras/iris/mvc"

	"github.com/kataras/iris"
)

type frontConfig interface {
	GetAPIJwtMiddleware() *jwt.Middleware
	GetRegisterLockerService() locker.CommonLockerService
	GetUserService() orm.UserService
	GetUserSessionService() orm.UserSessionService
	GetJwtValidationKey() string
}

// registerFrontUserRoutes 注册前台用户服务路由
func registerFrontUserRoutes(app *iris.Application, frontConf frontConfig) {
	apiAPI := app.Party("/api/v1")
	// 创建登录验证中间件
	needLoginMiddleware := api1FrontMiddleware.MustNewNeedLoginMiddleware(frontConf.GetAPIJwtMiddleware(), frontConf.GetUserSessionService())
	// 先是不需要login的Auth相关接口
	mvc.New(apiAPI.Party("/user")).
		Register(frontConf.GetAPIJwtMiddleware()).
		Register(frontConf.GetRegisterLockerService()).
		Register(frontConf.GetUserService()).
		Register(frontConf.GetUserSessionService()).
		Register(frontConf.GetJwtValidationKey()).
		Handle(new(front.AuthAPI))
	// 之后的用户信息相关的接口需要登录后才能调用
	apiAPI.Use(needLoginMiddleware.Serve)
	mvc.New(apiAPI.Party("/user")).
		Register(frontConf.GetAPIJwtMiddleware()).
		Register(frontConf.GetUserService()).
		Register(needLoginMiddleware).
		Handle(new(front.UserAPI))
}
