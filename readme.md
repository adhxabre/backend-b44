### Table of Contents

- [Handle Upload File](#handle-upload-file)
  - [Introduction](#introduction)
  - [Package](#Package)
  - [Routes](#routes)
  - [Handler](#Handler)
  - [Folder Store File](#folder-store-file)
  - [DotEnv](#dotenv)

---

# Handle Upload File

## Introduction

For this section:

- Handle File Upload for `Create Product` data

## Package

- Inside `pkg` folder, in `middleware` folder, inside it create `uploadFile.go` file, and write this below code

  > File: `pkg/middleware/upload_file.go`

  ```go
  package middleware

  import (
    "io"
    "io/ioutil"
    "net/http"

    "github.com/labstack/echo"
  )

  func UploadFile(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
      file, err := c.FormFile("image")
      if err != nil {
        return c.JSON(http.StatusBadRequest, err)
      }

      src, err := file.Open()
      if err != nil {
        return c.JSON(http.StatusBadRequest, err)
      }
      defer src.Close()

      tempFile, err := ioutil.TempFile("uploads", "image-*.png")
      if err != nil {
        return c.JSON(http.StatusBadRequest, err)
      }
      defer tempFile.Close()

      if _, err = io.Copy(tempFile, src); err != nil {
        return c.JSON(http.StatusBadRequest, err)
      }

      data := tempFile.Name()
      filename := data[8:] // split uploads/

      c.Set("dataFile", filename)
      return next(c)
    }
  }
  ```

## Routes

- In `routes` folder, inside `product.go` file, write `uploadFile` middleware on `/product` route

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
    e.POST("/product", middleware.Auth(middleware.UploadFile(h.CreateProduct)))
  }
  ```

## Handler

- In `handlers` folder, inside `product.go` file, write get `filename` and store like this below code

  > File: `handlers/product.go`

  ```go
  func (h *handlerProduct) CreateProduct(c echo.Context) error {
    dataFile := c.Get("dataFile").(string)
    fmt.Println("this is data file", dataFile)

    price, _ := strconv.Atoi(c.FormValue("price"))
    qty, _ := strconv.Atoi(c.FormValue("qty"))
    category_id, _ := strconv.Atoi(c.FormValue("category_id"))

    request := productdto.ProductRequest{
      Name:       c.FormValue("name"),
      Desc:       c.FormValue("desc"),
      Price:      price,
      Image:      dataFile,
      Qty:        qty,
      CategoryID: category_id,
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

- Embed Path file in `FindProducts` and `GetProduct` method
  > File: `handlers/product.go`
  - Create `path_file` Global variable
    ```go
    var path_file = "http://localhost:5000/uploads/"
    ```
  - `FindProducts` method
    ```go
    for i, p := range products {
      products[i].Image = path_file + p.Image
    }
    ```
  - `GetProduct` method
    ```go
    product.Image = path_file + product.Image
    ```

## Folder Store File

- Create `uploads` folder

  > File: `./uploads`

- Add this below code to make `uploads` can be used another client

  > File: `main.go`

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

    e.Static("/uploads", "./uploads")

    fmt.Println("server running localhost:5000")
    e.Logger.Fatal(e.Start("localhost:5000"))
  }
  ```

## DotEnv

- Installation

  ```bash
  go get github.com/joho/godotenv
  ```

- Create `.env` file and write this below code

  > File: `.env`

  ```env
  SECRET_KEY=suryaganteng
  ```

- In `main.go` file import `godotenv` and Init `godotenv` inside `main` function like this below code

  > File: `main.go`

  - Import `godotenv` package
    ```go
    import (
      // another package here ...
      "github.com/joho/godotenv" // import this package
    )
    ```
  - Init `godotenv`

    ```go
    func main() {

      	// env
        errEnv := godotenv.Load()
        if errEnv != nil {
          panic("Failed to load env file")
        }

        // Another code on this below ...
    }
    ```

- How to use Environment Variable, write this below code inside `jwt.go` file

  > File: `pkg/jwt/jwt.go`

  - Import `os` package
    ```go
    import (
      "fmt"
      "os" // import this package
      "github.com/golang-jwt/jwt/v4"
    )
    ```
  - Modify `SecretKey` variable like this below code
    ```go
    var SecretKey = os.Getenv("SECRET_KEY")
    ```
