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
}

type needLoginMiddleware struct {
	common.BaseController
	jwtHandler     *jwtmiddleware.Middleware
	sessionService orm.AdministratorSessionService
}

// MustNewNeedLoginMiddleware 创建登录态检查中间件
func MustNewNeedLoginMiddleware(jwtHandler *jwtmiddleware.Middleware, sessionService orm.AdministratorSessionService) NeedLoginMiddleware {
	return &needLoginMiddleware{
		jwtHandler:     jwtHandler,
		sessionService: sessionService,
	}
}

func (mid *needLoginMiddleware) Serve(ctx iris.Context) {
	administratorID, err := mid.AdministratorID(ctx)
	if err != nil {
		mid.Error(ctx, common.RetCodeRedisQueryFail, ctx.Translate("ValidSessionFail"), err, "mid.UID error: "+err.Error()).Dispatch(ctx)
		ctx.StopExecution()
		return
	}
	if administratorID == 0 {
		ctx.Application().Logger().Debugf("a request was block owing to needlogin")
		mid.NeedLogin(ctx.Translate("NeedLogin")).Dispatch(ctx)
		ctx.StopExecution()
		return
	}
	ctx.Next()
}

// AdministratorID 根据用户jwtToken获取administratorID
func (mid *needLoginMiddleware) AdministratorID(ctx iris.Context) (administratorID uint64, err error) {
	memstore := ctx.Values()
	inf := memstore.Get("administrator_id")
	if inf == nil {
		e := mid.jwtHandler.CheckJWT(ctx)
		if e != nil {
			ctx.Application().Logger().Debugf("administrator CheckJWT fail: %v", e)
			return
		}
		administratorToken := mid.jwtHandler.Get(ctx)
		claims, ok := administratorToken.Claims.(jwt.MapClaims)
		if ok && administratorToken.Valid {
			token := claims["administrator_token"].(string)
			adminID := uint64(claims["administrator_id"].(float64))
			var effective bool
			effective, err = mid.sessionService.ValidSessionToken(adminID, token)
			if err != nil || !effective {
				ctx.Application().Logger().Debugf("administrator [%d] session valid fail via token %s", adminID, token)
				return
			}
			administratorID = adminID
			ctx.Values().SetImmutable("administrator_id", administratorID)
			ctx.Values().SetImmutable("administrator_token", token)
		}
	} else {
		administratorID = inf.(uint64)
	}
	return
}
