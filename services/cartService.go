package services

import (
	"database/sql"
	"w4w/models"
	"w4w/store"
)

type ErrDupeCartItem struct {
}

func (e *ErrDupeCartItem) Error() string {
	return "Cannot add duplicate items to cart"
}

type ErrCartUnauthorized struct{}

func (e *ErrCartUnauthorized) Error() string {
	return "User attempted to modify contents of cart not belonging to them"
}

func AddToCart(sessId string, productId int) error {
	err := ensureCartCreated(sessId)
	if err != nil {
		return err
	}
	products, err := store.GetCartItemsBySessId(sessId)

	if err != nil {
		return err
	}

	for _, product := range products {
		if productId == product.Id {
			return &ErrDupeCartItem{}
		}
	}

	return store.AddToCart(sessId, productId)
}

func ensureCartCreated(sessId string) error {
	_, err := store.GetCartIdBySessId(sessId)
	if err == sql.ErrNoRows {
		return store.CreateCart(sessId)
	}
	return err
}

func GetCartItems(sessId string) (models.CartDisplayProducts, error) {
	return store.GetCartDisplayModelBySessId(sessId)
}

func DeleteFromCart(sessId string, cartItemId int) error {
	isOwner, err := store.CartItemOwenershipIsValid(sessId, cartItemId)

	if err != nil {
		return err
	}

	if !isOwner {
		return &ErrCartUnauthorized{}
	}

	return store.DeleteFromCart(sessId, cartItemId)
}

func ClearCart(sessId string) error {
	return store.DeleteAllCartItems(sessId)
}
