## Table of contents

- [Delete Query with Gorm](#delete-query-with-gorm)
  - [Repositories](#repositories)
  - [Handlers](#handlers)
  - [Routes](#routes)

# Delete Query with Gorm

### Repositories

> File: `repositories/user.go`

- Declare `DeleteUser` interface
  ```go
  type UserRepository interface {
    FindUsers() ([]models.User, error)
    GetUser(ID int) (models.User, error)
    CreateUser(user models.User) (models.User, error)
    UpdateUser(user models.User, ID int) (models.User, error)
    DeleteUser(user models.User, ID int) (models.User, error) // Write this code
  }
  ```
- Write `DeleteUser` function

  ```go
   // Write this code
  func (r *repository) DeleteUser(user models.User,ID int) (models.User, error) {
    err := r.db.Raw("DELETE FROM users WHERE id=?",ID).Scan(&user).Error

    return user, err
  }
  ```

### Handlers

> File: `handlers/user.go`

- Write `DeleteUser` function

  ```go
  func (h *handler) DeleteUser(c echo.Context) error {
    id, _ := strconv.Atoi(c.Param("id"))

    user, err := h.UserRepository.GetUser(id)
    if err != nil {
      return c.JSON(http.StatusBadRequest, dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()})
    }

    data, err := h.UserRepository.DeleteUser(user, id)
    if err != nil {
      return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()})
    }

    return c.JSON(http.StatusOK, dto.SuccessResult{Code: http.StatusOK, Data: convertResponse(data)})
  }
  ```

### Routes

> File: `routes/user.go`

- Write `Delete User` route with `DELETE` method

  ```go
  func UserRoutes(e *echo.Group) {
    userRepository := repositories.RepositoryUser(mysql.DB)
    h := handlers.HandlerUser(userRepository)

    e.GET("/users", h.FindUsers)
    e.GET("/user/:id", h.GetUser)
    e.POST("/user", h.CreateUser)
    e.PATCH("/user/:id", h.UpdateUser)
    e.DELETE("/user/:id", h.DeleteUser)
  }
  ```
