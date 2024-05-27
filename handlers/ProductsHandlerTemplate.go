package handlers

import (
	"database/sql"
	"net/http"
	"ward4woods.ca/helpers"
)

func (ph *ProductsHandler) ProductsDetails(w http.ResponseWriter, r *http.Request) {
	id, err := ph.getIdFromTemplateRequest(w, r)
	if err != nil {
		return
	}

	product, err := ph.productsStore.GetProductById(id)

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
	helpers.RenderTemplate(w, "html/templates/productDetails.html", product, ph.logger)
}

func (ph *ProductsHandler) ProductsList(w http.ResponseWriter, r *http.Request) {
	products, err := ph.productsStore.GetAllProducts()

	if err != nil {
		ph.logger.Error("Could not get products from store. Writing error response.")
		http.Error(w, "Server error when trying to load products.", http.StatusInternalServerError)
		return
	}

	helpers.RenderTemplate(w, "html/templates/productsList.html", products, ph.logger)
}

func (ph *ProductsHandler) getIdFromTemplateRequest(w http.ResponseWriter, r *http.Request) (int, error) {
	prefix := "/products/"
	return ph.getIdFromRequest(w, r, prefix)
}
