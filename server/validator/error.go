package validator

import (
	"fmt"
)

// ValidateError 字段验证错误
type ValidateError struct {
	FieldErrors []*FieldError
}

// FieldError 字段具体错误信息
type FieldError struct {
	Tag       string
	Field     string
	Value     interface{}
	Error     string
	Namespace string
}

func (e *ValidateError) Error() string {
	str := ""
	for _, f := range e.FieldErrors {
		str += fmt.Sprintln(f.Error) + ";"
	}
	return str
}
