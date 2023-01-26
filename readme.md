### Table of Contents

- [GORM Relation belongs to](#gorm-relation-belongs-to)
  - [Handlers](#handlers)
  - [Repository](#repository)
  - [Routes](#routes)

---

# GORM Relation Belongs to

Reference: [Official GORM Website](https://gorm.io/docs/belongs_to.html)

## Relation

For this section, example Belongs To relation:

- `Profile` &rarr; `User`: to get Profile User
- `Product` &rarr; `User`: to get Product Owner

## Handlers

- Inside `handlers` folder, create `profile.go` file, and write this below code

  > File: `handlers/profile.go`

  ```go
  package handlers

  import (
    profiledto "dumbmerch/dto/profile"
    dto "dumbmerch/dto/result"
    "dumbmerch/models"
    "dumbmerch/repositories"

    "net/http"
    "strconv"

    "github.com/labstack/echo"
  )

  type handlerProfile struct {
    ProfileRepository repositories.ProfileRepository
  }

  func HandlerProfile(ProfileRepository repositories.ProfileRepository) *handlerProfile {
    return &handlerProfile{ProfileRepository}
  }

  func (h *handlerProfile) GetProfile(c echo.Context) error {
    id, _ := strconv.Atoi(c.Param("id"))

    var profile models.Profile
    profile, err := h.ProfileRepository.GetProfile(id)
    if err != nil {
      return c.JSON(http.StatusBadRequest, dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()})
    }

    return c.JSON(http.StatusOK, dto.SuccessResult{Code: http.StatusOK, Data: convertResponseProfile(profile)})
  }

  func convertResponseProfile(u models.Profile) profiledto.ProfileResponse {
    return profiledto.ProfileResponse{
      ID:      u.ID,
      Phone:   u.Phone,
      Gender:  u.Gender,
      Address: u.Address,
      UserID:  u.UserID,
      User:    u.User,
    }
  }
  ```

- Inside `handlers` folder, create `product.go` file, and write this below code

  > File: `handlers/product.go`

  ```go
  package handlers

  import (
    productdto "dumbmerch/dto/product"
    dto "dumbmerch/dto/result"
    "dumbmerch/models"
    "dumbmerch/repositories"
    "net/http"
    "strconv"

    "github.com/go-playground/validator/v10"
    "github.com/labstack/echo"
  )

  type handlerProduct struct {
    ProductRepository repositories.ProductRepository
  }

  func HandlerProduct(ProductRepository repositories.ProductRepository) *handlerProduct {
    return &handlerProduct{ProductRepository}
  }

  func (h *handlerProduct) FindProducts(c echo.Context) error {
    products, err := h.ProductRepository.FindProducts()
    if err != nil {
      return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()})
    }

    return c.JSON(http.StatusOK, dto.SuccessResult{Code: http.StatusOK, Data: products})
  }

  func (h *handlerProduct) GetProduct(c echo.Context) error {
    id, _ := strconv.Atoi(c.Param("id"))

    var product models.Product
    product, err := h.ProductRepository.GetProduct(id)
    if err != nil {
      return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()})
    }

    return c.JSON(http.StatusOK, dto.SuccessResult{Code: http.StatusOK, Data: convertResponseProduct(product)})
  }

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

    product := models.Product{
      Name:   request.Name,
      Desc:   request.Desc,
      Price:  request.Price,
      Image:  request.Image,
      Qty:    request.Qty,
      UserID: request.UserID,
    }

    product, err = h.ProductRepository.CreateProduct(product)
    if err != nil {
      return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()})
    }

    product, _ = h.ProductRepository.GetProduct(product.ID)

    return c.JSON(http.StatusOK, dto.SuccessResult{Code: http.StatusOK, Data: convertResponseProduct(product)})
  }

  func convertResponseProduct(u models.Product) models.ProductResponse {
    return models.ProductResponse{
      Name:     u.Name,
      Desc:     u.Desc,
      Price:    u.Price,
      Image:    u.Image,
      Qty:      u.Qty,
      User:     u.User,
      Category: u.Category,
    }
  }
  ```

## Repository

- Inside `repositories` folder, create `profile.go` file, and write this below code

  > File: `repositories/profile.go`

  ```go
  package repositories

  import (
    "dumbmerch/models"

    "gorm.io/gorm"
  )

  type ProfileRepository interface {
    GetProfile(ID int) (models.Profile, error)
  }

  func RepositoryProfile(db *gorm.DB) *repository {
    return &repository{db}
  }

  func (r *repository) GetProfile(ID int) (models.Profile, error) {
    var profile models.Profile
    err := r.db.Preload("User").First(&profile, ID).Error

    return profile, err
  }
  ```

- Inside `repositories` folder, create `product.go` file, and write this below code

  > File: `repositories/product.go`

  ```go
  package repositories

  import (
    "dumbmerch/models"

    "gorm.io/gorm"
  )

  type ProductRepository interface {
    FindProducts() ([]models.Product, error)
    GetProduct(ID int) (models.Product, error)
    CreateProduct(product models.Product) (models.Product, error)
  }

  func RepositoryProduct(db *gorm.DB) *repository {
    return &repository{db}
  }

  func (r *repository) FindProducts() ([]models.Product, error) {
    var products []models.Product
    err := r.db.Preload("User").Find(&products).Error

    return products, err
  }

  func (r *repository) GetProduct(ID int) (models.Product, error) {
    var product models.Product
    // not yet using category relation, cause this step doesnt Belong to Many
    err := r.db.Preload("User").First(&product, ID).Error

    return product, err
  }

  func (r *repository) CreateProduct(product models.Product) (models.Product, error) {
    err := r.db.Create(&product).Error

    return product, err
  }
  ```

## Routes

- Inside `routes` folder, create `profile.go` file, and write this below code

  > File: `routes/profile.go`

  ```go
  package routes

  import (
    "dumbmerch/handlers"
    "dumbmerch/pkg/mysql"
    "dumbmerch/repositories"

    "github.com/labstack/echo"
  )

  func ProfileRoutes(e *echo.Group) {
    profileRepository := repositories.RepositoryProfile(mysql.DB)
    h := handlers.HandlerProfile(profileRepository)

    e.GET("/profile/:id", h.GetProfile)
  }
  ```

- Inside `routes` folder, create `profile.go` file, and write this below code

  > File: `routes/product.go`

  ```go
  package routes

  import (
    "dumbmerch/handlers"
    "dumbmerch/pkg/mysql"
    "dumbmerch/repositories"

    "github.com/labstack/echo"
  )

  func ProductRoutes(e *echo.Group) {
    productRepository := repositories.RepositoryProduct(mysql.DB)
    h := handlers.HandlerProduct(productRepository)

    e.GET("/products", h.FindProducts)
    e.GET("/product/:id", h.GetProduct)
    e.POST("/product", h.CreateProduct)
  }
  ```

- On `routes.go` file, write `ProfileRoutes` and `ProductRoutes`

  > File: `routes/routes.go`

  ```go
  package routes

  import "github.com/labstack/echo"

  func RouteInit(e *echo.Group) {
    UserRoutes(e)
    ProfileRoutes(e)
    ProductRoutes(e)
  }
  ```
