package manager_test

import (
	"github.com/kataras/iris/httptest"
	"github.com/tidwall/gjson"
	"testing"
)

func TestLogin(t *testing.T) {
	e := httptest.New(t, testApp)
	{
		// 未登录时期望获取到needlogin
		e.GET("/api/v1/manager/my").
			Expect().Status(httptest.StatusUnauthorized).
			Body().Contains(`"code": 2`).
			Contains(`"msg": "need login"`)
	}
	admin, err := getTestAdministrator()
	if err != nil {
		t.Fatalf("getTestAdministrator fail: %v", err)
	}
	// 获取到测试管理员，开始测试
	var jwtToken string
	{
		// 尝试登陆
		body := e.POST("/api/v1/manager/auth/login").
			WithFormField("account", admin.Account).
			WithFormField("password", admin.Password).Expect().Status(httptest.StatusOK).
			Body().Contains(`"code": 0`).
			Contains(`"login successful"`).Contains(`"token": "`).Contains(`"administrator": `).Raw()
		jwtToken = gjson.Get(body, "data.token").String()
	}
	{
		// 尝试获取信息
		e.GET("/api/v1/manager/my").WithHeader("Authorization", "bearermgr "+jwtToken).
			Expect().Status(httptest.StatusOK).
			Body().Contains(`"code": 0`).
			Contains(`"name": "testLogin"`)
	}
	{
		// 注销
		e.DELETE("/api/v1/manager/my").WithHeader("Authorization", "bearermgr "+jwtToken).
			Expect().Status(httptest.StatusOK)
	}
	{
		// 注销后应该返回未登录
		e.GET("/api/v1/manager/my").WithHeader("Authorization", "bearermgr "+jwtToken).
			Expect().Status(httptest.StatusUnauthorized)
	}
}
