package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"ward4woods.ca/data"
	"ward4woods.ca/helpers"
)

func (ph *ProductsHandler) ProductsList(w http.ResponseWriter, r *http.Request) {
	products, err := ph.productsStore.GetAllProducts()

	if err != nil {
		ph.logger.Error("Could not get products from store. Writing error response.")
		http.Error(w, "Server error when trying to load products.", http.StatusInternalServerError)
		return
	}

	helpers.RenderTemplate(w, "html/templates/productsList.html", products)
}

func (ph *ProductsHandler) getIdFromTemplateRequest(w http.ResponseWriter, r *http.Request) (int, error) {
	prefix := "/products/"
	return helpers.GetIdFromRequest(w, r, prefix)
}

// ProductsDetails handles creating the individual page for each product.
// Will be sent requests from the URL '/products/{id}'
func (ph *ProductsHandler) ProductsDetails(w http.ResponseWriter, r *http.Request, productsStore *data.ProductsStore) {
	id, err := helpers.GetIdFromRequest(w, r, "/products/")
	if err != nil {
		return
	}

	product, err := productsStore.GetProductById(id)

	if err == sql.ErrNoRows {
		ph.logger.Warn("Attempted to find product by id, but product didn't exist.")
		http.Error(w, "Product not found.", http.StatusNotFound)
		return
	}

	if err != nil {
		ph.logger.Warn("Error when finding product from database.")
		http.Error(w, "Error finding product.", http.StatusInternalServerError)
		return
	}

	ph.logger.Info(fmt.Sprintf("now serving details page for product: %+v", product))

	helpers.RenderTemplate(w, "templates/productDetails.html", product)

	//TODO: Setup details page to work with htmx rather than loading 2 templates form go
}
