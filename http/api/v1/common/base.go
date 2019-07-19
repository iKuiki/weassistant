package common

import (
	"github.com/kataras/iris"
	"weassistant/http/api/common"
	"weassistant/services/orm"
)

// BaseAPI 基础API
type BaseAPI struct {
	common.BaseController
}

// Sizeable 可以分页的存储服务
type Sizeable interface {
	GetCountByWhereOptions(whereOptions []orm.WhereOption) (count uint64, err error)
}

// Size 获取指定元素与搜索条件下，数据数量
func (api *BaseAPI) Size(ctx iris.Context, serv Sizeable, whereOptions []orm.WhereOption) (size uint64) {
	size, err := serv.GetCountByWhereOptions(whereOptions)
	if err != nil {
		panic(api.Error(ctx, common.RetCodeGormQueryFail, ctx.Translate("SizeableError"), err, "Sizeable.GetCountByWhereOptions error: "+err.Error()))
	}
	return size
}
