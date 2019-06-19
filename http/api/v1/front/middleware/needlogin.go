package middleware

import (
	"github.com/dgrijalva/jwt-go"
	jwtmiddleware "github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris"
	"weassistant/http/api/common"
	"weassistant/services/orm"
)

// NeedLoginMiddleware 检查登录状态的中间件
type NeedLoginMiddleware interface {
	Serve(iris.Context)
	UserID(ctx iris.Context) (userID uint64, err error)
}

type needLoginMiddleware struct {
	common.BaseController
	jwtHandler     *jwtmiddleware.Middleware
	sessionService orm.UserSessionService
}

// MustNewNeedLoginMiddleware 创建登录态检查中间件
func MustNewNeedLoginMiddleware(jwtHandler *jwtmiddleware.Middleware, sessionService orm.UserSessionService) NeedLoginMiddleware {
	return &needLoginMiddleware{
		jwtHandler:     jwtHandler,
		sessionService: sessionService,
	}
}

func (mid *needLoginMiddleware) Serve(ctx iris.Context) {
	userID, err := mid.UserID(ctx)
	if err != nil {
		mid.Error(ctx, common.RetCodeRedisQueryFail, ctx.Translate("ValidSessionFail"), err, "mid.UserID error: "+err.Error()).Dispatch(ctx)
		ctx.StopExecution()
		return
	}
	if userID == 0 {
		ctx.Application().Logger().Debugf("a request was block owing to needlogin")
		mid.NeedLogin(ctx.Translate("NeedLogin")).Dispatch(ctx)
		ctx.StopExecution()
		return
	}
	ctx.Next()
}

// UserID 根据用户jwtToken获取userID
func (mid *needLoginMiddleware) UserID(ctx iris.Context) (userID uint64, err error) {
	memstore := ctx.Values()
	inf := memstore.Get("user_id")
	if inf == nil {
		e := mid.jwtHandler.CheckJWT(ctx)
		if e != nil {
			ctx.Application().Logger().Debugf("user CheckJWT fail: %v", e)
			return
		}
		userToken := mid.jwtHandler.Get(ctx)
		claims, ok := userToken.Claims.(jwt.MapClaims)
		if ok && userToken.Valid {
			token := claims["user_token"].(string)
			uid := uint64(claims["user_id"].(float64))
			var effective bool
			effective, err = mid.sessionService.ValidSessionToken(uid, token)
			if err != nil || !effective {
				ctx.Application().Logger().Debugf("user [%d] session valid fail via token %s", uid, token)
				return
			}
			userID = uid
			ctx.Values().SetImmutable("user_id", userID)
		}
	} else {
		userID = inf.(uint64)
	}
	return
}
