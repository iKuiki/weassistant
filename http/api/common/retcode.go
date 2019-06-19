package common

// RespCode 若发生错误，错误的详细代码
type RespCode int32

const (
	// RespCodeNoError 无错误
	RespCodeNoError RespCode = 0
	// RespCodeNormalError 一般错误
	RespCodeNormalError RespCode = -1

	// RespCodeNotFound 目标未找到
	RespCodeNotFound RespCode = 1
	// RespCodeNeedLogin 需要登录
	RespCodeNeedLogin RespCode = 2
	// RespCodePermissionDeny  权限不足
	RespCodePermissionDeny RespCode = 3
)

// 输入错误使用正数，系统错误使用负数
const (

	// ------------ 系统错误 ----------------

	// RetCodeUnknownError 未知错误
	RetCodeUnknownError RespCode = -11
	// RetCodeGormQueryFail Gorm执行异常
	RetCodeGormQueryFail RespCode = -12
	// RetCodeLockFail 锁操作失败
	RetCodeLockFail RespCode = -13
	// RetCodeRedisQueryFail Redis执行异常
	RetCodeRedisQueryFail RespCode = -14
	// RetCodeXlsxFail xlsx操作失败
	RetCodeXlsxFail RespCode = -15
	// 短信验证码服务操作失败
	RetCodeSmgVerifierFail RespCode = -16
	// RetCodeJwtSignedFail jwt签名失败
	RetCodeJwtSignedFail RespCode = -17
	// RetCodeValidJwtTokenFail 验证登录凭证失败
	RetCodeValidJwtTokenFail RespCode = -18
	// RetCodeCreateSessionFail 创建session会话失败
	RetCodeCreateSessionFail RespCode = -19
	// RetCodeOAuthGetAuthorizeToken OAuth创建授权码失败
	RetCodeOAuthGetAuthorizeToken RespCode = -20

	// ------------ 输入参数错误 -------------

	// RetCodeValidFail 数据校验失败
	RetCodeValidFail RespCode = 11
	// RetCodeReadFormFail 读取表单失败
	RetCodeReadFormFail RespCode = 12
	// RetCodeFormFileFail 读取上传文件失败
	RetCodeFormFileFail RespCode = 13
)
