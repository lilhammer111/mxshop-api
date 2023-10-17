package validator

import (
	"github.com/go-playground/validator/v10"
	"regexp"
)

func ValidateMobile(fl validator.FieldLevel) bool {
	mobile := fl.Field().String()
	// 正则表达式判断是否合法
	ok, _ := regexp.MatchString(`^1[3-9]\d{9}$`, mobile)
	if !ok {
		return false
	}
	return true
}
