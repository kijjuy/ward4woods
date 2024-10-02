package handlers

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"w4w/models"
	"w4w/services"

	"github.com/labstack/echo/v4"
)

func GetAllProducts(c echo.Context) error {
	products, err := services.GetAllProducts()

	if err != nil {
		c.Logger().Error("Error getting products from database", "Error", err)
		return err
	}

	return c.Render(http.StatusOK, "productsList", products)
}

func ProductDetails(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		c.Logger().Error("Error getting id from URL path", "Error", err)
	}

	product, err := services.GetProductById(id)

	if err != nil {
		slog.Error("Error getting products from database", "Error", err)
		return err
	}

	return c.Render(http.StatusOK, "productDetails", product)
}

// TODO: Authorize endpoint
func DeleteProduct(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		return err
	}

	err = services.DeleteProduct(id)

	//TODO: properly handle error on frontend (perhaps with js alert)
	if errors.Is(err, &services.ErrNoRowsAffected{}) {
		slog.Warn("Attemped to delete a product that does not exist", "ProductID", id)
		c.NoContent(http.StatusNotFound)
	}

	if err != nil {
		return err
	}

	slog.Info("Deleted product from database", "Id", id)
	return c.NoContent(http.StatusOK)
}

// TODO: authorize endpoint
func NewProduct(c echo.Context) error {
	product, err := getProductFromForm(c)

	if err != nil {
		return err
	}

	err = services.CreateProduct(product)

	if err != nil {
		return err
	}

	logNewProduct(product)

	return c.NoContent(http.StatusOK)
}

func EditProduct(c echo.Context) error {
	productIdStr := c.Param("id")
	productId, err := strconv.Atoi(productIdStr)

	if err != nil {
		return err
	}

	product, err := services.GetProductById(productId)

	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "editProduct", product)
}

// TODO: authorize endpoint
func UpdateProduct(c echo.Context) error {
	productIdStr := c.Param("id")
	productId, err := strconv.Atoi(productIdStr)

	if err != nil {
		slog.Error("Error converting id param to string", "Error", err)
		return err
	}

	product, err := getProductFromForm(c)

	if err != nil {
		slog.Error("Error getting product from form", "Error", err)
		return err
	}

	err = services.UpdateProduct(productId, product)

	if err != nil {
		slog.Error("Error updating product in database", "Error", err)
		return err
	}

	return c.NoContent(http.StatusOK)
}

func AdminGetProductsList(c echo.Context) error {
	products, err := services.GetAllProducts()

	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "adminProductsList", products)
}

func GetCategories(c echo.Context) error {
	productIdStr := c.Param("id")
	productId, err := strconv.Atoi(productIdStr)

	if err != nil {
		return err
	}

	currentProduct, err := services.GetProductById(productId)

	if err != nil {
		return err
	}

	categories, err := services.GetCategories(currentProduct.Category)

	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "categorySelect", categories)
}

func getProductFromForm(c echo.Context) (models.Product, error) {
	name := c.FormValue("name")
	priceStr := c.FormValue("price")
	description := c.FormValue("description")
	category := c.FormValue("category")

	priceStr = strings.Replace(priceStr, ".", "", 1)
	price, err := strconv.Atoi(priceStr)

	if err != nil {
		return models.NewProduct(), err
	}

	product := models.Product{
		Name:        name,
		Price:       price,
		Description: description,
		Category:    category,
	}

	return product, nil
}

func logNewProduct(product models.Product) {
	slog.Info("Added new product to database",
		"name", product.Name,
		"price", product.Price,
		"category", product.Category,
	)
}
