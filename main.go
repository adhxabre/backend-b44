package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		c.Response().WriteHeader(http.StatusOK)
		return c.String(http.StatusOK, "Halo dunia!")
	})

	PORT := "5000"

	fmt.Println("Server is running on at localhost:" + PORT)
	e.Logger.Fatal(e.Start("localhost:" + PORT))
}
