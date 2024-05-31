package handlers

import (
	"github.com/gorilla/sessions"
	"net/http"
	"ward4woods.ca/application"
	"ward4woods.ca/data"
	"ward4woods.ca/services"
)

type productsHandler struct {
	router         *application.Router
	apiHandler     *application.ApiHandler
	productsStore  *data.ProductsStore
	sessionStore   *sessions.CookieStore
	productService *services.ProductService
	logger         *application.Logger
}

func HandleProducts(router *application.Router, productsStore *data.ProductsStore, logger *application.Logger, sessionStore *sessions.CookieStore) {

	apiHandler := application.NewApiHandler(logger)
	productsCartStore := data.NewProductsCartStore(data.CartSessionName)
	productService := services.NewProductService(productsStore, productsCartStore)

	ph := &productsHandler{
		router:         router,
		apiHandler:     apiHandler,
		productsStore:  productsStore,
		sessionStore:   sessionStore,
		productService: productService,
		logger:         logger,
	}

	router.AddRoute("/api/products", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			ph.getAllProducts(w)
			break
		}
	})

	router.AddRoute("/api/products/{id}", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			ph.productDetails(w, r)
			break
		}
	})

	router.AddRoute("/api/addToCart/{id}", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			ph.addToCart(w, r)
			break
		}
	})

}

func (ph *productsHandler) getAllProducts(w http.ResponseWriter) {
	templateWriter := application.NewApiTemplateWriter("html/templates/productsList.html")
	result, err := ph.productService.GetAllProducts()
	if ph.apiHandler.TryWriteError(application.WriteServerError, w, err) {
		return
	}

	ph.apiHandler.Handle(templateWriter.WriteTemplate, w, result)
}

func (ph *productsHandler) productDetails(w http.ResponseWriter, r *http.Request) {

	templateWriter := application.NewApiTemplateWriter("html/templates/productsDetails.html")
	id, err := application.GetIdFromApiRequest(r)
	if ph.apiHandler.TryWriteError(application.WriteServerError, w, err) {
		return
	}
	result, err := ph.productService.GetProductById(id)
	ph.apiHandler.Handle(templateWriter.WriteTemplate, w, result)
}

func (ph *productsHandler) addToCart(w http.ResponseWriter, r *http.Request) {

	id, err := application.GetIdFromApiRequest(r)
	if ph.apiHandler.TryWriteError(application.WriteServerError, w, err) {
		return
	}

	session, err := ph.sessionStore.Get(r, data.CartSessionName)
	if ph.apiHandler.TryWriteError(application.WriteServerError, w, err) {
		return
	}

	err = ph.productService.AddToCart(id, session)
}
