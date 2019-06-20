package manager

import (
	"weassistant/http/api/common"
	"weassistant/http/api/v1/common/retcode"
	"weassistant/models"
	"weassistant/services/orm"

	"golang.org/x/crypto/bcrypt"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	jwtmiddleware "github.com/iris-contrib/middleware/jwt"
	"github.com/jinzhu/gorm"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
)

// AuthAPI 登陆验证控制器
type AuthAPI struct {
	BaseAPI
}

// NoneRegister 注册
// 不开放注册
// 如有更新，请一同更新common_test文件下的getTestAdministrator测试方法
func (api *AuthAPI) NoneRegister(ctx iris.Context,
	administratorService orm.AdministratorService) mvc.Result {
	var formAdministrator, administrator models.Administrator
	api.ReadForm(ctx, &formAdministrator)
	formAdministrator.CreateTo(&administrator)
	api.Valid(ctx, administrator)
	// 后台并发量小，无需redis锁
	whereOptions := []orm.WhereOption{
		orm.WhereOption{Query: "account = ?", Item: []interface{}{administrator.Account}},
	}
	_, err := administratorService.GetByWhereOptions(whereOptions)
	if err != gorm.ErrRecordNotFound {
		// 用户已存在
		if err == nil {
			return api.InvalidParam(ctx, retcode.RetCodeAccountDuplicate, ctx.Translate("AccountAlreadyExist"), nil, "")
		}
		// 查询错误
		return api.Error(ctx, common.RetCodeGormQueryFail, ctx.Translate("CheckAccountExistFail"), err, "administratorService.GetByWhereOptions error: "+err.Error())
	}
	err = administratorService.Save(&administrator)
	if err != nil {
		return api.Error(ctx, common.RetCodeGormQueryFail, ctx.Translate("CreateAdministratorFail"), err, "administratorService.Save error: "+err.Error())
	}
	return api.Success(ctx.Translate("RegisterSuccess"), administrator)
}

// PostLogin 登陆
func (api *AuthAPI) PostLogin(ctx iris.Context, jwtHandler *jwtmiddleware.Middleware, jwtValidationKey string, administratorService orm.AdministratorService, sessionService orm.AdministratorSessionService) mvc.Result {
	account := ctx.FormValue("account")
	if account == "" {
		return api.InvalidParam(ctx, retcode.RetCodeLoginInfoIncorrect, ctx.Translate("LoginInfoIncorrect"), nil, "")
	}
	ctx.Application().Logger().Debugf("account %s trying login", account)
	whereOptions := []orm.WhereOption{
		orm.WhereOption{Query: "account = ?", Item: []interface{}{account}},
	}
	administrator, err := administratorService.GetByWhereOptions(whereOptions)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return api.InvalidParam(ctx, retcode.RetCodeLoginInfoIncorrect, ctx.Translate("LoginInfoIncorrect"), nil, "")
		}
		return api.Error(ctx, common.RetCodeGormQueryFail, ctx.Translate("QueryAdministratorInfoFail"), err, "administratorService.GetByWhereOptions error: "+err.Error())
	}
	ctx.Application().Logger().Debugf("account %s[%d] found", administrator.Account, administrator.ID)
	{
		password := ctx.FormValue("password")
		err = bcrypt.CompareHashAndPassword([]byte(administrator.Password), []byte(password))
		if err != nil {
			ctx.Application().Logger().Debugf("administrator %s login fail because password Incorrect", administrator.Account)
			return api.InvalidParam(ctx, retcode.RetCodeLoginInfoIncorrect, ctx.Translate("LoginInfoIncorrect"), nil, "")
		}
	}
	ctx.Application().Logger().Debugf("administrator %s[%d] login correct, making session...", administrator.Account, administrator.ID)
	session := models.AdministratorSession{
		AdministratorID: administrator.ID,
		Token:           uuid.New().String(),
		Effective:       true,
		LoginMethod:     models.LoginMethodAccountPassword,
		LoginIP:         ctx.RemoteAddr(),
	}
	err = sessionService.Save(&session)
	if err != nil {
		return api.Error(ctx, common.RetCodeGormQueryFail, ctx.Translate("SaveSessionFail"), err, "sessionService.Save error :"+err.Error())
	}
	token := jwt.NewWithClaims(jwtHandler.Config.SigningMethod, jwt.MapClaims{
		"administrator_id":    session.AdministratorID,
		"administrator_name":  administrator.Name,
		"administrator_token": session.Token,
	})
	tokenString, err := token.SignedString([]byte(jwtValidationKey))
	if err != nil {
		return api.Error(ctx, common.RetCodeJwtSignedFail, ctx.Translate("JwtSignedFail"), err, "token.SignedString error: "+err.Error())
	}
	ctx.Application().Logger().Debugf("administrator %s[%d] login success", administrator.Account, administrator.ID)
	return api.Success(ctx.Translate("LoginSuccess"), tokenString)
}
