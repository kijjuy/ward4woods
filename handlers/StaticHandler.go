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

// A StaticHandler is meant to handle requests that require an entire page to be loaded. These pages will be
// constructed by the html/template package and returned to the browser.
type StaticHandler struct {
	staticPath string
	template   string
	error      string
	logger     *slog.Logger
}

// NewStaticHandler creates a StaticHandler struct and returns a pointer to it.
func NewStaticHandler(staticPath string, template string, error string, logger *slog.Logger) *StaticHandler {
	return &StaticHandler{
		staticPath: staticPath,
		template:   template,
		error:      error,
		logger:     logger.With("Location", "StaticHandler"),
	}
}

// Render method attempts to write a templated view to the http.ResponseWriter. It takes a string as the path to the page
// to load into the template, and attempts to load that template into the StaticHandler.template field.
func (sh *StaticHandler) Render(w http.ResponseWriter, page string, data interface{}) {
	page = filepath.Join(sh.staticPath, page)

	err := sh.tryServePage(page, w, data)
	if err == nil {
		return
	}

	sh.logger.Warn("Error Rendering Page:", "Error", err)

	page = sh.nextPage(page, "index.html")

	if err := sh.tryServePage(page, w, data); err == nil {
		return
	}

	sh.logger.Info(fmt.Sprintf("Could not find page. Loading error view. Error: %s", err))
	errorPage := filepath.Join(sh.staticPath, "error.html")
	err = sh.tryServePage(errorPage, w, nil)
	if err != nil {
		sh.logger.Error(fmt.Sprintf("Could not find error page. Error: '%s'.", err))
		sh.logger.Error(fmt.Sprintf("Attempted to find template page at: %s. Attempted to find error page at %s.", sh.template, sh.error))
		return
	}
	return
}

// tryServePage attempts to execute the template based on the page it is provided.
// Returns early with an error if there is any trouble loading the page.
func (sh *StaticHandler) tryServePage(page string, w http.ResponseWriter, data interface{}) error {

	tmpl, err := template.New("layout").ParseFiles(sh.template)
	if err != nil {
		sh.logger.Warn("Error parsing layout template.", "Error", err)
		return err
	}

	tmpl, err = tmpl.New("content").ParseFiles(page)
	if err != nil {
		sh.logger.Warn("Error parsing content template.", "Error", err)
		return err
	}

	if err := tmpl.ExecuteTemplate(w, "layout", data); err != nil {
		sh.logger.Warn("Error executing layout template.", "Error", err)
		return err
	}

	sh.logger.Info(fmt.Sprintf("Now serving page: %s", page))
	return nil
}

// nextPage adds a specified suffix to the current page and returns it.
// Also logs that the page could not be found and is attempting to find it at another location.
func (sh *StaticHandler) nextPage(page, nextLocation string) string {
	message := "Could not find page: " + page + ". "
	page = path.Join(page, nextLocation)
	message += "Checking: " + page

	sh.logger.Info(message)
	return page
}

// HandleRequests is the default request handler for static file handling.
// It takes the path from the url and checks sends it to the Render method.
func (sh *StaticHandler) HandleRequests(w http.ResponseWriter, r *http.Request) {
	url := filepath.Clean(r.URL.String())
	sh.Render(w, url, nil)
}

// ProductsDetails handles creating the individual page for each product.
// Will be sent requests from the URL '/products/{id}'
func (sh *StaticHandler) ProductsDetails(w http.ResponseWriter, r *http.Request, productsStore *data.ProductsStore) {
	id, err := helpers.GetIdFromRequest(r, "/products/")
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

	//TODO: Setup details page to work with htmx rather than loading 2 templates form go
}
