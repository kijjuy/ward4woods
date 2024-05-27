package handlers

import (
	"database/sql"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"path"
	"path/filepath"
	"ward4woods.ca/data"
	"ward4woods.ca/helpers"
)

type StaticHandler struct {
	staticPath string
	template   string
	error      string
	logger     *slog.Logger
}

func NewStaticHandler(staticPath string, template string, error string, logger *slog.Logger) *StaticHandler {
	return &StaticHandler{
		staticPath: staticPath,
		template:   template,
		error:      error,
		logger:     logger.With("Location", "StaticHandler"),
	}
}

func (sh *StaticHandler) Render(w http.ResponseWriter, page string, dataArr ...interface{}) {
	var data interface{}
	if len(dataArr) == 1 {
		data = dataArr[0]
	}
	sh.logger.Info(fmt.Sprintf("Data is: %+v", data))
	page = filepath.Join(sh.staticPath, page)
	tmpl, err := template.ParseFiles(sh.template, page)
	if err == nil {
		tmpl.Execute(w, data)
		sh.logger.Info(fmt.Sprintf("Now serving page: %s", page))
		return
	}

	message := "Could not find page: " + page + ". "
	page = path.Join(page, "index.html")
	message += "Checking: " + page
	sh.logger.Info(fmt.Sprintf("%s Error: %s", message, err))

	tmpl, err = template.ParseFiles(sh.template, page)
	if err == nil {
		tmpl.Execute(w, data)
		sh.logger.Info(fmt.Sprintf("Now serving page: %s", page))
		return
	}

	sh.logger.Info(fmt.Sprintf("Could not find page. Loading error view. Error: %s", err))
	tmpl, err = tmpl.ParseFiles(sh.template, sh.error)
	if err != nil {
		sh.logger.Error(fmt.Sprintf("Could not find error page. Error: '%s'.", err))
		sh.logger.Error(fmt.Sprintf("Attempted to find template page at: %s. Attempted to find error page at %s.", sh.template, sh.error))
		return
	}
	tmpl.Execute(w, data)
	return
}

func (sh *StaticHandler) HandleRequests(w http.ResponseWriter, r *http.Request) {
	url := filepath.Clean(r.URL.String())
	sh.Render(w, url)
}
func (sh *StaticHandler) ProductsDetails(w http.ResponseWriter, r *http.Request, productsStore *data.ProductsStore) {
	id, err := helpers.GetIdFromRequest(w, r, "/products/")
	if err != nil {
		return
	}

	product, err := productsStore.GetProductById(id)

	if err == sql.ErrNoRows {
		sh.logger.Warn("Attempted to find product by id, but product didn't exist.")
		http.Error(w, "Product not found.", http.StatusNotFound)
		return
	}

	if err != nil {
		sh.logger.Warn("Error when finding product from database.")
		http.Error(w, "Error finding product.", http.StatusInternalServerError)
		return
	}

	sh.logger.Info(fmt.Sprintf("now serving details page for product: %+v", product))

	sh.Render(w, "templates/productDetails.html", product)
}
