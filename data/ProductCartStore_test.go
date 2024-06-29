package data

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"ward4woods.ca/models"
)

func generateTestProduct(id models.ProductId, price models.Cents) models.CartItem {
	return models.CartItem{
		Id:          id,
		Name:        fmt.Sprintf("Product %d", id),
		Price:       price,
		Description: fmt.Sprintf("Test product %d", id),
		Category:    "cutting board",
	}
}

func generateTestCartNoItems() models.Cart {
	return models.Cart{
		Id: 1,
	}
}

func generateTestCartSingleItem() models.Cart {
	item := generateTestProduct(1, 25000)
	items := []models.CartItem{
		item,
	}

	return models.Cart{
		Id:    1,
		Items: items,
	}
}

func generateTestCartThreeItems() models.Cart {
	item1 := generateTestProduct(1, 25000)
	item2 := generateTestProduct(2, 17000)
	item3 := generateTestProduct(3, 50000)
	items := []models.CartItem{item1, item2, item3}
	return models.Cart{
		Id:    1,
		Items: items,
	}
}

func TestCartStore(t *testing.T) {
	tests := []struct {
		name           string
		expectedResult interface{}
		expectedError  error
		testFunc       func(IProductCartStore) (any, error)
	}{
		{
			name:           "AddProductToCart - Success",
			expectedResult: nil,
			expectedError:  nil,
			testFunc: func(cs IProductCartStore) (any, error) {
				cartItem := generateTestProduct(1, 25000)
				return nil, cs.AddProductToCart(1, cartItem)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cartStore := NewProductsCartStore("abcd")

			result, err := test.testFunc(cartStore)

			if !reflect.DeepEqual(result, test.expectedResult) {
				t.Errorf("Expected result: %+v. Got: %+v.", test.expectedResult, result)
			}

			if !errors.Is(err, test.expectedError) {
				t.Errorf("Expected error: %+v. Got: %+v.", test.expectedError, err)
			}
		})

	}

}
