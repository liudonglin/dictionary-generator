package main

import (
	"github.com/labstack/echo"
)

func main() {
	e := echo.New()
	e.Static("/ui", "ui-dist")
	e.File("/ui", "ui-dist/index.html")

	e.Start(":8080")
}
