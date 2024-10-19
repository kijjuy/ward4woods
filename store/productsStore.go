package store

import (
	"database/sql"
	"log/slog"
	"w4w/models"
)

var db *sql.DB

func SetupProductsStore(newDb *sql.DB) {
	db = newDb
}

func GetAllProducts() (models.Products, error) {
	rows, err := db.Query("SELECT * FROM products")

	if err != nil {
		slog.Error("Error when getting products from database", "Error", err)
		return nil, err
	}

	products := models.NewProducts()

	for rows.Next() {
		product := models.NewProduct()
		err := rows.Scan(&product.Id, &product.Name, &product.Price, &product.Description, &product.Category)
		if err != nil {
			slog.Error("Error when adding product to list", "Error", err)
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}

func GetProductById(id int) (models.Product, error) {
	row := db.QueryRow("SELECT * FROM products WHERE product_id = $1", id)

	product := models.Product{}

	err := row.Scan(&product.Id, &product.Name, &product.Price, &product.Description, &product.Category)

	return product, err
}

func DeleteProductById(id int) (int, error) {
	slog.Info("Entered DeleteProductById")
	result, err := db.Exec("DELETE FROM products WHERE product_id = $1", id)

	slog.Info("Got result from db")
	slog.Info("", "Result", result)

	rowsAffected, err := result.RowsAffected()

	slog.Info("Got rows affected")

	if err != nil {
		return 0, err
	}

	slog.Info("Returning rows affected")

	return int(rowsAffected), err
}

func CreateProduct(product models.Product) error {
	_, err := db.Exec("INSERT INTO products (name, price, description, category) VALUES($1, $2, $3, $4)", product.Name, product.Price, product.Description, product.Category)

	return err
}

func UpdateProduct(id int, product models.Product) (int, error) {
	result, err := db.Exec("UPDATE products SET name=$1, price=$2, description=$3, category=$4 WHERE product_id = $5", product.Name, product.Price, product.Description, product.Category, id)

	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()

	return int(rowsAffected), err
}

func GetCategories() ([]string, error) {
	rows, err := db.Query("SELECT DISTINCT category FROM products")

	if err != nil {
		return nil, err
	}

	categories := make([]string, 0)

	for rows.Next() {
		var category string

		err = rows.Scan(&category)

		if err != nil {
			return nil, err
		}

		categories = append(categories, category)
	}

	return categories, err
}
