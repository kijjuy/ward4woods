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
