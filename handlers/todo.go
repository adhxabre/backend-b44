package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Todos struct {
	Id     string `json:"id"`
	Title  string `json:"title"`
	IsDone bool   `isDone:"isDone"`
}

var todos = []Todos{
	{
		Id:     "1",
		Title:  "Cuci tangan",
		IsDone: true,
	},
	{
		Id:     "2",
		Title:  "Jaga jarak",
		IsDone: false,
	},
}

func FindTodos(c echo.Context) error {
	c.Response().Header().Set("Content-Type", "application/json")
	c.Response().WriteHeader(http.StatusOK)
	return json.NewEncoder(c.Response()).Encode(todos)
}

func GetTodo(c echo.Context) error {
	c.Response().Header().Set("Content-Type", "application/json")
	id := c.Param("id")

	var todoData Todos
	var isGetTodo = false

	for _, todo := range todos {
		if id == todo.Id {
			isGetTodo = true
			todoData = todo
		}
	}

	if !isGetTodo {
		c.Response().WriteHeader(http.StatusNotFound)
		return json.NewEncoder(c.Response()).Encode("ID: " + id + " not found")
	}

	c.Response().WriteHeader(http.StatusOK)
	return json.NewEncoder(c.Response()).Encode(todoData)
}

func CreateTodo(c echo.Context) error {
	var data Todos

	json.NewDecoder(c.Request().Body).Decode(&data)

	todos = append(todos, data)

	c.Response().Header().Set("Content-Type", "application/json")
	c.Response().WriteHeader(http.StatusOK)
	return json.NewEncoder(c.Response()).Encode(todos)
}

func UpdateTodo(c echo.Context) error {
	id := c.Param("id")
	var data Todos
	var isGetTodo = false

	json.NewDecoder(c.Request().Body).Decode(&data)

	for idx, todo := range todos {
		if id == todo.Id {
			isGetTodo = true
			todos[idx] = data
		}
	}

	if !isGetTodo {
		c.Response().WriteHeader(http.StatusNotFound)
		return json.NewEncoder(c.Response()).Encode("ID: " + id + " not found")
	}

	c.Response().Header().Set("Content-Type", "application/json")
	c.Response().WriteHeader(http.StatusOK)
	return json.NewEncoder(c.Response()).Encode(todos)
}

func DeleteTodo(c echo.Context) error {
	id := c.Param("id")
	var isGetTodo = false
	var index = 0

	for idx, todo := range todos {
		if id == todo.Id {
			isGetTodo = true
			index = idx
		}
	}

	if !isGetTodo {
		c.Response().WriteHeader(http.StatusNotFound)
		return json.NewEncoder(c.Response()).Encode("ID: " + id + " not found")
	}

	todos = append(todos[:index], todos[index+1:]...)
	c.Response().Header().Set("Content-Type", "application/json")
	c.Response().WriteHeader(http.StatusOK)
	return json.NewEncoder(c.Response()).Encode(todos)
}
