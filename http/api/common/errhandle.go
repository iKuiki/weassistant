package common

import (
	"fmt"
	"github.com/getsentry/raven-go"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"runtime/debug"
)

// ErrHandler 错误控制器
var ErrHandler = iris.Handler(func(ctx iris.Context) {
	defer func() {
		if err := recover(); err != nil {
			resp, ok := err.(mvc.Response)
			if ok {
				// ctx.StatusCode(iris.StatusInternalServerError)
				resp.Dispatch(ctx)
				ctx.StopExecution()
			} else {
				ctx.StatusCode(iris.StatusInternalServerError)
				ctx.Text("Internal Server Error")
				ctx.Application().Logger().Errorf("[%s]%s Panic: %#v\n%s", ctx.Method(), ctx.Path(), err, debug.Stack())
				ctx.StopExecution()
				if sentryClient != nil {
					// 发送sentry事件
					e, ok := err.(error)
					if !ok {
						e = fmt.Errorf("Server Panic by some mystical power: %#v", e)
					}
					packet := raven.NewPacket(e.Error(),
						raven.NewException(e,
							raven.NewStacktrace(1, 7, []string{"clever"})),
						raven.NewHttp(ctx.Request()))
					sentryClient.Capture(packet, map[string]string{
						"method": ctx.Method(),
						"path":   ctx.Path(),
					})
				}
			}
		}
	}()
	ctx.Next()
})

// NewErrHandler 返回错误容器中间件
func NewErrHandler() iris.Handler {
	return ErrHandler
}
