package router

import (
	"github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris/core/router"
	"weassistant/http/api/v1/front"
	api1FrontMiddleware "weassistant/http/api/v1/front/middleware"
	"weassistant/services/locker"
	"weassistant/services/orm"

	"github.com/kataras/iris/mvc"
)

type frontConfig interface {
	GetAPIJwtMiddleware() *jwt.Middleware
	GetRegisterLockerService() locker.CommonLockerService
	GetUserService() orm.UserService
	GetUserSessionService() orm.UserSessionService
	GetJwtValidationKey() string
}

// registerFrontRoutes 前台服务路由总成
func registerFrontRoutes(apiParty router.Party, frontConf frontConfig) {
	// 先是不需要login的Auth相关接口
	registerFrontAuthRoutes(apiParty, frontConf)
	// 创建登录验证中间件
	needLoginMiddleware := api1FrontMiddleware.MustNewNeedLoginMiddleware(frontConf.GetAPIJwtMiddleware(), frontConf.GetUserSessionService())
	// 之后的用户信息相关的接口需要登录后才能调用
	apiParty.Use(needLoginMiddleware.Serve)
	registerFrontUserRoutes(apiParty, frontConf)
}

// 登陆认证相关
func registerFrontAuthRoutes(apiParty router.Party, frontConf frontConfig) {
	mvc.New(apiParty.Party("/user")).
		Register(frontConf.GetAPIJwtMiddleware()).
		Register(frontConf.GetRegisterLockerService()).
		Register(frontConf.GetUserService()).
		Register(frontConf.GetUserSessionService()).
		Register(frontConf.GetJwtValidationKey()).
		Handle(new(front.AuthAPI))
}

// registerFrontUserRoutes 注册前台用户服务路由
func registerFrontUserRoutes(apiParty router.Party, frontConf frontConfig) {
	mvc.New(apiParty.Party("/user")).
		Register(frontConf.GetAPIJwtMiddleware()).
		Register(frontConf.GetUserService()).
		Handle(new(front.UserAPI))
}
