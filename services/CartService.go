package services

import (
	"ward4woods.ca/data"
	"ward4woods.ca/models"
)

type ICartService interface {
	GetCartById(models.CartId) (models.Cart, error)
	AddItemToCart(models.Cart, models.CartItem) error
	RemoveItemFromCart(models.Cart, models.ProductId) error
	ClearCart(models.Cart) error
}

type CartService struct {
	CartStore data.IProductCartStore
}

var _ ICartService = &CartService{}

type CartDbError struct {
	Err error
}

func (e *CartDbError) Unwrap() error { return e.Err }

func (e *CartDbError) Error() string {
	return "Could not get cart from db."
}

func NewCartService(cartStore data.IProductCartStore) ICartService {
	cartService := CartService{cartStore}
	return &cartService
}

func (cs *CartService) GetCartById(cartId models.CartId) (models.Cart, error) {
	cart, err := cs.CartStore.GetCartById(cartId)

	if err != nil {
		return models.Cart{}, &CartDbError{err}
	}

	return cart, nil
}

func (cs *CartService) AddItemToCart(cart models.Cart, item models.CartItem) error {
	err := cs.CartStore.AddProductToCart(cart.Id, item)

	if err != nil {
		return &CartDbError{err}
	}

	return nil
}

func (cs *CartService) ClearCart(cart models.Cart) error {
	err := cs.CartStore.ClearCart(cart.Id)
	if err != nil {
		return &CartDbError{err}
	}
	return nil
}

func (cs *CartService) RemoveItemFromCart(cart models.Cart, prodId models.ProductId) error {
	err := cs.CartStore.RemoveProductFromCart(cart.Id, prodId)

	if err != nil {
		return &CartDbError{err}
	}

	return nil
}
