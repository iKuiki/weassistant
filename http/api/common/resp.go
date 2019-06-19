package common

import (
	"errors"
	"github.com/getsentry/raven-go"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"strconv"
)

var (
	sentryClient *raven.Client
)

// SetSentryClient 设置SentryClient
func SetSentryClient(client *raven.Client) {
	if sentryClient != nil {
		err := errors.New("sentry client duplicate definition")
		client.CaptureErrorAndWait(err, nil)
		panic(err)
	}
	sentryClient = client
}

// Output 输出json数据
func (ctl *BaseController) Output(data interface{}) mvc.Response {
	resp := RespondData{}
	resp.Assign(RespCodeNoError, "", data, nil)
	return mvc.Response{
		Object: resp,
	}
}

// Success 执行操作成功时调用，含info
func (ctl *BaseController) Success(info string, data interface{}) mvc.Response {
	resp := RespondData{}
	resp.Assign(RespCodeNoError, info, data, nil)
	return mvc.Response{
		Object: resp,
	}
}

func (ctl *BaseController) logErrToConsole(ctx iris.Context, logInfo string) {
	if logInfo != "" {
		ctx.Application().Logger().Warnf("[%s]%s | %s", ctx.Method(), ctx.Path(), logInfo)
	}
}

// NeedLogin 需要登录
func (ctl *BaseController) NeedLogin(info string) mvc.Response {
	resp := RespondData{}
	resp.Assign(RespCodeNeedLogin, info, nil, nil)
	return mvc.Response{
		Object: resp,
		Code:   iris.StatusUnauthorized,
	}
}

// InvalidParam 请求错误，可以填写参数，但返回仍未200
func (ctl *BaseController) InvalidParam(ctx iris.Context, code RespCode, info string, data interface{}, logInfo string) mvc.Response {
	resp := RespondData{}
	resp.Assign(code, info, data, nil)
	ctl.logErrToConsole(ctx, logInfo)
	return mvc.Response{
		Object: resp,
	}
}

// Error 出错时返回错误
func (ctl *BaseController) Error(ctx iris.Context, code RespCode, info string, err error, logInfo string) mvc.Response {
	// 发送sentry事件
	if sentryClient != nil {
		e := err
		if e == nil {
			e = errors.New("Server error by some mystical power")
		}
		packet := raven.NewPacket(e.Error(),
			raven.NewException(e,
				raven.NewStacktrace(1, 7, []string{"clever"})),
			raven.NewHttp(ctx.Request()))
		sentryClient.Capture(packet, map[string]string{
			"method":   ctx.Method(),
			"path":     ctx.Path(),
			"code":     strconv.FormatInt(int64(code), 10),
			"info":     info,
			"log_info": logInfo,
		})
	}
	// 构造返回
	resp := RespondData{}
	resp.Assign(code, info, nil, err)
	ctl.logErrToConsole(ctx, logInfo)
	return mvc.Response{
		Object: resp,
	}
}

// NotFound 找不到请求的对象时使用
func (ctl *BaseController) NotFound(info string, data interface{}) mvc.Response {
	resp := RespondData{}
	resp.Assign(RespCodeNormalError, info, data, nil)
	return mvc.Response{
		Object: resp,
		Code:   iris.StatusNotFound,
	}
}
