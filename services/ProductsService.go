package services

import (
	"github.com/gorilla/sessions"
	"ward4woods.ca/data"
	"ward4woods.ca/models"
)

type ProductService struct {
	productsStore    *data.ProductsStore
	productCartStore *data.ProductCartStore
}

func NewProductService(productStore *data.ProductsStore, productsCartStore *data.ProductCartStore) *ProductService {
	return &ProductService{productsStore: productStore, productCartStore: productsCartStore}
}

func (ps *ProductService) GetAllProducts() ([]models.Product, error) {
	products, err := ps.productsStore.GetAllProducts()
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (ps *ProductService) GetProductById(id int) (models.Product, error) {
	return ps.productsStore.GetProductById(id)

}

func (ps *ProductService) AddToCart(productId int, session *sessions.Session) error {
	product, err := ps.productsStore.GetProductById(productId)
	if err != nil {
		return err
	}

	err = ps.productCartStore.AddProductToCart(session, &product)
	if err != nil {
		return err
	}

	return nil
}
