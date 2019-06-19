package front_test

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/i18n"
	"math/rand"
	"weassistant/conf"
	apiCommon "weassistant/http/api/common"
	api1Router "weassistant/http/api/v1/router"
)

// 获取测试用的简易app
func getNewTestApp() *iris.Application {
	config := conf.MustNewConfig()
	err := config.Load("../../../../config.json")
	if err != nil {
		panic(err)
	}
	extraConf := conf.MustExtraNewConfig(config)
	app := iris.New()
	// 错误容器不可少，api内部往往有直接panic来报错的部分，如果没有错误容器将无法正常解析
	app.Use(apiCommon.ErrHandler)
	// 准备i18n
	app.Use(i18n.New(i18n.Config{
		Default:      "en",
		URLParameter: "lang",
		Languages: map[string]string{
			"en": "../../../locales/en-US.ini",
		},
	}))
	api1Router.RegisterAPI1Router(app, extraConf)
	return app
}

var (
	testApp *iris.Application
)

func init() {
	testApp = getNewTestApp()
}

// 用于产生随机字符串的函数
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
