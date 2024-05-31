package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"log/slog"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"ward4woods.ca/application"
	"ward4woods.ca/data"
	"ward4woods.ca/handlers"
)

func loadDotEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file.")
	}
}

func createDb(logger *application.Logger) *sql.DB {
	conString := os.Getenv("DATABASE_URL")
	logger.Info("Connecting to database...", nil)
	db, err := sql.Open("postgres", conString)
	if err != nil {
		logger.Error("Could not connect to database.", err)
		os.Exit(1)
	}
	logger.Info("Database connection established.", nil)
	return db
}

func createSessionStore(logger *application.Logger) *sessions.CookieStore {
	sessionKey := []byte(os.Getenv("SESSION_KEY"))

	sessionStore := sessions.NewCookieStore([]byte(sessionKey))
	logger.Info("Session store created successfully.", nil)
	return sessionStore
}

func createProductsStore(db *sql.DB, logger *application.Logger) *data.ProductsStore {
	return data.NewProductsStore(db, logger)
}

func setupProduct(router *application.Router, productsStore *data.ProductsStore, logger *slog.Logger, sessionStore *sessions.CookieStore) {

	productsHandler := handlers.NewProductsHandler(logger)
	productsCartStore := data.NewProductsCartStore(cartSessionName)
	productService := services.NewProductService(productsStore, productsCartStore)

	router.AddRoute("/api/products", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			result, err := productService.GetAllProducts()
			productsHandler.Handle(w, "html/templates/productsList.html", result, err)
			break
		}
	})

	router.AddRoute("/api/products/{id}", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			id, err := handlers.GetIdFromApiRequest(r)
			if productsHandler.TryWriteError(w, err) {
				break
			}
			result, err := productService.GetProductById(id)
			productsHandler.Handle(w, "html/templates/productsDetails.html", result, err)
			break
		}
	})

	router.AddRoute("/api/addToCart/{id}", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			id, err := handlers.GetIdFromApiRequest(r)
			if productsHandler.TryWriteError(w, err) {
				break
			}

			session, err := sessionStore.Get(r, productsCartStore.CartSessionId)
			if productsHandler.TryWriteError(w, err) {
				break
			}

			err = productService.AddToCart(id, session)
		}
	})

}

func setupStatic(router *application.Router, productsStore *data.ProductsStore, logger *application.Logger) {
	htmlPath := "html"
	templatePath := filepath.Join(htmlPath, "_layout.html")
	errorPath := filepath.Join(htmlPath, "error.html")

	logger.Info(fmt.Sprintf("New static handler created with template path: '%s' and error path: '%s'", templatePath, errorPath), nil)
	staticHandler := handlers.NewStaticHandler(htmlPath, templatePath, errorPath, logger)

	router.AddRoute("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		faviconPath := filepath.Join(htmlPath, "favicon.ico")
		http.ServeFile(w, r, faviconPath)
	})

	router.AddRoute("/products/{id}", func(w http.ResponseWriter, r *http.Request) {
		staticHandler.ProductsDetails(w, r, productsStore)
	})

	router.AddRoute("/", staticHandler.HandleRequests)
}

func main() {
	loadDotEnv()
	port := ":8080"
	logger := application.NewLogger(slog.Default())

	db := createDb(logger)
	sessionStore := createSessionStore(logger)

	productsStore := createProductsStore(db, logger)

	router := application.NewRouter()

	handlers.HandleProducts(router, productsStore, logger, sessionStore)

	setupStatic(router, productsStore, logger)

	logger.Info(fmt.Sprintf("Application now lisening at: localhost%s", port), nil)
	err := http.ListenAndServe(port, router.Serve())

	logger.Error("Application crashed.", err)
}
