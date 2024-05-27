package helpers

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
)

func GetIdFromRequest(w http.ResponseWriter, r *http.Request, prefix string) (int, error) {
	logger := slog.With("Location", "GetIdHelper")
	idStr := strings.TrimPrefix(r.URL.String(), prefix)
	id, err := strconv.Atoi(idStr)

	if err != nil {
		logger.Warn("Could not get products by id: invalid id.", "Error", err)
		fmt.Fprintf(w, "Invalid product id.")
	}

	return id, err
}
