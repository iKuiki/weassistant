package common

import (
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"

	"gopkg.in/go-playground/validator.v9"
	en_translations "gopkg.in/go-playground/validator.v9/translations/en"
	zh_translations "gopkg.in/go-playground/validator.v9/translations/zh"
	"strconv"

	"github.com/kataras/iris"
)

// 自动验证的语言相关
var (
	Validate         *validator.Validate
	uni              *ut.UniversalTranslator
	ValidateTransMap map[string]ut.Translator
	DefaultTrans     ut.Translator
)

func init() {
	Validate = validator.New()
	ValidateTransMap = make(map[string]ut.Translator)
	{
		en := en.New()
		uni = ut.New(en, en)
		trans, _ := uni.GetTranslator("en")
		en_translations.RegisterDefaultTranslations(Validate, trans)
		ValidateTransMap["en"] = trans
		ValidateTransMap["en-US"] = trans
		DefaultTrans = trans
	}
	{
		zh := zh.New()
		uni = ut.New(zh, zh)
		trans, _ := uni.GetTranslator("zh")
		zh_translations.RegisterDefaultTranslations(Validate, trans)
		ValidateTransMap["zh"] = trans
		ValidateTransMap["zh-CN"] = trans
	}
}

// BaseController 基础控制器，提供通用方法
type BaseController struct {
}

// ObtainLimitOffset 获取请求中的size与page参数，并转化为sql的limit与offset
func (ctl *BaseController) ObtainLimitOffset(ctx iris.Context, convertToOffset bool) (limit int64, offset int64) {
	var err error
	var page int64
	limit, err = strconv.ParseInt(string(ctx.FormValue("pageSize")), 10, 64)
	if err != nil || limit < 1 {
		limit = 10
	}
	page, err = strconv.ParseInt(string(ctx.FormValue("pageNo")), 10, 64)
	if err != nil || page < 1 {
		page = 1
	}
	if !convertToOffset {
		return limit, page
	}
	offset = PageToOffset(page, limit)
	return
}

// PageToOffset 将page与limit转化为offset
func PageToOffset(page, limit int64) (offset int64) {
	// calculate page
	if page < 1 {
		offset = 0
	} else {
		offset = (page - 1) * limit
	}
	return
}
