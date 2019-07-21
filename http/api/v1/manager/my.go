package manager

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	apiCommon "weassistant/http/api/common"
	"weassistant/services/orm"
)

// MyAPI 用户控制器
type MyAPI struct {
	BaseAPI
}

// Get 获取用户信息
func (api *MyAPI) Get(ctx iris.Context, administratorService orm.AdministratorService) mvc.Result {
	admin := api.admin(ctx, administratorService)
	return api.Output(admin)
}

// Delete 注销
func (api *MyAPI) Delete(ctx iris.Context, sessionService orm.AdministratorSessionService) mvc.Result {
	token := ctx.Values().GetString("administrator_token")
	session, err := sessionService.GetByWhereOptions([]orm.WhereOption{
		orm.WhereOption{Query: "token = ?", Item: []interface{}{token}},
	})
	if err != nil {
		return api.Error(ctx, apiCommon.RetCodeGormQueryFail, ctx.Translate("QueryAdministratorInfoFail"), err, "")
	}
	session.Effective = false
	err = sessionService.Save(&session)
	if err != nil {
		return api.Error(ctx, apiCommon.RetCodeGormQueryFail, ctx.Translate("LogoutFail"), err, "")
	}
	return api.Success(ctx.Translate("LogoutSuccess"), nil)
}
