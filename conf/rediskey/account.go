package rediskey

// 用户登录相关
const (
	// UserTokenSet 用户有效token库
	// type: set
	// %s target: user.ID
	// value: session.Token
	UserTokenSet = "user_%d_tokens"
)

// 后台账号登录相关
const (
	// AdministratorTokenSet 后台账号有效token库
	// type: set
	// %s target: administrator.ID
	// value: session.Token
	AdministratorTokenSet = "administrator_%d_tokens"
)
