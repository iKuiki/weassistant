package common

// RespondData API返回结果的框架格式
type RespondData struct {
	// Code 错误代码，若无错误发生则未0(RetCodeNoError)
	Code RespCode `json:"code"`
	// Info 附加信息
	Msg string `json:"msg"`
	// Data 若有数据返回，则在此对象内
	Data interface{} `json:"data"`
	// Err 如果有出错，处理出错
	Err error `json:"-"`
}

// Assign 将参数一次性赋值到返回数据中
func (r *RespondData) Assign(code RespCode, msg string, data interface{}, err error) {
	r.Code = code
	r.Msg = msg
	r.Data = data
	r.Err = err
}
