package front

import (
	"weassistant/http/api/common"
	"weassistant/http/api/v1/common/retcode"
	"weassistant/models"
	"weassistant/services/locker"
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

// PostRegister 注册
func (api *AuthAPI) PostRegister(ctx iris.Context,
	userService orm.UserService,
	registerLockerService locker.CommonLockerService) mvc.Result {
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
	return api.Success(ctx.Translate("RegisterSuccess"), user)
}

// PostLogin 登陆
func (api *AuthAPI) PostLogin(ctx iris.Context, jwtHandler *jwtmiddleware.Middleware, jwtValidationKey string, userService orm.UserService, sessionService orm.UserSessionService) mvc.Result {
	account := ctx.FormValue("account")
	if account == "" {
		ctx.Application().Logger().Debug("a user trying login but without account")
		return api.InvalidParam(ctx, retcode.RetCodeLoginInfoIncorrect, ctx.Translate("LoginInfoIncorrect"), nil, "")
	}
	ctx.Application().Logger().Debugf("account %s trying login", account)
	whereOptions := []orm.WhereOption{
		orm.WhereOption{Query: "account = ?", Item: []interface{}{account}},
	}
	user, err := userService.GetByWhereOptions(whereOptions)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return api.InvalidParam(ctx, retcode.RetCodeLoginInfoIncorrect, ctx.Translate("LoginInfoIncorrect"), nil, "")
		}
		return api.Error(ctx, common.RetCodeGormQueryFail, ctx.Translate("QueryUserInfoFail"), err, "userService.GetByWhereOptions error: "+err.Error())
	}
	ctx.Application().Logger().Debugf("account %s[%d] found", user.Account, user.ID)
	{
		password := ctx.FormValue("password")
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		if err != nil {
			ctx.Application().Logger().Debugf("user %s login fail because password Incorrect", user.Account)
			return api.InvalidParam(ctx, retcode.RetCodeLoginInfoIncorrect, ctx.Translate("LoginInfoIncorrect"), nil, "")
		}
	}
	ctx.Application().Logger().Debugf("user %s[%d] login correct, making session...", user.Account, user.ID)
	session := models.UserSession{
		UserID:      user.ID,
		Token:       uuid.New().String(),
		Effective:   true,
		LoginMethod: models.LoginMethodAccountPassword,
		LoginIP:     ctx.RemoteAddr(),
	}
	err = sessionService.Save(&session)
	if err != nil {
		return api.Error(ctx, common.RetCodeGormQueryFail, ctx.Translate("SaveSessionFail"), err, "sessionService.Save error :"+err.Error())
	}
	token := jwt.NewWithClaims(jwtHandler.Config.SigningMethod, jwt.MapClaims{
		"user_id":       session.UserID,
		"user_nickname": user.Nickname,
		"user_token":    session.Token,
	})
	tokenString, err := token.SignedString([]byte(jwtValidationKey))
	if err != nil {
		return api.Error(ctx, common.RetCodeJwtSignedFail, ctx.Translate("JwtSignedFail"), err, "token.SignedString error: "+err.Error())
	}
	ctx.Application().Logger().Debugf("user %s[%d] login success", user.Account, user.ID)
	return api.Success(ctx.Translate("LoginSuccess"), tokenString)
}
