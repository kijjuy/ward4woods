package models

type Product struct {
	Id          int
	Name        string
	Price       int
	Description string
	Category    string
}

func NewProduct() Product {
	return Product{}
}

type Products []Product

func NewProducts() Products {
	return make([]Product, 0)
}

type CartDisplayProduct struct {
	Product    Product
	CartItemId int
}

type CartDisplayProducts []CartDisplayProduct

func NewCartDisplayProducts() CartDisplayProducts {
	return make([]CartDisplayProduct, 0)
}
