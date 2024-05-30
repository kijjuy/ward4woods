package handlers

import (
	"net/http"
	"ward4woods.ca/helpers"
)

func (ph *ProductsHandler) getIdFromTemplateRequest(w http.ResponseWriter, r *http.Request) (int, error) {
	prefix := "/products/"
	return helpers.GetIdFromRequest(r, prefix)
}
