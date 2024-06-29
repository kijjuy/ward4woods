package services

import (
	"fmt"
	"ward4woods.ca/data"
	"ward4woods.ca/models"
)

type ProductService struct {
	productsStore data.IProductsStore
}

type InvalidIdError struct{}

func (ide *InvalidIdError) Error() string {
	return "Invalid ID. ID must be a positive integer."
}

type ProductNotExistError struct {
	wantedId models.ProductId
}

func (pne *ProductNotExistError) Error() string {
	return fmt.Sprintf("Product with id %d does not exits.", pne.wantedId)
}

type EmptyProductError struct{}

func (epe *EmptyProductError) Error() string {
	return "Product was empty."
}

func NewProductService(productStore data.IProductsStore) *ProductService {
	return &ProductService{productsStore: productStore}
}

func (ps *ProductService) GetAllProducts() (*[]models.Product, error) {
	products, err := ps.productsStore.GetAllProducts()
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (ps *ProductService) GetProductById(id models.ProductId) (models.Product, error) {
	if id < 1 {
		return models.Product{}, &InvalidIdError{}
	}
	product, err := ps.productsStore.GetProductById(id)
	if product.Id == 0 {
		return product, &ProductNotExistError{id}
	}
	return product, err
}

func (ps *ProductService) AddProduct(product models.Product) error {
	if product.Id == 0 {
		return &EmptyProductError{}
	}
	return ps.productsStore.AddProduct(product)
}
