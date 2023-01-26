### Table of Contents

- [Authentication JWT](#authentication-jwt)
  - [Introduction](#introduction)
  - [Installation](#installation)
  - [Package](#Package)
  - [Handler](#Handler)
  - [Repository](#repository)
  - [Routes](#routes)

---

# Authentication JWT

Reference: [Golang JWT](https://github.com/golang-jwt/jwt)

## Introduction

For this section:

- Generate Token using JWT if `User Login`
- Verify Token and Get User Data if `Create Product Data`

## Installation

- Golang Json Web Token (JWT)

  ```bash
  go get -u github.com/golang-jwt/jwt/v4
  ```

## Package

- Inside `pkg` folder, create `jwt` folder, inside it create `jwt.go` file, and write this below code

  > File: `pkg/jwt/jwt.go`

  ```go
  package jwtToken

  import (
    "fmt"

    "github.com/golang-jwt/jwt/v4"
  )

  var SecretKey = "SECRET_KEY"

  func GenerateToken(claims *jwt.MapClaims) (string, error) {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    webtoken, err := token.SignedString([]byte(SecretKey))
    if err != nil {
      return "", err
    }

    return webtoken, nil
  }

  func VerifyToken(tokenString string) (*jwt.Token, error) {
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
      if _, isValid := token.Method.(*jwt.SigningMethodHMAC); !isValid {
        return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
      }
      return []byte(SecretKey), nil
    })

    if err != nil {
      return nil, err
    }
    return token, nil
  }

  func DecodeToken(tokenString string) (jwt.MapClaims, error) {
    token, err := VerifyToken(tokenString)
    if err != nil {
      return nil, err
    }

    claims, isOk := token.Claims.(jwt.MapClaims)
    if isOk && token.Valid {
      return claims, nil
    }

    return nil, fmt.Errorf("invalid token")
  }
  ```

- Inside `pkg` folder, create `middleware` folder, inside it create `auth.go` file, and write this below code

  > File: `pkg/middleware/auth.go`

  ```go
  package middleware

  import (
    dto "dumbmerch/dto/result"
    jwtToken "dumbmerch/pkg/jwt"
    "net/http"
    "strings"

    "github.com/labstack/echo"
  )

  // Declare Result struct here ...
  type Result struct {
    Code    int         `json:"code"`
    Data    interface{} `json:"data"`
    Message string      `json:"message"`
  }

  // Create Auth function here ...
  func Auth(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
      token := c.Request().Header.Get("Authorization")

      if token == "" {
        return c.JSON(http.StatusUnauthorized, dto.ErrorResult{Code: http.StatusBadRequest, Message: "unauthorized"})
      }

      token = strings.Split(token, " ")[1]
      claims, err := jwtToken.DecodeToken(token)

      if err != nil {
        return c.JSON(http.StatusUnauthorized, Result{Code: http.StatusUnauthorized, Message: "unathorized"})
      }

      c.Set("userLogin", claims)
      return next(c)
    }
  }
  ```

## Handler

- Inside `handlers` folder, On `auth.go` file and write `Login` Function like this below code

  > File: `handlers/auth.go`

  ```go
  func (h *handlerAuth) Login(c echo.Context) error {
    request := new(authdto.LoginRequest)
    if err := c.Bind(request); err != nil {
      return c.JSON(http.StatusBadRequest, dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()})
    }

    user := models.User{
      Email:    request.Email,
      Password: request.Password,
    }

    // Check email
    user, err := h.AuthRepository.Login(user.Email)
    if err != nil {
      return c.JSON(http.StatusBadRequest, dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()})
    }

    // Check password
    isValid := bcrypt.CheckPasswordHash(request.Password, user.Password)
    if !isValid {
      return c.JSON(http.StatusBadRequest, dto.ErrorResult{Code: http.StatusBadRequest, Message: "wrong email or password"})
    }

    //generate token
    claims := jwt.MapClaims{}
    claims["id"] = user.ID
    claims["exp"] = time.Now().Add(time.Hour * 2).Unix() // 2 hours expired

    token, errGenerateToken := jwtToken.GenerateToken(&claims)
    if errGenerateToken != nil {
      log.Println(errGenerateToken)
      return echo.NewHTTPError(http.StatusUnauthorized)
    }

    loginResponse := authdto.LoginResponse{
      Name:     user.Name,
      Email:    user.Email,
      Password: user.Password,
      Token:    token,
    }

    return c.JSON(http.StatusOK, dto.SuccessResult{Code: http.StatusOK, Data: loginResponse})
  }
  ```

  - Inside `handlers` folder, On `product.go` file and write `CreateProduct` Function like this below code

  > File: `handlers/product.go`

  ```go
  func (h *handlerProduct) CreateProduct(c echo.Context) error {
    request := new(productdto.ProductRequest)
    if err := c.Bind(request); err != nil {
      return c.JSON(http.StatusBadRequest, dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()})
    }

    validation := validator.New()
    err := validation.Struct(request)
    if err != nil {
      return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()})
    }

    userLogin := c.Get("userLogin")
    userId := userLogin.(jwt.MapClaims)["id"].(float64)

    product := models.Product{
      Name:   request.Name,
      Desc:   request.Desc,
      Price:  request.Price,
      Image:  request.Image,
      Qty:    request.Qty,
      UserID: int(userId),
    }

    product, err = h.ProductRepository.CreateProduct(product)
    if err != nil {
      return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()})
    }

    product, _ = h.ProductRepository.GetProduct(product.ID)

    return c.JSON(http.StatusOK, dto.SuccessResult{Code: http.StatusOK, Data: convertResponseProduct(product)})
  }
  ```

## Repository

- Inside `repositories` folder, On `auth.go` file, write `Login` function like this below code

  > File: `repositories/auth.go`

  ```go
  func (r *repository) Login(email string) (models.User, error) {
    var user models.User
    err := r.db.First(&user, "email=?", email).Error

    return user, err
  }
  ```

## Routes

- Inside `routes` folder, in `auth.go` file and write `Login` route this below code

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
    e.POST("/login", h.Login) // add this code
  }
  ```

- Inside `routes` folder, in `product.go` file and write `product` route with `middleware` like this below code

  > File: `routes/product.go`

  ```go
  package routes

  import (
    "dumbmerch/handlers"
    "dumbmerch/pkg/middleware"
    "dumbmerch/pkg/mysql"
    "dumbmerch/repositories"

    "github.com/labstack/echo"
  )

  func ProductRoutes(e *echo.Group) {
    productRepository := repositories.RepositoryProduct(mysql.DB)
    h := handlers.HandlerProduct(productRepository)

    e.GET("/products", middleware.Auth(h.FindProducts))
    e.GET("/product/:id", middleware.Auth(h.GetProduct))
    e.POST("/product", middleware.Auth(h.CreateProduct))
  }
  ```
