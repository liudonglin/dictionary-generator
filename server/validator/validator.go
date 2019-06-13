package validator

import (
	zhongwen "github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"
	zh_translations "gopkg.in/go-playground/validator.v9/translations/zh"
)

var (
	validate *validator.Validate
	trans    ut.Translator
)

// ValidateStruct 验证结构体
func ValidateStruct(s interface{}) *ValidateError {
	err := validate.Struct(s)
	if err != nil {
		errs, _ := err.(validator.ValidationErrors)
		verr := &ValidateError{}

		for _, e := range errs {
			field := &FieldError{}
			field.Error = e.Translate(trans)
			field.Field = e.Field()
			field.Value = e.Value()
			field.Namespace = e.Namespace()
			field.Tag = e.Tag()
			verr.FieldErrors = append(verr.FieldErrors, field)
		}
		return verr
	}
	return nil
}

func init() {
	zh := zhongwen.New()
	uni := ut.New(zh, zh)
	trans, _ = uni.GetTranslator("zh")
	validate = validator.New()
	zh_translations.RegisterDefaultTranslations(validate, trans)
}
