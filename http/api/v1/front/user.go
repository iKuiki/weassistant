package front

import (
	"github.com/dgrijalva/jwt-go"
	jwtmiddleware "github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
)

// UserAPI 用户控制器
type UserAPI struct {
	BaseAPI
}

// Get 获取用户信息
func (api *UserAPI) Get(ctx iris.Context, jwtHandler *jwtmiddleware.Middleware) mvc.Result {
	userToken := jwtHandler.Get(ctx)
	var nickname string
	if claims, ok := userToken.Claims.(jwt.MapClaims); ok && userToken.Valid {
		nickname = claims["user_nickname"].(string)
	} else {
		nickname = "Claims Failed"
	}
	return api.Output("hello " + nickname)
}
