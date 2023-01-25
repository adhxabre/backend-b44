package main

import (
	"fmt"
	"fundamental-golang-result-new/routes"

	"github.com/labstack/echo"
)

func main() {
	e := echo.New()

	routes.RouteInit(e.Group("/api/v1"))

	fmt.Println("server running localhost:5000")
	e.Logger.Fatal(e.Start("localhost:5000"))
}
