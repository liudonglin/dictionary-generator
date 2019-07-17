package main

import (
	"dg-server/handler"
	"dg-server/tpl"
	"fmt"
)

func main() {
	port := 8080
	e := handler.New()

	tpl.TestGetTableScript(1)
	fmt.Printf("Access address http://localhost:%d/ui/#/dashboard \n", port)
	e.Start(fmt.Sprintf(":%d", port))
}
