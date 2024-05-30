package data

import (
	"fmt"

	"github.com/gorilla/sessions"
	"ward4woods.ca/models"
)

type ProductCartStore struct {
	CartSessionId string
}

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

func (cs *ProductCartStore) AddProductToCart(session *sessions.Session, product *models.Product) error {
	cart, err := cs.getCart(session)
	if err != nil {
		return err
	}

	cart = append(cart, product)
	return nil
}

func (cs *ProductCartStore) RemoveProductFromCart(session *sessions.Session, productIndex int) error {
	cart, err := cs.getCart(session)
	if err != nil {
		return err
	}

	if productIndex >= len(cart) {
		return fmt.Errorf("Attempted to remove cart item outside of cart slice bounds.")
	}

	cart = append(cart[:productIndex], cart[productIndex+1:]...)
	return nil
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
