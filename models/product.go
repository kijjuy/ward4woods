package models

type ProductId int

type Cents int

type Product struct {
	Id          ProductId `json:"id"`
	Name        string    `json:"name"`
	Price       Cents     `json:"price"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
}
