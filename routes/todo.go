package routes

import (
	"dumbmerch/handlers"

	"github.com/labstack/echo/v4"
)

func TodoRoutes(e *echo.Group) {
	e.GET("/todos", handlers.FindTodos)
	e.GET("/todo/:id", handlers.GetTodo)
	e.POST("/todo", handlers.CreateTodo)
	e.PATCH("/todo/:id", handlers.UpdateTodo)
	e.DELETE("/todo/:id", handlers.DeleteTodo)
}
