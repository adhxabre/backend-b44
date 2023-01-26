package handlers

import (
	productdto "dumbmerch/dto/product"
	dto "dumbmerch/dto/result"
	"dumbmerch/models"
	"dumbmerch/repositories"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
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
