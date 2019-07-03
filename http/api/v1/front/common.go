package front

import (
	"github.com/kataras/iris"
	apiCommon "weassistant/http/api/common"
	"weassistant/http/api/v1/common"
)

// BaseAPI 前台基本控制器
type BaseAPI struct {
	common.BaseAPI
}

// userID 获取当前用户的userID
func (api *BaseAPI) userID(ctx iris.Context) (userID uint64) {
	userID, err := ctx.Values().GetUint64("user_id")
	if err != nil {
		// 正常情况下needlogin模块应当已经将userID写入了memstore中，不应当发生此错误
		panic(api.Error(ctx, apiCommon.RetCodeUnknownError, ctx.Translate("ServerError"), err, "userID miss after needlogin middleware"))
	}
	return
}
