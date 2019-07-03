package front

import (
	jwtmiddleware "github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	apiCommon "weassistant/http/api/common"
	"weassistant/services/orm"
)

// UserAPI 用户控制器
type UserAPI struct {
	BaseAPI
}

// Get 获取用户信息
func (api *UserAPI) Get(ctx iris.Context, jwtHandler *jwtmiddleware.Middleware, userService orm.UserService) mvc.Result {
	userID := api.userID(ctx)
	user, err := userService.Get(userID)
	if err != nil {
		return api.Error(ctx, apiCommon.RetCodeGormQueryFail, ctx.Translate("QueryUserInfoFail"), err, "")
	}
	return api.Output(user)
}

// Delete 注销
func (api *UserAPI) Delete(ctx iris.Context, jwtHandler *jwtmiddleware.Middleware, sessionService orm.UserSessionService) mvc.Result {
	token := ctx.Values().GetString("user_token")
	session, err := sessionService.GetByWhereOptions([]orm.WhereOption{
		orm.WhereOption{Query: "token = ?", Item: []interface{}{token}},
	})
	if err != nil {
		return api.Error(ctx, apiCommon.RetCodeGormQueryFail, ctx.Translate("QueryUserInfoFail"), err, "")
	}
	session.Effective = false
	err = sessionService.Save(&session)
	if err != nil {
		return api.Error(ctx, apiCommon.RetCodeGormQueryFail, ctx.Translate("LogoutFail"), err, "")
	}
	return api.Success(ctx.Translate("LogoutSuccess"), nil)
}
