package manager

import (
	"github.com/jinzhu/gorm"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"weassistant/http/api/common"
	"weassistant/http/api/v1/common/retcode"
	"weassistant/models"
	"weassistant/services/locker"
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

// Post 创建用户
func (api *UserAPI) Post(ctx iris.Context,
	userService orm.UserService,
	registerLockerService locker.CommonLockerService) mvc.Result {
	// 此操作与用户注册几乎一个逻辑
	var formUser, user models.User
	api.ReadForm(ctx, &formUser)
	formUser.CreateTo(&user)
	api.Valid(ctx, user)
	// 锁此操作
	accountLocker := registerLockerService.ObtainLock(user.Account)
	accountLocker.Lock()
	defer accountLocker.Unlock()
	whereOptions := []orm.WhereOption{
		orm.WhereOption{Query: "account = ?", Item: []interface{}{user.Account}},
	}
	_, err := userService.GetByWhereOptions(whereOptions)
	if err != gorm.ErrRecordNotFound {
		// 用户已存在
		if err == nil {
			return api.InvalidParam(ctx, retcode.RetCodeAccountDuplicate, ctx.Translate("AccountAlreadyExist"), nil, "")
		}
		// 查询错误
		return api.Error(ctx, common.RetCodeGormQueryFail, ctx.Translate("CheckAccountExistFail"), err, "userService.GetByWhereOptions error: "+err.Error())
	}
	err = userService.Save(&user)
	if err != nil {
		return api.Error(ctx, common.RetCodeGormQueryFail, ctx.Translate("CreateUserFail"), err, "userService.Save error: "+err.Error())
	}
	return api.Success(ctx.Translate("CreateSuccess"), user)
}
