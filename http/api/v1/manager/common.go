package manager

import (
	"github.com/kataras/iris"
	apiCommon "weassistant/http/api/common"
	"weassistant/http/api/v1/common"
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
