package manager

import (
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"golang.org/x/crypto/bcrypt"
	apiCommon "weassistant/http/api/common"
	"weassistant/models"
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

// Patch 修改个人信息
// 因为修改个人信息属于比较特殊的情况，所以要对修改个人信息做更多处理
// TODO: 修复修改密码的时候需要靠account字段传旧密码的问题，可以考虑新建一个struct组合models.Administrator，扩展一个OldPassword字段
func (api *MyAPI) Patch(ctx iris.Context, administratorService orm.AdministratorService, sessionService orm.AdministratorSessionService) mvc.Result {
	admin := api.admin(ctx, administratorService)
	var formAdministrator models.Administrator
	api.ReadForm(ctx, &formAdministrator)
	// 由于验证规则是针对注册的规则，所以验证前需要针对注册规则适当修订

	// 是否修改了密码，因为验证时密码为必填项
	// 但是我们认为此处如果未填写密码则为未修改
	// 所以为了验证，如果密码为空则我们需要填入一个密码
	passwdHasChange := false
	{
		// 不允许修改Account，
		password := formAdministrator.Account
		formAdministrator.Account = admin.Account
		if formAdministrator.Password == "" {
			formAdministrator.Password = "This is a tmp Passwd."
		} else {
			// 如果需要修改密码，则需要验证旧密码
			{
				err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(password))
				if err != nil {
					ctx.Application().Logger().Debugf("administrator %s modify password fail because old password Incorrect", admin.Account)
					return api.InvalidParam(ctx, apiCommon.RetCodeReadFormFail, ctx.Translate("OldPasswdIncorrect"), nil, "")
				}
			}
			passwdHasChange = true
		}
	}
	api.Valid(ctx, formAdministrator)
	{
		// 验证通过后，如果原来未修改密码，则我们需要把我们预置的密码清除
		if !passwdHasChange {
			formAdministrator.Password = ""
		}
	}
	formAdministrator.UpdateTo(&admin)
	err := administratorService.Save(&admin)
	if err != nil {
		return api.Error(ctx, apiCommon.RetCodeGormQueryFail, ctx.Translate("SaveFail"), err, "")
	}
	// 如果修改了密码，应当清除登陆态？
	if passwdHasChange {
		currentToken := ctx.Values().GetString("administrator_token")
		sessions, err := sessionService.GetListByWhereOptions([]orm.WhereOption{
			orm.WhereOption{Query: "administrator_id = ?", Item: []interface{}{admin.ID}},
			orm.WhereOption{Query: "token != ?", Item: []interface{}{currentToken}},
			orm.WhereOption{Query: "Effective = ?", Item: []interface{}{true}},
		}, []string{}, 0, 0)
		if err != nil {
			return api.Error(ctx, apiCommon.RetCodeGormQueryFail, ctx.Translate("QueryAdministratorInfoFail"), err, "")
		}
		for _, s := range sessions {
			session := s
			session.Effective = false
			go func() {
				defer func() {
					if e := recover(); e != nil {
						ctx.Application().Logger().Error("panic recoverd at v1/manager/my.Patch: ", e)
					}
				}()
				err = sessionService.Save(&session)
				if err != nil {
					api.Error(ctx, apiCommon.RetCodeGormQueryFail, ctx.Translate("LogoutFail"), err, fmt.Sprintf("Destroy session[%s] fail while admin[%s] change passwd", session.Token, admin.Account))
				}
			}()
		}
	}
	return api.Success(ctx.Translate("SaveSuccess"), admin)
}
