package manager

import (
	"github.com/dgrijalva/jwt-go"
	jwtmiddleware "github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
)

// MyAPI 用户控制器
type MyAPI struct {
	BaseAPI
}

// Get 获取用户信息
func (api *MyAPI) Get(ctx iris.Context, jwtHandler *jwtmiddleware.Middleware) mvc.Result {
	administratorToken := jwtHandler.Get(ctx)
	var name string
	if claims, ok := administratorToken.Claims.(jwt.MapClaims); ok && administratorToken.Valid {
		name = claims["administrator_name"].(string)
	} else {
		name = "Claims Failed"
	}
	return api.Output("hello " + name)
}
