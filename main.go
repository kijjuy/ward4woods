package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"log/slog"

	_ "github.com/lib/pq"
	"ward4woods.ca/data"
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

type ProductsHandler struct {
	productsStore *data.ProductsStore
	logger        *slog.Logger
}

func NewProductsHandler(productsStore *data.ProductsStore, logger *slog.Logger) *ProductsHandler {
	return &ProductsHandler{productsStore: productsStore, logger: logger}
}

func (ph *ProductsHandler) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	products, err := ph.productsStore.GetAllProducts()

	if err != nil {
		ph.logger.Error("Could not get products from store. Writing error response.")
		fmt.Fprintf(w, "Server error when trying to load products.")
		return
	}

	json.NewEncoder(w).Encode(products)
}

func (ph *ProductsHandler) GetProductById(w http.ResponseWriter, r *http.Request) {
	ph.logger.Info("Url: ", r.URL)
	idStr := strings.TrimPrefix(r.URL.String(), "/api/products/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ph.logger.Warn("Could not get products by id: invalid id. Error: ", err)
		fmt.Fprintf(w, "Invalid product id.")
	}

	product, err := ph.productsStore.GetProductById(id)

	if err != nil {
		ph.logger.Error("Could not get product from database.")
	}

	json.NewEncoder(w).Encode(product)
}

func main() {
	port := ":8080"
	logger := slog.Default()

	conString := "postgres://postgres@172.17.0.2:5432?password=Password@1&sslmode=disable"
	db := createDb(conString, logger)

	productsStore := createProductsStore(db, logger)

	productsHandler := NewProductsHandler(productsStore, logger)

	http.HandleFunc("/api/products", productsHandler.GetAllProducts)
	http.HandleFunc("/api/products/{id}", productsHandler.GetProductById)

	logger.Info(fmt.Sprintf("Application now lisening at: localhost%s", port))
	err := http.ListenAndServe(port, nil)

	logger.Error("Application crashed. Error:", err)
}
