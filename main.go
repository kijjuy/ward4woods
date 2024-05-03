package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

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

func main() {
	port := ":8080"
	logger := slog.Default()

	conString := "postgres://postgres@172.17.0.2:5432?password=Password@1&sslmode=disable"
	db := createDb(conString, logger)

	productsStore := createProductsStore(db, logger)

	productsHandler := handlers.NewProductsHandler(productsStore, logger)

	mux := http.NewServeMux()

	mux.HandleFunc("/api/products", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			productsHandler.GetAllProducts(w, r)
		case http.MethodPost:
			productsHandler.CreateProduct(w, r)
		}
	})

	mux.HandleFunc("/api/products/{id}", productsHandler.GetProductById)

	logger.Info(fmt.Sprintf("Application now lisening at: localhost%s", port))
	err := http.ListenAndServe(port, mux)

	logger.Error("Application crashed. Error:", err)
}
