package main

import (
	"fmt"

	"github.com/labstack/echo"
)

func main() {
	e := echo.New()

	fmt.Println("server running localhost:5000")
	e.Logger.Fatal(e.Start("localhost:5000"))
}
