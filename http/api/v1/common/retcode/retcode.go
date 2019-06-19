package retcode

import (
	"weassistant/http/api/common"
)

// 输入错误使用正数，系统错误使用负数
const (

	// ------------ 系统错误 ----------------

	// ------------ Scene -------------

	// ------------ 输入参数错误 -------------

	// RetCodeAccountDuplicate 用户名重复
	RetCodeAccountDuplicate common.RespCode = 101
	// RetCodeLoginInfoIncorrect 登录信息错误
	RetCodeLoginInfoIncorrect common.RespCode = 102
)
