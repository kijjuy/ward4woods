package helpers

import (
	"net/http"
	"strconv"
	"strings"

	"ward4woods.ca/models"
)

func GetIdFromRequest(r *http.Request, prefix string) (models.ProductId, error) {
	url := strings.ToLower(r.URL.String())
	idStr := strings.TrimPrefix(url, prefix)
	id, err := strconv.Atoi(idStr)

	if err != nil {
		return 0, err
	}

	return models.ProductId(id), err
}
