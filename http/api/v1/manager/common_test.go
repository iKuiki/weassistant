package manager_test

import (
	"github.com/jinzhu/gorm"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/i18n"
	"github.com/pkg/errors"
	"math/rand"
	"weassistant/conf"
	apiCommon "weassistant/http/api/common"
	api1Router "weassistant/http/api/v1/router"
	"weassistant/http/common"
	"weassistant/models"
	"weassistant/services/orm"
)

// 获取测试用的简易app
func getNewTestApp() *iris.Application {
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
	testApp   *iris.Application
	config    conf.Config
	extraConf conf.ExtraConfig
)

func init() {
	config = conf.MustNewConfig()
	err := config.Load("../../../../config.json")
	if err != nil {
		panic(err)
	}
	extraConf = conf.MustExtraNewConfig(config)
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

// 获取一个用来测试的admin账号
// 此方法仿照auth/register方法重写，如register方法有改动，此方法也应当更新
func getTestAdministrator() (admin models.Administrator, err error) {
	formAdministrator := models.Administrator{
		Account:  "testR&L_" + RandStringBytes(10),
		Password: "123321",
		Name:     "testLogin",
	}
	var administrator models.Administrator
	formAdministrator.CreateTo(&administrator)
	// 验证
	err = common.Validate.Struct(administrator)
	if err != nil {
		return
	}
	// 后台并发量小，无需redis锁
	whereOptions := []orm.WhereOption{
		orm.WhereOption{Query: "account = ?", Item: []interface{}{administrator.Account}},
	}
	administratorService := extraConf.GetAdministratorService()
	_, err = administratorService.GetByWhereOptions(whereOptions)
	if err != gorm.ErrRecordNotFound {
		// 用户已存在
		if err == nil {
			err = errors.New("Account already exist")
			return
		}
		// 查询错误
		err = errors.Errorf("CheckAccountExistFail administratorService.GetByWhereOptions error: %v", err)
		return
	}
	err = administratorService.Save(&administrator)
	if err != nil {
		err = errors.Errorf("CreateAdministratorFail administratorService.Save error: %v", err)
		return
	}
	return formAdministrator, nil
}
