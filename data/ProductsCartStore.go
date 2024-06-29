package data

import (
	"fmt"

	"github.com/gorilla/sessions"
	"ward4woods.ca/models"
)

var CartSessionName = "castSession"

type IProductCartStore interface {
	AddProductToCart(models.CartId, models.CartItem) error
	RemoveProductFromCart(models.CartId, models.ProductId) error
	GetItemsFromCartId(models.CartId) ([]models.CartItem, error)
	ClearCart(models.CartId) error
	GetCartById(models.CartId) (models.Cart, error)
}

type ProductCartStore struct {
	CartSessionId string
}

var _ IProductCartStore = &ProductCartStore{}

type productsCart []*models.Product

func NewProductsCartStore(CartSessionId string) *ProductCartStore {
	return &ProductCartStore{CartSessionId: CartSessionId}
}

// tryCreateCart creates an instance of the shopping cart inside of the session store if it does not exist.
// If it already exists, nothing happens.
func (cs *ProductCartStore) tryCreateCart(session *sessions.Session) {
	if session.Values[cs.CartSessionId] == nil {
		session.Values[cs.CartSessionId] = make(productsCart, 0)
	}
}

func (cs *ProductCartStore) AddProductToCart(models.CartId, models.CartItem) error {
	return fmt.Errorf("Not implemented")
}

func (cs *ProductCartStore) RemoveProductFromCart(cartId models.CartId, itemId models.ProductId) error {
	return fmt.Errorf("Not implemented")
}

func (cs *ProductCartStore) GetItemsFromCartId(id models.CartId) ([]models.CartItem, error) {
	//TODO: Write select for cart
	emptyItem := make([]models.CartItem, 0)
	err := fmt.Errorf("not implemented")
	return emptyItem, err

}

func (cs *ProductCartStore) GetCartById(id models.CartId) (models.Cart, error) {
	emptyCart := models.Cart{}
	err := fmt.Errorf("not implemented")
	return emptyCart, err
}

func (cs *ProductCartStore) ClearCart(id models.CartId) error {
	return fmt.Errorf("Not implemented")
}

func (cs *ProductCartStore) getCart(session *sessions.Session) (productsCart, error) {
	cs.tryCreateCart(session)
	cart, ok := session.Values[cs.CartSessionId].(productsCart)

	if !ok {
		error := fmt.Errorf("Error when trying to get cart: Could not convert " +
			"session.Values[CartSessionId] to type productsCart.")
		return nil, error

	}

	return cart, nil
}
