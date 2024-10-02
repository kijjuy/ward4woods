package services

import (
	"w4w/models"
	"w4w/store"
)

type ErrNoRowsAffected struct{}

func (e *ErrNoRowsAffected) Error() string {
	return "No database entries were affected"
}

func GetAllProducts() (models.Products, error) {
	return store.GetAllProducts()
}

func GetProductById(id int) (models.Product, error) {
	return store.GetProductById(id)
}

func DeleteProduct(id int) error {
	numDeleted, err := store.DeleteProductById(id)

	if err != nil {
		return err
	}

	if numDeleted == 0 {
		return &ErrNoRowsAffected{}
	}

	return err
}

func CreateProduct(product models.Product) error {
	return store.CreateProduct(product)
}

func UpdateProduct(id int, product models.Product) error {
	rowsAffected, err := store.UpdateProduct(id, product)

	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return &ErrNoRowsAffected{}
	}

	return nil
}

func GetCategories(currentCategory string) (models.Categories, error) {
	categoriesStr, err := store.GetCategories()

	if err != nil {
		return nil, err
	}

	categories := models.NewCategories()

	for _, categoryStr := range categoriesStr {
		isCurrent := currentCategory == categoryStr

		var selected string

		if isCurrent {
			selected = "selected"
		}

		category := models.Category{
			Category: categoryStr,
			Selected: selected,
		}

		categories = append(categories, category)
	}

	return categories, err
}
