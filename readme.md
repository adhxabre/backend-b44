> This section we will using Database to store and manage user data (Create, Read, Update, Delete)

---

## Table of contents

- [Prepare](#prepare)
  - [Installation](#installation)
  - [Database](#database)
  - [Models](#models)
  - [Auto Migrate](#auto-migrate)
  - [Connection](#Connection)
  - [Data Transfer Object (DTO)](#data-transfer-object-dto)
- [Fetching Query with Gorm](#fetching-query-with-gorm)
  - [Repositories](#repositories)
  - [Handlers](#handlers)
  - [Routes](#routes)
  - [Root file main.go](#root-file-maingo)

# Prepare

### Installation:

- Gorm

  ```bash
  go get -u gorm.io/gorm
  ```

- MySql

  ```bash
  go get -u gorm.io/driver/mysql
  ```

- Validator

  ```go
  go get github.com/go-playground/validator/v10
  ```

### Database

- Create database named `dumbmerch`

### Models

- Create `models` folder, inside it Create `user.go` file, and write below code

  > File: `models/user.go`

  ```go
  package models

  import "time"

  // User model struct
  type User struct {
    ID          int			`json:"id"`
    Name 		    string		`json:"name" gorm:"type: varchar(255)"`
    Email		    string 		`json:"email" gorm:"type: varchar(255)"`
    Password 	  string		`json:"password" gorm:"type: varchar(255)"`
    CreatedAt 	time.Time	`json:"created_at"`
    UpdatedAt 	time.Time	`json:"updated_at"`
  }
  ```

### Auto Migrate

- Create `database` folder, inside it Create `migration.go` file, and write below code

  > File: `database/migration.go`

  ```go
  package database

  import (
    "dumbmerch/models"
    "dumbmerch/pkg/mysql"
    "fmt"
  )

  // Automatic Migration if Running App
  func RunMigration() {
    err := mysql.DB.AutoMigrate(&models.User{})

    if err != nil {
      fmt.Println(err)
      panic("Migration Failed")
    }

    fmt.Println("Migration Success")
  }
  ```

### Connection

- Create `pkg` folder, inside it Create `mysql` folder, inside it Create `mysql.go` file, and write below code

  > File: `pkg/mysql/mysql.go`

  ```go
  package mysql

  import (
    "fmt"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
  )

  var DB *gorm.DB

  // Connection Database
  func DatabaseInit() {
    var err error
    dsn := "{USER}:{PASSWORD}@tcp({HOST}:{POST})/{DATABASE}?charset=utf8mb4&parseTime=True&loc=Local"
    DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

    if err != nil {
      panic(err)
    }

    fmt.Println("Connected to Database")
  }
  ```

### Data Transfer Object (DTO)

- Create `dto` folder, inside it create `result` & `users` folder.

  > Folder: `dto/result`

  > Folder: `dto/users`

- Inside `dto/result` folder, create `result.go` file, and write this below code

  > File: `dto/result/result.go`

  ```go
  package dto

  type SuccessResult struct {
    Code int         `json:"code"`
    Data interface{} `json:"data"`
  }

  type ErrorResult struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
  }
  ```

- Inside `dto/users` folder, create `user_request.go` file, and write this below code

  > File: `dto/users/user_request.go`

  ```go
  package usersdto

  type CreateUserRequest struct {
    Name     string `json:"name" form:"name" validate:"required"`
    Email    string `json:"email" form:"email" validate:"required"`
    Password string `json:"password" form:"password" validate:"required"`
  }

  type UpdateUserRequest struct {
    Name     string `json:"name" form:"name"`
    Email    string `json:"email" form:"email"`
    Password string `json:"password" form:"password"`
  }
  ```

- Inside `dto/users` folder, create `user_response.go` file, and write this below code

  > File: `dto/users/user_response.go`

  ```go
  package usersdto

  type UserResponse struct {
    ID       int    `json:"id"`
    Name     string `json:"name" form:"name" validate:"required"`
    Email    string `json:"email" form:"email" validate:"required"`
    Password string `json:"password" form:"password" validate:"required"`
  }
  ```

# Fetching Query with Gorm

### Repositories

- Create `repositories` folder, inside it create `user.go` file, and write this below code

  > File: `repositories/user.go`

  ```go
  package repositories

  import (
    "dumbmerch/models"
    "gorm.io/gorm"
  )

  type UserRepository interface {
    FindUsers() ([]models.User, error)
    GetUser(ID int) (models.User, error)
  }

  type repository struct {
    db *gorm.DB
  }

  func RepositoryUser(db *gorm.DB) *repository {
    return &repository{db}
  }

  func (r *repository) FindUsers() ([]models.User, error) {
    var users []models.User
    err := r.db.Raw("SELECT * FROM users").Scan(&users).Error

    return users, err
  }

  func (r *repository) GetUser(ID int) (models.User, error) {
    var user models.User
    err := r.db.Raw("SELECT * FROM users WHERE id=?", ID).Scan(&user).Error

    return user, err
  }
  ```

### Handlers

- On `handlers` folder, create `user.go` file, and write this below code

  > File: `handlers/user.go`

  ```go
  package handlers

  import (
    dto "dumbmerch/dto/result"
    usersdto "dumbmerch/dto/users"
    "dumbmerch/models"
    "dumbmerch/repositories"
    "net/http"
    "strconv"

    "github.com/labstack/echo"
  )

  type handler struct {
    UserRepository repositories.UserRepository
  }

  func HandlerUser(UserRepository repositories.UserRepository) *handler {
    return &handler{UserRepository}
  }

  func (h *handler) FindUsers(c echo.Context) error {
    users, err := h.UserRepository.FindUsers()
    if err != nil {
      return c.JSON(http.StatusBadRequest, dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()})
    }

    return c.JSON(http.StatusOK, dto.SuccessResult{Code: http.StatusOK, Data: users})
  }

  func (h *handler) GetUser(c echo.Context) error {
    id, _ := strconv.Atoi(c.Param("id"))

    user, err := h.UserRepository.GetUser(id)
    if err != nil {
      return c.JSON(http.StatusBadRequest, dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()})
    }

    return c.JSON(http.StatusOK, dto.SuccessResult{Code: http.StatusOK, Data: convertResponse(user)})
  }

  func convertResponse(u models.User) usersdto.UserResponse {
    return usersdto.UserResponse{
      ID:       u.ID,
      Name:     u.Name,
      Email:    u.Email,
      Password: u.Password,
    }
  }
  ```

### Routes

- On `routes` folder, create `user.go`, and write this below code

  > File: `routes/user.go`

  ```go
  package routes

  import (
    "dumbmerch/handlers"
    "dumbmerch/pkg/mysql"
    "dumbmerch/repositories"

    "github.com/labstack/echo"
  )

  func UserRoutes(e *echo.Group) {
    userRepository := repositories.RepositoryUser(mysql.DB)
    h := handlers.HandlerUser(userRepository)

    e.GET("/users", h.FindUsers)
    e.GET("/user/:id", h.GetUser)
  }
  ```

### Root file `main.go`

Modify `main.go` file, adding `Initial Database` and Running `Auto Migration`

```go
package main

import (
	"dumbmerch/database"
	"dumbmerch/pkg/mysql"
	"dumbmerch/routes"
	"fmt"

	"github.com/labstack/echo"
)

func main() {
	e := echo.New()

	mysql.DatabaseInit()
	database.RunMigration()

	routes.RouteInit(e.Group("/api/v1"))

	fmt.Println("server running localhost:5000")
	e.Logger.Fatal(e.Start("localhost:5000"))
}
```
