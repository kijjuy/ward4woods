package helpers

import (
	"net/http"
	"strconv"
	"strings"
)

func GetIdFromRequest(r *http.Request, prefix string) (int, error) {
	idStr := strings.TrimPrefix(r.URL.String(), prefix)
	id, err := strconv.Atoi(idStr)

	if err != nil {
		return 0, err
	}

	return id, err
}
