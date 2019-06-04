package handler

import (
	"net/http"

	"../store/base/db"
	"../store/user"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// New ...
func New(db *db.DB) *echo.Echo {
	e := echo.New()

	// 配置静态文件
	e.Static("/ui", "ui-dist")
	e.File("/ui", "ui-dist/index.html")

	userStore := user.New(db)
	hlogin := &loginHandler{userStore}
	// 配置无需登录状态的请求
	e.POST("/api/login", hlogin.login)

	// 配置需要登录状态的请求
	config := middleware.JWTConfig{
		Claims:     &jwtUserClaims{},
		SigningKey: []byte(secret),
	}
	r := e.Group("/api")
	r.Use(middleware.JWTWithConfig(config))
	r.GET("/test", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	return e
}
