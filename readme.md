# Routing

Routing in `labstack/echo/v4` is pretty straightforward. Make sure you already understand the core concept of REST API like using GET/POST/PATCH/DELETE etc.

`API` or `Application Programming Interface` is an interface that can connect one application with another application. Thus, the API acts as an intermediary between different applications, either within the same platform or across platforms.

> File: `main.go`

- package name

  ```go
  package main
  ```

- import package

 	 ```go
	import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	)
  	```

- Struct

  	```go
	type Todos struct {
		Id     string `json:"id"`
		Title  string `json:"title"`
		IsDone bool   `isDone:"isDone"`
	}
  	```

- Dummy Data

  	```go
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

- Main function for Declare Route

  	```go
	func main() {
		e := echo.New()

		e.GET("/todos", FindTodos)
		e.GET("/todo/:id", GetTodo)
		e.POST("/todo", CreateTodo)
		e.PATCH("/todo/:id", UpdateTodo)
		e.DELETE("/todo/:id", DeleteTodo)

		fmt.Println("server running localhost:5000")
		e.Logger.Fatal(e.Start("localhost:5000"))
	}
 	 ```

- Get all Todo data

  	```go
	func FindTodos(c echo.Context) error {
		c.Response().Header().Set("Content-Type", "application/json")
		c.Response().WriteHeader(http.StatusOK)
		return json.NewEncoder(c.Response()).Encode(todos)
	}
  	```

- Get Todo data by Id

  	```go
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

- Create Todo Data

  	```go
	func CreateTodo(c echo.Context) error {
		var data Todos

		json.NewDecoder(c.Request().Body).Decode(&data)

		todos = append(todos, data)

		c.Response().Header().Set("Content-Type", "application/json")
		c.Response().WriteHeader(http.StatusOK)
		return json.NewEncoder(c.Response()).Encode(todos)
	}
 	```

- Update Todo Data

  	```go
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

- Delete Todo data

  	```go
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
