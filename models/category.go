package models

type Category struct {
	Category string
	Selected string
}

type Categories []Category

func NewCategories() Categories {
	return make([]Category, 0)
}
