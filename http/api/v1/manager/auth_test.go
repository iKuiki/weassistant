package manager_test

import (
	"github.com/kataras/iris/httptest"
	"github.com/tidwall/gjson"
	"testing"
)

func TestLogin(t *testing.T) {
	admin, err := getTestAdministrator()
	if err != nil {
		t.Fatalf("getTestAdministrator fail: %v", err)
	}
	// 获取到测试管理员，开始测试
	e := httptest.New(t, testApp)
	var jwtToken string
	{
		// 尝试登陆
		body := e.POST("/api/v1/manager/administrator/login").
			WithFormField("account", admin.Account).
			WithFormField("password", admin.Password).Expect().Status(httptest.StatusOK).
			Body().Contains(`"code": 0`).
			Contains(`"login successful"`).Raw()
		jwtToken = gjson.Get(body, "data").String()
	}
	{
		// 尝试获取信息
		e.GET("/api/v1/manager/administrator").WithHeader("Authorization", "bearermgr "+jwtToken).
			Expect().Status(httptest.StatusOK).
			Body().Contains(`"code": 0`).
			Contains(`"hello testLogin"`)
	}
}
