### Table of Contents

- [Hashing Password](#hashing-password)
  - [Introduction](#intoduction)
  - [Package](#Package)
  - [Handler](#Handler)
  - [Repository](#repository)
  - [Routes](#routes)

---

# Hashing Password

Reference: [Go Bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt)

## Introduction

For this section, Hashing password if User doing Register New Account

## Package

- Inside `pkg` folder, create `bcrypt` folder, inside it create `hash_password.go` file, and write this below code

  > File: `pkg/bcrypt/hash_password.go`

  ```go
  package bcrypt

  import "golang.org/x/crypto/bcrypt"

  func HashingPassword(password string) (string, error) {
    hashedByte, err := bcrypt.GenerateFromPassword([]byte(password), 10)
    if err != nil {
      return "", err
    }
    return string(hashedByte), nil
  }

  func CheckPasswordHash(password, hashedPassword string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
    return err == nil
  }
  ```

## Handler

- Inside `handlers` folder, create `auth.go` file and write this below code

  > File: `handlers/auth.go`

  ```go
  package handlers

  import (
    authdto "dumbmerch/dto/auth"
    dto "dumbmerch/dto/result"
    "net/http"

    "dumbmerch/models"
    "dumbmerch/pkg/bcrypt"
    "dumbmerch/repositories"

    "github.com/go-playground/validator/v10"
    "github.com/labstack/echo"
  )

  type handlerAuth struct {
    AuthRepository repositories.AuthRepository
  }

  func HandlerAuth(AuthRepository repositories.AuthRepository) *handlerAuth {
    return &handlerAuth{AuthRepository}
  }

  func (h *handlerAuth) Register(c echo.Context) error {
    request := new(authdto.AuthRequest)
    if err := c.Bind(request); err != nil {
      return c.JSON(http.StatusBadRequest, dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()})
    }

    validation := validator.New()
    err := validation.Struct(request)
    if err != nil {
      return c.JSON(http.StatusBadRequest, dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()})
    }

    password, err := bcrypt.HashingPassword(request.Password)
    if err != nil {
      return c.JSON(http.StatusBadRequest, dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()})
    }

    user := models.User{
      Name:     request.Name,
      Email:    request.Email,
      Password: password,
    }

    data, err := h.AuthRepository.Register(user)
    if err != nil {
      return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()})
    }

    return c.JSON(http.StatusOK, dto.SuccessResult{Code: http.StatusOK, Data: data})
  }
  ```

## Repository

- Inside `repositories` folder, create `auth.go` file and write this below code

  > File: `repositories/auth.go`

  ```go
  package repositories

  import (
    "dumbmerch/models"

    "gorm.io/gorm"
  )

  type AuthRepository interface {
    Register(user models.User) (models.User, error)
  }

  func RepositoryAuth(db *gorm.DB) *repository {
    return &repository{db}
  }

  func (r *repository) Register(user models.User) (models.User, error) {
    err := r.db.Create(&user).Error

    return user, err
  }
  ```

## Routes

- Inside `routes` folder, create `auth.go` file and write this below code

  > File: `routes/auth.go`

  ```go
  package routes

  import (
    "dumbmerch/handlers"
    "dumbmerch/pkg/mysql"
    "dumbmerch/repositories"

    "github.com/labstack/echo"
  )

  func AuthRoutes(e *echo.Group) {
    authRepository := repositories.RepositoryAuth(mysql.DB)
    h := handlers.HandlerAuth(authRepository)

    e.POST("/register", h.Register)
  }
  ```
