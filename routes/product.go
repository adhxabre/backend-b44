package routes

import (
	"dumbmerch/handlers"
	"dumbmerch/pkg/mysql"
	"dumbmerch/repositories"

	"github.com/labstack/echo"
)

func ProductRoutes(e *echo.Group) {
	productRepository := repositories.RepositoryProduct(mysql.DB)
	h := handlers.HandlerProduct(productRepository)

	e.GET("/products", h.FindProducts)
	e.GET("/product/:id", h.GetProduct)
	e.POST("/product", h.CreateProduct)
}
