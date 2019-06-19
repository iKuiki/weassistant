package rediskey

// 锁相关
const (
	// UserRegisterLocker 用户注册锁
	// type: *Redis Locker
	// %v target: user.Account
	UserRegisterLocker = "user_register_account_%v_locker"
)
