package vali

import (
	"QuickPass/pkg/util"
	"github.com/gin-gonic/gin/binding"
	zhongwen "github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
	"strings"
	"time"
)

var v *validator.Validate
var trans ut.Translator

func InitValidation() {
	// 中文翻译
	zh := zhongwen.New()
	uni := ut.New(zh, zh)
	trans, _ = uni.GetTranslator("zh")

	var ok bool
	v, ok = binding.Validator.Engine().(*validator.Validate)
	if ok {
		// 验证器注册翻译器
		zh_translations.RegisterDefaultTranslations(v, trans)
		// 自定义验证方法
		v.RegisterValidation("timeValidated", timeValidated)
		// 注册标签翻译
		v.RegisterTranslation("timeValidated", trans, func(ut ut.Translator) error {
			return ut.Add("timeValidated", "{0}格式错误!必须为[yyyy-MM-dd HH:mm:ss]格式", true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("timeValidated", fe.Field())
			return t
		})

		v.RegisterValidation("dateValidated", dateValidated)
		v.RegisterTranslation("dateValidated", trans, func(ut ut.Translator) error {
			return ut.Add("dateValidated", "{0}格式错误!必须为[yyyy-MM-dd]格式", true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("dateValidated", fe.Field())
			return t
		})

		// 自定义验证方法
		v.RegisterValidation("amountValidated", amountValidated)
		// 注册标签翻译
		v.RegisterTranslation("amountValidated", trans, func(ut ut.Translator) error {
			return ut.Add("amountValidated", "{0}格式错误!必须精确到分", true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("amountValidated", fe.Field())
			return t
		})

		// 卖出金额校验
		v.RegisterValidation("sellAmountValidated", sellAmountValidated)
		// 注册标签翻译
		v.RegisterTranslation("sellAmountValidated", trans, func(ut ut.Translator) error {
			return ut.Add("sellAmountValidated", "{0}格式错误!必须精确到元", true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("sellAmountValidated", fe.Field())
			return t
		})
	}
}

func TranslateOneError(errs validator.ValidationErrors) string {
	for _, e := range errs {
		// can translate each error one at a time.
		return e.Translate(trans)
	}

	return ""
}

func TranslateAllError(errs validator.ValidationErrors) string {
	var errList []string
	for _, e := range errs {
		// can translate each error one at a time.
		errList = append(errList, e.Translate(trans))
	}
	return strings.Join(errList, "|")
}

//时间格式校验
func timeValidated(fl validator.FieldLevel) bool {
	if timeString, ok := fl.Field().Interface().(string); ok {
		if timeString != "" {
			_, err := time.Parse(util.TIME_TEMPLATE_1, timeString)
			return err == nil
		}
	}

	return true
}

//日期格式校验
func dateValidated(fl validator.FieldLevel) bool {
	if timeString, ok := fl.Field().Interface().(string); ok {
		if timeString != "" {
			_, err := time.Parse(util.TIME_TEMPLATE_3, timeString)
			return err == nil
		}
	}

	return true
}

//金额格式校验
func amountValidated(fl validator.FieldLevel) bool {
	if amount, ok := fl.Field().Interface().(int64); ok {
		if amount != 0 {
			remainder := amount % 100
			return remainder == 0
		}
	}

	return true
}

//卖出金额格式校验
func sellAmountValidated(fl validator.FieldLevel) bool {
	if amount, ok := fl.Field().Interface().(int64); ok {
		if amount != 0 {
			remainder := amount % 10000
			return remainder == 0
		}
	}

	return true
}
