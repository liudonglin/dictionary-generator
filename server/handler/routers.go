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

	r.POST("/project/load", loadProject)
	r.POST("/project/save", saveProject)
	r.POST("/project/list", listProject)
	r.POST("/project/delete", deleteProject)
	r.POST("/conn/loadpid", loadConnectionsByPID)

	r.POST("/database/save", saveDataBase)
	r.POST("/database/list", listDataBase)
	r.POST("/database/delete", deleteDataBase)

	r.POST("/column/save", saveColumn)
	r.POST("/column/delete", deleteColumn)

	r.POST("/table/save", saveTable)
	r.POST("/table/list", listTable)
	r.POST("/table/listdetail", listTableDetail)
	r.POST("/table/delete", deleteTable)

	r.POST("/dbimport/loaddb", loadDBsByConnInfo)
	r.POST("/dbimport/savedbs", saveConnInfo)

	return e
}
