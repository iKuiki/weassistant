package manager

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"weassistant/http/api/common"
	"weassistant/services/orm"
)

// UserAPI 用户管理API服务
type UserAPI struct {
	BaseAPI
}

// Get 获取用户列表
func (api *UserAPI) Get(ctx iris.Context, userService orm.UserService) mvc.Result {
	limit, offset := api.ObtainLimitOffset(ctx, true)
	whereOptions := []orm.WhereOption{}
	users, err := userService.GetListByWhereOptions(whereOptions, []string{}, limit, offset)
	if err != nil {
		return api.Error(ctx, common.RetCodeGormQueryFail, ctx.Translate("QueryUsersError"), err, "userService.GetListByWhereOptions error: "+err.Error())
	}
	size := api.Size(ctx, userService, whereOptions)
	return api.Output(map[string]interface{}{
		"list": users,
		"pageinfo": map[string]interface{}{
			"total": size,
		},
	})
}
