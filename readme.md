# Group routes

Group Routes are needed in API development to differentiate a route for API or for standard website link.

- Create `routes` folder source inside it have `routes.go` and `todo.go` file

- Create `handlers` folder source inside it have `todo.go` file

---

- On `routes/routes.go` file, declare Grouping Function for all Route

  ```go
	package routes

	import "github.com/labstack/echo"

	func RouteInit(e *echo.Group) {
		TodoRoutes(e)
	}
  ```

- On `routes/todo.go` file, declare route and handler

  > File: `routes/todo.go`

  	```go
	package routes

	import (
		"dumbmerch/handlers"

		"github.com/labstack/echo"
	)

	func TodoRoutes(e *echo.Group) {
		e.GET("/todos", handlers.FindTodos)
		e.GET("/todo/:id", handlers.GetTodo)
		e.POST("/todo", handlers.CreateTodo)
		e.PATCH("/todo/:id", handlers.UpdateTodo)
		e.DELETE("/todo/:id", handlers.DeleteTodo)
	}
  	```

- On `handlers/todo.go` file, declare `struct`, `dummy data`, and the handlers function

  	```go
	package handlers

	import (
		"encoding/json"
		"net/http"

		"github.com/labstack/echo"
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
 	```

  	```go
  	// Get All Todo
	func FindTodos(c echo.Context) error {
		c.Response().Header().Set("Content-Type", "application/json")
		c.Response().WriteHeader(http.StatusOK)
		return json.NewEncoder(c.Response()).Encode(todos)
	}
  	```

  	```go
  	// Get Todo by Id
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
  	```

  	```go
  	// Create Todo
	func CreateTodo(c echo.Context) error {
		var data Todos

		json.NewDecoder(c.Request().Body).Decode(&data)

		todos = append(todos, data)

		c.Response().Header().Set("Content-Type", "application/json")
		c.Response().WriteHeader(http.StatusOK)
		return json.NewEncoder(c.Response()).Encode(todos)
	}
  	```

  	```go
 	// Update Todo
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
  	```

  	```go
  	// Delete Todo
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
  	```
  
- On `main.go` it become like this (we add e.Group("api/v1") to make versioning of our API) : 

  	```go
	package main

	import (
		"fmt"
		"dumbmerch/routes"

		"github.com/labstack/echo"
	)

	func main() {
		e := echo.New()

		routes.RouteInit(e.Group("/api/v1"))

		fmt.Println("server running localhost:5000")
		e.Logger.Fatal(e.Start("localhost:5000"))
	}
  	```
