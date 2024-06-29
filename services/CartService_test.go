package services

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"ward4woods.ca/data"
	"ward4woods.ca/models"
)

func genMockItems() []models.CartItem {
	items := []models.CartItem{
		{
			Id:          1,
			Name:        "Item 1",
			Price:       25000,
			Description: "Test description 1",
			Category:    "Cutting Boards",
		},
		{
			Id:          2,
			Name:        "Item 2",
			Price:       17000,
			Description: "Test description 2",
			Category:    "Cutting Boards",
		},
		{
			Id:          3,
			Name:        "Item 3",
			Price:       21000,
			Description: "Test description 3",
			Category:    "Cutting Boards",
		},
	}
	return items
}

func genMockCart() models.Cart {
	items := genMockItems()
	return models.Cart{
		Id:    1,
		Items: items,
	}
}

type mockCartStore struct {
	shouldReturnError bool
}

var _ data.IProductCartStore = &mockCartStore{}

func (mcs *mockCartStore) AddProductToCart(models.CartId, models.CartItem) error {
	if mcs.shouldReturnError {
		return fmt.Errorf("mock error")
	}
	return nil
}

func (mcs *mockCartStore) ClearCart(models.CartId) error {
	if mcs.shouldReturnError {
		return fmt.Errorf("mock error")
	}
	return nil
}

func (mcs *mockCartStore) GetItemsFromCartId(models.CartId) ([]models.CartItem, error) {
	if mcs.shouldReturnError {
		return nil, fmt.Errorf("mock error")
	}
	return genMockItems(), nil
}

func (mcs *mockCartStore) RemoveProductFromCart(models.CartId, models.ProductId) error {
	if mcs.shouldReturnError {
		return fmt.Errorf("mock error")
	}
	return nil
}

func (mcs *mockCartStore) GetCartById(models.CartId) (models.Cart, error) {
	if mcs.shouldReturnError {
		return models.Cart{}, fmt.Errorf("mock error")
	}
	return genMockCart(), nil
}

func TestCartService(t *testing.T) {
	tests := []struct {
		name           string
		expectedResult interface{}
		expectedErr    error
		testFunc       func(ICartService) (interface{}, error)
	}{
		{
			name:           "GetCartById - success",
			expectedResult: genMockCart(),
			expectedErr:    nil,
			testFunc: func(cs ICartService) (interface{}, error) {
				return cs.GetCartById(1)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {})
		shouldReturnError := false
		if test.expectedErr != nil {
			shouldReturnError = true
		}
		mcs := mockCartStore{shouldReturnError: shouldReturnError}
		cs := NewCartService(&mcs)

		result, errResult := test.testFunc(cs)

		if !reflect.DeepEqual(result, test.expectedResult) {
			t.Errorf("Expected result: %+v. Got result %+v.", test.expectedResult, result)
		}

		if !errors.Is(errResult, test.expectedErr) {
			t.Errorf("Expected Error: %s. Got error: %s", test.expectedErr, errResult)
		}
	}
}
