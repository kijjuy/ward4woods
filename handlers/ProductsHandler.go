package handlers

import (
	"fmt"
	"html/template"
	"log/slog"
	"net/http"

	"ward4woods.ca/application"
	"ward4woods.ca/helpers"
)

const (
	errorPath = "html/error.html"
)

type Result struct {
	data interface{}
	err  error
}

type ProductsHandler struct {
	routes
	logger *slog.Logger
}

type LogicFunc func(...interface{}) (interface{}, error)

type handlerAndPath struct {
	logicFunc LogicFunc
	path      string
}

type routes map[string]*handlerAndPath

func NewProductsHandler(logger *slog.Logger) *ProductsHandler {
	ph := &ProductsHandler{
		logger: logger.With("Location", "ProductsHandler"),
	}
	ph.logger.Info("Created product handler successfully.")
	return ph
}

// render renders the response from an html/template specified by templatePath using generic interface{} data.
func (ph *ProductsHandler) render(w http.ResponseWriter, templatePath string, data interface{}) {
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		ph.logger.Error("Could not load template files.", "Error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		ph.logger.Error("Could not execute template.", "Error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (ph *ProductsHandler) renderError(w http.ResponseWriter, err error) {
	ph.render(w, errorPath, err)
}

func (ph *ProductsHandler) AddHandler(route, htmlLocation string, logicFunc LogicFunc) {
	if ph.routes == nil {
		ph.routes = make(routes)
	}

	ph.routes[route] = &handlerAndPath{logicFunc: logicFunc, path: htmlLocation}
}

func (ph *ProductsHandler) Handle(w http.ResponseWriter, path string, data interface{}, err error) {
	if err != nil {
		ph.logger.Error("Error when trying to load response", "Error", err)
		ph.renderError(w, err)
		return
	}

	ph.logger.Info(fmt.Sprintf("Now generating response with data: [%+v]", data))
	ph.render(w, path, data)
}

func (ph *ProductsHandler) RegisterRoutes(router *application.Router) {
	for route, funcAndPath := range ph.routes {
		router.AddRoute(route, func(w http.ResponseWriter, r *http.Request) {
			result, err := funcAndPath.logicFunc("adsf")
			if err != nil {
				ph.renderError(w, err)
				return
			}
			ph.render(w, funcAndPath.path, result)
		})
	}
}

// TryWriteError checks if error is nil. If it is not, it writes the error response and returns true.
// If the error is nil, it doesn nothing then returns false.
func (ph *ProductsHandler) TryWriteError(w http.ResponseWriter, err error) bool {
	if err != nil {
		ph.renderError(w, err)
		return true
	}
	return false
}

func GetIdFromApiRequest(r *http.Request) (int, error) {
	prefix := "api/products"
	return helpers.GetIdFromRequest(r, prefix)
}
