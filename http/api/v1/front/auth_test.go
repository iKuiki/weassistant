package front_test

import (
	"github.com/kataras/iris/httptest"
	"github.com/tidwall/gjson"
	"testing"
)

func TestRegister(t *testing.T) {
	e := httptest.New(t, testApp)
	{
		// 验证提交为空的情况
		e.POST("/api/v1/user/register").Expect().Status(httptest.StatusOK).
			Body().
			Contains(`"code": 11`).
			Contains(`"Account is a required field"`).
			Contains(`"Nickname is a required field"`).
			Contains(`"Password must be at least 6 characters in length"`)
	}
	{
		// 测试重复用户名
		e.POST("/api/v1/user/register").
			WithFormField("nickname", "testName").
			WithFormField("account", "testExistAccount").
			WithFormField("password", "123321").Expect().Status(httptest.StatusOK)
		e.POST("/api/v1/user/register").
			WithFormField("nickname", "testName").
			WithFormField("account", "testExistAccount").
			WithFormField("password", "123321").Expect().Status(httptest.StatusOK).
			Body().Contains(`"code": 101`).
			Contains(`"account already exist"`)
	}
	{
		// 验证提交的昵称、用户名、密码长度不足
		e.POST("/api/v1/user/register").
			WithFormField("nickname", "t").
			WithFormField("account", "ttt").
			WithFormField("password", "ttttt").Expect().Status(httptest.StatusOK).
			Body().Contains(`"code": 11`).
			Contains(`"Nickname must be at least 2 characters in length"`).
			Contains(`"Account must be at least 4 characters in length"`).
			Contains(`"Password must be at least 6 characters in length"`)
	}
	{
		// 验证提交的昵称、用户名、密码长度过长
		e.POST("/api/v1/user/register").
			WithFormField("nickname", "ttttttttttttttttttttt").
			WithFormField("account", "ttttttttttttttttttttt").
			WithFormField("password", "ttttttttttttttttttttttttttttttttttttttttt").Expect().Status(httptest.StatusOK).
			Body().Contains(`"code": 11`).
			Contains(`"Nickname must be a maximum of 20 characters in length"`).
			Contains(`"Account must be a maximum of 20 characters in length"`).
			Contains(`"Password must be a maximum of 40 characters in length"`)
	}
	{
		// 验证正常提交时可以注册
		e.POST("/api/v1/user/register").
			WithFormField("nickname", "testRegister").
			WithFormField("account", "testReg_"+RandStringBytes(10)).
			WithFormField("password", "123321").Expect().Status(httptest.StatusOK).
			Body().Contains(`"code": 0`).
			Contains(`"register successful"`).
			Contains(`"nickname": "testRegister"`).
			Contains(`"last_login_at": null`)
	}
}

func TestRegisterAndLogin(t *testing.T) {
	e := httptest.New(t, testApp)
	account := "testR&L_" + RandStringBytes(10)
	password := "123321"
	{
		// 先注册
		e.POST("/api/v1/user/register").
			WithFormField("nickname", "testRegisterAndLogin").
			WithFormField("account", account).
			WithFormField("password", password).Expect().Status(httptest.StatusOK).
			Body().Contains(`"code": 0`).
			Contains(`"register successful"`).
			Contains(`"nickname": "testRegisterAndLogin"`).
			Contains(`"last_login_at": null`)
	}
	var jwtToken string
	{
		// 尝试登陆
		body := e.POST("/api/v1/user/login").
			WithFormField("account", account).
			WithFormField("password", password).Expect().Status(httptest.StatusOK).
			Body().Contains(`"code": 0`).
			Contains(`"login successful"`).Raw()
		jwtToken = gjson.Get(body, "data").String()
	}
	{
		// 尝试获取信息
		e.GET("/api/v1/user").WithHeader("Authorization", "Bearer "+jwtToken).
			Expect().Status(httptest.StatusOK).
			Body().Contains(`"code": 0`).
			Contains(`"hello testRegisterAndLogin"`)
	}
}
