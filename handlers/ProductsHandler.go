package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"ward4woods.ca/data"
	"ward4woods.ca/helpers"
	"ward4woods.ca/models"
)

type ProductsHandler struct {
	productsStore *data.ProductsStore
	logger        *slog.Logger
}

func NewProductsHandler(productsStore *data.ProductsStore, logger *slog.Logger) *ProductsHandler {
	return &ProductsHandler{
		productsStore: productsStore,
		logger:        logger.With("Location", "ProductsHandler"),
	}
}

func (ph *ProductsHandler) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	products, err := ph.productsStore.GetAllProducts()

	if err != nil {
		ph.logger.Error("Could not get products from store. Writing error response.")
		fmt.Fprintf(w, "Server error when trying to load products.")
		return
	}

	helpers.RenderTemplate(w, "html/templates/productsList.html", products)
}

func (ph *ProductsHandler) GetProductById(w http.ResponseWriter, r *http.Request) {
	id, err := ph.getIdFromApiRequest(w, r)
	if err != nil {
		return
	}

	product, err := ph.productsStore.GetProductById(id)

	if err == sql.ErrNoRows {
		fmt.Fprint(w, "Product not found.")
		ph.logger.Info("Attempted to find product by id, but product didn't exist.")
		return
	}

	if err != nil {
		fmt.Fprint(w, "Error finding product.")
		ph.logger.Warn("Error when finding product from database.")
		return
	}

	json.NewEncoder(w).Encode(product)
}

func (ph *ProductsHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	product := models.Product{}
	err := json.NewDecoder(r.Body).Decode(&product)

	if err != nil {
		ph.logger.Error("Could not parse product from body. Error:", err)
		fmt.Fprintf(w, "Error parsing product from json.")
		return
	}

	err = ph.productsStore.AddProduct(product)

	if err != nil {
		ph.logger.Warn("Error adding product to product store.")
		fmt.Fprintf(w, "Error adding product to database.")
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (ph *ProductsHandler) DeleteProductById(w http.ResponseWriter, r *http.Request) {
	id, err := ph.getIdFromApiRequest(w, r)

	if err != nil {
		return
	}

	err = ph.productsStore.DeleteProductById(id)

	if err != nil {
		fmt.Fprintf(w, "Error deleting product.")
		ph.logger.Warn("Could not delete product from products store.")
	}

	w.WriteHeader(http.StatusOK)
}

func (ph *ProductsHandler) getIdFromApiRequest(w http.ResponseWriter, r *http.Request) (int, error) {
	prefix := "api/products"
	return helpers.GetIdFromRequest(w, r, prefix)
}
