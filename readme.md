## Table of contents

- [Insert Query with Gorm](#insert-query-with-gorm)
  - [Repositories](#repositories)
  - [Handlers](#handlers)
  - [Routes](#routes)

# Insert Query with Gorm

### Repositories

> File: `repositories/user.go`

- Import `time`

  ```go
  import (
    "dumbmerch/models"
    "time"
    "gorm.io/gorm"
  )
  ```

- Declare `CreateUser` interface
  ```go
  type UserRepository interface {
    FindUsers() ([]models.User, error)
    GetUser(ID int) (models.User, error)
    CreateUser(user models.User) (models.User, error) // Write this code
  }
  ```
- Write `CreateUser` function

  ```go
   // Write this code
  func (r *repository) CreateUser(user models.User) (models.User, error) {
    err := r.db.Exec("INSERT INTO users(name,email,password,created_at,updated_at) VALUES (?,?,?,?,?)",user.Name,user.Email, user.Password, time.Now(), time.Now()).Error

    return user, err
  }
  ```

### Handlers

> File: `handlers/user.go`

- Import `Validator`

  ```go
  package handlers

  import (
    dto "dumbmerch/dto/result"
    usersdto "dumbmerch/dto/users"
    "dumbmerch/models"
    "dumbmerch/repositories"
    "net/http"
    "strconv"

    "github.com/go-playground/validator/v10"
    "github.com/labstack/echo/v4"
  )
  ```

- Write `CreateUser` function

  ```go
  func (h *handler) CreateUser(c echo.Context) error {
    request := new(usersdto.CreateUserRequest)
    if err := c.Bind(request); err != nil {
      return c.JSON(http.StatusBadRequest, dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()})
    }

    validation := validator.New()
    err := validation.Struct(request)
    if err != nil {
      return c.JSON(http.StatusBadRequest, dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()})
    }

    // data form pattern submit to pattern entity db user
    user := models.User{
      Name:     request.Name,
      Email:    request.Email,
      Password: request.Password,
    }

    data, err := h.UserRepository.CreateUser(user)
    if err != nil {
      return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()})
    }

    return c.JSON(http.StatusOK, dto.SuccessResult{Code: http.StatusOK, Data: convertResponse(data)})
  }
  ```

### Routes

> File: `routes/user.go`

- Write `Create User` route with `POST` method

  ```go
  func UserRoutes(e *echo.Group) {
    userRepository := repositories.RepositoryUser(mysql.DB)
    h := handlers.HandlerUser(userRepository)

    e.GET("/users", h.FindUsers)
    e.GET("/user/:id", h.GetUser)
    e.POST("/user", h.CreateUser)
  }
  ```
