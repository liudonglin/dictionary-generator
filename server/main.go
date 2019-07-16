package main

import (
	"dg-server/handler"
	"fmt"
)

func main() {
	port := 8080
	e := handler.New()

	fmt.Printf("Access address http://localhost:%d/ui/#/dashboard \n", port)
	e.Start(fmt.Sprintf(":%d", port))
}
