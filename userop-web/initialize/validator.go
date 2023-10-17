package initialize

import (
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
	"mxshop-api/userop-web/global"
	"reflect"
	"strings"
)

func Translation(locale string) error {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(field reflect.StructField) string {
			f := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
			if f == "-" {
				return ""
			}
			return f
		})

		zhT := zh.New()
		enT := en.New()
		uni := ut.New(enT, zhT)
		global.Trans, ok = uni.GetTranslator(locale)
		if !ok {
			return fmt.Errorf("uni.GetTranslator(%s)", locale)
		}
		switch locale {
		case "en":
			_ = enTranslations.RegisterDefaultTranslations(v, global.Trans)

		case "zh":
			_ = zhTranslations.RegisterDefaultTranslations(v, global.Trans)

		default:
			_ = enTranslations.RegisterDefaultTranslations(v, global.Trans)
		}
		return nil
	}
	return nil
}
