package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"log/slog"

	_ "github.com/lib/pq"
	"ward4woods.ca/data"
	"ward4woods.ca/handlers"
)

func createDb(conString string, logger *slog.Logger) *sql.DB {
	logger.Info("Connecting to database...")
	db, err := sql.Open("postgres", conString)
	if err != nil {
		logger.Error("Could not connect to database. Error:", err)
		os.Exit(1)
	}
	logger.Info("Database connection established.")
	return db
}

func createProductsStore(db *sql.DB, logger *slog.Logger) *data.ProductsStore {
	return data.NewProductsStore(db, logger)
}

func setupProduct(mux *http.ServeMux, productsStore *data.ProductsStore, logger *slog.Logger) {

	productsHandler := handlers.NewProductsHandler(productsStore, logger)

	mux.HandleFunc("/api/products", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			productsHandler.ProductsList(w, r)
			break
		case http.MethodPost:
			productsHandler.CreateProduct(w, r)
			break
		}
	})

	mux.HandleFunc("/api/products/{id}", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			productsHandler.GetProductById(w, r)
			break
		case http.MethodDelete:
			productsHandler.DeleteProductById(w, r)
			break
		}
	})

}

func setupStatic(mux *http.ServeMux, productsStore *data.ProductsStore, logger *slog.Logger) {
	htmlPath := "html"
	templatePath := filepath.Join(htmlPath, "_layout.html")
	errorPath := filepath.Join(htmlPath, "error.html")

	logger.Info(fmt.Sprintf("New static handler created with template path: '%s' and error path: '%s'", templatePath, errorPath))
	staticHandler := handlers.NewStaticHandler(htmlPath, templatePath, errorPath, logger)

	mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		faviconPath := filepath.Join(htmlPath, "favicon.ico")
		http.ServeFile(w, r, faviconPath)
	})

	mux.HandleFunc("/products/{id}", func(w http.ResponseWriter, r *http.Request) {
		staticHandler.ProductsDetails(w, r, productsStore)
	})

	mux.HandleFunc("/", staticHandler.HandleRequests)
}

func main() {
	port := ":8080"
	logger := slog.Default()

	conString := "postgres://postgres@172.17.0.2:5432?password=Password@1&sslmode=disable"
	db := createDb(conString, logger)

	productsStore := createProductsStore(db, logger)

	mux := http.NewServeMux()

	setupProduct(mux, productsStore, logger)

	setupStatic(mux, productsStore, logger)

	logger.Info(fmt.Sprintf("Application now lisening at: localhost%s", port))
	err := http.ListenAndServe(port, mux)

	logger.Error("Application crashed.", "Error", err)
}
