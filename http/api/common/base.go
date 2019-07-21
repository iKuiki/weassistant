package common

import (
	"fmt"
	"io/ioutil"
	"strings"
	"weassistant/http/common"

	"github.com/kataras/iris"
	"github.com/tealeg/xlsx"
	"gopkg.in/go-playground/validator.v9"
)

// BaseController 基础控制器，提供通用方法
type BaseController struct {
	common.BaseController
}

// ReadForm 读取表单
func (ctl *BaseController) ReadForm(ctx iris.Context, object interface{}) {
	err := ctx.ReadForm(object)
	if err != nil {
		ctx.Application().Logger().Debug("readform error: ", err)
		if err.Error() == "An empty form passed on ReadForm" {
			panic(ctl.InvalidParam(ctx, RetCodeValidFail, ctx.Translate("EmptyForm"), nil, ""))
		}
		// ctl.InvalidParam(ctx,RetCodeValidFail, "ReadForm Fail", nil, err.Error())
		panic(ctl.InvalidParam(ctx, RetCodeReadFormFail, ctx.Translate("ReadFormFail"), err, ""))
	}
}

// Valid 读取表单并验证数据
func (ctl *BaseController) Valid(ctx iris.Context, object interface{}) {
	err := common.Validate.Struct(object)
	if err != nil {
		if es, ok := err.(validator.ValidationErrors); ok {
			lang := ctx.Values().GetString(
				ctx.Application().ConfigurationReadOnly().GetTranslateLanguageContextKey())
			trans, ok := common.ValidateTransMap[lang]
			if !ok {
				trans = common.DefaultTrans
				ctl.logErrToConsole(ctx, fmt.Sprintf("lang %s validator trans not found, use default trans en-US", lang))
			}
			var esErrs []string
			for _, e := range es {
				esErrs = append(esErrs, e.Translate(trans))
			}
			panic(ctl.InvalidParam(ctx, RetCodeValidFail, strings.Join(esErrs, ";"), es.Translate(trans), ""))
		} else {
			panic(ctl.InvalidParam(ctx, RetCodeValidFail, err.Error(), nil, ""))
		}
	}
}

// ReadValid 读取表单并验证数据
func (ctl *BaseController) ReadValid(ctx iris.Context, object interface{}) {
	ctl.ReadForm(ctx, object)
	ctl.Valid(ctx, object)
}

// ReadXlsxFromForm 从上传文件中读取xlsx
func (ctl *BaseController) ReadXlsxFromForm(ctx iris.Context, field string) (sheet *xlsx.Sheet) {
	upFile, _, err := ctx.FormFile(field)
	if err != nil {
		panic(ctl.Error(ctx, RetCodeFormFileFail, ctx.Translate("UploadFileFail"), err, "FormFile error: "+err.Error()))
	}
	defer upFile.Close()
	body, err := ioutil.ReadAll(upFile)
	if err != nil {
		panic(ctl.Error(ctx, RetCodeUnknownError, ctx.Translate("IOFail"), err, "ioutil.ReadAll error: "+err.Error()))
	}
	excel, err := xlsx.OpenBinary(body)
	if err != nil {
		panic(ctl.Error(ctx, RetCodeXlsxFail, ctx.Translate("OpenXlsxFail"), err, "OpenBinary error: "+err.Error()))
	}
	if len(excel.Sheets) == 0 {
		panic(ctl.InvalidParam(ctx, RetCodeValidFail, ctx.Translate("UploadFileEmpty"), nil, ""))
	}
	sheet = excel.Sheets[0]
	if sheet == nil {
		panic(ctl.Error(ctx, RetCodeUnknownError, ctx.Translate("ServerError"), err, "ReadXlsxFromForm panic: sheet is nil"))
	}
	return
}
