package routes

import (
	"dumbmerch/handlers"
	"dumbmerch/pkg/middleware"
	"dumbmerch/pkg/mysql"
	"dumbmerch/repositories"

	"github.com/labstack/echo/v4"
)

func ProductRoutes(e *echo.Group) {
	productRepository := repositories.RepositoryProduct(mysql.DB)
	h := handlers.HandlerProduct(productRepository)

	e.GET("/products", middleware.Auth(h.FindProducts))
	e.GET("/product/:id", middleware.Auth(h.GetProduct))
	e.POST("/product", middleware.Auth(h.CreateProduct))
}
