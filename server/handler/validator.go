package handler

import (
	"sync"
	//validator "gopkg.in/go-playground/validator.v9"
)

var (
	once sync.Once
	//validate *validator.Validate
)

// Validator 获取单例
//func Validator() *validator.Validate {
//once.Do(func() {
//validate = validator.New()
//})
//return validate
//}
