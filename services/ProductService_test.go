package services

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"ward4woods.ca/data"
	"ward4woods.ca/models"
)

var products = []models.Product{
	{
		Id:          1,
		Name:        "Product1",
		Price:       15000,
		Description: "First product",
		Category:    "Cutting Board",
	}, {
		Id:          2,
		Name:        "Product2",
		Price:       17000,
		Description: "Second product",
		Category:    "Charcuterie Board",
	}, {
		Id:          3,
		Name:        "Product3",
		Price:       25000,
		Description: "Third product",
		Category:    "Chopping Block",
	},
}

type MockProductStore struct {
	shouldReturnError bool
}

var _ data.IProductsStore = &MockProductStore{}

func (mps *MockProductStore) AddProduct(models.Product) error {
	if mps.shouldReturnError {
		return fmt.Errorf("mock error")
	}
	return nil
}

func (mps *MockProductStore) DeleteProductById(id models.ProductId) error {
	if mps.shouldReturnError {
		return fmt.Errorf("mock error")
	}
	return nil
}

func (mps *MockProductStore) GetAllProducts() (*[]models.Product, error) {
	if mps.shouldReturnError {
		return nil, fmt.Errorf("mock error")
	}
	return &products, nil
}

func (mps *MockProductStore) GetProductById(id models.ProductId) (models.Product, error) {
	if mps.shouldReturnError {
		return models.Product{}, fmt.Errorf("mock error")
	}
	var product models.Product
	for _, p := range products {
		if p.Id == id {
			product = p
		}
	}
	return product, nil
}

func (mps *MockProductStore) UpdateProduct(product models.Product) error {
	if mps.shouldReturnError {
		return fmt.Errorf("mock error")
	}
	return nil
}

func TestProductService(t *testing.T) {
	tests := []struct {
		name           string
		expectedResult interface{}
		expectedErr    error
		testFunc       func(*ProductService) (interface{}, error)
	}{
		{
			name:           "GetAllProducts - Success",
			expectedResult: &products,
			expectedErr:    nil,
			testFunc: func(ps *ProductService) (interface{}, error) {
				return ps.GetAllProducts()
			},
		},
		{
			name:           "GetProductById - Success with valid product",
			expectedResult: products[0],
			expectedErr:    nil,
			testFunc: func(ps *ProductService) (interface{}, error) {
				return ps.GetProductById(1)
			},
		},
		{
			name:           "GetProductById - Error nil product",
			expectedResult: models.Product{},
			expectedErr:    &ProductNotExistError{4},
			testFunc: func(ps *ProductService) (interface{}, error) {
				return ps.GetProductById(4)
			},
		},
		{
			name:           "GetProductById - Error invalid id",
			expectedResult: models.Product{},
			expectedErr:    &InvalidIdError{},
			testFunc: func(ps *ProductService) (interface{}, error) {
				return ps.GetProductById(0)
			},
		},
		{
			name:           "AddProduct - Success",
			expectedResult: nil,
			expectedErr:    nil,
			testFunc: func(ps *ProductService) (interface{}, error) {
				return nil, ps.AddProduct(models.Product{
					Id:          4,
					Name:        "Product4",
					Price:       25000,
					Description: "test desc",
					Category:    "test cat",
				})
			},
		},
		{
			name:           "AddProduct - Error empty product",
			expectedResult: nil,
			expectedErr:    &EmptyProductError{},
			testFunc: func(ps *ProductService) (interface{}, error) {
				return nil, ps.AddProduct(models.Product{})
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			shouldReturnError := false
			if test.expectedErr != nil {
				shouldReturnError = true
			}
			mps := &MockProductStore{shouldReturnError: shouldReturnError}
			ps := NewProductService(mps)

			result, errResult := test.testFunc(ps)

			if !errors.Is(errResult, test.expectedErr) {
				t.Errorf("Expected error: [%s]. Got: [%s]", test.expectedErr, errResult)
			}

			if !reflect.DeepEqual(result, test.expectedResult) {
				t.Errorf("Expected result: [%v]. Got: [%v]", test.expectedResult, result)
			}

		})
	}

}
