package manager_test

import (
	"github.com/kataras/iris/httptest"
	"github.com/tidwall/gjson"
	"testing"
	"time"
)

// TestModifyMyAdminInfo 测试管理员修改自己的信息
// 设计：
// - 当未填入修改密码的选项时，密码应该保持不变
// - 当填入修改密码的选项时，应当顺便注销其他的session
func TestModifyMyAdminInfo(t *testing.T) {
	e := httptest.New(t, testApp)

	// 先获得一个测试admin
	admin, err := getTestAdministrator()
	if err != nil {
		t.Fatalf("getTestAdministrator fail: %v", err)
	}
	t.Log("admin got: ", admin)
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
		t.Log("jwtToken: ", jwtToken)
	}
	{
		// 尝试获取信息
		e.GET("/api/v1/manager/my").WithHeader("Authorization", "bearermgr "+jwtToken).
			Expect().Status(httptest.StatusOK).
			Body().Contains(`"code": 0`).
			Contains(`"name": "` + admin.Name + `"`)
	}
	{
		// 尝试不带原密码修改密码
		e.PATCH("/api/v1/manager/my").WithHeader("Authorization", "bearermgr "+jwtToken).
			WithFormField("name", "testModifyName").
			WithFormField("password", "testModifyName").
			Expect().Status(httptest.StatusOK).
			Body().Contains(`"code": 12`).
			Contains(`"msg": "old password incorrect"`)
	}
	// 制造多一份登陆token
	// 获取到测试管理员，开始测试
	var jwtToken2 string
	{
		// 尝试登陆
		body := e.POST("/api/v1/manager/auth/login").
			WithFormField("account", admin.Account).
			WithFormField("password", admin.Password).Expect().Status(httptest.StatusOK).
			Body().Contains(`"code": 0`).
			Contains(`"login successful"`).Contains(`"token": "`).Contains(`"administrator": `).Raw()
		jwtToken2 = gjson.Get(body, "data.token").String()
		t.Log("jwtToken2: ", jwtToken2)
	}
	{
		testName := "testModifyName_" + RandStringBytes(4)
		// 不修改密码，只修改名字后，名字应当两个token都可以登陆
		e.PATCH("/api/v1/manager/my").WithHeader("Authorization", "bearermgr "+jwtToken).
			WithFormField("name", testName).
			Expect().Status(httptest.StatusOK).
			Body().Contains(`"code": 0`)
		// 名字已经修改
		admin.Name = testName
		t.Log("admin.Name change: ", admin.Name)
		// 尝试使用token1获取信息
		e.GET("/api/v1/manager/my").WithHeader("Authorization", "bearermgr "+jwtToken).
			Expect().Status(httptest.StatusOK).
			Body().Contains(`"code": 0`).
			Contains(`"name": "` + admin.Name + `"`)
		// 尝试使用token2获取信息
		e.GET("/api/v1/manager/my").WithHeader("Authorization", "bearermgr "+jwtToken2).
			Expect().Status(httptest.StatusOK).
			Body().Contains(`"code": 0`).
			Contains(`"name": "` + admin.Name + `"`)
	}
	{
		testName := "testModifyName_" + RandStringBytes(4)
		testPasswd := "testModifyName_" + RandStringBytes(4)
		// 修改密码后，token2应当不可登陆
		e.PATCH("/api/v1/manager/my").WithHeader("Authorization", "bearermgr "+jwtToken).
			WithFormField("name", testName).
			WithFormField("password", testPasswd).
			WithFormField("account", admin.Password). // account字段用于旧密码
			Expect().Status(httptest.StatusOK).
			Body().Contains(`"code": 0`)
		// 名字已经修改
		admin.Name = testName
		t.Log("admin.Name change: ", admin.Name)
		// 密码已修改
		admin.Password = testPasswd
		t.Log("admin.Password change: ", admin.Password)
		time.Sleep(time.Second) // 休息1秒等待登出生效
		// 尝试使用token1获取信息
		e.GET("/api/v1/manager/my").WithHeader("Authorization", "bearermgr "+jwtToken).
			Expect().Status(httptest.StatusOK).
			Body().Contains(`"code": 0`).
			Contains(`"name": "` + admin.Name + `"`)
		// 尝试使用token2获取信息
		e.GET("/api/v1/manager/my").WithHeader("Authorization", "bearermgr "+jwtToken2).
			Expect().Status(httptest.StatusUnauthorized).
			Body().Contains(`"code": 2`).
			Contains(`"msg": "need login"`)
	}
}
