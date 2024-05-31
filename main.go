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

func createDb(logger *slog.Logger) *sql.DB {
	conString := os.Getenv("DATABASE_URL")
	logger.Info("Connecting to database...")
	db, err := sql.Open("postgres", conString)
	if err != nil {
		logger.Error("Could not connect to database. Error:", err)
		os.Exit(1)
	}
	logger.Info("Database connection established.")
	return db
}

func createSessionStore() *sessions.CookieStore {
	sessionKey := []byte(os.Getenv("SESSION_KEY"))

	sessionStore := sessions.NewCookieStore([]byte(sessionKey))
	slog.Info("Session store created successfully.")
	return sessionStore
}

func createProductsStore(db *sql.DB, logger *slog.Logger) *data.ProductsStore {
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

func setupStatic(router *application.Router, productsStore *data.ProductsStore, logger *slog.Logger) {
	htmlPath := "html"
	templatePath := filepath.Join(htmlPath, "_layout.html")
	errorPath := filepath.Join(htmlPath, "error.html")

	logger.Info(fmt.Sprintf("New static handler created with template path: '%s' and error path: '%s'", templatePath, errorPath))
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
	logger := slog.Default()

	db := createDb(logger)
	sessionStore := createSessionStore()

	productsStore := createProductsStore(db, logger)

	router := application.NewRouter()

	setupProduct(router, productsStore, logger, sessionStore)

	setupStatic(router, productsStore, logger)

	logger.Info(fmt.Sprintf("Application now lisening at: localhost%s", port))
	err := http.ListenAndServe(port, router.Serve())

	logger.Error("Application crashed.", "Error", err)
}
