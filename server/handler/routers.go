package handler

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// New ...
func New() *echo.Echo {
	e := echo.New()
	e.HTTPErrorHandler = customHTTPErrorHandler

	// 配置静态文件
	e.Static("/ui", "ui-dist")
	e.File("/ui", "ui-dist/index.html")

	// 配置无需登录状态的请求
	e.POST("/api/login", login)

	// 配置需要登录状态的请求
	config := middleware.JWTConfig{
		Claims:     &jwtUserClaims{},
		SigningKey: []byte(secret),
	}
	r := e.Group("/api")
	r.Use(middleware.JWTWithConfig(config))
	r.GET("/needtoken", func(c echo.Context) error {
		return c.JSON(http.StatusOK, &StandardResult{})
	})

	r.POST("/project/save", saveProject)
	r.POST("/project/list", listProject)
	r.POST("/project/delete", deleteProject)

	return e
}
