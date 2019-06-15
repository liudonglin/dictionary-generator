package handler

import (
	"net/http"

	"github.com/labstack/echo"
)

// BusinessError 业务错误
type BusinessError struct {
	Message  string
	Internal error // Stores the error returned by an external dependency
}

// Error makes it compatible with `error` interface.
func (he *BusinessError) Error() string {
	return he.Message
}

func customHTTPErrorHandler(err error, c echo.Context) {
	result := StandardResult{}
	if _, ok := err.(*BusinessError); ok {
		result.Code = 300
		result.Message = err.Error()
	} else if err1, ok := err.(*echo.HTTPError); ok {
		// 400 未带token
		// 401 无效token
		result.Code = err1.Code
		result.Message = err1.Message.(string)
	} else {
		result.Code = 500
		result.Message = "服务端发生错误,请稍后再试!"
	}
	c.JSON(http.StatusOK, result)
}

// StandardResult 统一结果返回
type StandardResult struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// PageResult ...
type PageResult struct {
	Total int         `json:"total"`
	List  interface{} `json:"list"`
}
