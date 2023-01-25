package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
)

func main() {
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		c.Response().WriteHeader(http.StatusOK)
		return c.String(http.StatusOK, "Hello World")
	})

	fmt.Println("Server running on localhost:5000")
	e.Logger.Fatal(e.Start("localhost:5000"))
}
