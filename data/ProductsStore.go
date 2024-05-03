package data

import (
	"database/sql"
	"fmt"
	"log/slog"

	"ward4woods.ca/models"
)

type IProductsStore interface {
	GetAllProducts() ([]models.Product, error)
	GetProductById(int) (models.Product, error)
	AddProduct(models.Product) error
	DeleteProductById(int) error
	UpdateProduct(models.Product) error
}

type ProductsStore struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewProductsStore(db *sql.DB, logger *slog.Logger) *ProductsStore {
	return &ProductsStore{db: db, logger: logger}
}

func (ps *ProductsStore) GetAllProducts() ([]models.Product, error) {
	rows, err := ps.db.Query("SELECT * FROM products")
	if err != nil {
		ps.logger.Error("Could not get all products from database. err: ", err)
	}
	defer rows.Close()

	products, err := ps.scanProducts(rows)

	if err != nil {
		ps.logger.Error("Could not get records from products table. Error:", err)
	}

	ps.logger.Info(fmt.Sprintf("Got %d records from products table.", len(products)))
	return products, nil
}

func (ps *ProductsStore) GetProductById(id int) (models.Product, error) {
	row := ps.db.QueryRow("SELECT * FROM products WHERE product_id = $1", id)
	product, err := scanProduct(row)

	if err == sql.ErrNoRows {
		ps.logger.Info("")
	}

	return product, err
}

func (ps *ProductsStore) AddProduct(product models.Product) error {
	panic("not implemented")
}

func (ps *ProductsStore) DeleteProductById(id int) error {
	panic("not implemented")
}

func (ps *ProductsStore) UpdateProduct(product models.Product) error {
	panic("not implemented")
}

func scanProduct(row *sql.Row) (models.Product, error) {
	product := models.Product{}
	err := row.Scan(&product.Id, &product.Name, &product.Price, &product.Description, &product.Category)

	return product, err
}

func (ps *ProductsStore) scanProducts(rows *sql.Rows) ([]models.Product, error) {
	var products []models.Product
	for rows.Next() {
		product := models.Product{}
		err := rows.Scan(&product.Id, &product.Name, &product.Price, &product.Description, &product.Category)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
}
