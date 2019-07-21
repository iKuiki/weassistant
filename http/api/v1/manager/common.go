package manager

import (
	"github.com/kataras/iris"
	apiCommon "weassistant/http/api/common"
	"weassistant/http/api/v1/common"
	"weassistant/models"
	"weassistant/services/orm"
)

// BaseAPI 基础API
type BaseAPI struct {
	common.BaseAPI
}

// adminID 获取当前用户的adminID
func (api *BaseAPI) adminID(ctx iris.Context) (adminID uint64) {
	adminID, err := ctx.Values().GetUint64("administrator_id")
	if err != nil {
		// 正常情况下needlogin模块应当已经将adminID写入了memstore中，不应当发生此错误
		panic(api.Error(ctx, apiCommon.RetCodeUnknownError, ctx.Translate("ServerError"), err, "adminID miss after needlogin middleware"))
	}
	return
}

func (api *BaseAPI) admin(ctx iris.Context, administratorService orm.AdministratorService) (admin models.Administrator) {
	adminID := api.adminID(ctx)
	admin, err := administratorService.Get(adminID)
	if err != nil {
		// 如果出现错误，则直接报错吧
		panic(api.Error(ctx, apiCommon.RetCodeGormQueryFail, ctx.Translate("QueryAdministratorInfoFail"), err, ""))
	}
	return
}
