package routes

import "github.com/labstack/echo"

func RouteInit(e *echo.Group) {
	UserRoutes(e)
}
