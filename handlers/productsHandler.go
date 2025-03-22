package handlers

import (
	"errors"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"w4w/models"
	"w4w/services"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/shopspring/decimal"
)

func GetAllProducts(c echo.Context) error {
	products, err := services.GetAllProducts()

	if err != nil {
		c.Logger().Error("Error getting products from database", "Error", err)
		return err
	}

	displayProducts := make([]models.ProductListDisplayModel, 0)

	for _, product := range products {
		imageId, err := services.GetMainProductImage(product.Id)
		if err != nil {
			imageId = "no-image.png"
			slog.Warn("Cound not get image for product.", "ProductId", product.Id, "Error", err)
		}
		displayProduct := models.ProductListDisplayModel{
			Product:          product,
			ProductMainImage: imageId,
		}
		displayProducts = append(displayProducts, displayProduct)

	}

	return c.Render(http.StatusOK, "productsList", displayProducts)
}

func ProductDetails(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		c.Logger().Error("Error getting id from URL path", "Error", err)
	}

	product, err := services.GetProductById(id)

	if err != nil {
		slog.Error("Error getting product from database", "Error", err)
		return err
	}

	images, err := services.GetImagesByProductId(id)

	if err != nil {
		slog.Error("Error getting images from database", "Error", err)
		return err
	}

	productDisplayModel := models.ProductDetailsDisplayModel{
		Product:     product,
		MainImage:   images[0],
		OtherImages: images[1:],
	}

	slog.Debug("Product is...", "Product", productDisplayModel)

	return c.Render(http.StatusOK, "productDetails", productDisplayModel)
}

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

func NewProduct(c echo.Context) error {
	product, err := getProductFromForm(c)

	if err != nil {
		slog.Error("Error getting count of images", "Error", err)
		return err
	}

	form, err := c.MultipartForm()
	if err != nil {
		return err
	}

	files := form.File["imageUploads[]"]

	for _, image := range files {
		src, err := image.Open()
		if err != nil {
			return err
		}

		defer src.Close()

		filename := uuid.New()

		dst, err := os.Create("uploads/" + filename.String())

		if err != nil {
			slog.Error("Error creating file for images.", "Error", err)
			return err
		}
		defer dst.Close()

		if _, err = io.Copy(dst, src); err != nil {
			slog.Error("Error copying image to new file")
			return err
		}

		newProductId, err := services.CreateProduct(product)

		if err != nil {
			return err
		}

		err = services.CreateNewProductImageDB(newProductId, filename)

		if err != nil {
			return err
		}

	}

	if err != nil {
		slog.Error("Error getting values from form submission.", "Error", err)
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
	price, err := decimal.NewFromString(priceStr)
	description := c.FormValue("description")
	category := c.FormValue("category")

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
