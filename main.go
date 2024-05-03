package main

import (
	"database/sql"
	"fmt"
	"os"

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

func main() {
	logger := slog.Default()

	conString := "postgres://postgres@172.17.0.2:5432?password=Password@1&sslmode=disable"
	db := createDb(conString, logger)

	productsStore := createProductsStore(db, logger)

	productsHandler := NewProductsHandler(productsStore, logger)

}
